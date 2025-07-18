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
	Short: "✨ TiLoKit – Khởi tạo project đa framework",
	Run: func(cmd *cobra.Command, args []string) {
		if projectName == "" {
			prompt := &survey.Input{Message: "📁 Nhập tên project:"}
			survey.AskOne(prompt, &projectName, survey.WithValidator(survey.Required))
		}

		if templateName == "" {
			options := template.GetSupportedTemplates()
			if len(options) == 0 {
				fmt.Println("⚠️ Không có template nào được cấu hình.")
				os.Exit(1)
			}

			prompt := &survey.Select{
				Message: "📦 Chọn framework muốn khởi tạo:",
				Options: options,
			}
			survey.AskOne(prompt, &templateName, survey.WithValidator(survey.Required))
		}

		if !template.Exists(templateName) {
			fmt.Println("❌ Template không hợp lệ:", templateName)
			os.Exit(1)
		}

		if err := template.Generate(templateName, projectName); err != nil {
			fmt.Println("❌ Lỗi:", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("❌ CLI lỗi:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&templateName, "template", "t", "", "Tên template (react, laravel...)")
	rootCmd.Flags().StringVarP(&projectName, "name", "n", "", "Tên project")
}
