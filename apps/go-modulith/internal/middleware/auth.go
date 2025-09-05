package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/zero-one-group/go-modulith/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	EmailKey  contextKey = "email"
)

func JWTAuth(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, span := otel.Tracer("middleware").Start(c.Request().Context(), "jwt_auth")
			defer span.End()

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				span.SetAttributes(attribute.String("error", "missing authorization header"))
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				span.SetAttributes(attribute.String("error", "invalid authorization header format"))
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
			}

			token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid signing method")
				}
				return []byte(cfg.JWT.Secret), nil
			})

			if err != nil {
				span.RecordError(err)
				span.SetAttributes(attribute.String("error", "invalid token"))
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			claims, ok := token.Claims.(*Claims)
			if !ok || !token.Valid {
				span.SetAttributes(attribute.String("error", "invalid token claims"))
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
			}

			ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, EmailKey, claims.Email)
			c.SetRequest(c.Request().WithContext(ctx))

			span.SetAttributes(
				attribute.String("user_id", claims.UserID),
				attribute.String("email", claims.Email),
			)

			return next(c)
		}
	}
}

func GetUserID(c echo.Context) string {
	if userID := c.Request().Context().Value(UserIDKey); userID != nil {
		return userID.(string)
	}
	return ""
}

func GetEmail(c echo.Context) string {
	if email := c.Request().Context().Value(EmailKey); email != nil {
		return email.(string)
	}
	return ""
}