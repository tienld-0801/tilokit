#!/bin/bash

# TiLoKit Changelog Generator
# Generates beautiful changelog from conventional commits like gofiber/nuxt

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

# Function to get commits since last tag
get_commits_since_last_tag() {
    local last_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
    if [ -z "$last_tag" ]; then
        git log --oneline --pretty=format:"%s" --reverse
    else
        git log ${last_tag}..HEAD --oneline --pretty=format:"%s" --reverse
    fi
}

# Function to categorize commits
categorize_commits() {
    local commits="$1"
    
    # Initialize category variables
    local feat_items fix_items perf_items refactor_items docs_items
    local test_items build_items ci_items chore_items style_items revert_items
    feat_items=""
    fix_items=""
    perf_items=""
    refactor_items=""
    docs_items=""
    test_items=""
    build_items=""
    ci_items=""
    chore_items=""
    style_items=""
    revert_items=""
    
    # Process each commit
    while IFS= read -r commit; do
        if [[ -z "$commit" ]]; then
            continue
        fi
        
        # Extract type and description using sed
        local type scope description
        
        # Check if commit has scope: feat(scope): message
        if [[ "$commit" == *"("* ]] && [[ "$commit" == *"): "* ]]; then
            type=$(echo "$commit" | sed 's/^\([a-zA-Z]*\)(.*/\1/')
            scope=$(echo "$commit" | sed 's/^[a-zA-Z]*(//' | sed 's/).*//')
            description=$(echo "$commit" | sed 's/^[^:]*: *//')
            scope="($scope)"
        # Check if commit is simple: feat: message
        elif [[ "$commit" == *": "* ]]; then
            type=$(echo "$commit" | sed 's/^\([a-zA-Z]*\):.*/\1/')
            scope=""
            description=$(echo "$commit" | sed 's/^[^:]*: *//')
        else
            continue
        fi
            
            # Add to appropriate category
            case "$type" in
                feat)
                    feat_items="${feat_items}- ✨ $description\n"
                    ;;
                fix)
                    fix_items="${fix_items}- 🐛 $description\n"
                    ;;
                perf)
                    perf_items="${perf_items}- ⚡ $description\n"
                    ;;
                refactor)
                    refactor_items="${refactor_items}- ♻️ $description\n"
                    ;;
                docs)
                    docs_items="${docs_items}- 📚 $description\n"
                    ;;
                test)
                    test_items="${test_items}- 🧪 $description\n"
                    ;;
                build)
                    build_items="${build_items}- 🔧 $description\n"
                    ;;
                ci)
                    ci_items="${ci_items}- 🔄 $description\n"
                    ;;
                chore)
                    # Skip chore commits (maintenance tasks) from release notes
                    continue
                    ;;
                style)
                    style_items="${style_items}- 💄 $description\n"
                    ;;
                revert)
                    revert_items="${revert_items}- ⏪ $description\n"
                    ;;
            esac
    done <<< "$commits"
    
    # Generate changelog sections
    local changelog=""
    
    if [ -n "$feat_items" ]; then
        changelog+="### Added
$feat_items
"
    fi
    
    if [ -n "$fix_items" ]; then
        changelog+="### Fixed
$fix_items
"
    fi
    
    if [ -n "$perf_items" ]; then
        changelog+="### Performance
$perf_items
"
    fi
    
    if [ -n "$refactor_items" ]; then
        changelog+="### Changed
$refactor_items
"
    fi
    
    if [ -n "$docs_items" ]; then
        changelog+="### Documentation
$docs_items
"
    fi
    
    if [ -n "$test_items" ]; then
        changelog+="### Tests
$test_items
"
    fi
    
    if [ -n "$build_items" ]; then
        changelog+="### Build System
$build_items
"
    fi
    
    if [ -n "$ci_items" ]; then
        changelog+="### CI/CD
$ci_items
"
    fi
    
    if [ -n "$style_items" ]; then
        changelog+="### Style
$style_items
"
    fi
    
    # Maintenance section removed - chore commits are excluded
    
    if [ -n "$revert_items" ]; then
        changelog+="### Reverted
$revert_items
"
    fi
    
    echo -e "$changelog"
}

# Function to update CHANGELOG.md
update_changelog() {
    local version=$1
    local date=$(date +%Y-%m-%d)
    local commits=$(get_commits_since_last_tag)
    
    if [ -z "$commits" ]; then
        print_warning "No commits found since last tag"
        return 0
    fi
    
    print_info "Generating changelog for version $version..."
    
    # Generate changelog content
    local changelog_content=$(categorize_commits "$commits")
    
    if [ -z "$changelog_content" ]; then
        print_warning "No conventional commits found"
        changelog_content="### Changed\n- Various improvements and bug fixes\n"
    fi
    
    # Create temporary files
    local temp_file=$(mktemp)
    local content_file=$(mktemp)
    
    # Write changelog content to temporary file
    echo "$changelog_content" > "$content_file"
    
    # Add new version entry
    awk -v version="$version" -v date="$date" -v content_file="$content_file" '
    /^## \[Unreleased\]$/ {
        print $0
        print ""
        print "### Changed"
        print "- Development continues..."
        print ""
        print "## [" substr(version, 2) "] - " date
        print ""
        # Read content from file
        while ((getline line < content_file) > 0) {
            print line
        }
        close(content_file)
        next
    }
    1
    ' CHANGELOG.md > "$temp_file"
    
    # Replace the original file
    mv "$temp_file" CHANGELOG.md
    
    # Cleanup temporary file
    rm -f "$content_file"
    
    # Clean up trailing spaces to prevent markdownlint issues
    if command -v sed >/dev/null 2>&1; then
        if [[ "$OSTYPE" == "darwin"* ]]; then
            # macOS sed
            sed -i '' 's/[[:space:]]*$//' CHANGELOG.md
        else
            # GNU sed
            sed -i 's/[[:space:]]*$//' CHANGELOG.md
        fi
        print_info "Cleaned trailing spaces from CHANGELOG.md"
    fi
    
    print_success "CHANGELOG.md updated for version $version"
}

# Main function
main() {
    local version=$1
    
    print_info "🔄 Generating Changelog"
    print_info "======================="
    
    if [ -z "$version" ]; then
        print_error "Version is required"
        print_info "Usage: $0 <version>"
        print_info "Example: $0 v0.1.1"
        exit 1
    fi
    
    update_changelog "$version"
    
    print_success "🎉 Changelog generation completed!"
}

# Run main function
main "$@"
