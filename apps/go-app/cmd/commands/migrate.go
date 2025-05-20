package commands

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/pressly/goose/v3"
)


func runMigration(db *sql.DB, dir string, args []string) error {

    err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

    mode := args[0]

	switch mode {
    case "create":
        if len(args) < 2 {
			return errors.New("migration name is required for 'create' command")
		}
        migrationName := args[1]
        err = goose.Create(db, dir, migrationName, "sql")
	case "up":
		err = goose.Up(db, dir)
	case "down":
		err = goose.Down(db, dir)
	case "reset":
		err = goose.Reset(db, dir)
    case "version":
        version := goose.Version(db, dir)
		fmt.Printf("Current migration version: %d\n", version)
	default:
		err = errors.New(mode + " is not Migrate function")
	}
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}
