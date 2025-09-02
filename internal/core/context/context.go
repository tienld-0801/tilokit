package tilocontext

import (
	"os"
	"path/filepath"
	"tilokit/internal/utils"
	"time"
)

type ProjectConfig struct {
	ProjectName    string                 `yaml:"project_name" mapstructure:"project_name"`
	Framework      string                 `yaml:"framework" mapstructure:"framework"`
	BuildTool      string                 `yaml:"build_tool" mapstructure:"build_tool"`
	PackageManager string                 `yaml:"package_manager" mapstructure:"package_manager"`
	OutputDir      string                 `yaml:"output_dir" mapstructure:"output_dir"`
	Template       string                 `yaml:"template" mapstructure:"template"`
	Features       []string               `yaml:"features" mapstructure:"features"`
	Variables      map[string]interface{} `yaml:"variables" mapstructure:"variables"`
	GitInit        bool                   `yaml:"git_init" mapstructure:"git_init"`
}

type ExecutionContext struct {
	Config      *ProjectConfig
	ProjectPath string
	TempDir     string
	StartTime   time.Time
	Variables   map[string]interface{}
	Metadata    map[string]interface{}
}

func NewExecutionContext(config *ProjectConfig) *ExecutionContext {
	projectPath := filepath.Join(config.OutputDir, config.ProjectName)

	ctx := &ExecutionContext{
		Config:      config,
		ProjectPath: projectPath,
		StartTime:   time.Now(),
		Variables:   make(map[string]interface{}),
		Metadata:    make(map[string]interface{}),
	}

	ctx.Variables["project_name"] = config.ProjectName
	ctx.Variables["framework"] = config.Framework
	ctx.Variables["build_tool"] = config.BuildTool
	ctx.Variables["package_manager"] = config.PackageManager
	ctx.Variables["timestamp"] = ctx.StartTime.Format("2006-01-02 15:04:05")
	ctx.Variables["welcome_message"] = config.ProjectName

	envConfig := utils.LoadEnvConfig()
	envVars := envConfig.ToVariables()

	for k, v := range envVars {
		ctx.Variables[k] = v
	}

	for k, v := range config.Variables {
		ctx.Variables[k] = v
	}

	return ctx
}

func (ctx *ExecutionContext) SetVariable(key string, value interface{}) {
	ctx.Variables[key] = value
}

func (ctx *ExecutionContext) GetVariable(key string) (interface{}, bool) {
	value, exists := ctx.Variables[key]
	return value, exists
}

func (ctx *ExecutionContext) SetMetadata(key string, value interface{}) {
	ctx.Metadata[key] = value
}

func (ctx *ExecutionContext) GetMetadata(key string) (interface{}, bool) {
	value, exists := ctx.Metadata[key]
	return value, exists
}

func (ctx *ExecutionContext) EnsureProjectDir() error {
	if ctx.ProjectPath == "" {
		return os.MkdirAll(ctx.ProjectPath, 0750)
	}
	return nil
}

func (ctx *ExecutionContext) CreateTempDir() error {
	tempDir, err := os.MkdirTemp("", "tilokit-*")
	if err != nil {
		return err
	}
	ctx.TempDir = tempDir
	return nil
}

func (ctx *ExecutionContext) Cleanup() error {
	if ctx.TempDir != "" {
		return os.RemoveAll(ctx.TempDir)
	}
	return nil
}
