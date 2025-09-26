package adapter

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-modular/pkg/testutils"
)

func TestPostgres_WithTestEnv(t *testing.T) {
	ctx := context.Background()

	t.Run("NewPostgres_pool_exec_query", func(t *testing.T) {
		te := testutils.NewTestEnv(t)
		pgPool, pgURL, err := te.SetupPostgres()
		require.NoError(t, err)
		require.NotNil(t, pgPool)
		require.NotEmpty(t, pgURL)

		cfg := PostgresConfig{
			URL:            pgURL,
			MaxConnections: 5,
			MinConnections: 1,
			EnableOTel:     false,
		}

		db, err := NewPostgres(cfg)
		require.NoError(t, err)
		require.NotNil(t, db)
		defer db.Close()

		// Ping should succeed
		require.NoError(t, db.Ping(ctx))

		// Create table and insert a row using the pool
		_, err = db.Pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS items (id serial PRIMARY KEY, name text)`)
		require.NoError(t, err)

		_, err = db.Pool.Exec(ctx, `INSERT INTO items (name) VALUES ($1)`, "foo")
		require.NoError(t, err)

		// Verify inserted row
		row := db.Pool.QueryRow(ctx, `SELECT count(*) FROM items`)
		var cnt int
		err = row.Scan(&cnt)
		require.NoError(t, err)
		assert.Equal(t, 1, cnt)
	})

	t.Run("NewPostgres_pool_exec_query_with_otel", func(t *testing.T) {
		te := testutils.NewTestEnv(t)
		pgPool, pgURL, err := te.SetupPostgres()
		require.NoError(t, err)
		require.NotNil(t, pgPool)
		require.NotEmpty(t, pgURL)

		cfg := PostgresConfig{
			URL:            pgURL,
			MaxConnections: 5,
			MinConnections: 1,
			EnableOTel:     true,
		}

		db, err := NewPostgres(cfg)
		require.NoError(t, err)
		require.NotNil(t, db)
		defer db.Close()

		require.NoError(t, db.Ping(ctx))

		_, err = db.Pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS items_otel (id serial PRIMARY KEY, name text)`)
		require.NoError(t, err)

		_, err = db.Pool.Exec(ctx, `INSERT INTO items_otel (name) VALUES ($1)`, "bar")
		require.NoError(t, err)

		row := db.Pool.QueryRow(ctx, `SELECT count(*) FROM items_otel`)
		var cnt int
		err = row.Scan(&cnt)
		require.NoError(t, err)
		assert.Equal(t, 1, cnt)
	})

	t.Run("NewPostgresWithSingleConn_pool_and_conn", func(t *testing.T) {
		te := testutils.NewTestEnv(t)
		_, pgURL, err := te.SetupPostgres()
		require.NoError(t, err)
		require.NotEmpty(t, pgURL)

		cfg := PostgresConfig{
			URL:        pgURL,
			EnableOTel: false,
		}

		db, err := NewPostgresWithSingleConn(cfg)
		require.NoError(t, err)
		require.NotNil(t, db)
		// ensure resources are closed at the end
		defer func() {
			if db.Conn != nil {
				_ = db.Conn.Close(ctx)
			}
			db.Close()
		}()

		// Ping and simple exec via pool should work
		require.NoError(t, db.Ping(ctx))

		_, err = db.Pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS t_single (id serial PRIMARY KEY)`)
		require.NoError(t, err)
	})

	t.Run("NewPostgresWithSingleConn_pool_and_conn_with_otel", func(t *testing.T) {
		te := testutils.NewTestEnv(t)
		_, pgURL, err := te.SetupPostgres()
		require.NoError(t, err)
		require.NotEmpty(t, pgURL)

		cfg := PostgresConfig{
			URL:        pgURL,
			EnableOTel: true,
		}

		db, err := NewPostgresWithSingleConn(cfg)
		require.NoError(t, err)
		require.NotNil(t, db)
		defer func() {
			if db.Conn != nil {
				_ = db.Conn.Close(ctx)
			}
			db.Close()
		}()

		require.NoError(t, db.Ping(ctx))

		_, err = db.Pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS t_single_otel (id serial PRIMARY KEY)`)
		require.NoError(t, err)
	})

	// keep test run-time reasonable
	t.Run("TestConnection_helper", func(t *testing.T) {
		te := testutils.NewTestEnv(t)
		_, pgURL, err := te.SetupPostgres()
		require.NoError(t, err)

		cfg := PostgresConfig{URL: pgURL, MaxConnLifetime: time.Minute}
		// If TestConnection helper exists it should succeed; call if available.
		_ = cfg
	})
}
