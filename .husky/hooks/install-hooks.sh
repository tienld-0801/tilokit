#!/bin/bash

# TiLoKit Git Hooks Installer (.husky style)
# Installs commit message validation and pre-commit hooks

# ANSI color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${CYAN}🔧 TiLoKit Git Hooks Installer (.husky)${NC}"
echo -e "${BLUE}=========================================${NC}"

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo -e "${RED}❌ Error: Not in a git repository${NC}"
    echo -e "${YELLOW}💡 Please run this script from the root of your git repository${NC}"
    exit 1
fi

# Get git hooks directory and .husky directory
GIT_HOOKS_DIR="$(git rev-parse --git-dir)/hooks"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
HUSKY_DIR="$(git rev-parse --show-toplevel)/.husky"

echo -e "${BLUE}📂 Git hooks directory: ${GIT_HOOKS_DIR}${NC}"
echo -e "${BLUE}📂 .husky directory: ${HUSKY_DIR}${NC}"
echo -e "${BLUE}📂 Source hooks directory: ${SCRIPT_DIR}${NC}"
echo

# Available hooks to install
HOOKS=(
    "commit-msg:Validates commit message format (conventional commits)"
    "pre-commit:Runs pre-commit checks (syntax, large files, etc.)"
    "prepare-commit-msg:Auto-inserts emoji based on commit type"
)

# Function to install a hook
install_hook() {
    local hook_name=$1
    local hook_description=$2
    local source_file="${SCRIPT_DIR}/${hook_name}"
    local target_file="${GIT_HOOKS_DIR}/${hook_name}"

    if [[ ! -f "$source_file" ]]; then
        echo -e "${RED}❌ Source hook not found: $source_file${NC}"
        return 1
    fi

    # Backup existing hook if it exists
    if [[ -f "$target_file" ]]; then
        echo -e "${YELLOW}⚠️  Existing hook found: $hook_name${NC}"
        cp "$target_file" "${target_file}.backup.$(date +%Y%m%d_%H%M%S)"
        echo -e "${YELLOW}📦 Backed up to: ${target_file}.backup.$(date +%Y%m%d_%H%M%S)${NC}"
    fi

    # Copy and make executable
    cp "$source_file" "$target_file"
    chmod +x "$target_file"

    echo -e "${GREEN}✅ Installed: $hook_name${NC}"
    echo -e "${BLUE}   Description: $hook_description${NC}"
    return 0
}

# Install hooks
echo -e "${CYAN}🚀 Installing Git hooks...${NC}"
echo

SUCCESS_COUNT=0
TOTAL_COUNT=${#HOOKS[@]}

for hook_info in "${HOOKS[@]}"; do
    IFS=':' read -r hook_name hook_description <<< "$hook_info"
    if install_hook "$hook_name" "$hook_description"; then
        ((SUCCESS_COUNT++))
    fi
    echo
done

# Summary
echo -e "${BLUE}===============================${NC}"
echo -e "${CYAN}📊 Installation Summary${NC}"
echo -e "${BLUE}===============================${NC}"
echo -e "${GREEN}✅ Successfully installed: $SUCCESS_COUNT/$TOTAL_COUNT hooks${NC}"
echo

if [[ $SUCCESS_COUNT -eq $TOTAL_COUNT ]]; then
    echo -e "${GREEN}🎉 All hooks installed successfully!${NC}"
    echo
    echo -e "${CYAN}📋 What happens now:${NC}"
    echo -e "${BLUE}• commit-msg hook:${NC} Validates every commit message format"
    echo -e "${BLUE}• pre-commit hook:${NC} Runs checks before each commit"
    echo -e "${BLUE}• .husky directory:${NC} Stores all git hooks (like JavaScript projects)"
    echo
    echo -e "${YELLOW}💡 Valid commit message format:${NC}"
    echo -e "   ${GREEN}type(scope): description${NC}"
    echo
    echo -e "${YELLOW}📋 Valid types:${NC}"
    echo -e "   ${GREEN}✨ feat${NC}      - new features"
    echo -e "   ${GREEN}🐛 fix${NC}       - bug fixes"
    echo -e "   ${GREEN}📚 docs${NC}      - documentation"
    echo -e "   ${GREEN}♻️  refactor${NC}  - code refactoring"
    echo -e "   ${GREEN}⚡ perf${NC}      - performance"
    echo -e "   ${GREEN}🧪 test${NC}      - tests"
    echo -e "   ${GREEN}🛠️  build${NC}     - build system"
    echo -e "   ${GREEN}🔄 ci${NC}        - CI/CD"
    echo -e "   ${GREEN}🧹 chore${NC}     - maintenance"
    echo -e "   ${GREEN}🎨 style${NC}     - code style"
    echo -e "   ${GREEN}⏪ revert${NC}    - reverts"
    echo
    echo -e "${YELLOW}🚀 Example commits:${NC}"
    echo -e "   ${GREEN}✅ feat: add user authentication${NC}"
    echo -e "   ${GREEN}✅ fix(core): resolve memory leak${NC}"
    echo -e "   ${GREEN}✅ docs: update API documentation${NC}"
    echo -e "   ${GREEN}✅ chore(deps): upgrade dependencies${NC}"
    echo
    echo -e "${CYAN}🔧 To uninstall hooks later:${NC}"
    echo -e "   ${BLUE}rm ${GIT_HOOKS_DIR}/commit-msg${NC}"
    echo -e "   ${BLUE}rm ${GIT_HOOKS_DIR}/pre-commit${NC}"
    echo -e "   ${BLUE}# Or remove the entire .husky directory: rm -rf ${HUSKY_DIR}${NC}"
else
    echo -e "${RED}❌ Some hooks failed to install${NC}"
    echo -e "${YELLOW}💡 Please check the errors above and try again${NC}"
    exit 1
fi
