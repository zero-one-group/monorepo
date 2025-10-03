//go:build prod || debug
// +build prod debug

package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"{{ package_name | kebab_case }}/internal/config"
)

var forceOverwrite bool

func init() {
	RootCmd.AddCommand(&cobra.Command{
		Use:   "generate:config",
		Short: "Generate sample configuration file (env.example)",
		RunE: func(cmd *cobra.Command, args []string) error {
			wd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			exampleFileName := ".env.example"
			envExamplePath := filepath.Join(wd, exampleFileName)

			info, err := os.Stat(envExamplePath)
			if err == nil && !info.IsDir() {
				if !forceOverwrite {
					fmt.Printf("File %s already exists. Overwrite? (y/N): ", exampleFileName)
					reader := bufio.NewReader(os.Stdin)
					input, _ := reader.ReadString('\n')
					input = strings.TrimSpace(strings.ToLower(input))
					if input != "y" && input != "yes" {
						fmt.Println("Aborted.")
						return nil
					}
				}
			} else if err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed to stat %s: %w", envExamplePath, err)
			}

			if err := config.GenerateExampleEnvFile(envExamplePath); err != nil {
				return fmt.Errorf("failed to generate %s: %w", exampleFileName, err)
			}

			fmt.Println("Generated", envExamplePath)
			return nil
		},
	})

	// add --force flag to allow non-interactive overwrite
	RootCmd.Commands()[len(RootCmd.Commands())-1].Flags().
		BoolVar(&forceOverwrite, "force", false, "Overwrite existing files without prompt")
}
