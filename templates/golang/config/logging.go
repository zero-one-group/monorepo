package config

import (
	"log/slog"
	"os"
	"strings"

	"{{package_name}}/internal/rest/middleware"

	"github.com/lmittmann/tint"
)

type LogConfig struct {
	Level       slog.Level
	Environment string
	Handler     slog.Handler
}

// GetLogLevel converts string log level to slog.Level
func GetLogLevel(level string) slog.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN", "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo // default to INFO
	}
}

// NewLogConfig creates a new logging configuration based on environment
func NewLogConfig() *LogConfig {
	env := os.Getenv("APP_ENVIRONMENT")
	if env == "" {
		env = "local"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		// Set default log levels per environment
		switch env {
		case "local", "development":
			logLevel = "DEBUG"
		case "testing":
			logLevel = "INFO"
		case "staging":
			logLevel = "WARN"
		case "production":
			logLevel = "ERROR"
		default:
			logLevel = "INFO"
		}
	}

	level := GetLogLevel(logLevel)
	handler := createHandler(env, level)

	return &LogConfig{
		Level:       level,
		Environment: env,
		Handler:     handler,
	}
}

// createHandler creates appropriate log handler based on environment
func createHandler(env string, level slog.Level) slog.Handler {
	w := os.Stdout
	opts := &slog.HandlerOptions{
		Level: level,
	}

	switch env {
	case "local", "development":
		// Use colored output for local development
		return tint.NewHandler(w, &tint.Options{
			Level:       level,
			ReplaceAttr: middleware.ColorizeLogging,
		})
	case "production":
		// Use JSON handler for production (structured logging)
		return slog.NewJSONHandler(w, opts)
	default:
		// Use text handler for other environments
		return slog.NewTextHandler(w, opts)
	}
}

// SetupLogging initializes the global logger with the configured handler
func SetupLogging() *LogConfig {
	config := NewLogConfig()
	logger := slog.New(config.Handler)
	slog.SetDefault(logger)

	slog.Info("Logging configured",
		slog.String("environment", config.Environment),
		slog.String("level", config.Level.String()),
	)

	return config
}
