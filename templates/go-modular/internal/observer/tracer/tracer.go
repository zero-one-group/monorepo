// Package tracer exposes helpers for working with the OpenTelemetry tracer.
//
// The application does not construct a TracerProvider itself. Traces are set up
// by compile-time auto-instrumentation: building with `otelc go build` injects
// the OpenTelemetry SDK into `main`, configures it from the standard OTEL_*
// environment variables, and registers the resulting provider globally.
//
// Because of that, everything here reads from the global provider rather than a
// provider we own. When the binary is built with a plain `go build`, no SDK is
// injected, the global provider is the no-op implementation, and each helper
// below degrades to doing nothing.
//
// See docs/OBSERVABILITY.md for the build and configuration workflow.
package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// instrumentationName identifies spans created through this package in the
// backend's "otel.library.name" attribute.
const instrumentationName = "{{ package_name | kebab_case }}"

// Tracer returns the named tracer from the global provider.
//
// It is safe to call before any SDK is installed: the global provider is a
// no-op until auto-instrumentation replaces it during process start-up.
func Tracer() trace.Tracer {
	return otel.Tracer(instrumentationName)
}

// Start begins a span for a unit of work that is not covered by
// auto-instrumentation, such as a domain service method.
//
// Auto-instrumentation already covers inbound and outbound HTTP, and otelpgx
// covers database queries, so reach for this only to add business-level detail.
// The returned context carries the span and must be passed downstream for
// child spans to nest correctly. The caller must end the span:
//
//	ctx, span := tracer.Start(ctx, "user.Register")
//	defer span.End()
func Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return Tracer().Start(ctx, name, opts...)
}

// RecordError marks the span in ctx as failed and attaches err to it.
//
// It is a no-op when err is nil, so it can be deferred or called on a shared
// error path without a preceding nil check.
func RecordError(ctx context.Context, err error) {
	if err == nil {
		return
	}
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

// AddAttributes attaches key/value attributes to the span in ctx.
//
// This is the cheapest way to make an existing auto-instrumented span more
// useful, e.g. tagging the HTTP server span with a tenant or user id.
func AddAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		return
	}
	trace.SpanFromContext(ctx).SetAttributes(attrs...)
}

// TraceID returns the trace id of the span in ctx as a hex string, or "" when
// ctx carries no sampled span. Useful for correlating a log line or an error
// response with a trace in Jaeger.
func TraceID(ctx context.Context) string {
	sc := trace.SpanContextFromContext(ctx)
	if !sc.IsValid() {
		return ""
	}
	return sc.TraceID().String()
}
