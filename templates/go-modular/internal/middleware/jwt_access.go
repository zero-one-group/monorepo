package middleware

import "github.com/labstack/echo/v4"

// JWTAccess is Echo middleware that requires a valid access JWT.
// Provided by the auth Module; consumed by Modules that protect routes.
type JWTAccess echo.MiddlewareFunc
