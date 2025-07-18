package scaffold

import (
	"fmt"
	"os"
)

// GenerateLaravel creates a new directory for a Laravel project with the specified name.
// Returns an error if the directory cannot be created.
func GenerateLaravel(projectName string) error {
	fmt.Println("🚧 Create template Laravel:", projectName)

	if err := os.MkdirAll(projectName, os.ModePerm); err != nil {
		return fmt.Errorf("internal error creating project directory: %w", err)
	}

	return nil
}
