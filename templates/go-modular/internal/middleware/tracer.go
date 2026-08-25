package middleware

import (
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

// TracingMiddleware returns Echo middleware that names the request span after
// the matched route.
//
// Compile-time auto-instrumentation (`otelc go build`) already creates a server
// span for every request, but it hooks net/http and therefore cannot see Echo's
// routing table. Those spans are named after the bare method ("GET") and carry
// url.path but no http.route, so requests to /users/1 and /users/2 land in
// separate buckets instead of aggregating under /users/:id.
//
// This middleware runs inside the auto-instrumented span and renames it to
// "GET /users/:id" once Echo has matched the route, which is what makes the
// latency aggregations in Jaeger meaningful.
//
// It reads the globally registered TracerProvider rather than taking one as an
// argument, because the SDK is installed by auto-instrumentation at start-up
// and is not owned by the application. Under a plain `go build` the global
// provider is a no-op and this middleware costs nothing.
func TracingMiddleware(serviceName string) echo.MiddlewareFunc {
	return otelecho.Middleware(serviceName)
}
