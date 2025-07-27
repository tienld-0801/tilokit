package scaffold

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/ti-lo/tilokit/internal/utils"
	"github.com/ti-lo/tilokit/internal/constants"
)

// GenerateAngular creates a new directory for an Angular project with the specified name.
// Returns an error if the directory cannot be created.
func GenerateAngularOptions(projectName string) error {
	utils.Log("🚧 Generating Angular project: %s", projectName)

	var pkgManager string
	survey.AskOne(&survey.Select{
		Message: "📦 Choose your package manager:",
		Options: constants.PackageManagers,
	}, &pkgManager, survey.WithValidator(survey.Required))

	// Use the Angular CLI via npx / yarn dlx / pnpm dlx
	var cmdName string
	var args []string
	switch pkgManager {
	case "npm":
		cmdName = "npx"
		args = []string{"-y", "@angular/cli", "new", projectName, "--routing", "--style", "scss"}
	case "yarn":
		cmdName = "yarn"
		args = []string{"dlx", "@angular/cli", "new", projectName, "--routing", "--style", "scss"}
	case "pnpm":
		cmdName = "pnpm"
		args = []string{"dlx", "@angular/cli", "new", projectName, "--routing", "--style", "scss"}
	case "bun":
		cmdName = "bunx"
		args = []string{"@angular/cli", "new", projectName, "--routing", "--style", "scss"}
	}

	if err := utils.RunCommand("", cmdName, args...); err != nil {
		return err
	}

	utils.Log("🎉 Angular project '%s' created!", projectName)
	return nil
}
