package vue

import "github.com/ti-lo/tilokit/internal/core/pkgmgr"

// CreateVue returns a scaffold function for Vue via create-vue.
func CreateVue(subTemplate string) func(string, pkgmgr.PkgMgr) error {
    return func(projectName string, pm pkgmgr.PkgMgr) error {
        // bunx uses create-vue@latest directly
        if pm == pkgmgr.Bun {
            return pkgmgr.ExecDLX("", pm, "create-vue@latest", projectName, "--", "--template", subTemplate)
        }
        return pkgmgr.ExecDLX("", pm, "create", "vue@latest", projectName, "--", "--template", subTemplate)
    }
}
