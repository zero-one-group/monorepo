package rest_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-standard/database"
	"go-standard/internal/metrics"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

// TestKit holds the pieces needed for an E2E test:
//   - Echo    – your router instance (register routes on this before Start())
//   - DB      – real Postgres connection pool (driven from DATABASE_URL)
//   - Metrics – real or stub metrics collector
//
// After calling Start, Server is listening and BaseURL is set.
type TestKit struct {
	Echo    *echo.Echo
	DB      *pgxpool.Pool
	Metrics *metrics.Metrics

	server  *httptest.Server
	BaseURL string
}

// loadEnv loads environment variables from a .env file located two levels up.
// It fatals the test if loading fails.
func loadEnv(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatal("failed to load .env (DATABASE_URL must be set): ", err)
	}
}

// NewTestKit initializes and returns a TestKit containing:
//  1. An Echo router (no routes registered yet)
//  2. A Postgres connection pool (using DATABASE_URL from env/.env)
//  3. A new Metrics instance
//
// It also registers a t.Cleanup to close the DB pool when the test ends.
// Note: you must register your handlers on kit.Echo before calling kit.Start().
func NewTestKit(t *testing.T) *TestKit {
	// 1) ensure DATABASE_URL is loaded
	loadEnv(t)

	// 2) new Echo instance
	e := echo.New()
	e.HideBanner = true

	// 3) setup Postgres pool
	dbPool, err := database.SetupPgxPool()
	require.NoError(t, err, "failed to connect to Postgres via DATABASE_URL")

	// 4) metrics collector
	m := metrics.NewMetrics()

	// 5) ensure we close the pool after test
	t.Cleanup(func() {
		dbPool.Close()
	})

	return &TestKit{
		Echo:    e,
		DB:      dbPool,
		Metrics: m,
	}
}

// Start takes the Echo router (with your routes already registered) and
// spins up an httptest.Server.  It sets kit.BaseURL to the server URL
// (e.g. "http://127.0.0.1:XXXXX") and registers a t.Cleanup to close it.
func (kit *TestKit) Start(t *testing.T) {
	kit.server = httptest.NewServer(kit.Echo)
	kit.BaseURL = kit.server.URL

	// ensure server is closed when test finishes
	t.Cleanup(func() {
		kit.server.Close()
	})
}

// doRequest sends an HTTP request to the given URL using the specified method.
// If payload != nil, it is JSON-marshaled and sent as the request body with
// “Content-Type: application/json”.
//
// The response body is read in full and unmarshaled into the generic type T.
// Returns the T instance and the HTTP status code.
//
// The test is fatally failed if any step errors.
func doRequest[T any](t *testing.T, method, url string, payload any) (T, int) {
	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		require.NoError(t, err, "failed to marshal payload")
		body = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, url, body)
	require.NoError(t, err, "failed to build HTTP request")

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err, "HTTP request failed")
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "failed to read response body")

	var out T
	require.NoError(t, json.Unmarshal(raw, &out),
		"failed to unmarshal response into %T, body=%s", out, string(raw),
	)
	return out, resp.StatusCode
}
