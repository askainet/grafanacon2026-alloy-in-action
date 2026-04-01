package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Mission3Handler struct {
	alloyEndpoint   string
	isMissionActive func() bool
	client          *http.Client
}

func NewMission3Handler(alloyEndpoint string, isMissionActive func() bool) *Mission3Handler {
	return &Mission3Handler{
		alloyEndpoint:   alloyEndpoint,
		isMissionActive: isMissionActive,
		client:          &http.Client{Timeout: 15 * time.Second},
	}
}

type mission3VerifyResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Details *mission3VerifyDetails `json:"details,omitempty"`
}

type mission3VerifyDetails struct {
	IncomingSpans   float64 `json:"incoming_spans"`
	OutgoingSpans   float64 `json:"outgoing_spans"`
	SamplingPercent float64 `json:"sampling_percent"`
}

// Verify handles GET /admin/mission3/verify
// It queries Alloy's /metrics endpoint for the probabilistic_sampler processor
// metrics to determine whether head sampling is correctly configured.
func (h *Mission3Handler) Verify(w http.ResponseWriter, r *http.Request) {
	if !h.isMissionActive() {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mission3VerifyResponse{
			Status:  "fail",
			Message: "mission 3 is not active - run 'make mission3' first",
		})
		return
	}

	metricsURL := fmt.Sprintf("%s/metrics", h.alloyEndpoint)
	resp, err := h.client.Get(metricsURL)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("failed to query Alloy metrics: %v", err),
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Alloy metrics returned status %d", resp.StatusCode),
		})
		return
	}

	// Parse Prometheus text format looking for probabilistic_sampler processor metrics.
	// We sum all matching lines in case there are multiple label sets.
	var incoming, outgoing float64
	var foundIncoming, foundOutgoing bool

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || !strings.Contains(line, "probabilistic_sampler") {
			continue
		}
		if strings.HasPrefix(line, "otelcol_processor_incoming_items_total") {
			if v, err := parsePrometheusValue(line); err == nil {
				incoming += v
				foundIncoming = true
			}
		} else if strings.HasPrefix(line, "otelcol_processor_outgoing_items_total") {
			if v, err := parsePrometheusValue(line); err == nil {
				outgoing += v
				foundOutgoing = true
			}
		}
	}

	if !foundIncoming && !foundOutgoing {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mission3VerifyResponse{
			Status:  "fail",
			Message: "no probabilistic_sampler metrics found - configure otelcol.processor.probabilistic_sampler in your Alloy pipeline",
		})
		return
	}

	if incoming == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mission3VerifyResponse{
			Status:  "fail",
			Message: "probabilistic_sampler found but no spans processed yet - ensure the OTLP receiver routes traces through the sampler",
		})
		return
	}

	samplingPercent := (outgoing / incoming) * 100

	details := &mission3VerifyDetails{
		IncomingSpans:   incoming,
		OutgoingSpans:   outgoing,
		SamplingPercent: samplingPercent,
	}

	var result mission3VerifyResponse

	switch {
	case outgoing == incoming:
		result = mission3VerifyResponse{
			Status:  "fail",
			Message: "sampler is not dropping any spans - check sampling_percentage is set correctly",
			Details: details,
		}
	case samplingPercent <= 20:
		result = mission3VerifyResponse{
			Status:  "pass",
			Message: fmt.Sprintf("head sampling is working - %.1f%% of spans sampled (%.0f incoming, %.0f outgoing)", samplingPercent, incoming, outgoing),
			Details: details,
		}
	default:
		result = mission3VerifyResponse{
			Status:  "fail",
			Message: fmt.Sprintf("sampling percentage too high - %.1f%% of spans sampled, target is ~5%%", samplingPercent),
			Details: details,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// parsePrometheusValue extracts the numeric value from a Prometheus text format line.
// Format: metric_name{labels} value
func parsePrometheusValue(line string) (float64, error) {
	idx := strings.LastIndex(line, " ")
	if idx < 0 {
		return 0, fmt.Errorf("no value found in metric line")
	}
	return strconv.ParseFloat(strings.TrimSpace(line[idx+1:]), 64)
}
