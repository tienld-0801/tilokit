package frameworks

import (
	"path/filepath"

	tilocontext "tilokit/internal/core/context"
	"tilokit/internal/plugins/templates"
	"tilokit/internal/templates/flutter"
	"tilokit/internal/utils"
	"tilokit/pkg/constants"

	"github.com/pkg/errors"
)

type FlutterPlugin struct {
	sdkFinder *utils.FlutterSDKFinder
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
	return []string{}
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

func (p *FlutterPlugin) generateSourceFiles(ctx *tilocontext.ExecutionContext) error {
	templateEngine := templates.NewTemplateEngine()

	files := map[string]string{
		"lib/main.dart":         flutter.MainDart,
		"pubspec.yaml":          flutter.PubspecYaml,
		"analysis_options.yaml": flutter.AnalysisOptions,
		"README.md":             flutter.Readme,
		"CHANGELOG.md":          flutter.Changelog,
		".metadata":             flutter.MetaData,
		"android/app/src/debug/AndroidManifest.xml":                                   flutter.AndroidDebugManifest,
		"android/app/src/main/AndroidManifest.xml":                                    flutter.AndroidManifest,
		"android/app/src/profile/AndroidManifest.xml":                                 flutter.AndroidProfileManifest,
		"android/app/src/main/kotlin/com/example/my_app/MainActivity.kt":              flutter.AndroidKotlinMainActivity,
		"android/app/src/main/java/io/flutter/plugins/GeneratedPluginRegistrant.java": flutter.AndroidGeneratedPluginRegistrant,
		"android/app/src/main/res/drawable/launch_background.xml":                     flutter.DrawableLaunchBackground,
		"android/app/src/main/res/drawable-v21/launch_background.xml":                 flutter.DrawableV21LaunchBackground,
		"android/app/src/main/res/values/styles.xml":                                  flutter.AndroidStyleValue,
		"android/app/src/main/res/values-night/styles.xml":                            flutter.AndroidStyleNightValue,
		"android/app/build.gradle.kts":                                                flutter.AndroidBuildGradle,
		"android/build.gradle.kts":                                                    flutter.AndroidRootBuildGradle,
		"android/.gitignore":                                                          flutter.AndroidGitignore,
		"android/gradle.properties":                                                   flutter.AndroidGradleProperties,
		"android/gradlew":                                                             flutter.AndroidGradlew,
		"android/gradlew.bat":                                                         flutter.AndroidGradlewBat,
		"android/gradle/wrapper/gradle-wrapper.properties":                            flutter.AndroidGradleWrapperProperties,
	}

	for path, content := range files {
		fullPath := filepath.Join(ctx.ProjectPath, path)

		processedContent, err := templateEngine.ProcessTemplateWithDelims(content, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
		if err != nil {
			return errors.Wrapf(err, "failed to process template for %s", path)
		}

		if err := utils.WriteFile(fullPath, processedContent); err != nil {
			return err
		}
	}

	projectName := ctx.Variables["project_name"].(string)
	imlPath := filepath.Join(ctx.ProjectPath, projectName+".iml")

	processedIml, err := templateEngine.ProcessTemplateWithDelims(flutter.MyAppIml, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to process template for %s.iml", projectName)
	}

	if err := utils.WriteFile(imlPath, processedIml); err != nil {
		return err
	}

	return nil
}

func (p *FlutterPlugin) generatePubspec(ctx *tilocontext.ExecutionContext) error {
	templateEngine := templates.NewTemplateEngine()

	pubspecTemplate := `name: {{.project_name}}
description: A new Flutter project.
publish_to: 'none'

version: 1.0.0+1

environment:
  sdk: '>=3.0.0 <4.0.0'

dependencies:
  flutter:
    sdk: flutter
  cupertino_icons: ^1.0.2

dev_dependencies:
  flutter_test:
    sdk: flutter
  flutter_lints: ^2.0.0

flutter:
  uses-material-design: true`

	pubspecContent, err := templateEngine.ProcessTemplateWithDelims(pubspecTemplate, constants.TiloLeftDelim, constants.TiloRightDelim, ctx)
	if err != nil {
		return err
	}

	pubspecPath := filepath.Join(ctx.ProjectPath, "pubspec.yaml")
	return utils.WriteFile(pubspecPath, pubspecContent)
}

func (p *FlutterPlugin) generateConfigFiles(ctx *tilocontext.ExecutionContext) error {
	return nil
}

func (p *FlutterPlugin) Validate() error {
	_, err := p.sdkFinder.FindFlutterSDK()
	return err
}

func (p *FlutterPlugin) PreGenerate(ctx *tilocontext.ExecutionContext) error {
	ctx.SetVariable("flutter_version", "^3.16.0")
	ctx.SetVariable("dart_version", "^3.2.0")
	ctx.SetVariable("cupertino_icons_version", "^1.0.2")
	ctx.SetVariable("flutter_lints_version", "^2.0.0")

	return nil
}

func (p *FlutterPlugin) Generate(ctx *tilocontext.ExecutionContext) error {
	if err := p.createDirectoryStructure(ctx); err != nil {
		return errors.Wrap(err, "failed to create directory structure")
	}

	if err := p.generatePubspec(ctx); err != nil {
		return errors.Wrap(err, "failed to generate pubspec.yaml")
	}

	if err := p.generateSourceFiles(ctx); err != nil {
		return errors.Wrap(err, "failed to generate source files")
	}

	if err := p.generateConfigFiles(ctx); err != nil {
		return errors.Wrap(err, "failed to generate config files")
	}

	return nil
}

func (p *FlutterPlugin) PostGenerate(ctx *tilocontext.ExecutionContext) error {
	ctx.SetMetadata("framework_generated", true)
	ctx.SetMetadata("start_command", "flutter run")

	return nil
}
