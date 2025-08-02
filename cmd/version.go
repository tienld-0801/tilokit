package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ti-lo/tilokit/internal/utils"
)

var (
	// Version is set during build time
	Version = "v0.1.7-dev"
	// BuildDate is set during build time
	BuildDate = "unknown"
	// GitCommit is set during build time
	GitCommit = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  "Display version, build date, and git commit information for TiLoKit",
	Run: func(cmd *cobra.Command, args []string) {
		utils.PrintBanner()

		fmt.Printf("TiLoKit Version: %s\n", Version)
		fmt.Printf("Build Date: %s\n", BuildDate)
		fmt.Printf("Git Commit: %s\n", GitCommit)
		fmt.Printf("Go Version: %s\n", "1.24.4")

		utils.Info("Visit https://github.com/tienld-0801/tilokit for more information")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
