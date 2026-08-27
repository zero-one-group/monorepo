// Package webhook is the replay harness every inbound webhook receiver
// (Xendit, Moka, WhatsApp) must pass. Fixtures are real payloads captured from
// the provider's test/trial environment and stored under apps/api/testdata/webhooks.
//
// The two assertions here encode the deck's receiver requirements:
//   - idempotent by event id: the same event twice is a no-op
//   - out-of-order tolerant: a late event never overwrites a terminal state
package webhook

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

// Outcome is the immutable processing result recorded on webhook_events.
type Outcome string

const (
	OutcomeProcessed Outcome = "PROCESSED" // applied to our ledger
	OutcomeIgnored   Outcome = "IGNORED"   // duplicate, late, or not for us — persisted, not applied
	OutcomeRejected  Outcome = "REJECTED"  // bad signature / malformed — persisted, not applied
)

// Result is what a receiver reports after handling one raw payload.
type Result struct {
	ProviderEventID string  // unique per provider; the dedupe key on webhook_events
	Outcome         Outcome // final, immutable
}

// Receiver is the minimal contract a provider receiver implements so it can be
// driven by this harness. Production receivers wrap this behind an Echo handler.
type Receiver interface {
	// Handle verifies, persists and processes one raw payload. It must be safe
	// to call with the same payload more than once.
	Handle(ctx context.Context, headers map[string]string, payload []byte) (Result, error)
}

// Fixture is one captured payload plus the headers the provider sent with it
// (signature headers matter — receivers verify them before anything else).
type Fixture struct {
	Name    string
	Headers map[string]string
	Payload []byte
}

// Load reads apps/api/testdata/webhooks/<provider>/<name>.json (and an optional
// <name>.headers.json) regardless of the package the test runs from.
func Load(t *testing.T, provider, name string) Fixture {
	t.Helper()
	dir := filepath.Join(moduleRoot(t), "testdata", "webhooks", provider)
	payload, err := os.ReadFile(filepath.Join(dir, name+".json"))
	require.NoError(t, err, "fixture %s/%s.json — record it from the provider's test environment, never hand-write it", provider, name)
	fx := Fixture{Name: name, Payload: payload, Headers: map[string]string{}}
	if hb, err := os.ReadFile(filepath.Join(dir, name+".headers.json")); err == nil {
		fx.Headers = parseHeaders(t, hb)
	}
	return fx
}

// AssertIdempotent delivers the same fixture twice: the first delivery must be
// PROCESSED, the second IGNORED, and both must report the same provider event id.
func AssertIdempotent(t *testing.T, r Receiver, fx Fixture) {
	t.Helper()
	ctx := context.Background()
	first, err := r.Handle(ctx, fx.Headers, fx.Payload)
	require.NoError(t, err)
	require.Equal(t, OutcomeProcessed, first.Outcome, "first delivery of %s", fx.Name)
	require.NotEmpty(t, first.ProviderEventID)

	second, err := r.Handle(ctx, fx.Headers, fx.Payload)
	require.NoError(t, err, "a replayed event must not error — the provider will retry on non-2xx")
	require.Equal(t, OutcomeIgnored, second.Outcome, "second delivery of %s must be a no-op", fx.Name)
	require.Equal(t, first.ProviderEventID, second.ProviderEventID)
}

// AssertLateEventIgnored delivers a sequence that ends in a terminal state, then a
// late event that arrived out of order. The late event must be IGNORED, not applied,
// and must not error (we still owe the provider a 2xx).
func AssertLateEventIgnored(t *testing.T, r Receiver, sequence []Fixture, late Fixture) {
	t.Helper()
	ctx := context.Background()
	for _, fx := range sequence {
		res, err := r.Handle(ctx, fx.Headers, fx.Payload)
		require.NoError(t, err, "sequence step %s", fx.Name)
		require.Equal(t, OutcomeProcessed, res.Outcome, "sequence step %s", fx.Name)
	}
	res, err := r.Handle(ctx, late.Headers, late.Payload)
	require.NoError(t, err, "late event %s must still be acknowledged", late.Name)
	require.Equal(t, OutcomeIgnored, res.Outcome, "late event %s must not overwrite a terminal state", late.Name)
}

// AssertRejectsBadSignature tampers with the payload and expects REJECTED without error
// (persist the attempt, return 2xx or 4xx per provider contract — but never process it).
func AssertRejectsBadSignature(t *testing.T, r Receiver, fx Fixture) {
	t.Helper()
	tampered := append([]byte{}, fx.Payload...)
	tampered = append(tampered, ' ')
	res, err := r.Handle(context.Background(), fx.Headers, tampered)
	require.NoError(t, err)
	require.Equal(t, OutcomeRejected, res.Outcome, "tampered %s must be rejected", fx.Name)
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok)
	// pkg/testutils/webhook/replay.go → apps/api
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", ".."))
}
