package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"{{ package_name | kebab_case }}/internal/adapter"
	"{{ package_name | kebab_case }}/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

var postgresModule = fx.Module("postgres",
	fx.Provide(newPostgresDB),
	fx.Provide(func(db *adapter.PostgresDB) *pgxpool.Pool { return db.Pool }),
)

func newPostgresDB(lc fx.Lifecycle, cfg *config.Config, log *slog.Logger) (*adapter.PostgresDB, error) {
	db, err := connectPostgresWithRetry(cfg, log)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			log.Info("Closing database connections")
			db.Close()
			return nil
		},
	})
	return db, nil
}

func connectPostgresWithRetry(cfg *config.Config, log *slog.Logger) (*adapter.PostgresDB, error) {
	const baseDelay = 2 * time.Second
	const maxDelay = 30 * time.Second
	const defaultMaxRetries = 5

	maxRetries := defaultMaxRetries
	if cfg.Database.PgMaxRetries == -1 {
		maxRetries = -1
	} else if cfg.Database.PgMaxRetries > 0 {
		maxRetries = cfg.Database.PgMaxRetries
	}

	log.Info("Initializing database connection", "max_retries", maxRetries)

	var lastErr error
	attempt := 1

	for {
		pg, err := adapter.NewPostgres(adapter.PostgresConfig{URL: cfg.GetDatabaseURL()})
		if err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			pingErr := pg.Ping(ctx)
			cancel()
			if pingErr == nil {
				log.Info("Database connection established", "attempt", attempt)
				return pg, nil
			}
			pg.Close()
			err = fmt.Errorf("ping failed: %w", pingErr)
		}

		lastErr = err
		log.Warn("Database connection attempt failed", "attempt", attempt, "err", lastErr)

		if maxRetries != -1 && attempt >= maxRetries {
			break
		}

		delay := baseDelay * time.Duration(attempt)
		if delay > maxDelay {
			delay = maxDelay
		}
		log.Info("Retrying database connection", "next_try_in", delay, "attempt", attempt+1)
		time.Sleep(delay)
		attempt++
	}

	return nil, fmt.Errorf("failed to establish database connection after %d attempts: %w", attempt, lastErr)
}
