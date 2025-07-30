package tilocontext

import (
	"os"
	"path/filepath"
	"time"
)

// ProjectConfig holds the configuration for project generation
type ProjectConfig struct {
	ProjectName   string            `yaml:"project_name" mapstructure:"project_name"`
	Framework     string            `yaml:"framework" mapstructure:"framework"`
	BuildTool     string            `yaml:"build_tool" mapstructure:"build_tool"`
	PackageManager string           `yaml:"package_manager" mapstructure:"package_manager"`
	OutputDir     string            `yaml:"output_dir" mapstructure:"output_dir"`
	Template      string            `yaml:"template" mapstructure:"template"`
	Features      []string          `yaml:"features" mapstructure:"features"`
	Variables     map[string]interface{} `yaml:"variables" mapstructure:"variables"`
	GitInit       bool              `yaml:"git_init" mapstructure:"git_init"`
	InstallDeps   bool              `yaml:"install_deps" mapstructure:"install_deps"`
}

// ExecutionContext provides runtime context for plugin execution
type ExecutionContext struct {
	Config        *ProjectConfig
	ProjectPath   string
	TempDir       string
	StartTime     time.Time
	Variables     map[string]interface{}
	Metadata      map[string]interface{}
}

// NewExecutionContext creates a new execution context
func NewExecutionContext(config *ProjectConfig) *ExecutionContext {
	projectPath := filepath.Join(config.OutputDir, config.ProjectName)
	
	ctx := &ExecutionContext{
		Config:      config,
		ProjectPath: projectPath,
		StartTime:   time.Now(),
		Variables:   make(map[string]interface{}),
		Metadata:    make(map[string]interface{}),
	}

	// Set default variables
	ctx.Variables["project_name"] = config.ProjectName
	ctx.Variables["framework"] = config.Framework
	ctx.Variables["build_tool"] = config.BuildTool
	ctx.Variables["package_manager"] = config.PackageManager
	ctx.Variables["timestamp"] = ctx.StartTime.Format("2006-01-02 15:04:05")
	
	// Merge user variables
	for k, v := range config.Variables {
		ctx.Variables[k] = v
	}

	return ctx
}

// SetVariable sets a variable in the execution context
func (ctx *ExecutionContext) SetVariable(key string, value interface{}) {
	ctx.Variables[key] = value
}

// GetVariable gets a variable from the execution context
func (ctx *ExecutionContext) GetVariable(key string) (interface{}, bool) {
	value, exists := ctx.Variables[key]
	return value, exists
}

// SetMetadata sets metadata in the execution context
func (ctx *ExecutionContext) SetMetadata(key string, value interface{}) {
	ctx.Metadata[key] = value
}

// GetMetadata gets metadata from the execution context
func (ctx *ExecutionContext) GetMetadata(key string) (interface{}, bool) {
	value, exists := ctx.Metadata[key]
	return value, exists
}

// EnsureProjectDir creates the project directory if it doesn't exist
func (ctx *ExecutionContext) EnsureProjectDir() error {
	return os.MkdirAll(ctx.ProjectPath, 0755)
}

// CreateTempDir creates a temporary directory for processing
func (ctx *ExecutionContext) CreateTempDir() error {
	tempDir, err := os.MkdirTemp("", "tilokit-*")
	if err != nil {
		return err
	}
	ctx.TempDir = tempDir
	return nil
}

// Cleanup removes temporary resources
func (ctx *ExecutionContext) Cleanup() error {
	if ctx.TempDir != "" {
		return os.RemoveAll(ctx.TempDir)
	}
	return nil
}
