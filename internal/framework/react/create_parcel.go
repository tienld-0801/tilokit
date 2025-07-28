package react

import "github.com/ti-lo/tilokit/internal/core/pkgmgr"

// CreateParcel returns a scaffold function for React + Parcel with provided sub template (react or react-ts).
// It invokes the official Parcel project scaffolder via `create-parcel@latest`.
func CreateParcel(subTemplate string) func(string, pkgmgr.PkgMgr) error {
    return func(projectName string, pm pkgmgr.PkgMgr) error {
        // Parcel templates support "react" and "react-ts" among others.
        // We rely on package manager's DLX/exec to download and run the scaffolder.
        return pkgmgr.ExecDLX("", pm, "create-parcel@latest", projectName, "--", "--template", subTemplate)
    }
}
