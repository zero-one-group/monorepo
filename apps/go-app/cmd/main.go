package main

import (
	"fmt"
	"go-app/cmd/commands"
	"go-app/config"
	"log/slog"
	"os"
)

func init() {
	config.LoadEnv()
}

func main() {
	if len(os.Args) < 2 {
		slog.Error("Expected a command")
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
