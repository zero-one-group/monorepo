package config

import (
	"log/slog"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		slog.Info("No env file found")
	}

}
