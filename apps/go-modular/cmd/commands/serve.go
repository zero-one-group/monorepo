package commands

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"go-modular/database"
	"go-modular/internal/server"
)

func init() {
	var argAutoMigrate bool

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the application HTTP server",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context() // Set context for the command

			if argAutoMigrate {
				slog.Info("Running database migrations before starting server")
				databaseURL := os.Getenv("DB_POSTGRES_URL")
				migrator := database.NewMigrator(databaseURL)
				if err := migrator.MigrateUp(ctx); err != nil {
					slog.Error("Failed to apply database migration", "err", err)
					os.Exit(1)
				}
				if err := migrator.Close(); err != nil {
					slog.Error("Failed to close database connection", "err", err)
					os.Exit(1)
				}
			}

			// Initialize HTTP server
			httpAddr := "0.0.0.0:8000"
			srv := server.NewHTTPServer(httpAddr)
			if err := srv.Start(); err != nil {
				slog.Error("HTTP server exited with error", "err", err)
				os.Exit(1)
			}
		},
	}

	// Add additional flags for the serve command
	serveCmd.Flags().BoolVar(&argAutoMigrate, "auto-migrate", false, "Run database migrations before starting the server")

	RootCmd.AddCommand(serveCmd)
}
