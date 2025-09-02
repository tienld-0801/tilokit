package cli

import (
	"fmt"

	"tilokit/internal/cli/generation"
	"tilokit/internal/utils"
	"tilokit/pkg/constants"

	"github.com/spf13/cobra"
)

type Manager struct {
	ProjectName    string
	Framework      string
	BuildTool      string
	Language       string
	RouterType     string
	RenderingMode  string
	Architecture   string
	OutputDir      string
	ListFrameworks bool
	ListBuildTools bool
	ShowVersion    bool
	Quiet          bool
	Force          bool
	Update         bool
	InitProject    bool
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) HasAnyFlags(cmd *cobra.Command) bool {
	return m.ProjectName != "" || m.Framework != "" || m.BuildTool != "" ||
		m.Language != "" || m.RouterType != "" || m.RenderingMode != "" || m.Architecture != "" ||
		m.ListFrameworks || m.ListBuildTools || m.Update || m.Quiet ||
		m.Force || m.ShowVersion || m.InitProject
}

func (m *Manager) HandleCommand(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf(constants.InvalidCommandMsg, args[0])
	}

	if m.ShowVersion {
		return ShowVersionInfo()
	}

	if m.Update {
		return m.RunUpdate()
	}

	if m.ListFrameworks {
		return m.ListSupportedFrameworks()
	}

	if m.ListBuildTools {
		return m.ListSupportedBuildTools()
	}

	if m.InitProject {
		return m.RunGenerateWithBanner()
	}

	if m.ProjectName != "" || m.Framework != "" {
		return m.RunGenerate()
	}

	if !m.HasAnyFlags(cmd) {
		return ShowUsageTable()
	}

	return nil
}

func (m *Manager) RunGenerateWithBanner() error {
	utils.PrintBanner()
	utils.SetQuiet(m.Quiet)
	return m.RunProjectGeneration()
}

func (m *Manager) RunGenerate() error {
	utils.SetQuiet(m.Quiet)
	return m.RunProjectGeneration()
}

func (m *Manager) SetupFlags(cmd *cobra.Command) {
	// Project creation flags
	cmd.Flags().StringVarP(&m.ProjectName, "name", "n", "", "Project name (required)")
	cmd.Flags().StringVarP(&m.Framework, "framework", "f", "", "Framework to use (react, vue, svelte, etc.)")
	cmd.Flags().StringVarP(&m.BuildTool, "build-tool", "b", "", "Build tool to use (vite, webpack, etc.)")
	cmd.Flags().StringVarP(&m.Language, "language", "L", "", "Language for templates (ts or js)")
	cmd.Flags().StringVarP(&m.RouterType, "router-type", "r", "", "Router type for Next.js (app or pages)")
	cmd.Flags().StringVarP(&m.RenderingMode, "rendering", "R", "", "Rendering mode for Angular (csr or ssr)")
	cmd.Flags().StringVarP(&m.Architecture, "architecture", "A", "", "Architecture for Angular (standalone or module)")
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

func (m *Manager) RunUpdate() error {
	return RunUpdateProcess()
}

func (m *Manager) ListSupportedFrameworks() error {
	utils.Info("🚀 Supported Frameworks:")
	frameworks := map[string][]string{
		"JavaScript/TypeScript": {"react", "vue", "angular", "svelte", "nextjs", "nuxtjs"},
		"Python":                {"django", "flask", "fastapi"},
		"PHP":                   {"laravel", "symfony", "cakephp", "codeigniter"},
		"Java":                  {"spring-boot", "quarkus"},
		"Go":                    {"gin", "echo", "fiber"},
		"Rust":                  {"actix", "rocket", "axum"},
		"C#":                    {"aspnetcore", "blazor"},
		"Ruby":                  {"rails", "sinatra"},
		"Node.js":               {"express", "nest", "fastify"},
		"Mobile":                {"react-native", "flutter", "ionic"},
		"Desktop":               {"electron", "tauri", "wails"},
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
		"JavaScript":       {"vite", "webpack", "rollup", "parcel"},
		"React Frameworks": {"next"},
		"Vue Frameworks":   {"nuxt"},
		"Package Managers": {"npm", "yarn", "pnpm"},
		"Python":           {"pip", "poetry", "pipenv"},
		"PHP":              {"composer"},
		"Java":             {"maven", "gradle"},
		"Go":               {"go-modules"},
		"Rust":             {"cargo"},
		"C#":               {"dotnet"},
		"Ruby":             {"bundler", "gem"},
		"Mobile/Desktop":   {"metro", "expo", "flutter-cli", "electron-builder"},
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
	orchestrator := generation.NewOrchestrator()
	config := generation.ProjectConfig{
		ProjectName:   m.ProjectName,
		Framework:     m.Framework,
		BuildTool:     m.BuildTool,
		Language:      m.Language,
		RouterType:    m.RouterType,
		RenderingMode: m.RenderingMode,
		Architecture:  m.Architecture,
		OutputDir:     m.OutputDir,
		Force:         m.Force,
	}
	return orchestrator.GenerateProject(config)
}
