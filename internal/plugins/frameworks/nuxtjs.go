package frameworks

import (
	"path/filepath"
	"strings"

	tilocontext "tilokit/internal/core/context"
	"tilokit/internal/plugins/templates"
	"tilokit/internal/templates/common"
	"tilokit/internal/templates/nuxtjs"
	"tilokit/internal/utils"
	"tilokit/pkg/constants"

	"github.com/pkg/errors"
)

// NuxtjsPlugin implements Nuxt.js framework support
type NuxtjsPlugin struct{}

// NewNuxtjsPlugin creates a new Nuxt.js plugin instance
func NewNuxtjsPlugin() *NuxtjsPlugin {
	return &NuxtjsPlugin{}
}

func (p *NuxtjsPlugin) Name() string {
	return "nuxtjs-framework"
}

func (p *NuxtjsPlugin) Version() string {
	return "1.0.0"
}

func (p *NuxtjsPlugin) Description() string {
	return "Nuxt.js Vue framework with SSR/SSG capabilities"
}

func (p *NuxtjsPlugin) SupportedFrameworks() []string {
	return []string{"nuxt"}
}

func (p *NuxtjsPlugin) SupportedBuildTools() []string {
	return []string{"nuxt"}
}

func (p *NuxtjsPlugin) PreGenerate(ctx *tilocontext.ExecutionContext) error {
	// Default to TypeScript if not provided
	if _, ok := ctx.Variables["language"]; !ok {
		ctx.SetVariable("language", "ts")
	}

	// Set Nuxt.js-specific variables with latest versions
	if _, ok := ctx.Variables["nuxt_version"]; !ok {
		ctx.SetVariable("nuxt_version", "^3.15.0")
	}
	if _, ok := ctx.Variables["vue_version"]; !ok {
		ctx.SetVariable("vue_version", "^3.5.0")
	}

	// Set project name from config
	if ctx.Config.ProjectName != "" {
		ctx.SetVariable("project_name", ctx.Config.ProjectName)
	}

	// Set package manager
	if ctx.Config.PackageManager != "" {
		ctx.SetVariable("package_manager", ctx.Config.PackageManager)
	} else {
		ctx.SetVariable("package_manager", "npm")
	}

	return nil
}

func (p *NuxtjsPlugin) Generate(ctx *tilocontext.ExecutionContext) error {
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

func (p *NuxtjsPlugin) PostGenerate(ctx *tilocontext.ExecutionContext) error {
	// Set post-generation metadata
	ctx.SetMetadata("framework_generated", true)
	ctx.SetMetadata("start_command", "npm run dev")

	return nil
}

func (p *NuxtjsPlugin) createDirectoryStructure(ctx *tilocontext.ExecutionContext) error {
	dirs := []string{
		"pages",
		"layouts",
		"components",
		"composables",
		"middleware",
		"plugins",
		"server/api",
		"assets/css",
		"static",
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

func (p *NuxtjsPlugin) generatePackageJson(ctx *tilocontext.ExecutionContext) error {
	// Get language from context
	language := "ts"
	if v, ok := ctx.Variables["language"]; ok {
		language = strings.ToLower(v.(string))
	}

	// Select appropriate package.json template based on language
	var packageJson string
	if language == "js" {
		packageJson = nuxtjs.NuxtjsPackageJsonJS
	} else {
		packageJson = nuxtjs.NuxtjsPackageJsonTS
	}

	// Use template engine with TiLoKit delimiters
	templateEngine := templates.NewTemplateEngine()

	// Process template with variables
	processedPackageJson, err := templateEngine.ProcessTemplateWithDelims(packageJson, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return errors.Wrap(err, "failed to process package.json template")
	}

	// Process common templates
	processedEnv, err := templateEngine.ProcessTemplateWithDelims(common.Env, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return errors.Wrap(err, "failed to process .env template")
	}

	processedGitignore, err := templateEngine.ProcessTemplateWithDelims(common.Gitignore, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return errors.Wrap(err, "failed to process .gitignore template")
	}

	files := map[string]string{
		"package.json": processedPackageJson,
		".env":         processedEnv,
		".gitignore":   processedGitignore,
	}

	for filename, content := range files {
		fullPath := filepath.Join(ctx.ProjectPath, filename)
		if err := utils.WriteFile(fullPath, content); err != nil {
			return err
		}
	}

	return nil
}

func (p *NuxtjsPlugin) generateSourceFiles(ctx *tilocontext.ExecutionContext) error {
	// Use template engine with TiLoKit delimiters
	templateEngine := templates.NewTemplateEngine()

	// Process templates with variables
	processedApp, err := templateEngine.ProcessTemplateWithDelims(nuxtjs.NuxtjsAppVue, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return errors.Wrap(err, "failed to process app.vue template")
	}

	processedIndex, err := templateEngine.ProcessTemplateWithDelims(nuxtjs.NuxtjsIndexPage, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return errors.Wrap(err, "failed to process index page template")
	}

	processedAbout, err := templateEngine.ProcessTemplateWithDelims(nuxtjs.NuxtjsAboutPage, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return errors.Wrap(err, "failed to process about page template")
	}

	processedLayout, err := templateEngine.ProcessTemplateWithDelims(nuxtjs.NuxtjsLayoutDefault, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return errors.Wrap(err, "failed to process layout template")
	}

	files := map[string]string{
		"app.vue":             processedApp,
		"pages/index.vue":     processedIndex,
		"pages/about.vue":     processedAbout,
		"layouts/default.vue": processedLayout,
	}

	for path, content := range files {
		fullPath := filepath.Join(ctx.ProjectPath, path)
		if err := utils.WriteFile(fullPath, content); err != nil {
			return err
		}
	}

	return nil
}

func (p *NuxtjsPlugin) generateConfigFiles(ctx *tilocontext.ExecutionContext) error {
	// Get language from context
	language := "ts"
	if v, ok := ctx.Variables["language"]; ok {
		language = strings.ToLower(v.(string))
	}

	// Nuxt.js configuration
	var nuxtConfig string
	if language == "js" {
		nuxtConfig = `// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  devtools: { enabled: true },
  modules: [
    '@nuxtjs/tailwindcss'
  ],
  css: ['~/assets/css/main.css']
})
`
	} else {
		nuxtConfig = `// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  devtools: { enabled: true },
  modules: [
    '@nuxtjs/tailwindcss'
  ],
  css: ['~/assets/css/main.css'],
  typescript: {
    strict: true
  }
})
`
	}

	// TypeScript configuration (only for TypeScript projects)
	tsConfig := `{
  "extends": "./.nuxt/tsconfig.json"
}
`

	// Tailwind CSS configuration
	tailwindConfig := `/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./components/**/*.{js,vue,ts}",
    "./layouts/**/*.vue",
    "./pages/**/*.vue",
    "./plugins/**/*.{js,ts}",
    "./nuxt.config.{js,ts}",
    "./app.vue"
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
`

	// Main CSS file with Tailwind directives
	mainCSS := `@tailwind base;
@tailwind components;
@tailwind utilities;

/* Custom styles */
body {
  font-family: 'Inter', sans-serif;
}
`

	// Use template engine with TiLoKit delimiters
	templateEngine := templates.NewTemplateEngine()

	// Process templates
	processedNuxtConfig, err := templateEngine.ProcessTemplateWithDelims(nuxtConfig, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return errors.Wrap(err, "failed to process nuxt.config template")
	}

	processedTailwind, err := templateEngine.ProcessTemplateWithDelims(tailwindConfig, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return errors.Wrap(err, "failed to process tailwind.config template")
	}

	processedCSS, err := templateEngine.ProcessTemplateWithDelims(mainCSS, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return errors.Wrap(err, "failed to process main.css template")
	}

	configs := map[string]string{
		"tailwind.config.js":  processedTailwind,
		"assets/css/main.css": processedCSS,
	}

	// Add nuxt.config based on language
	if language == "js" {
		configs["nuxt.config.js"] = processedNuxtConfig
	} else {
		configs["nuxt.config.ts"] = processedNuxtConfig
		processedTsConfig, err := templateEngine.ProcessTemplateWithDelims(tsConfig, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
		if err != nil {
			return errors.Wrap(err, "failed to process tsconfig template")
		}
		configs["tsconfig.json"] = processedTsConfig
	}

	for path, content := range configs {
		fullPath := filepath.Join(ctx.ProjectPath, path)
		if err := utils.WriteFile(fullPath, content); err != nil {
			return err
		}
	}

	return nil
}
