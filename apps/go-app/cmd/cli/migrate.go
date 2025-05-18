package cli

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/pressly/goose/v3"
)

func Migrate(db *sql.DB, dir string, mode string) error {
	err := goose.SetDialect("postgres")
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
