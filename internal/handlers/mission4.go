package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Mission4Handler struct {
	alloyEndpoint   string
	tempoEndpoint   string
	numChunks       int
	isMissionActive func() bool
	client          *http.Client
}

func NewMission4Handler(alloyEndpoint, tempoEndpoint string, numChunks int, isMissionActive func() bool) *Mission4Handler {
	return &Mission4Handler{
		alloyEndpoint:   alloyEndpoint,
		tempoEndpoint:   tempoEndpoint,
		numChunks:       numChunks,
		isMissionActive: isMissionActive,
		client:          &http.Client{Timeout: 15 * time.Second},
	}
}

type mission4VerifyResponse struct {
	Status   string                   `json:"status"`
	Message  string                   `json:"message"`
	Pipeline *mission4PipelineDetails `json:"pipeline,omitempty"`
	Keys     *mission4KeyDetails      `json:"keys,omitempty"`
}

type mission4PipelineDetails struct {
	TailSamplerFound      bool    `json:"tail_sampler_found"`
	ErrorPolicyFound      bool    `json:"error_policy_found"`
	SamplingPolicyFound   bool    `json:"sampling_policy_found"`
	TracesReceived        float64 `json:"traces_received"`
	GlobalSampled         float64 `json:"global_sampled"`
	GlobalNotSampled      float64 `json:"global_not_sampled"`
	GlobalSamplingPercent float64 `json:"global_sampling_percent"`
	ErrorsSampled         float64 `json:"errors_sampled"`
	PolicyErrors          float64 `json:"policy_errors"`
}

type mission4KeyDetails struct {
	TotalFound  int `json:"total_found"`
	TotalNeeded int `json:"total_needed"`
	AllFound    bool `json:"all_found"`
}

// Verify handles GET /admin/mission4/verify
// It checks Alloy metrics for correct tail_sampling configuration and queries
// Tempo to verify that all key fragments are recoverable from error traces.
func (h *Mission4Handler) Verify(w http.ResponseWriter, r *http.Request) {
	if !h.isMissionActive() {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mission4VerifyResponse{
			Status:  "fail",
			Message: "mission 4 is not active - run 'make mission4' first",
		})
		return
	}

	pipeline, pipelineErr := h.checkPipeline()
	keys := h.checkKeys()

	result := h.evaluate(pipeline, pipelineErr, keys)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// checkPipeline queries Alloy's /metrics endpoint and parses tail_sampling metrics.
func (h *Mission4Handler) checkPipeline() (*mission4PipelineDetails, error) {
	resp, err := h.client.Get(fmt.Sprintf("%s/metrics", h.alloyEndpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to query Alloy metrics: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Alloy metrics returned status %d", resp.StatusCode)
	}

	details := &mission4PipelineDetails{}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		if !strings.Contains(line, "tail_sampling") {
			continue
		}

		details.TailSamplerFound = true

		// Per-policy sampled counts
		if strings.HasPrefix(line, "otelcol_processor_tail_sampling_count_traces_sampled_total") {
			if strings.Contains(line, "keep_all_errors") {
				details.ErrorPolicyFound = true
				if strings.Contains(line, `sampled="true"`) {
					if v, err := parsePrometheusValue(line); err == nil {
						details.ErrorsSampled += v
					}
				}
			}
			if strings.Contains(line, "sample_normal_traffic") {
				details.SamplingPolicyFound = true
			}
		}

		// Global sampled/not_sampled counts
		if strings.HasPrefix(line, "otelcol_processor_tail_sampling_global_count_traces_sampled_total") {
			if strings.Contains(line, `sampled="true"`) {
				if v, err := parsePrometheusValue(line); err == nil {
					details.GlobalSampled += v
				}
			} else if strings.Contains(line, `sampled="false"`) {
				if v, err := parsePrometheusValue(line); err == nil {
					details.GlobalNotSampled += v
				}
			}
		}

		// Traces received
		if strings.HasPrefix(line, "otelcol_processor_tail_sampling_new_trace_id_received_total") {
			if v, err := parsePrometheusValue(line); err == nil {
				details.TracesReceived += v
			}
		}

		// Policy evaluation errors
		if strings.HasPrefix(line, "otelcol_processor_tail_sampling_sampling_policy_evaluation_error_total") {
			if v, err := parsePrometheusValue(line); err == nil {
				details.PolicyErrors += v
			}
		}
	}

	total := details.GlobalSampled + details.GlobalNotSampled
	if total > 0 {
		details.GlobalSamplingPercent = (details.GlobalSampled / total) * 100
	}

	// The error policy may not have a sampled="true" line yet if no errors
	// have been evaluated, but the policy is still present if we saw any
	// count_traces_sampled_total line with keep_all_errors in it.
	return details, nil
}

// checkKeys queries Tempo for error traces containing key fragments.
func (h *Mission4Handler) checkKeys() *mission4KeyDetails {
	traceQL := `{ span.key_sequence >= 1 && span.key_chunk != "" && status = error && resource.service.name = "mission-control" }`
	searchURL := fmt.Sprintf("%s/api/search?q=%s", h.tempoEndpoint, url.QueryEscape(traceQL))

	resp, err := h.client.Get(searchURL)
	if err != nil {
		return &mission4KeyDetails{TotalNeeded: h.numChunks}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &mission4KeyDetails{TotalNeeded: h.numChunks}
	}

	var tempoResp tempoSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&tempoResp); err != nil {
		return &mission4KeyDetails{TotalNeeded: h.numChunks}
	}

	found := make(map[int]bool)
	for _, trace := range tempoResp.Traces {
		for _, spanSet := range trace.SpanSets {
			for _, span := range spanSet.Spans {
				seq, chunk := extractKeyData(span.Attributes)
				if seq > 0 && chunk != "" {
					found[seq] = true
				}
			}
		}
	}

	return &mission4KeyDetails{
		TotalFound:  len(found),
		TotalNeeded: h.numChunks,
		AllFound:    len(found) == h.numChunks,
	}
}

func (h *Mission4Handler) evaluate(pipeline *mission4PipelineDetails, pipelineErr error, keys *mission4KeyDetails) mission4VerifyResponse {
	if pipelineErr != nil {
		return mission4VerifyResponse{
			Status:  "fail",
			Message: fmt.Sprintf("could not check Alloy pipeline: %v", pipelineErr),
			Keys:    keys,
		}
	}

	if !pipeline.TailSamplerFound {
		return mission4VerifyResponse{
			Status:   "fail",
			Message:  "no tail_sampling metrics found - configure otelcol.processor.tail_sampling in your Alloy pipeline",
			Pipeline: pipeline,
			Keys:     keys,
		}
	}

	if pipeline.TracesReceived == 0 {
		return mission4VerifyResponse{
			Status:   "fail",
			Message:  "tail_sampling processor found but no traces received - ensure the OTLP receiver routes traces through the tail sampler",
			Pipeline: pipeline,
			Keys:     keys,
		}
	}

	if !pipeline.SamplingPolicyFound {
		return mission4VerifyResponse{
			Status:   "fail",
			Message:  "missing 'sample_normal_traffic' policy - add a probabilistic policy to sample ~5% of normal traces",
			Pipeline: pipeline,
			Keys:     keys,
		}
	}

	if pipeline.PolicyErrors > 0 {
		return mission4VerifyResponse{
			Status:   "fail",
			Message:  fmt.Sprintf("%.0f policy evaluation errors detected - check your tail_sampling configuration", pipeline.PolicyErrors),
			Pipeline: pipeline,
			Keys:     keys,
		}
	}

	if keys.AllFound {
		return mission4VerifyResponse{
			Status:  "pass",
			Message: fmt.Sprintf("tail sampling is working - all %d key fragments recovered", keys.TotalNeeded),
		}
	}

	if pipeline.ErrorsSampled > 0 {
		return mission4VerifyResponse{
			Status:   "partial",
			Message:  fmt.Sprintf("tail sampling is configured and error traces are being kept, but only %d/%d key fragments recovered from Tempo - keep the mission running and check again", keys.TotalFound, keys.TotalNeeded),
			Pipeline: pipeline,
			Keys:     keys,
		}
	}

	return mission4VerifyResponse{
		Status:   "partial",
		Message:  fmt.Sprintf("tail sampling is receiving traces but no error traces sampled yet (%d/%d key fragments found) - ensure the 'keep_all_errors' status_code policy is configured and wait for error traces to arrive", keys.TotalFound, keys.TotalNeeded),
		Pipeline: pipeline,
		Keys:     keys,
	}
}
