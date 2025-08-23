package constants

import "runtime"

// Version information - set during build time
var (
	Version   = "v0.2.5-dev"
	BuildDate = "unknown"
	GitCommit = "unknown"
	GoVersion = runtime.Version()
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
	"nestjs", "express", "fastify",
	"django", "flask", "fastapi",
	"laravel", "symfony", "cakephp", "codeigniter",
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
	InvalidCommandMsg     = "invalid command '%s'. All commands must use flags with - or -- prefix. Use --help for available options"
	InvalidFlagMsg        = "invalid flag '%s'. Use --%s for long form or find the correct short form (e.g., -v for --version, -i for --init)"
	InvalidFlagGenericMsg = "invalid flag '%s'. Single dash flags must be exactly one character. Use double dash (--) for long form flags"
)

// CLI FRAMEWORK
const (
	ReactFramework      = "react"
	VueFramework        = "vue"
	SvelteFramework     = "svelte"
	AngularFramework    = "angular"
	NextFramework       = "next"
	NuxtFramework       = "nuxt"
	NestFramework       = "nestjs"
	ExpressFramework    = "express"
	FastifyFramework    = "fastify"
	DjangoFramework     = "django"
	FlaskFramework      = "flask"
	FastapiFramework    = "fastapi"
	LaravelFramework     = "laravel"
	SymfonyFramework     = "symfony"
	CakePHPFramework     = "cakephp"
	CodeIgniterFramework = "codeigniter"
	SpringBootFramework  = "spring-boot"
	QuarkusFramework    = "quarkus"
	GinFramework        = "gin"
	EchoFramework       = "echo"
	FiberFramework      = "fiber"
	RailsFramework      = "rails"
)

// CLI NextJS
const (
	AppRouter   = "app"
	PagesRouter = "pages"
)

// CLI Angular
const (
	AngularCsrMode                = "csr"
	AngularSsrMode                = "ssr"
	AngularArchitectureStandalone = "standalone"
	AngularArchitectureModule     = "module"
)

// Template Delimiters
const (
	TiloLeftDelim  = "<<TILO:"
	TiloRightDelim = ">>"
)

// CLI OS
const (
	VERSION = "1.0.0"
	ARM64   = "arm64"
	AMD64   = "amd64"
)
