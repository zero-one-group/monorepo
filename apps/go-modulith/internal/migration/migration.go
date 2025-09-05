package migration

import (
	"fmt"
	"log/slog"

	"github.com/pressly/goose/v3"
	"github.com/zero-one-group/go-modulith/internal/config"
	"github.com/zero-one-group/go-modulith/internal/database"
)

type Migrator struct {
	db  *database.DB
	cfg *config.Config
}

func NewMigrator(db *database.DB, cfg *config.Config) *Migrator {
	return &Migrator{
		db:  db,
		cfg: cfg,
	}
}

func (m *Migrator) Up() error {
	sqlDB, err := m.db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(sqlDB, "./migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	slog.Info("Database migrations completed successfully")
	return nil
}

func (m *Migrator) Down() error {
	sqlDB, err := m.db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Down(sqlDB, "./migrations"); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	slog.Info("Database migration rolled back successfully")
	return nil
}

func (m *Migrator) Status() error {
	sqlDB, err := m.db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Status(sqlDB, "./migrations"); err != nil {
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	return nil
}

func (m *Migrator) Reset() error {
	sqlDB, err := m.db.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Reset(sqlDB, "./migrations"); err != nil {
		return fmt.Errorf("failed to reset migrations: %w", err)
	}

	slog.Info("Database migrations reset successfully")
	return nil
}