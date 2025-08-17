package frameworks

import (
	"path/filepath"

	tilocontext "tilokit/internal/core/context"
	"tilokit/internal/plugins/templates"
	"tilokit/internal/templates/common"
	"tilokit/internal/templates/vue"
	"tilokit/internal/utils"

	"github.com/pkg/errors"

	"tilokit/pkg/constants"
)

// VuePlugin implements Vue framework support
type VuePlugin struct{}

// NewVuePlugin creates a new Vue plugin instance
func NewVuePlugin() *VuePlugin {
	return &VuePlugin{}
}

func (p *VuePlugin) Name() string {
	return "vue-framework"
}

func (p *VuePlugin) Version() string {
	return constants.VERSION
}

func (p *VuePlugin) Description() string {
	return "Vue 3 framework with Composition API and modern setup"
}

func (p *VuePlugin) SupportedFrameworks() []string {
	return []string{"vue"}
}

func (p *VuePlugin) SupportedBuildTools() []string {
	return []string{"vite", "webpack"}
}

func (p *VuePlugin) PreGenerate(ctx *tilocontext.ExecutionContext) error {
	ctx.SetVariable("vue_version", "^3.4.0")
	ctx.SetVariable("typescript_support", true)
	return nil
}

func (p *VuePlugin) Generate(ctx *tilocontext.ExecutionContext) error {
	if err := p.createDirectoryStructure(ctx); err != nil {
		return errors.Wrap(err, "failed to create directory structure")
	}

	if err := p.generatePackageJson(ctx); err != nil {
		return errors.Wrap(err, "failed to generate package.json")
	}

	if err := p.generateSourceFiles(ctx); err != nil {
		return errors.Wrap(err, "failed to generate source files")
	}

	if err := p.generateConfigFiles(ctx); err != nil {
		return errors.Wrap(err, "failed to generate config files")
	}

	return nil
}

func (p *VuePlugin) PostGenerate(ctx *tilocontext.ExecutionContext) error {
	ctx.SetMetadata("framework_generated", true)
	ctx.SetMetadata("start_command", "npm run dev")
	return nil
}

func (p *VuePlugin) createDirectoryStructure(ctx *tilocontext.ExecutionContext) error {
	dirs := []string{
		"src",
		"src/components",
		"src/components/icons",
		"src/composables",
		"src/stores",
		"src/views",
		"src/assets",
		"src/styles",
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

func (p *VuePlugin) generatePackageJson(ctx *tilocontext.ExecutionContext) error {
	language, _ := ctx.Variables["language"].(string)
	isJS := language == "js"
	var packageJson string
	if isJS {
		packageJson = vue.ViteJsPackageJson
	} else {
		packageJson = vue.ViteTsPackageJson
	}

	// Use template engine for consistent processing with TILOKit delimiters
	templateEngine := templates.NewTemplateEngine()
	processedContent, err := templateEngine.ProcessTemplateWithDelims(packageJson, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return errors.Wrap(err, "failed to process package.json template")
	}

	packageJsonPath := filepath.Join(ctx.ProjectPath, "package.json")
	return utils.WriteFile(packageJsonPath, processedContent)
}

func (p *VuePlugin) generateSourceFiles(ctx *tilocontext.ExecutionContext) error {
	// Safely extract language, default to "ts" on error
	rawLang, exists := ctx.Variables["language"]
	langStr, ok := rawLang.(string)
	if !exists || !ok {
		langStr = "ts"
	}
	isJS := langStr == "js"

	// Define file templates based on language
	var fileTemplates map[string]string
	if isJS {
		fileTemplates = p.getJavaScriptTemplates()
	} else {
		fileTemplates = p.getTypeScriptTemplates()
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

func (p *VuePlugin) getJavaScriptTemplates() map[string]string {
	return map[string]string{
		"src/main.js":                                vue.ViteJsMainFile,
		"src/App.vue":                                vue.ViteJsAppFile,
		"src/components/HelloWorld.vue":              vue.ViteJsHelloWorldVue,
		"src/components/WelcomeItem.vue":             vue.ViteJsWelcomeItem,
		"src/components/TheWelcome.vue":              vue.ViteJsTheWelcome,
		"src/components/icons/IconDocumentation.vue": vue.ViteIconDocumentation,
		"src/components/icons/IconTooling.vue":       vue.ViteIconTooling,
		"src/components/icons/IconEcosystem.vue":     vue.ViteIconEcoSystem,
		"src/components/icons/IconCommunity.vue":     vue.ViteIconCommunity,
		"src/components/icons/IconSupport.vue":       vue.ViteIconSupport,
		"src/router/index.js":                        vue.ViteJsRouter,
		"src/views/HomeView.vue":                     vue.ViteJsHomeView,
		"src/views/AboutView.vue":                    vue.ViteJsAboutView,
		"src/assets/main.css":                        vue.ViteJsMainCss,
		"src/assets/base.css":                        vue.ViteJsBaseCss,
	}
}

func (p *VuePlugin) getTypeScriptTemplates() map[string]string {
	return map[string]string{
		"src/main.ts":                                vue.ViteTsMainFile,
		"src/App.vue":                                vue.ViteTsAppFile,
		"src/components/HelloWorld.vue":              vue.ViteTSHelloWordVue,
		"src/components/WelcomeItem.vue":             vue.ViteTsWelcomeItem,
		"src/components/TheWelcome.vue":              vue.ViteTsTheWelcome,
		"src/components/icons/IconDocumentation.vue": vue.ViteIconDocumentation,
		"src/components/icons/IconTooling.vue":       vue.ViteIconTooling,
		"src/components/icons/IconEcosystem.vue":     vue.ViteIconEcoSystem,
		"src/components/icons/IconCommunity.vue":     vue.ViteIconCommunity,
		"src/components/icons/IconSupport.vue":       vue.ViteIconSupport,
		"src/router/index.ts":                        vue.ViteTsRouter,
		"src/views/HomeView.vue":                     vue.ViteTsHomeView,
		"src/views/AboutView.vue":                    vue.ViteTsAboutView,
		"src/assets/main.css":                        vue.ViteTsMainCss,
		"src/assets/base.css":                        vue.ViteTsBaseCss,
	}
}

func (p *VuePlugin) generateConfigFiles(ctx *tilocontext.ExecutionContext) error {
	language := ctx.Variables["language"].(string)
	isJS := language == "js"

	var configs map[string]string
	if isJS {
		configs = p.getJavaScriptConfigs()
	} else {
		configs = p.getTypeScriptConfigs()
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

func (p *VuePlugin) getJavaScriptConfigs() map[string]string {
	return map[string]string{
		"index.html": vue.ViteJsIndexHtml,
		".env":       common.Env,
	}
}

func (p *VuePlugin) getTypeScriptConfigs() map[string]string {
	return map[string]string{
		"tsconfig.json": vue.ViteTsConfig,
		"src/env.d.ts":  vue.ViteEnvDTs,
		"index.html":    vue.ViteTsIndexHtml,
		".env":          common.Env,
	}
}
