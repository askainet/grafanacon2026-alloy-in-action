package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/grafana/alloy-mission-control/internal/config"
	"github.com/grafana/alloy-mission-control/internal/handlers"
	"github.com/grafana/alloy-mission-control/internal/middleware"
	"github.com/grafana/alloy-mission-control/internal/missioncontrol"
	"github.com/grafana/alloy-mission-control/internal/missions"
	"github.com/grafana/alloy-mission-control/internal/telemetry"
	"github.com/grafana/alloy-mission-control/internal/traffic"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	// Initialize logging
	telemetry.InitLogging(cfg.LogLevel)

	slog.Info("mission control starting", "component", "main", "service", cfg.ServiceName)

	// Initialize tracing
	tp, err := telemetry.InitTracing(cfg)
	if err != nil {
		log.Fatal("failed to init tracing:", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			slog.Error("failed to shutdown tracer", "component", "main", "error", err)
		}
	}()

	// Initialize metrics
	metrics := telemetry.InitMetrics()

	baseURL := fmt.Sprintf("http://localhost:%s", cfg.Port)

	// Create mission controller
	controller := missioncontrol.NewController(metrics)

	// Create missions
	mission4 := missions.NewMission4(cfg, baseURL)

	// Register all missions
	registry := controller.Registry()
	registry.Register(missions.NewMission1(cfg, baseURL, metrics))
	registry.Register(missions.NewMission2())
	registry.Register(missions.NewMission3())
	registry.Register(mission4)

	// Create tracer for handlers
	tracer := tp.Tracer("mission-control")

	// Setup router with middleware
	r := chi.NewRouter()
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.MetricsMiddleware(metrics))

	// Health endpoint
	r.Get("/health", handlers.NewHealthHandler().ServeHTTP)

	// Metrics endpoint
	r.Handle("/metrics", promhttp.Handler())

	// Intel processing endpoint
	intelHandler := handlers.NewIntelHandler(metrics, tracer, controller)
	r.Post("/api/process-intel", intelHandler.ServeHTTP)

	// Operations endpoints (flavor/dummy endpoints for realistic traffic)
	opsHandler := handlers.NewOperationsHandler(metrics, tracer, controller, cfg, mission4)
	r.Post("/api/biometrics/retina-scan/verify", opsHandler.RetinaScanVerify)
	r.Post("/api/extraction/heli-pickup/request", opsHandler.HeliPickupRequest)
	r.Get("/api/message/burn-after-reading", opsHandler.BurnAfterReading)
	r.Post("/api/comms/dead-drop/upload", opsHandler.DeadDropUpload)
	r.Post("/api/sat-link/handshake/establish", opsHandler.SatLinkHandshake)
	r.Get("/api/comms/dead-drop/view", opsHandler.DeadDropView)

	// Agent endpoints
	agentManager := telemetry.NewAgentManager(metrics, baseURL)
	agentHandler := handlers.NewAgentHandler(agentManager, metrics)
	r.Route("/api/agents", func(r chi.Router) {
		r.Post("/heartbeat", agentHandler.Heartbeat)
		r.Get("/", agentHandler.ListLegitimateAgents)
		r.Get("/regex", agentHandler.GetLegitimateAgentRegex)
	})

	// Known legitimate API paths (used by Mission 1 path filtering)
	validPaths := []string{
		"/health",
		"/metrics",
		"/api/process-intel",
		"/api/biometrics/retina-scan/verify",
		"/api/extraction/heli-pickup/request",
		"/api/message/burn-after-reading",
		"/api/comms/dead-drop/upload",
		"/api/comms/dead-drop/view",
		"/api/sat-link/handshake/establish",
		"/api/agents/heartbeat",
		"/api/agents",
		"/api/agents/regex",
		"/api/metrics/allowed-paths",
	}

	// Allowed-paths endpoint for Mission 1
	pathsHandler := handlers.NewPathsHandler(validPaths)
	r.Get("/api/metrics/allowed-paths", pathsHandler.GetAllowedPathsRegex)

	// Admin endpoints for mission control
	adminHandler := handlers.NewAdminHandler(controller)
	s3Handler := handlers.NewS3Handler(cfg.S3Endpoint, "audit-logs", cfg.LokiEndpoint)
	tempoHandler := handlers.NewTempoHandler(cfg.TempoEndpoint, mission4.NumChunks())

	mimirHandler := handlers.NewMimirHandler(cfg.MimirEndpoint, validPaths, func() bool {
		return controller.IsActive("mission1")
	})
	mission3Handler := handlers.NewMission3Handler(cfg.AlloyEndpoint, func() bool {
		return controller.IsActive("mission3")
	})
	mission4Handler := handlers.NewMission4Handler(cfg.AlloyEndpoint, cfg.TempoEndpoint, mission4.NumChunks(), func() bool {
		return controller.IsActive("mission4")
	})

	r.Route("/admin", func(r chi.Router) {
		r.Get("/status", adminHandler.Status)
		r.Post("/missions/{id}/start", adminHandler.StartMission)
		r.Post("/reset", adminHandler.Reset)
		r.Get("/mission1/verify", mimirHandler.Verify)
		r.Get("/mission3/verify", mission3Handler.Verify)
		r.Get("/s3/audit-logs", s3Handler.ListObjects)
		r.Get("/s3/audit-logs/verify", s3Handler.Verify)
		r.Get("/mission4/verify", mission4Handler.Verify)
		r.Get("/mission4/access-token", tempoHandler.ServeHTTP)
	})

	// Start server in a goroutine so heartbeat-based services can connect
	addr := ":" + cfg.Port
	slog.Info("server listening", "component", "main", "addr", addr)
	go func() {
		if err := http.ListenAndServe(addr, r); err != nil {
			log.Fatal("server error:", err)
		}
	}()

	// Start agent manager (sends heartbeats to /api/agents/heartbeat)
	agentManager.Start(context.Background())
	defer agentManager.Stop()

	// Start synthetic traffic generator (3 workers creating baseline telemetry)
	trafficGen := traffic.NewGenerator(baseURL, 3)
	trafficGen.Start(context.Background())
	defer trafficGen.Stop()

	// Block forever (server runs in goroutine above)
	select {}
}
