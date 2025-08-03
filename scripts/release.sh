#!/bin/bash

# TiLoKit Release Script
# Usage: ./scripts/release.sh [version]
# Example: ./scripts/release.sh v0.1.0

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
    if [[ ! $version =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9\.-]+)?$ ]]; then
        print_error "Invalid version format: $version"
        print_info "Expected format: v1.0.0 or v1.0.0-beta.1"
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

# Function to check if we're on develop branch
check_develop_branch() {
    if [[ "$CURRENT_BRANCH" != "$DEVELOP_BRANCH" ]]; then
        print_error "You must be on the $DEVELOP_BRANCH branch to create a release."
        print_info "Current branch: $CURRENT_BRANCH"
        exit 1
    fi
}

# Function to update changelog
update_changelog() {
    local version=$1

    print_info "Generating changelog from conventional commits..."

    # Use the new changelog generator
    if [ -f "scripts/generate-changelog.sh" ]; then
        chmod +x scripts/generate-changelog.sh
        ./scripts/generate-changelog.sh "$version"
    else
        # Fallback to simple update
        local date=$(date +%Y-%m-%d)
        local temp_file=$(mktemp)

        awk -v version="$version" -v date="$date" '
        /^## \[Unreleased\]$/ {
            print $0
            print ""
            print "### Changed"
            print "- Development continues..."
            print ""
            print "## [" substr(version, 2) "] - " date
            print ""
            print "### Added"
            print "- New features and improvements"
            print ""
            next
        }
        1
        ' CHANGELOG.md > "$temp_file"

        mv "$temp_file" CHANGELOG.md
        print_success "CHANGELOG.md updated for version $version"
    fi
}

# Function to update version in code
update_version_in_code() {
    local version=$1

    print_info "Updating version in internal/cli/constants.go..."

    # Update constants.go
    sed -i "" "s/Version   = \".*\"/Version   = \"$version\"/" internal/cli/constants.go

    print_success "Version updated to $version in code"
}

# Function to create release branch
create_release_branch() {
    local version=$1
    local release_branch="release/$version"

    print_info "Creating release branch: $release_branch" >&2

    # Create and checkout release branch
    git checkout -b "$release_branch"

    print_success "Release branch $release_branch created" >&2
    echo "$release_branch"
}

# Function to commit release changes
commit_release_changes() {
    local version=$1

    print_info "Committing release changes..."

    git add CHANGELOG.md internal/cli/constants.go
    git commit -m "release: $version

- Update version to $version
- Update CHANGELOG.md with release notes
- Ready for release process"

    print_success "Release changes committed"
}

# Function to create and push tag
create_and_push_tag() {
    local version=$1
    local release_branch=$2

    print_info "Creating and pushing tag $version..."

    # Extract release notes to temporary file
    local tag_msg_file=$(mktemp)
    echo "Release $version" > "$tag_msg_file"
    echo "" >> "$tag_msg_file"
    awk '/^## \['"${version#v}"'\]/, /^## \[/ {
        if (/^## \['"${version#v}"'\]/) next
        if (/^## \[/ && !/^## \['"${version#v}"'\]/) exit
        print
    }' CHANGELOG.md | sed '/^$/d' >> "$tag_msg_file"

    # Create annotated tag
    git tag -a "$version" -F "$tag_msg_file"
    rm -f "$tag_msg_file"

    # Push branch and tag
    git push origin "$release_branch"
    git push origin "$version"

    print_success "Tag $version created and pushed"
}

# Function to merge back to develop
merge_release() {
    local version=$1
    local release_branch=$2

    print_info "Merging release back to develop..."
    git checkout "$DEVELOP_BRANCH"
    git pull origin "$DEVELOP_BRANCH"
    git merge --no-ff "$release_branch" -m "Merge release $version back to develop"

    # Keep version as-is after release (no auto-bump)
    print_info "Version remains $version for continued development..."

    git push origin "$DEVELOP_BRANCH"

    print_success "Merged back to $DEVELOP_BRANCH"

    # Clean up release branch
    print_info "Cleaning up release branch..."
    git branch -d "$release_branch"
    git push origin --delete "$release_branch"

    print_success "Release branch cleaned up"
}

# Main function
main() {
    local version=$1

    print_info "🚀 Starting TiLoKit Release Process"
    print_info "=================================="

    # Validate input
    if [[ -z "$version" ]]; then
        print_error "Version is required"
        print_info "Usage: $0 <version>"
        print_info "Example: $0 v0.1.0"
        exit 1
    fi

    validate_version "$version"

    # Pre-flight checks
    print_info "Running pre-flight checks..."
    check_clean_working_dir
    check_develop_branch

    # Ensure we're up to date
    print_info "Updating develop branch..."
    git pull origin "$DEVELOP_BRANCH"

    # Update files
    update_changelog "$version"
    update_version_in_code "$version"

    # Create release branch
    local release_branch
    release_branch=$(create_release_branch "$version")

    # Commit changes
    commit_release_changes "$version"

    # Create and push tag
    create_and_push_tag "$version" "$release_branch"

    # Merge release
    merge_release "$version" "$release_branch"

    print_success "🎉 Release $version completed successfully!"
    print_info "Release workflow will now build and deploy automatically."
    print_info "Monitor the progress at: https://github.com/tienld-0801/tilokit/actions"
}

# Run main function
main "$@"
