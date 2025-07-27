package template

import (
    "fmt"

    "github.com/ti-lo/tilokit/internal/common"
    "github.com/ti-lo/tilokit/internal/core/pkgmgr"
    "github.com/ti-lo/tilokit/internal/utils"
)

// Mixin represents an optional feature that can be applied after scaffold.
type Mixin interface {
    Name() string
    Apply(projectDir string, pm pkgmgr.PkgMgr) error
}

// Template describes scaffold recipe.
type Template struct {
    Base         func(projectName string, pm pkgmgr.PkgMgr) error
    CommonLibKey string
    Mixins       []Mixin
}

// Generate executes the template flow in a unified way.
func Generate(projectName string, tmpl Template) error {
    pm, err := pkgmgr.AskPackageManager()
    if err != nil {
        return err
    }

    if err := tmpl.Base(projectName, pm); err != nil {
        return err
    }

    // common libs
    if tmpl.CommonLibKey != "" {
        libs := common.ChooseCommonLibs(tmpl.CommonLibKey)
        if len(libs) > 0 {
            pkgs := utils.MapLibsToPackages(libs)
            if err := pkgmgr.Install(projectName, pm, pkgs...); err != nil {
                return err
            }
        }
    }

    // mixins selection
    if len(tmpl.Mixins) > 0 {
        var options []string
        for _, m := range tmpl.Mixins {
            options = append(options, m.Name())
        }
        var selected []string
        if err := utils.MultiSelect("🔧 Enable additional features?", &selected, options); err != nil {
            return err
        }
        for _, name := range selected {
            for _, m := range tmpl.Mixins {
                if m.Name() == name {
                    if err := m.Apply(projectName, pm); err != nil {
                        return fmt.Errorf("apply %s: %w", name, err)
                    }
                }
            }
        }
    }

    utils.Log("🎉 Project '%s' successfully created!", projectName)
    return nil
}
