package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/grafana/alloy-mission-control/internal/config"
	"github.com/grafana/alloy-mission-control/internal/missioncontrol"
	"github.com/grafana/alloy-mission-control/internal/telemetry"
)

// KeyChunkProvider provides key chunk data for Mission 4 error trace injection
type KeyChunkProvider interface {
	GetNextKeyChunk() (chunk string, sequence int, totalChunks int)
}

type OperationsHandler struct {
	metrics    *telemetry.Metrics
	tracer     trace.Tracer
	controller *missioncontrol.Controller
	config     *config.Config
	keyChunks  KeyChunkProvider
}

func NewOperationsHandler(metrics *telemetry.Metrics, tracer trace.Tracer, controller *missioncontrol.Controller, cfg *config.Config, keyChunks KeyChunkProvider) *OperationsHandler {
	return &OperationsHandler{
		metrics:    metrics,
		tracer:     tracer,
		controller: controller,
		config:     cfg,
		keyChunks:  keyChunks,
	}
}

// RetinaScanVerify handles biometric verification requests
func (h *OperationsHandler) RetinaScanVerify(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := h.tracer.Start(ctx, "biometrics_retina_scan")
	defer span.End()

	scanID := uuid.New().String()

	// Nested span: Capture image
	ctx, captureSpan := h.tracer.Start(ctx, "capture_retina_image")
	time.Sleep(h.randomDuration())
	captureSpan.End()

	// Nested span: Process image with sub-operations
	ctx, processSpan := h.tracer.Start(ctx, "process_image")
	{
		_, enhanceSpan := h.tracer.Start(ctx, "enhance_image_quality")
		time.Sleep(h.randomDuration())
		enhanceSpan.End()

		_, extractSpan := h.tracer.Start(ctx, "extract_pattern_features")
		time.Sleep(h.randomDuration())
		extractSpan.End()
	}
	processSpan.End()

	// Nested span: Database comparison
	ctx, dbSpan := h.tracer.Start(ctx, "database_pattern_match")
	time.Sleep(h.randomDuration())
	dbSpan.End()

	// Nested span: Validate and audit
	ctx, validateSpan := h.tracer.Start(ctx, "validate_and_audit")
	time.Sleep(h.randomDuration())
	validateSpan.End()

	h.metrics.IntelReportsProcessed.With(prometheus.Labels{
		"endpoint": "retina_scan",
	}).Inc()

	slog.Info("retina scan verified",
		"component", "operations",
		"scan_id", scanID,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "verified",
		"scan_id":    scanID,
		"match":      true,
		"confidence": 0.98,
	})
}

// HeliPickupRequest handles extraction requests
func (h *OperationsHandler) HeliPickupRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := h.tracer.Start(ctx, "extraction_heli_pickup")
	defer span.End()

	requestID := uuid.New().String()

	// Nested span: Validate extraction zone
	ctx, validateSpan := h.tracer.Start(ctx, "validate_extraction_zone")
	time.Sleep(h.randomDuration())
	validateSpan.End()

	// Nested span: Check airspace clearance
	ctx, airspaceSpan := h.tracer.Start(ctx, "check_airspace_clearance")
	{
		_, radarSpan := h.tracer.Start(ctx, "radar_sweep")
		time.Sleep(h.randomDuration())
		radarSpan.End()

		_, weatherSpan := h.tracer.Start(ctx, "weather_check")
		time.Sleep(h.randomDuration())
		weatherSpan.End()
	}
	airspaceSpan.End()

	// Nested span: Dispatch helicopter
	ctx, dispatchSpan := h.tracer.Start(ctx, "dispatch_helicopter")
	{
		_, pilotSpan := h.tracer.Start(ctx, "notify_pilot")
		time.Sleep(h.randomDuration())
		pilotSpan.End()

		_, routeSpan := h.tracer.Start(ctx, "calculate_route")
		time.Sleep(h.randomDuration())
		routeSpan.End()
	}
	dispatchSpan.End()

	h.metrics.IntelReportsProcessed.With(prometheus.Labels{
		"endpoint": "heli_pickup",
	}).Inc()

	slog.Info("heli pickup requested", "component", "operations", "request_id", requestID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "dispatched",
		"request_id":  requestID,
		"eta_minutes": 15,
	})
}

// BurnAfterReading handles secure message retrieval
func (h *OperationsHandler) BurnAfterReading(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := h.tracer.Start(ctx, "message_burn_after_reading")
	defer span.End()

	messageID := uuid.New().String()

	// Nested span: Authenticate request
	ctx, authSpan := h.tracer.Start(ctx, "authenticate_request")
	time.Sleep(h.randomDuration())
	authSpan.End()

	// Nested span: Retrieve encrypted message
	ctx, retrieveSpan := h.tracer.Start(ctx, "retrieve_encrypted_message")
	time.Sleep(h.randomDuration())
	retrieveSpan.End()

	// Nested span: Decrypt message
	ctx, decryptSpan := h.tracer.Start(ctx, "decrypt_message")
	{
		_, keySpan := h.tracer.Start(ctx, "fetch_encryption_key")
		time.Sleep(h.randomDuration())
		keySpan.End()

		_, decipherSpan := h.tracer.Start(ctx, "decipher_content")
		time.Sleep(h.randomDuration())
		decipherSpan.End()
	}
	decryptSpan.End()

	// Nested span: Secure deletion
	ctx, burnSpan := h.tracer.Start(ctx, "secure_delete")
	time.Sleep(h.randomDuration())
	burnSpan.End()

	h.metrics.IntelReportsProcessed.With(prometheus.Labels{
		"endpoint": "burn_message",
	}).Inc()

	slog.Info("secure message retrieved", "component", "operations", "message_id", messageID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "retrieved",
		"message_id": messageID,
		"burned":     true,
	})
}

// DeadDropUpload handles dead drop communication uploads
func (h *OperationsHandler) DeadDropUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := h.tracer.Start(ctx, "comms_dead_drop_upload")
	defer span.End()

	dropID := uuid.New().String()

	// Nested span: Verify drop location
	ctx, verifySpan := h.tracer.Start(ctx, "verify_drop_location")
	{
		_, coordSpan := h.tracer.Start(ctx, "validate_coordinates")
		time.Sleep(h.randomDuration())
		coordSpan.End()

		_, secureSpan := h.tracer.Start(ctx, "check_location_security")
		time.Sleep(h.randomDuration())
		secureSpan.End()
	}
	verifySpan.End()

	// Nested span: Encrypt payload
	ctx, encryptSpan := h.tracer.Start(ctx, "encrypt_payload")
	time.Sleep(h.randomDuration())
	encryptSpan.End()

	// Nested span: Upload to dead drop
	ctx, uploadSpan := h.tracer.Start(ctx, "upload_to_location")
	time.Sleep(h.randomDuration())
	uploadSpan.End()

	// Nested span: Confirm delivery
	ctx, confirmSpan := h.tracer.Start(ctx, "confirm_delivery")
	time.Sleep(h.randomDuration())
	confirmSpan.End()

	h.metrics.IntelReportsProcessed.With(prometheus.Labels{
		"endpoint": "dead_drop",
	}).Inc()

	slog.Info("dead drop uploaded", "component", "operations", "drop_id", dropID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "uploaded",
		"drop_id":  dropID,
		"location": "sector_7",
	})
}

// SatLinkHandshake handles satellite link establishment
func (h *OperationsHandler) SatLinkHandshake(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := h.tracer.Start(ctx, "satlink_handshake")
	defer span.End()

	linkID := uuid.New().String()
	frequencyOverride := r.Header.Get("X-Frequency-Override")

	// Nested span: Acquire satellite
	ctx, acquireSpan := h.tracer.Start(ctx, "acquire_satellite")
	{
		_, scanSpan := h.tracer.Start(ctx, "scan_orbit_positions")
		time.Sleep(h.randomDuration())
		scanSpan.End()

		_, targetSpan := h.tracer.Start(ctx, "target_optimal_satellite")
		time.Sleep(h.randomDuration())
		targetSpan.End()
	}
	acquireSpan.End()

	// Nested span: Establish uplink
	ctx, uplinkSpan := h.tracer.Start(ctx, "establish_uplink")
	{
		_, authSpan := h.tracer.Start(ctx, "authenticate_ground_station")
		time.Sleep(h.randomDuration())
		authSpan.End()

		_, syncSpan := h.tracer.Start(ctx, "sync_frequency")
		time.Sleep(h.randomDuration())
		syncSpan.End()
	}
	uplinkSpan.End()

	// Nested span: Verify link quality
	ctx, verifySpan := h.tracer.Start(ctx, "verify_link_quality")
	time.Sleep(h.randomDuration())
	verifySpan.End()

	// If frequency override header is present and we have a key chunk provider,
	// create an error span with key fragment data
	if frequencyOverride != "" && h.keyChunks != nil {
		chunk, sequence, totalChunks := h.keyChunks.GetNextKeyChunk()

		_, decodeSpan := h.tracer.Start(ctx, "signal_decode")
		decodeSpan.SetAttributes(
			attribute.Int("key_sequence", sequence),
			attribute.String("key_chunk", chunk),
			attribute.Int("total_chunks", totalChunks),
			attribute.String("error.type", "SignalInterference"),
			attribute.String("frequency_override", frequencyOverride),
		)
		decodeSpan.SetStatus(codes.Error, fmt.Sprintf("signal interference on frequency %s", frequencyOverride))
		decodeSpan.RecordError(fmt.Errorf("satellite signal interference detected"))
		decodeSpan.End()

		span.SetStatus(codes.Error, "satellite link interference detected")

		h.metrics.IntelReportsProcessed.With(prometheus.Labels{
			"endpoint": "sat_link",
		}).Inc()

		slog.Info("satellite link interference",
			"component", "operations",
			"link_id", linkID,
			"frequency_override", frequencyOverride)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "error",
			"error":   "satellite link interference detected",
			"link_id": linkID,
		})
		return
	}

	h.metrics.IntelReportsProcessed.With(prometheus.Labels{
		"endpoint": "sat_link",
	}).Inc()

	slog.Info("satellite link established", "component", "operations", "link_id", linkID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":          "established",
		"link_id":         linkID,
		"signal_strength": 0.95,
	})
}

// DeadDropView handles viewing the contents of a dead drop
// Requires authentication using the access token from Mission 4
func (h *OperationsHandler) DeadDropView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := h.tracer.Start(ctx, "deaddrop_view")
	defer span.End()

	// Extract authentication token from Authorization header or query parameter
	token := h.extractAuthToken(r)

	if token == "" {
		slog.Warn("deaddrop access denied - no authentication token provided", "component", "operations")
		h.metrics.IntelReportsProcessed.With(prometheus.Labels{
			"endpoint": "deaddrop_view_denied",
		}).Inc()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "access_denied",
			"message": "Authentication required. Provide access token in Authorization header or 'key' query parameter.",
		})
		return
	}

	// Validate token against the expected access token
	expectedToken := h.config.Mission4SecretMessage
	if token != expectedToken {
		slog.Warn("deaddrop access denied - invalid access token", "component", "operations")
		h.metrics.IntelReportsProcessed.With(prometheus.Labels{
			"endpoint": "deaddrop_view_denied",
		}).Inc()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "access_denied",
			"message": "Invalid access token. Assemble all token fragments from error traces.",
		})
		return
	}

	// Access granted - return the classified intelligence
	dropID := uuid.New().String()

	// Nested span: Verify authorization
	ctx, authSpan := h.tracer.Start(ctx, "verify_access_token")
	time.Sleep(h.randomDuration())
	authSpan.End()

	// Nested span: Retrieve classified data
	ctx, retrieveSpan := h.tracer.Start(ctx, "retrieve_classified_intel")
	time.Sleep(h.randomDuration())
	retrieveSpan.End()

	// Nested span: Decrypt contents
	ctx, decryptSpan := h.tracer.Start(ctx, "decrypt_deaddrop_contents")
	time.Sleep(h.randomDuration())
	decryptSpan.End()

	h.metrics.IntelReportsProcessed.With(prometheus.Labels{
		"endpoint": "deaddrop_view_success",
	}).Inc()

	slog.Info("deaddrop accessed successfully",
		"component", "operations",
		"drop_id", dropID,
		"authenticated", true)

	// Return the classified intelligence
	classifiedIntel := map[string]interface{}{
		"status":   "decrypted",
		"drop_id":  dropID,
		"location": "sector_7_alpha",
		"priority": "CRITICAL",
		"message": map[string]interface{}{
			"subject": "Operation Nightfall - Final Instructions",
			"content": "Congratulations, Agent. You have successfully intercepted and decoded the encrypted transmission. " +
				"Your mission is complete. The intelligence gathered here will be instrumental in stopping the threat. " +
				"Further instructions await at HQ. Mission status: SUCCESS.",
			"classification": "TOP SECRET",
			"expires":        "BURN AFTER READING",
		},
		"attachments": []string{
			"satellite_imagery_sector_7.enc",
			"agent_roster_compromised.enc",
			"extraction_routes.enc",
		},
		"next_steps": []string{
			"Report to HQ for debriefing",
			"Destroy all traces of this communication",
			"Await further orders",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(classifiedIntel)
}

// extractAuthToken extracts the authentication token from the request
// Supports both Authorization header (Bearer token) and query parameter
func (h *OperationsHandler) extractAuthToken(r *http.Request) string {
	// Try Authorization header first (Bearer token)
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		// Support both "Bearer TOKEN" and just "TOKEN"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
		// If no "Bearer" prefix, treat entire header as token
		return authHeader
	}

	// Fall back to query parameter
	return r.URL.Query().Get("key")
}

// randomDuration returns a random duration between 5-50ms
func (h *OperationsHandler) randomDuration() time.Duration {
	min := 5
	max := 50
	ms := rand.Intn(max-min+1) + min
	return time.Duration(ms) * time.Millisecond
}
