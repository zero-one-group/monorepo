package commands

import (
	"os"

	"go-modular/database"
	"go-modular/internal/config"
	"go-modular/internal/observer/logger"
	"go-modular/internal/server"

	"github.com/spf13/cobra"
)

func init() {
	var argAutoMigrate bool

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the application HTTP server",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context() // Set context for the command
			cfg := config.Get()

			// Initialize application logger
			logger := logger.SetupLogging(logger.LoggerOpts{
				Level:       cfg.GetSlogLevel(),
				Format:      cfg.Logging.Format,
				NoColor:     cfg.Logging.NoColor,
				Environment: cfg.App.Mode,
			})

			if argAutoMigrate {
				logger.Info("Running database migrations before starting server")
				migrator := database.NewMigrator(cfg.GetDatabaseURL())
				if err := migrator.MigrateUp(ctx); err != nil {
					logger.Error("Failed to apply database migration", "err", err)
					os.Exit(1)
				}
				if err := migrator.Close(); err != nil {
					logger.Error("Failed to close database connection", "err", err)
					os.Exit(1)
				}
			}

			// Initialize HTTP server
			httpAddr := "0.0.0.0:8000"
			srv := server.NewHTTPServer(httpAddr, logger)
			if err := srv.Start(); err != nil {
				logger.Error("HTTP server exited with error", "err", err)
				os.Exit(1)
			}
		},
	}

	// Add additional flags for the serve command
	serveCmd.Flags().BoolVar(&argAutoMigrate, "auto-migrate", false, "Run database migrations before starting the server")

	RootCmd.AddCommand(serveCmd)
}
