package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ti-lo/tilokit/internal/core/engine"
	"github.com/ti-lo/tilokit/internal/utils"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available frameworks, build tools, and plugins",
	Long:  "Display all available frameworks, build tools, and plugins that can be used with TiLoKit",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runList()
	},
}

var (
	listAll bool
)

func runList() error {
	utils.PrintBanner()
	
	// Initialize engine and register plugins to get accurate lists
	eng := engine.New()
	if err := registerPlugins(eng); err != nil {
		utils.Warning("Failed to load plugins: %v", err)
	}

	if listAll || listFrameworks {
		utils.Info("🚀 Supported Frameworks:")
		frameworks := map[string][]string{
			"JavaScript/TypeScript": {"react", "vue", "angular", "svelte", "nextjs", "nuxtjs"},
			"Python": {"django", "flask", "fastapi"},
			"PHP": {"laravel", "symfony"},
			"Java": {"spring-boot", "quarkus"},
			"Go": {"gin", "echo", "fiber"},
			"Rust": {"actix", "rocket", "axum"},
			"C#": {"aspnetcore", "blazor"},
			"Ruby": {"rails", "sinatra"},
			"Node.js": {"express", "nestjs", "fastify"},
			"Mobile": {"react-native", "flutter", "ionic"},
			"Desktop": {"electron", "tauri", "wails"},
		}
		for category, fws := range frameworks {
			utils.Log("  %s:", category)
			for _, fw := range fws {
				utils.Log("    • %s", fw)
			}
		}
		fmt.Println()
	}

	if listAll || listBuildTools {
		utils.Info("🔧 Supported Build Tools:")
		buildTools := map[string][]string{
			"JavaScript": {"vite", "webpack", "rollup", "parcel"},
			"Package Managers": {"npm", "yarn", "pnpm"},
			"Python": {"pip", "poetry", "pipenv"},
			"PHP": {"composer"},
			"Java": {"maven", "gradle"},
			"Go": {"go-modules"},
			"Rust": {"cargo"},
			"C#": {"dotnet"},
			"Ruby": {"bundler", "gem"},
			"Mobile/Desktop": {"metro", "expo", "flutter-cli", "electron-builder"},
		}
		for category, tools := range buildTools {
			utils.Log("  %s:", category)
			for _, tool := range tools {
				utils.Log("    • %s", tool)
			}
		}
		fmt.Println()
	}

	if listAll {
		utils.Info("🔌 Available Plugins:")
		
		// Framework plugins
		utils.Log("  Framework Plugins:")
		utils.Log("    • react-framework - React with TypeScript and modern setup")
		utils.Log("    • vue-framework - Vue 3 with Composition API")
		
		// Build tool plugins
		utils.Log("  Build Tool Plugins:")
		utils.Log("    • vite-builder - Vite with optimized configuration")
		
		// Tool plugins
		utils.Log("  Tool Plugins:")
		utils.Log("    • git-integration - Git repository initialization")
		
		fmt.Println()
		
		utils.Info("💡 Usage Examples:")
		utils.Log("  tilokit --name my-app --framework react --build-tool vite")
		utils.Log("  tilokit --name vue-app --framework vue")
		utils.Log("  tilokit --list-frameworks")
		utils.Log("  tilokit --list-build-tools")
	}

	return nil
}

func init() {
	rootCmd.AddCommand(listCmd)
	
	listCmd.Flags().BoolVar(&listAll, "all", false, "List all available options")
	listCmd.Flags().BoolVar(&listFrameworks, "frameworks", false, "List available frameworks")
	listCmd.Flags().BoolVar(&listBuildTools, "build-tools", false, "List available build tools")
}
