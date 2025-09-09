package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	var argAutoMigrate bool

	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the application HTTP server",
		Run: func(cmd *cobra.Command, args []string) {
			// ctx := cmd.Context() // Set context for the command

			if argAutoMigrate {
				fmt.Println("Running database migrations...")
			}

			// Log server start
			fmt.Printf("Starting HTTP server...")
		},
	}

	serveCmd.Flags().BoolVar(&argAutoMigrate, "auto-migrate", false, "Run database migrations before starting the server")
	RootCmd.AddCommand(serveCmd)
}
