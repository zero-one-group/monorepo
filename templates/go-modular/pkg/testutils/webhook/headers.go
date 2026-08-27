package webhook

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func parseHeaders(t *testing.T, b []byte) map[string]string {
	t.Helper()
	var h map[string]string
	require.NoError(t, json.Unmarshal(b, &h))
	return h
}
