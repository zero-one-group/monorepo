package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func SetupPgxPool() (*pgxpool.Pool, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable not set")
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}
	config.ConnConfig.Tracer = otelpgx.NewTracer()

	if err := otelpgx.RecordStats(dbPool); err != nil {
		return nil, fmt.Errorf("unable to record database stats: %w", err)
	}

	if err := dbPool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Connected to DB Postgresql...")

	return dbPool, nil
}

func SetupSQLDatabase() (*sql.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable not set")
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, fmt.Errorf("sql open error: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("sql ping error: %w", err)
	}

	fmt.Println("Connected to DB via *sql.DB...")
	return db, nil
}
