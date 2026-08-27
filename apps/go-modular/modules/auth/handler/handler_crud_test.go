package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-modular/modules/auth/models"
	"go-modular/modules/user/repository"

	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newCrudRouter mirrors module.go's route table (minus the JWT guard, tested in modules/auth).
func newCrudRouter(svc *crudFake) *echo.Echo {
	h := NewHandler(&HandlerOpts{Logger: slog.New(slog.NewTextHandler(io.Discard, nil)), AuthService: svc})
	e := echo.New()
	g := e.Group("/api/v1/auth")
	g.GET("/verify-email", h.ValidateEmailVerificationByLink)
	g.POST("/refresh-token", h.CreateRefreshToken)
	g.PUT("/refresh-token", h.UpdateRefreshToken)
	g.GET("/refresh-token/:tokenId", h.GetRefreshToken)
	g.DELETE("/refresh-token/:tokenId", h.DeleteRefreshToken)
	g.POST("/verification/email/initiate", h.InitiateEmailVerification)
	g.POST("/verification/email/validate", h.ValidateEmailVerification)
	g.POST("/password", h.SetUserPassword)
	g.PUT("/password/:userId", h.UpdateUserPassword)
	g.POST("/session", h.CreateSession)
	g.PUT("/session", h.UpdateSession)
	g.GET("/session/:sessionId", h.GetSession)
	g.DELETE("/session/:sessionId", h.DeleteSession)
	g.POST("/verification/email/revoke", h.RevokeEmailVerification)
	g.POST("/verification/email/resend", h.ResendEmailVerification)
	return e
}

func call(e *echo.Echo, method, path, body string) (*httptest.ResponseRecorder, map[string]any) {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	var out map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &out)
	return rec, out
}

var (
	uid    = uuid.Must(uuid.NewV7())
	future = time.Now().Add(time.Hour).UTC().Format(time.RFC3339)
)

// ---------------------------------------------------------------- password

func TestSetUserPassword(t *testing.T) {
	t.Run("201 maps user_id and forwards the plain password for the service to hash", func(t *testing.T) {
		svc := &crudFake{}
		rec, body := call(newCrudRouter(svc), http.MethodPost, "/api/v1/auth/password",
			`{"user_id":"`+uid.String()+`","password":"secure.password","password_confirmation":"secure.password"}`)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "Password set successfully", body["message"])
		assert.Equal(t, uid, svc.lastPassword.UserID)
		assert.Equal(t, "secure.password", svc.lastPassword.PasswordHash)
	})
	t.Run("400 on malformed json, validation (short / mismatch / bad uuid)", func(t *testing.T) {
		e := newCrudRouter(&crudFake{})
		for _, body := range []string{
			`{`,
			`{"user_id":"` + uid.String() + `","password":"short","password_confirmation":"short"}`,
			`{"user_id":"` + uid.String() + `","password":"secure.password","password_confirmation":"different.one"}`,
			`{"user_id":"nope","password":"secure.password","password_confirmation":"secure.password"}`,
		} {
			rec, _ := call(e, http.MethodPost, "/api/v1/auth/password", body)
			assert.Equal(t, http.StatusBadRequest, rec.Code, body)
		}
	})
	t.Run("500 hides the service error", func(t *testing.T) {
		svc := &crudFake{crudErr: errors.New("db down")}
		rec, body := call(newCrudRouter(svc), http.MethodPost, "/api/v1/auth/password",
			`{"user_id":"`+uid.String()+`","password":"secure.password","password_confirmation":"secure.password"}`)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, "Failed to set password", body["error"])
	})
}

func TestUpdateUserPassword(t *testing.T) {
	valid := `{"current_password":"current.password","new_password":"secure.password","password_confirmation":"secure.password"}`
	t.Run("200 forwards current and new password for the given user", func(t *testing.T) {
		svc := &crudFake{}
		rec, body := call(newCrudRouter(svc), http.MethodPut, "/api/v1/auth/password/"+uid.String(), valid)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Password updated successfully", body["message"])
		assert.Equal(t, []string{uid.String(), "current.password", "secure.password"}, svc.lastUpdatePw)
	})
	t.Run("400s", func(t *testing.T) {
		e := newCrudRouter(&crudFake{})
		rec, _ := call(e, http.MethodPut, "/api/v1/auth/password/not-a-uuid", valid)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		rec, _ = call(e, http.MethodPut, "/api/v1/auth/password/"+uid.String(), `{`)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		rec, body := call(e, http.MethodPut, "/api/v1/auth/password/"+uid.String(), `{"current_password":"x","new_password":"y","password_confirmation":"z"}`)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "Validation failed", body["error"])
	})
	t.Run("a wrong current password surfaces as 400 with the reason", func(t *testing.T) {
		svc := &crudFake{crudErr: errors.New("current password is incorrect")}
		rec, body := call(newCrudRouter(svc), http.MethodPut, "/api/v1/auth/password/"+uid.String(), valid)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "current password is incorrect", body["error"])
	})
}

// ---------------------------------------------------------------- session

func TestCreateSession(t *testing.T) {
	t.Run("201 maps every optional field", func(t *testing.T) {
		svc := &crudFake{}
		rec, _ := call(newCrudRouter(svc), http.MethodPost, "/api/v1/auth/session",
			`{"user_id":"`+uid.String()+`","token_hash":"h","user_agent":"UA","device_name":"Pixel","device_fingerprint":"fp","ip_address":"10.0.0.7","expires_at":"`+future+`"}`)
		require.Equal(t, http.StatusCreated, rec.Code)
		s := svc.lastSession
		assert.Equal(t, uid, s.UserID)
		assert.Equal(t, "h", s.TokenHash)
		assert.Equal(t, "UA", *s.UserAgent)
		assert.Equal(t, "Pixel", *s.DeviceName)
		assert.Equal(t, "fp", *s.DeviceFingerprint)
		assert.Equal(t, "10.0.0.7", s.IPAddress.String())
		assert.WithinDuration(t, time.Now().Add(time.Hour), s.ExpiresAt, 2*time.Second)
	})
	t.Run("400s", func(t *testing.T) {
		e := newCrudRouter(&crudFake{})
		for name, body := range map[string]string{
			"malformed":      `{`,
			"missing hash":   `{"user_id":"` + uid.String() + `","expires_at":"` + future + `"}`,
			"bad ip":         `{"user_id":"` + uid.String() + `","token_hash":"h","ip_address":"not-an-ip","expires_at":"` + future + `"}`,
			"bad expires_at": `{"user_id":"` + uid.String() + `","token_hash":"h","expires_at":"tomorrow"}`,
			"non-uuid user":  `{"user_id":"abc","token_hash":"h","expires_at":"` + future + `"}`,
		} {
			rec, _ := call(e, http.MethodPost, "/api/v1/auth/session", body)
			assert.Equal(t, http.StatusBadRequest, rec.Code, name)
		}
	})
	t.Run("500 on service error", func(t *testing.T) {
		rec, _ := call(newCrudRouter(&crudFake{crudErr: errors.New("x")}), http.MethodPost, "/api/v1/auth/session",
			`{"user_id":"`+uid.String()+`","token_hash":"h","expires_at":"`+future+`"}`)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestUpdateSession(t *testing.T) {
	sid := uuid.Must(uuid.NewV7())
	t.Run("200 maps timestamps, ip and revoked_by", func(t *testing.T) {
		svc := &crudFake{}
		rec, _ := call(newCrudRouter(svc), http.MethodPut, "/api/v1/auth/session",
			`{"session_id":"`+sid.String()+`","ip_address":"::1","refreshed_at":"`+future+`","revoked_at":"`+future+`","revoked_by":"`+uid.String()+`"}`)
		require.Equal(t, http.StatusOK, rec.Code)
		s := svc.lastSession
		assert.Equal(t, sid, s.ID)
		assert.Equal(t, "::1", s.IPAddress.String())
		assert.NotNil(t, s.RefreshedAt)
		assert.NotNil(t, s.RevokedAt)
		assert.Equal(t, uid, *s.RevokedBy)
	})
	t.Run("400s", func(t *testing.T) {
		e := newCrudRouter(&crudFake{})
		for name, body := range map[string]string{
			"malformed":        `{`,
			"missing id":       `{}`,
			"bad ip":           `{"session_id":"` + sid.String() + `","ip_address":"zzz"}`,
			"bad refreshed_at": `{"session_id":"` + sid.String() + `","refreshed_at":"soon"}`,
			"bad revoked_at":   `{"session_id":"` + sid.String() + `","revoked_at":"soon"}`,
			"bad revoked_by":   `{"session_id":"` + sid.String() + `","revoked_by":"not-uuid"}`,
		} {
			rec, _ := call(e, http.MethodPut, "/api/v1/auth/session", body)
			assert.Equal(t, http.StatusBadRequest, rec.Code, name)
		}
	})
	t.Run("500 on service error", func(t *testing.T) {
		rec, _ := call(newCrudRouter(&crudFake{crudErr: errors.New("x")}), http.MethodPut, "/api/v1/auth/session", `{"session_id":"`+sid.String()+`"}`)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestGetAndDeleteSession(t *testing.T) {
	sid := uuid.Must(uuid.NewV7())
	found := &crudFake{session: &models.Session{ID: sid, UserID: uid}}
	rec, body := call(newCrudRouter(found), http.MethodGet, "/api/v1/auth/session/"+sid.String(), "")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, sid.String(), body["id"])

	rec, body = call(newCrudRouter(&crudFake{}), http.MethodGet, "/api/v1/auth/session/"+sid.String(), "")
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, "Session not found", body["error"])

	rec, _ = call(newCrudRouter(&crudFake{}), http.MethodGet, "/api/v1/auth/session/bad", "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	rec, _ = call(newCrudRouter(&crudFake{crudErr: errors.New("x")}), http.MethodGet, "/api/v1/auth/session/"+sid.String(), "")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	svc := &crudFake{}
	rec, body = call(newCrudRouter(svc), http.MethodDelete, "/api/v1/auth/session/"+sid.String(), "")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Session deleted successfully", body["message"])
	assert.Equal(t, []uuid.UUID{sid}, svc.deleted)

	rec, _ = call(newCrudRouter(&crudFake{}), http.MethodDelete, "/api/v1/auth/session/bad", "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	rec, _ = call(newCrudRouter(&crudFake{crudErr: errors.New("x")}), http.MethodDelete, "/api/v1/auth/session/"+sid.String(), "")
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

// ---------------------------------------------------------------- refresh token

func TestCreateRefreshToken(t *testing.T) {
	sid := uuid.Must(uuid.NewV7())
	t.Run("201 maps fields incl. optional session and ip", func(t *testing.T) {
		svc := &crudFake{}
		rec, _ := call(newCrudRouter(svc), http.MethodPost, "/api/v1/auth/refresh-token",
			`{"user_id":"`+uid.String()+`","session_id":"`+sid.String()+`","token_hash":"h","ip_address":"192.168.1.1","user_agent":"UA","expires_at":"`+future+`"}`)
		require.Equal(t, http.StatusCreated, rec.Code)
		tok := svc.lastToken
		assert.Equal(t, uid, tok.UserID)
		assert.Equal(t, sid, *tok.SessionID)
		assert.Equal(t, []byte("h"), tok.TokenHash)
		assert.Equal(t, "192.168.1.1", tok.IPAddress.String())
		assert.Equal(t, "UA", *tok.UserAgent)
	})
	t.Run("400s", func(t *testing.T) {
		e := newCrudRouter(&crudFake{})
		for name, body := range map[string]string{
			"malformed":      `{`,
			"missing hash":   `{"user_id":"` + uid.String() + `","expires_at":"` + future + `"}`,
			"bad user":       `{"user_id":"x","token_hash":"h","expires_at":"` + future + `"}`,
			"bad session":    `{"user_id":"` + uid.String() + `","session_id":"x","token_hash":"h","expires_at":"` + future + `"}`,
			"bad ip":         `{"user_id":"` + uid.String() + `","token_hash":"h","ip_address":"x","expires_at":"` + future + `"}`,
			"bad expires_at": `{"user_id":"` + uid.String() + `","token_hash":"h","expires_at":"x"}`,
		} {
			rec, _ := call(e, http.MethodPost, "/api/v1/auth/refresh-token", body)
			assert.Equal(t, http.StatusBadRequest, rec.Code, name)
		}
	})
	t.Run("500", func(t *testing.T) {
		rec, _ := call(newCrudRouter(&crudFake{crudErr: errors.New("x")}), http.MethodPost, "/api/v1/auth/refresh-token",
			`{"user_id":"`+uid.String()+`","token_hash":"h","expires_at":"`+future+`"}`)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestUpdateGetDeleteRefreshToken(t *testing.T) {
	tid := uuid.Must(uuid.NewV7())
	t.Run("update 200 + 400s + 500", func(t *testing.T) {
		svc := &crudFake{}
		rec, _ := call(newCrudRouter(svc), http.MethodPut, "/api/v1/auth/refresh-token",
			`{"token_id":"`+tid.String()+`","ip_address":"10.1.1.1","user_agent":"UA","revoked_at":"`+future+`","revoked_by":"`+uid.String()+`"}`)
		require.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, tid, svc.lastToken.ID)
		assert.NotNil(t, svc.lastToken.RevokedAt)
		assert.Equal(t, uid, *svc.lastToken.RevokedBy)

		e := newCrudRouter(&crudFake{})
		for name, body := range map[string]string{
			"malformed":      `{`,
			"missing id":     `{}`,
			"bad token id":   `{"token_id":"x"}`,
			"bad ip":         `{"token_id":"` + tid.String() + `","ip_address":"x"}`,
			"bad revoked_at": `{"token_id":"` + tid.String() + `","revoked_at":"x"}`,
			"bad revoked_by": `{"token_id":"` + tid.String() + `","revoked_by":"x"}`,
		} {
			rec, _ := call(e, http.MethodPut, "/api/v1/auth/refresh-token", body)
			assert.Equal(t, http.StatusBadRequest, rec.Code, name)
		}
		rec, _ = call(newCrudRouter(&crudFake{crudErr: errors.New("x")}), http.MethodPut, "/api/v1/auth/refresh-token", `{"token_id":"`+tid.String()+`"}`)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
	t.Run("get 200 / 404 / 400 / 500", func(t *testing.T) {
		rec, body := call(newCrudRouter(&crudFake{token: &models.RefreshToken{ID: tid}}), http.MethodGet, "/api/v1/auth/refresh-token/"+tid.String(), "")
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, tid.String(), body["id"])
		rec, _ = call(newCrudRouter(&crudFake{}), http.MethodGet, "/api/v1/auth/refresh-token/"+tid.String(), "")
		assert.Equal(t, http.StatusNotFound, rec.Code)
		rec, _ = call(newCrudRouter(&crudFake{}), http.MethodGet, "/api/v1/auth/refresh-token/bad", "")
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		rec, _ = call(newCrudRouter(&crudFake{crudErr: errors.New("x")}), http.MethodGet, "/api/v1/auth/refresh-token/"+tid.String(), "")
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
	t.Run("delete 200 / 400 / 500", func(t *testing.T) {
		svc := &crudFake{}
		rec, _ := call(newCrudRouter(svc), http.MethodDelete, "/api/v1/auth/refresh-token/"+tid.String(), "")
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, []uuid.UUID{tid}, svc.deleted)
		rec, _ = call(newCrudRouter(&crudFake{}), http.MethodDelete, "/api/v1/auth/refresh-token/bad", "")
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		rec, _ = call(newCrudRouter(&crudFake{crudErr: errors.New("x")}), http.MethodDelete, "/api/v1/auth/refresh-token/"+tid.String(), "")
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

// ---------------------------------------------------------------- verification

func TestInitiateAndResendEmailVerification(t *testing.T) {
	for _, path := range []string{"/api/v1/auth/verification/email/initiate", "/api/v1/auth/verification/email/resend"} {
		t.Run(path, func(t *testing.T) {
			svc := &crudFake{}
			rec, body := call(newCrudRouter(svc), http.MethodPost, path, `{"email":"jane@example.com","redirect_to":"https://app.example.com/verified"}`)
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, body["message"], "if the email is registered", "never confirms whether an email exists")
			assert.Equal(t, []string{"jane@example.com", "https://app.example.com/verified"}, svc.lastVerify)

			e := newCrudRouter(&crudFake{})
			rec, _ = call(e, http.MethodPost, path, `{`)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			rec, _ = call(e, http.MethodPost, path, `{"email":"nope","redirect_to":"not a url"}`)
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			rec, _ = call(newCrudRouter(&crudFake{verifyErr: errors.New("user not found")}), http.MethodPost, path, `{"email":"jane@example.com"}`)
			assert.Equal(t, http.StatusNotFound, rec.Code)
			rec, _ = call(newCrudRouter(&crudFake{verifyErr: repository.ErrNotFound}), http.MethodPost, path, `{"email":"jane@example.com"}`)
			assert.Equal(t, http.StatusNotFound, rec.Code)
			rec, _ = call(newCrudRouter(&crudFake{verifyErr: errors.New("token still valid")}), http.MethodPost, path, `{"email":"jane@example.com"}`)
			assert.Equal(t, http.StatusConflict, rec.Code)
			rec, _ = call(newCrudRouter(&crudFake{verifyErr: errors.New("smtp down")}), http.MethodPost, path, `{"email":"jane@example.com"}`)
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		})
	}
}

func TestValidateEmailVerification(t *testing.T) {
	path := "/api/v1/auth/verification/email/validate"
	rec, body := call(newCrudRouter(&crudFake{verifyOK: true}), http.MethodPost, path, `{"token":"abc"}`)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Email successfully verified", body["message"])

	rec, body = call(newCrudRouter(&crudFake{verifyOK: false}), http.MethodPost, path, `{"token":"abc"}`)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Equal(t, "Invalid or expired token", body["error"])

	rec, _ = call(newCrudRouter(&crudFake{verifyErr: errors.New("expired")}), http.MethodPost, path, `{"token":"abc"}`)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	e := newCrudRouter(&crudFake{})
	rec, _ = call(e, http.MethodPost, path, `{`)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	rec, _ = call(e, http.MethodPost, path, `{}`)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRevokeEmailVerification(t *testing.T) {
	path := "/api/v1/auth/verification/email/revoke"
	svc := &crudFake{}
	rec, body := call(newCrudRouter(svc), http.MethodPost, path, `{"token":"abc"}`)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Verification token revoked", body["message"])
	assert.Equal(t, []string{"abc"}, svc.lastVerify)

	rec, _ = call(newCrudRouter(&crudFake{verifyErr: repository.ErrNotFound}), http.MethodPost, path, `{"token":"abc"}`)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	rec, _ = call(newCrudRouter(&crudFake{verifyErr: errors.New("x")}), http.MethodPost, path, `{"token":"abc"}`)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	e := newCrudRouter(&crudFake{})
	rec, _ = call(e, http.MethodPost, path, `{`)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	rec, _ = call(e, http.MethodPost, path, `{}`)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestValidateEmailVerificationByLink(t *testing.T) {
	base := "/api/v1/auth/verify-email"
	redirect := "https://app.example.com/verified?src=mail"

	t.Run("json mode", func(t *testing.T) {
		rec, body := call(newCrudRouter(&crudFake{verifyOK: true}), http.MethodGet, base+"?token=abc", "")
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Email successfully verified", body["message"])
		rec, _ = call(newCrudRouter(&crudFake{verifyOK: false}), http.MethodGet, base+"?token=abc", "")
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		rec, body = call(newCrudRouter(&crudFake{verifyErr: errors.New("expired")}), http.MethodGet, base+"?token=abc", "")
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, "expired", body["details"])
		rec, body = call(newCrudRouter(&crudFake{}), http.MethodGet, base, "")
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "token query parameter is required", body["error"])
	})

	t.Run("redirect mode appends the outcome to redirect_to", func(t *testing.T) {
		q := "&redirect_to=" + strings.ReplaceAll(redirect, "&", "%26")
		rec, _ := call(newCrudRouter(&crudFake{verifyOK: true}), http.MethodGet, base+"?token=abc"+q, "")
		assert.Equal(t, http.StatusFound, rec.Code)
		loc := rec.Header().Get(echo.HeaderLocation)
		assert.Contains(t, loc, "verified=true")
		assert.Contains(t, loc, "src=mail", "existing query params are preserved")

		rec, _ = call(newCrudRouter(&crudFake{verifyErr: errors.New("expired")}), http.MethodGet, base+"?token=abc"+q, "")
		assert.Equal(t, http.StatusFound, rec.Code)
		assert.Contains(t, rec.Header().Get(echo.HeaderLocation), "verified=false")
		assert.Contains(t, rec.Header().Get(echo.HeaderLocation), "error=expired")

		rec, _ = call(newCrudRouter(&crudFake{verifyOK: false}), http.MethodGet, base+"?token=abc"+q, "")
		assert.Contains(t, rec.Header().Get(echo.HeaderLocation), "verified=false")

		rec, _ = call(newCrudRouter(&crudFake{}), http.MethodGet, base+"?redirect_to=https://app.example.com/x", "")
		assert.Equal(t, http.StatusFound, rec.Code)
		assert.Contains(t, rec.Header().Get(echo.HeaderLocation), "error=token_required")

		rec, _ = call(newCrudRouter(&crudFake{verifyOK: true}), http.MethodGet, base+"?token=abc&redirect_to=%3A%2F%2Fbad", "")
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		rec, _ = call(newCrudRouter(&crudFake{}), http.MethodGet, base+"?redirect_to=%3A%2F%2Fbad", "")
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
