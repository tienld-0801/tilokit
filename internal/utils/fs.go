package utils

import "os"

// PathExists returns true if the given path exists on the filesystem.
func PathExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}
	// For other errors, assume it exists to be safe.
	return true
}
