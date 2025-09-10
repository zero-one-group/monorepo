package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

type LoggerOpts struct {
	// Level is the concrete slog.Level to use (optional; defaults to INFO)
	Level slog.Level

	// Format controls output format: "json" or "pretty" (defaults to "pretty")
	Format string

	// NoColor disables colorized output for pretty formatter
	NoColor bool

	// Environment is an informational label (e.g. "development", "production")
	Environment string
}

type AppLogger struct {
	Handler     slog.Handler
	Logger      *slog.Logger
	Level       slog.Level
	Environment string
}

// SetupLogging initializes the global logger using provided options
// and returns the created *slog.Logger.
func SetupLogging(opts LoggerOpts) *slog.Logger {
	lc := NewAppLogger(opts)

	// set global default logger
	slog.SetDefault(lc.Logger)

	slog.Info("Logging configured",
		slog.String("environment", lc.Environment),
		slog.String("level", lc.Level.String()),
	)

	return lc.Logger
}

// NewAppLogger builds a AppLogger from LoggerOpts without depending on app config package.
func NewAppLogger(opts LoggerOpts) *AppLogger {
	env := opts.Environment
	if env == "" {
		env = "development"
	}

	// determine log level (default to INFO)
	level := slog.LevelInfo
	if opts.Level != 0 {
		level = opts.Level
	}

	// determine format and color preference (defaults)
	format := opts.Format
	if format == "" {
		format = "pretty"
	}
	noColor := opts.NoColor

	handler := createHandler(env, format, noColor, level)
	logger := slog.New(handler)

	return &AppLogger{
		Handler:     handler,
		Logger:      logger,
		Level:       level,
		Environment: env,
	}
}

// createHandler creates appropriate log handler based on requested format and environment.
// format: "json" => JSON structured output
//
//	"pretty" (default) => colorized human-friendly output in dev (tint) or plain text when noColor=true
func createHandler(env, format string, noColor bool, level slog.Level) slog.Handler {
	w := os.Stdout
	opts := &slog.HandlerOptions{
		Level: level,
	}

	// JSON takes precedence when requested
	if format == "json" {
		return slog.NewJSONHandler(w, opts)
	}

	// pretty format
	// in development/local prefer tint (colorized) unless noColor is true
	if format == "pretty" {
		if (env == "local" || env == "development") && !noColor {
			return tint.NewHandler(w, &tint.Options{
				Level:       level,
				ReplaceAttr: colorizeSlog,
				AddSource:   level == slog.LevelDebug,
				TimeFormat:  time.Stamp,
				NoColor:     noColor,
			})
		}
		// fallback to text handler (no colors / production)
		return slog.NewTextHandler(w, opts)
	}

	// unknown format -> fallback: production-friendly JSON when not local; otherwise text
	if env == "production" {
		return slog.NewJSONHandler(w, opts)
	}
	return slog.NewTextHandler(w, opts)
}

func colorizeSlog(groups []string, a slog.Attr) slog.Attr {
	// Only transform top-level level attribute
	if a.Key == slog.LevelKey && len(groups) == 0 {
		if level, ok := a.Value.Any().(slog.Level); ok {
			switch level {
			case slog.LevelError:
				return tint.Attr(9, slog.String(a.Key, "ERR"))
			case slog.LevelWarn:
				return tint.Attr(12, slog.String(a.Key, "WRN"))
			case slog.LevelInfo:
				return tint.Attr(10, slog.String(a.Key, "INF"))
			case slog.LevelDebug:
				return tint.Attr(10, slog.String(a.Key, "DBG"))
			}
		}
	}
	return a
}
