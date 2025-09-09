//go:build debug
// +build debug

package commands

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"go-modular/database"
)

var migrateDownCmd = &cobra.Command{
	Use:   "migrate:down [steps]",
	Short: "Rollback the last or N database migrations",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		databaseURL := os.Getenv("DB_POSTGRES_URL")
		migrator := database.NewMigrator(databaseURL)
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
