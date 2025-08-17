package utils

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"strings"
)

// EnvConfig holds all environment variables used by TiLoKit
type EnvConfig struct {
	// Template variables
	WelcomeMessage string
	ProjectAuthor  string
	ProjectEmail   string

	// Build settings
	DefaultNodeVersion string
	DefaultGoVersion   string

	// Feature flags
	EnableTypeScript bool
	EnableTesting    bool
	EnableLinting    bool

	// Development settings
	DevMode     bool
	VerboseMode bool
}

// LoadEnvConfig loads all environment variables into a structured config
func LoadEnvConfig() *EnvConfig {
	_ = godotenv.Load()

	return &EnvConfig{
		// Template variables
		WelcomeMessage: getEnvWithDefault("TILOKIT_WELCOME_MESSAGE", "Welcome to TiLoKit!"),
		ProjectAuthor:  getEnvWithDefault("TILOKIT_PROJECT_AUTHOR", ""),
		ProjectEmail:   getEnvWithDefault("TILOKIT_PROJECT_EMAIL", ""),

		// Build settings
		DefaultNodeVersion: getEnvWithDefault("TILOKIT_NODE_VERSION", "20"),
		DefaultGoVersion:   getEnvWithDefault("TILOKIT_GO_VERSION", "1.21"),

		// Feature flags
		EnableTypeScript: getEnvBool("TILOKIT_ENABLE_TYPESCRIPT", true),
		EnableTesting:    getEnvBool("TILOKIT_ENABLE_TESTING", true),
		EnableLinting:    getEnvBool("TILOKIT_ENABLE_LINTING", true),

		// Development settings
		DevMode:     getEnvBool("TILOKIT_DEV_MODE", false),
		VerboseMode: getEnvBool("TILOKIT_VERBOSE", false),
	}
}

// ToVariables converts EnvConfig to a map for template variables
func (e *EnvConfig) ToVariables() map[string]interface{} {
	return map[string]interface{}{
		"welcome_message":      e.WelcomeMessage,
		"project_author":       e.ProjectAuthor,
		"project_email":        e.ProjectEmail,
		"default_node_version": e.DefaultNodeVersion,
		"default_go_version":   e.DefaultGoVersion,
		"enable_typescript":    e.EnableTypeScript,
		"enable_testing":       e.EnableTesting,
		"enable_linting":       e.EnableLinting,
		"dev_mode":            e.DevMode,
		"verbose_mode":        e.VerboseMode,
	}
}

// Helper functions
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		// Parse boolean from string
		switch strings.ToLower(value) {
		case "true", "1", "yes", "on":
			return true
		case "false", "0", "no", "off":
			return false
		default:
			// Try to parse as boolean
			if parsed, err := strconv.ParseBool(value); err == nil {
				return parsed
			}
		}
	}
	return defaultValue
}
