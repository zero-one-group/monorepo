package middleware

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"{{ package_name | kebab_case }}/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func okHandler(c echo.Context) error { return c.String(http.StatusOK, "ok") }

func serve(e *echo.Echo, req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

// --- rate limit -------------------------------------------------------------

func TestRateLimit_AllowsBurstThenRejects(t *testing.T) {
	e := echo.New()
	e.Use(RateLimitMiddleware(1, 3)) // 1 req/s refill, burst of 3
	e.GET("/", okHandler)

	for i := 1; i <= 3; i++ {
		rec := serve(e, httptest.NewRequest(http.MethodGet, "/", nil))
		assert.Equal(t, http.StatusOK, rec.Code, "request %d within burst", i)
	}
	rec := serve(e, httptest.NewRequest(http.MethodGet, "/", nil))
	assert.Equal(t, http.StatusTooManyRequests, rec.Code, "4th request in the same second is rejected")
}

func TestRateLimit_IsPerClientIP(t *testing.T) {
	e := echo.New()
	e.Use(RateLimitMiddleware(1, 1))
	e.GET("/", okHandler)

	a := httptest.NewRequest(http.MethodGet, "/", nil)
	a.Header.Set(echo.HeaderXRealIP, "10.0.0.1")
	b := httptest.NewRequest(http.MethodGet, "/", nil)
	b.Header.Set(echo.HeaderXRealIP, "10.0.0.2")

	assert.Equal(t, http.StatusOK, serve(e, a).Code)
	assert.Equal(t, http.StatusTooManyRequests, serve(e, a).Code, "same ip exhausted")
	assert.Equal(t, http.StatusOK, serve(e, b).Code, "another ip has its own bucket")
}

func TestRateLimiterStore_ReusesLimiterPerIP(t *testing.T) {
	store := &RateLimiterStore{visitors: map[string]*rateLimiter{}}
	first := store.getRateLimiter("1.1.1.1", 1, 1)
	second := store.getRateLimiter("1.1.1.1", 1, 1)
	assert.Same(t, first, second)
	assert.NotSame(t, first, store.getRateLimiter("2.2.2.2", 1, 1))
}

// --- request id -------------------------------------------------------------

func TestRequestID_GeneratedWhenMissingOrInvalid(t *testing.T) {
	e := echo.New()
	e.Use(RequestIDMiddleware())
	var fromCtx, fromEcho string
	e.GET("/", func(c echo.Context) error {
		fromCtx = GetRequestID(c.Request().Context())
		fromEcho = GetRequestIDFromEcho(c)
		return okHandler(c)
	})

	for _, incoming := range []string{"", "not-a-typeid"} {
		rec := serve(e, withHeader(httptest.NewRequest(http.MethodGet, "/", nil), RequestIDHeader, incoming))
		got := rec.Header().Get(RequestIDHeader)
		assert.True(t, strings.HasPrefix(got, "req_"), "generated id has the req prefix, got %q", got)
		assert.Equal(t, got, fromCtx, "same id on the request context")
		assert.Equal(t, got, fromEcho, "same id on the echo context")
	}
}

func TestRequestID_PreservesValidIncomingID(t *testing.T) {
	e := echo.New()
	e.Use(RequestIDMiddleware())
	e.GET("/", okHandler)
	valid := generateTypeID()
	rec := serve(e, withHeader(httptest.NewRequest(http.MethodGet, "/", nil), RequestIDHeader, valid))
	assert.Equal(t, valid, rec.Header().Get(RequestIDHeader))
}

func TestRequestID_HelpersWithoutMiddleware(t *testing.T) {
	assert.Equal(t, "", GetRequestID(context.Background()))
	assert.NotNil(t, LogWithRequestID(context.Background()))
	ctx := context.WithValue(context.Background(), requestIDCtxKey, "req_x")
	assert.Equal(t, "req_x", GetRequestID(ctx))
	assert.NotNil(t, LogWithRequestID(ctx))
	e := echo.New()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
	assert.Equal(t, "", GetRequestIDFromEcho(c))
}

// --- security headers -------------------------------------------------------

func TestSecurityHeaders(t *testing.T) {
	e := echo.New()
	e.Use(SecurityHeadersMiddleware())
	e.GET("/", okHandler)
	rec := serve(e, httptest.NewRequest(http.MethodGet, "/", nil))
	assert.Equal(t, "nosniff", rec.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", rec.Header().Get("X-Frame-Options"))
	assert.Contains(t, rec.Header().Get("Strict-Transport-Security"), "max-age=31536000")
	assert.Contains(t, rec.Header().Get("Content-Security-Policy"), "frame-ancestors 'none'")
}

// --- timeout ----------------------------------------------------------------

func TestTimeout(t *testing.T) {
	e := echo.New()
	e.Use(TimeoutMiddleware(30 * time.Millisecond))
	e.GET("/fast", okHandler)
	e.GET("/slow", func(c echo.Context) error {
		select {
		case <-c.Request().Context().Done():
		case <-time.After(time.Second):
		}
		return nil
	})
	assert.Equal(t, http.StatusOK, serve(e, httptest.NewRequest(http.MethodGet, "/fast", nil)).Code)
	assert.Equal(t, http.StatusRequestTimeout, serve(e, httptest.NewRequest(http.MethodGet, "/slow", nil)).Code)
}

// --- logger -----------------------------------------------------------------

func TestLogger_LogsEveryRequestWithLevelByStatus(t *testing.T) {
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buf, nil)))
	t.Cleanup(func() { slog.SetDefault(prev) })

	e := echo.New()
	e.Use(RequestIDMiddleware(), LoggerMiddleware(slog.Default()))
	e.GET("/ok", okHandler)
	e.GET("/bad", func(c echo.Context) error { return c.String(http.StatusBadRequest, "bad") })
	e.GET("/boom", func(c echo.Context) error { return c.String(http.StatusInternalServerError, "boom") })

	serve(e, httptest.NewRequest(http.MethodGet, "/ok?x=1", nil))
	serve(e, httptest.NewRequest(http.MethodGet, "/bad", nil))
	serve(e, httptest.NewRequest(http.MethodGet, "/boom", nil))

	out := buf.String()
	assert.Contains(t, out, `"level":"INFO"`)
	assert.Contains(t, out, `"level":"WARN"`)
	assert.Contains(t, out, `"level":"ERROR"`)
	assert.Contains(t, out, `"query":"x=1"`)
	assert.Contains(t, out, `"request_id":"req_`, "request id is on every log line")
}

func TestFormatLatency(t *testing.T) {
	assert.Equal(t, "500ns", formatLatency(500*time.Nanosecond))
	assert.Equal(t, "1.50µs", formatLatency(1500*time.Nanosecond))
	assert.Equal(t, "2.00ms", formatLatency(2*time.Millisecond))
	assert.Equal(t, "1.50s", formatLatency(1500*time.Millisecond))
}

// --- cors + compression -----------------------------------------------------

func TestCORS_UsesConfiguredOrigins(t *testing.T) {
	cfg := &config.Config{}
	cfg.App.CORSOrigins = []string{"https://app.example.com"}
	cfg.App.CORSCredentials = true
	cfg.App.CORSMaxAge = 300

	e := echo.New()
	e.Use(CORSMiddleware(cfg))
	e.GET("/", okHandler)

	allowed := httptest.NewRequest(http.MethodOptions, "/", nil)
	allowed.Header.Set(echo.HeaderOrigin, "https://app.example.com")
	allowed.Header.Set(echo.HeaderAccessControlRequestMethod, http.MethodGet)
	rec := serve(e, allowed)
	assert.Equal(t, "https://app.example.com", rec.Header().Get(echo.HeaderAccessControlAllowOrigin))
	assert.Equal(t, "true", rec.Header().Get(echo.HeaderAccessControlAllowCredentials))

	actual := httptest.NewRequest(http.MethodGet, "/", nil)
	actual.Header.Set(echo.HeaderOrigin, "https://app.example.com")
	rec = serve(e, actual)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get(echo.HeaderAccessControlExposeHeaders), "X-App-Audience")

	denied := httptest.NewRequest(http.MethodGet, "/", nil)
	denied.Header.Set(echo.HeaderOrigin, "https://evil.example")
	rec = serve(e, denied)
	assert.Empty(t, rec.Header().Get(echo.HeaderAccessControlAllowOrigin))
}

func TestCompression_GzipsWhenAccepted(t *testing.T) {
	e := echo.New()
	e.Use(CompressionMiddleware())
	e.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, strings.Repeat("hello ", 200)) })

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderAcceptEncoding, "gzip")
	rec := serve(e, req)
	require.Equal(t, "gzip", rec.Header().Get(echo.HeaderContentEncoding))
	zr, err := gzip.NewReader(rec.Body)
	require.NoError(t, err)
	plain, _ := io.ReadAll(zr)
	assert.Contains(t, string(plain), "hello")
}

func withHeader(r *http.Request, k, v string) *http.Request {
	if v != "" {
		r.Header.Set(k, v)
	}
	return r
}
