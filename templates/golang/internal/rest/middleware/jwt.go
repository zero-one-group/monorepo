package middleware

import (
	"net/http"
	"strings"
	"{{ package_name }}/utils"

	"github.com/labstack/echo/v4"
)

func ValidateUserToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.ErrUnauthorized
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) < 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid bearer token")
			}

			tokenString := parts[1]
			claims, err := utils.ValidateToken(tokenString)
			if err != nil {
				return echo.ErrUnauthorized
			}

			c.Set("user", claims)
			return next(c)
		}
	}
}
