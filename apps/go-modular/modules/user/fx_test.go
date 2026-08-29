package user

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-modular/internal/server"
	"go-modular/modules/user/services"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewUserHandler(t *testing.T) {
	svc := services.NewUserService(services.UserServiceOpts{})
	h := newUserHandler(slog.New(slog.NewTextHandler(io.Discard, nil)), svc)
	assert.NotNil(t, h)
}

func TestRegisterRoutes_TableAndJWTGuard(t *testing.T) {
	h := newUserHandler(slog.New(slog.NewTextHandler(io.Discard, nil)), services.NewUserService(services.UserServiceOpts{}))
	guard := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error { return c.NoContent(http.StatusTeapot) } // stand-in for the JWT guard
	}
	e := echo.New()
	registerRoutes(server.APIV1{Root: e.Group("/api/v1")}, h, guard)

	routes := map[string]bool{}
	for _, r := range e.Routes() {
		routes[r.Method+" "+r.Path] = true
	}
	for _, want := range []string{
		"POST /api/v1/users", "GET /api/v1/users", "GET /api/v1/users/:userId",
		"PUT /api/v1/users/:userId", "DELETE /api/v1/users/:userId",
	} {
		assert.True(t, routes[want], want)
	}

	for _, r := range []struct{ method, path string }{
		{http.MethodGet, "/api/v1/users"}, {http.MethodPost, "/api/v1/users"}, {http.MethodDelete, "/api/v1/users/abc"},
	} {
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, httptest.NewRequest(r.method, r.path, nil))
		assert.Equal(t, http.StatusTeapot, rec.Code, "JWT guard wraps %s %s", r.method, r.path)
	}
}
