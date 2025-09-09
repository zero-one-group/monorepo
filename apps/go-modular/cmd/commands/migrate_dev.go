//go:build debug
// +build debug

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var migrateCreateCmd = &cobra.Command{
	Use:   "migrate:create [migration_name]",
	Short: "Create new database migration file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		migrationName := args[0]
		fmt.Printf("Creating migration %s", migrationName)
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "migrate:down",
	Short: "Rollback the latest database migration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Rolling back the latest database migration...")
	},
}

var migrateResetCmd = &cobra.Command{
	Use:   "migrate:reset",
	Short: "Reset the database to the latest migration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Resetting the database to the latest migration...")
	},
}

func init() {
	RootCmd.AddCommand(migrateCreateCmd)
	RootCmd.AddCommand(migrateDownCmd)
	RootCmd.AddCommand(migrateResetCmd)
}
