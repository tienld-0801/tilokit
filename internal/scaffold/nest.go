package scaffold

import (
	"fmt"
	"os"
)

// GenerateNestOptions creates a new directory for a Nest project with the specified name.
func GenerateNestOptions(projectName string) error {
	fmt.Println("🚧 Create template Nest:", projectName)

	if err := os.MkdirAll(projectName, os.ModePerm); err != nil {
		return fmt.Errorf("internal error creating project directory: %w", err)
	}

	return nil
}
