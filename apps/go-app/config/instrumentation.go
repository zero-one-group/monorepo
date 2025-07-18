package config

import (
	"context"
	"fmt"
	"go-app/internal/metrics"
	"go-app/internal/rest/middleware"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// ApplyInstrumentation configures Prometheus metrics and OpenTelemetry tracing
// for the Echo instance. It returns a shutdown function for the tracer provider
// and an error if the setup fails.
func ApplyInstrumentation(
	ctx context.Context,
	e *echo.Echo,
	appMetrics *metrics.Metrics,
) (func(context.Context) error, error) {
	enableInstrumentationStr := os.Getenv("ENABLE_INSTRUMENTATION")
	enableInstrumentation, err := strconv.ParseBool(enableInstrumentationStr)
	if err != nil {
		enableInstrumentation = false
		fmt.Printf("Warning: ENABLE_INSTRUMENTATION environment variable '%s' could not be parsed as boolean. Defaulting to false. Error: %v\n", enableInstrumentationStr, err)
	}

	if !enableInstrumentation {
		fmt.Println("Instrumentation is disabled by ENABLE_INSTRUMENTATION environment variable.")
		return func(context.Context) error { return nil }, nil
	}

	fmt.Println("Instrumentation is enabled.")

	// Apply Prometheus middleware and metrics endpoint.
	err = initMetrics(e, appMetrics)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize metrics: %w", err)
	}

	// Initialize the OpenTelemetry tracer provider.
	tp, shutdown, err := initTracer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tracer: %w", err)
	}

	e.Use(middleware.AttachTraceProvider(tp))
	return shutdown, nil
}

// initMetrics configures Prometheus metrics for the Echo instance.
func initMetrics(e *echo.Echo, appMetrics *metrics.Metrics) error {
	serviceName := os.Getenv("SERVICE_NAME")
	// @see: https://echo.labstack.com/docs/middleware/prometheus#custom-configuration
	e.Use(echoprometheus.NewMiddleware(serviceName))
	e.GET("/metrics", echoprometheus.NewHandler())

	// Register all custom metrics from the appMetrics struct
	if err := prometheus.Register(appMetrics.UserRepoCalls); err != nil {
		return fmt.Errorf("failed to register UserRepoCalls metric: %w", err)
	}
	// Register other metrics here if you add them to your metrics.Metrics struct
	// if err := prometheus.Register(appMetrics.OtherMetric); err != nil {
	// 	return fmt.Errorf("failed to register OtherMetric: %w", err)
	// }

	fmt.Println("Prometheus metrics initialized and registered.")
	return nil
}

// initTracer initializes an OTel tracer provider. In non-production
// environments, it uses a no-op provider. It returns a shutdown function
// and an error.
func initTracer(
	ctx context.Context,
) (*sdktrace.TracerProvider, func(context.Context) error, error) {
	env := os.Getenv("APP_ENVIRONMENT")

	// In non-prod, install a NeverSample provider so no spans are exported.
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
		// Return a no-op shutdown function
		return tp, func(context.Context) error { return nil }, nil
	}

	// In production, configure a real OTLP/gRPC exporter.
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
		return nil, nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
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
		return nil, nil, fmt.Errorf("failed to create OTel resources: %w", err)
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
		fmt.Println("Shutting down OpenTelemetry tracer provider...")
		return tp.Shutdown(ctx)
	}

	return tp, shutdown, nil
}

