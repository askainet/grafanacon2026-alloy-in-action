package missions

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/grafana/alloy-mission-control/internal/config"
	"github.com/grafana/alloy-mission-control/internal/telemetry"
)

// Mission1 implements the High Cardinality Explosion scenario.
// An adversary probes the server with random URL paths, creating thousands of
// unique path label values in http_requests_total and http_request_duration_seconds.
type Mission1 struct {
	baseCardinality int
	maxCardinality  int
	growthRate      int
	growthInterval  time.Duration

	baseURL string
	client  *http.Client
	metrics *telemetry.Metrics

	currentCardinality atomic.Int32
	garbagePaths       []string

	ticker *time.Ticker
	done   chan struct{}
}

func NewMission1(cfg *config.Config, baseURL string, metrics *telemetry.Metrics) *Mission1 {
	m := &Mission1{
		baseCardinality: 100,
		maxCardinality:  cfg.Mission1MaxCardinality,
		growthRate:      cfg.Mission1GrowthRate,
		growthInterval:  10 * time.Second,
		baseURL:         baseURL,
		client:          &http.Client{Timeout: 10 * time.Second},
		metrics:         metrics,
		garbagePaths:    make([]string, 0, cfg.Mission1MaxCardinality),
	}

	// Pre-generate garbage paths that look like adversary endpoint probing
	for i := 0; i < m.maxCardinality; i++ {
		m.garbagePaths = append(m.garbagePaths, randomPath())
	}

	return m
}

func randomPath() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("/api/%x", b)
}

func (m *Mission1) ID() string   { return "mission1" }
func (m *Mission1) Name() string { return "High Cardinality Explosion" }
func (m *Mission1) Description() string {
	return "Adversary probes random URL paths, creating thousands of unique path label values in HTTP metrics"
}

func (m *Mission1) OnActivate(ctx context.Context) error {
	m.done = make(chan struct{})
	m.currentCardinality.Store(int32(m.baseCardinality))

	slog.Info("adversary path probing started",
		"component", "mission1",
		"initial_paths", m.baseCardinality,
		"max_paths", m.maxCardinality,
		"growth_rate", m.growthRate,
		"growth_interval_seconds", m.growthInterval.Seconds())

	// Send initial batch of probe requests
	m.sendProbes(0, m.baseCardinality)

	// Start background goroutine to grow cardinality
	m.ticker = time.NewTicker(m.growthInterval)
	go m.growCardinality()

	return nil
}

func (m *Mission1) OnDeactivate(ctx context.Context) error {
	if m.ticker != nil {
		m.ticker.Stop()
	}
	close(m.done)

	// Clean up garbage path metrics from the Prometheus registry so
	// subsequent scrapes no longer include them.
	count := int(m.currentCardinality.Load())
	for i := range count {
		path := m.garbagePaths[i]
		m.metrics.RequestsTotal.DeleteLabelValues("GET", path, "404")
		m.metrics.RequestDuration.DeleteLabelValues("GET", path)
	}

	slog.Info("adversary path probing stopped",
		"component", "mission1",
		"paths_cleaned", count)

	return nil
}

func (m *Mission1) growCardinality() {
	for {
		select {
		case <-m.ticker.C:
			current := m.currentCardinality.Load()
			if current < int32(m.maxCardinality) {
				newCardinality := min(current+int32(m.growthRate), int32(m.maxCardinality))
				m.sendProbes(int(current), int(newCardinality))
				m.currentCardinality.Store(newCardinality)

				slog.Info("adversary probe expansion",
					"component", "mission1",
					"previous_paths", current,
					"current_paths", newCardinality)
			}
		case <-m.done:
			return
		}
	}
}

// sendProbes sends GET requests to garbage paths in the range [from, to).
// Each request hits the server, goes through the metrics middleware, and creates
// a new http_requests_total + http_request_duration_seconds series with a unique path.
func (m *Mission1) sendProbes(from, to int) {
	for i := from; i < to; i++ {
		url := m.baseURL + m.garbagePaths[i]
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			continue
		}
		resp, err := m.client.Do(req)
		if err != nil {
			continue
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}
