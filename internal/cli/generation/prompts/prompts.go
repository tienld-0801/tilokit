package prompts

import (
	"tilokit/internal/cli/generation/registry"
	"tilokit/internal/config"
	"tilokit/pkg/constants"

	"github.com/AlecAivazis/survey/v2"
)

// PromptHandler handles interactive user prompts for project configuration
type PromptHandler struct{}

// NewPromptHandler creates a new PromptHandler instance
func NewPromptHandler() *PromptHandler {
	return &PromptHandler{}
}

// PromptForProjectName prompts user for project name if not provided
func (p *PromptHandler) PromptForProjectName(currentName string) (string, error) {
	if currentName != "" {
		return currentName, nil
	}

	var projectName string
	prompt := &survey.Input{
		Message: "📁 Project name:",
		Help:    "Enter the name for your new project",
	}
	if err := survey.AskOne(prompt, &projectName, survey.WithValidator(survey.Required)); err != nil {
		return "", err
	}
	return projectName, nil
}

// PromptForFramework prompts user for framework selection if not provided
func (p *PromptHandler) PromptForFramework(currentFramework string, cfg *config.Config) (string, error) {
	if currentFramework != "" {
		return currentFramework, nil
	}

	var framework string
	supportedFrameworks := constants.SupportedFrameworks
	prompt := &survey.Select{
		Message: "🚀 Choose framework:",
		Options: supportedFrameworks,
		Default: cfg.DefaultFramework,
	}
	if err := survey.AskOne(prompt, &framework); err != nil {
		return "", err
	}
	return framework, nil
}

// PromptForBuildTool prompts user for build tool selection if not provided
func (p *PromptHandler) PromptForBuildTool(currentBuildTool, framework string) (string, error) {
	if currentBuildTool != "" {
		return currentBuildTool, nil
	}

	// Backend frameworks don't need build tools
	if framework == constants.ExpressFramework || framework == constants.NestFramework || framework == constants.FastifyFramework {
		return "", nil
	}

	pluginRegistry := &registry.PluginRegistry{}
	supportedBuildTools := pluginRegistry.GetBuildToolsForFramework(framework)
	if len(supportedBuildTools) > 1 {
		var buildTool string
		prompt := &survey.Select{
			Message: "🔧 Choose build tool:",
			Options: supportedBuildTools,
			Default: supportedBuildTools[0],
		}
		if err := survey.AskOne(prompt, &buildTool); err != nil {
			return "", err
		}
		return buildTool, nil
	} else if len(supportedBuildTools) == 1 {
		return supportedBuildTools[0], nil
	}

	// Use framework-appropriate default
	return pluginRegistry.GetDefaultBuildTool(framework), nil
}

// PromptForLanguage prompts user for language selection (TS/JS) for supported frameworks
func (p *PromptHandler) PromptForLanguage(currentLanguage, framework string) (string, error) {
	if currentLanguage != "" {
		return currentLanguage, nil
	}

	switch framework {
	case "react", "vue", "svelte":
		var language string
		langs := []string{"ts", "js"}
		prompt := &survey.Select{
			Message: "🗂 Choose language:",
			Options: langs,
			Default: "ts",
		}
		if err := survey.AskOne(prompt, &language); err != nil {
			return "", err
		}
		return language, nil
	}

	return "ts", nil // Default for other frameworks
}

// PromptForNextJSRouter prompts user for Next.js router type
func (p *PromptHandler) PromptForNextJSRouter(currentRouter, framework string) (string, error) {
	if currentRouter != "" || framework != "next" {
		return currentRouter, nil
	}

	var routerType string
	routers := []string{"app", "pages"}
	prompt := &survey.Select{
		Message: "🛣️ Choose Next.js router:",
		Options: routers,
		Default: "app",
		Help:    "App Router (recommended) uses the new app directory structure. Pages Router uses the traditional pages directory.",
	}
	if err := survey.AskOne(prompt, &routerType); err != nil {
		return "", err
	}
	return routerType, nil
}

// PromptForAngularOptions prompts user for Angular-specific options
func (p *PromptHandler) PromptForAngularOptions(currentRenderingMode, currentArchitecture, framework string) (string, string, error) {
	if framework != "angular" {
		return currentRenderingMode, currentArchitecture, nil
	}

	renderingMode := currentRenderingMode
	architecture := currentArchitecture

	// Rendering mode
	if renderingMode == "" {
		renderingModes := []string{"csr", "ssr"}
		prompt := &survey.Select{
			Message: "🎨 Choose rendering mode:",
			Options: renderingModes,
			Default: "csr",
			Help:    "CSR (Client-Side Rendering) for SPA. SSR (Server-Side Rendering) for better SEO and performance.",
		}
		if err := survey.AskOne(prompt, &renderingMode); err != nil {
			return "", "", err
		}
	}

	// Architecture
	if architecture == "" {
		architectures := []string{"standalone", "module"}
		prompt := &survey.Select{
			Message: "🏗️ Choose architecture:",
			Options: architectures,
			Default: "standalone",
			Help:    "Standalone (Angular 17+) uses modern standalone components. Module uses traditional NgModule architecture.",
		}
		if err := survey.AskOne(prompt, &architecture); err != nil {
			return "", "", err
		}
	}

	return renderingMode, architecture, nil
}

// PromptForReactNativeTemplate prompts user for React Native template selection
func (p *PromptHandler) PromptForReactNativeTemplate(currentTemplate, framework string) (string, error) {
	if currentTemplate != "" || framework != "react-native" {
		return currentTemplate, nil
	}

	var template string
	expoTemplates := []string{"expo"}
	prompt := &survey.Select{
		Message: "📱 Choose React Native template:",
		Options: expoTemplates,
		Default: "expo",
		Help:    "Expo: Managed workflow with Expo SDK. Expo Router: File-based routing. Bare: Minimal React Native. TypeScript: TypeScript template.",
	}
	if err := survey.AskOne(prompt, &template); err != nil {
		return "", err
	}
	return template, nil
}
