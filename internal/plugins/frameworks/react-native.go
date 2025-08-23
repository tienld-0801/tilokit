package frameworks

import (
	"os"
	"path/filepath"

	tilocontext "tilokit/internal/core/context"
	"tilokit/internal/plugins/templates"
	"tilokit/internal/templates/common"
	reactNative "tilokit/internal/templates/react-native"
	"tilokit/internal/utils"
	"tilokit/pkg/constants"

	"github.com/pkg/errors"
)

// ReactNativePlugin implements React Native framework support
type ReactNativePlugin struct{}

func NewReactNativePlugin() *ReactNativePlugin {
	return &ReactNativePlugin{}
}

func (p *ReactNativePlugin) Name() string {
	return "react-native"
}

func (p *ReactNativePlugin) Version() string {
	return constants.VERSION
}

func (p *ReactNativePlugin) Description() string {
	return "React Native mobile app framework"
}

func (p *ReactNativePlugin) SupportedFrameworks() []string {
	return []string{"react-native", "rn"}
}

func (p *ReactNativePlugin) SupportedBuildTools() []string {
	return []string{"metro", "expo"}
}

func (p *ReactNativePlugin) createDirectoryStructure(ctx *tilocontext.ExecutionContext) error {
	dirs := []string{
		".expo",
		".expo/types",
		"src",
		"src/components",
		"src/components/ui",
		"src/app",
		"src/app/(tabs)",
		"src/assets",
		"src/constants",
		"src/hooks",
		"src/utils",
		"src/styles",
		"src/assets",
		"src/scripts",
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

func (p *ReactNativePlugin) generateModuleFiles(ctx *tilocontext.ExecutionContext) error {
	// Define file contents
	fileContents := map[string]string{
		constants.PackageJsonFileName:  reactNative.PackageJson,
		constants.EnvFileName:          common.Env,
		constants.GitignoreFileName:    reactNative.ExpoGitignore,
		constants.EslintConfigFileName: reactNative.Eslint,
		constants.TsConfigFileName:     reactNative.TsConfig,
		constants.AppConfigFileName:    reactNative.AppConfig,
		constants.ExpoEnvFileName:      reactNative.ExpoEnv,
	}

	templateEngine := templates.NewTemplateEngine()

	for filename, content := range fileContents {
		fullPath := filepath.Join(ctx.ProjectPath, filename)

		// Process template content with TILOKit delimiters
		processedContent, err := templateEngine.ProcessTemplateWithDelims(
			content,
			constants.TiloLeftDelim,
			constants.TiloRightDelim,
			ctx,
		)
		if err != nil {
			return errors.Wrapf(err, "failed to process template for file %s", filename)
		}

		if err := utils.WriteFile(fullPath, processedContent); err != nil {
			return errors.Wrapf(err, "failed to write file %s", filename)
		}
	}

	return nil
}

func (p *ReactNativePlugin) generateSourceFiles(ctx *tilocontext.ExecutionContext) error {
	// Define file groups with their base paths and contents
	fileGroups := map[string]map[string]string{
		"src/app": {
			"_layout.tsx":    reactNative.ExpoMainLayout,
			"+not-found.tsx": reactNative.ExpoNotFound,
		},
		"src/app/(tabs)": {
			"_layout.tsx": reactNative.ExpoTabsLayout,
			"explore.tsx": reactNative.ExpoTabsExplore,
			"index.tsx":   reactNative.ExpoTabsMainFile,
		},
		".expo": {
			"devices.json": reactNative.ExpoDevices,
			"README.md":    reactNative.ExpoReadMe,
		},
		".expo/types": {
			"router.d.ts": reactNative.ExpoRouterTypes,
		},
		"src/components": {
			"Collapsible.tsx":        reactNative.ExpoCollapsible,
			"ExternalLink.tsx":       reactNative.ExpoExternalLink,
			"HapticTab.tsx":          reactNative.ExpoTabsMainFile,
			"ParallaxScrollView.tsx": reactNative.ExpoParallaxScrollView,
			"ThemedText.tsx":         reactNative.ExpoThemedText,
			"ThemedView.tsx":         reactNative.ExpoThemedView,
		},
		"src/components/ui": {
			"IconSymbol.ios.tsx":       reactNative.ExpoIconSymbolIos,
			"IconSymbol.tsx":           reactNative.ExpoIconSymbol,
			"TabBarBackground.ios.tsx": reactNative.ExpoTabBarBackgroundIos,
			"TabBarBackground.tsx":     reactNative.ExpoTabBarBackground,
		},
		"src/constants": {
			"Colors.ts": reactNative.ExpoIconSymbolIos,
		},
		"src/hooks": {
			"useColorScheme.ts": reactNative.ExpoUseColorScheme,
		},
		"src/scripts": {
			"reset-project.js": reactNative.ExpoResetProject,
		},
	}

	templateEngine := templates.NewTemplateEngine()

	for basePath, files := range fileGroups {
		if err := p.processFileGroup(ctx, templateEngine, basePath, files); err != nil {
			return err
		}
	}

	// Copy assets from the framework's asset directory
	if err := p.copyAssets(ctx); err != nil {
		return err
	}

	return nil
}

// processFileGroup handles the template processing and file writing for a group of files
func (p *ReactNativePlugin) processFileGroup(
	ctx *tilocontext.ExecutionContext,
	templateEngine *templates.TemplateEngine,
	basePath string,
	files map[string]string,
) error {
	for filename, content := range files {
		fullPath := filepath.Join(ctx.ProjectPath, basePath, filename)

		processedContent, err := templateEngine.ProcessTemplateWithDelims(
			content,
			constants.TiloLeftDelim,
			constants.TiloRightDelim,
			ctx,
		)
		if err != nil {
			return errors.Wrapf(err, "failed to process template for %s/%s", basePath, filename)
		}

		if err := utils.WriteFile(fullPath, processedContent); err != nil {
			return errors.Wrapf(err, "failed to write file %s/%s", basePath, filename)
		}
	}

	return nil
}

// copyAssets copies assets from the framework's asset directory to the project's assets directory
func (p *ReactNativePlugin) copyAssets(ctx *tilocontext.ExecutionContext) error {
	// Define the source and destination paths for assets
	assetMappings := map[string]string{
		"assets/react-native/expo/icon.png":               "src/assets/images/icon.png",
		"assets/react-native/expo/adaptive-icon.png":      "src/assets/images/adaptive-icon.png",
		"assets/react-native/expo/favicon.png":            "src/assets/images/favicon.png",
		"assets/react-native/expo/splash-icon.png":        "src/assets/images/splash-icon.png",
		"assets/react-native/expo/react-logo.png":         "src/assets/images/react-logo.png",
		"assets/react-native/expo/react-logo@2x.png":      "src/assets/images/react-logo@2x.png",
		"assets/react-native/expo/react-logo@3x.png":      "src/assets/images/react-logo@3x.png",
		"assets/react-native/expo/partial-react-logo.png": "src/assets/images/partial-react-logo.png",
		"assets/react-native/expo/SpaceMono-Regular.ttf":  "src/assets/fonts/SpaceMono-Regular.ttf",
	}

	// Get the current working directory (assuming the framework is running from the workspace root)
	workspaceRoot, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "failed to get current working directory")
	}

	for srcPath, destPath := range assetMappings {
		srcFullPath := filepath.Join(workspaceRoot, srcPath)
		destFullPath := filepath.Join(ctx.ProjectPath, destPath)

		// Check if source file exists
		if !utils.FileExists(srcFullPath) {
			continue // Skip if source doesn't exist
		}

		// Copy the file
		if err := utils.CopyFile(srcFullPath, destFullPath); err != nil {
			return errors.Wrapf(err, "failed to copy asset %s to %s", srcPath, destPath)
		}
	}

	return nil
}

func (p *ReactNativePlugin) PreGenerate(ctx *tilocontext.ExecutionContext) error {
	// Set React-specific variables
	ctx.SetVariable("react_version", "^18.2.0")
	ctx.SetVariable("react_dom_version", "^18.2.0")
	ctx.SetVariable("typescript_support", true)

	return nil
}

func (p *ReactNativePlugin) Generate(ctx *tilocontext.ExecutionContext) error {
	// Create directory structure
	if err := p.createDirectoryStructure(ctx); err != nil {
		return errors.Wrap(err, "failed to create directory structure")
	}

	if err := p.generateModuleFiles(ctx); err != nil {
		return errors.Wrap(err, "failed to create module files")
	}

	if err := p.generateSourceFiles(ctx); err != nil {
		return errors.Wrap(err, "failed to create source files")
	}

	return nil
}

func (p *ReactNativePlugin) PostGenerate(ctx *tilocontext.ExecutionContext) error {
	// TODO: Implement React Native post-generation logic
	return nil
}
