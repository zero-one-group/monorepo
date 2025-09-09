package commands

import (
	"fmt"
	"go-modular/database"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
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

			// Log server start
			fmt.Printf("Starting HTTP server...")
		},
	}

	// Add additional flags for the serve command
	serveCmd.Flags().BoolVar(&argAutoMigrate, "auto-migrate", false, "Run database migrations before starting the server")

	RootCmd.AddCommand(serveCmd)
}
