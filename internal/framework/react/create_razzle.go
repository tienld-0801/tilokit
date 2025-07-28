package react

import "github.com/ti-lo/tilokit/internal/core/pkgmgr"

// CreateRazzle returns a scaffold function for React + Razzle (SSR framework).
// If useTS is true, we will pass the --typescript flag.
func CreateRazzle(useTS bool) func(string, pkgmgr.PkgMgr) error {
    return func(projectName string, pm pkgmgr.PkgMgr) error {
        args := []string{projectName}
        if useTS {
            args = append(args, "--typescript")
        }
        // Razzle uses create-razzle-app package.
        return pkgmgr.ExecDLX("", pm, "create-razzle-app", args...)
    }
}
