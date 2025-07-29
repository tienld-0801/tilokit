package react

import (
	"github.com/AlecAivazis/survey/v2"
	common "github.com/ti-lo/tilokit/internal/common"
	templateCore "github.com/ti-lo/tilokit/internal/core/template"
	reactfw "github.com/ti-lo/tilokit/internal/framework/react"
)

// Generate scaffolds a React project with selectable bundler (Vite / Parcel / Razzle) and language (JS/TS) via unified Template flow
func Generate(projectName string) error {
	var bundler string
	if err := survey.AskOne(&survey.Select{
		Message: "Choose build tool / framework:",
		Options: []string{"Vite", "Parcel", "Razzle"},
	}, &bundler); err != nil {
		return err
	}

	var variant string
	if err := survey.AskOne(&survey.Select{
		Message: "Choose language for React template:",
		Options: []string{"JavaScript", "TypeScript"},
	}, &variant); err != nil {
		return err
	}

	var tmpl templateCore.Template

	switch bundler {
	case "Vite":
		if variant == "TypeScript" {
			tmpl = templateCore.Template{
				Base:         common.CreateVite("react-ts"),
				CommonLibKey: "react-vite-ts",
			}
		} else {
			tmpl = templateCore.Template{
				Base:         common.CreateVite("react"),
				CommonLibKey: "react",
			}
		}
	case "Parcel":
		if variant == "TypeScript" {
			tmpl = templateCore.Template{
				Base:         reactfw.CreateParcel("react-ts"),
				CommonLibKey: "react",
			}
		} else {
			tmpl = templateCore.Template{
				Base:         reactfw.CreateParcel("react"),
				CommonLibKey: "react",
			}
		}
	case "Razzle":
		tmpl = templateCore.Template{
			Base:         reactfw.CreateRazzle(variant == "TypeScript"),
			CommonLibKey: "react",
		}
	}
	return templateCore.Generate(projectName, tmpl)
}
