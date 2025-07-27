package vue

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/ti-lo/tilokit/internal/common"
	pkgmgr "github.com/ti-lo/tilokit/internal/core/pkgmgr"
	vuefw "github.com/ti-lo/tilokit/internal/framework/vue"
	"github.com/ti-lo/tilokit/internal/utils"
)

// GenerateVueOptions scaffolds a Vue project using create-vue
func Generate(projectName string) error {
	var lang string
	survey.AskOne(&survey.Select{
		Message: "Choose language for Vue template:",
		Options: []string{"JavaScript", "TypeScript"},
	}, &lang, survey.WithValidator(survey.Required))

	var sub string
	var libsKey string
	if lang == "TypeScript" {
		sub = "vue-ts"
		libsKey = "vue-vite-ts"
	} else {
		sub = "vue"
		libsKey = "vue"
	}

	// reuse new base
	pm, err := pkgmgr.AskPackageManager()
	if err != nil {
		return err
	}

	if err := vuefw.CreateVue(sub)(projectName, pm); err != nil {
		return err
	}

	libs := common.ChooseCommonLibs(libsKey)
	if len(libs) > 0 {
		pkgs := utils.MapLibsToPackages(libs)
		if err := pkgmgr.Install(projectName, pm, pkgs...); err != nil {
			return err
		}
	}

	utils.Log("🎉 Vue %s project '%s' successfully created!", lang, projectName)
	return nil
}
