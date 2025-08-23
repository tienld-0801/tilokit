package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomSecret generates a cryptographically secure random secret
func GenerateRandomSecret(byteLength int) (string, error) {
	bytes := make([]byte, byteLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateBase64Secret generates a cryptographically secure random secret in base64 format
func GenerateBase64Secret(byteLength int) (string, error) {
	bytes := make([]byte, byteLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	// Laravel expects base64 format with base64: prefix
	return "base64:" + hex.EncodeToString(bytes), nil
}
