package react

import (
	"github.com/AlecAivazis/survey/v2"
	templateCore "github.com/ti-lo/tilokit/internal/core/template"
	reactfw "github.com/ti-lo/tilokit/internal/framework/react"
)

// GenerateReact scaffolds a React + Vite project (JS or TS) via unified Template flow
func Generate(projectName string) error {
	var variant string
	if err := survey.AskOne(&survey.Select{
		Message: "Choose language for React template:",
		Options: []string{"JavaScript", "TypeScript"},
	}, &variant); err != nil {
		return err
	}

	var tmpl templateCore.Template
	if variant == "TypeScript" {
		tmpl = templateCore.Template{
			Base:         reactfw.CreateVite("react-ts"),
			CommonLibKey: "react-vite-ts",
			Mixins:       nil,
		}
	} else {
		tmpl = templateCore.Template{
			Base:         reactfw.CreateVite("react"),
			CommonLibKey: "react",
			Mixins:       nil,
		}
	}
	return templateCore.Generate(projectName, tmpl)
}
