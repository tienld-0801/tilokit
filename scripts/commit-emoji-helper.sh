#!/bin/bash

# TiLoKit Commit Emoji Helper
# Utility script to help with commit message emoji management

# ANSI color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Emoji mapping function
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

# Valid commit types
VALID_TYPES=("feat" "fix" "docs" "refactor" "perf" "test" "build" "ci" "chore" "style" "revert" "release")

# Function to show help
show_help() {
    echo -e "${BLUE}🎨 TiLoKit Commit Emoji Helper${NC}"
    echo
    echo -e "${YELLOW}Usage:${NC}"
    echo "  $0 [command] [options]"
    echo
    echo -e "${YELLOW}Commands:${NC}"
    echo "  add <message>           - Add emoji to a commit message"
    echo "  detect <message>        - Detect what emoji would be added"
    echo "  list                    - List all available emojis and types"
    echo "  validate <message>      - Validate a commit message format"
    echo "  interactive             - Interactive mode to create commit message"
    echo
    echo -e "${YELLOW}Examples:${NC}"
    echo "  $0 add \"fix memory leak issue\""
    echo "  $0 detect \"add new authentication feature\""
    echo "  $0 list"
    echo "  $0 validate \"✨ feat: add new feature\""
    echo "  $0 interactive"
}

# Function to detect commit type from message
detect_commit_type() {
    local msg="$1"
    local lower_msg=$(echo "$msg" | tr '[:upper:]' '[:lower:]')
    
    # Check for explicit type prefix (type: message)
    if [[ $msg =~ ^([a-z]+): ]]; then
        echo "${BASH_REMATCH[1]}"
        return
    fi
    
    # Auto-detect based on keywords
    if [[ $lower_msg =~ (add|new|feature|implement|create|introduce) ]]; then
        echo "feat"
    elif [[ $lower_msg =~ (fix|bug|issue|resolve|correct|repair) ]]; then
        echo "fix"
    elif [[ $lower_msg =~ (doc|readme|guide|manual|comment) ]]; then
        echo "docs"
    elif [[ $lower_msg =~ (refactor|restructure|reorganize|clean) ]]; then
        echo "refactor"
    elif [[ $lower_msg =~ (perf|performance|optimize|speed|fast) ]]; then
        echo "perf"
    elif [[ $lower_msg =~ (test|spec|testing|coverage) ]]; then
        echo "test"
    elif [[ $lower_msg =~ (build|compile|make|makefile) ]]; then
        echo "build"
    elif [[ $lower_msg =~ (ci|github|workflow|action|pipeline) ]]; then
        echo "ci"
    elif [[ $lower_msg =~ (style|format|lint|prettier) ]]; then
        echo "style"
    elif [[ $lower_msg =~ (revert|undo|rollback) ]]; then
        echo "revert"
    elif [[ $lower_msg =~ (release|version|v[0-9]) ]]; then
        echo "release"
    elif [[ $lower_msg =~ (chore|maintenance|update|upgrade|deps) ]]; then
        echo "chore"
    else
        echo "chore"  # Default fallback
    fi
}

# Function to add emoji to message
add_emoji() {
    local message="$1"
    
    if [[ -z "$message" ]]; then
        echo -e "${RED}❌ Error: Please provide a commit message${NC}"
        return 1
    fi
    
    # Skip if already has emoji
    if echo "$message" | grep -qE '^[^[:space:]]+ [a-z]+:'; then
        echo -e "${YELLOW}⚠️  Message already has emoji:${NC} $message"
        return 0
    fi
    
    local detected_type=$(detect_commit_type "$message")
    local emoji=$(get_emoji_for_type "$detected_type")
    
    if [[ -n "$emoji" ]]; then
        local new_message="$emoji $detected_type: $message"
        echo -e "${GREEN}✅ Added emoji:${NC} $new_message"
        echo "$new_message"
    else
        echo -e "${RED}❌ Could not detect commit type for: $message${NC}"
        return 1
    fi
}

# Function to detect emoji for message
detect_emoji() {
    local message="$1"
    
    if [[ -z "$message" ]]; then
        echo -e "${RED}❌ Error: Please provide a commit message${NC}"
        return 1
    fi
    
    local detected_type=$(detect_commit_type "$message")
    local emoji=$(get_emoji_for_type "$detected_type")
    
    echo -e "${BLUE}🔍 Analysis for:${NC} \"$message\""
    echo -e "${BLUE}📝 Detected type:${NC} $detected_type"
    echo -e "${BLUE}🎨 Emoji:${NC} $emoji"
    
    if [[ -n "$emoji" ]]; then
        echo -e "${BLUE}📋 Suggested format:${NC} $emoji $detected_type: $message"
    else
        echo -e "${RED}❌ Could not detect commit type${NC}"
    fi
}

# Function to list all emojis
list_emojis() {
    echo -e "${BLUE}🎨 Available Commit Types and Emojis:${NC}"
    echo
    for type in "${VALID_TYPES[@]}"; do
        local emoji="${EMOJI_MAP[$type]}"
        local description=""
        
        case "$type" in
            "feat") description="new features" ;;
            "fix") description="bug fixes" ;;
            "docs") description="documentation changes" ;;
            "refactor") description="code refactoring" ;;
            "perf") description="performance improvements" ;;
            "test") description="adding or updating tests" ;;
            "build") description="build system changes" ;;
            "ci") description="CI/CD changes" ;;
            "chore") description="maintenance tasks" ;;
            "style") description="code style changes" ;;
            "revert") description="reverting previous commits" ;;
            "release") description="version releases" ;;
        esac
        
        echo -e "  ${GREEN}$emoji $type${NC}      - $description"
    done
    echo
    echo -e "${YELLOW}💡 Usage:${NC} $emoji $type: your commit message"
}

# Function to validate commit message
validate_message() {
    local message="$1"
    
    if [[ -z "$message" ]]; then
        echo -e "${RED}❌ Error: Please provide a commit message${NC}"
        return 1
    fi
    
    # Check format: emoji type: description
    if echo "$message" | grep -qE '^[^[:space:]]+ [a-z]+: .+'; then
        local emoji=$(echo "$message" | sed -E 's/^([^[:space:]]+) [a-z]+: .*/\1/')
        local type=$(echo "$message" | sed -E 's/^[^[:space:]]+ ([a-z]+): .*/\1/')
        local expected_emoji=$(get_emoji_for_type "$type")
        
        if [[ "$emoji" == "$expected_emoji" ]]; then
            echo -e "${GREEN}✅ Valid commit message:${NC} $message"
            echo -e "${GREEN}✅ Emoji matches type:${NC} $emoji for $type"
        else
            echo -e "${RED}❌ Invalid emoji for type '$type'${NC}"
            echo -e "${YELLOW}📋 Expected:${NC} $expected_emoji"
            echo -e "${YELLOW}📋 Found:${NC} $emoji"
        fi
    else
        echo -e "${RED}❌ Invalid format. Expected: emoji type: description${NC}"
        echo -e "${YELLOW}📋 Example:${NC} ✨ feat: add new feature"
    fi
}

# Function for interactive mode
interactive_mode() {
    echo -e "${BLUE}🎨 Interactive Commit Message Creator${NC}"
    echo
    echo -e "${YELLOW}Enter your commit message (or press Ctrl+C to cancel):${NC}"
    read -r message
    
    if [[ -z "$message" ]]; then
        echo -e "${RED}❌ Empty message, exiting${NC}"
        return 1
    fi
    
    echo
    detect_emoji "$message"
    echo
    
    local detected_type=$(detect_commit_type "$message")
    local emoji=$(get_emoji_for_type "$detected_type")
    
    if [[ -n "$emoji" ]]; then
        local suggested="$emoji $detected_type: $message"
        echo -e "${YELLOW}Use this message? (y/n):${NC} $suggested"
        read -r confirm
        
        if [[ "$confirm" =~ ^[Yy]$ ]]; then
            echo "$suggested"
        else
            echo -e "${BLUE}ℹ️  You can manually edit the message${NC}"
        fi
    fi
}

# Main script logic
case "$1" in
    "add")
        add_emoji "$2"
        ;;
    "detect")
        detect_emoji "$2"
        ;;
    "list")
        list_emojis
        ;;
    "validate")
        validate_message "$2"
        ;;
    "interactive")
        interactive_mode
        ;;
    "help"|"-h"|"--help"|"")
        show_help
        ;;
    *)
        echo -e "${RED}❌ Unknown command: $1${NC}"
        echo
        show_help
        exit 1
        ;;
esac
