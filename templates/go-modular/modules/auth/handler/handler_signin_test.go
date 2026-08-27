package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"{{ package_name | kebab_case }}/modules/auth/models"
	"{{ package_name | kebab_case }}/modules/auth/services"
	"{{ package_name | kebab_case }}/pkg/apputils"

	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeAuthService only implements the sign-in surface for real; everything else
// returns errNotImplemented so an unexpected call fails loudly.
type fakeAuthService struct {
	services.AuthServiceInterface // nil embed: any un-overridden method panics = "not expected here"
	err                           error
	audienceSeen                  string
	identifier                    string
}

func (f *fakeAuthService) signin(ctx context.Context, identifier string) (*models.AuthenticatedUser, error) {
	f.identifier = identifier
	if h, ok := ctx.Value(apputils.HeadersContextKey).(map[string]string); ok {
		f.audienceSeen = h["X-App-Audience"]
	}
	if f.err != nil {
		return nil, f.err
	}
	sid := uuid.Must(uuid.NewV7())
	return &models.AuthenticatedUser{
		UserWithCredentials: models.UserWithCredentials{AccessToken: "access", RefreshToken: "refresh"},
		SessionID:           &sid,
	}, nil
}
func (f *fakeAuthService) SignInWithEmail(ctx context.Context, email, _ string) (*models.AuthenticatedUser, error) {
	return f.signin(ctx, email)
}
func (f *fakeAuthService) SignInWithUsername(ctx context.Context, username, _ string) (*models.AuthenticatedUser, error) {
	return f.signin(ctx, username)
}

func newSigninRouter(svc *fakeAuthService) *echo.Echo {
	h := NewHandler(&HandlerOpts{Logger: slog.New(slog.NewTextHandler(io.Discard, nil)), AuthService: svc})
	e := echo.New()
	e.POST("/api/v1/auth/signin/email", h.SignInWithEmail)
	e.POST("/api/v1/auth/signin/username", h.SignInWithUsername)
	return e
}

func post(e *echo.Echo, path, body string, headers map[string]string) (*httptest.ResponseRecorder, map[string]any) {
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	var out map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &out)
	return rec, out
}

func TestSignIn_Success(t *testing.T) {
	for _, tc := range []struct{ path, body, wantIdentifier string }{
		{"/api/v1/auth/signin/email", `{"email":"jane@example.com","password":"secret"}`, "jane@example.com"},
		{"/api/v1/auth/signin/username", `{"username":"jane","password":"secret"}`, "jane"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			svc := &fakeAuthService{}
			rec, body := post(newSigninRouter(svc), tc.path, tc.body, map[string]string{"X-App-Audience": "admin-app"})
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "access", body["access_token"])
			assert.Equal(t, "refresh", body["refresh_token"])
			assert.NotEmpty(t, body["session_id"])
			assert.Equal(t, tc.wantIdentifier, svc.identifier)
			assert.Equal(t, "admin-app", svc.audienceSeen, "X-App-Audience must reach the service via the context")
		})
	}
}

func TestSignIn_ErrorMapping(t *testing.T) {
	cases := []struct {
		name       string
		err        error
		wantStatus int
		wantError  string
	}{
		{"invalid credentials → 401", services.ErrInvalidCredentials, http.StatusUnauthorized, "Invalid email or password"},
		{"unverified email → 401", services.ErrEmailNotVerified, http.StatusUnauthorized, "Email is not verified"},
		{"anything else → 500", errors.New("db down"), http.StatusInternalServerError, "Internal server error"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			svc := &fakeAuthService{err: c.err}
			rec, body := post(newSigninRouter(svc), "/api/v1/auth/signin/email", `{"email":"jane@example.com","password":"secret"}`, nil)
			assert.Equal(t, c.wantStatus, rec.Code)
			assert.Equal(t, c.wantError, body["error"])
		})
	}
	t.Run("username variant wording", func(t *testing.T) {
		svc := &fakeAuthService{err: services.ErrInvalidCredentials}
		rec, body := post(newSigninRouter(svc), "/api/v1/auth/signin/username", `{"username":"jane","password":"x"}`, nil)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, "Invalid username or password", body["error"])
	})
}

func TestSignIn_RequestValidation(t *testing.T) {
	e := newSigninRouter(&fakeAuthService{})
	for _, tc := range []struct{ name, path, body, wantError string }{
		{"malformed json", "/api/v1/auth/signin/email", `{"email":`, "Invalid request payload"},
		{"bad email format", "/api/v1/auth/signin/email", `{"email":"nope","password":"x"}`, "Validation failed"},
		{"missing password", "/api/v1/auth/signin/email", `{"email":"a@b.co"}`, "Validation failed"},
		{"missing username", "/api/v1/auth/signin/username", `{"password":"x"}`, "Validation failed"},
		{"malformed json (username)", "/api/v1/auth/signin/username", `{`, "Invalid request payload"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec, body := post(e, tc.path, tc.body, nil)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, tc.wantError, body["error"])
			if tc.wantError == "Validation failed" {
				details, ok := body["details"].(map[string]any)
				require.True(t, ok)
				assert.NotEmpty(t, details)
			}
		})
	}
}
