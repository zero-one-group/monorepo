package server

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"{{ package_name | kebab_case }}/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testLogger() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }

func spaFS(withIndex bool) http.FileSystem {
	m := fstest.MapFS{"app.js": {Data: []byte("console.log('app')")}}
	if withIndex {
		m["index.html"] = &fstest.MapFile{Data: []byte("<!doctype html><title>Example Admin</title>")}
	}
	return http.FS(m)
}

// RootHandler is the embedded-SPA fallback.
func TestRootHandler_ServesIndexForAnyAppRoute(t *testing.T) {
	h := &ServerHandler{Logger: testLogger()}
	e := echo.New()
	e.GET("/", h.RootHandler(spaFS(true)))
	e.GET("/*", h.RootHandler(spaFS(true)))

	for _, path := range []string{"/", "/orders/123", "/deep/link?x=1"} {
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		assert.Equal(t, http.StatusOK, rec.Code, path)
		assert.Contains(t, rec.Header().Get(echo.HeaderContentType), "text/html", path)
		assert.Contains(t, rec.Body.String(), "Example Admin", path)
	}
}

func TestRootHandler_DoesNotShadowStaticOrHealthz(t *testing.T) {
	h := &ServerHandler{Logger: testLogger()}
	e := echo.New()
	e.GET("/*", h.RootHandler(spaFS(true)))
	for _, path := range []string{"/static/missing.js", "/healthz"} {
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		assert.Equal(t, http.StatusNotFound, rec.Code, path)
	}
}

func TestRootHandler_404WhenNoIndexBuilt(t *testing.T) {
	h := &ServerHandler{Logger: testLogger()}
	e := echo.New()
	e.GET("/*", h.RootHandler(spaFS(false)))
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/anything", nil))
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "index.html not found")
}

func TestGetFileSystem_ExposesEmbeddedStaticDir(t *testing.T) {
	h := NewServerHandler(nil, testLogger())
	f, err := getFileSystem(h.WebFS, "static").Open("index.html")
	require.NoError(t, err, "the embedded SPA shell ships with the binary")
	_ = f.Close()
	assert.Panics(t, func() { getFileSystem(h.WebFS, "../escape") }, "fs.Sub rejects invalid paths")
}

// API docs are gated by ENABLE_API_DOCS; the spec is the embedded swagger.json.
func TestAPIDocs_EnabledAndDisabled(t *testing.T) {
	t.Setenv("ENABLE_API_DOCS", "true")
	_, err := config.Load("")
	require.NoError(t, err)

	h := NewServerHandler(nil, testLogger())
	e := echo.New()
	h.RegisterRoutes(e)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/openapi.json", nil))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get(echo.HeaderContentType), "application/json")
	assert.Contains(t, rec.Body.String(), `"swagger"`)

	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api-docs", nil))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "/api/openapi.json", "docs UI points at the served spec")

	t.Setenv("ENABLE_API_DOCS", "false")
	_, err = config.Load("")
	require.NoError(t, err)
	for _, path := range []string{"/api/openapi.json", "/api-docs"} {
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		assert.Equal(t, http.StatusNotFound, rec.Code, path)
	}
}

func TestRegisterRoutes_StaticAndSPARoutesExist(t *testing.T) {
	t.Setenv("ENABLE_API_DOCS", "true")
	_, err := config.Load("")
	require.NoError(t, err)
	h := NewServerHandler(nil, testLogger())
	e := echo.New()
	h.RegisterRoutes(e)

	paths := map[string]bool{}
	for _, r := range e.Routes() {
		paths[r.Method+" "+r.Path] = true
	}
	for _, want := range []string{"GET /", "GET /*", "GET /static/*", "GET /healthz", "GET /api-docs", "GET /api/openapi.json"} {
		assert.True(t, paths[want], "route %s registered", want)
	}

	// The embedded index.html ships with the binary: the SPA shell is served for /.
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	assert.Equal(t, http.StatusOK, rec.Code)
}
