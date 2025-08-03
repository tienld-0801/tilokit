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

# Function to prepare release directly on develop
prepare_release_on_develop() {
    local version=$1
    
    print_info "Preparing release $version directly on develop branch"
    # Stay on develop branch - no need to create release branch
    return 0
}

# Function to commit release changes
commit_release_changes() {
    local version=$1

    print_info "Committing release changes..."

    git add CHANGELOG.md internal/cli/constants.go
    git commit -m "🚀 release: $version

- Update version to $version
- Update CHANGELOG.md with release notes
- Ready for release process"

    print_success "Release changes committed"
}

# Function to create and push tag
create_and_push_tag() {
    local version=$1

    print_info "Creating tag $version..."

    # Create tag message file
    local tag_msg_file="/tmp/tag_message_$version.txt"
    echo "Release $version" > "$tag_msg_file"
    echo "" >> "$tag_msg_file"
    echo "Changes in this release:" >> "$tag_msg_file"
    echo "" >> "$tag_msg_file"

    # Extract changes for this version from CHANGELOG.md
    awk -v version="$version" '
    BEGIN { in_version = 0; found_version = 0 }
    /^## \[/ {
        if (found_version && in_version) {
            exit
        }
        if ($0 ~ "\\[" substr(version, 2) "\\]") {
            in_version = 1
            found_version = 1
            next
        } else {
            in_version = 0
        }
    }
    in_version && !/^## \[/ && !/^$/ {
        print $0
    }' CHANGELOG.md | sed '/^$/d' >> "$tag_msg_file"

    # Create annotated tag from current commit
    git tag -a "$version" -F "$tag_msg_file"
    rm -f "$tag_msg_file"

    # Push develop branch and tag together
    print_info "Pushing develop branch and tag..."
    git push origin "$DEVELOP_BRANCH"
    git push origin "$version"

    print_success "Tag $version created and pushed with develop branch"
}

# Function to finalize release (no merge needed since we're already on develop)
finalize_release() {
    local version=$1

    print_info "Finalizing release $version..."
    
    # No merge needed - we committed directly to develop
    # Version remains as-is for continued development
    print_info "Version remains $version for continued development..."
    
    print_success "Release finalized on $DEVELOP_BRANCH"
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

    # Prepare release on develop (no branch creation)
    prepare_release_on_develop "$version"

    # Commit changes directly to develop
    commit_release_changes "$version"

    # Create and push tag from develop
    create_and_push_tag "$version"

    # Finalize release (no merge needed)
    finalize_release "$version"

    print_success "🎉 Release $version completed successfully!"
    print_info "Release workflow will now build and deploy automatically."
    print_info "Monitor the progress at: https://github.com/tienld-0801/tilokit/actions"
}

# Run main function
main "$@"
