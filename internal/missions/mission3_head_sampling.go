package missions

import (
	"context"
	"log/slog"
	"math/rand"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Mission3 implements the Trace Sampling - Part 1 scenario
// When active, generates detailed traces with multiple child spans
// to demonstrate the need for probabilistic sampling
type Mission3 struct {
	operations []string
	tracer     trace.Tracer
	ticker     *time.Ticker
	done       chan struct{}
}

// NewMission3 creates a new Trace Sampling (Probabilistic) mission
func NewMission3() *Mission3 {
	return &Mission3{
		operations: []string{
			"validate_credentials",
			"decrypt_payload",
			"verify_signature",
			"check_permissions",
			"update_database",
			"send_notification",
			"audit_log",
		},
	}
}

func (m *Mission3) ID() string {
	return "mission3"
}

func (m *Mission3) Name() string {
	return "Trace Sampling - Probabilistic"
}

func (m *Mission3) Description() string {
	return "Generates detailed traces with multiple spans to demonstrate probabilistic sampling"
}

func (m *Mission3) OnActivate(ctx context.Context) error {
	// Recreate done channel for reactivation
	m.done = make(chan struct{})

	m.tracer = otel.Tracer("mission3")

	slog.Info("detailed request tracing enabled", "component", "tracing", "trace_interval_ms", 200)

	// Start background goroutine to emit detailed traces every 200ms
	m.ticker = time.NewTicker(200 * time.Millisecond)
	go m.emitDetailedTraces()

	return nil
}

func (m *Mission3) OnDeactivate(ctx context.Context) error {
	if m.ticker != nil {
		m.ticker.Stop()
	}
	close(m.done)

	slog.Info("detailed request tracing disabled", "component", "tracing")
	return nil
}

func (m *Mission3) emitDetailedTraces() {
	for {
		select {
		case <-m.ticker.C:
			// Create a root span for intel processing operation
			ctx := context.Background()
			ctx, rootSpan := m.tracer.Start(ctx, "process_intel_report")
			rootSpan.SetAttributes(attribute.String("mission", "mission3"))

			// Phase 1: Validation (with nested operations)
			ctx, validateSpan := m.tracer.Start(ctx, "validate_credentials")
			validateSpan.SetAttributes(attribute.String("phase", "authentication"))
			{
				_, checkAuthSpan := m.tracer.Start(ctx, "check_authorization")
				time.Sleep(time.Duration(10+rand.Intn(20)) * time.Millisecond)
				checkAuthSpan.End()

				_, verifyTokenSpan := m.tracer.Start(ctx, "verify_token")
				time.Sleep(time.Duration(5+rand.Intn(15)) * time.Millisecond)
				verifyTokenSpan.End()
			}
			validateSpan.End()

			// Phase 2: Decryption (with multiple sub-steps)
			ctx, decryptSpan := m.tracer.Start(ctx, "decrypt_payload")
			decryptSpan.SetAttributes(attribute.String("phase", "decryption"))
			{
				_, fetchKeySpan := m.tracer.Start(ctx, "fetch_encryption_key")
				time.Sleep(time.Duration(15+rand.Intn(25)) * time.Millisecond)
				fetchKeySpan.End()

				_, decipherSpan := m.tracer.Start(ctx, "decipher_content")
				{
					_, decompressSpan := m.tracer.Start(ctx, "decompress_data")
					time.Sleep(time.Duration(8+rand.Intn(12)) * time.Millisecond)
					decompressSpan.End()

					_, parseSpan := m.tracer.Start(ctx, "parse_structure")
					time.Sleep(time.Duration(5+rand.Intn(10)) * time.Millisecond)
					parseSpan.End()
				}
				decipherSpan.End()
			}
			decryptSpan.End()

			// Phase 3: Security checks
			ctx, securitySpan := m.tracer.Start(ctx, "verify_signature")
			securitySpan.SetAttributes(attribute.String("phase", "security"))
			time.Sleep(time.Duration(12+rand.Intn(18)) * time.Millisecond)
			securitySpan.End()

			// Phase 4: Authorization and database operations
			ctx, permSpan := m.tracer.Start(ctx, "check_permissions")
			permSpan.SetAttributes(attribute.String("phase", "authorization"))
			{
				_, roleCheckSpan := m.tracer.Start(ctx, "check_role_access")
				time.Sleep(time.Duration(6+rand.Intn(10)) * time.Millisecond)
				roleCheckSpan.End()

				_, policySpan := m.tracer.Start(ctx, "validate_policy")
				time.Sleep(time.Duration(8+rand.Intn(12)) * time.Millisecond)
				policySpan.End()
			}
			permSpan.End()

			// Phase 5: Database update with transaction
			ctx, dbSpan := m.tracer.Start(ctx, "update_database")
			dbSpan.SetAttributes(attribute.String("phase", "persistence"))
			{
				_, beginTxSpan := m.tracer.Start(ctx, "begin_transaction")
				time.Sleep(time.Duration(3+rand.Intn(7)) * time.Millisecond)
				beginTxSpan.End()

				_, writeSpan := m.tracer.Start(ctx, "write_intel_record")
				time.Sleep(time.Duration(10+rand.Intn(20)) * time.Millisecond)
				writeSpan.End()

				_, commitSpan := m.tracer.Start(ctx, "commit_transaction")
				time.Sleep(time.Duration(5+rand.Intn(10)) * time.Millisecond)
				commitSpan.End()
			}
			dbSpan.End()

			// Phase 6: Notifications
			ctx, notifySpan := m.tracer.Start(ctx, "send_notification")
			notifySpan.SetAttributes(attribute.String("phase", "notification"))
			{
				_, prepareSpan := m.tracer.Start(ctx, "prepare_notification")
				time.Sleep(time.Duration(5+rand.Intn(10)) * time.Millisecond)
				prepareSpan.End()

				_, sendSpan := m.tracer.Start(ctx, "send_to_queue")
				time.Sleep(time.Duration(8+rand.Intn(15)) * time.Millisecond)
				sendSpan.End()
			}
			notifySpan.End()

			// Phase 7: Audit log
			ctx, auditSpan := m.tracer.Start(ctx, "audit_log")
			auditSpan.SetAttributes(attribute.String("phase", "audit"))
			time.Sleep(time.Duration(7+rand.Intn(13)) * time.Millisecond)
			auditSpan.End()

			rootSpan.End()

		case <-m.done:
			return
		}
	}
}
