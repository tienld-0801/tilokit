#!/bin/bash

# TiLoKit CI Commit Message Checker
# Validates all commits in PR follow conventional commit format

# ANSI color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Valid conventional commit types with required emojis
VALID_TYPES=(
    "feat" "fix" "docs" "refactor" "perf" "test"
    "build" "ci" "chore" "style" "revert" "release"
)

# Emoji mapping function for each commit type
get_emoji_for_type() {
    case "$1" in
        "feat") echo "✨" ;;
        "fix") echo "🐛" ;;
        "docs") echo "📚" ;;
        "refactor") echo "♻️" ;;
        "perf") echo "⚡" ;;
        "test") echo "🧪" ;;
        "build") echo "🛠️" ;;
        "ci") echo "🔄" ;;
        "chore") echo "🧹" ;;
        "style") echo "🎨" ;;
        "revert") echo "⏪" ;;
        "release") echo "🚀" ;;
        *) echo "" ;;
    esac
}

# Function to validate commit message
validate_commit() {
    local commit_msg="$1"
    local commit_sha="$2"

    # NO EXCEPTIONS: ALL commits must have emoji format
    # Even merge and revert commits must follow emoji convention

    # Check if commit has emoji format
    local emoji_pattern='^[^[:space:]] [a-z]+: .+'
    local no_emoji_pattern='^[a-z]+: .+'
    
    if [[ $commit_msg =~ $emoji_pattern ]]; then
        # Extract type and emoji from: emoji type: description
        local commit_type=$(echo "$commit_msg" | sed -E 's/^[^[:space:]]+ ([a-z]+): .*/\1/')
        local commit_emoji=$(echo "$commit_msg" | sed -E 's/^([^[:space:]]+) [a-z]+: .*/\1/')
    elif [[ $commit_msg =~ $no_emoji_pattern ]]; then
        echo -e "${RED}❌ Missing emoji: ${commit_sha:0:8}${NC}"
        echo -e "${RED}   Message: $commit_msg${NC}"
        echo -e "${RED}   Required: emoji type: description${NC}"
        echo -e "${RED}   Example: ✨ feat: add new feature${NC}"
        return 1
    else
        echo -e "${RED}❌ Invalid format: ${commit_sha:0:8}${NC}"
        echo -e "${RED}   Message: $commit_msg${NC}"
        echo -e "${RED}   Required: emoji type: description${NC}"
        return 1
    fi

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

    # Check if emoji matches the commit type
    local expected_emoji=$(get_emoji_for_type "$commit_type")
    if [[ "$commit_emoji" != "$expected_emoji" ]]; then
        echo -e "${RED}❌ Wrong emoji for type '$commit_type': ${commit_sha:0:8}${NC}"
        echo -e "${RED}   Message: $commit_msg${NC}"
        echo -e "${RED}   Found emoji: $commit_emoji${NC}"
        echo -e "${RED}   Expected emoji: $expected_emoji${NC}"
        echo -e "${RED}   Correct format: $expected_emoji $commit_type: description${NC}"
        return 1
    fi

    echo -e "${GREEN}✅ Valid: ${commit_sha:0:8} ($expected_emoji $commit_type)${NC}"
    return 0
}

# Main function
main() {
    echo -e "${BLUE}🔍 Checking commit messages for MANDATORY EMOJI format in CI/CD...${NC}"

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
        echo -e "${YELLOW}💡 Please fix commit messages to follow MANDATORY EMOJI format:${NC}"
        echo -e "${YELLOW}   emoji type: description${NC}"
        echo -e "${YELLOW}   Valid types: ${VALID_TYPES[*]}${NC}"
        echo
        echo -e "${YELLOW}📋 Examples with REQUIRED emojis:${NC}"
        echo -e "${GREEN}   ✅ ✨ feat: add user authentication${NC}"
        echo -e "${GREEN}   ✅ 🐛 fix: resolve memory leak${NC}"
        echo -e "${GREEN}   ✅ 📚 docs: update README${NC}"
        echo -e "${GREEN}   ✅ 🧹 chore: upgrade dependencies${NC}"
        echo
        echo -e "${RED}❌ Invalid examples:${NC}"
        echo -e "${RED}   ❌ feat: add feature (missing emoji)${NC}"
        echo -e "${RED}   ❌ 🐛 docs: update (wrong emoji for docs)${NC}"
        exit 1
    else
        echo -e "${GREEN}🎉 All commits follow MANDATORY EMOJI commit format!${NC}"
        exit 0
    fi
}

# Run main function
main "$@"
