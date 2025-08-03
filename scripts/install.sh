#!/bin/bash

# TiLoKit Installation Script
# Universal installer for Linux, macOS, and Windows (via Git Bash/WSL)
set -e

REPO="tienld-0801/tilokit"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
BINARY_NAME="tilokit"
BASE_URL="${TILOKIT_BASE_URL:-https://tienld-0801.github.io/tilokit}"
VERSION="${TILOKIT_VERSION:-latest}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Helper functions
log_info() { echo -e "${BLUE}ℹ️  $1${NC}"; }
log_success() { echo -e "${GREEN}✅ $1${NC}"; }
log_warning() { echo -e "${YELLOW}⚠️  $1${NC}"; }
log_error() { echo -e "${RED}❌ $1${NC}"; }

echo -e "${PURPLE}🚀 TiLoKit Universal Installer${NC}"
echo -e "${CYAN}===============================================${NC}"
echo -e "${BLUE}📦 Installing TiLoKit CLI Toolkit${NC}"
echo ""

# Check for required tools
check_dependencies() {
    local missing_deps=()
    
    # Essential tools
    for cmd in curl uname; do
        if ! command -v "$cmd" &> /dev/null; then
            missing_deps+=("$cmd")
        fi
    done
    
    # Hash utilities (at least one must be available)
    if ! command -v sha256sum &> /dev/null && ! command -v shasum &> /dev/null; then
        missing_deps+=("sha256sum or shasum")
    fi
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        log_error "Missing required dependencies: ${missing_deps[*]}"
        log_info "Please install these tools and try again"
        exit 1
    fi
}

# Enhanced OS and architecture detection
detect_platform() {
    local os arch

    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    arch=$(uname -m)

    # Architecture mapping
    case $arch in
        x86_64|amd64) arch="amd64" ;;
        arm64|aarch64) arch="arm64" ;;
        armv7l) arch="arm" ;;
        i386|i686) arch="386" ;;
        *) log_error "Unsupported architecture: $arch"; exit 1 ;;
    esac

    # OS detection with Windows support
    case $os in
        linux*)
            os="linux"
            if grep -q Microsoft /proc/version 2>/dev/null; then
                log_info "Detected WSL environment"
            fi
            ;;
        darwin*)
            os="darwin"
            # Check for Apple Silicon specifically
            if [ "$arch" = "arm64" ]; then
                log_info "Detected Apple Silicon Mac"
            fi
            ;;
        mingw*|msys*|cygwin*|win*)
            os="windows"
            log_info "Detected Windows environment (Git Bash/MSYS2)"
            INSTALL_DIR="${INSTALL_DIR:-$HOME/bin}"
            ;;
        *)
            log_error "Unsupported OS: $os"
            log_info "Supported platforms: Linux, macOS, Windows (Git Bash/WSL)"
            exit 1
            ;;
    esac

    export DETECTED_OS="$os"
    export DETECTED_ARCH="$arch"

    log_info "Platform: $os-$arch"
}

check_dependencies
detect_platform

# Download and install binary
download_and_install() {
    local binary_url temp_file binary_name

    # Construct download URL
    binary_name="tilokit-$DETECTED_OS-$DETECTED_ARCH"
    if [ "$DETECTED_OS" = "windows" ]; then
        binary_name="${binary_name}.exe"
    fi

    binary_url="$BASE_URL/$binary_name"
    temp_file=$(mktemp)

    log_info "Downloading from: $binary_url"

    # Download with retry logic
    local max_retries=3
    local retry_count=0

    while [ $retry_count -lt $max_retries ]; do
        if curl -fsSL --connect-timeout 10 --max-time 60 -o "$temp_file" "$binary_url"; then
            break
        else
            retry_count=$((retry_count + 1))
            if [ $retry_count -lt $max_retries ]; then
                log_warning "Download failed, retrying... ($retry_count/$max_retries)"
                sleep 2
            else
                log_error "Failed to download TiLoKit after $max_retries attempts"
                log_info "Please check your internet connection and try again"
                log_info "Manual download: $binary_url"
                exit 1
            fi
        fi
    done

    # Verify download integrity
    checksum_url="${binary_url}.sha256"
    if curl -fsSL "$checksum_url" -o "${temp_file}.sha256"; then
        # Extract expected hash and compare with actual
        expected=$(cut -d' ' -f1 "${temp_file}.sha256")
        if command -v sha256sum >/dev/null 2>&1; then
            actual=$(sha256sum "$temp_file" | cut -d' ' -f1)
        elif command -v shasum >/dev/null 2>&1; then
            actual=$(shasum -a 256 "$temp_file" | cut -d' ' -f1)
        fi
        
        if [ "$expected" != "$actual" ]; then
            log_error "Checksum verification failed"
            log_error "Expected: $expected"
            log_error "Actual: $actual"
            exit 1
        fi
        log_success "Checksum verification passed"
    else
        log_warning "Checksum file unavailable – skipping integrity check"
    fi

    log_success "Download completed ($(du -h "$temp_file" | cut -f1))"

    # Make executable (not needed for Windows but harmless)
    chmod +x "$temp_file"

    # Prepare installation directory
    if [ "$DETECTED_OS" = "windows" ]; then
        # For Windows, ensure ~/bin exists and is in PATH
        if [ ! -d "$HOME/bin" ]; then
            mkdir -p "$HOME/bin"
            log_info "Created $HOME/bin directory"
        fi
        INSTALL_DIR="$HOME/bin"
    fi

    # Install binary with proper naming
    log_info "Installing TiLoKit to $INSTALL_DIR..."
    
    # Set correct install target name
    local install_target="$INSTALL_DIR/$BINARY_NAME"
    if [ "$DETECTED_OS" = "windows" ]; then
        install_target="$INSTALL_DIR/${BINARY_NAME}.exe"
    fi
    
    if [ -w "$INSTALL_DIR" ] || [ "$DETECTED_OS" = "windows" ]; then
        if ! mv "$temp_file" "$install_target"; then
            log_error "Failed to move binary to $install_target"
            exit 1
        fi
    else
        log_warning "Requesting sudo access to install to $INSTALL_DIR..."
        if ! sudo mv "$temp_file" "$install_target"; then
            log_error "Failed to install with sudo"
            exit 1
        fi
    fi

    log_success "TiLoKit installed to $install_target"
}

# Verify installation and provide usage info
verify_installation() {
    # Use consistent install path logic
    local install_path="$INSTALL_DIR/$BINARY_NAME"
    if [ "$DETECTED_OS" = "windows" ]; then
        install_path="$INSTALL_DIR/${BINARY_NAME}.exe"
    fi

    if [ -f "$install_path" ] && [ -x "$install_path" ]; then
        log_success "TiLoKit installed successfully!"
        echo ""

        # Check if it's in PATH
        if command -v tilokit &> /dev/null; then
            log_success "TiLoKit is available in your PATH"
            echo -e "${BLUE}🎉 You can now use TiLoKit:${NC}"
            echo -e "   ${CYAN}tilokit --help${NC}       # Show help"
            echo -e "   ${CYAN}tilokit list${NC}         # List available templates"
            echo -e "   ${CYAN}tilokit create app${NC}   # Create new project"
            echo ""

            # Show version
            if tilokit version 2>/dev/null; then
                echo ""
            else
                log_warning "Could not get version info"
            fi
        else
            log_warning "TiLoKit installed but not in PATH"
            if [ "$DETECTED_OS" = "windows" ]; then
                echo -e "${YELLOW}To use TiLoKit, add $HOME/bin to your PATH or run:${NC}"
                echo -e "   ${CYAN}$install_path --help${NC}"
            else
                echo -e "${YELLOW}To use TiLoKit, add $INSTALL_DIR to your PATH or run:${NC}"
                echo -e "   ${CYAN}$install_path --help${NC}"
            fi
        fi
    else
        log_error "Installation verification failed"
        log_info "Expected location: $install_path"
        exit 1
    fi
}

# Perform installation
log_info "Starting installation process..."
download_and_install
verify_installation

log_success "Installation completed! 🎊"
log_info "For help and documentation: https://github.com/tienld-0801/tilokit"
