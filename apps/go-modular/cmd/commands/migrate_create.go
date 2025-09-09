//go:build debug
// +build debug

package commands

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"go-modular/database"
)

var migrateCreateCmd = &cobra.Command{
	Use:   "migrate:create [migration_name]",
	Short: "Create new database migration file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		migrationName := args[0]
		databaseURL := os.Getenv("DB_POSTGRES_URL")
		migrator := database.NewMigrator(databaseURL)
		if err := migrator.MigrateCreate(cmd.Context(), migrationName); err != nil {
			log.Fatalf("Failed to create new migration: %v", err)
		}
		if err := migrator.Close(); err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(migrateCreateCmd)
}
