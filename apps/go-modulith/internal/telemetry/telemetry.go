package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/zero-one-group/go-modulith/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.31.0"
)

type Telemetry struct {
	tracerProvider *trace.TracerProvider
}

func NewTelemetry(cfg *config.Config) (*Telemetry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.OpenTelemetry.ServiceName),
			semconv.ServiceVersion(cfg.OpenTelemetry.ServiceVersion),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	var exporter trace.SpanExporter

	if cfg.OpenTelemetry.Endpoint != "" {
		exporter, err = otlptracegrpc.New(ctx,
			otlptracegrpc.WithEndpoint(cfg.OpenTelemetry.Endpoint),
			otlptracegrpc.WithInsecure(),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
		}
	}

	var spanProcessors []trace.SpanProcessor
	if exporter != nil {
		spanProcessors = append(spanProcessors, trace.NewBatchSpanProcessor(exporter))
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithSpanProcessor(spanProcessors[0]),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	slog.Info("OpenTelemetry initialized", 
		"service", cfg.OpenTelemetry.ServiceName,
		"version", cfg.OpenTelemetry.ServiceVersion,
		"endpoint", cfg.OpenTelemetry.Endpoint)

	return &Telemetry{
		tracerProvider: tracerProvider,
	}, nil
}

func (t *Telemetry) Shutdown(ctx context.Context) error {
	if t.tracerProvider != nil {
		return t.tracerProvider.Shutdown(ctx)
	}
	return nil
}

func Tracer(name string) {
	otel.Tracer(name)
}