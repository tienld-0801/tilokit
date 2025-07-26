package scaffold

import (
	"fmt"
	"os"
)

// GenerateAngular creates a new directory for an Angular project with the specified name.
// Returns an error if the directory cannot be created.
func GenerateAngularOptions(projectName string) error {
	fmt.Println("🚧 Create template Angular:", projectName)

	if err := os.MkdirAll(projectName, os.ModePerm); err != nil {
		return fmt.Errorf("internal error creating project directory: %w", err)
	}

	return nil
}
