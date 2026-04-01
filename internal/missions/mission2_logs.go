package missions

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"
)

type debugLogTemplate struct {
	msg       string
	component string
	attrs     func() []any
}

// Mission2 implements the Log Routing scenario
// When active, emits DEBUG level logs that need to be routed to S3
type Mission2 struct {
	templates []debugLogTemplate
	ticker    *time.Ticker
	done      chan struct{}
}

// NewMission2 creates a new Log Routing mission
func NewMission2() *Mission2 {
	tables := []string{"users", "sessions", "orders", "inventory", "audit_log"}
	queues := []string{"order-events", "notifications", "audit-stream", "telemetry-ingest"}
	upstreams := []string{"auth-service", "inventory-api", "payment-gateway", "notification-svc"}
	ciphers := []string{"TLS_AES_128_GCM_SHA256", "TLS_AES_256_GCM_SHA384", "TLS_CHACHA20_POLY1305_SHA256"}
	buckets := []string{"api-global", "api-per-user", "webhook-inbound", "batch-jobs"}
	configHashes := []string{"a3f8c2d1", "b7e4a9f0", "c1d6b3e8", "d9f2c7a4"}

	return &Mission2{
		templates: []debugLogTemplate{
			{
				msg:       "request body parsed",
				component: "http",
				attrs: func() []any {
					return []any{
						"content_length", rand.Intn(8192) + 64,
						"content_type", "application/json",
						"parse_duration_us", rand.Intn(500) + 10,
					}
				},
			},
			{
				msg:       "cache lookup completed",
				component: "cache",
				attrs: func() []any {
					hit := rand.Float64() > 0.3
					return []any{
						"cache_hit", hit,
						"key_prefix", fmt.Sprintf("sess:%04x", rand.Intn(0xFFFF)),
						"lookup_duration_us", rand.Intn(200) + 5,
					}
				},
			},
			{
				msg:       "connection pool stats",
				component: "db",
				attrs: func() []any {
					max := 25
					active := rand.Intn(max-2) + 1
					idle := max - active
					return []any{
						"active_conns", active,
						"idle_conns", idle,
						"max_conns", max,
					}
				},
			},
			{
				msg:       "database query executed",
				component: "db",
				attrs: func() []any {
					return []any{
						"query_duration_ms", rand.Intn(45) + 1,
						"rows_returned", rand.Intn(100),
						"table", tables[rand.Intn(len(tables))],
					}
				},
			},
			{
				msg:       "outbound request completed",
				component: "http",
				attrs: func() []any {
					codes := []int{200, 200, 200, 200, 201, 204, 304}
					return []any{
						"upstream", upstreams[rand.Intn(len(upstreams))],
						"status_code", codes[rand.Intn(len(codes))],
						"response_time_ms", rand.Intn(120) + 5,
					}
				},
			},
			{
				msg:       "tls handshake completed",
				component: "security",
				attrs: func() []any {
					return []any{
						"cipher_suite", ciphers[rand.Intn(len(ciphers))],
						"protocol_version", "TLSv1.3",
					}
				},
			},
			{
				msg:       "gc stats",
				component: "runtime",
				attrs: func() []any {
					return []any{
						"heap_alloc_mb", float64(rand.Intn(128)+32) + rand.Float64(),
						"heap_objects", rand.Intn(500000) + 50000,
						"gc_pause_us", rand.Intn(800) + 50,
					}
				},
			},
			{
				msg:       "worker pool status",
				component: "pool",
				attrs: func() []any {
					return []any{
						"active_workers", rand.Intn(8) + 1,
						"queue_depth", rand.Intn(50),
						"total_processed", rand.Intn(10000) + 500,
					}
				},
			},
			{
				msg:       "rate limiter check",
				component: "ratelimit",
				attrs: func() []any {
					tokens := rand.Intn(100)
					return []any{
						"bucket", buckets[rand.Intn(len(buckets))],
						"tokens_remaining", tokens,
						"allowed", tokens > 0,
					}
				},
			},
			{
				msg:       "session validated",
				component: "auth",
				attrs: func() []any {
					return []any{
						"session_id", fmt.Sprintf("%08x", rand.Int31()),
						"ttl_remaining_s", rand.Intn(3600) + 60,
						"refresh_needed", rand.Float64() > 0.7,
					}
				},
			},
			{
				msg:       "message queue poll",
				component: "queue",
				attrs: func() []any {
					return []any{
						"queue_name", queues[rand.Intn(len(queues))],
						"messages_available", rand.Intn(25),
						"consumer_lag", rand.Intn(100),
					}
				},
			},
			{
				msg:       "config reload check",
				component: "config",
				attrs: func() []any {
					hash := configHashes[rand.Intn(len(configHashes))]
					return []any{
						"config_hash", hash,
						"last_modified", time.Now().Add(-time.Duration(rand.Intn(3600)) * time.Second).Format(time.RFC3339),
						"changed", rand.Float64() > 0.9,
					}
				},
			},
		},
	}
}

func (m *Mission2) ID() string {
	return "mission2"
}

func (m *Mission2) Name() string {
	return "Log Routing"
}

func (m *Mission2) Description() string {
	return "Generates DEBUG level logs that should be routed to S3 storage"
}

func (m *Mission2) OnActivate(ctx context.Context) error {
	m.done = make(chan struct{})

	slog.Info("verbose logging enabled", "component", "diagnostics", "level", "DEBUG")

	m.ticker = time.NewTicker(250 * time.Millisecond)
	go m.emitDebugLogs()

	return nil
}

func (m *Mission2) OnDeactivate(ctx context.Context) error {
	if m.ticker != nil {
		m.ticker.Stop()
	}
	close(m.done)

	slog.Info("verbose logging disabled", "component", "diagnostics")
	return nil
}

func (m *Mission2) emitDebugLogs() {
	for {
		select {
		case <-m.ticker.C:
			numLogs := 1 + rand.Intn(3)

			for range numLogs {
				tmpl := m.templates[rand.Intn(len(m.templates))]
				attrs := append([]any{"component", tmpl.component}, tmpl.attrs()...)
				slog.Debug(tmpl.msg, attrs...)
			}
		case <-m.done:
			return
		}
	}
}
