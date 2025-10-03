//go:build debug
// +build debug

package commands

import (
	"log"

	"github.com/spf13/cobra"
	"{{ package_name | kebab_case }}/database"
	"{{ package_name | kebab_case }}/internal/config"
)

var migrateCreateCmd = &cobra.Command{
	Use:   "migrate:create [migration_name]",
	Short: "Create new database migration file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Get()
		migrationName := args[0]
		migrator := database.NewMigrator(cfg.GetDatabaseURL())
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
