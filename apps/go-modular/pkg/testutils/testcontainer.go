package testutils

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestEnv groups all test dependencies for integration tests.
type TestEnv struct {
	T           *testing.T
	Ctx         context.Context
	PGPool      *pgxpool.Pool
	RedisClient *redis.Client
	PGURL       string
	RedisAddr   string

	postgresC testcontainers.Container
	redisC    testcontainers.Container
}

// NewTestEnv returns a new TestEnv.
func NewTestEnv(t *testing.T) *TestEnv {
	return &TestEnv{
		T:   t,
		Ctx: context.Background(),
	}
}

// SetupPostgres starts a Postgres container and initializes the PGPool.
func (te *TestEnv) SetupPostgres() (*pgxpool.Pool, string, error) {
	t := te.T
	ctx := te.Ctx

	var err error
	te.postgresC, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:18rc1-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_PASSWORD": "testpass",
				"POSTGRES_USER":     "testuser",
				"POSTGRES_DB":       "testdb",
			},
			WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(120 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		return nil, "", fmt.Errorf("start postgres container: %w", err)
	}
	t.Cleanup(func() { _ = te.postgresC.Terminate(ctx) })

	pgEndpoint, err := te.postgresC.Endpoint(ctx, "")
	if err != nil {
		return nil, "", fmt.Errorf("get postgres endpoint: %w", err)
	}
	pgHost, pgPortStr, err := net.SplitHostPort(strings.TrimPrefix(pgEndpoint, "tcp://"))
	if err != nil {
		return nil, "", fmt.Errorf("parse postgres endpoint: %w", err)
	}
	pgURL := fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", pgHost, pgPortStr)

	poolConfig, err := pgxpool.ParseConfig(pgURL)
	if err != nil {
		return nil, "", fmt.Errorf("parse pg config: %w", err)
	}
	poolConfig.MaxConns = 5
	poolConfig.MinConns = 1
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.ConnConfig.RuntimeParams["search_path"] = "public,reference,scheduler"
	poolConfig.ConnConfig.RuntimeParams["timezone"] = "UTC"

	pgPool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, "", fmt.Errorf("create pg pool: %w", err)
	}
	t.Cleanup(func() { pgPool.Close() })

	te.PGPool = pgPool
	te.PGURL = pgURL
	return pgPool, pgURL, nil
}

// SetupRedis starts a Redis container and initializes the RedisClient.
func (te *TestEnv) SetupRedis() (*redis.Client, string, error) {
	t := te.T
	ctx := te.Ctx

	var err error
	te.redisC, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "redis:8-alpine",
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor:   wait.ForListeningPort("6379/tcp").WithStartupTimeout(60 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		return nil, "", fmt.Errorf("start redis container: %w", err)
	}
	t.Cleanup(func() { _ = te.redisC.Terminate(ctx) })

	redisEndpoint, err := te.redisC.Endpoint(ctx, "")
	if err != nil {
		return nil, "", fmt.Errorf("get redis endpoint: %w", err)
	}
	redisHost, redisPortStr, err := net.SplitHostPort(strings.TrimPrefix(redisEndpoint, "tcp://"))
	if err != nil {
		return nil, "", fmt.Errorf("parse redis endpoint: %w", err)
	}
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPortStr)

	redisClient := redis.NewClient(&redis.Options{Addr: redisAddr, DB: 15})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, "", fmt.Errorf("redis ping: %w", err)
	}
	t.Cleanup(func() {
		_ = redisClient.Close()
	})

	te.RedisClient = redisClient
	te.RedisAddr = redisAddr
	return redisClient, redisAddr, nil
}

// SetupConfig sets the environment variables for the test configuration.
func (te *TestEnv) SetupConfig() {
	if te.PGURL != "" {
		if err := os.Setenv("DATABASE_URL", te.PGURL); err != nil {
			te.T.Logf("Failed to set DATABASE_URL: %v", err)
		}
	}
	if te.RedisAddr != "" {
		if err := os.Setenv("REDIS_URL", fmt.Sprintf("redis://%s", te.RedisAddr)); err != nil {
			te.T.Logf("Failed to set REDIS_URL: %v", err)
		}
		if err := os.Setenv("REDIS_ENABLED", "true"); err != nil {
			te.T.Logf("Failed to set REDIS_ENABLED: %v", err)
		}
	}
	if err := os.Setenv("LOG_LEVEL", "debug"); err != nil {
		te.T.Logf("Failed to set LOG_LEVEL: %v", err)
	}
}

// RunAppMigrations runs the database migrations for the application.
func (te *TestEnv) RunAppMigrations() {
	cmd := exec.Command("go", "run", "-tags", "debug", "./cmd/", "migrate:reset", "--force", "--up", "--seed")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	// ensure command runs from module root (where go.mod lives)
	if root := findModuleRoot(); root != "" {
		cmd.Dir = root
		te.T.Logf("running migrations from module root: %s", root)
	} else {
		te.T.Log("module root not found, running migrations from current working directory")
	}

	te.T.Log("Running test database migrations")
	err := cmd.Run()
	require.NoError(te.T, err, "failed to run test database migrations")
}

// findModuleRoot walks up from the current working directory to find the directory containing go.mod.
func findModuleRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// reached filesystem root
			return ""
		}
		dir = parent
	}
}
