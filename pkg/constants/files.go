package constants

// Common file names used across the application
const (
	// Configuration files
	PackageJsonFileName  = "package.json"
	TsConfigFileName     = "tsconfig.json"
	EslintConfigFileName = "eslint.config.js"
	EnvFileName          = ".env"
	GitignoreFileName    = ".gitignore"
	ReadmeFileName       = "README.md"
	DockerFileName       = "Dockerfile"
	MakeFileName         = "Makefile"

	// Node.js specific
	NestCliFileName = "nest-cli.json"

	// Go specific
	GoModFileName = "go.mod"
	GoSumFileName = "go.sum"

	// Python specific
	RequirementsFileName = "requirements.txt"
	SetupPyFileName      = "setup.py"

	// Ruby specific
	GemFileName = "Gemfile"

	// PHP specific
	ComposerFileName = "composer.json"

	// Java/Kotlin specific
	PomXmlFileName         = "pom.xml"
	BuildGradleFileName    = "build.gradle"
	SettingsGradleFileName = "settings.gradle"

	// C# specific - Note: These are patterns, not literal filenames
	CsprojFilePattern = "*.csproj"
	SlnFilePattern    = "*.sln"

	// Expo app config
	AppConfigFileName = "app.json"
	ExpoEnvFileName   = "expo-env.d.ts"
)
