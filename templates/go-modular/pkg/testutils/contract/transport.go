package contract

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

// replayTransport answers every request from a fixture named after the request:
// "<METHOD> <path>" → "<method>_<path-with-slashes-as-dashes>.json".
// Tests that need a different mapping call Load directly instead of Client.
type replayTransport struct {
	t        *testing.T
	provider string
}

func (rt *replayTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	name := FixtureNameFor(req.Method, req.URL.Path)
	body := Load(rt.t, rt.provider, name)
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewReader(body)),
		Request:    req,
	}, nil
}

// FixtureNameFor maps an HTTP call to a fixture file name: GET /v3/outlets → get_v3-outlets.
func FixtureNameFor(method, path string) string {
	p := strings.Trim(path, "/")
	p = strings.ReplaceAll(p, "/", "-")
	if p == "" {
		p = "root"
	}
	return strings.ToLower(method) + "_" + p
}
