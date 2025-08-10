#!/bin/bash

# TiLoKit Emoji Validation Test Suite
# Tests all commit message formats and emojis to ensure compatibility

# ANSI color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${CYAN}🧪 TiLoKit Emoji Validation Test Suite${NC}"
echo -e "${CYAN}======================================${NC}"
echo

# Test emoji mapping from PR auto-label rules
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
        *) echo "❓" ;;
    esac
}

# All valid commit types
VALID_TYPES=("feat" "fix" "docs" "refactor" "perf" "test" "build" "ci" "chore" "style" "revert" "release")

# Test commit messages
TEST_MESSAGES=(
    "add new authentication feature"
    "resolve memory leak issue"
    "update installation guide"
    "rename module path to tilokit for simplified imports"
    "optimize database queries for better performance"
    "add unit tests for user authentication"
    "update build configuration for production"
    "optimize GitHub Actions workflow"
    "clean up unused dependencies"
    "improve code formatting and style"
    "revert previous breaking changes"
    "version 1.0.0 with new features"
)

# Current regex patterns from commit-msg hook
COMMIT_PATTERN_WITH_EMOJI='^[^[:space:]]+[[:space:]][a-z]+:[[:space:]].+'
COMMIT_PATTERN_WITHOUT_EMOJI='^[a-z]+: .+'

echo -e "${BLUE}📋 Testing All Emoji + Type Combinations${NC}"
echo "================================================"

TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Test each emoji with each type
for type in "${VALID_TYPES[@]}"; do
    emoji=$(get_emoji_for_type "$type")
    test_message="test ${type} functionality"
    commit_msg="${emoji} ${type}: ${test_message}"

    TOTAL_TESTS=$((TOTAL_TESTS + 1))

    echo -n "Testing: ${commit_msg} ... "

    # Test with regex pattern
    if echo "$commit_msg" | grep -qE "$COMMIT_PATTERN_WITH_EMOJI"; then
        echo -e "${GREEN}✅ PASS${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}❌ FAIL${NC}"
        FAILED_TESTS=$((FAILED_TESTS + 1))

        # Debug information
        echo -e "  ${YELLOW}Debug info:${NC}"
        echo -e "  Message: '$commit_msg'"
        echo -e "  Hex dump: $(echo -n "$commit_msg" | hexdump -C | head -1)"
        echo
    fi
done

echo
echo -e "${BLUE}📋 Testing Real-world Commit Messages${NC}"
echo "======================================"

# Test with realistic commit messages
for i in "${!TEST_MESSAGES[@]}"; do
    message="${TEST_MESSAGES[$i]}"
    types=("feat" "fix" "docs" "refactor" "perf" "test" "build" "ci" "chore" "style")
    type="${types[$((i % ${#types[@]}))]}"
    emoji=$(get_emoji_for_type "$type")
    commit_msg="${emoji} ${type}: ${message}"

    TOTAL_TESTS=$((TOTAL_TESTS + 1))

    echo -n "Testing: ${emoji} ${type}: $(echo "$message" | cut -c1-40)... "

    if echo "$commit_msg" | grep -qE "$COMMIT_PATTERN_WITH_EMOJI"; then
        echo -e "${GREEN}✅ PASS${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}❌ FAIL${NC}"
        FAILED_TESTS=$((FAILED_TESTS + 1))

        echo -e "  ${YELLOW}Full message:${NC} '$commit_msg'"
        echo -e "  ${YELLOW}Hex dump:${NC} $(echo -n "$commit_msg" | hexdump -C | head -1)"
        echo
    fi
done

echo
echo -e "${BLUE}📋 Testing Edge Cases${NC}"
echo "===================="

# Test edge cases
EDGE_CASES=(
    "♻️ refactor: rename module path to tilokit for simplified imports"  # Composite emoji
    "📚 docs: update README with current Go version"                      # Single emoji
    "🚀 release: version 1.0.0"                                         # Short message
    "✨ feat: add comprehensive authentication system with JWT tokens and role-based access control" # Long message
    "🐛 fix: resolve critical security vulnerability in user input validation" # Security fix
)

for commit_msg in "${EDGE_CASES[@]}"; do
    TOTAL_TESTS=$((TOTAL_TESTS + 1))

    echo -n "Testing edge case: $(echo "$commit_msg" | cut -c1-50)... "

    if echo "$commit_msg" | grep -qE "$COMMIT_PATTERN_WITH_EMOJI"; then
        echo -e "${GREEN}✅ PASS${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}❌ FAIL${NC}"
        FAILED_TESTS=$((FAILED_TESTS + 1))

        echo -e "  ${YELLOW}Full message:${NC} '$commit_msg'"
        echo -e "  ${YELLOW}Hex dump:${NC} $(echo -n "$commit_msg" | hexdump -C | head -1)"
        echo
    fi
done

echo
echo -e "${BLUE}📋 Testing Invalid Formats (Should Fail)${NC}"
echo "=========================================="

# Test invalid formats that should fail
INVALID_CASES=(
    "feat: missing emoji"                    # No emoji
    "FEAT: uppercase type"                   # Uppercase
    "feat:missing space after colon"        # No space after colon
    "feat(scope): scopes not allowed"       # Scope usage
    " feat: leading whitespace"             # Leading whitespace
    "feat: trailing whitespace "           # Trailing whitespace
)

for commit_msg in "${INVALID_CASES[@]}"; do
    TOTAL_TESTS=$((TOTAL_TESTS + 1))

    echo -n "Testing invalid: '$commit_msg' ... "

    if echo "$commit_msg" | grep -qE "$COMMIT_PATTERN_WITH_EMOJI"; then
        echo -e "${RED}❌ FAIL (Should have been rejected)${NC}"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    else
        echo -e "${GREEN}✅ PASS (Correctly rejected)${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    fi
done

echo
echo -e "${CYAN}📊 Test Results Summary${NC}"
echo "======================="
echo -e "Total Tests: ${TOTAL_TESTS}"
echo -e "Passed: ${GREEN}${PASSED_TESTS}${NC}"
echo -e "Failed: ${RED}${FAILED_TESTS}${NC}"
echo -e "Success Rate: $(( (PASSED_TESTS * 100) / TOTAL_TESTS ))%"

if [ $FAILED_TESTS -eq 0 ]; then
    echo
    echo -e "${GREEN}🎉 All tests passed! Emoji validation is working correctly.${NC}"
    exit 0
else
    echo
    echo -e "${RED}⚠️ Some tests failed. Please review the regex pattern.${NC}"
    exit 1
fi
