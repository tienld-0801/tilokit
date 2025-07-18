package utils

import (
	"fmt"
	"os"
)

func IsProduction() bool {
	return os.Getenv("TILOKIT_ENV") == "production"
}

func Log(msg string, args ...any) {
	if !IsProduction() {
		fmt.Printf(msg+"\n", args...)
	}
}

func Error(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
}
