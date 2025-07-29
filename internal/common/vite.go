package common

import "github.com/ti-lo/tilokit/internal/core/pkgmgr"

// CreateVite returns a scaffold function for any framework
func CreateVite(subTemplate string) func(string, pkgmgr.PkgMgr) error {
	return func(projectName string, pm pkgmgr.PkgMgr) error {
		args := []string{"create", "vite@latest", projectName, "--", "--template", subTemplate}
		return pkgmgr.ExecDLX("", pm, args[0], args[1:]...)
	}
}
