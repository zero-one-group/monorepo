package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwa"
	"{{ package_name | kebab_case }}/pkg/apputils"
)

// JWTMiddleware verifies a Bearer JWT and stores the parsed token claims in echo.Context.
// Usage:
//
//	e.Use(auth.JWTMiddleware(opts.JWTSecretKey, opts.SigningAlg))
func JWTMiddleware(secret []byte, alg jwa.SignatureAlgorithm) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
			}

			tokenStr := strings.TrimSpace(parts[1])

			// Use the shared JWT helper to parse & validate (validates exp/nbf etc).
			jwtGen := apputils.NewJWTGenerator(apputils.JWTConfig{
				SecretKey:  secret,
				SigningAlg: alg,
			})

			claims, err := jwtGen.ParseAndValidate(c.Request().Context(), tokenStr)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("invalid token: %v", err))
			}

			// Enforce token type to be "access" (defensive)
			if t, ok := claims["typ"]; ok {
				if ts, ok := t.(string); ok && ts != "access" {
					return echo.NewHTTPError(http.StatusUnauthorized, "token is not an access token")
				}
			}

			// store claims for handlers to use
			c.Set("jwt_claims", claims)

			// Extract common fields into shortcuts
			if sub, ok := claims["sub"]; ok {
				c.Set("user_id", fmt.Sprint(sub))
			}
			// session id may be stored as "sid" or "SID" (signing code used "SID")
			if sid, ok := claims["sid"]; ok {
				c.Set("session_id", fmt.Sprint(sid))
			} else if sid2, ok := claims["SID"]; ok {
				c.Set("session_id", fmt.Sprint(sid2))
			}
			if aud, ok := claims["aud"]; ok {
				c.Set("audience", fmt.Sprint(aud))
			}

			// also place the raw token string if needed
			c.Set("jwt_raw", tokenStr)

			// propagate claims into request context as well (for services that read ctx.Value("headers")/claims)
			ctx := context.WithValue(c.Request().Context(), apputils.JWTClaimsContextKey, claims)
			*c.Request() = *c.Request().Clone(ctx)

			return next(c)
		}
	}
}

// GetJWTClaims retrieves parsed JWT claims from the context (if present).
func GetJWTClaims(c echo.Context) (map[string]any, bool) {
	if v := c.Get("jwt_claims"); v != nil {
		if m, ok := v.(map[string]any); ok {
			return m, true
		}
	}
	return nil, false
}

// GetUserID attempts to return a user id from context (user_id key or token 'sub' claim).
func GetUserID(c echo.Context) (string, bool) {
	if v := c.Get("user_id"); v != nil {
		if s, ok := v.(string); ok {
			return s, true
		}
		return fmt.Sprint(v), true
	}
	if claims, ok := GetJWTClaims(c); ok {
		if sub, ok := claims["sub"]; ok {
			return fmt.Sprint(sub), true
		}
	}
	return "", false
}
