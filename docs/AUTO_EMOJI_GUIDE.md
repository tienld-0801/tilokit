# 🎨 Auto Emoji Guide for TiLoKit

This guide explains how to use the automatic emoji functionality in TiLoKit's Git workflow.

## 🚀 Overview

TiLoKit automatically adds appropriate emojis to your commit messages based on the content and type. This ensures consistent, professional-looking commit messages that follow conventional commit standards.

## ✨ Features

- **Automatic Detection**: Analyzes commit message content to determine the appropriate type and emoji
- **Manual Override**: Use explicit type prefixes (e.g., `fix: message`) for precise control
- **Validation**: Ensures emojis match the commit type
- **Helper Tools**: Command-line utilities for emoji management

## 🛠️ How It Works

### 1. Automatic Emoji Addition

When you commit with a message like:

```bash
git commit -m "add new authentication feature"
```

The system automatically converts it to:

```
✨ feat: add new authentication feature
```

### 2. Manual Type Specification

You can also specify the type explicitly:

```bash
git commit -m "fix: resolve memory leak issue"
```

This becomes:

```
🐛 fix: resolve memory leak issue
```

## 📋 Supported Commit Types

| Type       | Emoji | Description      | Keywords                                        |
| ---------- | ----- | ---------------- | ----------------------------------------------- |
| `feat`     | ✨    | New features     | add, new, feature, implement, create, introduce |
| `fix`      | 🐛    | Bug fixes        | fix, bug, issue, resolve, correct, repair       |
| `docs`     | 📚    | Documentation    | doc, readme, guide, manual, comment             |
| `refactor` | ♻️    | Code refactoring | refactor, restructure, reorganize, clean        |
| `perf`     | ⚡    | Performance      | perf, performance, optimize, speed, fast        |
| `test`     | 🧪    | Tests            | test, spec, testing, coverage                   |
| `build`    | 🛠️    | Build system     | build, compile, make, makefile                  |
| `ci`       | 🔄    | CI/CD            | ci, github, workflow, action, pipeline          |
| `style`    | 🎨    | Code style       | style, format, lint, prettier                   |
| `revert`   | ⏪    | Revert           | revert, undo, rollback                          |
| `release`  | 🚀    | Releases         | release, version, v[0-9]                        |
| `chore`    | 🧹    | Maintenance      | chore, maintenance, update, upgrade, deps       |

## 🎯 Usage Examples

### Basic Usage

```bash
# Auto-detected as feat
git commit -m "add user authentication"

# Auto-detected as fix
git commit -m "fix login bug"

# Auto-detected as docs
git commit -m "update installation guide"
```

### Manual Type Specification

```bash
# Explicit type specification
git commit -m "feat: add new dashboard"
git commit -m "fix: resolve memory leak"
git commit -m "docs: update API documentation"
```

### Using the Helper Script

```bash
# List all available types and emojis
git emoji list

# Detect what emoji would be added
git emoji detect "add new feature"

# Add emoji to a message
git emoji add "fix critical bug"

# Validate a commit message
git emoji validate "✨ feat: add new feature"

# Interactive mode
git emoji interactive
```

## 🔧 Configuration

### Git Alias

The helper script is available as a Git alias:

```bash
git emoji [command] [options]
```

### Direct Script Usage

You can also use the script directly:

```bash
./scripts/commit-emoji-helper.sh [command] [options]
```

## 🚫 What Gets Skipped

The auto-emoji system skips:

- Merge commits (automatically generated)
- Squash commits
- Cherry-pick commits
- Empty messages
- Messages that already have emojis

## 🔍 Detection Logic

The system uses keyword matching to detect commit types:

1. **Explicit Type**: If message starts with `type:`, uses that type
2. **Keyword Matching**: Analyzes message content for relevant keywords
3. **Fallback**: Defaults to `chore` if no type can be determined

## ✅ Validation

All commit messages are validated to ensure:

- Proper emoji-type matching
- Correct format: `emoji type: description`
- No leading/trailing whitespace
- Valid commit types

## 🎨 Customization

### Adding New Types

To add new commit types, edit:

1. `.git/hooks/prepare-commit-msg` - Add detection logic
2. `.git/hooks/commit-msg` - Add validation rules
3. `scripts/commit-emoji-helper.sh` - Add helper functions

### Modifying Emojis

Update the `get_emoji_for_type()` function in both hook files and the helper script.

## 🐛 Troubleshooting

### Common Issues

1. **Wrong Emoji Detected**

   - Use explicit type: `git commit -m "fix: your message"`
   - Check keyword usage in your message

2. **Validation Errors**

   - Ensure proper format: `emoji type: description`
   - Check for leading/trailing whitespace

3. **Helper Script Not Working**
   - Ensure script is executable: `chmod +x scripts/commit-emoji-helper.sh`
   - Check Git alias: `git config alias.emoji`

### Debug Mode

Enable debug output by setting:

```bash
export GIT_DEBUG=1
```

## 📚 Best Practices

1. **Be Descriptive**: Write clear, descriptive commit messages
2. **Use Keywords**: Include relevant keywords for better detection
3. **Consistent Format**: Follow the `emoji type: description` format
4. **Review Messages**: Always review auto-generated messages before committing

## 🔗 Related Files

- `.git/hooks/prepare-commit-msg` - Auto-emoji addition
- `.git/hooks/commit-msg` - Message validation
- `scripts/commit-emoji-helper.sh` - Helper utilities
- `.github/workflows/setup-labels.yml` - Label definitions

## 💡 Tips

- Use `git emoji interactive` for guided message creation
- Check `git emoji detect` before committing to see what will be added
- Use `git emoji validate` to check existing messages
- Keep commit messages concise but descriptive

---

**Happy Committing! 🎉**
