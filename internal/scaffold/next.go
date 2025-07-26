package scaffold

import (
	"fmt"
	"os"
)

// GenerateNextOptions creates a new directory for a Next project with the specified name.
// Returns an error if the directory cannot be created.
func GenerateNextOptions(projectName string) error {
	fmt.Println("🚧 Create template Next", projectName)

	if err := os.MkdirAll(projectName, os.ModePerm); err != nil {
		return fmt.Errorf("internal error creating project directory: %w", err)
	}

	return nil
}
