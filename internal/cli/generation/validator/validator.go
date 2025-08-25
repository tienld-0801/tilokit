package validator

import (
	"fmt"
	"path/filepath"

	"tilokit/internal/cli/generation/registry"
	"tilokit/internal/utils"
)

// Validator handles input validation for project generation
type Validator struct{}

// NewValidator creates a new Validator instance
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateProjectInputs validates all project generation inputs
func (v *Validator) ValidateProjectInputs(projectName, outputDir string, force bool) error {
	if err := v.ValidateProjectName(projectName); err != nil {
		return err
	}

	if err := v.ValidateProjectDirectory(projectName, outputDir, force); err != nil {
		return err
	}

	return nil
}

// ValidateProjectName validates the project name
func (v *Validator) ValidateProjectName(projectName string) error {
	return utils.ValidateProjectName(projectName)
}

// ValidateProjectDirectory checks if project directory already exists
func (v *Validator) ValidateProjectDirectory(projectName, outputDir string, force bool) error {
	projectPath := projectName
	if outputDir != "." {
		projectPath = filepath.Join(outputDir, projectName)
	}

	if utils.DirExists(projectPath) && !force {
		return fmt.Errorf("directory '%s' already exists. Use --force to overwrite", projectPath)
	}

	return nil
}

// ValidateFramework validates if the framework is supported
func (v *Validator) ValidateFramework(framework string, supportedFrameworks []string) error {
	if framework == "" {
		return fmt.Errorf("framework is required")
	}

	for _, supported := range supportedFrameworks {
		if framework == supported {
			return nil
		}
	}

	return fmt.Errorf("unsupported framework: %s", framework)
}

// ValidateLanguage validates if the language is supported for the given framework
func (v *Validator) ValidateLanguage(language, framework string) error {
	supportedLanguages := map[string][]string{
		"react":  {"ts", "js"},
		"vue":    {"ts", "js"},
		"svelte": {"ts", "js"},
		"next":   {"ts", "js"},
		"nuxt":   {"ts", "js"},
	}

	if languages, exists := supportedLanguages[framework]; exists {
		for _, lang := range languages {
			if language == lang {
				return nil
			}
		}
		return fmt.Errorf("unsupported language '%s' for framework '%s'. Supported: %v", language, framework, languages)
	}

	// For frameworks that don't have language options, any language is acceptable
	return nil
}

// ValidateBuildTool validates if the build tool is supported for the given framework
func (v *Validator) ValidateBuildTool(buildTool, framework string) error {
	if buildTool == "" {
		return nil // Build tool is optional for some frameworks
	}

	pluginRegistry := &registry.PluginRegistry{}
	supportedBuildTools := pluginRegistry.GetBuildToolsForFramework(framework)
	for _, tool := range supportedBuildTools {
		if buildTool == tool {
			return nil
		}
	}

	return fmt.Errorf("unsupported build tool '%s' for framework '%s'. Supported: %v", buildTool, framework, supportedBuildTools)
}
