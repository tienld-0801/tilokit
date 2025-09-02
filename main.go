package main

import (
	"os"

	"tilokit/cmd"
	"tilokit/internal/cli"
	"tilokit/internal/utils"
)

func main() {
	if err := cli.ValidateFlagUsage(os.Args[1:]); err != nil {
		utils.Error("%v", err)
		os.Exit(1)
	}

	defer func() {
		if r := recover(); r != nil {
			utils.Error("Fatal error: %v", r)
			os.Exit(1)
		}
	}()

	cmd.Execute()
}
