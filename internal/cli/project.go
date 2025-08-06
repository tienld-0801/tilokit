package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ti-lo/tilokit/internal/config"
	"github.com/ti-lo/tilokit/internal/core/engine"
	"github.com/ti-lo/tilokit/internal/core/registry"
	"github.com/ti-lo/tilokit/internal/plugins/builders"
	"github.com/ti-lo/tilokit/internal/plugins/frameworks"
	"github.com/ti-lo/tilokit/internal/plugins/tools"
	"github.com/ti-lo/tilokit/internal/utils"
	"github.com/ti-lo/tilokit/pkg/constants"
)

// RunProjectGenerationProcess handles the project generation logic
func (m *Manager) RunProjectGenerationProcess() error {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		// Only continue with defaults for specific expected errors
		if os.IsNotExist(err) {
			utils.Info("No config file found, using defaults")
			cfg = &config.Config{}
		} else {
			return fmt.Errorf("failed to load config: %w", err)
		}
	}

	// Interactive prompts if values not provided
	if err := m.promptForMissingValues(cfg); err != nil {
		return err
	}

	// Validate inputs
	if err := m.validateInputs(); err != nil {
		return err
	}

	// Create project configuration
	projectConfig := config.CreateProjectConfig(m.ProjectName, m.Framework, m.BuildTool, m.OutputDir)

	// Initialize engine and register plugins
	eng := engine.New()
	if err := m.registerPlugins(eng); err != nil {
		return err
	}

	// Execute project generation
	ctx := context.Background()
	if err := eng.Execute(ctx, projectConfig); err != nil {
		utils.Error("Project generation failed: %v", err)
		return err
	}

	// Success message
	utils.Success("%s project '%s' created successfully!", m.Framework, m.ProjectName)
	utils.Info("Project location: %s", m.OutputDir)

	// Provide framework-specific next steps
	switch m.Framework {
	case "react", "vue", "angular", "svelte":
		utils.Info("Next steps:")
		utils.Info("   cd %s", m.ProjectName)
		utils.Info("   npm install")
		utils.Info("   npm run dev")
	case "django", "flask", "fastapi":
		utils.Info("Next steps:")
		utils.Info("   cd %s", m.ProjectName)
		utils.Info("   python -m venv venv")
		utils.Info("   source venv/bin/activate")
		utils.Info("   pip install -r requirements.txt")
	default:
		utils.Info("Check the README.md for setup instructions")
	}

	utils.Info("Happy coding!")
	return nil
}

func (m *Manager) promptForMissingValues(cfg *config.Config) error {
	// Project name
	if m.ProjectName == "" {
		prompt := &survey.Input{
			Message: "📁 Project name:",
			Help:    "Enter the name for your new project",
		}
		if err := survey.AskOne(prompt, &m.ProjectName, survey.WithValidator(survey.Required)); err != nil {
			return err
		}
	}

	// Framework
	if m.Framework == "" {
		supportedFrameworks := constants.SupportedFrameworks
		prompt := &survey.Select{
			Message: "🚀 Choose framework:",
			Options: supportedFrameworks,
			Default: cfg.DefaultFramework,
		}
		if err := survey.AskOne(prompt, &m.Framework); err != nil {
			return err
		}
	}

	// Build tool
	if m.BuildTool == "" {
		supportedBuildTools := m.getBuildToolsForFramework(m.Framework)
		if len(supportedBuildTools) > 1 {
			prompt := &survey.Select{
				Message: "🔧 Choose build tool:",
				Options: supportedBuildTools,
				Default: supportedBuildTools[0],
			}
			if err := survey.AskOne(prompt, &m.BuildTool); err != nil {
				return err
			}
		} else if len(supportedBuildTools) == 1 {
			m.BuildTool = supportedBuildTools[0]
		} else {
			// Use framework-appropriate default
			m.BuildTool = m.getDefaultBuildTool(m.Framework)
		}
	}

	// Output directory
	if m.OutputDir == "" {
		m.OutputDir = "."
	}

	return nil
}

func (m *Manager) validateInputs() error {
	if err := utils.ValidateProjectName(m.ProjectName); err != nil {
		return err
	}

	// Check if project directory already exists
	projectPath := m.ProjectName
	if m.OutputDir != "." {
		projectPath = filepath.Join(m.OutputDir, m.ProjectName)
	}

	if utils.DirExists(projectPath) && !m.Force {
		return fmt.Errorf("directory '%s' already exists. Use --force to overwrite", projectPath)
	}

	return nil
}

func (m *Manager) registerPlugins(eng *engine.Engine) error {
	// Register plugins with error handling
	plugins := []registry.Plugin{
		// JavaScript/TypeScript Frameworks
		frameworks.NewReactPlugin(),
		frameworks.NewVuePlugin(),
		// More JS frameworks can be added here

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

func (m *Manager) getBuildToolsForFramework(framework string) []string {
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

func (m *Manager) getDefaultBuildTool(framework string) string {
	defaults := map[string]string{
		"django":      "pip",
		"flask":       "pip",
		"fastapi":     "pip",
		"spring-boot": "maven",
		"quarkus":     "maven",
		"rails":       "bundler",
		"gin":         "go-modules",
		"echo":        "go-modules",
		"fiber":       "go-modules",
		"laravel":     "composer",
		"symfony":     "composer",
		// JavaScript frameworks default to vite
	}

	if tool, exists := defaults[framework]; exists {
		return tool
	}
	return "vite" // fallback for JS frameworks
}
