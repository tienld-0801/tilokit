package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ReadFile reads content from a file with path validation
func ReadFile(path string) (string, error) {
	// Validate and clean the path to prevent path traversal
	cleanPath, err := validatePath(path)
	if err != nil {
		return "", err
	}
	
	// #nosec G304 - Path is validated above
	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteFile writes content to a file, creating directories if needed
func WriteFile(path, content string) error {
	dir := filepath.Dir(path)
	if err := EnsureDir(dir); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

// EnsureDir creates a directory and all parent directories if they don't exist
func EnsureDir(path string) error {
	// Use more restrictive permissions (0750 instead of 0755)
	return os.MkdirAll(path, 0750)
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// DirExists checks if a directory exists
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// CopyFile copies a file from src to dst with path validation
func CopyFile(src, dst string) error {
	// Validate paths
	cleanSrc, err := validatePath(src)
	if err != nil {
		return fmt.Errorf("invalid source path: %w", err)
	}
	
	cleanDst, err := validatePath(dst)
	if err != nil {
		return fmt.Errorf("invalid destination path: %w", err)
	}
	
	// #nosec G304 - Path is validated above
	data, err := os.ReadFile(cleanSrc)
	if err != nil {
		return err
	}
	return WriteFile(cleanDst, string(data))
}

// RemoveFile removes a file if it exists
func RemoveFile(path string) error {
	if FileExists(path) {
		return os.Remove(path)
	}
	return nil
}

// RemoveDir removes a directory and all its contents
func RemoveDir(path string) error {
	if DirExists(path) {
		return os.RemoveAll(path)
	}
	return nil
}

// validatePath validates and cleans a file path to prevent path traversal attacks
func validatePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}
	
	// Clean the path to resolve any . or .. components
	cleanPath := filepath.Clean(path)
	
	// Check for path traversal attempts
	if strings.Contains(cleanPath, "..") {
		return "", fmt.Errorf("path traversal detected: %s", path)
	}
	
	// Ensure path is absolute or relative to current directory
	if !filepath.IsAbs(cleanPath) {
		absPath, err := filepath.Abs(cleanPath)
		if err != nil {
			return "", fmt.Errorf("failed to resolve absolute path: %w", err)
		}
		cleanPath = absPath
	}
	
	return cleanPath, nil
}
