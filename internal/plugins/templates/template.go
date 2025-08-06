package templates

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"

	"github.com/ti-lo/tilokit/internal/core/context"
	"github.com/ti-lo/tilokit/internal/utils"
)

// TemplateEngine handles template processing
type TemplateEngine struct {
	templates map[string]*template.Template
}

// NewTemplateEngine creates a new template engine
func NewTemplateEngine() *TemplateEngine {
	return &TemplateEngine{
		templates: make(map[string]*template.Template),
	}
}

// ProcessTemplate processes a template string with context variables
func (te *TemplateEngine) ProcessTemplate(templateContent string, ctx *tilocontext.ExecutionContext) (string, error) {
	tmpl, err := template.New("template").Parse(templateContent)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse template")
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, ctx.Variables); err != nil {
		return "", errors.Wrap(err, "failed to execute template")
	}

	return result.String(), nil
}

// ProcessTemplateFile processes a template file and writes the result
func (te *TemplateEngine) ProcessTemplateFile(templatePath, outputPath string, ctx *tilocontext.ExecutionContext) error {
	// Read template file
	templateContent, err := utils.ReadFile(templatePath)
	if err != nil {
		return errors.Wrap(err, "failed to read template file")
	}

	// Process template
	result, err := te.ProcessTemplate(templateContent, ctx)
	if err != nil {
		return errors.Wrap(err, "failed to process template")
	}

	// Write result
	if err := utils.WriteFile(outputPath, result); err != nil {
		return errors.Wrap(err, "failed to write processed template")
	}

	return nil
}

// CopyTemplateDirectory copies and processes all templates in a directory
func (te *TemplateEngine) CopyTemplateDirectory(templateDir, outputDir string, ctx *tilocontext.ExecutionContext) error {
	return filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path
		relPath, err := filepath.Rel(templateDir, path)
		if err != nil {
			return err
		}

		outputPath := filepath.Join(outputDir, relPath)

		if info.IsDir() {
			return utils.EnsureDir(outputPath)
		}

		// Process template files
		if strings.HasSuffix(path, ".tmpl") {
			// Remove .tmpl extension from output
			outputPath = strings.TrimSuffix(outputPath, ".tmpl")
			return te.ProcessTemplateFile(path, outputPath, ctx)
		}

		// Copy non-template files as-is
		return utils.CopyFile(path, outputPath)
	})
}
