package commands

import (
	"errors"
	"fmt"
)


func Execute(command string, args []string) error {
    subcommand := ""
	if len(args) > 0 {
		subcommand = args[0]
	}

    switch command {
	case "migrate":
		dir := "./migrations"
		if err := runMigration(dir, subcommand); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	case "seed":
		target := "all"
		if subcommand != "" {
			target = subcommand
		}
		if err := runSeeder(target); err != nil {
			return fmt.Errorf("seeding failed: %w", err)
		}
	default:
		return errors.New("unknown command: " + command)
	}

	return nil
}
