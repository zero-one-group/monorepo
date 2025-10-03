package middleware

import (
	"go-modular/internal/config"

	"github.com/labstack/echo/v4"
	emw "github.com/labstack/echo/v4/middleware"
)

func CORSMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return emw.CORSWithConfig(emw.CORSConfig{
		AllowOrigins: cfg.App.CORSOrigins,
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.PATCH,
			echo.DELETE,
			echo.OPTIONS,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
		},
		ExposeHeaders: []string{
			echo.HeaderAccept,
			echo.HeaderAcceptEncoding,
			echo.HeaderAuthorization,
			echo.HeaderCacheControl,
			echo.HeaderConnection,
			echo.HeaderContentLength,
			echo.HeaderContentType,
			echo.HeaderOrigin,
			echo.HeaderXCSRFToken,
			echo.HeaderXRequestID,
			"Pragma",
			"User-Agent",
			"X-App-Audience",
			"X-Signature",
		},
		AllowCredentials: cfg.App.CORSCredentials,
		MaxAge:           cfg.App.CORSMaxAge,
	})
}
