package registry

import (
	"fmt"

	"tilokit/internal/core/engine"
	"tilokit/internal/core/registry"
	"tilokit/internal/plugins/builders"
	"tilokit/internal/plugins/frameworks"
	"tilokit/internal/plugins/tools"
)

// PluginRegistry handles plugin registration for the engine
type PluginRegistry struct{}

// NewPluginRegistry creates a new PluginRegistry instance
func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{}
}

// RegisterAllPlugins registers all available plugins with the engine
func (r *PluginRegistry) RegisterAllPlugins(eng *engine.Engine) error {
	plugins := []registry.Plugin{
		// JavaScript/TypeScript Frameworks
		frameworks.NewReactPlugin(),
		frameworks.NewVuePlugin(),
		frameworks.NewSveltePlugin(),
		frameworks.NewAngularPlugin(),
		frameworks.NewNextjsPlugin(),
		frameworks.NewNuxtjsPlugin(),

		// Backend Frameworks
		// Python
		frameworks.NewPythonDjangoPlugin(),
		frameworks.NewPythonFlaskPlugin(),
		frameworks.NewPythonFastAPIPlugin(),

		// PHP
		frameworks.NewPHPLaravelPlugin(),
		frameworks.NewPHPSymfonyPlugin(),
		frameworks.NewPHPCakePlugin(),
		frameworks.NewPHPCodeIgniterPlugin(),

		// Java
		frameworks.NewJavaSpringBootPlugin(),
		frameworks.NewJavaQuarkusPlugin(),

		// Go
		frameworks.NewGoGinPlugin(),
		frameworks.NewGoEchoPlugin(),
		frameworks.NewGoFiberPlugin(),

		// Rust
		frameworks.NewRustActixPlugin(),
		frameworks.NewRustRocketPlugin(),
		frameworks.NewRustAxumPlugin(),

		// C#
		frameworks.NewCSharpASPNetCorePlugin(),
		frameworks.NewCSharpBlazorPlugin(),

		// Ruby
		frameworks.NewRubyRailsPlugin(),
		frameworks.NewRubySinatraPlugin(),

		// Node.js
		frameworks.NewNodeExpressPlugin(),
		frameworks.NewNodeNestJSPlugin(),
		frameworks.NewNodeFastifyPlugin(),

		// Mobile Frameworks
		frameworks.NewReactNativePlugin(),
		frameworks.NewFlutterPlugin(),
		frameworks.NewIonicPlugin(),

		// Desktop Frameworks
		frameworks.NewElectronPlugin(),
		frameworks.NewTauriPlugin(),
		frameworks.NewWailsPlugin(),

		// Build Tools
		builders.NewVitePlugin(),
		builders.NewWebpackPlugin(),
		builders.NewRollupPlugin(),

		// Tools
		tools.NewGitPlugin(),
	}

	// Register all plugins with error handling
	for _, plugin := range plugins {
		if err := eng.RegisterPlugin(plugin); err != nil {
			return fmt.Errorf("failed to register plugin %s: %w", plugin.Name(), err)
		}
	}

	return nil
}

// GetBuildToolsForFramework returns supported build tools for a given framework
func (r *PluginRegistry) GetBuildToolsForFramework(framework string) []string {
	buildToolMap := map[string][]string{
		"react":        {"vite", "webpack", "rollup"},
		"vue":          {"vite", "webpack"},
		"svelte":       {"vite"},
		"angular":      {"angular-cli"},
		"react-native": {"expo"},
		"next":         {"next"},
		"nuxt":         {"nuxt"},
	}
	if tools, exists := buildToolMap[framework]; exists {
		return tools
	}
	return []string{"vite"}
}

// GetDefaultBuildTool returns the default build tool for a given framework
func (r *PluginRegistry) GetDefaultBuildTool(framework string) string {
	defaults := map[string]string{
		"django":      "pip",
		"flask":       "pip",
		"fastapi":     "pip",
		"spring-boot": "maven",
		"quarkus":     "maven",
		"rails":       "bundler",
		"gin":         "go-modules",
		"echo":        "go-modules",
		"fiber":       "go-modules",
		"laravel":     "composer",
		"symfony":     "composer",
		"next":        "next",
		"nuxt":        "nuxt",
	}
	if tool, exists := defaults[framework]; exists {
		return tool
	}
	return "vite"
}
