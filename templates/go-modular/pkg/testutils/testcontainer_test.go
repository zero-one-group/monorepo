package testutils

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

func TestPostgresConnection(t *testing.T) {
	env := NewTestEnv(t)
	pgPool, pgURL, err := env.SetupPostgres()
	if err != nil {
		t.Fatalf("failed to setup postgres: %v", err)
	}
	ctx := context.Background()

	pgConn, err := pgx.Connect(ctx, pgURL)
	if err != nil {
		t.Fatalf("failed to connect to Postgres: %v", err)
	}
	defer func() {
		if err := pgConn.Close(ctx); err != nil {
			t.Errorf("failed to close Postgres connection: %v", err)
		}
	}()
	defer pgPool.Close()

	var one int
	err = pgConn.QueryRow(ctx, "SELECT 1").Scan(&one)
	if err != nil {
		t.Fatalf("failed to query Postgres: %v", err)
	}
	if one != 1 {
		t.Errorf("unexpected Postgres query result: got %d, want 1", one)
	}
}

func TestRedisConnection(t *testing.T) {
	env := NewTestEnv(t)
	redisClient, redisAddr, err := env.SetupRedis()
	if err != nil {
		t.Fatalf("failed to setup redis: %v", err)
	}
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	defer func() {
		if err := rdb.Close(); err != nil {
			t.Errorf("failed to close Redis client: %v", err)
		}
	}()
	defer func() {
		if err := redisClient.Close(); err != nil {
			t.Errorf("failed to close testcontainer Redis client: %v", err)
		}
	}()

	status := rdb.Ping(ctx)
	if status.Err() != nil {
		t.Fatalf("failed to connect to Redis: %v", status.Err())
	}
}
