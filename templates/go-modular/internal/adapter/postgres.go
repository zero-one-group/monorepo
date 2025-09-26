package adapter

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/exaring/otelpgx"
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
	EnableOTel      bool
}

// NewPostgres creates a new database connection pool.
// Only URL is mandatory, other fields are optional and will use defaults if zero.
func NewPostgres(cfg PostgresConfig) (*PostgresDB, error) {
	// apply sensible defaults
	if cfg.MaxConnections == 0 {
		cfg.MaxConnections = 5
	}
	if cfg.MinConnections == 0 {
		cfg.MinConnections = 1
	}
	if cfg.SearchPath == "" {
		cfg.SearchPath = "public"
	}
	if cfg.Timezone == "" {
		cfg.Timezone = "UTC"
	}

	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pg pool config: %w", err)
	}

	// apply pool settings
	poolConfig.MaxConns = cfg.MaxConnections
	poolConfig.MinConns = cfg.MinConnections
	if cfg.MaxConnLifetime > 0 {
		poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
	}
	if cfg.MaxConnIdleTime > 0 {
		poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
	}

	// ensure runtime params exist and set search_path/timezone
	if poolConfig.ConnConfig.RuntimeParams == nil {
		poolConfig.ConnConfig.RuntimeParams = map[string]string{}
	}
	// prefer provided values; ensure keys set in a form Postgres accepts
	poolConfig.ConnConfig.RuntimeParams["search_path"] = cfg.SearchPath
	// Postgres accepts "TimeZone" param name; set both variants for safety
	poolConfig.ConnConfig.RuntimeParams["TimeZone"] = cfg.Timezone
	poolConfig.ConnConfig.RuntimeParams["timezone"] = cfg.Timezone

	// Initialize OpenTelemetry tracing for pgx if enabled
	// Uses otelpgx package to create a tracer that trims SQL in span names
	// and uses a custom function to generate span names from statements
	// See pgxSpanNameFunc below for details
	if cfg.EnableOTel {
		poolConfig.ConnConfig.Tracer = otelpgx.NewTracer(
			otelpgx.WithTrimSQLInSpanName(),
			otelpgx.WithSpanNameFunc(pgxSpanNameFunc),
		)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgxpool: %w", err)
	}

	return &PostgresDB{Pool: pool}, nil
}

// NewPostgresWithSingleConn creates a new PostgresDB with both pool and single connection exported.
func NewPostgresWithSingleConn(cfg PostgresConfig) (*PostgresDB, error) {
	// create pool first (will apply defaults and runtime params)
	pg, err := NewPostgres(cfg)
	if err != nil {
		return nil, err
	}

	// create a single connection using centralized helper
	conn, err := NewSingleConnection(cfg)
	if err != nil {
		pg.Close()
		return nil, fmt.Errorf("failed to create single connection: %w", err)
	}

	return &PostgresDB{
		Pool: pg.Pool,
		Conn: conn,
	}, nil
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

// CheckPostgresConnection tests database connectivity with timeout
func CheckPostgresConnection(cfg PostgresConfig) error {
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

// NewSingleConnection creates and returns a single pgx.Conn using the provided PostgresConfig.
// This is used by CheckPostgresConnection and other helpers that need a standalone connection.
func NewSingleConnection(cfg PostgresConfig) (*pgx.Conn, error) {
	// apply sensible defaults used elsewhere
	if cfg.SearchPath == "" {
		cfg.SearchPath = "public"
	}
	if cfg.Timezone == "" {
		cfg.Timezone = "UTC"
	}

	connCfg, err := pgx.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse single connection config: %w", err)
	}

	if connCfg.RuntimeParams == nil {
		connCfg.RuntimeParams = map[string]string{}
	}
	connCfg.RuntimeParams["search_path"] = cfg.SearchPath
	// set both variants for compatibility
	connCfg.RuntimeParams["TimeZone"] = cfg.Timezone
	connCfg.RuntimeParams["timezone"] = cfg.Timezone

	// set a reasonable connect timeout
	if connCfg.ConnectTimeout == 0 {
		connCfg.ConnectTimeout = 10 * time.Second
	}

	// Initialize OpenTelemetry tracing for pgx if enabled
	// Uses otelpgx package to create a tracer that trims SQL in span names
	// and uses a custom function to generate span names from statements
	// See pgxSpanNameFunc below for details
	if cfg.EnableOTel {
		connCfg.Tracer = otelpgx.NewTracer(
			otelpgx.WithTrimSQLInSpanName(),
			otelpgx.WithSpanNameFunc(pgxSpanNameFunc),
		)
	}

	conn, err := pgx.ConnectConfig(context.Background(), connCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create single connection: %w", err)
	}
	return conn, nil
}

// pgxSpanNameFunc trims "-- name: " prefix and returns the statement name for tracing
func pgxSpanNameFunc(stmt string) string {
	// If stmt is of the sqlc form "-- name: Example :one\n...",
	// extract "Example". Otherwise, leave as-is.
	stmt = strings.TrimPrefix(stmt, "-- name: ")
	if i := strings.IndexByte(stmt, ' '); i != -1 {
		return stmt[:i]
	}
	return stmt
}
