package frameworks

import (
	"fmt"
	"os"
	"path/filepath"

	tilocontext "tilokit/internal/core/context"
	"tilokit/internal/plugins/templates"
	"tilokit/internal/templates/php"
	"tilokit/internal/utils"
	"tilokit/pkg/constants"

	"github.com/pkg/errors"
)

// SymfonyPlugin implements Symfony framework support
type SymfonyPlugin struct{}

// NewPHPSymfonyPlugin creates a new Symfony plugin instance
func NewPHPSymfonyPlugin() *SymfonyPlugin {
	return &SymfonyPlugin{}
}

func (p *SymfonyPlugin) Name() string {
	return "symfony-framework"
}

func (p *SymfonyPlugin) Version() string {
	return constants.VERSION
}

func (p *SymfonyPlugin) Description() string {
	return "Symfony PHP framework with modern setup and best practices"
}

func (p *SymfonyPlugin) SupportedFrameworks() []string {
	return []string{"symfony"}
}

func (p *SymfonyPlugin) SupportedBuildTools() []string {
	return []string{}
}

func (p *SymfonyPlugin) PreGenerate(ctx *tilocontext.ExecutionContext) error {
	// Set Symfony-specific variables
	ctx.SetVariable("symfony_version", "^6.0")
	ctx.SetVariable("php_version", ">=8.1")

	// Generate random APP_SECRET for security
	secret, err := utils.GenerateRandomSecret(32)
	if err != nil {
		return errors.Wrap(err, "failed to generate APP_SECRET")
	}
	ctx.SetVariable("app_secret", secret)

	return nil
}

func (p *SymfonyPlugin) Generate(ctx *tilocontext.ExecutionContext) error {
	return GenerateSymfony(ctx)
}

func (p *SymfonyPlugin) PostGenerate(ctx *tilocontext.ExecutionContext) error {
	return nil
}

func GenerateSymfony(ctx *tilocontext.ExecutionContext) error {
	templateEngine := templates.NewTemplateEngine()
	// Create Symfony project structure
	dirs := []string{
		"bin",
		"config",
		"public",
		"src/Controller",
		"templates",
		"var/cache",
		"var/log",
		"vendor",
		"migrations",
		"tests",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(ctx.ProjectPath, dir), 0750); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Generate files using template engine
	files := map[string]string{
		"composer.json":                     php.SymfonyComposer,
		"src/Controller/HomeController.php": php.SymfonyController,
		"src/Kernel.php":                    php.SymfonyKernel,
		".env":                              php.SymfonyEnv,
		"public/index.php":                  php.SymfonyPublicIndex,
	}

	for filePath, templateContent := range files {
		fullPath := filepath.Join(ctx.ProjectPath, filePath)

		// Process template with TILOKit delimiters
		processedContent, err := templateEngine.ProcessTemplateWithDelims(templateContent, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
		if err != nil {
			return errors.Wrapf(err, "failed to process template for %s", filePath)
		}

		if err := utils.WriteFile(fullPath, processedContent); err != nil {
			return errors.Wrapf(err, "failed to write file %s", filePath)
		}
	}

	return nil
}
