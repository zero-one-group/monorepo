package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

// MigrationFS contains the embedded database migration files.
//
//go:embed migrations/*.sql
var MigrationFS embed.FS

type Migrator struct {
	db   *sql.DB
	pool *pgxpool.Pool
}

// createConnection creates a new database connection pool and a SQL database instance with retry mechanism
func createConnection(dbURL string) (*pgxpool.Pool, *sql.DB, error) {
	const baseDelay = 2 * time.Second
	const maxDelay = 30 * time.Second
	const maxRetries = 5

	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		poolConfig, err := pgxpool.ParseConfig(dbURL)
		if err != nil {
			lastErr = fmt.Errorf("failed to parse database URL on attempt %d: %w", attempt, err)
			slog.Warn("Database connection attempt failed", "attempt", attempt)
			if attempt == maxRetries {
				break
			}
			delay := min(time.Duration(attempt-1)*baseDelay, maxDelay)
			time.Sleep(delay)
			continue
		}

		poolConfig.ConnConfig.ConnectTimeout = time.Second * 10
		poolConfig.ConnConfig.RuntimeParams = map[string]string{
			"search_path": "public,reference,scheduler",
			"timezone":    "UTC",
		}

		pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err != nil {
			lastErr = fmt.Errorf("%s", err.Error())
			slog.Warn("Database connection attempt failed", "attempt", attempt)
			if attempt == maxRetries {
				break
			}
			delay := min(time.Duration(attempt-1)*baseDelay, maxDelay)
			time.Sleep(delay)
			continue
		}

		db := stdlib.OpenDBFromPool(pool)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err = db.PingContext(ctx)
		cancel()
		if err != nil {
			lastErr = fmt.Errorf("%s", err.Error())
			pool.Close()
			slog.Warn("Database connection attempt failed", "attempt", attempt)
			if attempt == maxRetries {
				break
			}
			delay := min(time.Duration(attempt-1)*baseDelay, maxDelay)
			time.Sleep(delay)
			continue
		}

		slog.Info("Database connection established", "attempt", attempt)
		return pool, db, nil
	}

	return nil, nil, fmt.Errorf("failed to establish database connection after %d attempts: %w", maxRetries, lastErr)
}

// min helper for durations
func min(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

func NewMigrator(dbURL string) *Migrator {
	pool, db, err := createConnection(dbURL)
	if err != nil {
		slog.Error("Failed to create database connection", "err", err.Error())
		os.Exit(1) // Exit with error code, no panic/stack trace
	}

	// Configure goose migrator
	_ = goose.SetDialect("postgres")
	goose.SetTableName("app_migrations")
	goose.SetBaseFS(MigrationFS)
	goose.SetSequential(true)

	slog.Info("Database migrator initialized")
	return &Migrator{db: db, pool: pool}
}

func (m *Migrator) MigrateUp(ctx context.Context) error {
	slog.Info("Applying Goose migrations", "direction", "up")
	if err := goose.UpContext(ctx, m.db, "migrations"); err != nil {
		slog.Error("Goose migration up failed", "err", err)
		return fmt.Errorf("goose up: %w", err)
	}

	slog.Info("All migrations applied successfully")
	return nil
}

func (m *Migrator) MigrateDown(ctx context.Context, steps string) error {
	slog.Info("Rolling back Goose migrations", "direction", "down", "steps", steps)
	downSteps := 1
	if steps != "" {
		n, err := strconv.Atoi(steps)
		if err != nil || n < 1 {
			slog.Error("Invalid steps argument for migration down", "steps", steps, "err", err)
			return fmt.Errorf("invalid steps argument: %w", err)
		}
		downSteps = n
	}

	for i := 0; i < downSteps; i++ {
		if err := goose.Down(m.db, "migrations"); err != nil {
			slog.Warn("No more Goose migrations to rollback", "iteration", i+1)
			break
		}
		slog.Info("Goose migration rolled back", "iteration", i+1)
	}

	if err := m.MigrateStatus(ctx); err != nil {
		slog.Error("Failed to get migration status after down", "err", err)
	}

	return nil
}

func (m *Migrator) MigrateReset(ctx context.Context) error {
	if err := m.MigrateStatus(ctx); err != nil {
		slog.Error("Failed to get migration status after down", "err", err)
	}

	slog.Info("Resetting all Goose migrations")
	if err := goose.Reset(m.db, "migrations"); err != nil {
		slog.Error("Failed to reset Goose migrations", "err", err)
		panic(fmt.Errorf("failed to reset migrations: %w", err))
	}

	return nil
}

func (m *Migrator) MigrateStatus(ctx context.Context) error {
	slog.Info("Checking current Goose migration status")
	if err := goose.Status(m.db, "migrations"); err != nil {
		slog.Error("Failed to show Goose migration status", "err", err)
		panic(fmt.Errorf("failed to show migration status: %w", err))
	}
	return nil
}

func (m *Migrator) MigrateVersion(ctx context.Context) error {
	slog.Info("Checking current Goose migration version")
	if err := goose.Version(m.db, "migrations"); err != nil {
		slog.Error("Failed to show Goose migration version", "err", err)
		panic(fmt.Errorf("failed to show migration version: %w", err))
	}
	return nil
}

func (m *Migrator) MigrateCreate(ctx context.Context, migrationName string) error {
	if err := m.MigrateStatus(ctx); err != nil {
		slog.Error("Failed to show migration status before create", "err", err)
	}

	execDir, err := os.Getwd() // Use path relative to the current file location
	if err != nil {
		slog.Error("Failed to get working directory", "err", err)
		return err
	}
	migrationsDir := filepath.Join(execDir, "migrations")
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		// Try relative to this source file (for go run or test)
		migrationsDir = filepath.Join(execDir, "database", "migrations")
		if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
			slog.Error("Migrations directory does not exist", "checked", []string{
				filepath.Join(execDir, "migrations"),
				filepath.Join(execDir, "database", "migrations"),
			})
			return fmt.Errorf("migrations directory does not exist")
		}
	}

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		slog.Error("Failed to read migrations directory", "err", err)
		return err
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		parts := strings.SplitN(name, "_", 2)
		if len(parts) != 2 {
			continue
		}
		if strings.TrimSuffix(parts[1], ".sql") == migrationName {
			slog.Warn("Migration already exists", "migration", name)
			return fmt.Errorf("migration already exists: %s", name)
		}
	}

	slog.Info("Creating new Goose migration file", "name", migrationName)
	if err := goose.Create(m.db, migrationsDir, migrationName, "sql"); err != nil {
		slog.Error("Failed to create migration file", "err", err)
		return err
	}

	return nil
}

func (m *Migrator) Close() error {
	if m.db != nil {
		slog.Info("Closing database connection")
		_ = m.db.Close()
		m.db = nil // prevent double close/log
	}
	return nil
}
