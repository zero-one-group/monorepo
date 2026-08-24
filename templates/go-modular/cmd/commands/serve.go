package commands

import (
	"os"

	"{{ package_name | kebab_case }}/database"
	"{{ package_name | kebab_case }}/internal/app"
	"{{ package_name | kebab_case }}/internal/config"
	"{{ package_name | kebab_case }}/internal/observer/logger"

	"github.com/spf13/cobra"
)

func init() {
	var argAutoMigrate bool

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the application HTTP server",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			cfg := config.Get()

			log := logger.SetupLogging(logger.LoggerOpts{
				Level:       cfg.GetSlogLevel(),
				Format:      cfg.Logging.Format,
				NoColor:     cfg.Logging.NoColor,
				Environment: cfg.App.Mode,
			})

			if argAutoMigrate {
				log.Info("Running database migrations before starting server")
				migrator := database.NewMigrator(cfg.GetDatabaseURL())
				if err := migrator.MigrateUp(ctx); err != nil {
					log.Error("Failed to apply database migration", "err", err)
					os.Exit(1)
				}
				if err := migrator.Close(); err != nil {
					log.Error("Failed to close database connection", "err", err)
					os.Exit(1)
				}
			}

			app.New(cfg, log).Run()
		},
	}

	serveCmd.Flags().BoolVar(&argAutoMigrate, "auto-migrate", false, "Run database migrations before starting the server")
	RootCmd.AddCommand(serveCmd)
}
