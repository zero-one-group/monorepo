package main

import (
	"fmt"
	"log/slog"
	"os"
	"{{package_name}}/cmd/commands"
	"{{package_name}}/config"
)

func init() {
	config.LoadEnv()
}

func main() {
	if len(os.Args) < 2 {
		slog.Error("Expected a command")
		fmt.Fprintf(os.Stderr, "Usage: %s <command> [args...]\n", os.Args[0])
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	err := commands.Execute(command, args)
	if err != nil {
		slog.Error("Command execution failed",
			slog.String("command", command),
			slog.String("error", err.Error()),
		)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
