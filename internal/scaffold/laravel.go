package scaffold

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/ti-lo/tilokit/internal/utils"
)

// GenerateLaravelOptions scaffolds a Laravel application using the Laravel installer or composer.
func GenerateLaravelOptions(projectName string) error {
	var toolChoice string
	survey.AskOne(&survey.Select{
		Message: "Choose install method:",
		Options: []string{"laravel installer (preferred)", "composer"},
	}, &toolChoice, survey.WithValidator(survey.Required))

	utils.Log("🚧 Generating Laravel project: %s", projectName)

	switch toolChoice {
	case "laravel installer (preferred)":
		if err := utils.RunCommand("", "laravel", "new", projectName, "--jet", "--dark"); err != nil {
			return err
		}
	case "composer":
		if err := utils.RunCommand("", "composer", "create-project", "laravel/laravel", projectName); err != nil {
			return err
		}
	}

	utils.Log("🎉 Laravel project '%s' successfully created!", projectName)
	return nil
}
