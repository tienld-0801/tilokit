package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/ti-lo/tilokit/internal/template"
)

var (
	templateName string
	projectName  string
)

var rootCmd = &cobra.Command{
	Use:   "tilokit",
	Short: "✨ TiLoKit – Init project multiple framework",
	Run: func(cmd *cobra.Command, args []string) {
		if projectName == "" {
			prompt := &survey.Input{Message: "📁 Input name project:"}
			survey.AskOne(prompt, &projectName, survey.WithValidator(survey.Required))
		}

		if templateName == "" {
			options := template.GetSupportedTemplates()
			if len(options) == 0 {
				fmt.Println("⚠️ No template config.")
				os.Exit(1)
			}

			prompt := &survey.Select{
				Message: "📦 Choose framework init",
				Options: options,
			}
			survey.AskOne(prompt, &templateName, survey.WithValidator(survey.Required))
		}

		if !template.Exists(templateName) {
			fmt.Println("❌ Template isvalid:", templateName)
			os.Exit(1)
		}

		if err := template.Generate(templateName, projectName); err != nil {
			fmt.Println("❌ Error:", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("❌ CLI Error:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&templateName, "template", "t", "", "Name template (react, laravel...)")
	rootCmd.Flags().StringVarP(&projectName, "name", "n", "", "Name project")
}
