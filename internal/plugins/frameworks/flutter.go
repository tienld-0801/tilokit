package frameworks

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	tilocontext "tilokit/internal/core/context"
	"tilokit/internal/plugins/templates"
	"tilokit/internal/templates/flutter"
	"tilokit/internal/utils"
	"tilokit/pkg/constants"

	"github.com/pkg/errors"
)

// FlutterPlugin implements Flutter framework support
type FlutterPlugin struct {
	sdkFinder *utils.FlutterSDKFinder
	sdkPath   string
}

func NewFlutterPlugin() *FlutterPlugin {
	return &FlutterPlugin{
		sdkFinder: &utils.FlutterSDKFinder{},
	}
}

func (p *FlutterPlugin) Name() string {
	return "flutter"
}

func (p *FlutterPlugin) Version() string {
	return constants.VERSION
}

func (p *FlutterPlugin) Description() string {
	return "Flutter cross-platform mobile framework"
}

func (p *FlutterPlugin) SupportedFrameworks() []string {
	return []string{"flutter"}
}

func (p *FlutterPlugin) SupportedBuildTools() []string {
	return []string{"flutter-cli", "dart"}
}

// ensureFlutterSDK finds and validates Flutter SDK before generating files
func (p *FlutterPlugin) ensureFlutterSDK(ctx *tilocontext.ExecutionContext) error {
	if p.sdkPath != "" {
		// SDK already found and cached
		return nil
	}

	sdkPath, err := p.sdkFinder.FindFlutterSDK()
	if err != nil {
		return fmt.Errorf("flutter SDK detection failed: %w\n\nPlease ensure Flutter is properly installed and either:\n1. Add Flutter to your PATH\n2. Set FLUTTER_ROOT environment variable\n3. Install Flutter in a standard location", err)
	}

	p.sdkPath = sdkPath
	
	// Add Flutter SDK path to context variables for template processing
	ctx.SetVariable("flutter_sdk_path", sdkPath)
	
	// Get Flutter info and add to context
	info := p.sdkFinder.GetFlutterInfo(sdkPath)
	if version := info["version"]; version != "" {
		ctx.SetVariable("flutter_version", version)
	}
	if channel := info["channel"]; channel != "" {
		ctx.SetVariable("flutter_channel", channel)
	}
	if dartVersion := info["dart_version"]; dartVersion != "" {
		ctx.SetVariable("dart_version", dartVersion)
	}
	
	// Try to get Android SDK path using flutter doctor
	if androidSDK := p.sdkFinder.GetAndroidSDKFromFlutterDoctor(sdkPath); androidSDK != "" {
		ctx.SetVariable("android_sdk_path", androidSDK)
	} else {
		// Fallback to environment variables
		p.detectAndroidSDKFromEnv(ctx)
	}

	return nil
}

// GetFlutterBinPath returns the path to Flutter bin directory
func (p *FlutterPlugin) GetFlutterBinPath() string {
	if p.sdkPath == "" {
		return ""
	}
	return filepath.Join(p.sdkPath, "bin")
}

// GetFlutterExecutable returns the full path to Flutter executable
func (p *FlutterPlugin) GetFlutterExecutable() string {
	binPath := p.GetFlutterBinPath()
	if binPath == "" {
		return ""
	}
	
	flutterExe := "flutter"
	if runtime.GOOS == "windows" {
		flutterExe = "flutter.bat"
	}
	
	return filepath.Join(binPath, flutterExe)
}

// Define allowed Flutter commands
var allowedFlutterCommands = map[string]bool{
	"run": true, "build": true, "test": true, "pub": true, 
	"doctor": true, "clean": true, "analyze": true,
}

// runFlutterCommand executes a Flutter command with proper SDK path
func (p *FlutterPlugin) runFlutterCommand(ctx *tilocontext.ExecutionContext, args ...string) error {
	if err := p.ensureFlutterSDK(ctx); err != nil {
		return err
	}

	// Validate command
	if len(args) > 0 && !allowedFlutterCommands[args[0]] {
		return fmt.Errorf("disallowed flutter command: %s", args[0])
	}

	flutterExe := p.GetFlutterExecutable()
	// nolint:gosec // Flutter executable path is validated and controlled
	cmd := exec.Command(flutterExe, args...)
	cmd.Dir = ctx.ProjectPath
	
	// Set environment variables
	cmd.Env = append(os.Environ(), 
		fmt.Sprintf("FLUTTER_ROOT=%s", p.sdkPath),
		fmt.Sprintf("PATH=%s%c%s", p.GetFlutterBinPath(), os.PathListSeparator, os.Getenv("PATH")),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("flutter command failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

func (p *FlutterPlugin) createDirectoryStructure(ctx *tilocontext.ExecutionContext) error {
	dirs := []string{
		// IDE directories
		".idea",
		".idea/libraries",
		".idea/runConfigurations",
		".vscode",

		// Android platform
		"android",
		"android/app",
		"android/app/src",
		"android/app/src/main",
		"android/app/src/main/kotlin",
		"android/app/src/main/kotlin/com/example/my_app",
		"android/app/src/main/java/io/flutter/plugins",
		"android/app/src/main/res",
		"android/app/src/main/res/drawable",
		"android/app/src/main/res/drawable-v21",
		"android/app/src/main/res/mipmap-hdpi",
		"android/app/src/main/res/mipmap-mdpi",
		"android/app/src/main/res/mipmap-xhdpi",
		"android/app/src/main/res/mipmap-xxhdpi",
		"android/app/src/main/res/mipmap-xxxhdpi",
		"android/app/src/main/res/values",
		"android/app/src/main/res/values-night",
		"android/app/src/debug",
		"android/app/src/profile",
		"android/gradle",
		"android/gradle/wrapper",

		// iOS platform
		"ios",
		"ios/Runner",
		"ios/Runner/Assets.xcassets",
		"ios/Runner/Assets.xcassets/AppIcon.appiconset",
		"ios/Runner/Assets.xcassets/LaunchImage.imageset",
		"ios/Runner/Base.lproj",
		"ios/Runner.xcodeproj",
		"ios/Runner.xcodeproj/project.xcworkspace",
		"ios/Runner.xcodeproj/project.xcworkspace/xcshareddata",
		"ios/Runner.xcodeproj/xcshareddata",
		"ios/Runner.xcodeproj/xcshareddata/xcschemes",
		"ios/Runner.xcworkspace",
		"ios/Runner.xcworkspace/xcshareddata",

		// Linux platform
		"linux",
		"linux/flutter",

		// macOS platform
		"macos",
		"macos/Runner",
		"macos/Runner/Assets.xcassets",
		"macos/Runner/Assets.xcassets/AppIcon.appiconset",
		"macos/Runner/Base.lproj",
		"macos/Runner/Configs",
		"macos/Runner.xcodeproj",
		"macos/Runner.xcodeproj/project.xcworkspace",
		"macos/Runner.xcodeproj/project.xcworkspace/xcshareddata",
		"macos/Runner.xcodeproj/xcshareddata",
		"macos/Runner.xcodeproj/xcshareddata/xcschemes",
		"macos/Runner.xcworkspace",
		"macos/Runner.xcworkspace/xcshareddata",

		// Web platform
		"web",

		// Windows platform
		"windows",
		"windows/flutter",
		"windows/runner",

		// Main Flutter directories
		"lib",
		"lib/src",
		"lib/models",
		"lib/services",
		"lib/widgets",
		"lib/screens",
		"lib/utils",

		// Test directories
		"test",
		"test/widget_test",
		"integration_test",

		// Build directories
		"build",

		// Assets directories
		"assets",
		"assets/images",
		"assets/fonts",
		"assets/icons",

		// Documentation
		"docs",
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(ctx.ProjectPath, dir)
		if err := utils.EnsureDir(dirPath); err != nil {
			return err
		}
	}

	return nil
}

func (p *FlutterPlugin) generateModuleFiles(ctx *tilocontext.ExecutionContext) error {
	// Ensure Flutter SDK is available for template processing
	if err := p.ensureFlutterSDK(ctx); err != nil {
		return err
	}

	// Define file contents for Flutter project
	fileContents := map[string]string{
		// src files
		"pubspec.yaml":          flutter.PubspecYaml,
		"analysis_options.yaml": flutter.AnalysisOptions,
		"README.md":             flutter.Readme,
		"CHANGELOG.md":          flutter.Changelog,
		"my_app.iml":            flutter.MyAppIml,
		".metadata":             flutter.MetaData,
		".gitignore":            flutter.Gitignore,
	}

	templateEngine := templates.NewTemplateEngine()

	for filename, content := range fileContents {
		fullPath := filepath.Join(ctx.ProjectPath, filename)

		// Process template content with TILOKit delimiters
		// Note: You'll need to pass Flutter SDK info to template engine if needed
		processedContent, err := templateEngine.ProcessTemplateWithDelims(
			content,
			constants.TiloLeftDelim,
			constants.TiloRightDelim,
			ctx,
		)
		if err != nil {
			return errors.Wrapf(err, "failed to process template for file %s", filename)
		}

		// Ensure the directory exists before writing the file
		if err := utils.EnsureDir(filepath.Dir(fullPath)); err != nil {
			return errors.Wrapf(err, "failed to create directory for file %s", filename)
		}

		if err := utils.WriteFile(fullPath, processedContent); err != nil {
			return errors.Wrapf(err, "failed to write file %s", filename)
		}
	}

	return nil
}

// Your enhanced generateSourceFiles method with SDK detection
func (p *FlutterPlugin) generateSourceFiles(ctx *tilocontext.ExecutionContext) error {
	// First, ensure Flutter SDK is available
	if err := p.ensureFlutterSDK(ctx); err != nil {
		return err
	}

	// Define file groups with their base paths and contents 
	fileGroups := map[string]map[string]string{
		// lib files
		"lib": {
			"main.dart": flutter.MainDart,
		},
		// android files
		"android/app/src/debug": {
			"AndroidManifest.xml": flutter.AndroidDebugManifest,
		},
		"android/app/src/main": {
			"AndroidManifest.xml": flutter.AndroidManifest,
		},
		"android/app/src/main/java/io/flutter/plugins": {
			"GeneratedPluginRegistrant.java": flutter.AndroidGeneratedPluginRegistrant,
		},
		"android/app/src/main/kotlin/com/example/my_app": {
			"MainActivity.kt": flutter.AndroidKotlinMainActivity,
		},
		"android/app/src/main/res/drawable": {
			"launch_background.xml": flutter.DrawableLaunchBackground,
		},
		"android/app/src/main/res/drawable-v21": {
			"launch_background.xml": flutter.DrawableV21LaunchBackground,
		},
		"android/app/src/main/res/values": {
			"styles.xml": flutter.AndroidStyleValue,
		},
		"android/app/src/main/res/values-night": {
			"styles.xml": flutter.AndroidStyleNightValue,
		},
		"android/app/src/profile": {
			"AndroidManifest.xml": flutter.AndroidProfileManifest,
		},
		"android/app": {
			"build.gradle.kts": flutter.AndroidBuildGradle,
		},
		"android": {
			"build.gradle.kts":  flutter.AndroidRootBuildGradle,
			".gitignore":        flutter.AndroidGitignore,
			"gradle.properties": flutter.AndroidGradleProperties,
			"gradlew":           flutter.AndroidGradlew,
			"gradlew.bat":       flutter.AndroidGradlewBat,
			"local.properties":  flutter.AndroidLocalProperties,
		},
		"android/gradle/wrapper": {
			"gradle-wrapper.properties": flutter.AndroidGradleWrapperProperties,
		},
		// ios files
		"ios/Runner": {
			"Info.plist":        flutter.IosInfoPlist,
			"AppDelegate.swift": flutter.IosAppDelegate,
		},
		// web files
		"web": {
			"index.html":    flutter.WebIndexHtml,
			"manifest.json": flutter.WebManifestJson,
		},
		// windows files
		"windows/runner": {
			"main.cpp": flutter.WindowsMainCpp,
		},
		// macos files
		"macos/Runner": {
			"AppDelegate.swift": flutter.MacosAppDelegate,
		},
		// linux files
		"linux": {
			"main.cpp": flutter.LinuxMainCpp,
		},
		// test files
		"test": {
			"widget_test.dart": flutter.WidgetTest,
		},
		"integration_test": {
			"app_test.dart": flutter.IntegrationTest,
		},
	}

	templateEngine := templates.NewTemplateEngine()

	// Note: If you need Flutter SDK info in templates, you'll need to modify
	// the template engine or context to support passing this data

	for basePath, files := range fileGroups {
		if err := p.processFileGroup(ctx, templateEngine, basePath, files); err != nil {
			return fmt.Errorf("failed to process file group %s: %w", basePath, err)
		}
	}

	// Copy assets from the framework's asset directory
	if err := p.copyAssets(ctx); err != nil {
		return fmt.Errorf("failed to copy assets: %w", err)
	}

	// After generating files, optionally run flutter pub get
	if err := p.runFlutterCommand(ctx, "pub", "get"); err != nil {
		// Warning: Failed to run flutter pub get - continuing anyway
		// Don't return error here as file generation was successful
		_ = err // explicitly ignore the error
	}

	// Flutter project files generated successfully
	return nil
}


// processFileGroup handles the template processing and file writing for a group of files
func (p *FlutterPlugin) processFileGroup(
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

func (p *FlutterPlugin) copyAssets(ctx *tilocontext.ExecutionContext) error {
	// Define the source and destination paths for assets
	assetMappings := map[string]string{
		"assets/flutter/mipmap-hdpi/ic_launcher.png":    "android/app/src/main/res/mipmap-hdpi/ic_launcher.png",
		"assets/flutter/mipmap-mdpi/ic_launcher.png":    "android/app/src/main/res/mipmap-mdpi/ic_launcher.png",
		"assets/flutter/mipmap-xhdpi/ic_launcher.png":   "android/app/src/main/res/mipmap-xhdpi/ic_launcher.png",
		"assets/flutter/mipmap-xxhdpi/ic_launcher.png":  "android/app/src/main/res/mipmap-xxhdpi/ic_launcher.png",
		"assets/flutter/mipmap-xxxhdpi/ic_launcher.png": "android/app/src/main/res/mipmap-xxxhdpi/ic_launcher.png",
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

// ValidateFlutterInstallation checks if Flutter is properly installed and configured
func (p *FlutterPlugin) ValidateFlutterInstallation(ctx *tilocontext.ExecutionContext) error {
	if err := p.ensureFlutterSDK(ctx); err != nil {
		return err
	}

	// Run flutter doctor to check installation
	return p.runFlutterCommand(ctx, "doctor")
}

// GetFlutterSDKPath returns the Flutter SDK path
func (p *FlutterPlugin) GetFlutterSDKPath(ctx *tilocontext.ExecutionContext) (string, error) {
	if err := p.ensureFlutterSDK(ctx); err != nil {
		return "", err
	}
	return p.sdkPath, nil
}

// GetFlutterInfo returns comprehensive Flutter SDK information
func (p *FlutterPlugin) GetFlutterInfo(ctx *tilocontext.ExecutionContext) (map[string]string, error) {
	if err := p.ensureFlutterSDK(ctx); err != nil {
		return nil, err
	}

	return p.sdkFinder.GetFlutterInfo(p.sdkPath), nil
}

// detectAndroidSDKFromEnv detects Android SDK from environment variables (fallback)
func (p *FlutterPlugin) detectAndroidSDKFromEnv(ctx *tilocontext.ExecutionContext) {
	// Check environment variables in order of preference
	envVars := []string{"ANDROID_SDK_ROOT", "ANDROID_HOME"}
	
	for _, envVar := range envVars {
		if androidSDK := os.Getenv(envVar); androidSDK != "" {
			if p.sdkFinder.IsValidAndroidSDK(androidSDK) {
				ctx.SetVariable("android_sdk_path", androidSDK)
				return
			}
		}
	}
}

// PostGenerationSetup runs Flutter commands after file generation
func (p *FlutterPlugin) PostGenerationSetup(ctx *tilocontext.ExecutionContext) error {
	// Run flutter pub get to fetch dependencies
	if err := p.runFlutterCommand(ctx, "pub", "get"); err != nil {
		return fmt.Errorf("failed to run 'flutter pub get': %w", err)
	}

	// Optionally run flutter analyze to check for issues
	if err := p.runFlutterCommand(ctx, "analyze"); err != nil {
		// Flutter analyze found issues - continuing anyway
		// Don't return error as this is just a warning
		_ = err // explicitly ignore the error
	}

	return nil
}

func (p *FlutterPlugin) PreGenerate(ctx *tilocontext.ExecutionContext) error {
	// Validate Flutter SDK early in the process
	if err := p.ensureFlutterSDK(ctx); err != nil {
		return errors.Wrap(err, "Flutter SDK validation failed")
	}

	// Add Flutter SDK path to context variables for template processing
	ctx.SetVariable("flutter_sdk_path", p.sdkPath)

	// Flutter SDK validated successfully
	return nil
}

func (p *FlutterPlugin) Generate(ctx *tilocontext.ExecutionContext) error {
	// Create directory structure
	if err := p.createDirectoryStructure(ctx); err != nil {
		return errors.Wrap(err, "failed to create directory structure")
	}
	
	if err := p.generateModuleFiles(ctx); err != nil {
		return errors.Wrap(err, "failed to create moduleFiles")
	}
	
	if err := p.generateSourceFiles(ctx); err != nil {
		return errors.Wrap(err, "failed to create source files")
	}

	return nil
}

func (p *FlutterPlugin) PostGenerate(ctx *tilocontext.ExecutionContext) error {
	// Run post-generation Flutter setup
	if err := p.PostGenerationSetup(ctx); err != nil {
		// Post-generation setup had issues - continuing anyway
		// Don't fail the entire process for post-generation issues
		_ = err // explicitly ignore the error
	}

	return nil
}