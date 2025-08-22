package config

import (
	"context"
	"go-app/internal/logging"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		logging.LogInfo(context.Background(), "No env file found")
	}

}
