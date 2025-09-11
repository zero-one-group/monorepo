//go:build debug
// +build debug

package commands

import (
	"log"

	"github.com/spf13/cobra"
	"go-modular/database"
	"go-modular/internal/config"
)

var migrateDownCmd = &cobra.Command{
	Use:   "migrate:down [steps]",
	Short: "Rollback the last or N database migrations",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Get()
		migrator := database.NewMigrator(cfg.GetDatabaseURL())
		steps := ""
		if len(args) > 0 {
			steps = args[0]
		}
		if err := migrator.MigrateDown(cmd.Context(), steps); err != nil {
			log.Fatalf("Failed to rollback database migration: %v", err)
		}
		if err := migrator.Close(); err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(migrateDownCmd)
}
