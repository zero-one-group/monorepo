package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/zero-one-group/go-modulith/internal/config"
	"go.opentelemetry.io/otel/trace"
)

func NewLogger(cfg *config.Config) *slog.Logger {
	var handler slog.Handler

	level := parseLogLevel(cfg.Logging.Level)
	opts := &slog.HandlerOptions{
		Level: level,
		AddSource: cfg.IsDevelopment(),
	}

	if cfg.Logging.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(&TracingHandler{Handler: handler})
	slog.SetDefault(logger)

	return logger
}

func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type TracingHandler struct {
	slog.Handler
}

func (h *TracingHandler) Handle(ctx context.Context, record slog.Record) error {
	if span := trace.SpanFromContext(ctx); span.IsRecording() {
		spanCtx := span.SpanContext()
		if spanCtx.HasTraceID() {
			record.Add("trace_id", spanCtx.TraceID().String())
		}
		if spanCtx.HasSpanID() {
			record.Add("span_id", spanCtx.SpanID().String())
		}
	}

	return h.Handler.Handle(ctx, record)
}

func (h *TracingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &TracingHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *TracingHandler) WithGroup(name string) slog.Handler {
	return &TracingHandler{Handler: h.Handler.WithGroup(name)}
}

func FromContext(ctx context.Context) *slog.Logger {
	logger := slog.Default()
	
	if span := trace.SpanFromContext(ctx); span.IsRecording() {
		spanCtx := span.SpanContext()
		attrs := make([]slog.Attr, 0, 2)
		
		if spanCtx.HasTraceID() {
			attrs = append(attrs, slog.String("trace_id", spanCtx.TraceID().String()))
		}
		if spanCtx.HasSpanID() {
			attrs = append(attrs, slog.String("span_id", spanCtx.SpanID().String()))
		}
		
		if len(attrs) > 0 {
			anyAttrs := make([]any, len(attrs))
			for i, attr := range attrs {
				anyAttrs[i] = attr
			}
			logger = logger.With(anyAttrs...)
		}
	}
	
	return logger
}