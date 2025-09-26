package middleware

import (
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func AttachTraceProvider(provider *sdktrace.TracerProvider, serviceName string) echo.MiddlewareFunc {
	return otelecho.Middleware(serviceName, otelecho.WithTracerProvider(provider))
}
