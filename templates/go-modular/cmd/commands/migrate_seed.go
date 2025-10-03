//go:build debug
// +build debug

package commands

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"go-modular/database"
	"go-modular/internal/config"
)

var forceSeed bool

var migrateSeedCmd = &cobra.Command{
	Use:   "migrate:seed",
	Short: "Seed the database with initial data",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Get()

		if !forceSeed {
			fmt.Print("Are you sure you want to seed the database? (y/N): ")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(strings.ToLower(input))
			if input != "y" && input != "yes" {
				fmt.Println("Aborted.")
				return
			}
		}

		// Call SeedInitialData to seed initial data
		migrator := database.NewMigrator(cfg.GetDatabaseURL())
		if err := migrator.SeedInitialData(cmd.Context()); err != nil {
			log.Fatalf("Failed to seed initial data: %v", err)
		}

		// Close database connection after seeding
		if err := migrator.Close(); err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}

	},
}

func init() {
	migrateSeedCmd.Flags().BoolVar(&forceSeed, "force", false, "Force seed without confirmation")
	RootCmd.AddCommand(migrateSeedCmd)
}
