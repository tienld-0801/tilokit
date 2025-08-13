package frameworks

import (
	"path/filepath"
	"strings"

	tilocontext "tilokit/internal/core/context"
	"tilokit/internal/templates/nuxtjs"
	"tilokit/internal/templates/common"
	"tilokit/internal/utils"

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
	// Set Nuxt.js-specific variables
	ctx.SetVariable("nuxt_version", "^3.8.0")
	ctx.SetVariable("vue_version", "^3.3.0")
	ctx.SetVariable("typescript_support", true)

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
	ctx.SetMetadata("install_command", "npm install")
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
	packageJson := nuxtjs.NuxtjsPackageJson
	envContent := common.Env
	gitIgnoreContent := common.Gitignore

	// Replace project name placeholder
	packageJson = strings.ReplaceAll(packageJson, "{{.ProjectName}}", ctx.Config.ProjectName)

	files := map[string]string{
		"package.json": packageJson,
		".env":         envContent,
		".gitignore":   gitIgnoreContent,
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
	appContent := nuxtjs.NuxtjsAppVue
	indexContent := nuxtjs.NuxtjsIndexPage
	aboutContent := nuxtjs.NuxtjsAboutPage
	layoutContent := nuxtjs.NuxtjsLayoutDefault

	// Replace project name placeholders
	indexContent = strings.ReplaceAll(indexContent, "{{.ProjectName}}", ctx.Config.ProjectName)
	aboutContent = strings.ReplaceAll(aboutContent, "{{.ProjectName}}", ctx.Config.ProjectName)
	layoutContent = strings.ReplaceAll(layoutContent, "{{.ProjectName}}", ctx.Config.ProjectName)

	files := map[string]string{
		"app.vue":              appContent,
		"pages/index.vue":      indexContent,
		"pages/about.vue":      aboutContent,
		"layouts/default.vue":  layoutContent,
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
	// Nuxt.js configuration
	nuxtConfig := `// https://nuxt.com/docs/api/configuration/nuxt-config
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

	// TypeScript configuration
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

	configs := map[string]string{
		"nuxt.config.ts":      nuxtConfig,
		"tsconfig.json":       tsConfig,
		"tailwind.config.js":  tailwindConfig,
		"assets/css/main.css": mainCSS,
	}

	for path, content := range configs {
		fullPath := filepath.Join(ctx.ProjectPath, path)
		if err := utils.WriteFile(fullPath, content); err != nil {
			return err
		}
	}

	return nil
}
