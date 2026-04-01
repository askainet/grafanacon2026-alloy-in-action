package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/grafana/alloy-mission-control/internal/telemetry"
)

// MetricsMiddleware records HTTP request metrics
func MetricsMiddleware(metrics *telemetry.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture status code
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Call next handler
			next.ServeHTTP(rw, r)

			// Record metrics
			duration := time.Since(start)
			statusStr := strconv.Itoa(rw.statusCode)

			metrics.RequestsTotal.WithLabelValues(r.Method, r.URL.Path, statusStr).Inc()
			metrics.RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration.Seconds())
		})
	}
}
