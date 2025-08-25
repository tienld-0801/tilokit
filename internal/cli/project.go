package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sirupsen/logrus"

	"tilokit/internal/config"
	tilocontext "tilokit/internal/core/context"
	"tilokit/internal/core/engine"
	"tilokit/internal/core/registry"
	"tilokit/internal/plugins/builders"
	"tilokit/internal/plugins/frameworks"
	"tilokit/internal/plugins/tools"
	"tilokit/internal/ui"
	"tilokit/internal/utils"
	"tilokit/pkg/constants"
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
	// Language selection to variables (default ts)
	lang := m.Language
	if lang == "" {
		lang = "ts"
	}
	if projectConfig.Variables == nil {
		projectConfig.Variables = make(map[string]interface{})
	}
	projectConfig.Variables["language"] = lang
	// Router type for Next.js (default app router)
	if m.Framework == "next" {
		routerType := m.RouterType
		if routerType == "" {
			routerType = "app"
		}
		projectConfig.Variables["router_type"] = routerType
	}

	// Angular options
	if m.Framework == "angular" {
		renderingMode := m.RenderingMode
		if renderingMode == "" {
			renderingMode = "csr"
		}
		projectConfig.Variables["rendering_mode"] = renderingMode

		architecture := m.Architecture
		if architecture == "" {
			architecture = "standalone"
		}
		projectConfig.Variables["architecture"] = architecture
	}
	// Initialize engine and register plugins
	eng := engine.New()
	if err := m.registerPlugins(eng); err != nil {
		return err
	}
	// Execute project generation with animated progress
	ctx := context.Background()

	// Define generation steps
	steps := []string{
		"Initializing project structure",
		"Setting up framework configuration",
		"Installing build tools",
		"Creating source files",
		"Configuring development environment",
		"Finalizing project setup",
	}

	// Create progress channel
	progressChan := make(chan ui.ProgressMsg, 10)

	// Start animated progress UI
	doneCh, err := ui.RunProjectProgress(fmt.Sprintf("🚀 Creating %s Project: %s", m.Framework, m.ProjectName), steps, progressChan)
	if err != nil {
		logrus.Errorf("Progress UI error: %v", err)
		return err
	}

	// Send initial progress
	progressChan <- ui.ProgressMsg{Step: "Starting project generation...", Progress: 0.0}

	// Execute project generation with progress updates
	if err := m.executeWithProgress(eng, ctx, projectConfig, progressChan); err != nil {
		progressChan <- ui.ProgressMsg{Step: "❌ Project generation failed", Progress: 1.0, Done: true}
		close(progressChan)
		<-doneCh
		logrus.Errorf("Project generation failed: %v", err)
		return err
	}

	// Send completion
	progressChan <- ui.ProgressMsg{Step: "✅ Project created successfully!", Progress: 1.0, Done: true}
	close(progressChan)

	// Wait for UI to clean up
	<-doneCh
	fmt.Printf("\n📋 Next steps:\n")
	switch m.Framework {
	case constants.ReactFramework, constants.VueFramework, constants.AngularFramework, constants.SvelteFramework:
		fmt.Printf("   cd %s\n", m.ProjectName)
		fmt.Printf("   npm install\n")
		fmt.Printf("   npm run dev\n")
	case constants.NextFramework:
		fmt.Printf("   cd %s\n", m.ProjectName)
		fmt.Printf("   npm install\n")
		fmt.Printf("   npm run dev\n")
		fmt.Printf("   Open http://localhost:3000 to view your Next.js app\n")
	case constants.NuxtFramework:
		fmt.Printf("   cd %s\n", m.ProjectName)
		fmt.Printf("   npm install\n")
		fmt.Printf("   npm run dev\n")
		fmt.Printf("   Open http://localhost:3000 to view your Nuxt.js app\n")
	case constants.ReactNativeFramework:
		fmt.Printf("   cd %s\n", m.ProjectName)
		fmt.Printf("   npm install\n")
		fmt.Printf("   npx expo start\n")
		fmt.Printf("   Scan QR code with Expo Go app or run on simulator\n")
	default:
		fmt.Printf("   cd %s\n", m.ProjectName)
		fmt.Printf("   Follow framework-specific setup instructions\n")
	}
	fmt.Printf("\n🎉 Happy coding!\n")
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
	// Build tool (skip for backend frameworks)
	if m.BuildTool == "" {
		// Backend frameworks don't need build tools
		if m.Framework == constants.ExpressFramework || m.Framework == constants.NestFramework || m.Framework == constants.FastifyFramework {
			m.BuildTool = "" // No build tool needed for backend frameworks
		} else {
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
	}
	// Set default
	// Language (TS or JS) for React/Vue/Svelte
	if m.Language == "" {
		switch m.Framework {
		case "react", "vue", "svelte":
			langs := []string{"ts", "js"}
			defLang := "ts"
			prompt := &survey.Select{
				Message: "🗂 Choose language:",
				Options: langs,
				Default: defLang,
			}
			if err := survey.AskOne(prompt, &m.Language); err != nil {
				return err
			}
		}
	}

	// Router type for Next.js
	if m.RouterType == "" && m.Framework == "next" {
		routers := []string{"app", "pages"}
		prompt := &survey.Select{
			Message: "🛣️ Choose Next.js router:",
			Options: routers,
			Default: "app",
			Help:    "App Router (recommended) uses the new app directory structure. Pages Router uses the traditional pages directory.",
		}
		if err := survey.AskOne(prompt, &m.RouterType); err != nil {
			return err
		}
	}

	// Angular options
	if m.Framework == "angular" {
		// Rendering mode
		if m.RenderingMode == "" {
			renderingModes := []string{"csr", "ssr"}
			prompt := &survey.Select{
				Message: "🎨 Choose rendering mode:",
				Options: renderingModes,
				Default: "csr",
				Help:    "CSR (Client-Side Rendering) for SPA. SSR (Server-Side Rendering) for better SEO and performance.",
			}
			if err := survey.AskOne(prompt, &m.RenderingMode); err != nil {
				return err
			}
		}

		// Architecture
		if m.Architecture == "" {
			architectures := []string{"standalone", "module"}
			prompt := &survey.Select{
				Message: "🏗️ Choose architecture:",
				Options: architectures,
				Default: "standalone",
				Help:    "Standalone (Angular 17+) uses modern standalone components. Module uses traditional NgModule architecture.",
			}
			if err := survey.AskOne(prompt, &m.Architecture); err != nil {
				return err
			}
		}
	}

	// React Native options
	if m.Framework == "react-native" {
		// Expo template selection
		if m.Architecture == "" {
			expoTemplates := []string{"expo"}
			prompt := &survey.Select{
				Message: "📱 Choose React Native template:",
				Options: expoTemplates,
				Default: "expo",
				Help:    "Expo: Managed workflow with Expo SDK. Expo Router: File-based routing. Bare: Minimal React Native. TypeScript: TypeScript template.",
			}
			if err := survey.AskOne(prompt, &m.Architecture); err != nil {
				return err
			}
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
		frameworks.NewSveltePlugin(),
		frameworks.NewAngularPlugin(),
		frameworks.NewNextjsPlugin(),
		frameworks.NewNuxtjsPlugin(),
		// More JS frameworks can be added here
		// Backend Frameworks
		// Python
		frameworks.NewPythonDjangoPlugin(),
		frameworks.NewPythonFlaskPlugin(),
		frameworks.NewPythonFastAPIPlugin(),
		// PHP
		frameworks.NewPHPLaravelPlugin(),
		frameworks.NewPHPSymfonyPlugin(),
		frameworks.NewPHPCakePlugin(),
		frameworks.NewPHPCodeIgniterPlugin(),
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
		"react":        {"vite", "webpack", "rollup"},
		"vue":          {"vite", "webpack"},
		"svelte":       {"vite"},
		"angular":      {"angular-cli"},
		"react-native": {"expo"},
		"next":         {"next"},
		"nuxt":         {"nuxt"},
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
		"next":        "next",
		"nuxt":        "nuxt",
	}
	if tool, exists := defaults[framework]; exists {
		return tool
	}
	return "vite"
}

// executeWithProgress runs project generation with progress updates
func (m *Manager) executeWithProgress(eng *engine.Engine, ctx context.Context, projectConfig *tilocontext.ProjectConfig, progressChan chan<- ui.ProgressMsg) error {
	steps := []struct {
		name     string
		progress float64
	}{
		{"Initializing project structure", 0.16},
		{"Setting up framework configuration", 0.33},
		{"Installing build tools", 0.50},
		{"Creating source files", 0.66},
		{"Configuring development environment", 0.83},
		{"Finalizing project setup", 0.99},
	}

	for i, step := range steps {
		progressChan <- ui.ProgressMsg{Step: step.name, Progress: step.progress}
		time.Sleep(350 * time.Millisecond)

		// Execute actual generation on the last step
		if i == len(steps)-1 {
			if err := eng.Execute(ctx, projectConfig); err != nil {
				return err
			}
		}
	}

	return nil
}
