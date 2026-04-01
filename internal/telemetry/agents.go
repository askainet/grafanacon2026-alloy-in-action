package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"math"
	"math/rand"
	"net/http"
	"time"
)

// Agent represents a field agent
type Agent struct {
	ID          string
	CountryCode string
	Active      bool
}

// AgentManager manages legitimate field agents and their metrics
type AgentManager struct {
	metrics          *Metrics
	legitimateAgents []Agent
	baseURL          string
	client           *http.Client
	ticker           *time.Ticker
	done             chan struct{}
}

// NewAgentManager creates a new agent manager
func NewAgentManager(metrics *Metrics, baseURL string) *AgentManager {
	legitimateAgents := []Agent{
		// United States - 5 agents
		{ID: "ALPHA-001", CountryCode: "US", Active: true},
		{ID: "ALPHA-002", CountryCode: "US", Active: true},
		{ID: "ALPHA-003", CountryCode: "US", Active: true},
		{ID: "ALPHA-004", CountryCode: "US", Active: true},
		{ID: "ALPHA-005", CountryCode: "US", Active: true},
		// United Kingdom - 4 agents
		{ID: "BRAVO-001", CountryCode: "GB", Active: true},
		{ID: "BRAVO-002", CountryCode: "GB", Active: true},
		{ID: "BRAVO-003", CountryCode: "GB", Active: true},
		{ID: "BRAVO-004", CountryCode: "GB", Active: true},
		// Germany - 4 agents
		{ID: "CHARLIE-001", CountryCode: "DE", Active: true},
		{ID: "CHARLIE-002", CountryCode: "DE", Active: true},
		{ID: "CHARLIE-003", CountryCode: "DE", Active: true},
		{ID: "CHARLIE-004", CountryCode: "DE", Active: true},
		// Japan - 3 agents
		{ID: "DELTA-001", CountryCode: "JP", Active: true},
		{ID: "DELTA-002", CountryCode: "JP", Active: true},
		{ID: "DELTA-003", CountryCode: "JP", Active: true},
		// France - 3 agents
		{ID: "ECHO-001", CountryCode: "FR", Active: true},
		{ID: "ECHO-002", CountryCode: "FR", Active: true},
		{ID: "ECHO-003", CountryCode: "FR", Active: true},
		// Australia - 3 agents
		{ID: "FOXTROT-001", CountryCode: "AU", Active: true},
		{ID: "FOXTROT-002", CountryCode: "AU", Active: true},
		{ID: "FOXTROT-003", CountryCode: "AU", Active: true},
	}

	return &AgentManager{
		metrics:          metrics,
		legitimateAgents: legitimateAgents,
		baseURL:          baseURL,
		client:           &http.Client{Timeout: 10 * time.Second},
		done:             make(chan struct{}),
	}
}

// Start begins updating agent metrics in the background
func (am *AgentManager) Start(ctx context.Context) {
	slog.Info("starting agent manager",
		"component", "agents",
		"legitimate_agents", len(am.legitimateAgents))

	// Send initial heartbeats for all agents
	for _, agent := range am.legitimateAgents {
		am.sendHeartbeat(agent.ID, agent.CountryCode, "active")
	}

	// Update agent status periodically to simulate field activity
	am.ticker = time.NewTicker(30 * time.Second)
	go am.updateAgentStatus()

	// Also update mainframe CPU utilization
	go am.updateMainframeCPU()
}

// Stop halts the agent manager
func (am *AgentManager) Stop() {
	if am.ticker != nil {
		am.ticker.Stop()
	}
	close(am.done)
	slog.Info("agent manager stopped", "component", "agents")
}

// updateAgentStatus randomly toggles agent status to simulate field operations
func (am *AgentManager) updateAgentStatus() {
	for {
		select {
		case <-am.ticker.C:
			// Toggle 1-4 agents with weighted probability: 1(40%) 2(30%) 3(20%) 4(10%)
			roll := rand.Intn(100)
			var numToToggle int
			switch {
			case roll < 40:
				numToToggle = 1
			case roll < 70:
				numToToggle = 2
			case roll < 90:
				numToToggle = 3
			default:
				numToToggle = 4
			}

			for i := 0; i < numToToggle; i++ {
				idx := rand.Intn(len(am.legitimateAgents))
				agent := &am.legitimateAgents[idx]
				agent.Active = !agent.Active

				status := "inactive"
				if agent.Active {
					status = "active"
				}

				am.sendHeartbeat(agent.ID, agent.CountryCode, status)

				slog.Info("agent status updated",
					"component", "agents",
					"agent_id", agent.ID,
					"active", agent.Active)
			}

			// Jitter the next tick: 20-40s instead of fixed 30s
			am.ticker.Reset(time.Duration(20+rand.Intn(21)) * time.Second)
		case <-am.done:
			return
		}
	}
}

// updateMainframeCPU updates the mainframe CPU utilization metric using a
// smooth sine wave with layered noise so the graph looks organic.
func (am *AgentManager) updateMainframeCPU() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	start := time.Now()

	for {
		select {
		case <-ticker.C:
			t := time.Since(start).Seconds()

			// Slow sine wave (period ~3min) as the base trend: 35-65 range
			base := 50 + 15*math.Sin(t/90*math.Pi)
			// Faster wobble (period ~30s): +/- 8
			wobble := 8 * math.Sin(t/15*math.Pi)
			// Random noise: +/- 5
			noise := float64(rand.Intn(11)) - 5

			cpuUtil := math.Max(5, math.Min(95, base+wobble+noise))
			am.metrics.MainframeCPUUtilization.Set(cpuUtil)

			slog.Info("mainframe CPU updated",
				"component", "agents",
				"utilization", cpuUtil)
		case <-am.done:
			return
		}
	}
}

// sendHeartbeat sends a heartbeat POST request to the heartbeat endpoint
func (am *AgentManager) sendHeartbeat(agentID, countryCode, status string) {
	body, _ := json.Marshal(map[string]string{
		"agent_id":     agentID,
		"country_code": countryCode,
		"status":       status,
	})

	url := am.baseURL + "/api/agents/heartbeat"
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := am.client.Do(req)
	if err != nil {
		slog.Warn("heartbeat request failed", "component", "agents", "agent_id", agentID, "error", err)
		return
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
}

// GetLegitimateAgents returns the list of legitimate agent IDs
func (am *AgentManager) GetLegitimateAgents() []Agent {
	return am.legitimateAgents
}

// GetLegitimateAgentRegex returns a regex pattern that matches only legitimate agent IDs
func (am *AgentManager) GetLegitimateAgentRegex() string {
	if len(am.legitimateAgents) == 0 {
		return ""
	}

	pattern := "^("
	for i, agent := range am.legitimateAgents {
		if i > 0 {
			pattern += "|"
		}
		pattern += agent.ID
	}
	pattern += ")$"

	return pattern
}
