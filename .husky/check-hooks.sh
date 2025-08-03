#!/bin/bash

# TiLoKit Hooks Checker
# Automatically checks and installs Git hooks if missing

# ANSI color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to check if hooks are installed
check_hooks_installed() {
    local git_hooks_dir="$(git rev-parse --git-dir)/hooks"
    
    if [[ -f "$git_hooks_dir/commit-msg" ]] && [[ -f "$git_hooks_dir/pre-commit" ]]; then
        return 0  # Hooks are installed
    else
        return 1  # Hooks are missing
    fi
}

# Function to install hooks automatically
auto_install_hooks() {
    echo -e "${YELLOW}⚠️  Git hooks not found!${NC}"
    echo -e "${BLUE}🔧 Auto-installing hooks for commit message validation...${NC}"
    
    # Make installer executable and run it
    chmod +x .husky/hooks/install-hooks.sh
    ./.husky/hooks/install-hooks.sh
    
    if [[ $? -eq 0 ]]; then
        echo -e "${GREEN}✅ Hooks installed successfully!${NC}"
        return 0
    else
        echo -e "${RED}❌ Failed to install hooks${NC}"
        return 1
    fi
}

# Main logic
main() {
    # Check if we're in a git repository
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        echo -e "${RED}❌ Not in a git repository${NC}"
        exit 1
    fi
    
    # Check if .husky directory exists
    if [[ ! -d ".husky" ]]; then
        echo -e "${YELLOW}⚠️  .husky directory not found${NC}"
        echo -e "${BLUE}ℹ️  This project uses Git hooks for commit validation${NC}"
        exit 0
    fi
    
    # Check if hooks are installed
    if check_hooks_installed; then
        echo -e "${GREEN}✅ Git hooks are installed${NC}"
        exit 0
    else
        # Auto-install hooks
        auto_install_hooks
        exit $?
    fi
}

# Run only if called directly (not sourced)
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
