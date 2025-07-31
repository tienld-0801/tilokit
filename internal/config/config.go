package config

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/ti-lo/tilokit/internal/core/context"
)

// Config holds the application configuration
type Config struct {
	DefaultFramework     string            `mapstructure:"default_framework"`
	DefaultBuildTool     string            `mapstructure:"default_build_tool"`
	DefaultPackageManager string           `mapstructure:"default_package_manager"`
	DefaultOutputDir     string            `mapstructure:"default_output_dir"`
	Templates            map[string]string `mapstructure:"templates"`
	Plugins              []string          `mapstructure:"plugins"`
	Features             map[string]bool   `mapstructure:"features"`
}

// LoadConfig loads configuration from file and environment
func LoadConfig() (*Config, error) {
	viper.SetConfigName("tilokit")
	viper.SetConfigType("yaml")
	
	// Add config paths
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.tilokit")
	viper.AddConfigPath("/etc/tilokit")

	// Set defaults
	setDefaults()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, errors.Wrap(err, "failed to read config file")
		}
	}

	// Bind environment variables
	viper.SetEnvPrefix("TILOKIT")
	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}

	return &config, nil
}

// CreateProjectConfig creates a project configuration from CLI inputs
func CreateProjectConfig(projectName, framework, buildTool, outputDir string) *tilocontext.ProjectConfig {
	config := &tilocontext.ProjectConfig{
		ProjectName:    projectName,
		Framework:      framework,
		BuildTool:      buildTool,
		OutputDir:      outputDir,
		Variables:      make(map[string]interface{}),
		GitInit:        true,
		InstallDeps:    true,
	}

	// Set defaults if not provided
	if config.OutputDir == "" {
		config.OutputDir = "."
	}

	if config.BuildTool == "" {
		config.BuildTool = getDefaultBuildTool(framework)
	}

	// Set package manager based on build tool
	config.PackageManager = getDefaultPackageManager(config.BuildTool)

	return config
}

// SaveConfig saves the current configuration to file
func SaveConfig(config *Config) error {
	configDir := filepath.Join(os.Getenv("HOME"), ".tilokit")
	// Use more restrictive permissions (0750 instead of 0755)
	if err := os.MkdirAll(configDir, 0750); err != nil {
		return errors.Wrap(err, "failed to create config directory")
	}

	configFile := filepath.Join(configDir, "tilokit.yaml")
	
	viper.Set("default_framework", config.DefaultFramework)
	viper.Set("default_build_tool", config.DefaultBuildTool)
	viper.Set("default_package_manager", config.DefaultPackageManager)
	viper.Set("default_output_dir", config.DefaultOutputDir)
	viper.Set("templates", config.Templates)
	viper.Set("plugins", config.Plugins)
	viper.Set("features", config.Features)

	return viper.WriteConfigAs(configFile)
}

func setDefaults() {
	viper.SetDefault("default_framework", "react")
	viper.SetDefault("default_build_tool", "vite")
	viper.SetDefault("default_package_manager", "npm")
	viper.SetDefault("default_output_dir", ".")
	viper.SetDefault("templates", map[string]string{})
	viper.SetDefault("plugins", []string{})
	viper.SetDefault("features", map[string]bool{
		"typescript":     true,
		"eslint":         true,
		"prettier":       true,
		"testing":        true,
		"git_init":       true,
		"install_deps":   true,
	})
}

func getDefaultBuildTool(framework string) string {
	buildTools := map[string]string{
		"react":   "vite",
		"vue":     "vite",
		"svelte":  "vite",
		"angular": "angular-cli",
		"next":    "next",
		"nuxt":    "nuxt",
		"gatsby":  "gatsby",
	}

	if tool, exists := buildTools[framework]; exists {
		return tool
	}
	return "vite"
}

func getDefaultPackageManager(buildTool string) string {
	packageManagers := map[string]string{
		"vite":        "npm",
		"webpack":     "npm",
		"rollup":      "npm",
		"angular-cli": "npm",
		"next":        "npm",
		"nuxt":        "npm",
		"gatsby":      "npm",
	}

	if pm, exists := packageManagers[buildTool]; exists {
		return pm
	}
	return "npm"
}
