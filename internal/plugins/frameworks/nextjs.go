package frameworks

import (
	"path/filepath"
	"strings"

	tilocontext "tilokit/internal/core/context"
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
	return "1.0.0"
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
	// Set Next.js-specific variables
	ctx.SetVariable("next_version", "14.0.0")
	ctx.SetVariable("react_version", "^18.0.0")
	ctx.SetVariable("typescript_support", true)

	// Get router type from context (default to app router)
	routerType, exists := ctx.GetVariable("router_type")
	if !exists || routerType == nil {
		routerType = constants.AppRouter
	}
	ctx.SetVariable("router_type", routerType)

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
	routerType, _ := ctx.GetVariable("router_type")
	routerTypeStr := routerType.(string)

	var dirs []string
	if routerTypeStr == constants.AppRouter {
		// App Router structure
		dirs = []string{
			"app",
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
	routerType, _ := ctx.GetVariable("router_type")
	routerTypeStr := routerType.(string)
	envContent := common.Env
	gitIgnoreContent := common.Gitignore

	// Choose package.json based on router type
	var packageJson string
	if routerTypeStr == constants.AppRouter {
		packageJson = nextjs.NextjsAppPackageJson
	} else {
		packageJson = nextjs.NextjsPagesPackageJson
	}

	// Replace project name placeholder
	packageJson = strings.ReplaceAll(packageJson, "{{.ProjectName}}", ctx.Config.ProjectName)

	files := map[string]string{
		"package.json": packageJson,
		".env.local":   envContent,
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

func (p *NextjsPlugin) generateSourceFiles(ctx *tilocontext.ExecutionContext) error {
	routerType, _ := ctx.GetVariable("router_type")
	routerTypeStr := routerType.(string)
	projectName := ctx.Config.ProjectName

	var files map[string]string

	if routerTypeStr == "app" {
		// App Router files
		layoutContent := strings.ReplaceAll(nextjs.NextjsAppLayout, "{{.ProjectName}}", projectName)
		pageContent := strings.ReplaceAll(nextjs.NextjsAppPage, "{{.ProjectName}}", projectName)
		aboutContent := strings.ReplaceAll(nextjs.NextjsAppAboutPage, "{{.ProjectName}}", projectName)
		globalCSS := nextjs.NextjsAppGlobalCSS

		files = map[string]string{
			"app/layout.tsx":     layoutContent,
			"app/page.tsx":       pageContent,
			"app/about/page.tsx": aboutContent,
			"app/globals.css":    globalCSS,
		}
	} else {
		// Pages Router files
		indexContent := strings.ReplaceAll(nextjs.NextjsPagesIndexPage, "{{.ProjectName}}", projectName)
		appContent := nextjs.NextjsPagesAppPage
		aboutContent := strings.ReplaceAll(nextjs.NextjsPagesAboutPage, "{{.ProjectName}}", projectName)
		globalCSS := nextjs.NextjsPagesGlobalCSS

		files = map[string]string{
			"pages/index.tsx":    indexContent,
			"pages/_app.tsx":     appContent,
			"pages/about.tsx":    aboutContent,
			"styles/globals.css": globalCSS,
		}
	}

	for path, content := range files {
		fullPath := filepath.Join(ctx.ProjectPath, path)
		if err := utils.WriteFile(fullPath, content); err != nil {
			return err
		}
	}

	return nil
}

func (p *NextjsPlugin) generateConfigFiles(ctx *tilocontext.ExecutionContext) error {
	routerType, _ := ctx.GetVariable("router_type")
	routerTypeStr := routerType.(string)

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

	// Choose Tailwind config based on router type
	var tailwindConfig, postCSSConfig string
	if routerTypeStr == "app" {
		tailwindConfig = nextjs.NextjsAppTailwindConfig
		postCSSConfig = nextjs.NextjsAppPostCSSConfig
	} else {
		tailwindConfig = nextjs.NextjsPagesTailwindConfig
		postCSSConfig = nextjs.NextjsPagesPostCSSConfig
	}

	configs := map[string]string{
		"next.config.js":     nextConfig,
		"tsconfig.json":      tsConfig,
		"tailwind.config.js": tailwindConfig,
		"postcss.config.js":  postCSSConfig,
	}

	for path, content := range configs {
		fullPath := filepath.Join(ctx.ProjectPath, path)
		if err := utils.WriteFile(fullPath, content); err != nil {
			return err
		}
	}

	return nil
}
