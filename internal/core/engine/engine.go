// Package engine provides the core execution engine for TiLoKit project generation.
package engine

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	tilocontext "github.com/ti-lo/tilokit/internal/core/context"
	"github.com/ti-lo/tilokit/internal/core/registry"
)

// Engine represents the core execution engine for TiLoKit
type Engine struct {
	registry *registry.PluginRegistry
	logger   *logrus.Logger
}

// New creates a new Engine instance
func New() *Engine {
	return &Engine{
		registry: registry.New(),
		logger:   logrus.New(),
	}
}

// RegisterPlugin registers a plugin with the engine
func (e *Engine) RegisterPlugin(plugin registry.Plugin) error {
	return e.registry.Register(plugin)
}

// Execute runs the project generation process
func (e *Engine) Execute(ctx context.Context, config *tilocontext.ProjectConfig) error {
	e.logger.Info("Starting project generation...")

	// Create execution context
	execCtx := tilocontext.NewExecutionContext(config)

	// Validate configuration
	if err := e.validateConfig(config); err != nil {
		return errors.Wrap(err, "configuration validation failed")
	}

	// Load required plugins
	plugins, err := e.registry.LoadPlugins(config.Framework, config.BuildTool)
	if err != nil {
		return errors.Wrap(err, "failed to load plugins")
	}

	// Execute lifecycle hooks
	for _, plugin := range plugins {
		if err := plugin.PreGenerate(execCtx); err != nil {
			return errors.Wrapf(err, "pre-generate hook failed for plugin %s", plugin.Name())
		}
	}

	// Generate project structure
	if err := e.generateProject(execCtx, plugins); err != nil {
		return errors.Wrap(err, "project generation failed")
	}

	// Execute post-generation hooks
	for _, plugin := range plugins {
		if err := plugin.PostGenerate(execCtx); err != nil {
			return errors.Wrapf(err, "post-generate hook failed for plugin %s", plugin.Name())
		}
	}

	e.logger.Info("Project generation completed successfully!")
	return nil
}

func (e *Engine) validateConfig(config *tilocontext.ProjectConfig) error {
	if config.ProjectName == "" {
		return fmt.Errorf("project name is required")
	}
	if config.Framework == "" {
		return fmt.Errorf("framework is required")
	}
	return nil
}

func (e *Engine) generateProject(ctx *tilocontext.ExecutionContext, plugins []registry.Plugin) error {
	for _, plugin := range plugins {
		if err := plugin.Generate(ctx); err != nil {
			return errors.Wrapf(err, "generation failed for plugin %s", plugin.Name())
		}
	}
	return nil
}
