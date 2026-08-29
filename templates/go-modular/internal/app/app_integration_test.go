package app

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"{{ package_name | kebab_case }}/internal/config"
	"{{ package_name | kebab_case }}/internal/server"
	modAuth "{{ package_name | kebab_case }}/modules/auth"
	modUser "{{ package_name | kebab_case }}/modules/user"
	"{{ package_name | kebab_case }}/pkg/testutils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

// The one end-to-end wiring test: real Postgres (testcontainers) → migrations + seed →
// the fx graph exactly as New() builds it (minus the listener) → healthz, seeded admin
// sign-in, JWT-guarded route. Everything the unit tests fake is real here. Copy this
// shape for every new module's wiring test.
func TestApp_EndToEnd(t *testing.T) {
	te := testutils.NewTestEnv(t)
	pool, pgURL, err := te.SetupPostgres()
	require.NoError(t, err)
	t.Cleanup(pool.Close)
	te.SetupConfig()
	te.RunAppMigrations()

	t.Setenv("DATABASE_URL", pgURL)
	t.Setenv("ENABLE_API_DOCS", "true")
	cfg, err := config.Load("")
	require.NoError(t, err)

	var e *echo.Echo
	app := fxtest.New(t,
		fx.NopLogger,
		fx.Supply(cfg, slog.New(slog.NewTextHandler(io.Discard, nil))),
		postgresModule,
		mailerModule,
		// httpModule without runHTTPServer: the test drives echo directly.
		fx.Provide(newEcho, newAPIV1, server.NewServerHandler),
		fx.Invoke(registerServerRoutes),
		modUser.Module,
		modAuth.Module,
		fx.Populate(&e),
	)
	app.RequireStart()
	t.Cleanup(app.RequireStop)

	do := func(method, path, body, bearer string) (*httptest.ResponseRecorder, map[string]any) {
		req := httptest.NewRequest(method, path, strings.NewReader(body))
		if body != "" {
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		}
		if bearer != "" {
			req.Header.Set(echo.HeaderAuthorization, "Bearer "+bearer)
		}
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		var out map[string]any
		_ = json.Unmarshal(rec.Body.Bytes(), &out)
		return rec, out
	}

	t.Run("healthz reports the database up", func(t *testing.T) {
		rec, body := do(http.MethodGet, "/healthz", "", "")
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "up", body["status"])
	})

	t.Run("protected user routes require a token", func(t *testing.T) {
		rec, _ := do(http.MethodGet, "/api/v1/users", "", "")
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	var token string
	t.Run("seeded admin can sign in", func(t *testing.T) {
		rec, body := do(http.MethodPost, "/api/v1/auth/signin/username", `{"username":"admin","password":"secure.password"}`, "")
		require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
		token, _ = body["access_token"].(string)
		require.NotEmpty(t, token)
		assert.NotEmpty(t, body["refresh_token"])
		assert.NotEmpty(t, body["session_id"])
	})

	t.Run("wrong password is a 401, not a 500", func(t *testing.T) {
		rec, _ := do(http.MethodPost, "/api/v1/auth/signin/username", `{"username":"admin","password":"nope"}`, "")
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("the access token opens the user routes", func(t *testing.T) {
		rec, _ := do(http.MethodGet, "/api/v1/users", "", token)
		assert.Equal(t, http.StatusOK, rec.Code)
		var users []map[string]any
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &users))
		assert.GreaterOrEqual(t, len(users), 2, "seeded admin + johndoe")
	})

	t.Run("global middleware is wired: cors, gzip", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		req.Header.Set(echo.HeaderOrigin, "https://anything.example")
		req.Header.Set(echo.HeaderAcceptEncoding, "gzip")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.NotEmpty(t, rec.Header().Get(echo.HeaderAccessControlAllowOrigin), "CORS middleware answered (origins come from config/.env)")
		assert.Equal(t, "gzip", rec.Header().Get(echo.HeaderContentEncoding))
	})
}

func TestConnectPostgresWithRetry_GivesUpAfterMaxRetries(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Database.PgMaxRetries = 1 // no sleep between attempts when there is only one
	cfg.Database.PostgresURL = "postgres://nobody:nothing@127.0.0.1:1/none?sslmode=disable&connect_timeout=1"

	start := time.Now()
	pg, err := connectPostgresWithRetry(&cfg, slog.New(slog.NewTextHandler(io.Discard, nil)))
	assert.Nil(t, pg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "after 1 attempts")
	assert.Less(t, time.Since(start), 15*time.Second)
	_ = context.Background()
}
