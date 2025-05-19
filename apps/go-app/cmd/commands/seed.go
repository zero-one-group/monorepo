package commands

import (
	"errors"
	"fmt"
	"log"
	"go-app/database"
	"go-app/seeders"
)

func runSeeder(target string) error {
	log.Printf("Seeding target: %s", target)
    db, err := database.SetupSQLDatabase()
	if err != nil {
		log.Fatal("Failed to set up database: " + err.Error())
	}
	defer db.Close()

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
