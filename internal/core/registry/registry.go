package registry

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/ti-lo/tilokit/internal/core/context"
)

// Plugin interface defines the contract for all plugins
type Plugin interface {
	Name() string
	Version() string
	Description() string
	SupportedFrameworks() []string
	SupportedBuildTools() []string
	PreGenerate(ctx *tilocontext.ExecutionContext) error
	Generate(ctx *tilocontext.ExecutionContext) error
	PostGenerate(ctx *tilocontext.ExecutionContext) error
}

// PluginRegistry manages plugin registration and loading
type PluginRegistry struct {
	plugins map[string]Plugin
	mutex   sync.RWMutex
}

// New creates a new plugin registry
func New() *PluginRegistry {
	return &PluginRegistry{
		plugins: make(map[string]Plugin),
	}
}

// Register registers a plugin with the registry
func (r *PluginRegistry) Register(plugin Plugin) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	name := plugin.Name()
	if _, exists := r.plugins[name]; exists {
		return fmt.Errorf("plugin %s is already registered", name)
	}

	r.plugins[name] = plugin
	return nil
}

// GetPlugin retrieves a plugin by name
func (r *PluginRegistry) GetPlugin(name string) (Plugin, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	plugin, exists := r.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	return plugin, nil
}

// LoadPlugins loads plugins based on framework and build tool
func (r *PluginRegistry) LoadPlugins(framework, buildTool string) ([]Plugin, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var selectedPlugins []Plugin

	for _, plugin := range r.plugins {
		// Check if plugin supports the framework
		if framework != "" && !r.supportsFramework(plugin, framework) {
			continue
		}

		// Check if plugin supports the build tool
		if buildTool != "" && !r.supportsBuildTool(plugin, buildTool) {
			continue
		}

		selectedPlugins = append(selectedPlugins, plugin)
	}

	if len(selectedPlugins) == 0 {
		return nil, errors.Errorf("no plugins found for framework: %s, build tool: %s", framework, buildTool)
	}

	return selectedPlugins, nil
}

// ListPlugins returns all registered plugins
func (r *PluginRegistry) ListPlugins() map[string]Plugin {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	plugins := make(map[string]Plugin)
	for name, plugin := range r.plugins {
		plugins[name] = plugin
	}

	return plugins
}

// GetSupportedFrameworks returns all supported frameworks
func (r *PluginRegistry) GetSupportedFrameworks() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	frameworkSet := make(map[string]bool)
	for _, plugin := range r.plugins {
		for _, framework := range plugin.SupportedFrameworks() {
			frameworkSet[framework] = true
		}
	}

	var frameworks []string
	for framework := range frameworkSet {
		frameworks = append(frameworks, framework)
	}

	return frameworks
}

// GetSupportedBuildTools returns all supported build tools
func (r *PluginRegistry) GetSupportedBuildTools() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	buildToolSet := make(map[string]bool)
	for _, plugin := range r.plugins {
		for _, buildTool := range plugin.SupportedBuildTools() {
			buildToolSet[buildTool] = true
		}
	}

	var buildTools []string
	for buildTool := range buildToolSet {
		buildTools = append(buildTools, buildTool)
	}

	return buildTools
}

func (r *PluginRegistry) supportsFramework(plugin Plugin, framework string) bool {
	for _, supported := range plugin.SupportedFrameworks() {
		if supported == framework {
			return true
		}
	}
	return false
}

func (r *PluginRegistry) supportsBuildTool(plugin Plugin, buildTool string) bool {
	for _, supported := range plugin.SupportedBuildTools() {
		if supported == buildTool {
			return true
		}
	}
	return false
}
