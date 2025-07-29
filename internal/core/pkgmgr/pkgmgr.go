package pkgmgr

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ti-lo/tilokit/internal/constants"
	"github.com/ti-lo/tilokit/internal/utils"
)

// PkgMgr represents a JavaScript package manager selected by the user.
type PkgMgr string

const (
	Npm  PkgMgr = "npm"
	Yarn PkgMgr = "yarn"
	Pnpm PkgMgr = "pnpm"
	Bun  PkgMgr = "bun"
)

// AskPackageManager prompts the user to choose a package manager.
func AskPackageManager() (PkgMgr, error) {
	var pm string
	if err := survey.AskOne(&survey.Select{
		Message: "📦 Choose your package manager:",
		Options: constants.PackageManagers,
	}, &pm, survey.WithValidator(survey.Required)); err != nil {
		return "", err
	}
	return PkgMgr(pm), nil
}

// internal helpers for DLX-style exec.
func runner(pm PkgMgr) (string, []string) {
	switch pm {
	case Npm:
		return "npm", nil
	case Yarn:
		return "yarn", []string{"dlx"}
	case Pnpm:
		return "pnpm", []string{"dlx"}
	case Bun:
		return "bunx", nil
	default:
		return string(pm), nil
	}
}

// ExecDLX runs a CLI package through the package manager.
func ExecDLX(dir string, pm PkgMgr, pkg string, extraArgs ...string) error {
	name, pre := runner(pm)
	args := append(append(pre, pkg), extraArgs...)
	return utils.RunCommand(dir, name, args...)
}

// Install installs dependencies (prod).
func Install(projectDir string, pm PkgMgr, pkgs ...string) error {
	var name string
	var args []string
	switch pm {
	case Npm:
		name = "npm"
		args = append([]string{"install"}, pkgs...)
	case Yarn:
		name = "yarn"
		args = append([]string{"add"}, pkgs...)
	case Pnpm:
		name = "pnpm"
		args = append([]string{"add"}, pkgs...)
	case Bun:
		name = "bun"
		args = append([]string{"add"}, pkgs...)
	default:
		return fmt.Errorf("unknown package manager %s", pm)
	}
	return utils.RunCommand(projectDir, name, args...)
}

// InstallDev installs dev dependencies.
func InstallDev(projectDir string, pm PkgMgr, pkgs ...string) error {
	var name string
	var args []string
	switch pm {
	case Npm:
		name = "npm"
		args = append([]string{"install", "-D"}, pkgs...)
	case Yarn:
		name = "yarn"
		args = append([]string{"add", "-D"}, pkgs...)
	case Pnpm:
		name = "pnpm"
		args = append([]string{"add", "-D"}, pkgs...)
	case Bun:
		name = "bun"
		args = append([]string{"add", "-d"}, pkgs...)
	default:
		return fmt.Errorf("unknown package manager %s", pm)
	}
	return utils.RunCommand(projectDir, name, args...)
}
