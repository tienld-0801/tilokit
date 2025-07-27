package react

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ti-lo/tilokit/internal/utils"
	"github.com/ti-lo/tilokit/internal/constants"
	"github.com/ti-lo/tilokit/internal/common"
)

func GenerateReactJSOptions(projectName string) error {

	// Ask for package manager
	var pkgManager string
	prompt := &survey.Select{
		Message: "📦 Choose your package manager:",
		Options: constants.PackageManagers,
	}
	if err := survey.AskOne(prompt, &pkgManager, survey.WithValidator(survey.Required)); err != nil {
		return fmt.Errorf("failed to choose package manager: %w", err)
	}

	utils.Log("✨ Generating React + Vite (JavaScript) project: %s", projectName)

	// Scaffold project using create-vite
	var cmdName string
	var args []string
	switch pkgManager {
	case "npm":
		cmdName = "npm"
		args = []string{"create", "vite@latest", projectName, "--", "--template", "react"}
	case "yarn":
		cmdName = "yarn"
		args = []string{"create", "vite", projectName, "--template", "react"}
	case "pnpm":
		cmdName = "pnpm"
		args = []string{"create", "vite", projectName, "--", "--template", "react"}
	case "bun":
		cmdName = "bunx"
		args = []string{"create-vite@latest", projectName, "--", "--template", "react"}
	}

	if err := utils.RunCommand("", cmdName, args...); err != nil {
		return err
	}

	// Offer to install common libs
	libs := common.ChooseCommonLibs("react")
	if len(libs) > 0 {
		pkgArgs := utils.MapLibsToPackages(libs)
		// install libs
		switch pkgManager {
		case "npm":
			if err := utils.RunCommand(projectName, "npm", append([]string{"install"}, pkgArgs...)...); err != nil {
				return err
			}
		case "yarn":
			if err := utils.RunCommand(projectName, "yarn", append([]string{"add"}, pkgArgs...)...); err != nil {
				return err
			}
		case "pnpm":
			if err := utils.RunCommand(projectName, "pnpm", append([]string{"add"}, pkgArgs...)...); err != nil {
				return err
			}
		case "bun":
			if err := utils.RunCommand(projectName, "bun", append([]string{"add"}, pkgArgs...)...); err != nil {
				return err
			}
		}
	}

	utils.Log("🎉 React project '%s' successfully created!", projectName)
	return nil
}
