package telemetry

import (
	"log/slog"
	"math/rand"
	"net/http"
	"time"
)

// Backend represents a dependency that mission-control monitors.
type Backend struct {
	Address   string // e.g. "loki:3100"
	HealthURL string // full URL to hit for health checks
}

// BackendMonitor periodically health-checks backend services
// and records the results as Prometheus metrics.
type BackendMonitor struct {
	metrics  *Metrics
	backends []Backend
	client   *http.Client
	ticker   *time.Ticker
	done     chan struct{}
}

func NewBackendMonitor(metrics *Metrics) *BackendMonitor {
	return &BackendMonitor{
		metrics: metrics,
		backends: []Backend{
			{Address: "grafana:3000", HealthURL: "http://grafana:3000/api/health"},
			{Address: "loki:3100", HealthURL: "http://loki:3100/ready"},
			{Address: "tempo:3200", HealthURL: "http://tempo:3200/ready"},
			{Address: "mimir:9009", HealthURL: "http://mimir:9009/ready"},
		},
		client: &http.Client{Timeout: 5 * time.Second},
		done:   make(chan struct{}),
	}
}

func (bm *BackendMonitor) Start() {
	slog.Info("starting backend monitor", "component", "backends", "targets", len(bm.backends))

	// Run an initial check immediately
	bm.checkAll()

	bm.ticker = time.NewTicker(10 * time.Second)
	go bm.run()
}

func (bm *BackendMonitor) Stop() {
	if bm.ticker != nil {
		bm.ticker.Stop()
	}
	close(bm.done)
	slog.Info("backend monitor stopped", "component", "backends")
}

func (bm *BackendMonitor) run() {
	for {
		select {
		case <-bm.ticker.C:
			bm.checkAll()
			// Jitter: 8-12s
			bm.ticker.Reset(time.Duration(8+rand.Intn(5)) * time.Second)
		case <-bm.done:
			return
		}
	}
}

func (bm *BackendMonitor) checkAll() {
	for _, b := range bm.backends {
		bm.metrics.BackendHealthChecks.WithLabelValues(b.Address).Inc()

		resp, err := bm.client.Get(b.HealthURL)
		if err != nil || resp.StatusCode >= 400 {
			bm.metrics.BackendUp.WithLabelValues(b.Address).Set(0)
			slog.Debug("backend health check failed",
				"component", "backends",
				"backend", b.Address,
				"error", err)
		} else {
			bm.metrics.BackendUp.WithLabelValues(b.Address).Set(1)
		}

		if resp != nil {
			resp.Body.Close()
		}
	}
}
