package handlers

import (
	"encoding/json"
	"net/http"
)

// SDTarget represents a Prometheus HTTP SD target group.
type SDTarget struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

// SDHandler serves Prometheus HTTP Service Discovery format responses
// containing the LGTM stack services with enrichment metadata.
type SDHandler struct {
	targets []SDTarget
}

func NewSDHandler() *SDHandler {
	return &SDHandler{
		targets: []SDTarget{
			{
				Targets: []string{"grafana:3000"},
				Labels: map[string]string{
					"service_name":      "grafana",
					"team":              "intel-analytics",
					"aws_account_id":    "194851037629",
					"aws_region":        "us-east-1",
					"availability_zone": "us-east-1a",
					"vpc_id":            "vpc-0a3c7d1e9f42b8d56",
					"instance_type":     "m5.xlarge",
					"environment":       "classified",
				},
			},
			{
				Targets: []string{"loki:3100"},
				Labels: map[string]string{
					"service_name":      "loki",
					"team":              "sigint",
					"aws_account_id":    "194851037629",
					"aws_region":        "us-east-1",
					"availability_zone": "us-east-1b",
					"vpc_id":            "vpc-0a3c7d1e9f42b8d56",
					"instance_type":     "r6g.2xlarge",
					"environment":       "classified",
				},
			},
			{
				Targets: []string{"tempo:3200"},
				Labels: map[string]string{
					"service_name":      "tempo",
					"team":              "sigint",
					"aws_account_id":    "194851037629",
					"aws_region":        "us-east-1",
					"availability_zone": "us-east-1c",
					"vpc_id":            "vpc-0a3c7d1e9f42b8d56",
					"instance_type":     "r6g.2xlarge",
					"environment":       "classified",
				},
			},
			{
				Targets: []string{"mimir:9009"},
				Labels: map[string]string{
					"service_name":      "mimir",
					"team":              "sigint",
					"aws_account_id":    "194851037629",
					"aws_region":        "us-east-1",
					"availability_zone": "us-east-1a",
					"vpc_id":            "vpc-0a3c7d1e9f42b8d56",
					"instance_type":     "r6g.4xlarge",
					"environment":       "classified",
				},
			},
		},
	}
}

// ServeHTTP returns the service discovery targets in Prometheus HTTP SD format.
// GET /api/sd/targets
func (h *SDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.targets)
}
