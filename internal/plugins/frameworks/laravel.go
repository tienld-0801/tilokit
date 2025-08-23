package frameworks

import (
	"os"
	"path/filepath"
	"strings"

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
	return []string{}
}

func (p *LaravelPlugin) PreGenerate(ctx *tilocontext.ExecutionContext) error {
	// Set Laravel-specific variables
	if ctx.Variables == nil {
		ctx.Variables = make(map[string]interface{})
	}
	ctx.Variables["laravel_version"] = "^10.0"
	ctx.Variables["php_version"] = ">=8.1"

	// Generate random APP_KEY for security (Laravel expects base64 format)
	appKey, err := utils.GenerateBase64Secret(32)
	if err != nil {
		return errors.Wrap(err, "failed to generate APP_KEY")
	}
	ctx.Variables["app_key"] = appKey

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
		"bootstrap/cache",
		"config",
		"database/migrations",
		"database/seeders",
		"public",
		"resources/views",
		"routes",
		"storage/app",
		"storage/framework",
		"storage/framework/cache",
		"storage/framework/sessions",
		"storage/framework/views",
		"storage/logs",
		"tests/Feature",
		"tests/Unit",
	}

	for _, dir := range dirs {
		full := filepath.Join(ctx.ProjectPath, dir)
		mode := os.FileMode(0750)
		switch {
		case dir == "public":
			mode = 0755
		case strings.HasPrefix(dir, "bootstrap") || dir == "bootstrap/cache":
			mode = 0770
		case strings.HasPrefix(dir, "storage"):
			mode = 0770
		}
		if err := os.MkdirAll(full, mode); err != nil {
			return errors.Wrapf(err, "failed to create directory %s", dir)
		}
	}

	// Generate files using template engine
	files := map[string]string{
		"composer.json": php.LaravelComposer,
		"app/Http/Controllers/HomeController.php": php.LaravelController,
		"routes/web.php":   php.LaravelRoutes,
		".env.example":     php.LaravelEnvExample,
		"public/index.php": php.LaravelPublicIndex,
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
