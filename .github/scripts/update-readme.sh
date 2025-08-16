#!/bin/bash

# TiLoKit README Update Script
# Updates version badge and contributors in README.md during release
# Usage: ./update-readme.sh [version]

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

# Get version parameter
VERSION=$1
if [[ -z "$VERSION" ]]; then
    print_error "Version parameter is required"
    print_info "Usage: ./update-readme.sh v0.1.0"
    exit 1
fi

README_FILE=".github/README.md"

# Check if README.md exists
if [[ ! -f "$README_FILE" ]]; then
    print_error "README.md not found"
    exit 1
fi

print_info "Updating README.md for version $VERSION..."

# 1. Update version badge in README.md (line 9)
print_info "Updating version badge..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS - escape special characters in version
    ESCAPED_VERSION=$(echo "$VERSION" | sed 's/[[\.\*^$(){}?+|]/\\&/g')
    sed -i '' "s|release-v[^-]*--[^-]*-green|release-${ESCAPED_VERSION}--green|g" "$README_FILE"
else
    # Linux - escape special characters in version
    ESCAPED_VERSION=$(echo "$VERSION" | sed 's/[[\.\*^$(){}?+|]/\\&/g')
    sed -i "s|release-v[^-]*--[^-]*-green|release-${ESCAPED_VERSION}--green|g" "$README_FILE"
fi

# 2. Get contributors from GitHub API
print_info "Fetching contributors from GitHub API..."

# Call GitHub API to get contributors
API_RESPONSE=$(curl -s --location 'https://api.github.com/repos/tienld-0801/tilokit/contributors')

if [[ $? -ne 0 ]] || [[ -z "$API_RESPONSE" ]]; then
    print_error "Failed to fetch contributors from GitHub API"
    print_warning "Falling back to default contributors"
    ALL_USERS="tienld-0801"
else
    print_success "Successfully fetched contributors from GitHub API"

    # Parse JSON response to extract login names, excluding bots
    if command -v jq &> /dev/null; then
        # Use jq if available (preferred method) - only get Users, not Bots
        ALL_USERS=$(echo "$API_RESPONSE" | jq -r '.[] | select(.type == "User") | .login' | head -10)
        print_info "Using jq to parse JSON response"
    else
        # Fallback: use grep and sed to parse JSON (less reliable but works)
        print_info "jq not found, using fallback JSON parsing"

        # Extract entries that are Users (not Bots)
        ALL_USERS=$(echo "$API_RESPONSE" | \
            grep -B5 -A5 '"type": "User"' | \
            grep '"login":' | \
            sed 's/.*"login": *"//' | \
            sed 's/".*//' | \
            head -10)
    fi

    # Remove any empty lines
    ALL_USERS=$(echo "$ALL_USERS" | grep -v '^$')

    if [[ -z "$ALL_USERS" ]]; then
        print_warning "No valid contributors found in API response"
        ALL_USERS="tienld-0801"
    fi
fi

print_info "Found GitHub contributors (Users only, Bots excluded):"
echo "$ALL_USERS"

# Debug: Show total count
CONTRIBUTOR_COUNT=$(echo "$ALL_USERS" | wc -l)
print_info "Total contributors found: $CONTRIBUTOR_COUNT"

# 3. Check existing contributors in README to avoid duplicates
print_info "Checking existing contributors in README..."

EXISTING_CONTRIBUTORS=""
if [[ -f "$README_FILE" ]]; then
    # Extract existing GitHub usernames from contributor avatar links only
    # Look for pattern: <a href="https://github.com/username"> followed by <img with .png
    EXISTING_CONTRIBUTORS=$(grep -A1 '<a href="https://github.com/[^/"]*">' "$README_FILE" | \
        grep -B1 'img src="https://github.com/' | \
        grep '<a href=' | \
        sed 's|.*<a href="https://github.com/||' | \
        sed 's|">.*||' | \
        sort | uniq)

    if [[ -n "$EXISTING_CONTRIBUTORS" ]]; then
        print_info "Found existing contributors in README:"
        echo "$EXISTING_CONTRIBUTORS"
    else
        print_info "No existing contributors found in README"
    fi
fi

# 4. Generate contributors section with avatars (left-aligned)
print_info "Generating contributors section..."

# Filter out existing contributors from API response
NEW_CONTRIBUTORS=""
if [[ -n "$ALL_USERS" ]]; then
    while IFS= read -r github_user; do
        if [[ -n "$github_user" ]]; then
            # Check if user already exists in README
            if echo "$EXISTING_CONTRIBUTORS" | grep -q "^${github_user}$"; then
                print_info "Skipping existing contributor: $github_user"
            else
                NEW_CONTRIBUTORS+="$github_user\n"
            fi
        fi
    done <<< "$ALL_USERS"
fi

# Remove trailing newline and convert to proper format
NEW_CONTRIBUTORS=$(echo -e "$NEW_CONTRIBUTORS" | grep -v '^$' || true)

# Add GitHub users with avatars (left-aligned, no center div)
if [[ -n "$NEW_CONTRIBUTORS" ]]; then
    CONTRIBUTOR_COUNT=0
    while IFS= read -r github_user; do
        if [[ -n "$github_user" ]] && [[ $CONTRIBUTOR_COUNT -lt 8 ]]; then
            # Add spacing before first contributor
            if [[ $CONTRIBUTOR_COUNT -eq 0 ]]; then
                CONTRIBUTORS_SECTION+="

<a href=\"https://github.com/${github_user}\">
  <img src=\"https://github.com/${github_user}.png?size=50&mask=circle\" width=\"50\" height=\"50\" alt=\"${github_user}\"/>
</a>"
            else
                CONTRIBUTORS_SECTION+="
<a href=\"https://github.com/${github_user}\">
  <img src=\"https://github.com/${github_user}.png?size=50&mask=circle\" width=\"50\" height=\"50\" alt=\"${github_user}\"/>
</a>"
            fi
            ((CONTRIBUTOR_COUNT++))
        fi
    done <<< "$NEW_CONTRIBUTORS"

    print_success "Added $CONTRIBUTOR_COUNT new contributors to README"
else
    print_warning "No new contributors to add (all already exist in README)"
    # Set empty contributors section to avoid issues in processing
    CONTRIBUTORS_SECTION=""
fi

# 5. Find and replace the contributors section in README.md
# Add contributors section right after the Contributing header
print_info "Updating contributors section..."

# Create a temporary file with the updated content
TEMP_FILE=$(mktemp)

# Flags to track state
IN_CONTRIBUTING_SECTION=false
CONTRIBUTORS_ADDED=false
SKIP_UNTIL_NEXT_SECTION=false

while IFS= read -r line; do
    # Check if we've reached the Contributing section
    if [[ "$line" == "## 🤝 **Contributing**" ]]; then
        echo "$line" >> "$TEMP_FILE"
        # Add contributors section right after Contributing header only if there are new contributors
        if [[ -n "$CONTRIBUTORS_SECTION" ]]; then
            echo "$CONTRIBUTORS_SECTION" >> "$TEMP_FILE"
        fi
        CONTRIBUTORS_ADDED=true
        IN_CONTRIBUTING_SECTION=true
        SKIP_UNTIL_NEXT_SECTION=false
    # Check if we hit an existing Contributors subsection - skip it entirely
    elif [[ "$line" =~ ^###.*Contributors.*$ ]] && [[ "$IN_CONTRIBUTING_SECTION" == true ]]; then
        # Start skipping until we hit the next ### section
        SKIP_UNTIL_NEXT_SECTION=true
        continue
    # Check if we hit another ### section - stop skipping
    elif [[ "$line" =~ ^###.*$ ]] && [[ "$SKIP_UNTIL_NEXT_SECTION" == true ]]; then
        # We've reached the next subsection, stop skipping
        SKIP_UNTIL_NEXT_SECTION=false
        echo "$line" >> "$TEMP_FILE"
    # Check if we hit another ## section - exit Contributing section
    elif [[ "$line" =~ ^##.*$ ]] && [[ "$line" != "## 🤝 **Contributing**" ]]; then
        IN_CONTRIBUTING_SECTION=false
        SKIP_UNTIL_NEXT_SECTION=false
        echo "$line" >> "$TEMP_FILE"
    # Skip lines if we're in skip mode
    elif [[ "$SKIP_UNTIL_NEXT_SECTION" == true ]]; then
        continue
    # Normal line - add it
    else
        echo "$line" >> "$TEMP_FILE"
    fi
done < "$README_FILE"

# Replace the original file
mv "$TEMP_FILE" "$README_FILE"

# 5. Update last modified date in acknowledgments
print_info "Updating last modified date..."
CURRENT_DATE=$(date +"%Y-%m-%d")
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    sed -i '' "s|Updated: [0-9-]*|Updated: $CURRENT_DATE|g" "$README_FILE"
else
    # Linux
    sed -i "s|Updated: [0-9-]*|Updated: $CURRENT_DATE|g" "$README_FILE"
fi

print_success "README.md updated successfully!"
print_info "Changes made:"
print_info "  - Version badge updated to $VERSION"
print_info "  - Contributors section updated with $(echo "$CONTRIBUTORS" | wc -l) contributors"
print_info "  - Last modified date updated to $CURRENT_DATE"

# Show the diff for verification
print_info "Verifying changes..."
if command -v git &> /dev/null; then
    git diff --no-index /dev/null "$README_FILE" | head -20 || true
fi

print_success "README.md update completed!"
