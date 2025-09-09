//go:build debug
// +build debug

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var forceSeed bool

var migrateSeedCmd = &cobra.Command{
	Use:   "migrate:seed",
	Short: "Seed the database with initial data",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Seeding the database with initial data...")
	},
}

func init() {
	migrateSeedCmd.Flags().BoolVar(&forceSeed, "force", false, "Force seed without confirmation")
	RootCmd.AddCommand(migrateSeedCmd)
}
