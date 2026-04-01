package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type MimirHandler struct {
	mimirEndpoint   string
	legitimatePaths map[string]bool
	isMissionActive func() bool
	client          *http.Client
}

func NewMimirHandler(mimirEndpoint string, legitimatePaths []string, isMissionActive func() bool) *MimirHandler {
	pathSet := make(map[string]bool, len(legitimatePaths))
	for _, p := range legitimatePaths {
		pathSet[p] = true
	}
	return &MimirHandler{
		mimirEndpoint:   mimirEndpoint,
		legitimatePaths: pathSet,
		isMissionActive: isMissionActive,
		client:          &http.Client{Timeout: 15 * time.Second},
	}
}

// isLegitimatePath checks if a path is a known route.
// Admin paths use a prefix check since they contain dynamic segments.
func (h *MimirHandler) isLegitimatePath(path string) bool {
	if strings.HasPrefix(path, "/admin/") {
		return true
	}
	return h.legitimatePaths[path]
}

type mission1VerifyResponse struct {
	Status          string   `json:"status"`
	Message         string   `json:"message"`
	Window          string   `json:"window,omitempty"`
	TotalSeries     int      `json:"total_series,omitempty"`
	LegitimateCount int      `json:"legitimate_count,omitempty"`
	RogueCount      int      `json:"rogue_count,omitempty"`
	RoguePaths      []string `json:"rogue_paths,omitempty"`
}

// Verify handles GET /admin/mission1/verify
func (h *MimirHandler) Verify(w http.ResponseWriter, r *http.Request) {
	if !h.isMissionActive() {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mission1VerifyResponse{
			Status:  "fail",
			Message: "mission 1 is not active - run 'make mission1' first",
		})
		return
	}

	// Parse optional window parameter (default 30s)
	window := "30s"
	if w := r.URL.Query().Get("window"); w != "" {
		window = w
	}

	// Query http_requests_total for recent series. Each garbage path the adversary
	// probed creates a unique {path} label value visible here.
	query := fmt.Sprintf("last_over_time(http_requests_total[%s])", window)
	queryURL := fmt.Sprintf("%s/prometheus/api/v1/query?query=%s", h.mimirEndpoint, url.QueryEscape(query))

	resp, err := h.client.Get(queryURL)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("failed to query Mimir: %v", err),
		})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("failed to read Mimir response: %v", err),
		})
		return
	}

	if resp.StatusCode != http.StatusOK {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]string{
			"error":          fmt.Sprintf("Mimir returned status %d", resp.StatusCode),
			"mimir_response": string(body),
		})
		return
	}

	var promResp prometheusQueryResponse
	if err := json.Unmarshal(body, &promResp); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("failed to parse Mimir response: %v", err),
		})
		return
	}

	if promResp.Status != "success" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Mimir query failed: %s", promResp.Status),
		})
		return
	}

	// Extract path labels and classify as legitimate or rogue
	var roguePaths []string
	legitimateCount := 0

	for _, result := range promResp.Data.Result {
		path, ok := result.Metric["path"]
		if !ok {
			continue
		}
		if h.isLegitimatePath(path) {
			legitimateCount++
		} else {
			roguePaths = append(roguePaths, path)
		}
	}

	totalSeries := len(promResp.Data.Result)

	switch {
	case totalSeries == 0:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mission1VerifyResponse{
			Status:  "fail",
			Message: "no http_requests_total series found in Mimir - is the metrics pipeline configured?",
		})
	case len(roguePaths) == 0 && legitimateCount > 0:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mission1VerifyResponse{
			Status:  "pass",
			Message: fmt.Sprintf("only legitimate paths found in Mimir (%d series) - rogue paths successfully filtered", legitimateCount),
		})
	case len(roguePaths) > 0:
		result := mission1VerifyResponse{
			Status:          "fail",
			Message:         fmt.Sprintf("%d rogue path(s) still reaching Mimir - relabel rules need work", len(roguePaths)),
			Window:          window,
			TotalSeries:     totalSeries,
			LegitimateCount: legitimateCount,
			RogueCount:      len(roguePaths),
		}
		if len(roguePaths) > 10 {
			result.RoguePaths = roguePaths[:10]
		} else {
			result.RoguePaths = roguePaths
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	default:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mission1VerifyResponse{
			Status:  "fail",
			Message: "http_requests_total series found but no path labels - check the metric labels",
		})
	}
}

// Prometheus API response structures

type prometheusQueryResponse struct {
	Status string              `json:"status"`
	Data   prometheusQueryData `json:"data"`
}

type prometheusQueryData struct {
	ResultType string                  `json:"resultType"`
	Result     []prometheusQueryResult `json:"result"`
}

type prometheusQueryResult struct {
	Metric map[string]string `json:"metric"`
	Value  json.RawMessage   `json:"value"`
}
