package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all Prometheus metrics for the application
type Metrics struct {
	RequestsTotal           *prometheus.CounterVec
	IntelReportsProcessed   *prometheus.CounterVec
	MainframeCPUUtilization prometheus.Gauge
	ActiveAgents            *prometheus.GaugeVec
	AgentCommsTotal         *prometheus.CounterVec
	MissionActive           *prometheus.GaugeVec
	RequestDuration         *prometheus.HistogramVec
}

// InitMetrics creates and registers all Prometheus metrics
func InitMetrics() *Metrics {
	return &Metrics{
		RequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests to mission control",
			},
			[]string{"method", "path", "status"},
		),
		IntelReportsProcessed: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "intel_reports_processed_total",
				Help: "Total number of intelligence reports processed",
			},
			[]string{"endpoint"},
		),
		MainframeCPUUtilization: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "mainframe_cpu_utilization",
				Help: "Current CPU utilization of the mainframe (0-100)",
			},
		),
		ActiveAgents: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "active_agents",
				Help: "Active field agents reporting status (1=active, 0=inactive)",
			},
			[]string{"agent_id", "country_code"},
		),
		AgentCommsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "agent_comms_total",
				Help: "Total communications sent by field agents",
			},
			[]string{"id", "region"},
		),
		MissionActive: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "mission_active",
				Help: "Whether a mission is currently active (1=active, 0=inactive)",
			},
			[]string{"mission_id"},
		),
		RequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
	}
}
