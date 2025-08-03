#!/bin/bash

# TiLoKit CI Commit Message Checker
# Validates all commits in PR follow conventional commit format

# ANSI color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Valid conventional commit types
VALID_TYPES=(
    "feat" "fix" "docs" "refactor" "perf" "test"
    "build" "ci" "chore" "style" "revert"
)

# Function to validate commit message
validate_commit() {
    local commit_msg="$1"
    local commit_sha="$2"

    # Skip merge commits
    if [[ $commit_msg =~ ^Merge\ .*|^Revert\ .* ]]; then
        echo -e "${BLUE}ℹ️  Skipping merge/revert commit: ${commit_sha:0:8}${NC}"
        return 0
    fi

    # Check conventional commit format
    local commit_pattern='^[a-z]+(\([^)]+\))?:[ ].+'
    if [[ ! $commit_msg =~ $commit_pattern ]]; then
        echo -e "${RED}❌ Invalid format: ${commit_sha:0:8}${NC}"
        echo -e "${RED}   Message: $commit_msg${NC}"
        echo -e "${RED}   Required: type(scope): description${NC}"
        return 1
    fi

    # Extract type
    local commit_type=$(echo "$commit_msg" | sed -E 's/^([a-z]+)(\([^)]*\))?: .*/\1/')

    # Check if type is valid
    local valid_type=false
    for type in "${VALID_TYPES[@]}"; do
        if [[ "$commit_type" == "$type" ]]; then
            valid_type=true
            break
        fi
    done

    if [[ "$valid_type" == false ]]; then
        echo -e "${RED}❌ Invalid type '$commit_type': ${commit_sha:0:8}${NC}"
        echo -e "${RED}   Message: $commit_msg${NC}"
        echo -e "${RED}   Valid types: ${VALID_TYPES[*]}${NC}"
        return 1
    fi

    echo -e "${GREEN}✅ Valid: ${commit_sha:0:8} ($commit_type)${NC}"
    return 0
}

# Main function
main() {
    echo -e "${BLUE}🔍 Checking commit messages in CI/CD...${NC}"

    # Get the range of commits to check
    local base_ref="${GITHUB_BASE_REF:-main}"
    local head_ref="${GITHUB_HEAD_REF:-HEAD}"

    echo -e "${BLUE}📋 Checking commits from $base_ref to $head_ref${NC}"

    # Get commits in the current branch/PR
    local commits
    if [[ -n "$GITHUB_BASE_REF" ]]; then
        # In PR context
        commits=$(git log --format="%H %s" origin/$base_ref..HEAD)
    else
        # Local context - check last 10 commits
        commits=$(git log --format="%H %s" -10)
    fi

    if [[ -z "$commits" ]]; then
        echo -e "${YELLOW}⚠️  No commits found to check${NC}"
        exit 0
    fi

    local total_commits=0
    local invalid_commits=0

    # Check each commit
    while IFS= read -r line; do
        [[ -z "$line" ]] && continue

        local commit_sha=$(echo "$line" | cut -d' ' -f1)
        local commit_msg=$(echo "$line" | cut -d' ' -f2-)

        ((total_commits++))

        if ! validate_commit "$commit_msg" "$commit_sha"; then
            ((invalid_commits++))
        fi

    done <<< "$commits"

    # Summary
    echo
    echo -e "${BLUE}======================================${NC}"
    echo -e "${BLUE}📊 Commit Validation Summary${NC}"
    echo -e "${BLUE}======================================${NC}"
    echo -e "${BLUE}Total commits checked: $total_commits${NC}"
    echo -e "${GREEN}Valid commits: $((total_commits - invalid_commits))${NC}"
    echo -e "${RED}Invalid commits: $invalid_commits${NC}"

    if [[ $invalid_commits -gt 0 ]]; then
        echo
        echo -e "${RED}❌ Commit validation failed!${NC}"
        echo -e "${YELLOW}💡 Please fix commit messages to follow conventional commits:${NC}"
        echo -e "${YELLOW}   type(scope): description${NC}"
        echo -e "${YELLOW}   Valid types: ${VALID_TYPES[*]}${NC}"
        echo
        echo -e "${YELLOW}📋 Examples:${NC}"
        echo -e "${GREEN}   ✅ feat: add user authentication${NC}"
        echo -e "${GREEN}   ✅ fix(core): resolve memory leak${NC}"
        echo -e "${GREEN}   ✅ docs: update README${NC}"
        exit 1
    else
        echo -e "${GREEN}🎉 All commits follow conventional commit format!${NC}"
        exit 0
    fi
}

# Run main function
main "$@"
