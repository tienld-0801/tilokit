package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"tilokit/internal/cli"
	"tilokit/pkg/constants"
)

var (
	cliManager = cli.NewManager()
)

var rootCmd = &cobra.Command{
	Use:   constants.AppName,
	Short: constants.AppShort,
	Long:  constants.AppDescription,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cliManager.HandleCommand(cmd, args)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cliManager.SetupFlags(rootCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.SetUsageTemplate("")
}
