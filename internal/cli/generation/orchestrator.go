package generation

import (
	"context"
	"fmt"

	"tilokit/internal/cli/generation/progress"
	"tilokit/internal/cli/generation/prompts"
	"tilokit/internal/cli/generation/registry"
	"tilokit/internal/cli/generation/validator"
	"tilokit/internal/config"
	tilocontext "tilokit/internal/core/context"
	"tilokit/internal/core/engine"
	"tilokit/internal/utils"
	"tilokit/pkg/constants"
)

// Orchestrator coordinates the entire project generation process
type Orchestrator struct {
	promptHandler   *prompts.PromptHandler
	validator       *validator.Validator
	pluginRegistry  *registry.PluginRegistry
	progressHandler *progress.ProgressHandler
}

// NewOrchestrator creates a new Orchestrator instance
func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		promptHandler:   prompts.NewPromptHandler(),
		validator:       validator.NewValidator(),
		pluginRegistry:  registry.NewPluginRegistry(),
		progressHandler: progress.NewProgressHandler(),
	}
}

// ProjectConfig holds all project generation parameters
type ProjectConfig struct {
	ProjectName   string
	Framework     string
	BuildTool     string
	Language      string
	RouterType    string
	RenderingMode string
	Architecture  string
	OutputDir     string
	Force         bool
}

// GenerateProject handles the complete project generation workflow
func (o *Orchestrator) GenerateProject(projectConfig ProjectConfig) error {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		utils.Info("No config file found, using defaults")
		cfg = &config.Config{}
	}

	// Interactive prompts for missing values
	if err := o.promptForMissingValues(cfg, &projectConfig); err != nil {
		return err
	}

	// Validate inputs
	if err := o.validator.ValidateProjectInputs(projectConfig.ProjectName, projectConfig.OutputDir, projectConfig.Force); err != nil {
		return err
	}

	// Validate language and build tool compatibility
	if err := o.validator.ValidateLanguage(projectConfig.Language, projectConfig.Framework); err != nil {
		return err
	}

	if err := o.validator.ValidateBuildTool(projectConfig.BuildTool, projectConfig.Framework); err != nil {
		return err
	}

	// Create project configuration
	tilokitProjectConfig := config.CreateProjectConfig(projectConfig.ProjectName, projectConfig.Framework, projectConfig.BuildTool, projectConfig.OutputDir)

	// Set up project variables
	o.setupProjectVariables(tilokitProjectConfig, projectConfig)

	// Initialize engine and register plugins
	eng := engine.New()
	if err := o.pluginRegistry.RegisterAllPlugins(eng); err != nil {
		return err
	}

	// Execute project generation with progress
	ctx := context.Background()
	if err := o.progressHandler.ExecuteWithProgress(eng, ctx, tilokitProjectConfig, projectConfig.ProjectName, projectConfig.Framework); err != nil {
		return err
	}

	// Display next steps
	o.displayNextSteps(projectConfig.ProjectName, projectConfig.Framework)

	return nil
}

// promptForMissingValues handles all interactive prompts
func (o *Orchestrator) promptForMissingValues(cfg *config.Config, projectConfig *ProjectConfig) error {
	var err error

	// Project name
	projectConfig.ProjectName, err = o.promptHandler.PromptForProjectName(projectConfig.ProjectName)
	if err != nil {
		return err
	}

	// Framework
	projectConfig.Framework, err = o.promptHandler.PromptForFramework(projectConfig.Framework, cfg)
	if err != nil {
		return err
	}

	// Build tool
	projectConfig.BuildTool, err = o.promptHandler.PromptForBuildTool(projectConfig.BuildTool, projectConfig.Framework)
	if err != nil {
		return err
	}

	// Language
	projectConfig.Language, err = o.promptHandler.PromptForLanguage(projectConfig.Language, projectConfig.Framework)
	if err != nil {
		return err
	}

	// Next.js router
	projectConfig.RouterType, err = o.promptHandler.PromptForNextJSRouter(projectConfig.RouterType, projectConfig.Framework)
	if err != nil {
		return err
	}

	// Angular options
	projectConfig.RenderingMode, projectConfig.Architecture, err = o.promptHandler.PromptForAngularOptions(projectConfig.RenderingMode, projectConfig.Architecture, projectConfig.Framework)
	if err != nil {
		return err
	}

	// React Native template
	if projectConfig.Framework == "react-native" {
		projectConfig.Architecture, err = o.promptHandler.PromptForReactNativeTemplate(projectConfig.Architecture, projectConfig.Framework)
		if err != nil {
			return err
		}
	}

	// Output directory
	if projectConfig.OutputDir == "" {
		projectConfig.OutputDir = "."
	}

	return nil
}

// setupProjectVariables configures project variables based on framework and options
func (o *Orchestrator) setupProjectVariables(tilokitProjectConfig *tilocontext.ProjectConfig, projectConfig ProjectConfig) {
	if tilokitProjectConfig.Variables == nil {
		tilokitProjectConfig.Variables = make(map[string]interface{})
	}

	// Language selection (default ts)
	language := projectConfig.Language
	if language == "" {
		language = "ts"
	}
	tilokitProjectConfig.Variables["language"] = language

	// Build tool
	tilokitProjectConfig.Variables["BuildTool"] = projectConfig.BuildTool

	// Framework-specific variables
	switch projectConfig.Framework {
	case "next":
		routerType := projectConfig.RouterType
		if routerType == "" {
			routerType = "app"
		}
		tilokitProjectConfig.Variables["router_type"] = routerType

	case "angular":
		renderingMode := projectConfig.RenderingMode
		if renderingMode == "" {
			renderingMode = "csr"
		}
		tilokitProjectConfig.Variables["rendering_mode"] = renderingMode

		architecture := projectConfig.Architecture
		if architecture == "" {
			architecture = "standalone"
		}
		tilokitProjectConfig.Variables["architecture"] = architecture
	}
}

// displayNextSteps shows framework-specific next steps to the user
func (o *Orchestrator) displayNextSteps(projectName, framework string) {
	fmt.Printf("\n📋 Next steps:\n")
	switch framework {
	case constants.ReactFramework, constants.VueFramework, constants.AngularFramework, constants.SvelteFramework:
		fmt.Printf("   cd %s\n", projectName)
		fmt.Printf("   npm install\n")
		fmt.Printf("   npm run dev\n")
	case constants.NextFramework:
		fmt.Printf("   cd %s\n", projectName)
		fmt.Printf("   npm install\n")
		fmt.Printf("   npm run dev\n")
		fmt.Printf("   Open http://localhost:3000 to view your Next.js app\n")
	case constants.NuxtFramework:
		fmt.Printf("   cd %s\n", projectName)
		fmt.Printf("   npm install\n")
		fmt.Printf("   npm run dev\n")
		fmt.Printf("   Open http://localhost:3000 to view your Nuxt.js app\n")
	case constants.ReactNativeFramework:
		fmt.Printf("   cd %s\n", projectName)
		fmt.Printf("   npm install\n")
		fmt.Printf("   npx expo start\n")
		fmt.Printf("   Scan QR code with Expo Go app or run on simulator\n")
	default:
		fmt.Printf("   cd %s\n", projectName)
		fmt.Printf("   Follow framework-specific setup instructions\n")
	}
	fmt.Printf("\n🎉 Happy coding!\n")
}
