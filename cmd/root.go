package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ti-lo/tilokit/internal/cli"
)

var (
	// CLI manager instance
	cliManager = cli.NewManager()
)

var rootCmd = &cobra.Command{
	Use:   cli.AppName,
	Short: cli.AppShort,
	Long:  cli.AppDescription,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cliManager.HandleCommand(cmd, args)
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Set up flags using CLI manager
	cliManager.SetupFlags(rootCmd)

	// Disable default completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Custom help template
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.SetUsageTemplate("")
}
