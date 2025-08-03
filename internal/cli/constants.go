package cli

// Version information - set during build time
var (
	Version   = "v0.1.2-dev"
	BuildDate = "unknown"
	GitCommit = "unknown"
	GoVersion = "1.24.5"
)

// Known long flags for validation
var KnownLongFlags = []string{
	"version", "init", "name", "framework", "build-tool",
	"output", "list-frameworks", "list-build-tools",
	"quiet", "force", "update", "help",
}

// Supported Frameworks - central registry
var SupportedFrameworks = []string{
	"react", "vue", "svelte", "angular", "next", "nuxt",
	"django", "flask", "fastapi",
	"laravel", "symfony",
	"spring-boot", "quarkus",
	"gin", "echo", "fiber",
	"rails",
}

// CLI Messages
const (
	AppName        = "tilokit"
	AppShort       = "✨ TiLoKit – Modern Multi-Framework Project Generator"
	AppDescription = "Universal CLI toolkit for multi-framework project generation"

	// Error messages
	InvalidCommandMsg = "invalid command '%s'. All commands must use flags with - or -- prefix. Use --help for available options"
	InvalidFlagMsg    = "invalid flag '%s'. Use --%s for long form or find the correct short form (e.g., -v for --version, -i for --init)"
	InvalidFlagGenericMsg = "invalid flag '%s'. Single dash flags must be exactly one character. Use double dash (--) for long form flags"
)
