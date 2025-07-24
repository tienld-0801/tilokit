package cmd

import (
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/ti-lo/tilokit/internal/platform"
	"github.com/ti-lo/tilokit/internal/utils"
)

var (
	templateName string
	projectName  string
)

var rootCmd = &cobra.Command{
	Use:   "tilokit",
	Short: "✨ TiLoKit – Init project multiple framework",
	Run: func(cmd *cobra.Command, args []string) {
		if projectName == "" && !utils.IsProduction() {
			prompt := &survey.Input{Message: "📁 Input name project:"}
			survey.AskOne(prompt, &projectName, survey.WithValidator(survey.Required))
		}

		if templateName == "" && !utils.IsProduction() {
			options := platform.GetSupportedTemplates()
			if len(options) == 0 {
				utils.Error("⚠️ No template config.")
				os.Exit(1)
			}

			prompt := &survey.Select{
				Message: "📦 Choose framework init",
				Options: options,
			}
			survey.AskOne(prompt, &templateName, survey.WithValidator(survey.Required))
		}

		if !platform.Exists(templateName) {
			utils.Error("❌ Template invalid: %s", templateName)
			os.Exit(1)
		}

		if err := platform.Generate(templateName, projectName); err != nil {
			utils.Error("❌ Error: %v", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.Error("❌ CLI Error: %v", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&templateName, "template", "t", "", "Name template (react, laravel...)")
	rootCmd.Flags().StringVarP(&projectName, "name", "n", "", "Name project")
}
