package missions

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/grafana/alloy-mission-control/internal/config"
)

// Mission4 implements the Trace Sampling - Part 2 (Tail Sampling) scenario
// A field agent discovered that injecting a specific frequency override into the
// sat-link handshake always causes interference errors. They embed token fragments
// in the error telemetry. Participants must use tail sampling to collect all error
// traces and assemble the access token.
type Mission4 struct {
	accessToken     string
	keyChunks       []string
	numChunks       int
	currentChunkIdx atomic.Int32
	emitInterval    time.Duration
	baseURL         string
	client          *http.Client
	ticker          *time.Ticker
	done            chan struct{}
	mu              sync.Mutex
}

// NewMission4 creates a new Tail Sampling mission
func NewMission4(cfg *config.Config, baseURL string) *Mission4 {
	numChunks := 5
	token := cfg.Mission4SecretMessage
	chunks := splitIntoChunks(token, numChunks)

	return &Mission4{
		accessToken: token,
		keyChunks:     chunks,
		numChunks:     numChunks,
		emitInterval:  15 * time.Second,
		baseURL:       baseURL,
		client:        &http.Client{Timeout: 10 * time.Second},
	}
}

// splitIntoChunks divides a string into N chunks as evenly as possible
func splitIntoChunks(s string, n int) []string {
	if n <= 0 {
		return []string{s}
	}
	if len(s) == 0 {
		return []string{}
	}

	chunks := make([]string, 0, n)
	chunkSize := len(s) / n
	remainder := len(s) % n

	start := 0
	for i := range n {
		size := chunkSize
		if i < remainder {
			size++
		}

		end := min(start+size, len(s))

		chunks = append(chunks, s[start:end])
		start = end
	}

	return chunks
}

func (m *Mission4) ID() string {
	return "mission4"
}

func (m *Mission4) Name() string {
	return "Tail Sampling - Access Token Recovery"
}

func (m *Mission4) Description() string {
	return fmt.Sprintf("Agent transmitting access token in %d pieces via error traces (emits every %v)", len(m.keyChunks), m.emitInterval)
}

func (m *Mission4) OnActivate(ctx context.Context) error {
	m.done = make(chan struct{})
	m.currentChunkIdx.Store(0)

	slog.Info("extended trace diagnostics enabled", "component", "tracing", "mode", "detailed")

	m.ticker = time.NewTicker(m.emitInterval)
	go m.sendInterferenceRequests()

	return nil
}

func (m *Mission4) OnDeactivate(ctx context.Context) error {
	if m.ticker != nil {
		m.ticker.Stop()
	}
	close(m.done)

	slog.Info("extended trace diagnostics disabled", "component", "tracing")
	return nil
}

// GetNextKeyChunk returns the current key chunk, its sequence number (1-based),
// and the total number of chunks. It advances the index for the next call.
func (m *Mission4) GetNextKeyChunk() (chunk string, sequence int, totalChunks int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	idx := int(m.currentChunkIdx.Load())
	chunk = m.keyChunks[idx]
	sequence = idx + 1
	totalChunks = m.numChunks

	nextIdx := (idx + 1) % m.numChunks
	m.currentChunkIdx.Store(int32(nextIdx))

	return chunk, sequence, totalChunks
}

// NumChunks returns the total number of key chunks
func (m *Mission4) NumChunks() int {
	return m.numChunks
}

// AccessToken returns the full access token
func (m *Mission4) AccessToken() string {
	return m.accessToken
}

func (m *Mission4) sendInterferenceRequests() {
	for {
		select {
		case <-m.ticker.C:
			url := m.baseURL + "/api/sat-link/handshake/establish"
			req, err := http.NewRequest("POST", url, bytes.NewReader([]byte("{}")))
			if err != nil {
				slog.Warn("failed to create interference request", "component", "tracing", "error", err)
				continue
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Frequency-Override", "47.3MHz")

			resp, err := m.client.Do(req)
			if err != nil {
				slog.Warn("interference request failed", "component", "tracing", "error", err)
				continue
			}
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()

		case <-m.done:
			return
		}
	}
}
