package config

import (
	"log/slog"
	"os"
	"time"
)

func LoadJWTConfig() (string, time.Duration) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		slog.Error("Missing JWT_SECRET environment variable")
		os.Exit(1)
	}

	jwtTTLStr := os.Getenv("JWT_TTL")
	if jwtTTLStr == "" {
		slog.Error("Missing JWT_TTL environment variable")
		os.Exit(1)
	}

	jwtTTL, err := time.ParseDuration(jwtTTLStr)
	if err != nil {
		slog.Error("Invalid JWT_TTL value", "error", err)
		os.Exit(1)
	}

	return jwtSecret, jwtTTL
}
