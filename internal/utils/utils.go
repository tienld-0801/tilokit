package utils

import (
	"fmt"

	"github.com/ti-lo/tilokit/internal/scaffold"
)

var Generators = map[string]func(string) error{
	"react":   scaffold.GenerateReact,
	"laravel": scaffold.GenerateLaravel,
}

func IsValidTemplate(name string) bool {
	_, ok := Generators[name]
	return ok
}

func Generate(template, projectName string) error {
	if gen, ok := Generators[template]; ok {
		return gen(projectName)
	}
	return fmt.Errorf("template %s not found", template)
}
