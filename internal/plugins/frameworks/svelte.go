package frameworks

import (
	"path/filepath"

	tilocontext "tilokit/internal/core/context"
	"tilokit/internal/plugins/templates"
	"tilokit/internal/templates/common"
	"tilokit/internal/templates/svelte"
	"tilokit/internal/utils"
	"tilokit/pkg/constants"

	"github.com/pkg/errors"
)

// SveltePlugin implements Svelte framework support
type SveltePlugin struct{}

// NewSveltePlugin creates a new Svelte plugin instance
func NewSveltePlugin() *SveltePlugin {
	return &SveltePlugin{}
}

func (p *SveltePlugin) Name() string {
	return "svelte-framework"
}

func (p *SveltePlugin) Version() string {
	return constants.VERSION
}

func (p *SveltePlugin) Description() string {
	return "Svelte framework with Vite and modern setup"
}

func (p *SveltePlugin) SupportedFrameworks() []string {
	return []string{"svelte"}
}

func (p *SveltePlugin) SupportedBuildTools() []string {
	return []string{"vite"}
}

func (p *SveltePlugin) PreGenerate(ctx *tilocontext.ExecutionContext) error {
	// Set Svelte-specific variables
	ctx.SetVariable("svelte_version", "^5.38.2")
	ctx.SetVariable("vite_version", "^5.0.3")
	ctx.SetVariable("typescript_support", true)

	return nil
}

func (p *SveltePlugin) Generate(ctx *tilocontext.ExecutionContext) error {
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

func (p *SveltePlugin) PostGenerate(ctx *tilocontext.ExecutionContext) error {
	// Set post-generation metadata
	ctx.SetMetadata("framework_generated", true)
	ctx.SetMetadata("start_command", "npm run dev")

	return nil
}

func (p *SveltePlugin) createDirectoryStructure(ctx *tilocontext.ExecutionContext) error {
	dirs := []string{
		"src",
		"src/lib",
		"src/assets",
		"public",
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(ctx.ProjectPath, dir)
		if err := utils.EnsureDir(dirPath); err != nil {
			return err
		}
	}

	return nil
}

func (p *SveltePlugin) generatePackageJson(ctx *tilocontext.ExecutionContext) error {
	var packageJson string

	language, _ := ctx.Variables["language"].(string)
	if language == "js" {
		packageJson = svelte.ViteJsPackageJson
	} else {
		packageJson = svelte.ViteTsPackageJson
	}

	templateEngine := templates.NewTemplateEngine()

	files := map[string]string{
		"package.json": packageJson,
		".env":         common.Env,
		".gitignore":   common.Gitignore,
	}

	for filename, content := range files {
		fullPath := filepath.Join(ctx.ProjectPath, filename)

		// Process template content with TILOKit delimiters
		processedContent, err := templateEngine.ProcessTemplateWithDelims(content, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
		if err != nil {
			return errors.Wrapf(err, "failed to process template for %s", filename)
		}

		if err := utils.WriteFile(fullPath, processedContent); err != nil {
			return err
		}
	}

	return nil
}

func (p *SveltePlugin) generateSourceFiles(ctx *tilocontext.ExecutionContext) error {
	language, _ := ctx.Variables["language"].(string)
	isJS := language == "js"

	var fileTemplates map[string]string
	if isJS {
		fileTemplates = p.getJavaScriptTemplates()
	} else {
		fileTemplates = p.getTypeScriptTemplates()
	}

	templateEngine := templates.NewTemplateEngine()

	for relativePath, content := range fileTemplates {
		fullPath := filepath.Join(ctx.ProjectPath, relativePath)

		// Process template content with TILOKit delimiters
		processedContent, err := templateEngine.ProcessTemplateWithDelims(content, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
		if err != nil {
			return errors.Wrapf(err, "failed to process template for %s", relativePath)
		}

		if err := utils.WriteFile(fullPath, processedContent); err != nil {
			return err
		}
	}

	return nil
}

func (p *SveltePlugin) getJavaScriptTemplates() map[string]string {
	return map[string]string{
		"src/main.js":            svelte.ViteJsMainFile,
		"src/App.svelte":         svelte.ViteJsAppFile,
		"src/lib/Counter.svelte":  svelte.ViteJsCounterComponent,
		"src/app.css":            svelte.ViteJsAppCss,
		"src/assets/svelte.svg":  svelte.SvelteSvg,
		"public/vite.svg":        svelte.ViteSvg,
	}
}

func (p *SveltePlugin) getTypeScriptTemplates() map[string]string {
	return map[string]string{
		"src/main.ts":            svelte.ViteTsMainFile,
		"src/App.svelte":         svelte.ViteTsAppFile,
		"src/lib/Counter.svelte":  svelte.ViteTsCounterComponent,
		"src/app.css":            svelte.ViteTsAppCss,
		"src/app.d.ts":           svelte.ViteTsAppDTs,
		"src/assets/svelte.svg":  svelte.SvelteSvg,
		"public/vite.svg":        svelte.ViteSvg,
	}
}

func (p *SveltePlugin) generateConfigFiles(ctx *tilocontext.ExecutionContext) error {
	language, _ := ctx.Variables["language"].(string)
	isJS := language == "js"

	var configs map[string]string
	if isJS {
		configs = p.getJavaScriptConfigs()
	} else {
		configs = p.getTypeScriptConfigs()
	}

	templateEngine := templates.NewTemplateEngine()

	for relativePath, content := range configs {
		fullPath := filepath.Join(ctx.ProjectPath, relativePath)

		// Signal that Vite config already exists to avoid duplication by the Vite builder
		if relativePath == "vite.config.js" || relativePath == "vite.config.ts" {
			ctx.SetMetadata("vite_config_generated", true)
		}

		// Process template content with TILOKit delimiters
		processedContent, err := templateEngine.ProcessTemplateWithDelims(content, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
		if err != nil {
			return errors.Wrapf(err, "failed to process config template for %s", relativePath)
		}

		if err := utils.WriteFile(fullPath, processedContent); err != nil {
			return err
		}
	}

	return nil
}

func (p *SveltePlugin) getJavaScriptConfigs() map[string]string {
	return map[string]string{
		"index.html":     svelte.ViteJsIndexHtml,
		"vite.config.js": svelte.ViteJsViteConfig,
	}
}

func (p *SveltePlugin) getTypeScriptConfigs() map[string]string {
	return map[string]string{
		"index.html":         svelte.ViteTsIndexHtml,
		"vite.config.ts":     svelte.ViteTsViteConfig,
		"tsconfig.json":      svelte.ViteTsTsConfig,
		"tsconfig.node.json": svelte.ViteTsTsConfigNode,
	}
}
