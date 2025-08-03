package main

import (
	"os"

	"github.com/ti-lo/tilokit/cmd"
	"github.com/ti-lo/tilokit/internal/utils"
)

// main is the entry point of the application and invokes the command execution logic.
func main() {
	// Set up graceful error handling
	defer func() {
		if r := recover(); r != nil {
			utils.Error("Fatal error: %v", r)
			os.Exit(1)
		}
	}()

	// Execute the CLI
	cmd.Execute()
}
