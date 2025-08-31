package frameworks

import (
	tilocontext "tilokit/internal/core/context"
	"tilokit/pkg/constants"
)

// IonicPlugin implements Ionic framework support
type IonicPlugin struct{}

func NewIonicPlugin() *IonicPlugin {
	return &IonicPlugin{}
}

func (p *IonicPlugin) Name() string {
	return "ionic"
}

func (p *IonicPlugin) Version() string {
	return constants.VERSION
}

func (p *IonicPlugin) Description() string {
	return "Ionic hybrid mobile app framework"
}

func (p *IonicPlugin) SupportedFrameworks() []string {
	return []string{"ionic"}
}

func (p *IonicPlugin) SupportedBuildTools() []string {
	return []string{"ionic-cli", "capacitor", "cordova"}
}

func (p *IonicPlugin) PreGenerate(ctx *tilocontext.ExecutionContext) error {
	// TODO: Implement Ionic pre-generation logic
	return nil
}

func (p *IonicPlugin) Generate(ctx *tilocontext.ExecutionContext) error {
	// TODO: Implement Ionic project generation
	return nil
}

func (p *IonicPlugin) PostGenerate(ctx *tilocontext.ExecutionContext) error {
	// TODO: Implement Ionic post-generation logic
	return nil
}
