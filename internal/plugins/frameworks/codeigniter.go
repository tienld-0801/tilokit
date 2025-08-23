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
	return []string{}
}

func (p *CodeIgniterPlugin) PreGenerate(ctx *tilocontext.ExecutionContext) error {
	// Set CodeIgniter-specific variables
	if ctx.Variables == nil {
		ctx.Variables = make(map[string]interface{})
	}
	ctx.Variables["codeigniter_version"] = "^4.0"
	ctx.Variables["php_version"] = ">=7.4"
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
		full := filepath.Join(ctx.ProjectPath, dir)
		mode := os.FileMode(0750)
		if dir == "public" {
			mode = 0755
		} else if strings.HasPrefix(dir, "writable") {
			mode = 0770
		}
		if err := os.MkdirAll(full, mode); err != nil {
			return errors.Wrapf(err, "failed to create directory %s", full)
		}
	}

	// Generate files using template engine
	files := map[string]string{
		"composer.json":            php.CodeIgniterComposer,
		"app/Controllers/Home.php": php.CodeIgniterController,
		"app/Config/Routes.php":    php.CodeIgniterRoutes,
		"app/Config/Paths.php":     php.CodeIgniterPaths,
		"app/Config/App.php":       php.CodeIgniterApp,
		".env":                     php.CodeIgniterEnv,
		"public/index.php":         php.CodeIgniterPublicIndex,
		"public/.htaccess":         php.CodeIgniterPublicHtaccess,
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

		// Make files world-readable for web server compatibility
		// #nosec G302 - 0644 permissions are required for web server files
		if err := os.Chmod(fullPath, 0644); err != nil {
			return errors.Wrapf(err, "failed to set file permissions for %s", filePath)
		}
	}

	return nil
}
