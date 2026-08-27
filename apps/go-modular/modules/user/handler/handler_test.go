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

	"go-modular/modules/user/models"

	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeUserService struct {
	err    error
	users  map[uuid.UUID]*models.User
	last   *models.User
	filter *models.FilterUser
}

func newFakeUserService() *fakeUserService {
	return &fakeUserService{users: map[uuid.UUID]*models.User{}}
}

func (f *fakeUserService) CreateUser(_ context.Context, u *models.User) error {
	f.last = u
	if f.err == nil {
		u.ID = uuid.Must(uuid.NewV7())
	}
	return f.err
}
func (f *fakeUserService) GetUserByID(_ context.Context, id uuid.UUID) (*models.User, error) {
	if f.err != nil {
		return nil, f.err
	}
	u, ok := f.users[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}
func (f *fakeUserService) ListUsers(_ context.Context, filter *models.FilterUser) ([]*models.User, error) {
	f.filter = filter
	if f.err != nil {
		return nil, f.err
	}
	out := []*models.User{}
	for _, u := range f.users {
		out = append(out, u)
	}
	return out, nil
}
func (f *fakeUserService) UpdateUser(_ context.Context, u *models.User) error {
	f.last = u
	return f.err
}
func (f *fakeUserService) DeleteUser(context.Context, uuid.UUID) error { return f.err }
func (f *fakeUserService) GetUserByEmail(context.Context, string) (*models.User, error) {
	return nil, f.err
}
func (f *fakeUserService) GetUserByUsername(context.Context, string) (*models.User, error) {
	return nil, f.err
}
func (f *fakeUserService) MarkEmailVerified(context.Context, uuid.UUID) error { return f.err }

// newRouter wires the handler through a real Echo router the same way module.go does,
// so path params, binding and JSON rendering are exercised for real.
func newRouter(svc *fakeUserService) *echo.Echo {
	h := NewHandler(&HandlerOpts{Logger: slog.New(slog.NewTextHandler(io.Discard, nil)), UserService: svc})
	e := echo.New()
	g := e.Group("/api/v1/users")
	g.POST("", h.CreateUser)
	g.GET("", h.ListUsers)
	g.GET("/:userId", h.GetUser)
	g.PUT("/:userId", h.UpdateUser)
	g.DELETE("/:userId", h.DeleteUser)
	return e
}

func do(e *echo.Echo, method, path, body string) (*httptest.ResponseRecorder, map[string]any) {
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

func TestCreateUser(t *testing.T) {
	t.Run("201 with the created user", func(t *testing.T) {
		svc := newFakeUserService()
		rec, body := do(newRouter(svc), http.MethodPost, "/api/v1/users", `{"name":"Jane","email":"jane@example.com"}`)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "jane@example.com", body["email"])
		assert.Equal(t, "Jane", svc.last.DisplayName)
		assert.Nil(t, svc.last.EmailVerifiedAt, "new users start unverified")
	})
	t.Run("400 on malformed JSON", func(t *testing.T) {
		rec, body := do(newRouter(newFakeUserService()), http.MethodPost, "/api/v1/users", `{"name":`)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "Invalid request payload", body["error"])
	})
	t.Run("400 with field details on validation failure", func(t *testing.T) {
		rec, body := do(newRouter(newFakeUserService()), http.MethodPost, "/api/v1/users", `{"name":"","email":"not-an-email"}`)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "Validation failed", body["error"])
		details, _ := body["details"].(map[string]any)
		assert.NotEmpty(t, details)
	})
	t.Run("500 when the service fails, without leaking the error", func(t *testing.T) {
		svc := newFakeUserService()
		svc.err = errors.New("unique violation on email")
		rec, body := do(newRouter(svc), http.MethodPost, "/api/v1/users", `{"name":"Jane","email":"jane@example.com"}`)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.NotContains(t, body["error"], "unique violation")
	})
}

func TestListUsers(t *testing.T) {
	svc := newFakeUserService()
	id := uuid.Must(uuid.NewV7())
	svc.users[id] = &models.User{ID: id, Email: "a@example.com"}

	t.Run("200 with the list and bound query filters", func(t *testing.T) {
		rec, _ := do(newRouter(svc), http.MethodGet, "/api/v1/users?search=am&limit=5&offset=10", "")
		assert.Equal(t, http.StatusOK, rec.Code)
		var list []map[string]any
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &list))
		assert.Len(t, list, 1)
		require.NotNil(t, svc.filter.Search)
		assert.Equal(t, "am", *svc.filter.Search)
		assert.Equal(t, 5, svc.filter.Limit)
		assert.Equal(t, 10, svc.filter.Offset)
	})
	t.Run("400 on an unbindable filter", func(t *testing.T) {
		rec, body := do(newRouter(svc), http.MethodGet, "/api/v1/users?limit=abc", "")
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "Invalid filter parameters", body["error"])
	})
	t.Run("500 when the service fails", func(t *testing.T) {
		failing := newFakeUserService()
		failing.err = errors.New("db down")
		rec, _ := do(newRouter(failing), http.MethodGet, "/api/v1/users", "")
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestGetUser(t *testing.T) {
	svc := newFakeUserService()
	id := uuid.Must(uuid.NewV7())
	svc.users[id] = &models.User{ID: id, Email: "a@example.com"}
	e := newRouter(svc)

	rec, body := do(e, http.MethodGet, "/api/v1/users/"+id.String(), "")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, id.String(), body["id"])

	rec, body = do(e, http.MethodGet, "/api/v1/users/not-a-uuid", "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "User ID in path must be a valid UUID", body["error"])

	rec, body = do(e, http.MethodGet, "/api/v1/users/"+uuid.Must(uuid.NewV7()).String(), "")
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, "User not found", body["error"])
}

func TestUpdateUser(t *testing.T) {
	id := uuid.Must(uuid.NewV7())
	t.Run("200 and maps the payload onto the path id", func(t *testing.T) {
		svc := newFakeUserService()
		rec, body := do(newRouter(svc), http.MethodPut, "/api/v1/users/"+id.String(), `{"name":"New","email":"new@example.com"}`)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "User updated successfully", body["message"])
		assert.Equal(t, id, svc.last.ID)
		assert.Equal(t, "New", svc.last.DisplayName)
	})
	t.Run("400s", func(t *testing.T) {
		e := newRouter(newFakeUserService())
		rec, _ := do(e, http.MethodPut, "/api/v1/users/bad", `{"name":"New","email":"new@example.com"}`)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		rec, _ = do(e, http.MethodPut, "/api/v1/users/"+id.String(), `{`)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		rec, body := do(e, http.MethodPut, "/api/v1/users/"+id.String(), `{"name":"","email":"x"}`)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, "Validation failed", body["error"])
	})
	t.Run("500 when the service fails", func(t *testing.T) {
		svc := newFakeUserService()
		svc.err = errors.New("db down")
		rec, _ := do(newRouter(svc), http.MethodPut, "/api/v1/users/"+id.String(), `{"name":"New","email":"new@example.com"}`)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestDeleteUser(t *testing.T) {
	id := uuid.Must(uuid.NewV7())
	rec, body := do(newRouter(newFakeUserService()), http.MethodDelete, "/api/v1/users/"+id.String(), "")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "User deleted successfully", body["message"])

	rec, _ = do(newRouter(newFakeUserService()), http.MethodDelete, "/api/v1/users/bad", "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	svc := newFakeUserService()
	svc.err = errors.New("no rows")
	rec, body = do(newRouter(svc), http.MethodDelete, "/api/v1/users/"+id.String(), "")
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, "User not found", body["error"])
}
