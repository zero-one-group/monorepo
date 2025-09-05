package app

import (
	"context"
	"log/slog"

	"github.com/zero-one-group/go-modulith/internal/config"
	"github.com/zero-one-group/go-modulith/internal/database"
	"github.com/zero-one-group/go-modulith/internal/logger"
	"github.com/zero-one-group/go-modulith/internal/migration"
	"github.com/zero-one-group/go-modulith/internal/module/auth"
	"github.com/zero-one-group/go-modulith/internal/module/product"
	"github.com/zero-one-group/go-modulith/internal/module/user"
	"github.com/zero-one-group/go-modulith/internal/telemetry"
	"github.com/zero-one-group/go-modulith/internal/validator"
	"go.uber.org/fx"
)

func NewApp() *fx.App {
	return fx.New(
		// Configuration
		fx.Provide(config.Load),

		// Infrastructure
		fx.Provide(logger.NewLogger),
		fx.Provide(database.NewDatabase),
		fx.Provide(telemetry.NewTelemetry),
		fx.Provide(validator.New),
		fx.Provide(migration.NewMigrator),

		// Repositories
		fx.Provide(auth.NewRepository),
		fx.Provide(user.NewRepository),
		fx.Provide(product.NewRepository),

		// Services
		fx.Provide(auth.NewService),
		fx.Provide(user.NewService),
		fx.Provide(product.NewService),

		// Handlers
		fx.Provide(auth.NewHandler),
		fx.Provide(user.NewHandler),
		fx.Provide(product.NewHandler),

		// Server
		fx.Provide(NewServer),

		// Lifecycle
		fx.Invoke(RegisterHooks),
	)
}

func RegisterHooks(
	lc fx.Lifecycle,
	server *Server,
	db *database.DB,
	tel *telemetry.Telemetry,
	log *slog.Logger,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Starting application...")
			return server.Start()
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Shutting down application...")

			if err := server.Shutdown(ctx); err != nil {
				log.Error("Failed to shutdown server gracefully", "error", err)
			}

			if err := db.Close(); err != nil {
				log.Error("Failed to close database connection", "error", err)
			}

			if err := tel.Shutdown(ctx); err != nil {
				log.Error("Failed to shutdown telemetry", "error", err)
			}

			return nil
		},
	})
}
