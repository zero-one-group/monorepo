// Package contract supports record/replay contract tests against third-party
// APIs (Moka trial account, Xendit test mode).
//
// Default mode replays fixtures from apps/api/testdata/contract/<provider>/ so the
// suite is fast and offline. Set CONTRACT_RECORD=1 to hit the real test environment
// and re-record — a deliberate act, done when the provider changes, never in CI.
//
// Contract tests live behind the `contract` build tag:
//
//	//go:build contract
//
// and run with `moon run api:test-contract`.
package contract

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

// Recording is set when CONTRACT_RECORD=1: tests must call the real API and
// then Save() the response so the next replay run sees it.
func Recording() bool { return os.Getenv("CONTRACT_RECORD") == "1" }

// Client returns the HTTP client contract tests should use, plus the base URL.
// In replay mode the transport serves fixtures instead of touching the network.
func Client(t *testing.T, provider, baseURL string) (*http.Client, string) {
	t.Helper()
	if Recording() {
		return &http.Client{}, baseURL
	}
	return &http.Client{Transport: &replayTransport{t: t, provider: provider}}, baseURL
}

// Save writes a recorded response body as the fixture for name. No-op in replay mode.
func Save(t *testing.T, provider, name string, body []byte) {
	t.Helper()
	if !Recording() {
		return
	}
	path := fixturePath(t, provider, name)
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	require.NoError(t, os.WriteFile(path, body, 0o644))
	t.Logf("recorded %s", path)
}

// Load returns a previously recorded fixture; it fails with guidance when missing.
func Load(t *testing.T, provider, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(fixturePath(t, provider, name))
	require.NoError(t, err, "no fixture for %s/%s — run with CONTRACT_RECORD=1 against the %s test environment to record it", provider, name, provider)
	return b
}

// RequireCredential skips the test in record mode when a credential is absent,
// so a missing MOKA_API_TOKEN produces a clear skip instead of a confusing 401.
func RequireCredential(t *testing.T, env string) string {
	t.Helper()
	v := os.Getenv(env)
	if Recording() && v == "" {
		t.Skipf("%s not set; cannot record contract fixtures", env)
	}
	return v
}

func fixturePath(t *testing.T, provider, name string) string {
	t.Helper()
	return filepath.Join(moduleRoot(t), "testdata", "contract", provider, name+".json")
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok)
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", ".."))
}
