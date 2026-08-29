package contract

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFixtureNameFor(t *testing.T) {
	require.Equal(t, "get_v3-outlets", FixtureNameFor("GET", "/v3/outlets/"))
	require.Equal(t, "post_root", FixtureNameFor("POST", "/"))
}

func TestReplayTransportServesFixture(t *testing.T) {
	t.Setenv("CONTRACT_RECORD", "")
	dir := filepath.Join(moduleRoot(t), "testdata", "contract", "_selftest")
	require.NoError(t, os.MkdirAll(dir, 0o755))
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	require.NoError(t, os.WriteFile(filepath.Join(dir, "get_ping.json"), []byte(`{"ok":true}`), 0o644))

	client, base := Client(t, "_selftest", "https://example.invalid")
	resp, err := client.Get(base + "/ping")
	require.NoError(t, err)
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.JSONEq(t, `{"ok":true}`, string(b))
}

func TestSaveIsNoopUnlessRecording(t *testing.T) {
	t.Setenv("CONTRACT_RECORD", "")
	Save(t, "_selftest", "should-not-exist", []byte("{}"))
	_, err := os.Stat(filepath.Join(moduleRoot(t), "testdata", "contract", "_selftest", "should-not-exist.json"))
	require.True(t, os.IsNotExist(err))
}

func TestRequireCredentialSkipsOnlyWhenRecording(t *testing.T) {
	t.Setenv("CONTRACT_RECORD", "")
	t.Setenv("CONTRACT_TEST_TOKEN", "")
	require.Equal(t, "", RequireCredential(t, "CONTRACT_TEST_TOKEN"))
}
