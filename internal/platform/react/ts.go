package react

import (
    "fmt"

    "github.com/AlecAivazis/survey/v2"
    "github.com/ti-lo/tilokit/internal/common"
    "github.com/ti-lo/tilokit/internal/utils"
    "github.com/ti-lo/tilokit/internal/constants"
)

func GenerateReactTSOptions(projectName string) error {
    var pkgManager string
    prompt := &survey.Select{
        Message: "📦 Choose your package manager:",
        Options: constants.PackageManagers,
    }
    if err := survey.AskOne(prompt, &pkgManager, survey.WithValidator(survey.Required)); err != nil {
        return fmt.Errorf("failed to choose package manager: %w", err)
    }

    utils.Log("✨ Generating React + Vite (TypeScript) project: %s", projectName)

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

    libs := common.ChooseCommonLibs("react-vite-ts")
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

    utils.Log("🎉 React (TS) project '%s' successfully created!", projectName)
    return nil
}
