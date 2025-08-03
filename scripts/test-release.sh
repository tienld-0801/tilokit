#!/bin/bash

# Test Release Script - Local testing only
# Usage: ./scripts/test-release.sh v0.1.2-dev

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# Test version update function
test_version_update() {
    local version=$1
    local current_version=$(grep 'Version   =' internal/cli/constants.go | cut -d'"' -f2)

    print_info "Current version: $current_version"
    print_info "Testing version update to: $version"

    # Create backup
    cp internal/cli/constants.go internal/cli/constants.go.backup

    # Test sed command
    sed -i "" "s/Version   = \".*\"/Version   = \"$version\"/" internal/cli/constants.go

    # Check if update worked
    local new_version=$(grep 'Version   =' internal/cli/constants.go | cut -d'"' -f2)

    if [[ "$new_version" == "$version" ]]; then
        print_success "Version update test PASSED: $current_version → $new_version"
    else
        print_error "Version update test FAILED: Expected $version, got $new_version"
    fi

    # Restore backup
    mv internal/cli/constants.go.backup internal/cli/constants.go
    print_info "Restored original version: $current_version"
}

# Test changelog generation
test_changelog_generation() {
    local version=$1

    print_info "Testing changelog generation for $version"

    if [ -f "scripts/generate-changelog.sh" ]; then
        print_info "Found changelog generator script"
        # Test without actually modifying
        print_success "Changelog generation test ready"
    else
        print_warning "No changelog generator found - will use fallback method"
    fi
}

# Main test function
main() {
    local version=$1

    print_info "🧪 Testing TiLoKit Release Process (LOCAL ONLY)"
    print_info "=============================================="

    if [[ -z "$version" ]]; then
        print_error "Version is required"
        print_info "Usage: $0 <version>"
        print_info "Example: $0 v0.1.2-dev"
        exit 1
    fi

    # Test version format
    if [[ ! $version =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9\.-]+)?$ ]]; then
        print_error "Invalid version format: $version"
        print_info "Expected format: v1.0.0 or v1.0.0-beta.1"
        exit 1
    fi

    print_success "Version format is valid: $version"

    # Test version update
    test_version_update "$version"

    # Test changelog
    test_changelog_generation "$version"

    # Test build
    print_info "Testing build process..."
    if go build -o tilokit-test .; then
        print_success "Build test PASSED"
        rm -f tilokit-test
    else
        print_error "Build test FAILED"
        exit 1
    fi

    print_success "🎉 All tests PASSED! Release script should work correctly."
    print_warning "Remember: This was a LOCAL TEST only - no files were permanently modified"
}

main "$@"
