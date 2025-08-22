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

// LaravelPlugin implements Laravel framework support
type LaravelPlugin struct{}

// NewLaravelPlugin creates a new Laravel plugin instance
func NewPHPLaravelPlugin() *LaravelPlugin {
	return &LaravelPlugin{}
}

func (p *LaravelPlugin) Name() string {
	return "laravel-framework"
}

func (p *LaravelPlugin) Version() string {
	return constants.VERSION
}

func (p *LaravelPlugin) Description() string {
	return "Laravel PHP framework with modern setup and best practices"
}

func (p *LaravelPlugin) SupportedFrameworks() []string {
	return []string{"laravel"}
}

func (p *LaravelPlugin) SupportedBuildTools() []string {
	return []string{"composer", "vite", "webpack"}
}

func (p *LaravelPlugin) PreGenerate(ctx *tilocontext.ExecutionContext) error {
	// Set Laravel-specific variables
	ctx.SetVariable("laravel_version", "^10.0")
	ctx.SetVariable("php_version", ">=8.1")
	return nil
}

func (p *LaravelPlugin) Generate(ctx *tilocontext.ExecutionContext) error {
	return GenerateLaravel(ctx)
}

func (p *LaravelPlugin) PostGenerate(ctx *tilocontext.ExecutionContext) error {
	return nil
}

func GenerateLaravel(ctx *tilocontext.ExecutionContext) error {
	templateEngine := templates.NewTemplateEngine()

	// Create Laravel project structure
	dirs := []string{
		"app/Http/Controllers",
		"app/Models",
		"bootstrap",
		"config",
		"database/migrations",
		"database/seeders",
		"public",
		"resources/views",
		"routes",
		"storage/app",
		"storage/framework",
		"storage/logs",
		"tests/Feature",
		"tests/Unit",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(ctx.ProjectPath, dir), 0750); err != nil {
			return errors.Wrapf(err, "failed to create directory %s", dir)
		}
	}

	// Generate files using template engine
	files := map[string]string{
		"composer.json":                        php.LaravelComposer,
		"routes/web.php":                       php.LaravelRoutes,
		"app/Http/Controllers/HomeController.php": php.LaravelController,
		".env.example":                         php.LaravelEnvExample,
		"artisan":                              php.LaravelArtisan,
		"public/index.php":                     php.LaravelPublicIndex,
	}

	for filePath, templateContent := range files {
		fullPath := filepath.Join(ctx.ProjectPath, filePath)

		// Process template content with TILOKit delimiters
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
