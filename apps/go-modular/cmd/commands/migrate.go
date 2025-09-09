//go:build prod || debug
// +build prod debug

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var migrateUpCmd = &cobra.Command{
	Use:   "migrate:up",
	Short: "Apply the latest database migration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Applying latest database migration...")
	},
}

var migrateStatusCmd = &cobra.Command{
	Use:   "migrate:status",
	Short: "Show the status of database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Showing status of database migrations...")
	},
}

var migrateVersionCmd = &cobra.Command{
	Use:   "migrate:version",
	Short: "Show the current database migration version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Showing current database migration version...")
	},
}

func init() {
	RootCmd.AddCommand(migrateUpCmd)
	RootCmd.AddCommand(migrateStatusCmd)
	RootCmd.AddCommand(migrateVersionCmd)
}
