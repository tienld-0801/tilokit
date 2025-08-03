package cli

import (
	"fmt"

	"github.com/ti-lo/tilokit/internal/utils"
)

// ShowVersionInfo displays version information without banner
func ShowVersionInfo() error {
	fmt.Printf("TiLoKit Version: %s\n", Version)
	fmt.Printf("Build Date: %s\n", BuildDate)
	fmt.Printf("Git Commit: %s\n", GitCommit)
	fmt.Printf("Go Version: %s\n", GoVersion)
	return nil
}

// ShowUsageTable displays a colorful command usage table
func ShowUsageTable() error {
	fmt.Printf(`
%s
`, utils.ColorizeString(AppShort, "cyan"))
	fmt.Printf("%s\n\n", utils.ColorizeString(AppDescription, "gray"))
	
	fmt.Printf("%s\n", utils.ColorizeString("USAGE", "yellow"))
	fmt.Printf("  %s [flags]\n\n", AppName)
	
	fmt.Printf("%s\n", utils.ColorizeString("⚠️  IMPORTANT: All commands must use - or -- prefix", "red"))
	fmt.Printf("  %s\n\n", "Bare commands like 'version' or 'new' are not allowed")
	
	fmt.Printf("%s\n", utils.ColorizeString("PROJECT CREATION OPTIONS", "yellow"))
	fmt.Printf("  %-20s %s\n", "-n, --name", "Project name (required)")
	fmt.Printf("  %-20s %s\n", "-f, --framework", "Framework to use")
	fmt.Printf("  %-20s %s\n", "-b, --build-tool", "Build tool to use")
	fmt.Printf("  %-20s %s\n\n", "-o, --output", "Output directory")
	
	fmt.Printf("%s\n", utils.ColorizeString("INFORMATION OPTIONS", "yellow"))
	fmt.Printf("  %-20s %s\n", "-l, --list-frameworks", "List supported frameworks")
	fmt.Printf("  %-20s %s\n", "-t, --list-build-tools", "List supported build tools")
	fmt.Printf("  %-20s %s\n\n", "-v, --version", "Show version")
	
	fmt.Printf("%s\n", utils.ColorizeString("PROJECT INITIALIZATION", "yellow"))
	fmt.Printf("  %-20s %s\n\n", "-i, --init", "Initialize new project (with banner)")
	
	fmt.Printf("%s\n", utils.ColorizeString("OTHER OPTIONS", "yellow"))
	fmt.Printf("  %-20s %s\n", "-q, --quiet", "Quiet mode")
	fmt.Printf("  %-20s %s\n", "-F, --force", "Force overwrite")
	fmt.Printf("  %-20s %s\n", "-u, --update", "Update to latest version")
	fmt.Printf("  %-20s %s\n\n", "-h, --help", "Show this help")
	
	fmt.Printf("%s\n", utils.ColorizeString("EXAMPLES", "yellow"))
	fmt.Printf("  %s\n", utils.ColorizeString("tilokit -i", "green"))
	fmt.Printf("  %s\n", utils.ColorizeString("tilokit -n my-app -f react -b vite", "green"))
	fmt.Printf("  %s\n", utils.ColorizeString("tilokit --list-frameworks", "green"))
	fmt.Printf("  %s\n", utils.ColorizeString("tilokit --version", "green"))
	fmt.Printf("  %s\n\n", utils.ColorizeString("tilokit --update", "green"))
	
	return nil
}
