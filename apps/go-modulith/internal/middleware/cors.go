package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/zero-one-group/go-modulith/internal/config"
)

func CORS(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			origin := req.Header.Get("Origin")
			if origin != "" && isAllowedOrigin(origin, cfg.CORS.AllowedOrigins) {
				res.Header().Set("Access-Control-Allow-Origin", origin)
			}

			res.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.CORS.AllowedMethods, ", "))
			res.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.CORS.AllowedHeaders, ", "))
			res.Header().Set("Access-Control-Allow-Credentials", "true")
			res.Header().Set("Access-Control-Max-Age", "86400")

			if req.Method == http.MethodOptions {
				return c.NoContent(http.StatusNoContent)
			}

			return next(c)
		}
	}
}

func isAllowedOrigin(origin string, allowedOrigins []string) bool {
	for _, allowedOrigin := range allowedOrigins {
		if allowedOrigin == "*" || allowedOrigin == origin {
			return true
		}
	}
	return false
}