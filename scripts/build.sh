#!/bin/bash

# TiLoKit Build Script
set -e

echo "🚀 Building TiLoKit..."

# Get version info
VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS="-X main.Version=$VERSION -X main.BuildDate=$BUILD_DATE -X main.GitCommit=$GIT_COMMIT"

# Clean previous builds
echo "🧹 Cleaning previous builds..."
rm -rf dist/

# Create dist directory
mkdir -p dist/

# Build for multiple platforms
echo "🔨 Building for multiple platforms..."

# Linux
echo "  📦 Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o dist/tilokit-linux-amd64 .

echo "  📦 Building for Linux (arm64)..."
GOOS=linux GOARCH=arm64 go build -ldflags "$LDFLAGS" -o dist/tilokit-linux-arm64 .

# macOS
echo "  📦 Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFLAGS" -o dist/tilokit-darwin-amd64 .

echo "  📦 Building for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -ldflags "$LDFLAGS" -o dist/tilokit-darwin-arm64 .

# Windows
echo "  📦 Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -ldflags "$LDFLAGS" -o dist/tilokit-windows-amd64.exe .

# Create checksums
echo "🔐 Creating checksums..."
cd dist/
for file in tilokit-*; do
    if [[ -f "$file" ]]; then
        sha256sum "$file" > "$file.sha256"
    fi
done
cd ..

echo "✅ Build completed successfully!"
echo "📁 Binaries available in dist/ directory"

# List built files
echo "📋 Built files:"
ls -la dist/
