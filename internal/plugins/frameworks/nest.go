package frameworks

import (
	"path/filepath"

	tilocontext "tilokit/internal/core/context"
	"tilokit/internal/plugins/templates"
	"tilokit/internal/templates/common"
	nodejs "tilokit/internal/templates/node"
	"tilokit/internal/utils"
	"tilokit/pkg/constants"

	"github.com/pkg/errors"
)

// NodeNestJSPlugin implements NestJS framework support (TypeScript only)
type NodeNestJSPlugin struct{}

func NewNodeNestJSPlugin() *NodeNestJSPlugin {
	return &NodeNestJSPlugin{}
}

func (p *NodeNestJSPlugin) Name() string {
	return "node-nestjs"
}

func (p *NodeNestJSPlugin) Version() string {
	return constants.VERSION
}

func (p *NodeNestJSPlugin) Description() string {
	return "NestJS progressive Node.js framework (TypeScript)"
}

func (p *NodeNestJSPlugin) SupportedFrameworks() []string {
	return []string{"nestjs", "nest"}
}

func (p *NodeNestJSPlugin) SupportedBuildTools() []string {
	return []string{} // Backend frameworks don't need build tools
}

func (p *NodeNestJSPlugin) PreGenerate(ctx *tilocontext.ExecutionContext) error {
	// NestJS is always TypeScript
	ctx.SetVariable("language", "ts")
	ctx.SetVariable("nestjs_version", "^11.1.4")
	return nil
}

func (p *NodeNestJSPlugin) Generate(ctx *tilocontext.ExecutionContext) error {
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

func (p *NodeNestJSPlugin) PostGenerate(ctx *tilocontext.ExecutionContext) error {
	// Set post-generation metadata
	ctx.SetMetadata("framework_generated", true)
	ctx.SetMetadata("start_command", "npm run start:dev")
	return nil
}

func (p *NodeNestJSPlugin) createDirectoryStructure(ctx *tilocontext.ExecutionContext) error {
	dirs := []string{
		"src",
		"test",
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(ctx.ProjectPath, dir)
		if err := utils.EnsureDir(dirPath); err != nil {
			return err
		}
	}

	return nil
}

func (p *NodeNestJSPlugin) generatePackageJson(ctx *tilocontext.ExecutionContext) error {
	packageJson := nodejs.NestJSPackageJsonTS
	packageJsonFile := constants.PackageJsonFileName
	templateEngine := templates.NewTemplateEngine()

	fullPath := filepath.Join(ctx.ProjectPath, packageJsonFile)

	// Process template content with TILOKit delimiters
	processedContent, err := templateEngine.ProcessTemplateWithDelims(packageJson, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to process template for %s", packageJsonFile)
	}

	if err := utils.WriteFile(fullPath, processedContent); err != nil {
		return err
	}

	return nil
}

func (p *NodeNestJSPlugin) generateSourceFiles(ctx *tilocontext.ExecutionContext) error {
	files := map[string]string{
		"src/main.ts":           nodejs.NestJSMainTS,
		"src/app.module.ts":     nodejs.NestJSAppModuleTS,
		"src/app.controller.ts": nodejs.NestJSAppControllerTS,
		"src/app.service.ts":    nodejs.NestJSAppServiceTS,
	}

	templateEngine := templates.NewTemplateEngine()

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

func (p *NodeNestJSPlugin) generateConfigFiles(ctx *tilocontext.ExecutionContext) error {
	configs := map[string]string{
		constants.TsConfigFileName: nodejs.NestJSTsConfig,
		constants.NestCliFileName: nodejs.NestJSNestCliJson,
		".env": common.Env,
		".gitignore": common.Gitignore,
	}

	templateEngine := templates.NewTemplateEngine()

	for filename, content := range configs {
		fullPath := filepath.Join(ctx.ProjectPath, filename)

		// Process template content with TILOKit delimiters
		processedContent, err := templateEngine.ProcessTemplateWithDelims(content, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
		if err != nil {
			return errors.Wrapf(err, "failed to process config template for %s", filename)
		}

		if err := utils.WriteFile(fullPath, processedContent); err != nil {
			return err
		}
	}

	return nil
}
