package scaffold

import (
    "github.com/AlecAivazis/survey/v2"
    "github.com/ti-lo/tilokit/internal/common"
    "github.com/ti-lo/tilokit/internal/utils"
    "github.com/ti-lo/tilokit/internal/constants"
)

// GenerateNextOptions scaffolds a Next.js application using create-next-app
func GenerateNextOptions(projectName string) error {
    var pkgManager string
    survey.AskOne(&survey.Select{
        Message: "📦 Choose your package manager:",
        Options: constants.PackageManagers,
    }, &pkgManager, survey.WithValidator(survey.Required))

    utils.Log("🚧 Generating Next.js project: %s", projectName)

    var cmdName string
    var args []string
    switch pkgManager {
    case "npm":
        cmdName = "npm"
        args = []string{"create", "next-app@latest", projectName}
    case "yarn":
        cmdName = "yarn"
        args = []string{"create", "next-app", projectName}
    case "pnpm":
        cmdName = "pnpm"
        args = []string{"create", "next-app", projectName}
    case "bun":
        cmdName = "bunx"
        args = []string{"create-next-app@latest", projectName}
    }

    if err := utils.RunCommand("", cmdName, args...); err != nil {
        return err
    }

    libs := common.ChooseCommonLibs("react") // share libs with react
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

    utils.Log("🎉 Next.js project '%s' successfully created!", projectName)
    return nil
}
