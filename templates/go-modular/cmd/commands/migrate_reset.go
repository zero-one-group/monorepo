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
	"{{ package_name | kebab_case }}/database"
	"{{ package_name | kebab_case }}/internal/config"
)

var forceReset bool
var argUp bool
var argSeed bool

var migrateResetCmd = &cobra.Command{
	Use:   "migrate:reset",
	Short: "Rollback all database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Get()

		if !forceReset {
			fmt.Print("Are you sure you want to rollback all database migrations? (y/N): ")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(strings.ToLower(input))
			if input != "y" && input != "yes" {
				fmt.Println("Aborted.")
				return
			}
		}

		migratorReset := database.NewMigrator(cfg.GetDatabaseURL())
		err := migratorReset.MigrateReset(cmd.Context())
		if err != nil {
			log.Fatalf("Failed to reset database migration: %v", err)
		}
		if err := migratorReset.Close(); err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}

		// If seed called but not up, return an error
		if argSeed && !argUp {
			log.Println("Cannot run seeders without running migrations up first, please use --up flag.")
			return
		}

		if argUp {
			migratorUp := database.NewMigrator(cfg.GetDatabaseURL())
			if err := migratorUp.MigrateUp(cmd.Context()); err != nil {
				log.Fatalf("Failed to apply database migration: %v", err)
			}
			if err := migratorUp.Close(); err != nil {
				log.Fatalf("Failed to close database connection: %v", err)
			}

			if argSeed {
				seedArgs := make([]string, len(args))
				copy(seedArgs, args)
				if forceReset {
					seedArgs = append(seedArgs, "--force")
				}
				// Set the "force" flag for migrateSeedCmd and check error
				if err := migrateSeedCmd.Flags().Set("force", fmt.Sprintf("%v", forceReset)); err != nil {
					log.Printf("Failed to set force flag for seed command: %v", err)
				}
				migrateSeedCmd.Run(cmd, seedArgs)
			}
		}
	},
}

func init() {
	migrateResetCmd.Flags().BoolVar(&forceReset, "force", false, "Force reset without confirmation")
	migrateResetCmd.Flags().BoolVar(&argUp, "up", false, "Run migrations up after reset")
	migrateResetCmd.Flags().BoolVar(&argSeed, "seed", false, "Run seeders after migration up")
	RootCmd.AddCommand(migrateResetCmd)
}
