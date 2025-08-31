package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.New()
	quiet  = false
)

// SetQuiet sets the quiet mode for logging
func SetQuiet(q bool) {
	quiet = q
	if quiet {
		logger.SetLevel(logrus.ErrorLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}
}

// IsProduction checks if running in production mode
func IsProduction() bool {
	return os.Getenv("TILOKIT_ENV") == "production"
}

// Log prints a formatted log message
func Log(format string, args ...interface{}) {
	if !quiet {
		color.Cyan(format, args...)
	}
}

// Success prints a success message
func Success(format string, args ...interface{}) {
	if !quiet {
		color.Green("✅ "+format, args...)
	}
}

// Warning prints a warning message
func Warning(format string, args ...interface{}) {
	if !quiet {
		color.Yellow("⚠️  "+format, args...)
	}
}

// Error prints an error message
func Error(format string, args ...interface{}) {
	color.Red("❌ "+format, args...)
}

// Info prints an info message
func Info(format string, args ...interface{}) {
	if !quiet {
		color.Blue("ℹ️  "+format, args...)
	}
}

// PrintBanner prints the TiLoKit banner
func PrintBanner() {
	if quiet {
		return
	}

	banner := `
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║        _______   _   _                _  __  _   _           ║
║       |__   __| (_) | |              | |/ / (_) | |          ║
║          | |     _  | |        ___   | ' /   _  | |__        ║
║          | |    | | | |       / _ \  |  <   | | | __|        ║
║          | |    | | | |____  | (_) | | . \  | | | |__        ║
║          |_|    |_| |______|  \___/  |_|\_\ |_| \____|       ║
║                                                              ║
║         🚀 Multi-Framework Project Generator 🚀              ║
║                                                              ║
║    💬 Join our Discord: https://discord.gg/BqTqZ46uT9        ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝
`
	color.Cyan(banner)
}

// Contains checks if a slice contains a string
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ToKebabCase converts a string to kebab-case
func ToKebabCase(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	return s
}

// ToPascalCase converts a string to PascalCase
func ToPascalCase(s string) string {
	words := strings.FieldsFunc(s, func(c rune) bool {
		return c == ' ' || c == '-' || c == '_'
	})

	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}

	return strings.Join(words, "")
}

// ValidateProjectName validates a project name
func ValidateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	if len(name) < 2 {
		return fmt.Errorf("project name must be at least 2 characters long")
	}

	if strings.Contains(name, " ") {
		return fmt.Errorf("project name cannot contain spaces")
	}

	return nil
}

// CommandExists checks if a command exists in PATH
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// ColorizeString colorizes a string with the given color
func ColorizeString(text, colorName string) string {
	switch colorName {
	case "red":
		return color.RedString(text)
	case "green":
		return color.GreenString(text)
	case "yellow":
		return color.YellowString(text)
	case "blue":
		return color.BlueString(text)
	case "magenta":
		return color.MagentaString(text)
	case "cyan":
		return color.CyanString(text)
	case "white":
		return color.WhiteString(text)
	case "gray", "grey":
		return color.HiBlackString(text)
	default:
		return text
	}
}
