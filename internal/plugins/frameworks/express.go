package frameworks

import (
	"path/filepath"
	"strings"

	tilocontext "tilokit/internal/core/context"
	"tilokit/internal/plugins/templates"
	"tilokit/internal/templates/common"
	nodejs "tilokit/internal/templates/node"
	"tilokit/internal/utils"
	"tilokit/pkg/constants"

	"github.com/pkg/errors"
)

// NodeExpressPlugin implements Express.js framework support
type NodeExpressPlugin struct{}

func NewNodeExpressPlugin() *NodeExpressPlugin {
	return &NodeExpressPlugin{}
}

func (p *NodeExpressPlugin) Name() string {
	return "node-express"
}

func (p *NodeExpressPlugin) Version() string {
	return constants.VERSION
}

func (p *NodeExpressPlugin) Description() string {
	return "Express.js web framework for Node.js"
}

func (p *NodeExpressPlugin) SupportedFrameworks() []string {
	return []string{"express", "expressjs"}
}

func (p *NodeExpressPlugin) SupportedBuildTools() []string {
	return []string{} // Backend frameworks don't need build tools
}

func (p *NodeExpressPlugin) PreGenerate(ctx *tilocontext.ExecutionContext) error {
	// Default to JavaScript if not provided by the caller
	if _, ok := ctx.Variables["language"]; !ok {
		ctx.SetVariable("language", "js")
	}
	// Don't override an explicitly provided version
	if _, ok := ctx.Variables["express_version"]; !ok {
		ctx.SetVariable("express_version", "^4.21.2")
	}

	// Set required template variables
	ctx.SetVariable("project_name", ctx.Config.ProjectName)
	ctx.SetVariable("package_manager", ctx.Config.PackageManager)

	return nil
}

func (p *NodeExpressPlugin) Generate(ctx *tilocontext.ExecutionContext) error {
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

func (p *NodeExpressPlugin) PostGenerate(ctx *tilocontext.ExecutionContext) error {
	// Set post-generation metadata
	ctx.SetMetadata("framework_generated", true)
	ctx.SetMetadata("start_command", "npm run dev")
	return nil
}

func (p *NodeExpressPlugin) createDirectoryStructure(ctx *tilocontext.ExecutionContext) error {
	dirs := []string{
		"routes",
		"middleware",
		"models",
		"controllers",
		"config",
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(ctx.ProjectPath, dir)
		if err := utils.EnsureDir(dirPath); err != nil {
			return err
		}
	}

	return nil
}

func (p *NodeExpressPlugin) generatePackageJson(ctx *tilocontext.ExecutionContext) error {
	var packageJson string
	lang := "js"
	if v, ok := ctx.Variables["language"].(string); ok && v != "" {
		lang = strings.ToLower(strings.TrimSpace(v))
	}
	if lang == "ts" || lang == "typescript" {
		packageJson = nodejs.ExpressPackageJsonTS
	} else {
		packageJson = nodejs.ExpressPackageJsonJS
	}

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

func (p *NodeExpressPlugin) generateSourceFiles(ctx *tilocontext.ExecutionContext) error {
	lang := "js"
	if v, ok := ctx.Variables["language"].(string); ok && v != "" {
		lang = strings.ToLower(v)
	}

	files := map[string]string{}
	if lang == "ts" || lang == "typescript" {
		// Add src/ for TypeScript projects (dev script expects src/index.ts)
		if err := utils.EnsureDir(filepath.Join(ctx.ProjectPath, "src")); err != nil {
			return err
		}
		files["src/index.ts"] = nodejs.ExpressIndexTS
	} else {
		files["index.js"] = nodejs.ExpressIndexJS
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

func (p *NodeExpressPlugin) generateConfigFiles(ctx *tilocontext.ExecutionContext) error {
	// Base configs
	configs := map[string]string{
		constants.EnvFileName:       common.Env,
		constants.GitignoreFileName: common.Gitignore,
	}

	// Add tsconfig.json when TypeScript is selected
	lang := "js"
	if v, ok := ctx.Variables["language"].(string); ok && v != "" {
		lang = strings.ToLower(strings.TrimSpace(v))
	}
	if lang == "ts" || lang == "typescript" {
		configs[constants.TsConfigFileName] = nodejs.ExpressTsConfig
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
