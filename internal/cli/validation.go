package cli

import (
	"fmt"
	"strings"

	"tilokit/pkg/constants"
)

func ValidateFlagUsage(args []string) error {
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") || arg == "--help" || arg == "-h" {
			continue
		}

		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") {
			if len(arg) > 2 {
				longFlag := arg[1:]
				for _, known := range constants.KnownLongFlags {
					if longFlag == known {
						return fmt.Errorf(constants.InvalidFlagMsg, arg, longFlag)
					}
				}
				return fmt.Errorf(constants.InvalidFlagGenericMsg, arg)
			}
		}
	}
	return nil
}
