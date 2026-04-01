package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"time"
)

type TempoHandler struct {
	tempoEndpoint string
	numChunks     int
	client        *http.Client
}

func NewTempoHandler(tempoEndpoint string, numChunks int) *TempoHandler {
	return &TempoHandler{
		tempoEndpoint: tempoEndpoint,
		numChunks:     numChunks,
		client:        &http.Client{Timeout: 15 * time.Second},
	}
}

type keyFragment struct {
	Sequence int    `json:"sequence"`
	Chunk    string `json:"chunk,omitempty"`
	Found    bool   `json:"found"`
}

type accessTokenResponse struct {
	Status      string        `json:"status"`
	Fragments   []keyFragment `json:"fragments"`
	TotalFound  int           `json:"total_found"`
	TotalNeeded int           `json:"total_needed"`
	Token       string        `json:"token"`
}

// ServeHTTP handles GET /admin/mission4/access-token
func (h *TempoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Build TraceQL query to find error spans with key_sequence attribute
	traceQL := `{ span.key_sequence >= 1 && span.key_chunk != "" && status = error && resource.service.name = "mission-control" }`

	searchURL := fmt.Sprintf("%s/api/search?q=%s", h.tempoEndpoint, url.QueryEscape(traceQL))

	resp, err := h.client.Get(searchURL)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("failed to query Tempo: %v", err),
		})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("failed to read Tempo response: %v", err),
		})
		return
	}

	if resp.StatusCode != http.StatusOK {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]string{
			"error":         fmt.Sprintf("Tempo returned status %d", resp.StatusCode),
			"tempo_response": string(body),
		})
		return
	}

	// Parse Tempo search response
	var tempoResp tempoSearchResponse
	if err := json.Unmarshal(body, &tempoResp); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("failed to parse Tempo response: %v", err),
		})
		return
	}

	// Extract and deduplicate fragments by sequence number
	found := make(map[int]string)
	for _, trace := range tempoResp.Traces {
		for _, spanSet := range trace.SpanSets {
			for _, span := range spanSet.Spans {
				seq, chunk := extractKeyData(span.Attributes)
				if seq > 0 && chunk != "" {
					found[seq] = chunk
				}
			}
		}
	}

	// Build fragment list
	fragments := make([]keyFragment, 0, h.numChunks)
	for i := 1; i <= h.numChunks; i++ {
		if chunk, ok := found[i]; ok {
			fragments = append(fragments, keyFragment{Sequence: i, Chunk: chunk, Found: true})
		} else {
			fragments = append(fragments, keyFragment{Sequence: i, Found: false})
		}
	}

	result := accessTokenResponse{
		Status:      "incomplete",
		Fragments:   fragments,
		TotalFound:  len(found),
		TotalNeeded: h.numChunks,
		Token:       "INCOMPLETE",
	}

	if len(found) == h.numChunks {
		result.Status = "complete"
		seqs := make([]int, 0, len(found))
		for seq := range found {
			seqs = append(seqs, seq)
		}
		sort.Ints(seqs)
		var assembled string
		for _, seq := range seqs {
			assembled += found[seq]
		}
		result.Token = assembled
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Tempo API response structures

type tempoSearchResponse struct {
	Traces []tempoTrace `json:"traces"`
}

type tempoTrace struct {
	TraceID  string         `json:"traceID"`
	SpanSets []tempoSpanSet `json:"spanSets"`
}

type tempoSpanSet struct {
	Spans []tempoSpan `json:"spans"`
}

type tempoSpan struct {
	SpanID     string           `json:"spanID"`
	Attributes []tempoAttribute `json:"attributes"`
}

type tempoAttribute struct {
	Key   string         `json:"key"`
	Value tempoAttrValue `json:"value"`
}

type tempoAttrValue struct {
	StringValue string `json:"stringValue,omitempty"`
	IntValue    string `json:"intValue,omitempty"`
}

func extractKeyData(attrs []tempoAttribute) (int, string) {
	var seq int
	var chunk string

	for _, attr := range attrs {
		switch attr.Key {
		case "key_sequence":
			if attr.Value.IntValue != "" {
				fmt.Sscanf(attr.Value.IntValue, "%d", &seq)
			}
		case "key_chunk":
			chunk = attr.Value.StringValue
		}
	}

	return seq, chunk
}
