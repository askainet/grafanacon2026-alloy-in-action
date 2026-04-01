package traffic

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// Endpoint represents a target endpoint with weight for distribution
type Endpoint struct {
	Method string
	Path   string
	Weight int // relative weight for distribution
}

// Generator creates synthetic traffic to intel endpoints
type Generator struct {
	baseURL     string
	client      *http.Client
	workers     int
	endpoints   []Endpoint
	stopCh      chan struct{}
	wg          sync.WaitGroup
}

// NewGenerator creates a new traffic generator with multiple workers
func NewGenerator(baseURL string, workers int) *Generator {
	if workers <= 0 {
		workers = 3 // Default to 3 workers
	}

	// Define endpoints with weighted distribution
	// 40% retina scan, 30% intel, 15% heli pickup, 10% burn message, 5% dead drop, 5% sat link
	endpoints := []Endpoint{
		{Method: "POST", Path: "/api/biometrics/retina-scan/verify", Weight: 40},
		{Method: "POST", Path: "/api/process-intel", Weight: 25},
		{Method: "POST", Path: "/api/extraction/heli-pickup/request", Weight: 15},
		{Method: "GET", Path: "/api/message/burn-after-reading", Weight: 10},
		{Method: "POST", Path: "/api/comms/dead-drop/upload", Weight: 5},
		{Method: "POST", Path: "/api/sat-link/handshake/establish", Weight: 5},
	}

	return &Generator{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		workers:   workers,
		endpoints: endpoints,
		stopCh:    make(chan struct{}),
	}
}

// Start begins generating synthetic traffic with multiple workers
func (g *Generator) Start(ctx context.Context) {
	for i := 0; i < g.workers; i++ {
		g.wg.Add(1)
		go g.generateTraffic(ctx, i+1)
	}
	slog.Info("synthetic traffic generator started", "component", "traffic", "workers", g.workers)
}

// Stop halts traffic generation and waits for workers to finish
func (g *Generator) Stop() {
	close(g.stopCh)
	g.wg.Wait()
	slog.Info("synthetic traffic generator stopped", "component", "traffic")
}

// generateTraffic runs in a goroutine and generates periodic requests
func (g *Generator) generateTraffic(ctx context.Context, workerID int) {
	defer g.wg.Done()

	ticker := time.NewTicker(g.randomInterval())
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-g.stopCh:
			return
		case <-ticker.C:
			// Occasionally fire a burst of 3-6 rapid requests (~15% chance)
			burstCount := 1
			if rand.Intn(100) < 15 {
				burstCount = rand.Intn(4) + 3
			}

			for i := 0; i < burstCount; i++ {
				endpoint := g.selectEndpoint()
				g.makeRequest(endpoint, workerID)
				if i < burstCount-1 {
					time.Sleep(time.Duration(100+rand.Intn(400)) * time.Millisecond)
				}
			}

			// Reset ticker with random interval for more realistic traffic
			ticker.Reset(g.randomInterval())
		}
	}
}

// selectEndpoint selects an endpoint based on weighted distribution
func (g *Generator) selectEndpoint() Endpoint {
	// Calculate total weight
	totalWeight := 0
	for _, ep := range g.endpoints {
		totalWeight += ep.Weight
	}

	// Select random value in range [0, totalWeight)
	r := rand.Intn(totalWeight)

	// Find the endpoint that corresponds to this value
	cumulative := 0
	for _, ep := range g.endpoints {
		cumulative += ep.Weight
		if r < cumulative {
			return ep
		}
	}

	// Fallback (should never happen)
	return g.endpoints[0]
}

// makeRequest sends a request to the specified endpoint
func (g *Generator) makeRequest(endpoint Endpoint, workerID int) {
	url := g.baseURL + endpoint.Path

	var req *http.Request
	var err error

	if endpoint.Method == "POST" {
		req, err = http.NewRequest("POST", url, bytes.NewReader([]byte("{}")))
	} else {
		req, err = http.NewRequest(endpoint.Method, url, nil)
	}

	if err != nil {
		slog.Warn("failed to create synthetic request",
			"component", "traffic",
			"error", err,
			"worker", workerID,
			"endpoint", endpoint.Path,
		)
		return
	}

	if endpoint.Method == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("User-Agent", "IMF-FieldAgent/1.0")

	resp, err := g.client.Do(req)
	if err != nil {
		slog.Warn("synthetic traffic request failed",
			"component", "traffic",
			"error", err,
			"worker", workerID,
			"endpoint", endpoint.Path,
		)
		return
	}
	defer resp.Body.Close()

	// Drain and close the response body
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusOK {
		slog.Warn("synthetic traffic request returned non-200 status",
			"component", "traffic",
			"status", resp.StatusCode,
			"worker", workerID,
			"endpoint", endpoint.Path,
		)
	}
}

// randomInterval returns a random duration between 2-8 seconds
// This creates more realistic, varied traffic patterns
func (g *Generator) randomInterval() time.Duration {
	// Random duration between 2-8 seconds
	min := 2000
	max := 8000
	ms := rand.Intn(max-min) + min
	return time.Duration(ms) * time.Millisecond
}
