//go:build debug
// +build debug

package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var migrateSeedCmd = &cobra.Command{
	Use:   "migrate:seed",
	Short: "Seed the database with initial data",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Seeding the database with initial data...")
	},
}

func init() {
	RootCmd.AddCommand(migrateSeedCmd)
}
