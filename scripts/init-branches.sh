#!/bin/bash

# TiLoKit Branch Initialization Script
# Sets up the proper branch structure for GitFlow

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
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

# Check if we're in a git repository
check_git_repo() {
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        print_error "Not in a git repository"
        exit 1
    fi
}

# Create branch if it doesn't exist
create_branch() {
    local branch_name=$1
    local base_branch=${2:-main}
    
    if git show-ref --verify --quiet refs/heads/$branch_name; then
        print_warning "Branch $branch_name already exists"
    else
        print_info "Creating branch $branch_name from $base_branch"
        git checkout $base_branch
        git checkout -b $branch_name
        git push -u origin $branch_name
        print_success "Created branch $branch_name"
    fi
}

# Set up remote tracking
setup_remote_tracking() {
    local branch_name=$1
    
    print_info "Setting up remote tracking for $branch_name"
    git checkout $branch_name
    git push -u origin $branch_name 2>/dev/null || true
}

# Main function
main() {
    print_info "🌟 Initializing TiLoKit Branch Structure"
    print_info "======================================="
    
    # Check prerequisites
    check_git_repo
    
    # Ensure we have the main branches
    print_info "Setting up main branches..."
    
    # Create main branch if it doesn't exist (usually exists)
    if ! git show-ref --verify --quiet refs/heads/main; then
        print_warning "Main branch doesn't exist, creating it"
        git checkout -b main
        git push -u origin main
    fi
    
    # Create develop branch
    create_branch "develop" "main"
    
    # Set up remote tracking for existing branches
    for branch in main develop; do
        if git show-ref --verify --quiet refs/heads/$branch; then
            setup_remote_tracking $branch
        fi
    done
    
    # Switch to develop as the default working branch
    git checkout develop
    
    print_success "🎉 Branch structure initialized successfully!"
    print_info ""
    print_info "Branch structure:"
    print_info "  main     - Production-ready code"
    print_info "  develop  - Integration and development"
    print_info ""
    print_info "You're now on the develop branch and ready to start working!"
    print_info ""
    print_info "Quick commands:"
    print_info "  Create feature: git checkout -b feature/your-feature-name"
    print_info "  Start release: ./scripts/release.sh v0.1.0"
    print_info "  Create hotfix: ./scripts/hotfix.sh v0.1.1"
}

# Run main function
main "$@"
