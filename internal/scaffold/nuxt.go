package scaffold

import (
	"fmt"
	"os"
)

func GenerateNuxtOptions(projectName string) error {
	fmt.Println("🚧 Create template Nuxt:", projectName)

	if err := os.MkdirAll(projectName, os.ModePerm); err != nil {
		return fmt.Errorf("internal error creating project directory: %w", err)
	}

	return nil
}
