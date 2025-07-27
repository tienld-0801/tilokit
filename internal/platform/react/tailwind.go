package react

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ti-lo/tilokit/internal/common"
	"github.com/ti-lo/tilokit/internal/utils"
	"github.com/ti-lo/tilokit/internal/constants"
)

func GenerateReactTailwind(projectName string) error {
	var pkgManager string
	prompt := &survey.Select{
		Message: "📦 Choose your package manager:",
		Options: constants.PackageManagers,
	}
	if err := survey.AskOne(prompt, &pkgManager, survey.WithValidator(survey.Required)); err != nil {
		return fmt.Errorf("failed to choose package manager: %w", err)
	}

	utils.Log("✨ Generating React + Vite + TailwindCSS project: %s", projectName)

	// First scaffold the basic TypeScript template (easiest for tailwind)
	var cmdName string
	var args []string

	switch pkgManager {
	case "npm":
		cmdName = "npm"
		args = []string{"create", "vite@latest", projectName, "--", "--template", "react-ts"}
	case "yarn":
		cmdName = "yarn"
		args = []string{"create", "vite", projectName, "--template", "react-ts"}
	case "pnpm":
		cmdName = "pnpm"
		args = []string{"create", "vite", projectName, "--", "--template", "react-ts"}
	case "bun":
		cmdName = "bunx"
		args = []string{"create-vite@latest", projectName, "--", "--template", "react-ts"}
	}

	if err := utils.RunCommand("", cmdName, args...); err != nil {
		return err
	}

	installArgs := []string{"-D", "tailwindcss", "postcss", "autoprefixer"}
	switch pkgManager {
	case "npm":
		if err := utils.RunCommand(projectName, "npm", append([]string{"install"}, installArgs...)...); err != nil {
			return err
		}
		if err := utils.RunCommand(projectName, "npx", "tailwindcss", "init", "-p"); err != nil {
			return err
		}
	case "yarn":
		if err := utils.RunCommand(projectName, "yarn", append([]string{"add"}, installArgs...)...); err != nil {
			return err
		}
		if err := utils.RunCommand(projectName, "yarn", "tailwindcss", "init", "-p"); err != nil {
			return err
		}
	case "pnpm":
		if err := utils.RunCommand(projectName, "pnpm", append([]string{"add", "-D"}, installArgs[1:]...)...); err != nil {
			return err
		}
		if err := utils.RunCommand(projectName, "pnpx", "tailwindcss", "init", "-p"); err != nil {
			return err
		}
	case "bun":
		if err := utils.RunCommand(projectName, "bun", append([]string{"add", "-d"}, installArgs...)...); err != nil {
			return err
		}
		if err := utils.RunCommand(projectName, "bunx", "tailwindcss", "init", "-p"); err != nil {
			return err
		}
	}

	// Offer to install additional libs
	libs := common.ChooseCommonLibs("react-vite-tailwind")
	if len(libs) > 0 {
		pkgs := utils.MapLibsToPackages(libs)
		switch pkgManager {
		case "npm":
			if err := utils.RunCommand(projectName, "npm", append([]string{"install"}, pkgs...)...); err != nil {
				return err
			}
		case "yarn":
			if err := utils.RunCommand(projectName, "yarn", append([]string{"add"}, pkgs...)...); err != nil {
				return err
			}
		case "pnpm":
			if err := utils.RunCommand(projectName, "pnpm", append([]string{"add"}, pkgs...)...); err != nil {
				return err
			}
		case "bun":
			if err := utils.RunCommand(projectName, "bun", append([]string{"add"}, pkgs...)...); err != nil {
				return err
			}
		}
	}

	utils.Log("🎉 React + TailwindCSS project '%s' successfully created!", projectName)
	return nil
}
