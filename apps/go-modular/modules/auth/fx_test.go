package auth

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-modular/internal/config"
	appMiddleware "go-modular/internal/middleware"
	"go-modular/internal/server"
	"go-modular/modules/auth/handler"
	"go-modular/modules/auth/repository"
	user_services "go-modular/modules/user/services"
	"go-modular/pkg/apputils"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var secret = []byte("unit-test-secret-do-not-use-in-prod")

func testConfig() *config.Config {
	cfg := config.DefaultConfig()
	cfg.App.JWTSecretKey = string(secret)
	cfg.App.BaseURL = "https://api.example.com"
	return &cfg
}

type noopUserService struct {
	user_services.UserServiceInterface
}

// A zero pool is enough to construct the module: none of the requests in these tests get
// past validation or the JWT guard, so the DB is never touched.
func testHandler(t *testing.T) *handler.Handler {
	t.Helper()
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	repo := repository.NewAuthRepository(&pgxpool.Pool{}, log)
	svc, err := newAuthService(repo, noopUserService{}, nil, testConfig())
	require.NoError(t, err)
	return newAuthHandler(log, svc)
}

func accessToken(t *testing.T, sid string) string {
	t.Helper()
	gen := apputils.NewJWTGenerator(apputils.JWTConfig{SecretKey: secret, AccessTokenExpiry: time.Hour, Issuer: "https://api.example.com"})
	tok, err := gen.Sign(context.Background(), map[string]any{"sid": sid, "email": "jane@example.com"}, uuid.Must(uuid.NewV7()).String())
	require.NoError(t, err)
	return tok
}

func refreshToken(t *testing.T) string {
	t.Helper()
	gen := apputils.NewJWTGenerator(apputils.JWTConfig{SecretKey: secret, RefreshTokenExpiry: time.Hour})
	tok, err := gen.GenerateRefreshTokenJWT(context.Background(), uuid.Must(uuid.NewV7()).String(), "client-app", uuid.Must(uuid.NewV7()).String())
	require.NoError(t, err)
	return tok
}

func TestValidateAuthConfig(t *testing.T) {
	require.NoError(t, validateAuthConfig(testConfig()))

	for name, mutate := range map[string]func(*config.Config){
		"JWTSecretKey is required": func(c *config.Config) { c.App.JWTSecretKey = "  " },
		"BaseURL is required":      func(c *config.Config) { c.App.BaseURL = "" },
	} {
		cfg := testConfig()
		mutate(cfg)
		err := validateAuthConfig(cfg)
		require.Error(t, err, name)
		assert.Contains(t, err.Error(), name)

		_, err = newJWTAccess(cfg)
		assert.Error(t, err, "newJWTAccess re-validates: %s", name)
		_, err = newAuthService(nil, nil, nil, cfg)
		assert.Error(t, err, "newAuthService re-validates: %s", name)
	}
}

func TestJWTSigningAlg(t *testing.T) {
	cfg := testConfig()
	cfg.App.JWTAlgorithm = ""
	assert.Equal(t, jwa.HS256, jwtSigningAlg(cfg), "empty defaults to HS256")
	cfg.App.JWTAlgorithm = " rs256 "
	assert.Equal(t, jwa.RS256, jwtSigningAlg(cfg), "trimmed and upper-cased")
}

func TestRegisterRoutes_PublicVsProtected(t *testing.T) {
	e := echo.New()
	jwt, err := newJWTAccess(testConfig())
	require.NoError(t, err)
	registerRoutes(server.APIV1{Root: e.Group("/api/v1")}, testHandler(t), jwt)

	routes := map[string]bool{}
	for _, r := range e.Routes() {
		routes[r.Method+" "+r.Path] = true
	}
	for _, want := range []string{
		"POST /api/v1/auth/signin/email", "POST /api/v1/auth/signin/username", "GET /api/v1/auth/verify-email",
		"POST /api/v1/auth/refresh-token", "POST /api/v1/auth/password", "GET /api/v1/auth/session/:sessionId",
		"POST /api/v1/auth/verification/email/resend",
	} {
		assert.True(t, routes[want], want)
	}

	// Protected routes reject anonymous callers before touching any service.
	for _, r := range []struct{ method, path string }{
		{http.MethodPost, "/api/v1/auth/password"},
		{http.MethodPost, "/api/v1/auth/session"},
		{http.MethodGet, "/api/v1/auth/session/" + uuid.Must(uuid.NewV7()).String()},
		{http.MethodPost, "/api/v1/auth/verification/email/resend"},
	} {
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, httptest.NewRequest(r.method, r.path, nil))
		assert.Equal(t, http.StatusUnauthorized, rec.Code, "%s %s", r.method, r.path)
	}

	// Public sign-in reaches the handler (400 = validation, not 401).
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/v1/auth/signin/email", nil))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestJWTMiddleware(t *testing.T) {
	e := echo.New()
	var gotUser, gotSID string
	var claimsOK bool
	e.GET("/me", func(c echo.Context) error {
		gotUser, _ = GetUserID(c)
		gotSID, _ = c.Get("session_id").(string)
		_, claimsOK = GetJWTClaims(c)
		_, inCtx := c.Request().Context().Value(apputils.JWTClaimsContextKey).(map[string]any)
		assert.True(t, inCtx, "claims propagated to request context")
		return c.NoContent(http.StatusOK)
	}, JWTMiddleware(secret, jwa.HS256))

	hit := func(authz string) int {
		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		if authz != "" {
			req.Header.Set(echo.HeaderAuthorization, authz)
		}
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		return rec.Code
	}

	sid := uuid.Must(uuid.NewV7()).String()
	assert.Equal(t, http.StatusOK, hit("Bearer "+accessToken(t, sid)))
	assert.NotEmpty(t, gotUser)
	assert.Equal(t, sid, gotSID)
	assert.True(t, claimsOK)

	assert.Equal(t, http.StatusOK, hit("bearer "+accessToken(t, sid)), "scheme is case-insensitive")
	assert.Equal(t, http.StatusUnauthorized, hit(""), "missing header")
	assert.Equal(t, http.StatusUnauthorized, hit("Token abc"), "wrong scheme")
	assert.Equal(t, http.StatusUnauthorized, hit("Bearer"), "no token")
	assert.Equal(t, http.StatusUnauthorized, hit("Bearer not.a.jwt"), "garbage")
	assert.Equal(t, http.StatusUnauthorized, hit("Bearer "+refreshToken(t)), "a refresh token is not an access token")

	other := apputils.NewJWTGenerator(apputils.JWTConfig{SecretKey: []byte("someone-else"), AccessTokenExpiry: time.Hour})
	forged, _ := other.Sign(context.Background(), map[string]any{}, "u")
	assert.Equal(t, http.StatusUnauthorized, hit("Bearer "+forged), "wrong key")

	expired := apputils.NewJWTGenerator(apputils.JWTConfig{SecretKey: secret, AccessTokenExpiry: -time.Minute})
	old, _ := expired.Sign(context.Background(), map[string]any{}, "u")
	assert.Equal(t, http.StatusUnauthorized, hit("Bearer "+old), "expired")
}

func TestJWTAccess_IsTheSameGuard(t *testing.T) {
	// The fx-provided JWTAccess must behave exactly like JWTMiddleware: other modules mount it.
	jwt, err := newJWTAccess(testConfig())
	require.NoError(t, err)
	e := echo.New()
	e.GET("/x", func(c echo.Context) error { return c.NoContent(http.StatusOK) }, echo.MiddlewareFunc(jwt))
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/x", nil))
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+accessToken(t, "s"))
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	var _ appMiddleware.JWTAccess = jwt
}

func TestContextHelpersWithoutMiddleware(t *testing.T) {
	e := echo.New()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
	_, ok := GetJWTClaims(c)
	assert.False(t, ok)
	_, ok = GetUserID(c)
	assert.False(t, ok)

	c.Set("jwt_claims", map[string]any{"sub": "u-1"})
	id, ok := GetUserID(c)
	assert.True(t, ok)
	assert.Equal(t, "u-1", id)

	c.Set("user_id", 42)
	id, _ = GetUserID(c)
	assert.Equal(t, "42", id)
}
