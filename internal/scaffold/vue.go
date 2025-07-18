package scaffold

import (
	"fmt"
	"os"
)

// GenerateLaravel creates a new directory for a Laravel project with the specified name.
// Returns an error if the directory cannot be created.
func GenerateVue(projectName string) error {
	fmt.Println("🚧 Create template Vue:", projectName)

	if err := os.MkdirAll(projectName, os.ModePerm); err != nil {
		return fmt.Errorf("internal error creating project directory: %w", err)
	}

	return nil
}
