package frameworks

import (
	"path/filepath"
	"strings"

	tilocontext "tilokit/internal/core/context"
	"tilokit/internal/plugins/templates"
	"tilokit/internal/templates/common"
	"tilokit/internal/templates/nextjs"
	"tilokit/internal/utils"
	"tilokit/pkg/constants"

	"github.com/pkg/errors"
)

// NextjsPlugin implements Next.js framework support
type NextjsPlugin struct{}

// NewNextjsPlugin creates a new Next.js plugin instance
func NewNextjsPlugin() *NextjsPlugin {
	return &NextjsPlugin{}
}

func (p *NextjsPlugin) Name() string {
	return "nextjs-framework"
}

func (p *NextjsPlugin) Version() string {
	return constants.VERSION
}

func (p *NextjsPlugin) Description() string {
	return "Next.js React framework with SSR/SSG capabilities"
}

func (p *NextjsPlugin) SupportedFrameworks() []string {
	return []string{"next"}
}

func (p *NextjsPlugin) SupportedBuildTools() []string {
	return []string{"next"}
}

func (p *NextjsPlugin) PreGenerate(ctx *tilocontext.ExecutionContext) error {
	// Next.js only supports TypeScript
	ctx.SetVariable("language", "ts")

	// Set Next.js-specific variables with latest versions
	if _, ok := ctx.Variables["next_version"]; !ok {
		ctx.SetVariable("next_version", "^15.5.0")
	}
	if _, ok := ctx.Variables["react_version"]; !ok {
		ctx.SetVariable("react_version", "^19.0.0")
	}
	if _, ok := ctx.Variables["react_dom_version"]; !ok {
		ctx.SetVariable("react_dom_version", "^19.0.0")
	}

	// Get router type from context (default to app router)
	if _, ok := ctx.Variables["router_type"]; !ok {
		ctx.SetVariable("router_type", constants.AppRouter)
	}

	// Set required template variables
	ctx.SetVariable("project_name", ctx.Config.ProjectName)
	ctx.SetVariable("package_manager", ctx.Config.PackageManager)

	return nil
}

func (p *NextjsPlugin) Generate(ctx *tilocontext.ExecutionContext) error {
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

func (p *NextjsPlugin) PostGenerate(ctx *tilocontext.ExecutionContext) error {
	// Set post-generation metadata
	ctx.SetMetadata("framework_generated", true)
	ctx.SetMetadata("start_command", "npm run dev")

	return nil
}

func (p *NextjsPlugin) createDirectoryStructure(ctx *tilocontext.ExecutionContext) error {
	routerType := constants.AppRouter
	if v, ok := ctx.Variables["router_type"].(string); ok && v != "" {
		routerType = v
	}

	var dirs []string
	if routerType == constants.AppRouter {
		// App Router structure
		dirs = []string{
			constants.AppRouter,
			"app/about",
			"components",
			"lib",
			"public",
			"utils",
		}
	} else {
		// Pages Router structure
		dirs = []string{
			"pages",
			"pages/api",
			"components",
			"styles",
			"public",
			"lib",
			"utils",
		}
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(ctx.ProjectPath, dir)
		if err := utils.EnsureDir(dirPath); err != nil {
			return err
		}
	}

	return nil
}

func (p *NextjsPlugin) generatePackageJson(ctx *tilocontext.ExecutionContext) error {
	routerType := constants.AppRouter
	if v, ok := ctx.Variables["router_type"].(string); ok && v != "" {
		routerType = v
	}

	// Next.js only supports TypeScript
	var packageJson string
	if routerType == constants.AppRouter {
		packageJson = nextjs.NextjsAppPackageJsonTS
	} else {
		packageJson = nextjs.NextjsPagesPackageJsonTS
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

func (p *NextjsPlugin) generateSourceFiles(ctx *tilocontext.ExecutionContext) error {
	routerType := "app"
	if v, ok := ctx.Variables["router_type"].(string); ok && v != "" {
		routerType = v
	}

	templateEngine := templates.NewTemplateEngine()
	var files map[string]string

	if routerType == "app" {
		// App Router files - TypeScript only
		files = map[string]string{
			"app/layout.tsx":     nextjs.NextjsAppLayoutPage,
			"app/page.tsx":       nextjs.NextjsAppHomePage,
			"app/about/page.tsx": nextjs.NextjsAppAboutPage,
			"app/globals.css":    nextjs.NextjsAppGlobalCSS,
		}
	} else {
		// Pages Router files - TypeScript only
		files = map[string]string{
			"pages/index.tsx":    nextjs.NextjsPagesIndexPage,
			"pages/_app.tsx":     nextjs.NextjsPagesAppPage,
			"pages/about.tsx":    nextjs.NextjsPagesAboutPage,
			"styles/globals.css": nextjs.NextjsPagesGlobalCSS,
		}
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

func (p *NextjsPlugin) generateConfigFiles(ctx *tilocontext.ExecutionContext) error {
	routerType := constants.AppRouter
	if v, ok := ctx.Variables["router_type"].(string); ok && v != "" {
		routerType = v
	}

	lang := "ts"
	if v, ok := ctx.Variables["language"].(string); ok && v != "" {
		lang = strings.ToLower(v)
	}

	// Next.js configuration
	nextConfig := `/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  swcMinify: true,
}

module.exports = nextConfig
`

	// TypeScript configuration
	tsConfig := `{
  "compilerOptions": {
    "target": "es5",
    "lib": ["dom", "dom.iterable", "es6"],
    "allowJs": true,
    "skipLibCheck": true,
    "strict": true,
    "forceConsistentCasingInFileNames": true,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "node",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "preserve",
    "incremental": true,
    "plugins": [
      {
        "name": "next"
      }
    ],
    "paths": {
      "@/*": ["./*"]
    }
  },
  "include": ["next-env.d.ts", "**/*.ts", "**/*.tsx", ".next/types/**/*.ts"],
  "exclude": ["node_modules"]
}
`

	// Base configs
	configs := map[string]string{
		constants.EnvFileName:       common.Env,
		constants.GitignoreFileName: common.Gitignore,
		"next.config.js":            nextConfig,
		"postcss.config.js":         nextjs.NextjsAppPostCSSConfig,
	}

	// Add TypeScript config only if using TypeScript
	if lang == "ts" {
		configs[constants.TsConfigFileName] = tsConfig
	}

	// Choose Tailwind config based on router type
	if routerType == constants.AppRouter {
		configs["tailwind.config.js"] = nextjs.NextjsAppTailwindConfig
	} else {
		configs["tailwind.config.js"] = nextjs.NextjsPagesTailwindConfig
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
