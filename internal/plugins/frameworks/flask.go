package frameworks

import (
	"path/filepath"
	"regexp"
	"strings"

	tilocontext "tilokit/internal/core/context"
	"tilokit/internal/plugins/templates"
	pythonTemplates "tilokit/internal/templates/python"
	"tilokit/internal/utils"
	"tilokit/pkg/constants"

	"github.com/pkg/errors"
)

// FlaskPlugin implements Flask framework support
type FlaskPlugin struct{}

func NewFlaskPlugin() *FlaskPlugin {
	return &FlaskPlugin{}
}

func (p *FlaskPlugin) Name() string {
	return "flask"
}

func (p *FlaskPlugin) Version() string {
	return constants.VERSION
}

func (p *FlaskPlugin) Description() string {
	return "Flask micro web framework for Python"
}

func (p *FlaskPlugin) SupportedFrameworks() []string {
	return []string{"flask"}
}

func (p *FlaskPlugin) SupportedBuildTools() []string {
	return []string{
		constants.BuildToolPip,
		constants.BuildToolPoetry,
		constants.BuildToolPipenv,
		constants.BuildToolConda,
	}
}

func (p *FlaskPlugin) PreGenerate(ctx *tilocontext.ExecutionContext) error {
	if _, ok := ctx.Variables["python_version"]; !ok {
		ctx.SetVariable("python_version", "3.11")
	}

	if _, ok := ctx.Variables["flask_version"]; !ok {
		ctx.SetVariable("flask_version", "3.0.0")
	}

	if _, ok := ctx.Variables["BuildTool"]; !ok {
		ctx.SetVariable("BuildTool", "pip")
	}

	// Ensure valid Python package name for templates/imports
	if _, ok := ctx.Variables["project_name"]; !ok {
		name := ctx.Config.ProjectName
		safe := strings.ToLower(name)
		safe = strings.ReplaceAll(safe, "-", "_")
		re := regexp.MustCompile(`[^a-z0-9_]+`)
		safe = re.ReplaceAllString(safe, "_")
		if len(safe) == 0 || safe[0] < 'a' || safe[0] > 'z' {
			safe = "app_" + safe
		}
		ctx.SetVariable("project_name", safe)
	}

	return nil
}

func (p *FlaskPlugin) Generate(ctx *tilocontext.ExecutionContext) error {
	projectPath := ctx.ProjectPath
	projectName := ctx.Config.ProjectName
	buildTool := "pip"

	if bt, ok := ctx.Variables["BuildTool"].(string); ok {
		buildTool = bt
	}

	// Create directory structure
	dirs := []string{
		"app", "app/templates", "app/static", "app/static/css",
		"app/static/js", "tests", "config",
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(projectPath, dir)
		if err := utils.EnsureDir(dirPath); err != nil {
			return errors.Wrapf(err, "failed to create directory %s", dir)
		}
	}

	// Generate all files
	if err := p.generateFlaskApp(ctx, projectPath, projectName); err != nil {
		return errors.Wrap(err, "failed to generate Flask app")
	}

	if err := p.generateConfig(ctx, projectPath); err != nil {
		return errors.Wrap(err, "failed to generate config")
	}

	if err := p.generateRequirements(ctx, projectPath, buildTool); err != nil {
		return errors.Wrap(err, "failed to generate requirements")
	}

	if err := p.generateTemplates(ctx, projectPath, projectName); err != nil {
		return errors.Wrap(err, "failed to generate templates")
	}

	if err := p.generateMiscFiles(ctx, projectPath, projectName); err != nil {
		return errors.Wrap(err, "failed to generate misc files")
	}

	return nil
}

func (p *FlaskPlugin) PostGenerate(ctx *tilocontext.ExecutionContext) error {
	return nil
}

func (p *FlaskPlugin) generateFlaskApp(ctx *tilocontext.ExecutionContext, projectPath, projectName string) error {
	engine := templates.NewTemplateEngine()

	processedApp, err := engine.ProcessTemplateWithDelims(pythonTemplates.FlaskAppPy, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return err
	}

	return utils.WriteFile(filepath.Join(projectPath, "app.py"), processedApp)
}

func (p *FlaskPlugin) generateConfig(ctx *tilocontext.ExecutionContext, projectPath string) error {
	engine := templates.NewTemplateEngine()

	processedConfig, err := engine.ProcessTemplateWithDelims(pythonTemplates.FlaskConfigPy, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return err
	}

	if err := utils.WriteFile(filepath.Join(projectPath, "config", "config.py"), processedConfig); err != nil {
		return err
	}

	return utils.WriteFile(filepath.Join(projectPath, "config", "__init__.py"), "")
}

func (p *FlaskPlugin) generateRequirements(ctx *tilocontext.ExecutionContext, projectPath, buildTool string) error {
	engine := templates.NewTemplateEngine()

	switch buildTool {
	case constants.BuildToolPoetry:
		processedToml, err := engine.ProcessTemplateWithDelims(pythonTemplates.FlaskPyprojectToml, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
		if err != nil {
			return err
		}
		return utils.WriteFile(filepath.Join(projectPath, "pyproject.toml"), processedToml)

	case constants.BuildToolPipenv:
		processedPipfile, err := engine.ProcessTemplateWithDelims(pythonTemplates.FlaskPipfile, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
		if err != nil {
			return err
		}
		return utils.WriteFile(filepath.Join(projectPath, "Pipfile"), processedPipfile)

	case constants.BuildToolConda:
		processedEnv, err := engine.ProcessTemplateWithDelims(pythonTemplates.FlaskEnvironmentYml, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
		if err != nil {
			return err
		}
		return utils.WriteFile(filepath.Join(projectPath, "environment.yml"), processedEnv)

	default: // pip
		processedReq, err := engine.ProcessTemplateWithDelims(pythonTemplates.FlaskRequirementsTxt, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
		if err != nil {
			return err
		}
		return utils.WriteFile(filepath.Join(projectPath, "requirements.txt"), processedReq)
	}
}

func (p *FlaskPlugin) generateTemplates(ctx *tilocontext.ExecutionContext, projectPath, projectName string) error {
	engine := templates.NewTemplateEngine()

	// Base template
	processedBase, err := engine.ProcessTemplateWithDelims(pythonTemplates.FlaskBaseHtml, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return err
	}
	if err := utils.WriteFile(filepath.Join(projectPath, "app", "templates", "base.html"), processedBase); err != nil {
		return err
	}

	// Index template
	processedIndex, err := engine.ProcessTemplateWithDelims(pythonTemplates.FlaskIndexHtml, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return err
	}
	return utils.WriteFile(filepath.Join(projectPath, "app", "templates", "index.html"), processedIndex)
}

func (p *FlaskPlugin) generateMiscFiles(ctx *tilocontext.ExecutionContext, projectPath, projectName string) error {
	engine := templates.NewTemplateEngine()

	processedReadme, err := engine.ProcessTemplateWithDelims(pythonTemplates.FlaskReadmeMd, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return err
	}

	return utils.WriteFile(filepath.Join(projectPath, "README.md"), processedReadme)
}
