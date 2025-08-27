package frameworks

import (
	"path/filepath"

	tilocontext "tilokit/internal/core/context"
	"tilokit/internal/plugins/templates"
	pythonTemplates "tilokit/internal/templates/python"
	"tilokit/internal/utils"
	"tilokit/pkg/constants"

	"github.com/pkg/errors"
)

// FastAPIPlugin implements FastAPI framework support
type FastAPIPlugin struct{}

func NewFastAPIPlugin() *FastAPIPlugin {
	return &FastAPIPlugin{}
}

func (p *FastAPIPlugin) Name() string {
	return "fastapi"
}

func (p *FastAPIPlugin) Version() string {
	return constants.VERSION
}

func (p *FastAPIPlugin) Description() string {
	return "FastAPI modern Python web framework"
}

func (p *FastAPIPlugin) SupportedFrameworks() []string {
	return []string{"fastapi"}
}

func (p *FastAPIPlugin) SupportedBuildTools() []string {
	return []string{constants.BuildToolPip, constants.BuildToolPoetry, constants.BuildToolPipenv, constants.BuildToolConda}
}

func (p *FastAPIPlugin) PreGenerate(ctx *tilocontext.ExecutionContext) error {
	if _, ok := ctx.Variables["python_version"]; !ok {
		ctx.SetVariable("python_version", "3.11")
	}

	if _, ok := ctx.Variables["fastapi_version"]; !ok {
		ctx.SetVariable("fastapi_version", "0.104.1")
	}

	if _, ok := ctx.Variables["BuildTool"]; !ok {
		ctx.SetVariable("BuildTool", "pip")
	}

	return nil
}

func (p *FastAPIPlugin) Generate(ctx *tilocontext.ExecutionContext) error {
	projectPath := ctx.ProjectPath
	projectName := ctx.Config.ProjectName
	buildTool := "pip"

	if bt, ok := ctx.Variables["BuildTool"].(string); ok {
		buildTool = bt
	}

	// Create directory structure
	dirs := []string{
		"app", "app/api", "app/api/v1", "app/core", "app/models",
		"app/schemas", "app/crud", "app/db", "tests", "alembic",
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(projectPath, dir)
		if err := utils.EnsureDir(dirPath); err != nil {
			return errors.Wrapf(err, "failed to create directory %s", dir)
		}
	}

	// Generate all files
	if err := p.generateFastAPIApp(ctx, projectPath, projectName); err != nil {
		return errors.Wrap(err, "failed to generate FastAPI app")
	}

	if err := p.generateCore(ctx, projectPath); err != nil {
		return errors.Wrap(err, "failed to generate core")
	}

	if err := p.generateAPI(ctx, projectPath); err != nil {
		return errors.Wrap(err, "failed to generate API")
	}

	if err := p.generateRequirements(ctx, projectPath, buildTool); err != nil {
		return errors.Wrap(err, "failed to generate requirements")
	}

	if err := p.generateDatabase(ctx, projectPath); err != nil {
		return errors.Wrap(err, "failed to generate database")
	}

	if err := p.generateTests(ctx, projectPath); err != nil {
		return errors.Wrap(err, "failed to generate tests")
	}

	if err := p.generateMiscFiles(ctx, projectPath, projectName); err != nil {
		return errors.Wrap(err, "failed to generate misc files")
	}

	return nil
}

func (p *FastAPIPlugin) PostGenerate(ctx *tilocontext.ExecutionContext) error {
	return nil
}

func (p *FastAPIPlugin) generateFastAPIApp(ctx *tilocontext.ExecutionContext, projectPath, projectName string) error {
	engine := templates.NewTemplateEngine()

	processedMain, err := engine.ProcessTemplateWithDelims(pythonTemplates.FastAPIMainPy, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return err
	}

	if err := utils.WriteFile(filepath.Join(projectPath, "main.py"), processedMain); err != nil {
		return err
	}

	// App __init__.py
	return utils.WriteFile(filepath.Join(projectPath, "app", "__init__.py"), "")
}

func (p *FastAPIPlugin) generateCore(ctx *tilocontext.ExecutionContext, projectPath string) error {
	engine := templates.NewTemplateEngine()

	processedConfig, err := engine.ProcessTemplateWithDelims(pythonTemplates.FastAPIConfigPy, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return err
	}

	if err := utils.WriteFile(filepath.Join(projectPath, "app", "core", "config.py"), processedConfig); err != nil {
		return err
	}

	// Core __init__.py
	return utils.WriteFile(filepath.Join(projectPath, "app", "core", "__init__.py"), "")
}

func (p *FastAPIPlugin) generateAPI(ctx *tilocontext.ExecutionContext, projectPath string) error {
	engine := templates.NewTemplateEngine()

	// API router
	processedApi, err := engine.ProcessTemplateWithDelims(pythonTemplates.FastAPIApiPy, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return err
	}

	if err := utils.WriteFile(filepath.Join(projectPath, "app", "api", "v1", "api.py"), processedApi); err != nil {
		return err
	}

	// Create endpoints directory
	endpointsDir := filepath.Join(projectPath, "app", "api", "v1", "endpoints")
	if err := utils.EnsureDir(endpointsDir); err != nil {
		return err
	}

	// Items endpoint
	processedItems, err := engine.ProcessTemplateWithDelims(pythonTemplates.FastAPIItemsPy, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return err
	}

	if err := utils.WriteFile(filepath.Join(endpointsDir, "items.py"), processedItems); err != nil {
		return err
	}

	// Users endpoint
	processedUsers, err := engine.ProcessTemplateWithDelims(pythonTemplates.FastAPIUsersPy, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return err
	}

	if err := utils.WriteFile(filepath.Join(endpointsDir, "users.py"), processedUsers); err != nil {
		return err
	}

	// API dependencies
	processedDeps, err := engine.ProcessTemplateWithDelims(pythonTemplates.FastAPIDepsPy, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return err
	}

	if err := utils.WriteFile(filepath.Join(projectPath, "app", "api", "deps.py"), processedDeps); err != nil {
		return err
	}

	// __init__.py files
	if err := utils.WriteFile(filepath.Join(projectPath, "app", "api", "__init__.py"), ""); err != nil {
		return err
	}
	if err := utils.WriteFile(filepath.Join(projectPath, "app", "api", "v1", "__init__.py"), ""); err != nil {
		return err
	}
	return utils.WriteFile(filepath.Join(endpointsDir, "__init__.py"), "")
}

func (p *FastAPIPlugin) generateRequirements(ctx *tilocontext.ExecutionContext, projectPath, buildTool string) error {
	engine := templates.NewTemplateEngine()

	switch buildTool {
	case constants.BuildToolPoetry:
		processedToml, err := engine.ProcessTemplateWithDelims(pythonTemplates.FastAPIPyprojectToml, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
		if err != nil {
			return err
		}
		return utils.WriteFile(filepath.Join(projectPath, "pyproject.toml"), processedToml)

	case constants.BuildToolPipenv:
		processedPipfile, err := engine.ProcessTemplateWithDelims(pythonTemplates.FastAPIPipfile, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
		if err != nil {
			return err
		}
		return utils.WriteFile(filepath.Join(projectPath, "Pipfile"), processedPipfile)

	case constants.BuildToolConda:
		processedEnv, err := engine.ProcessTemplateWithDelims(pythonTemplates.FastAPIEnvironmentYml, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
		if err != nil {
			return err
		}
		return utils.WriteFile(filepath.Join(projectPath, "environment.yml"), processedEnv)

	default: // pip
		processedReq, err := engine.ProcessTemplateWithDelims(pythonTemplates.FastAPIRequirementsTxt, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
		if err != nil {
			return err
		}
		return utils.WriteFile(filepath.Join(projectPath, "requirements.txt"), processedReq)
	}
}

func (p *FastAPIPlugin) generateDatabase(ctx *tilocontext.ExecutionContext, projectPath string) error {
	engine := templates.NewTemplateEngine()

	// Database session
	processedSession, err := engine.ProcessTemplateWithDelims(pythonTemplates.FastAPISessionPy, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return err
	}

	if err := utils.WriteFile(filepath.Join(projectPath, "app", "db", "session.py"), processedSession); err != nil {
		return err
	}

	// Generate models
	if err := p.generateModels(ctx, projectPath); err != nil {
		return err
	}

	// Generate schemas
	if err := p.generateSchemas(ctx, projectPath); err != nil {
		return err
	}

	return utils.WriteFile(filepath.Join(projectPath, "app", "db", "__init__.py"), "")
}

func (p *FastAPIPlugin) generateModels(ctx *tilocontext.ExecutionContext, projectPath string) error {
	engine := templates.NewTemplateEngine()

	// Item model
	processedItem, err := engine.ProcessTemplateWithDelims(pythonTemplates.FastAPIItemModelPy, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return err
	}

	if err := utils.WriteFile(filepath.Join(projectPath, "app", "models", "item.py"), processedItem); err != nil {
		return err
	}

	// User model
	processedUser, err := engine.ProcessTemplateWithDelims(pythonTemplates.FastAPIUserModelPy, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return err
	}

	if err := utils.WriteFile(filepath.Join(projectPath, "app", "models", "user.py"), processedUser); err != nil {
		return err
	}

	return utils.WriteFile(filepath.Join(projectPath, "app", "models", "__init__.py"), "")
}

func (p *FastAPIPlugin) generateSchemas(ctx *tilocontext.ExecutionContext, projectPath string) error {
	engine := templates.NewTemplateEngine()

	// Item schema
	processedItem, err := engine.ProcessTemplateWithDelims(pythonTemplates.FastAPIItemSchemaPy, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return err
	}

	if err := utils.WriteFile(filepath.Join(projectPath, "app", "schemas", "item.py"), processedItem); err != nil {
		return err
	}

	// User schema
	processedUser, err := engine.ProcessTemplateWithDelims(pythonTemplates.FastAPIUserSchemaPy, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return err
	}

	if err := utils.WriteFile(filepath.Join(projectPath, "app", "schemas", "user.py"), processedUser); err != nil {
		return err
	}

	return utils.WriteFile(filepath.Join(projectPath, "app", "schemas", "__init__.py"), "")
}

func (p *FastAPIPlugin) generateTests(ctx *tilocontext.ExecutionContext, projectPath string) error {
	engine := templates.NewTemplateEngine()

	// Test configuration
	processedConftest, err := engine.ProcessTemplateWithDelims(pythonTemplates.FastAPIConftestPy, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return err
	}

	if err := utils.WriteFile(filepath.Join(projectPath, "tests", "conftest.py"), processedConftest); err != nil {
		return err
	}

	// Basic tests
	processedTestMain, err := engine.ProcessTemplateWithDelims(pythonTemplates.FastAPITestMainPy, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return err
	}

	if err := utils.WriteFile(filepath.Join(projectPath, "tests", "test_main.py"), processedTestMain); err != nil {
		return err
	}

	return utils.WriteFile(filepath.Join(projectPath, "tests", "__init__.py"), "")
}

func (p *FastAPIPlugin) generateMiscFiles(ctx *tilocontext.ExecutionContext, projectPath, projectName string) error {
	engine := templates.NewTemplateEngine()

	processedReadme, err := engine.ProcessTemplateWithDelims(pythonTemplates.FastAPIReadmeMd, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return err
	}

	return utils.WriteFile(filepath.Join(projectPath, "README.md"), processedReadme)
}
