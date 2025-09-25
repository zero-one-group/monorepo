package middleware

import (
	"github.com/labstack/echo/v4"
)

func SecurityHeadersMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// X-Content-Type-Options - prevents MIME type sniffing
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")

			// X-Frame-Options - prevents clickjacking
			c.Response().Header().Set("X-Frame-Options", "DENY")

			// Strict-Transport-Security (HSTS) - forces HTTPS
			c.Response().Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

			// Content-Security-Policy - restrictive for API
			c.Response().Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")

			return next(c)
		}
	}
}
