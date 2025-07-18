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

// GetSupportedTemplates trả về danh sách template
func GetSupportedTemplates() []string {
	keys := []string{}
	for k := range registry {
		keys = append(keys, k)
	}
	return keys
}

// Exists kiểm tra template có tồn tại không
func Exists(name string) bool {
	_, ok := registry[name]
	return ok
}

// Generate gọi template tương ứng
func Generate(name, projectName string) error {
	if gen, ok := registry[name]; ok {
		return gen(projectName)
	}
	return fmt.Errorf("template '%s' không tồn tại", name)
}
