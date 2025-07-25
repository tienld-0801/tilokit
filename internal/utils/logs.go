package utils

import (
	"fmt"
	"os"
)

var quiet bool

// SetQuiet toggles quiet mode; when enabled Log will be silenced.
func SetQuiet(q bool) {
	quiet = q
}

func IsProduction() bool {
	return os.Getenv("TILOKIT_ENV") == "production"
}

func Log(msg string, args ...any) {
	if quiet {
		return
	}
	if !IsProduction() {
		fmt.Printf(msg+"\n", args...)
	}
}

func Error(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
}
