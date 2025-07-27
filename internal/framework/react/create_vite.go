package react

import "github.com/ti-lo/tilokit/internal/core/pkgmgr"

// CreateVite returns a scaffold function for React + Vite with provided sub template (react or react-ts)
func CreateVite(subTemplate string) func(string, pkgmgr.PkgMgr) error {
    return func(projectName string, pm pkgmgr.PkgMgr) error {
        return pkgmgr.ExecDLX("", pm, "create", "vite@latest", projectName, "--", "--template", subTemplate)
    }
}
