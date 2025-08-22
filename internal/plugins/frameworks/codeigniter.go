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

// CodeIgniterPlugin implements CodeIgniter framework support
type CodeIgniterPlugin struct{}

// NewCodeIgniterPlugin creates a new CodeIgniter plugin instance
func NewPHPCodeIgniterPlugin() *CodeIgniterPlugin {
	return &CodeIgniterPlugin{}
}

func (p *CodeIgniterPlugin) Name() string {
	return "codeigniter-framework"
}

func (p *CodeIgniterPlugin) Version() string {
	return constants.VERSION
}

func (p *CodeIgniterPlugin) Description() string {
	return "CodeIgniter PHP framework with modern setup and best practices"
}

func (p *CodeIgniterPlugin) SupportedFrameworks() []string {
	return []string{"codeigniter"}
}

func (p *CodeIgniterPlugin) SupportedBuildTools() []string {
	return []string{"composer"}
}

func (p *CodeIgniterPlugin) PreGenerate(ctx *tilocontext.ExecutionContext) error {
	// Set CodeIgniter-specific variables
	ctx.SetVariable("codeigniter_version", "^4.0")
	ctx.SetVariable("php_version", ">=7.4")
	return nil
}

func (p *CodeIgniterPlugin) Generate(ctx *tilocontext.ExecutionContext) error {
	return GenerateCodeIgniter(ctx)
}

func (p *CodeIgniterPlugin) PostGenerate(ctx *tilocontext.ExecutionContext) error {
	return nil
}

func GenerateCodeIgniter(ctx *tilocontext.ExecutionContext) error {
	templateEngine := templates.NewTemplateEngine()
	// Create CodeIgniter project structure
	dirs := []string{
		"app/Controllers",
		"app/Models",
		"app/Views",
		"app/Config",
		"app/Database/Migrations",
		"app/Database/Seeds",
		"public",
		"writable/cache",
		"writable/logs",
		"writable/session",
		"writable/uploads",
		"tests/unit",
		"tests/feature",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(ctx.ProjectPath, dir), 0750); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Generate files using template engine
	files := map[string]string{
		"composer.json":            php.CodeIgniterComposer,
		"app/Controllers/Home.php": php.CodeIgniterController,
		"app/Config/Routes.php":    php.CodeIgniterRoutes,
		".env":                     php.CodeIgniterEnv,
		"public/index.php":         php.CodeIgniterPublicIndex,
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
