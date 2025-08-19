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

	mode, _ := ctx.GetVariable("rendering_mode")
	start := "ng serve"
	if v, ok := mode.(string); ok && v == constants.AngularSsrMode {
		// Matches AngularSSRPackageJson scripts
		start = "npm run serve:ssr"
	}
	ctx.SetMetadata("start_command", start)

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
	mode := constants.AngularCsrMode
	if v, ok := renderingMode.(string); ok && v != "" {
		mode = v
	}

	if mode == constants.AngularSsrMode {
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

	// Resolve architecture/mode safely
	archStr := constants.AngularArchitectureStandalone
	if v, ok := architecture.(string); ok && v != "" {
		archStr = v
	}

	mode := constants.AngularCsrMode
	if v, ok := renderingMode.(string); ok && v != "" {
		mode = v
	}

	// Define file templates based on architecture and rendering mode
	var fileTemplates map[string]string
	if mode == constants.AngularSsrMode {
		if archStr == constants.AngularArchitectureStandalone {
			fileTemplates = map[string]string{
				"src/app/app.component.ts":   angular.StandaloneSSRAppComponent,
				"src/app/app.component.html": angular.AngularAppComponentHtml,
				"src/app/app.component.css":  angular.AngularAppComponentCss,
				"src/main.ts":                angular.StandaloneSSRMainTs,
				"src/main.server.ts":         angular.StandaloneSSRMainServerTs,
				"server.ts":                  angular.StandaloneSSRServerTs,
			}
		} else {
			fileTemplates = map[string]string{
				"src/app/app.component.ts":     angular.ModuleSSRAppComponent,
				"src/app/app.component.html":   angular.AngularAppComponentHtml,
				"src/app/app.component.css":    angular.AngularAppComponentCss,
				"src/app/app.module.ts":        angular.ModuleSSRAppModule,
				"src/app/app.server.module.ts": angular.ModuleSSRAppServerModule,
				"src/main.ts":                  angular.ModuleSSRMainTs,
				"src/main.server.ts":           angular.ModuleSSRMainServerTs,
				"server.ts":                    angular.ModuleSSRServerTs,
			}
		}
	} else { // CSR
		if archStr == constants.AngularArchitectureStandalone {
			fileTemplates = map[string]string{
				"src/app/app.component.ts":   angular.StandaloneAppComponent,
				"src/app/app.component.html": angular.AngularAppComponentHtml,
				"src/app/app.component.css":  angular.AngularAppComponentCss,
				"src/main.ts":                angular.StandaloneMainTs,
			}
		} else {
			fileTemplates = map[string]string{
				"src/app/app.component.ts":   angular.ModuleAppComponent,
				"src/app/app.component.html": angular.AngularAppComponentHtml,
				"src/app/app.component.css":  angular.AngularAppComponentCss,
				"src/app/app.module.ts":      angular.ModuleAppModule,
				"src/main.ts":                angular.ModuleMainTs,
			}
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

	// Resolve mode/arch safely
	mode := constants.AngularCsrMode
	if v, ok := renderingMode.(string); ok && v != "" {
		mode = v
	}

	arch := constants.AngularArchitectureStandalone
	if v, ok := architecture.(string); ok && v != "" {
		arch = v
	}

	// Select angular.json
	angularJson := angular.AngularJson
	if mode == constants.AngularSsrMode {
		if arch == constants.AngularArchitectureStandalone {
			angularJson = angular.StandaloneSSRAngularJson
		} else {
			angularJson = angular.ModuleSSRAngularJson
		}
	}

	// Define config templates
	configs := map[string]string{
		"angular.json":       angularJson,
		"tsconfig.json":      angular.AngularTsConfig,
		"tsconfig.app.json":  angular.AngularTsConfigApp,
		"tsconfig.spec.json": angular.AngularTsConfigSpec,
		"src/index.html":     angular.AngularIndexHtml,
		"src/styles.css":     angular.AngularStylesCss,
		".env":               common.Env,
	}

	if mode == constants.AngularSsrMode {
		configs["tsconfig.server.json"] = angular.AngularTsConfigServer
		if arch == constants.AngularArchitectureStandalone {
			configs["src/main.server.ts"] = angular.StandaloneSSRMainServerTs
			configs["server.ts"] = angular.StandaloneSSRServerTs
		} else {
			configs["src/main.server.ts"] = angular.ModuleSSRMainServerTs
			configs["server.ts"] = angular.ModuleSSRServerTs
		}
	}

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
