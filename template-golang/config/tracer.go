package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// InitTracer returns a TracerProvider, a shutdown func, and an error.
// In non-prod (ENVIRONMENT != "production") it installs a NeverSample
// provider so no spans are ever exported.
func InitTracer(ctx context.Context) (
	*sdktrace.TracerProvider,
	func(context.Context) error,
) {
	env := os.Getenv("APP_ENVIRONMENT")

	// 1) Non-prod: no-op tracer
	if env != "production" {
		tp := sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.NeverSample()),
		)
		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(
			propagation.NewCompositeTextMapPropagator(
				propagation.TraceContext{},
				propagation.Baggage{},
			),
		)
		return tp, func(context.Context) error { return nil }
	}

	// 2) Production: real OTLP/gRPC exporter
	serviceName := os.Getenv("SERVICE_NAME")
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:4317"
	}

	exporter, err := otlptrace.New(
		ctx,
		otlptracegrpc.NewClient(
			otlptracegrpc.WithEndpoint(endpoint),
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithTimeout(5*time.Second),
		),
	)
	if err != nil {
		fmt.Printf("Err creating OTLP exporter: %w\n", err)
		os.Exit(1)
		return nil, nil
	}

	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("deployment.environment", env),
			attribute.String("telemetry.sdk.language", "go"),
		),
		resource.WithHost(),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		fmt.Printf("Err creating resources: %w\n", err)
		os.Exit(1)
		return nil, nil
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	shutdown := func(ctx context.Context) error {
		fmt.Println("shutting down tracer provider")
		return tp.Shutdown(ctx)
	}

	return tp, shutdown
}

