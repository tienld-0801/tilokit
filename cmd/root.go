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
	listOnly     bool
	force        bool
	quiet        bool
)

var rootCmd = &cobra.Command{
	SilenceUsage: true,
	Use:          "tilokit",
	Short:        "✨ TiLoKit – Init project multiple framework",
	Run: func(cmd *cobra.Command, args []string) {
		// Init Banner
		utils.PrintBanner()

		if listOnly {
			utils.Log("Supported templates: %v", platform.GetSupportedTemplates())
			return
		}

		utils.SetQuiet(quiet)

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

		if utils.IsProduction() && (templateName == "" || projectName == "") {
			utils.Error("missing --template/-t and/or --name/-n. Run 'tilokit --help' for usage")
			os.Exit(2)
		}

		if !platform.Exists(templateName) {
			utils.Error("❌ Template invalid: %s", templateName)
			os.Exit(1)
		}

		if utils.PathExists(projectName) && !force {
			utils.Error("directory %s already exists (use --force to overwrite)", projectName)
			os.Exit(2)
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
	rootCmd.Flags().BoolVarP(&listOnly, "list", "l", false, "List available templates")
	rootCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing directory")
	rootCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Silence normal output")
}
