package commands

import (
	"database/sql"
	"errors"
	"fmt"
	"go-app/seeders"
	"log/slog"
)

func runSeeder(db *sql.DB, target string) error {
   slog.Info("Seeding target", "target", target)

	switch target {
	case "all":
		if err := seeders.SeedUsers(db); err != nil {
			return fmt.Errorf("seeding users failed: %w", err)
		}
		// continue for other tables
	case "users":
		if err := seeders.SeedUsers(db); err != nil {
			return fmt.Errorf("seeding users failed: %w", err)
        }
    // continue for other tables
	default:
		return errors.New("unknown seed target: " + target)
	}

	return nil
}
