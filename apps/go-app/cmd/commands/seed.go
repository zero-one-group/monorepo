package commands

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"go-app/seeders"
)

func runSeeder(db *sql.DB, target string) error {
	log.Printf("Seeding target: %s", target)

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
