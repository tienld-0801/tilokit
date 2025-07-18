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

func GetSupportedTemplates() []string {
	keys := []string{}
	for k := range registry {
		keys = append(keys, k)
	}
	return keys
}

func Exists(name string) bool {
	_, ok := registry[name]
	return ok
}

func Generate(name, projectName string) error {
	if gen, ok := registry[name]; ok {
		return gen(projectName)
	}
	return fmt.Errorf("template '%s' không tồn tại", name)
}
