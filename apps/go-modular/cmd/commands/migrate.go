//go:build prod || debug
// +build prod debug

package commands

import (
	"log"

	"github.com/spf13/cobra"
	"go-modular/database"
	"go-modular/internal/config"
)

var migrateUpCmd = &cobra.Command{
	Use:   "migrate:up",
	Short: "Apply the latest database migration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Get()
		migrator := database.NewMigrator(cfg.GetDatabaseURL())
		if err := migrator.MigrateUp(cmd.Context()); err != nil {
			log.Fatalf("Failed to apply database migration: %v", err)
		}
		if err := migrator.Close(); err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
	},
}

var migrateStatusCmd = &cobra.Command{
	Use:   "migrate:status",
	Short: "Show the status of database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Get()
		migrator := database.NewMigrator(cfg.GetDatabaseURL())
		if err := migrator.MigrateStatus(cmd.Context()); err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}
		if err := migrator.Close(); err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
	},
}

var migrateVersionCmd = &cobra.Command{
	Use:   "migrate:version",
	Short: "Show the current database migration version",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Get()
		migrator := database.NewMigrator(cfg.GetDatabaseURL())
		if err := migrator.MigrateVersion(cmd.Context()); err != nil {
			log.Fatalf("Failed to get migration version: %v", err)
		}
		if err := migrator.Close(); err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(migrateUpCmd)
	RootCmd.AddCommand(migrateStatusCmd)
	RootCmd.AddCommand(migrateVersionCmd)
}
