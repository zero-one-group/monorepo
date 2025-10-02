package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"go-standard/internal/logging"
	"go-standard/seeders"
)

func runSeeder(db *sql.DB, target string) error {
	ctx := context.Background()
	logging.LogInfo(ctx, "Seeding target", slog.String("target", target))

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
