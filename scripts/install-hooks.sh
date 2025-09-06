#!/bin/bash

# TiLoKit Git Hooks Installer
# Installs commit message validation and auto-emoji hooks

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔧 Installing TiLoKit Git Hooks...${NC}"

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo -e "${RED}❌ Error: Not in a git repository${NC}"
    exit 1
fi

# Create hooks directory if it doesn't exist
mkdir -p .git/hooks

# Copy hooks from .githooks to .git/hooks
if [[ -f ".githooks/commit-msg" ]]; then
    cp .githooks/commit-msg .git/hooks/commit-msg
    chmod +x .git/hooks/commit-msg
    echo -e "${GREEN}✅ Installed commit-msg hook${NC}"
else
    echo -e "${YELLOW}⚠️  commit-msg hook not found in .githooks/${NC}"
fi

if [[ -f ".githooks/prepare-commit-msg" ]]; then
    cp .githooks/prepare-commit-msg .git/hooks/prepare-commit-msg
    chmod +x .git/hooks/prepare-commit-msg
    echo -e "${GREEN}✅ Installed prepare-commit-msg hook${NC}"
else
    echo -e "${YELLOW}⚠️  prepare-commit-msg hook not found in .githooks/${NC}"
fi

echo -e "${GREEN}🎉 Git hooks installed successfully!${NC}"
echo -e "${BLUE}📋 Hooks will now:${NC}"
echo -e "  • Require valid commit types (feat, fix, docs, etc.)"
echo -e "  • Auto-add emojis to valid commit messages"
echo -e "  • Reject commits without proper tags"
echo
echo -e "${YELLOW}💡 To test:${NC} git commit -m \"feat: add new feature\""
