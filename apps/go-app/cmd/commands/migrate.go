package commands

import (
	"errors"
	"fmt"
	"log"
	"go-app/database"

	"github.com/pressly/goose/v3"
)

func runMigration(dir string, mode string) error {

    db, err := database.SetupSQLDatabase()
	if err != nil {
		log.Fatal("Failed to set up database: " + err.Error())
	}
	defer db.Close()

	err = goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	switch mode {
	case "up":
		err = goose.Up(db, dir)
	case "down":
		err = goose.Down(db, dir)
	case "reset":
		err = goose.Reset(db, dir)
	default:
		err = errors.New(mode + " is not Migrate function")
	}
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}
