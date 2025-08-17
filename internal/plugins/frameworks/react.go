package frameworks

import (
	"path/filepath"

	tilocontext "tilokit/internal/core/context"
	"tilokit/internal/plugins/templates"
	"tilokit/internal/templates/common"
	"tilokit/internal/templates/react"
	"tilokit/internal/utils"
	"tilokit/pkg/constants"

	"github.com/pkg/errors"
)

// ReactPlugin implements React framework support
type ReactPlugin struct{}

// NewReactPlugin creates a new React plugin instance
func NewReactPlugin() *ReactPlugin {
	return &ReactPlugin{}
}

func (p *ReactPlugin) Name() string {
	return "react-framework"
}

func (p *ReactPlugin) Version() string {
	return constants.VERSION
}

func (p *ReactPlugin) Description() string {
	return "React framework with modern setup and best practices"
}

func (p *ReactPlugin) SupportedFrameworks() []string {
	return []string{"react"}
}

func (p *ReactPlugin) SupportedBuildTools() []string {
	return []string{"vite", "webpack", "rollup"}
}

func (p *ReactPlugin) PreGenerate(ctx *tilocontext.ExecutionContext) error {
	// Set React-specific variables
	ctx.SetVariable("react_version", "^18.2.0")
	ctx.SetVariable("react_dom_version", "^18.2.0")
	ctx.SetVariable("typescript_support", true)

	return nil
}

func (p *ReactPlugin) Generate(ctx *tilocontext.ExecutionContext) error {
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

func (p *ReactPlugin) PostGenerate(ctx *tilocontext.ExecutionContext) error {
	// Set post-generation metadata
	ctx.SetMetadata("framework_generated", true)
	ctx.SetMetadata("start_command", "npm run dev")

	return nil
}

func (p *ReactPlugin) createDirectoryStructure(ctx *tilocontext.ExecutionContext) error {
	dirs := []string{
		"src",
		"src/components",
		"src/hooks",
		"src/utils",
		"src/styles",
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

func (p *ReactPlugin) generatePackageJson(ctx *tilocontext.ExecutionContext) error {
	var packageJson, packageJsonFile string
	var envContent, envFile string
	var gitIgnoreContent, gitignoreFile string
	var eslintContent, eslintFile string

	if ctx.Variables["language"].(string) == "js" {
		packageJson = react.ViteJsPackageJson
		eslintContent = react.ViteJsEslint
	} else {
		packageJson = react.ViteTsPackageJson
		eslintContent = react.ViteTsEslint
	}

	envContent = common.Env
	gitIgnoreContent = common.Gitignore

	packageJsonFile = "package.json"
	eslintFile = "eslint.config.js"
	envFile = ".env"
	gitignoreFile = ".gitignore"

	templateEngine := templates.NewTemplateEngine()

	for filename, content := range map[string]string{
		packageJsonFile: packageJson,
		envFile:         envContent,
		gitignoreFile:   gitIgnoreContent,
		eslintFile:      eslintContent,
	} {
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

func (p *ReactPlugin) generateSourceFiles(ctx *tilocontext.ExecutionContext) error {
	var mainContent, appContent string
	var mainFile, appFile string

	if ctx.Variables["language"].(string) == "js" {
		mainContent = react.ViteJsMainFile
		appContent = react.ViteJsAppFile
		mainFile = "main.js"
		appFile = "App.js"
	} else {
		mainContent = react.ViteTsMainFile
		appContent = react.ViteTsAppFile
		mainFile = "main.tsx"
		appFile = "App.tsx"
	}

	files := map[string]string{
		mainFile: mainContent,
		appFile:  appContent,
	}

	templateEngine := templates.NewTemplateEngine()

	for path, content := range files {
		fullPath := filepath.Join(ctx.ProjectPath, "src", path)

		// Process template content with TILOKit delimiters
		processedContent, err := templateEngine.ProcessTemplateWithDelims(content, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
		if err != nil {
			return errors.Wrapf(err, "failed to process template for %s", path)
		}

		if err := utils.WriteFile(fullPath, processedContent); err != nil {
			return err
		}
	}

	return nil
}

func (p *ReactPlugin) generateConfigFiles(ctx *tilocontext.ExecutionContext) error {
	var indexHtml string
	if ctx.Variables["language"].(string) == "js" {
		indexHtml = react.ViteJsIndexHtml
	} else {
		indexHtml = react.ViteTsIndexHtml
	}

	templateEngine := templates.NewTemplateEngine()

	for path, content := range map[string]string{
		"index.html": indexHtml,
	} {
		fullPath := filepath.Join(ctx.ProjectPath, path)

		// Process template content with TILOKit delimiters
		processedContent, err := templateEngine.ProcessTemplateWithDelims(content, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
		if err != nil {
			return errors.Wrapf(err, "failed to process template for %s", path)
		}

		if err := utils.WriteFile(fullPath, processedContent); err != nil {
			return err
		}
	}

	return nil
}
