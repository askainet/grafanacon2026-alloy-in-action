package handlers

import (
	"encoding/json"
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/trace"

	"github.com/grafana/alloy-mission-control/internal/missioncontrol"
	"github.com/grafana/alloy-mission-control/internal/telemetry"
)

type IntelHandler struct {
	metrics    *telemetry.Metrics
	tracer     trace.Tracer
	controller *missioncontrol.Controller
}

func NewIntelHandler(metrics *telemetry.Metrics, tracer trace.Tracer, controller *missioncontrol.Controller) *IntelHandler {
	return &IntelHandler{
		metrics:    metrics,
		tracer:     tracer,
		controller: controller,
	}
}

func (h *IntelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Start root span
	ctx, span := h.tracer.Start(ctx, "process_intel")
	defer span.End()

	// Generate intel report ID
	reportID := uuid.New().String()

	// Simulate intel report processing with varied duration
	duration := randomDuration()
	time.Sleep(duration)

	// Record metrics
	h.metrics.IntelReportsProcessed.With(prometheus.Labels{
		"endpoint": "process_intel",
	}).Inc()

	// Log intel processing completion
	slog.Info("intel report processed",
		"component", "intel",
		"report_id", reportID,
		"duration_ms", duration.Milliseconds(),
	)

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":    "success",
		"report_id": reportID,
	})
}

// randomDuration returns a random duration between 5-50ms
func randomDuration() time.Duration {
	min := 5
	max := 50
	ms := rand.Intn(max-min+1) + min
	return time.Duration(ms) * time.Millisecond
}
