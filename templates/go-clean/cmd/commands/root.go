package commands

import (
	"errors"
	"fmt"
	"{{ package_name | kebab_case }}/database"
)

func Execute(command string, args []string) error {
	subcommand := ""
	if len(args) > 0 {
		subcommand = args[0]
	}

	db, err := database.SetupSQLDatabase()
	if err != nil {
		return fmt.Errorf("failed to connect to DB: %w", err)
	}
	defer db.Close()

	switch command {
	case "migrate":
		dir := "./migrations"
		if err := runMigration(db, dir, args); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	case "seed":
		target := "all"
		if subcommand != "" {
			target = subcommand
		}
		if err := runSeeder(db, target); err != nil {
			return fmt.Errorf("seeding failed: %w", err)
		}
	default:
		return errors.New("unknown command: " + command)
	}

	return nil
}
