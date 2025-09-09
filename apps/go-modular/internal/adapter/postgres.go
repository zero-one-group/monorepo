package adapter

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresDB represents the database connection pool and optionally a single connection
type PostgresDB struct {
	Pool *pgxpool.Pool // Use connection pool for most PostgreSQL operations
	Conn *pgx.Conn     // Optional: single connection, can be nil if not used
}

// Connection configuration
type PostgresConfig struct {
	URL             string
	MaxConnections  int32
	MinConnections  int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
	SearchPath      string
	Timezone        string
}

// NewPostgres creates a new database connection pool.
// Only URL is mandatory, other fields are optional and will use defaults if zero.
// PgSearchPath and PgTimezone are now configurable, default: "public" and "UTC"
func NewPostgres(cfg PostgresConfig) (*PostgresDB, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("database URL is required")
	}

	// Set defaults if zero
	if cfg.MaxConnections == 0 {
		cfg.MaxConnections = 25
	}
	if cfg.MinConnections == 0 {
		cfg.MinConnections = 5
	}
	if cfg.MaxConnLifetime == 0 {
		cfg.MaxConnLifetime = time.Hour
	}
	if cfg.MaxConnIdleTime == 0 {
		cfg.MaxConnIdleTime = 30 * time.Minute
	}
	if cfg.SearchPath == "" {
		cfg.SearchPath = "public"
	}
	if cfg.Timezone == "" {
		cfg.Timezone = "UTC"
	}

	pool, err := createPgPool(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	return &PostgresDB{Pool: pool}, nil
}

// NewPostgresWithSingleConn creates a new PostgresDB with both pool and single connection exported.
func NewPostgresWithSingleConn(cfg PostgresConfig) (*PostgresDB, error) {
	db, err := NewPostgres(cfg)
	if err != nil {
		return nil, err
	}
	conn, err := NewSingleConnection(cfg)
	if err != nil {
		return nil, err
	}
	db.Conn = conn
	return db, nil
}

// Create a single database connection, useful for one-off tasks or migrations
func NewSingleConnection(cfg PostgresConfig) (*pgx.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Configure connection settings similar to pool
	poolConfig.ConnConfig.ConnectTimeout = time.Second * 10
	poolConfig.ConnConfig.RuntimeParams = map[string]string{
		"search_path": cfg.SearchPath,
		"timezone":    cfg.Timezone,
	}

	conn, err := pgx.ConnectConfig(ctx, poolConfig.ConnConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create single connection: %w", err)
	}

	return conn, nil
}

// Ping checks if the database connection is alive
func (db *PostgresDB) Ping(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}

// Close closes all connections in the pool
func (db *PostgresDB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

// Stats returns connection pool statistics
func (db *PostgresDB) Stats() *pgxpool.Stat {
	return db.Pool.Stat()
}

// Begin starts a new transaction
func (db *PostgresDB) Begin(ctx context.Context) (pgx.Tx, error) {
	return db.Pool.Begin(ctx)
}

// BeginTx starts a new transaction with options
func (db *PostgresDB) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return db.Pool.BeginTx(ctx, txOptions)
}

// Query executes a query that returns rows
func (db *PostgresDB) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return db.Pool.Query(ctx, sql, args...)
}

// QueryRow executes a query that returns at most one row
func (db *PostgresDB) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return db.Pool.QueryRow(ctx, sql, args...)
}

// Exec executes a query that doesn't return rows
func (db *PostgresDB) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return db.Pool.Exec(ctx, sql, args...)
}

// WithTx executes a function within a transaction
func (db *PostgresDB) WithTx(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			// Only log error if it's not already closed
			fmt.Printf("failed to rollback transaction: %v\n", err)
		}
	}()

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetPoolConnection gets a connection from the pool
// Returns *pgxpool.Conn which should be released back to pool
func (db *PostgresDB) GetPoolConnection(ctx context.Context) (*pgxpool.Conn, error) {
	return db.Pool.Acquire(ctx)
}

// TestConnection tests database connectivity with timeout
func TestConnection(cfg PostgresConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := NewSingleConnection(cfg)
	if err != nil {
		return fmt.Errorf("failed to create test connection: %w", err)
	}
	defer func() {
		if err := conn.Close(ctx); err != nil {
			fmt.Printf("failed to close test connection: %v\n", err)
		}
	}()

	if err := conn.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// createPgPool creates a new pgxpool with configuration
func createPgPool(config *PostgresConfig) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Configure connection pool
	poolConfig.MaxConns = config.MaxConnections
	poolConfig.MinConns = config.MinConnections
	poolConfig.MaxConnLifetime = config.MaxConnLifetime
	poolConfig.MaxConnIdleTime = config.MaxConnIdleTime

	// Configure connection settings
	poolConfig.ConnConfig.ConnectTimeout = time.Second * 10
	poolConfig.ConnConfig.RuntimeParams = map[string]string{
		"search_path": config.SearchPath,
		"timezone":    config.Timezone,
	}

	// Create the pool
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return pool, nil
}
