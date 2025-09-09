package main

import (
	"fmt"
	"os"

	"go-modular/cmd/commands"
)

// @title	    Go Application API
// @description	Go Application API documentation
// @version		1.0

// @securityDefinitions.http bearerAuth
// @scheme bearer
// @bearerFormat JWT

// @BasePath /api
func main() {
	if err := setupConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setupConfig() error {
	_, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// TODO: Do something!

	return nil
}
