package utils

import (
	"os"
	"path/filepath"
)

// ReadFile reads content from a file
func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
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
	return os.MkdirAll(path, 0755)
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

// CopyFile copies a file from src to dst
func CopyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return WriteFile(dst, string(data))
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
