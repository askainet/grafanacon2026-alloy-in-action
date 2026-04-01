package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/grafana/alloy-mission-control/internal/telemetry"
)

type AgentHandler struct {
	agentManager *telemetry.AgentManager
	metrics      *telemetry.Metrics
}

func NewAgentHandler(agentManager *telemetry.AgentManager, metrics *telemetry.Metrics) *AgentHandler {
	return &AgentHandler{
		agentManager: agentManager,
		metrics:      metrics,
	}
}

// Heartbeat handles agent heartbeat reports and updates metrics
// POST /api/agents/heartbeat
// Body: {"agent_id": "ALPHA-001", "country_code": "US", "status": "active"}
func (h *AgentHandler) Heartbeat(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AgentID     string `json:"agent_id"`
		CountryCode string `json:"country_code"`
		Status      string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	if req.AgentID == "" || req.CountryCode == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "agent_id and country_code are required"})
		return
	}

	// Update active_agents gauge
	var gaugeValue float64
	if req.Status == "active" {
		gaugeValue = 1
	}
	h.metrics.ActiveAgents.With(prometheus.Labels{
		"agent_id":     req.AgentID,
		"country_code": req.CountryCode,
	}).Set(gaugeValue)

	// Increment agent_comms_total counter
	h.metrics.AgentCommsTotal.With(prometheus.Labels{
		"id":     req.AgentID,
		"region": req.CountryCode,
	}).Inc()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// GetLegitimateAgentRegex returns a regex pattern for filtering legitimate agents
func (h *AgentHandler) GetLegitimateAgentRegex(w http.ResponseWriter, r *http.Request) {
	regex := h.agentManager.GetLegitimateAgentRegex()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"regex": regex,
	})
}

// ListLegitimateAgents returns the list of legitimate agents
func (h *AgentHandler) ListLegitimateAgents(w http.ResponseWriter, r *http.Request) {
	agents := h.agentManager.GetLegitimateAgents()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"count":  len(agents),
		"agents": agents,
	})
}
