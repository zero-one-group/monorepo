package main

import (
	"context"
	"fmt"
	"os"
	"go-app/cmd/commands"
	"go-app/config"
	"go-app/internal/logging"
)

func init() {
	config.LoadEnv()
}

func main() {
	if len(os.Args) < 2 {
		logging.LogErrorMessage(context.Background(), "Expected a command")
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	err := commands.Execute(command, args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
