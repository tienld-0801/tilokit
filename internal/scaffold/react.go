package scaffold

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	react "github.com/ti-lo/tilokit/internal/platform/react"
)

func GenerateReactOptions(projectName string) error {
	var templateChoice string
	prompt := &survey.Select{
		Message: "🧩 Choose React template:",
		Options: []string{
			"React + Vite (JavaScript)",
			"React + Vite (TypeScript)",
			"React + Vite + TailwindCSS",
		},
	}
	err := survey.AskOne(prompt, &templateChoice, survey.WithValidator(survey.Required))
	if err != nil {
		return fmt.Errorf("failed to read React template selection: %w", err)
	}
	switch templateChoice {
	case "React + Vite (JavaScript)":
		return react.GenerateReactJSOptions(projectName)
	case "React + Vite (TypeScript)":
		return react.GenerateReactTSOptions(projectName)
	case "React + Vite + TailwindCSS":
		return react.GenerateReactTailwind(projectName)
	default:
		return fmt.Errorf("invalid template selection")
	}
}
