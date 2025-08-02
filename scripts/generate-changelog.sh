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
    
    # Initialize arrays
    declare -A categories
    categories[feat]=""
    categories[fix]=""
    categories[perf]=""
    categories[refactor]=""
    categories[docs]=""
    categories[test]=""
    categories[build]=""
    categories[ci]=""
    categories[chore]=""
    categories[style]=""
    categories[revert]=""
    
    # Process each commit
    while IFS= read -r commit; do
        if [[ -z "$commit" ]]; then
            continue
        fi
        
        # Extract type and description  
        local type scope description
        
        # Try to match with scope first: feat(scope): message
        if [[ "$commit" == *"("* ]] && [[ "$commit" =~ ^([a-zA-Z]+)\(([^)]*)\): ]]; then
            type="${BASH_REMATCH[1]}"
            scope="(${BASH_REMATCH[2]})"
            description="${commit#*: }"
        # Then try without scope: feat: message
        elif [[ "$commit" =~ ^([a-zA-Z]+):[[:space:]]*(.+)$ ]]; then
            type="${BASH_REMATCH[1]}"
            scope=""
            description="${BASH_REMATCH[2]}"
        else
            continue
        fi
            
            # Add to appropriate category
            case "$type" in
                feat)
                    categories[feat]+="- ✨ $description\n"
                    ;;
                fix)
                    categories[fix]+="- 🐛 $description\n"
                    ;;
                perf)
                    categories[perf]+="- ⚡ $description\n"
                    ;;
                refactor)
                    categories[refactor]+="- ♻️ $description\n"
                    ;;
                docs)
                    categories[docs]+="- 📚 $description\n"
                    ;;
                test)
                    categories[test]+="- 🧪 $description\n"
                    ;;
                build)
                    categories[build]+="- 🔧 $description\n"
                    ;;
                ci)
                    categories[ci]+="- 🔄 $description\n"
                    ;;
                chore)
                    categories[chore]+="- 🏠 $description\n"
                    ;;
                style)
                    categories[style]+="- 💄 $description\n"
                    ;;
                revert)
                    categories[revert]+="- ⏪ $description\n"
                    ;;
            esac
        else
            # Non-conventional commit
            categories[chore]+="- 🏠 $commit\n"
        fi
    done <<< "$commits"
    
    # Generate changelog sections
    local changelog=""
    
    if [ -n "${categories[feat]}" ]; then
        changelog+="### Added\n${categories[feat]}\n"
    fi
    
    if [ -n "${categories[fix]}" ]; then
        changelog+="### Fixed\n${categories[fix]}\n"
    fi
    
    if [ -n "${categories[perf]}" ]; then
        changelog+="### Performance\n${categories[perf]}\n"
    fi
    
    if [ -n "${categories[refactor]}" ]; then
        changelog+="### Changed\n${categories[refactor]}\n"
    fi
    
    if [ -n "${categories[docs]}" ]; then
        changelog+="### Documentation\n${categories[docs]}\n"
    fi
    
    if [ -n "${categories[test]}" ]; then
        changelog+="### Tests\n${categories[test]}\n"
    fi
    
    if [ -n "${categories[build]}" ]; then
        changelog+="### Build System\n${categories[build]}\n"
    fi
    
    if [ -n "${categories[ci]}" ]; then
        changelog+="### CI/CD\n${categories[ci]}\n"
    fi
    
    if [ -n "${categories[revert]}" ]; then
        changelog+="### Reverted\n${categories[revert]}\n"
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
    
    # Create temporary file
    local temp_file=$(mktemp)
    
    # Add new version entry
    awk -v version="$version" -v date="$date" -v content="$changelog_content" '
    /^## \[Unreleased\]$/ {
        print $0
        print ""
        print "### Changed"
        print "- Development continues..."
        print ""
        print "## [" substr(version, 2) "] - " date
        print ""
        printf "%s", content
        next
    }
    1
    ' CHANGELOG.md > "$temp_file"
    
    # Replace the original file
    mv "$temp_file" CHANGELOG.md
    
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
