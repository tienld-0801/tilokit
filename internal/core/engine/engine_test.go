package engine

import (
	"testing"

	tilocontext "github.com/ti-lo/tilokit/internal/core/context"
)

func TestEngineNew(t *testing.T) {
	engine := New()
	if engine == nil {
		t.Fatal("Expected engine to be created, got nil")
	}
}

func TestEngineRegisterPlugin(t *testing.T) {
	engine := New()

	// Mock plugin for testing
	mockPlugin := &MockPlugin{}

	err := engine.RegisterPlugin(mockPlugin)
	if err != nil {
		t.Fatalf("Expected no error registering plugin, got: %v", err)
	}
}

// MockPlugin for testing
type MockPlugin struct{}

func (p *MockPlugin) Name() string                                         { return "mock" }
func (p *MockPlugin) Version() string                                      { return "1.0.0" }
func (p *MockPlugin) Description() string                                  { return "Mock plugin for testing" }
func (p *MockPlugin) SupportedFrameworks() []string                        { return []string{"mock"} }
func (p *MockPlugin) SupportedBuildTools() []string                        { return []string{"mock"} }
func (p *MockPlugin) PreGenerate(ctx *tilocontext.ExecutionContext) error  { return nil }
func (p *MockPlugin) Generate(ctx *tilocontext.ExecutionContext) error     { return nil }
func (p *MockPlugin) PostGenerate(ctx *tilocontext.ExecutionContext) error { return nil }
