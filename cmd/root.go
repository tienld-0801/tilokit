package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/ti-lo/tilokit/internal/config"
	"github.com/ti-lo/tilokit/internal/core/engine"
	"github.com/ti-lo/tilokit/internal/core/registry"
	"github.com/ti-lo/tilokit/internal/plugins/builders"
	"github.com/ti-lo/tilokit/internal/plugins/frameworks"
	"github.com/ti-lo/tilokit/internal/plugins/tools"
	"github.com/ti-lo/tilokit/internal/utils"
)

var (
	projectName    string
	framework      string
	buildTool      string
	outputDir      string
	listFrameworks bool
	listBuildTools bool
	quiet          bool
	force          bool
	update         bool
)

var rootCmd = &cobra.Command{
	Use:   "tilokit",
	Short: "✨ TiLoKit – Modern Multi-Framework Project Generator",
	Long: `TiLoKit is a powerful CLI tool for generating modern web projects
with support for multiple frameworks, build tools, and best practices.

Features:
  • Plugin-based architecture
  • Modern build tool integration (Vite, Webpack, etc.)
  • TypeScript support
  • ESLint & Prettier configuration
  • Testing setup
  • Git initialization
  • Dependency management`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Handle update flag
		if update {
			return updateCmd.RunE(cmd, args)
		}
		return runGenerate()
	},
}

func runGenerate() error {
	// Print banner
	utils.PrintBanner()
	utils.SetQuiet(quiet)

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		utils.Warning("Failed to load config, using defaults: %v", err)
		cfg = &config.Config{}
	}

	// Handle list commands
	if listFrameworks {
		return listSupportedFrameworks()
	}

	if listBuildTools {
		return listSupportedBuildTools()
	}

	// Interactive prompts if values not provided
	if err := promptForMissingValues(cfg); err != nil {
		return err
	}

	// Validate inputs
	if err := validateInputs(); err != nil {
		return err
	}

	// Create project configuration
	projectConfig := config.CreateProjectConfig(projectName, framework, buildTool, outputDir)

	// Initialize engine and register plugins
	eng := engine.New()
	if err := registerPlugins(eng); err != nil {
		return err
	}

	// Execute project generation
	ctx := context.Background()
	if err := eng.Execute(ctx, projectConfig); err != nil {
		utils.Error("Project generation failed: %v", err)
		return err
	}

	// Success message
	utils.Success("%s project '%s' created successfully!", framework, projectName)
	utils.Info("Project location: %s", outputDir)

	// Provide framework-specific next steps
	switch framework {
	case "react", "vue", "angular", "svelte":
		utils.Info("Next steps:")
		utils.Info("   cd %s", projectName)
		utils.Info("   npm install")
		utils.Info("   npm run dev")
	case "django", "flask", "fastapi":
		utils.Info("Next steps:")
		utils.Info("   cd %s", projectName)
		utils.Info("   python -m venv venv")
		utils.Info("   source venv/bin/activate")
		utils.Info("   pip install -r requirements.txt")
	default:
		utils.Info("Check the README.md for setup instructions")
	}

	utils.Info("Happy coding!")

	return nil
}

func promptForMissingValues(cfg *config.Config) error {
	// Project name
	if projectName == "" {
		prompt := &survey.Input{
			Message: "📁 Project name:",
			Help:    "Enter the name for your new project",
		}
		if err := survey.AskOne(prompt, &projectName, survey.WithValidator(survey.Required)); err != nil {
			return err
		}
	}

	// Framework
	if framework == "" {
		supportedFrameworks := []string{"react", "vue", "svelte", "angular", "next", "nuxt", "flask", "fastapi", "quarkus", "spring-boot", "rails"}
		prompt := &survey.Select{
			Message: "🚀 Choose framework:",
			Options: supportedFrameworks,
			Default: cfg.DefaultFramework,
		}
		if err := survey.AskOne(prompt, &framework); err != nil {
			return err
		}
	}

	// Build tool
	if buildTool == "" {
		supportedBuildTools := getBuildToolsForFramework(framework)
		if len(supportedBuildTools) > 1 {
			prompt := &survey.Select{
				Message: "🔧 Choose build tool:",
				Options: supportedBuildTools,
				Default: supportedBuildTools[0],
			}
			if err := survey.AskOne(prompt, &buildTool); err != nil {
				return err
			}
		} else if len(supportedBuildTools) == 1 {
			buildTool = supportedBuildTools[0]
		} else {
			buildTool = "vite" // fallback
		}
	}

	// Output directory
	if outputDir == "" {
		outputDir = "."
	}

	return nil
}

func validateInputs() error {
	if err := utils.ValidateProjectName(projectName); err != nil {
		return err
	}

	// Check if project directory already exists
	projectPath := projectName
	if outputDir != "." {
		projectPath = outputDir + "/" + projectName
	}

	if utils.DirExists(projectPath) && !force {
		utils.Error("Directory '%s' already exists. Use --force to overwrite.", projectPath)
		return os.ErrExist
	}

	return nil
}

func registerPlugins(eng *engine.Engine) error {
	// Register plugins with error handling
	plugins := []registry.Plugin{
		// JavaScript/TypeScript Frameworks
		frameworks.NewReactPlugin(),
		frameworks.NewVuePlugin(),
		// TODO: Add more JS frameworks (Angular, Svelte, Next.js, Nuxt.js)

		// Backend Frameworks
		// Python
		frameworks.NewPythonDjangoPlugin(),
		frameworks.NewPythonFlaskPlugin(),
		frameworks.NewPythonFastAPIPlugin(),

		// PHP
		frameworks.NewPHPLaravelPlugin(),
		frameworks.NewPHPSymfonyPlugin(),

		// Java
		frameworks.NewJavaSpringBootPlugin(),
		frameworks.NewJavaQuarkusPlugin(),

		// Go
		frameworks.NewGoGinPlugin(),
		frameworks.NewGoEchoPlugin(),
		frameworks.NewGoFiberPlugin(),

		// Rust
		frameworks.NewRustActixPlugin(),
		frameworks.NewRustRocketPlugin(),
		frameworks.NewRustAxumPlugin(),

		// C#
		frameworks.NewCSharpASPNetCorePlugin(),
		frameworks.NewCSharpBlazorPlugin(),

		// Ruby
		frameworks.NewRubyRailsPlugin(),
		frameworks.NewRubySinatraPlugin(),

		// Node.js
		frameworks.NewNodeExpressPlugin(),
		frameworks.NewNodeNestJSPlugin(),
		frameworks.NewNodeFastifyPlugin(),

		// Mobile Frameworks
		frameworks.NewReactNativePlugin(),
		frameworks.NewFlutterPlugin(),
		frameworks.NewIonicPlugin(),

		// Desktop Frameworks
		frameworks.NewElectronPlugin(),
		frameworks.NewTauriPlugin(),
		frameworks.NewWailsPlugin(),

		// Build Tools
		builders.NewVitePlugin(),
		builders.NewWebpackPlugin(),
		builders.NewRollupPlugin(),

		// Tools
		tools.NewGitPlugin(),
	}

	// Register all plugins with error handling
	for _, plugin := range plugins {
		if err := eng.RegisterPlugin(plugin); err != nil {
			return fmt.Errorf("failed to register plugin %s: %w", plugin.Name(), err)
		}
	}

	return nil
}

func listSupportedFrameworks() error {
	frameworkOptions := []string{
		// JavaScript/TypeScript
		"react", "vue", "angular", "svelte", "nextjs", "nuxtjs",
		// Python
		"django", "flask", "fastapi",
		// PHP
		"laravel", "symfony",
		// Java
		"spring-boot", "quarkus",
		// Go
		"gin", "echo", "fiber",
		// Rust
		"actix", "rocket", "axum",
		// C#
		"aspnetcore", "blazor",
		// Ruby
		"rails", "sinatra",
		// Node.js
		"express", "nestjs", "fastify",
		// Mobile
		"react-native", "flutter", "ionic",
		// Desktop
		"electron", "tauri", "wails",
	}
	utils.Info("Supported frameworks:")
	for _, fw := range frameworkOptions {
		utils.Log("  • %s", fw)
	}
	return nil
}

func listSupportedBuildTools() error {
	buildToolOptions := []string{
		// JavaScript Build Tools
		"vite", "webpack", "rollup", "parcel",
		// Package Managers
		"npm", "yarn", "pnpm",
		// Language-specific
		"pip", "poetry", "pipenv", // Python
		"composer",        // PHP
		"maven", "gradle", // Java
		"go-modules",     // Go
		"cargo",          // Rust
		"dotnet",         // C#
		"bundler", "gem", // Ruby
		// Mobile/Desktop
		"metro", "expo", "flutter-cli", "electron-builder",
		"none",
	}
	utils.Info("Supported build tools:")
	for _, bt := range buildToolOptions {
		utils.Log("  • %s", bt)
	}
	return nil
}

func getBuildToolsForFramework(framework string) []string {
	buildToolMap := map[string][]string{
		"react":   {"vite", "webpack", "rollup"},
		"vue":     {"vite", "webpack"},
		"svelte":  {"vite", "rollup"},
		"angular": {"angular-cli"},
		"next":    {"next"},
		"nuxt":    {"nuxt"},
	}

	if tools, exists := buildToolMap[framework]; exists {
		return tools
	}
	return []string{"vite"}
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&projectName, "name", "n", "", "Project name")
	rootCmd.Flags().StringVarP(&framework, "framework", "f", "", "Framework to use (react, vue, svelte, etc.)")
	rootCmd.Flags().StringVarP(&buildTool, "build-tool", "b", "", "Build tool to use (vite, webpack, etc.)")
	rootCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory")
	rootCmd.Flags().BoolVar(&listFrameworks, "list-frameworks", false, "List supported frameworks")
	rootCmd.Flags().BoolVar(&listBuildTools, "list-build-tools", false, "List supported build tools")
	rootCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Quiet mode")
	rootCmd.Flags().BoolVar(&force, "force", false, "Force overwrite existing directory")
	rootCmd.Flags().BoolVar(&update, "update", false, "Update TiLoKit to the latest version")
}
