package platform

import (
	"fmt"

	"github.com/ti-lo/tilokit/internal/scaffold"
)

// GeneratorFunc defines the function signature for template generators.
type GeneratorFunc func(projectName string) error

// orderedTemplates defines the order in which templates should be listed in CLI.
var orderedTemplates = []string{
	"react",
	"vue",
	"angular",
	"next",
	"nuxt",
	"nest",
	"laravel",
}

// registry maps template names to their generator functions.
var registry = map[string]GeneratorFunc{
	"react":   scaffold.GenerateReactOptions,
	"vue":     scaffold.GenerateVueOptions,
	"angular": scaffold.GenerateAngularOptions,
	"next":    scaffold.GenerateNextOptions,
	"nuxt":    scaffold.GenerateNuxtOptions,
	"nest":    scaffold.GenerateNestOptions,
	"laravel": scaffold.GenerateLaravelOptions,
}

// GetSupportedTemplates returns template names in the order defined in orderedTemplates.
func GetSupportedTemplates() []string {
	templates := []string{}
	for _, name := range orderedTemplates {
		if _, ok := registry[name]; ok {
			templates = append(templates, name)
		}
	}
	return templates
}

// Exists reports whether a template with the given name is registered.
func Exists(name string) bool {
	_, ok := registry[name]
	return ok
}

// Generate creates a new project using the specified template and project name.
func Generate(name, projectName string) error {
	if gen, ok := registry[name]; ok {
		return gen(projectName)
	}
	return fmt.Errorf("template '%s' does not exist", name)
}
