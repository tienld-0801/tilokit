package scaffold

import (
	"fmt"
	"os"
)

func GenerateLaravelOptions(projectName string) error {
	fmt.Println("🚧 Create template Laravel:", projectName)

	if err := os.MkdirAll(projectName, os.ModePerm); err != nil {
		return fmt.Errorf("internal error creating project directory: %w", err)
	}

	return nil
}
