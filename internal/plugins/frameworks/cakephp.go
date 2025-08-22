package frameworks

import (
	"os"
	"path/filepath"

	tilocontext "tilokit/internal/core/context"
	"tilokit/internal/plugins/templates"
	"tilokit/internal/templates/php"
	"tilokit/internal/utils"
	"tilokit/pkg/constants"

	"github.com/pkg/errors"
)

// CakePHPPlugin implements CakePHP framework support
type CakePHPPlugin struct{}

// NewPHPCakePlugin creates a new CakePHP plugin instance
func NewPHPCakePlugin() *CakePHPPlugin {
	return &CakePHPPlugin{}
}

func (p *CakePHPPlugin) Name() string {
	return "cakephp-framework"
}

func (p *CakePHPPlugin) Version() string {
	return constants.VERSION
}

func (p *CakePHPPlugin) Description() string {
	return "CakePHP framework with modern setup and best practices"
}

func (p *CakePHPPlugin) SupportedFrameworks() []string {
	return []string{"cakephp"}
}

func (p *CakePHPPlugin) SupportedBuildTools() []string {
	return []string{"composer", "bake"}
}

func (p *CakePHPPlugin) PreGenerate(ctx *tilocontext.ExecutionContext) error {
	// Set CakePHP-specific variables
	ctx.SetVariable("cakephp_version", "^5.0")
	ctx.SetVariable("php_version", ">=8.1")
	return nil
}

func (p *CakePHPPlugin) Generate(ctx *tilocontext.ExecutionContext) error {
	return GenerateCakePHP(ctx)
}

func (p *CakePHPPlugin) PostGenerate(ctx *tilocontext.ExecutionContext) error {
	return nil
}

func GenerateCakePHP(ctx *tilocontext.ExecutionContext) error {
	templateEngine := templates.NewTemplateEngine()

	// Create CakePHP project structure
	dirs := []string{
		"bin",
		"config",
		"logs",
		"plugins",
		"src/Controller",
		"src/Model/Table",
		"src/Model/Entity",
		"templates/Home",
		"tests/TestCase/Controller",
		"tmp/cache",
		"vendor",
		"webroot",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(ctx.ProjectPath, dir), 0750); err != nil {
			return errors.Wrapf(err, "failed to create directory %s", dir)
		}
	}

	// Generate files using template engine

	files := map[string]string{
		"composer.json":                     php.CakePHPComposer,
		"src/Controller/HomeController.php": php.CakePHPController,
		"src/Controller/AppController.php":  php.CakePHPAppController,
		"config/routes.php":                 php.CakePHPRoutes,
		".env":                              php.CakePHPEnv,
	}

	for filePath, content := range files {
		fullPath := filepath.Join(ctx.ProjectPath, filePath)

		// Process template content with TILOKit delimiters
		processedContent, err := templateEngine.ProcessTemplateWithDelims(content, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
		if err != nil {
			return errors.Wrapf(err, "failed to process template for %s", filePath)
		}

		if err := utils.WriteFile(fullPath, processedContent); err != nil {
			return errors.Wrapf(err, "failed to write file %s", filePath)
		}
	}

	return nil
}
