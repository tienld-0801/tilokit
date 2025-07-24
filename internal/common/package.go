package common

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

func ChooseCommonLibs(framework string) []string {
	var libChoices []string
	var options []string

	switch framework {
	case "react", "react-vite-ts", "react-vite-tailwind":
		options = []string{
			"ESLint",
			"Prettier",
			"TailwindCSS",
			"React Router",
			"Zustand",
			"Axios",
			"Jest",
		}
	case "vue", "vue-vite-ts", "vue-vite-tailwind":
		options = []string{
			"ESLint",
			"Prettier",
			"TailwindCSS",
			"Vue Router",
			"Pinia",
			"Axios",
			"Vitest",
		}
	default:
		options = []string{
			"ESLint", "Prettier", "TailwindCSS", "Axios",
		}
	}

	prompt := &survey.MultiSelect{
		Message: "🔧 Select common libraries for your project:",
		Options: options,
	}

	err := survey.AskOne(prompt, &libChoices)
	if err != nil {
		fmt.Println("❌ Cancelled or error occurred.")
		return nil
	}

	fmt.Printf("✅ Selected libraries: %v\n", libChoices)
	return libChoices
}
