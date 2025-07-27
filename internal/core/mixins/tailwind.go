package mixins

import "github.com/ti-lo/tilokit/internal/core/pkgmgr"

type tailwindMixin struct{}

// Tailwind reusable instance
var Tailwind = tailwindMixin{}

func (tailwindMixin) Name() string { return "TailwindCSS" }

func (tailwindMixin) Apply(projectDir string, pm pkgmgr.PkgMgr) error {
    if err := pkgmgr.InstallDev(projectDir, pm, "tailwindcss", "postcss", "autoprefixer"); err != nil {
        return err
    }
    return pkgmgr.ExecDLX(projectDir, pm, "tailwindcss", "init", "-p")
}
