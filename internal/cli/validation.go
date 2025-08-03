package cli

import (
	"fmt"
	"strings"
)

// ValidateFlagUsage validates that flags follow proper format
func ValidateFlagUsage(args []string) error {
	for _, arg := range args {
		// Skip non-flag arguments and help flag
		if !strings.HasPrefix(arg, "-") || arg == "--help" || arg == "-h" {
			continue
		}

		// Single dash flags must be exactly 2 characters (-x)
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") {
			if len(arg) > 2 {
				// Check if it's a known long flag being used with single dash
				longFlag := arg[1:]
				for _, known := range KnownLongFlags {
					if longFlag == known {
						return fmt.Errorf(InvalidFlagMsg, arg, longFlag)
					}
				}
				return fmt.Errorf(InvalidFlagGenericMsg, arg)
			}
		}
	}
	return nil
}
