#!/bin/bash

# TiLoKit Status Check Script
# Shows comprehensive project status

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Function to print colored output
print_header() {
    echo -e "${BOLD}${CYAN}$1${NC}"
    echo -e "${CYAN}$(printf '=%.0s' {1..50})${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Main function
main() {
    print_header "🚀 TiLoKit Project Status"
    echo
    
    # Project Information
    print_header "📋 Project Information"
    echo "Project: TiLoKit CLI Toolkit"
    echo "Repository: $(git config --get remote.origin.url 2>/dev/null || echo 'Not configured')"
    echo "Directory: $(pwd)"
    echo
    
    # Version Information
    print_header "📦 Version Information"
    if [ -f "cmd/version.go" ]; then
        VERSION=$(grep 'Version = ' cmd/version.go | sed 's/.*Version = "\(.*\)".*/\1/')
        echo "Code Version: $VERSION"
    fi
    
    if [ -f "tilokit" ]; then
        echo "Binary Version:"
        ./tilokit version 2>/dev/null || echo "  Binary not working"
    else
        print_warning "Binary not built (run 'make build')"
    fi
    echo
    
    # Git Status
    print_header "🌿 Git Status"
    echo "Current Branch: $(git branch --show-current 2>/dev/null || echo 'Unknown')"
    echo "Latest Commit: $(git log --oneline -1 2>/dev/null || echo 'No commits')"
    echo "Working Directory:"
    if [[ -n $(git status --porcelain 2>/dev/null) ]]; then
        git status --short
    else
        print_success "Clean"
    fi
    
    echo "Recent Tags:"
    git tag --sort=-version:refname | head -5 2>/dev/null || echo "  No tags found"
    echo
    
    # Release System Status
    print_header "🚀 Release System"
    echo "Release Scripts:"
    for script in release.sh hotfix.sh init-branches.sh status.sh; do
        if [ -f "scripts/$script" ] && [ -x "scripts/$script" ]; then
            print_success "$script"
        else
            print_error "$script missing or not executable"
        fi
    done
    
    echo
    echo "GitHub Workflows:"
    for workflow in ci.yml release.yml; do
        if [ -f ".github/workflows/$workflow" ]; then
            print_success "$workflow"
        else
            print_error "$workflow missing"
        fi
    done
    echo
    
    # Next Steps
    print_header "🎯 Suggested Next Steps"
    CURRENT_BRANCH=$(git branch --show-current 2>/dev/null || echo "unknown")
    
    if [ "$CURRENT_BRANCH" != "develop" ] && [ "$CURRENT_BRANCH" != "main" ]; then
        print_info "Initialize branches: make init-branches"
    fi
    
    if [ ! -f "tilokit" ]; then
        print_info "Build the project: make build"
    fi
    
    if ! git tag >/dev/null 2>&1 || [ -z "$(git tag)" ]; then
        print_info "Create first release: make release VERSION=v0.1.0"
    fi
    
    print_info "Check release readiness: make check-release"
    print_info "View all commands: make help"
    echo
    
    print_success "🎉 TiLoKit release system is ready!"
}

# Run main function
main "$@"
