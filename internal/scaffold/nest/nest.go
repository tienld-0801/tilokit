package nest

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/ti-lo/tilokit/internal/utils"
	"github.com/ti-lo/tilokit/internal/constants"
)

// GenerateNestOptions scaffolds a NestJS project using nest CLI
func Generate(projectName string) error {
	var pkgManager string
	survey.AskOne(&survey.Select{
		Message: "📦 Choose your package manager:",
		Options: constants.PackageManagers,
	}, &pkgManager, survey.WithValidator(survey.Required))

	utils.Log("🚧 Generating NestJS project: %s", projectName)

	var cmdName string
	var args []string
	// nest new via npx etc.
	switch pkgManager {
	case "npm":
		cmdName = "npx"
		args = []string{"-y", "@nestjs/cli", "new", projectName}
	case "yarn":
		cmdName = "yarn"
		args = []string{"dlx", "@nestjs/cli", "new", projectName}
	case "pnpm":
		cmdName = "pnpm"
		args = []string{"dlx", "@nestjs/cli", "new", projectName}
	case "bun":
		cmdName = "bunx"
		args = []string{"@nestjs/cli", "new", projectName}
	}

	if err := utils.RunCommand("", cmdName, args...); err != nil {
		return err
	}

	utils.Log("🎉 NestJS project '%s' successfully created!", projectName)
	return nil
}
