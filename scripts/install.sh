#!/bin/bash

# TiLoKit Installation Script
set -e

REPO="ti-lo/tilokit"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="tilokit"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 TiLoKit Installer${NC}"
echo "=================================="

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo -e "${RED}❌ Unsupported architecture: $ARCH${NC}"; exit 1 ;;
esac

case $OS in
    linux) OS="linux" ;;
    darwin) OS="darwin" ;;
    *) echo -e "${RED}❌ Unsupported OS: $OS${NC}"; exit 1 ;;
esac

echo -e "${BLUE}📋 Detected platform: ${OS}-${ARCH}${NC}"

# Get latest release
echo -e "${YELLOW}🔍 Fetching latest release...${NC}"
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
    echo -e "${RED}❌ Failed to fetch latest release${NC}"
    exit 1
fi

echo -e "${GREEN}📦 Latest version: $LATEST_RELEASE${NC}"

# Download binary
BINARY_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/tilokit-$OS-$ARCH"
if [ "$OS" = "windows" ]; then
    BINARY_URL="${BINARY_URL}.exe"
fi

echo -e "${YELLOW}⬇️  Downloading TiLoKit...${NC}"
TEMP_FILE=$(mktemp)
curl -L -o "$TEMP_FILE" "$BINARY_URL"

# Check if download was successful
if [ ! -f "$TEMP_FILE" ] || [ ! -s "$TEMP_FILE" ]; then
    echo -e "${RED}❌ Failed to download TiLoKit${NC}"
    exit 1
fi

# Make executable
chmod +x "$TEMP_FILE"

# Install binary
echo -e "${YELLOW}📦 Installing TiLoKit to $INSTALL_DIR...${NC}"
if [ -w "$INSTALL_DIR" ]; then
    mv "$TEMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
else
    echo -e "${YELLOW}🔐 Requesting sudo access to install to $INSTALL_DIR...${NC}"
    sudo mv "$TEMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
fi

# Verify installation
if command -v tilokit &> /dev/null; then
    echo -e "${GREEN}✅ TiLoKit installed successfully!${NC}"
    echo ""
    echo -e "${BLUE}🎉 You can now use TiLoKit:${NC}"
    echo "   tilokit --help"
    echo "   tilokit list"
    echo "   tilokit create my-project"
    echo ""
    tilokit version
else
    echo -e "${RED}❌ Installation failed${NC}"
    exit 1
fi
