package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ti-lo/tilokit/internal/utils"
)

// Manager handles all CLI operations
type Manager struct {
	// CLI flags
	ProjectName    string
	Framework      string
	BuildTool      string
	OutputDir      string
	ListFrameworks bool
	ListBuildTools bool
	ShowVersion    bool
	Quiet          bool
	Force          bool
	Update         bool
	InitProject    bool
}

// NewManager creates a new CLI manager
func NewManager() *Manager {
	return &Manager{}
}

// HasAnyFlags checks if any flags are provided
func (m *Manager) HasAnyFlags(cmd *cobra.Command) bool {
	return m.ProjectName != "" || m.Framework != "" || m.BuildTool != "" ||
		m.ListFrameworks || m.ListBuildTools || m.Update || m.Quiet ||
		m.Force || m.ShowVersion || m.InitProject
}

// HandleCommand processes the main command logic
func (m *Manager) HandleCommand(cmd *cobra.Command, args []string) error {
	// Reject any arguments without - or -- prefix
	if len(args) > 0 {
		return fmt.Errorf(InvalidCommandMsg, args[0])
	}

	// Handle version flag first
	if m.ShowVersion {
		return ShowVersionInfo()
	}

	// Handle init flag - this is the only case that shows banner
	if m.InitProject {
		return m.RunGenerateWithBanner()
	}

	// Handle other flags
	if m.Update {
		return m.RunUpdate()
	}

	if m.ListFrameworks {
		return m.ListSupportedFrameworks()
	}

	if m.ListBuildTools {
		return m.ListSupportedBuildTools()
	}

	// If project creation flags provided, run generation without banner
	if m.ProjectName != "" || m.Framework != "" {
		return m.RunGenerate()
	}

	// If no flags provided, show usage
	if !m.HasAnyFlags(cmd) {
		return ShowUsageTable()
	}

	return nil
}

// RunGenerateWithBanner runs project generation with banner (only for -i/--init)
func (m *Manager) RunGenerateWithBanner() error {
	// Print banner and run full project generation flow
	utils.PrintBanner()
	utils.SetQuiet(m.Quiet)
	return m.RunProjectGeneration()
}

// RunGenerate runs project generation without banner
func (m *Manager) RunGenerate() error {
	// No banner for direct flag usage
	utils.SetQuiet(m.Quiet)
	return m.RunProjectGeneration()
}

// SetupFlags configures all CLI flags on the root command
func (m *Manager) SetupFlags(cmd *cobra.Command) {
	// Project creation flags
	cmd.Flags().StringVarP(&m.ProjectName, "name", "n", "", "Project name (required)")
	cmd.Flags().StringVarP(&m.Framework, "framework", "f", "", "Framework to use (react, vue, svelte, etc.)")
	cmd.Flags().StringVarP(&m.BuildTool, "build-tool", "b", "", "Build tool to use (vite, webpack, etc.)")
	cmd.Flags().StringVarP(&m.OutputDir, "output", "o", ".", "Output directory")

	// Information flags
	cmd.Flags().BoolVarP(&m.ListFrameworks, "list-frameworks", "l", false, "List all supported frameworks")
	cmd.Flags().BoolVarP(&m.ListBuildTools, "list-build-tools", "t", false, "List all supported build tools")
	cmd.Flags().BoolVarP(&m.ShowVersion, "version", "v", false, "Show version information")

	// Project initialization
	cmd.Flags().BoolVarP(&m.InitProject, "init", "i", false, "Initialize a new project (with banner)")

	// Other options
	cmd.Flags().BoolVarP(&m.Quiet, "quiet", "q", false, "Quiet mode (suppress output)")
	cmd.Flags().BoolVarP(&m.Force, "force", "F", false, "Force overwrite existing directory")
	cmd.Flags().BoolVarP(&m.Update, "update", "u", false, "Update TiLoKit to the latest version")
}

// Placeholder methods - these will delegate to existing logic
func (m *Manager) RunUpdate() error {
	return RunUpdateProcess()
}

func (m *Manager) ListSupportedFrameworks() error {
	utils.Info("🚀 Supported Frameworks:")
	frameworks := map[string][]string{
		"JavaScript/TypeScript": {"react", "vue", "angular", "svelte", "nextjs", "nuxtjs"},
		"Python": {"django", "flask", "fastapi"},
		"PHP": {"laravel", "symfony"},
		"Java": {"spring-boot", "quarkus"},
		"Go": {"gin", "echo", "fiber"},
		"Rust": {"actix", "rocket", "axum"},
		"C#": {"aspnetcore", "blazor"},
		"Ruby": {"rails", "sinatra"},
		"Node.js": {"express", "nestjs", "fastify"},
		"Mobile": {"react-native", "flutter", "ionic"},
		"Desktop": {"electron", "tauri", "wails"},
	}
	for category, fws := range frameworks {
		utils.Log("  %s:", category)
		for _, fw := range fws {
			utils.Log("    • %s", fw)
		}
	}
	return nil
}

func (m *Manager) ListSupportedBuildTools() error {
	utils.Info("🔧 Supported Build Tools:")
	buildTools := map[string][]string{
		"JavaScript": {"vite", "webpack", "rollup", "parcel"},
		"Package Managers": {"npm", "yarn", "pnpm"},
		"Python": {"pip", "poetry", "pipenv"},
		"PHP": {"composer"},
		"Java": {"maven", "gradle"},
		"Go": {"go-modules"},
		"Rust": {"cargo"},
		"C#": {"dotnet"},
		"Ruby": {"bundler", "gem"},
		"Mobile/Desktop": {"metro", "expo", "flutter-cli", "electron-builder"},
	}
	for category, tools := range buildTools {
		utils.Log("  %s:", category)
		for _, tool := range tools {
			utils.Log("    • %s", tool)
		}
	}
	return nil
}

func (m *Manager) RunProjectGeneration() error {
	return m.RunProjectGenerationProcess()
}
