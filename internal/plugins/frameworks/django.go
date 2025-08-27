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

// DjangoPlugin implements Django framework support
type DjangoPlugin struct{}

func NewDjangoPlugin() *DjangoPlugin {
	return &DjangoPlugin{}
}

func (p *DjangoPlugin) Name() string {
	return "django"
}

func (p *DjangoPlugin) Version() string {
	return constants.VERSION
}

func (p *DjangoPlugin) Description() string {
	return "Django web framework for Python with modern setup"
}

func (p *DjangoPlugin) SupportedFrameworks() []string {
	return []string{"django"}
}

func (p *DjangoPlugin) SupportedBuildTools() []string {
	return []string{"pip", "poetry", "pipenv"}
}

func (p *DjangoPlugin) PreGenerate(ctx *tilocontext.ExecutionContext) error {
	if _, ok := ctx.Variables["python_version"]; !ok {
		ctx.SetVariable("python_version", "3.11")
	}

	if _, ok := ctx.Variables["django_version"]; !ok {
		ctx.SetVariable("django_version", "4.2.7")
	}

	if _, ok := ctx.Variables["BuildTool"]; !ok {
		ctx.SetVariable("BuildTool", constants.BuildToolPip)
	}

	return nil
}

func (p *DjangoPlugin) Generate(ctx *tilocontext.ExecutionContext) error {
	projectPath := ctx.ProjectPath
	projectName := ctx.Config.ProjectName
	buildTool := constants.BuildToolPip

	if bt, ok := ctx.Variables["BuildTool"].(string); ok {
		buildTool = bt
	}

	// Create Django project structure
	dirs := []string{
		projectName, projectName + "/apps", projectName + "/apps/core",
		projectName + "/settings", "static", "templates", "tests",
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(projectPath, dir)
		if err := utils.EnsureDir(dirPath); err != nil {
			return errors.Wrapf(err, "failed to create directory %s", dir)
		}
	}

	// Generate all files
	if err := p.generateDjangoProject(ctx, projectPath, projectName); err != nil {
		return errors.Wrap(err, "failed to generate Django project")
	}

	if err := p.generateSettings(ctx, projectPath, projectName); err != nil {
		return errors.Wrap(err, "failed to generate settings")
	}

	if err := p.generateApps(ctx, projectPath, projectName); err != nil {
		return errors.Wrap(err, "failed to generate apps")
	}

	if err := p.generateRequirements(ctx, projectPath, buildTool); err != nil {
		return errors.Wrap(err, "failed to generate requirements")
	}

	if err := p.generateMiscFiles(ctx, projectPath, projectName); err != nil {
		return errors.Wrap(err, "failed to generate misc files")
	}

	return nil
}

func (p *DjangoPlugin) PostGenerate(ctx *tilocontext.ExecutionContext) error {
	return nil
}

func (p *DjangoPlugin) generateDjangoProject(ctx *tilocontext.ExecutionContext, projectPath, projectName string) error {
	engine := templates.NewTemplateEngine()

	// Main manage.py
	processedManage, err := engine.ProcessTemplate(pythonTemplates.DjangoManagePy, ctx)
	if err != nil {
		return err
	}

	if err := utils.WriteFile(filepath.Join(projectPath, "manage.py"), processedManage); err != nil {
		return err
	}

	// Main URLs
	processedUrls, err := engine.ProcessTemplate(pythonTemplates.DjangoUrlsPy, ctx)
	if err != nil {
		return err
	}

	return utils.WriteFile(filepath.Join(projectPath, projectName, "urls.py"), processedUrls)
}

func (p *DjangoPlugin) generateSettings(ctx *tilocontext.ExecutionContext, projectPath, projectName string) error {
	engine := templates.NewTemplateEngine()

	processedSettings, err := engine.ProcessTemplateWithDelims(pythonTemplates.DjangoBasePy, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return err
	}

	if err := utils.WriteFile(filepath.Join(projectPath, projectName, "settings", "base.py"), processedSettings); err != nil {
		return err
	}

	processedDev, err := engine.ProcessTemplateWithDelims(pythonTemplates.DjangoDevelopmentPy, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return err
	}

	if err := utils.WriteFile(filepath.Join(projectPath, projectName, "settings", "development.py"), processedDev); err != nil {
		return err
	}

	// Create __init__.py files
	initFiles := []string{
		filepath.Join(projectPath, projectName, "__init__.py"),
		filepath.Join(projectPath, projectName, "settings", "__init__.py"),
	}

	for _, initFile := range initFiles {
		if err := utils.WriteFile(initFile, ""); err != nil {
			return err
		}
	}

	return nil
}

func (p *DjangoPlugin) generateApps(ctx *tilocontext.ExecutionContext, projectPath, projectName string) error {
	engine := templates.NewTemplateEngine()

	// Core app config
	processedApps, err := engine.ProcessTemplate(pythonTemplates.DjangoCoreAppsPy, ctx)
	if err != nil {
		return err
	}

	if err := utils.WriteFile(filepath.Join(projectPath, projectName, "apps", "core", "apps.py"), processedApps); err != nil {
		return err
	}

	// Core models
	processedModels, err := engine.ProcessTemplate(pythonTemplates.DjangoCoreModelsPy, ctx)
	if err != nil {
		return err
	}

	if err := utils.WriteFile(filepath.Join(projectPath, projectName, "apps", "core", "models.py"), processedModels); err != nil {
		return err
	}

	// Core views
	processedViews, err := engine.ProcessTemplate(pythonTemplates.DjangoCoreViewsPy, ctx)
	if err != nil {
		return err
	}

	if err := utils.WriteFile(filepath.Join(projectPath, projectName, "apps", "core", "views.py"), processedViews); err != nil {
		return err
	}

	// Core URLs
	processedUrls, err := engine.ProcessTemplate(pythonTemplates.DjangoCoreUrlsPy, ctx)
	if err != nil {
		return err
	}

	if err := utils.WriteFile(filepath.Join(projectPath, projectName, "apps", "core", "urls.py"), processedUrls); err != nil {
		return err
	}

	// Create __init__.py files
	initFiles := []string{
		filepath.Join(projectPath, projectName, "apps", "__init__.py"),
		filepath.Join(projectPath, projectName, "apps", "core", "__init__.py"),
		filepath.Join(projectPath, projectName, "apps", "core", "migrations", "__init__.py"),
	}

	for _, initFile := range initFiles {
		if err := utils.WriteFile(initFile, ""); err != nil {
			return err
		}
	}

	return nil
}

func (p *DjangoPlugin) generateRequirements(ctx *tilocontext.ExecutionContext, projectPath, buildTool string) error {
	engine := templates.NewTemplateEngine()

	switch buildTool {
	case constants.BuildToolPoetry:
		processedToml, err := engine.ProcessTemplateWithDelims(pythonTemplates.DjangoPyprojectToml, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
		if err != nil {
			return err
		}
		return utils.WriteFile(filepath.Join(projectPath, "pyproject.toml"), processedToml)

	case constants.BuildToolPipenv:
		processedPipfile, err := engine.ProcessTemplateWithDelims(pythonTemplates.DjangoPipfile, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
		if err != nil {
			return err
		}
		return utils.WriteFile(filepath.Join(projectPath, "Pipfile"), processedPipfile)

	case constants.BuildToolConda:
		processedEnv, err := engine.ProcessTemplateWithDelims(pythonTemplates.DjangoEnvironmentYml, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
		if err != nil {
			return err
		}
		return utils.WriteFile(filepath.Join(projectPath, "environment.yml"), processedEnv)

	default: // pip
		processedReq, err := engine.ProcessTemplateWithDelims(pythonTemplates.DjangoRequirementsTxt, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
		if err != nil {
			return err
		}
		return utils.WriteFile(filepath.Join(projectPath, "requirements.txt"), processedReq)
	}
}

func (p *DjangoPlugin) generateMiscFiles(ctx *tilocontext.ExecutionContext, projectPath, projectName string) error {
	engine := templates.NewTemplateEngine()

	// Home template
	processedHome, err := engine.ProcessTemplate(pythonTemplates.DjangoHomeHtml, ctx)
	if err != nil {
		return err
	}

	if err := utils.WriteFile(filepath.Join(projectPath, "templates", "home.html"), processedHome); err != nil {
		return err
	}

	// README.md
	processedReadme, err := engine.ProcessTemplate(pythonTemplates.DjangoReadmeMd, ctx)
	if err != nil {
		return err
	}

	return utils.WriteFile(filepath.Join(projectPath, "README.md"), processedReadme)
}
