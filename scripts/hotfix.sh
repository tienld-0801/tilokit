#!/bin/bash

# TiLoKit Hotfix Script
# Usage: ./scripts/hotfix.sh [patch_version]
# Example: ./scripts/hotfix.sh v0.1.1

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
CURRENT_BRANCH=$(git branch --show-current)
DEVELOP_BRANCH="develop"
MAIN_BRANCH="main"

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

# Function to validate version format
validate_version() {
    local version=$1
    if [[ ! $version =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        print_error "Invalid hotfix version format: $version"
        print_info "Expected format: v1.0.1 (patch version only)"
        exit 1
    fi
}

# Function to check if working directory is clean
check_clean_working_dir() {
    if [[ -n $(git status --porcelain) ]]; then
        print_error "Working directory is not clean. Please commit or stash your changes."
        git status --short
        exit 1
    fi
}

# Function to check if we're on main branch
check_main_branch() {
    if [[ "$CURRENT_BRANCH" != "$MAIN_BRANCH" ]]; then
        print_error "You must be on the $MAIN_BRANCH branch to create a hotfix."
        print_info "Current branch: $CURRENT_BRANCH"
        exit 1
    fi
}

# Function to get latest version tag
get_latest_version() {
    git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0"
}

# Function to validate hotfix version is next patch
validate_hotfix_version() {
    local new_version=$1
    local latest_version
    latest_version=$(get_latest_version)
    
    print_info "Latest version: $latest_version"
    print_info "Hotfix version: $new_version"
    
    # Extract version numbers
    if [[ $latest_version =~ v([0-9]+)\.([0-9]+)\.([0-9]+) ]]; then
        local latest_major=${BASH_REMATCH[1]}
        local latest_minor=${BASH_REMATCH[2]}
        local latest_patch=${BASH_REMATCH[3]}
    else
        print_error "Could not parse latest version: $latest_version"
        exit 1
    fi
    
    if [[ $new_version =~ v([0-9]+)\.([0-9]+)\.([0-9]+) ]]; then
        local new_major=${BASH_REMATCH[1]}
        local new_minor=${BASH_REMATCH[2]}
        local new_patch=${BASH_REMATCH[3]}
    else
        print_error "Could not parse new version: $new_version"
        exit 1
    fi
    
    # Validate it's a proper patch increment
    if [[ $new_major -ne $latest_major ]] || [[ $new_minor -ne $latest_minor ]] || [[ $new_patch -ne $((latest_patch + 1)) ]]; then
        print_error "Hotfix version must be a patch increment of the latest version"
        print_info "Expected: v$latest_major.$latest_minor.$((latest_patch + 1))"
        print_info "Got: $new_version"
        exit 1
    fi
}

# Function to update changelog for hotfix
update_changelog_hotfix() {
    local version=$1
    local date=$(date +%Y-%m-%d)
    
    print_info "Updating CHANGELOG.md for hotfix..."
    
    # Create temporary file
    local temp_file=$(mktemp)
    
    # Add hotfix entry after Unreleased section
    awk -v version="$version" -v date="$date" '
    /^## \[Unreleased\]$/ {
        print $0
        getline
        print $1
        while ((getline) && !/^## \[/) {
            print
        }
        print "## [" substr(version, 2) "] - " date
        print ""
        print "### Fixed"
        print "- Critical bug fixes"
        print ""
        print "## [" $0
        next
    }
    1
    ' CHANGELOG.md > "$temp_file"
    
    # Replace the original file
    mv "$temp_file" CHANGELOG.md
    
    print_success "CHANGELOG.md updated for hotfix $version"
}

# Function to create hotfix branch
create_hotfix_branch() {
    local version=$1
    local hotfix_branch="hotfix/$version"
    
    print_info "Creating hotfix branch: $hotfix_branch"
    
    # Create and checkout hotfix branch from main
    git checkout -b "$hotfix_branch"
    
    print_success "Hotfix branch $hotfix_branch created"
    echo "$hotfix_branch"
}

# Function to update version in code
update_version_in_code() {
    local version=$1
    
    print_info "Updating version in cmd/version.go..."
    
    # Update version.go
    sed -i "s/Version = \".*\"/Version = \"$version\"/" cmd/version.go
    
    print_success "Version updated to $version in code"
}

# Function to commit hotfix changes
commit_hotfix_changes() {
    local version=$1
    
    print_info "Committing hotfix changes..."
    
    git add CHANGELOG.md cmd/version.go
    git commit -m "hotfix: release $version

- Critical bug fixes
- Update CHANGELOG.md
- Bump version to $version"
    
    print_success "Hotfix changes committed"
}

# Function to create and push tag
create_and_push_tag() {
    local version=$1
    local hotfix_branch=$2
    
    print_info "Creating and pushing tag $version..."
    
    # Create annotated tag
    git tag -a "$version" -m "Hotfix $version

Critical bug fixes and patches."
    
    # Push branch and tag
    git push origin "$hotfix_branch"
    git push origin "$version"
    
    print_success "Tag $version created and pushed"
}

# Function to merge hotfix
merge_hotfix() {
    local version=$1
    local hotfix_branch=$2
    
    print_info "Merging hotfix to main branch..."
    
    # Switch to main and merge
    git checkout "$MAIN_BRANCH"
    git pull origin "$MAIN_BRANCH"
    git merge --no-ff "$hotfix_branch" -m "Merge hotfix $version"
    git push origin "$MAIN_BRANCH"
    
    print_success "Merged to $MAIN_BRANCH"
    
    # Switch to develop and merge
    print_info "Merging hotfix back to develop..."
    git checkout "$DEVELOP_BRANCH"
    git pull origin "$DEVELOP_BRANCH"
    git merge --no-ff "$hotfix_branch" -m "Merge hotfix $version back to develop"
    git push origin "$DEVELOP_BRANCH"
    
    print_success "Merged back to $DEVELOP_BRANCH"
    
    # Clean up hotfix branch
    print_info "Cleaning up hotfix branch..."
    git branch -d "$hotfix_branch"
    git push origin --delete "$hotfix_branch"
    
    print_success "Hotfix branch cleaned up"
}

# Main function
main() {
    local version=$1
    
    print_info "🔧 Starting TiLoKit Hotfix Process"
    print_info "=================================="
    
    # Validate input
    if [[ -z "$version" ]]; then
        print_error "Version is required"
        print_info "Usage: $0 <patch_version>"
        print_info "Example: $0 v0.1.1"
        exit 1
    fi
    
    validate_version "$version"
    
    # Pre-flight checks
    print_info "Running pre-flight checks..."
    check_clean_working_dir
    check_main_branch
    validate_hotfix_version "$version"
    
    # Ensure we're up to date
    print_info "Updating main branch..."
    git pull origin "$MAIN_BRANCH"
    
    # Create hotfix branch
    local hotfix_branch
    hotfix_branch=$(create_hotfix_branch "$version")
    
    print_warning "Please make your hotfix changes now and commit them."
    print_warning "Press Enter when ready to continue with the release process..."
    read -r
    
    # Update files
    update_changelog_hotfix "$version"
    update_version_in_code "$version"
    
    # Commit changes
    commit_hotfix_changes "$version"
    
    # Create and push tag
    create_and_push_tag "$version" "$hotfix_branch"
    
    # Merge hotfix
    merge_hotfix "$version" "$hotfix_branch"
    
    print_success "🎉 Hotfix $version completed successfully!"
    print_info "Release workflow will now build and deploy automatically."
    print_info "Monitor the progress at: https://github.com/ti-lo/tilokit/actions"
}

# Run main function
main "$@"
