package middleware

import (
	"os"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func AttachTraceProvider(provider *sdktrace.TracerProvider) echo.MiddlewareFunc {
	return otelecho.Middleware(
		os.Getenv("SERVICE_NAME"),
		otelecho.WithTracerProvider(provider),
	)
}
