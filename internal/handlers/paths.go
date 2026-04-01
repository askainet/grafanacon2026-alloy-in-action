package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

// PathsHandler serves the allowed-paths regex for Mission 1 cardinality filtering.
type PathsHandler struct {
	validPaths []string
}

func NewPathsHandler(validPaths []string) *PathsHandler {
	return &PathsHandler{validPaths: validPaths}
}

// GetAllowedPathsRegex returns a regex matching all legitimate HTTP paths.
// GET /api/metrics/allowed-paths
func (h *PathsHandler) GetAllowedPathsRegex(w http.ResponseWriter, r *http.Request) {
	// ^$ matches metrics without a path label (active_agents, go_*, etc.)
	// /admin/.+ matches all admin endpoints (dynamic IDs in routes)
	regex := "^$|^(" + strings.Join(h.validPaths, "|") + "|/admin/.+)$"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"regex": regex,
	})
}
