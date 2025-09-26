package commands

import (
	"fmt"
	"log"

	"{{ package_name | kebab_case }}/internal"
	"{{ package_name | kebab_case }}/internal/config"

	"github.com/spf13/cobra"
)

var cfgFile string
var argVersionShort bool
var argVersionSemantic bool

var RootCmd = &cobra.Command{
	Use:   "{{ package_name | kebab_case }}",
	Short: "Modular monolith Go application",
	Long:  `Modular monolith Go application with a focus on simplicity and maintainability.`,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the application version",
	Run: func(cmd *cobra.Command, args []string) {
		if argVersionShort {
			fmt.Printf("%s (%s)\n", internal.Version, internal.BuildHash)
			return
		} else if argVersionSemantic {
			fmt.Printf("%s\n", internal.Version)
			return
		} else {
			fmt.Printf("Go Minimal %s (%s) %s %s\n", internal.Version, internal.BuildHash, internal.BuildDate, internal.Platform)
		}
	},
}

func init() {
	// Initialize the configuration
	_, err := config.Load(cfgFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Disable the default help subcommand
	RootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	// Add version subcommand
	RootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&argVersionShort, "short", "s", false, "Show short version")
	versionCmd.Flags().BoolVarP(&argVersionSemantic, "semantic", "S", false, "Show semantic version")
}
