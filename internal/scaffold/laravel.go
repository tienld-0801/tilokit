package scaffold

import (
	"fmt"
	"os"
)

// GenerateLaravelOptions creates a new directory for a Laravel project with the specified name.
func GenerateLaravelOptions(projectName string) error {
	fmt.Println("🚧 Create template Laravel:", projectName)

	if err := os.MkdirAll(projectName, os.ModePerm); err != nil {
		return fmt.Errorf("internal error creating project directory: %w", err)
	}

	return nil
}
