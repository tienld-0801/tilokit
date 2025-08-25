package progress

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"tilokit/internal/core/engine"
	tilocontext "tilokit/internal/core/context"
	"tilokit/internal/ui"
)

// ProgressHandler handles project generation progress tracking and UI
type ProgressHandler struct{}

// NewProgressHandler creates a new ProgressHandler instance
func NewProgressHandler() *ProgressHandler {
	return &ProgressHandler{}
}

// ExecuteWithProgress runs project generation with animated progress updates
func (p *ProgressHandler) ExecuteWithProgress(
	eng *engine.Engine,
	ctx context.Context,
	projectConfig *tilocontext.ProjectConfig,
	projectName, framework string,
) error {
	// Define generation steps
	steps := []string{
		"Initializing project structure",
		"Setting up framework configuration",
		"Installing build tools",
		"Creating source files",
		"Configuring development environment",
		"Finalizing project setup",
	}

	// Create progress channel
	progressChan := make(chan ui.ProgressMsg, 10)

	// Start animated progress UI
	doneCh, err := ui.RunProjectProgress(
		fmt.Sprintf("🚀 Creating %s Project: %s", framework, projectName),
		steps,
		progressChan,
	)
	if err != nil {
		logrus.Errorf("Progress UI error: %v", err)
		return err
	}

	// Send initial progress
	progressChan <- ui.ProgressMsg{Step: "Starting project generation...", Progress: 0.0}

	// Execute project generation with progress updates
	if err := p.executeStepsWithProgress(eng, ctx, projectConfig, progressChan); err != nil {
		progressChan <- ui.ProgressMsg{Step: "❌ Project generation failed", Progress: 1.0, Done: true}
		close(progressChan)
		<-doneCh
		logrus.Errorf("Project generation failed: %v", err)
		return err
	}

	// Send completion
	progressChan <- ui.ProgressMsg{Step: "✅ Project created successfully!", Progress: 1.0, Done: true}
	close(progressChan)

	// Wait for UI to clean up
	<-doneCh

	return nil
}

// executeStepsWithProgress executes the actual generation steps with progress updates
func (p *ProgressHandler) executeStepsWithProgress(
	eng *engine.Engine,
	ctx context.Context,
	projectConfig *tilocontext.ProjectConfig,
	progressChan chan<- ui.ProgressMsg,
) error {
	steps := []struct {
		name     string
		progress float64
	}{
		{"Initializing project structure", 0.16},
		{"Setting up framework configuration", 0.33},
		{"Installing build tools", 0.50},
		{"Creating source files", 0.66},
		{"Configuring development environment", 0.83},
		{"Finalizing project setup", 0.99},
	}

	for i, step := range steps {
		progressChan <- ui.ProgressMsg{Step: step.name, Progress: step.progress}
		time.Sleep(350 * time.Millisecond)

		// Execute actual generation on the last step
		if i == len(steps)-1 {
			if err := eng.Execute(ctx, projectConfig); err != nil {
				return err
			}
		}
	}

	return nil
}
