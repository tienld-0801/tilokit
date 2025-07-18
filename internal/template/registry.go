package template

import (
	"fmt"

	"github.com/ti-lo/tilokit/internal/scaffold"
)

type GeneratorFunc func(projectName string) error

var registry = map[string]GeneratorFunc{
	"react":   scaffold.GenerateReact,
	"laravel": scaffold.GenerateLaravel,
}

// GetSupportedTemplates returns a slice of all registered template names available for project generation.
func GetSupportedTemplates() []string {
	keys := []string{}
	for k := range registry {
		keys = append(keys, k)
	}
	return keys
}

// Exists reports whether a template with the given name is registered.
func Exists(name string) bool {
	_, ok := registry[name]
	return ok
}

// Generate creates a new project using the specified template and project name.
// Returns an error if the template does not exist or if project generation fails.
func Generate(name, projectName string) error {
	if gen, ok := registry[name]; ok {
		return gen(projectName)
	}
	return fmt.Errorf("template '%s' không tồn tại", name)
}
