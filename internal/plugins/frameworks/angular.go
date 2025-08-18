package frameworks

import (
	"path/filepath"

	tilocontext "tilokit/internal/core/context"
	"tilokit/internal/plugins/templates"
	"tilokit/internal/templates/angular"
	"tilokit/internal/templates/common"
	"tilokit/internal/utils"
	"tilokit/pkg/constants"

	"github.com/pkg/errors"
)

// AngularPlugin implements Angular framework support
type AngularPlugin struct{}

// NewAngularPlugin creates a new Angular plugin instance
func NewAngularPlugin() *AngularPlugin {
	return &AngularPlugin{}
}

func (p *AngularPlugin) Name() string {
	return "angular-framework"
}

func (p *AngularPlugin) Version() string {
	return constants.VERSION
}

func (p *AngularPlugin) Description() string {
	return "Angular framework with modern setup and best practices"
}

func (p *AngularPlugin) SupportedFrameworks() []string {
	return []string{"angular"}
}

func (p *AngularPlugin) SupportedBuildTools() []string {
	return []string{"angular-cli"}
}

func (p *AngularPlugin) PreGenerate(ctx *tilocontext.ExecutionContext) error {
	// Set Angular-specific variables
	ctx.SetVariable("angular_version", "^17.0.0")
	ctx.SetVariable("typescript_support", true)

	return nil
}

func (p *AngularPlugin) Generate(ctx *tilocontext.ExecutionContext) error {
	// Create directory structure
	if err := p.createDirectoryStructure(ctx); err != nil {
		return errors.Wrap(err, "failed to create directory structure")
	}

	// Generate package.json
	if err := p.generatePackageJson(ctx); err != nil {
		return errors.Wrap(err, "failed to generate package.json")
	}

	// Generate source files
	if err := p.generateSourceFiles(ctx); err != nil {
		return errors.Wrap(err, "failed to generate source files")
	}

	// Generate configuration files
	if err := p.generateConfigFiles(ctx); err != nil {
		return errors.Wrap(err, "failed to generate config files")
	}

	return nil
}

func (p *AngularPlugin) PostGenerate(ctx *tilocontext.ExecutionContext) error {
	// Set post-generation metadata
	ctx.SetMetadata("framework_generated", true)
	ctx.SetMetadata("start_command", "ng serve")

	return nil
}

func (p *AngularPlugin) createDirectoryStructure(ctx *tilocontext.ExecutionContext) error {
	dirs := []string{
		"src",
		"src/app",
		"src/assets",
		"src/environments",
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(ctx.ProjectPath, dir)
		if err := utils.EnsureDir(dirPath); err != nil {
			return err
		}
	}

	return nil
}

func (p *AngularPlugin) generatePackageJson(ctx *tilocontext.ExecutionContext) error {
	// Get rendering mode to determine which package.json to use
	renderingMode, _ := ctx.GetVariable("rendering_mode")

	var packageTemplate string
	if renderingMode != nil && renderingMode.(string) == "ssr" {
		packageTemplate = angular.AngularSSRPackageJson
	} else {
		packageTemplate = angular.AngularPackageJson
	}

	// Use template engine for consistent processing with TILOKit delimiters
	templateEngine := templates.NewTemplateEngine()
	processedContent, err := templateEngine.ProcessTemplateWithDelims(packageTemplate, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return errors.Wrap(err, "failed to process package.json template")
	}

	packageJsonPath := filepath.Join(ctx.ProjectPath, "package.json")
	return utils.WriteFile(packageJsonPath, processedContent)
}

func (p *AngularPlugin) generateSourceFiles(ctx *tilocontext.ExecutionContext) error {
	// Get architecture and rendering mode from context
	architecture, _ := ctx.GetVariable("architecture")
	renderingMode, _ := ctx.GetVariable("rendering_mode")

	// Default values
	archStr := constants.AngularArchitectureStandalone
	if architecture != nil {
		archStr = architecture.(string)
	}

	_ = renderingMode // avoid unused variable error for now

	// Define file templates based on architecture and rendering mode
	var fileTemplates map[string]string

	// For now, only support CSR mode to avoid undefined template errors
	if archStr == constants.AngularArchitectureStandalone {
		fileTemplates = map[string]string{
			"src/app/app.component.ts":    angular.StandaloneAppComponent,
			"src/app/app.component.html":  angular.AngularAppComponentHtml,
			"src/app/app.component.css":   angular.AngularAppComponentCss,
			"src/main.ts":                 angular.StandaloneMainTs,
		}
	} else { // module architecture
		fileTemplates = map[string]string{
			"src/app/app.component.ts":    angular.ModuleAppComponent,
			"src/app/app.component.html":  angular.AngularAppComponentHtml,
			"src/app/app.component.css":   angular.AngularAppComponentCss,
			"src/app/app.module.ts":       angular.ModuleAppModule,
			"src/main.ts":                 angular.ModuleMainTs,
		}
	}

	// Write all files using template engine with custom delimiters
	templateEngine := templates.NewTemplateEngine()

	for relativePath, raw := range fileTemplates {
		// Process template content with TILOKit delimiters
		processedContent, err := templateEngine.ProcessTemplateWithDelims(raw, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
		if err != nil {
			return errors.Wrapf(err, "failed to process template for %s", relativePath)
		}

		fullPath := filepath.Join(ctx.ProjectPath, relativePath)
		if err := utils.WriteFile(fullPath, processedContent); err != nil {
			return err
		}
	}

	return nil
}

func (p *AngularPlugin) generateConfigFiles(ctx *tilocontext.ExecutionContext) error {
	// Get rendering mode and architecture to determine which configs to generate
	renderingMode, _ := ctx.GetVariable("rendering_mode")
	architecture, _ := ctx.GetVariable("architecture")

	_ = renderingMode // avoid unused variable error for now

	archStr := "standalone"
	if architecture != nil {
		archStr = architecture.(string)
	}

	// Define config templates
	configs := map[string]string{
		"angular.json":       angular.AngularJson,
		"tsconfig.json":      angular.AngularTsConfig,
		"tsconfig.app.json":  angular.AngularTsConfigApp,
		"tsconfig.spec.json": angular.AngularTsConfigSpec,
		"src/index.html":     angular.AngularIndexHtml,
		"src/styles.css":     angular.AngularStylesCss,
		".env":               common.Env,
	}

	// For now, only support CSR mode - SSR templates will be added later
	_ = archStr // avoid unused variable error

	// Write all config files using template engine
	templateEngine := templates.NewTemplateEngine()

	for relativePath, content := range configs {
		// Process template content with TILOKit delimiters
		processedContent, err := templateEngine.ProcessTemplateWithDelims(content, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
		if err != nil {
			return errors.Wrapf(err, "failed to process config template for %s", relativePath)
		}

		fullPath := filepath.Join(ctx.ProjectPath, relativePath)
		if err := utils.WriteFile(fullPath, processedContent); err != nil {
			return err
		}
	}

	return nil
}
