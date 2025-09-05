package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/labstack/echo/v4"
	"github.com/zero-one-group/go-modulith/internal/errors"
	"github.com/zero-one-group/go-modulith/internal/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func Recovery() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, span := otel.Tracer("middleware").Start(c.Request().Context(), "recovery")
			defer span.End()

			defer func() {
				if r := recover(); r != nil {
					log := logger.FromContext(ctx)
					stack := debug.Stack()

					err := fmt.Errorf("panic: %v", r)
					span.RecordError(err)
					span.SetAttributes(
						attribute.String("panic.value", fmt.Sprintf("%v", r)),
						attribute.String("panic.stack", string(stack)),
					)

					log.Error("Panic recovered",
						"error", r,
						"stack", string(stack),
						"path", c.Request().URL.Path,
						"method", c.Request().Method,
					)

					if !c.Response().Committed {
						c.JSON(http.StatusInternalServerError, errors.ErrorResponse{
							Error: "Internal server error",
							Code:  "INTERNAL_ERROR",
						})
					}
				}
			}()

			return next(c)
		}
	}
}