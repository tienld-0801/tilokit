# TiLoKit Git Hooks

This directory contains Git hooks for TiLoKit that enforce commit message standards and automatically add emojis.

## Hooks Included

### `commit-msg`

- **Purpose**: Validates commit messages
- **Requirements**:
  - Must have valid commit type (feat, fix, docs, etc.)
  - Must have emoji (auto-added by prepare-commit-msg)
- **Rejects**: Commits without proper format

### `prepare-commit-msg`

- **Purpose**: Auto-adds emojis to commit messages
- **Behavior**:
  - Only works with valid commit types
  - Automatically prepends appropriate emoji
  - Rejects invalid commit types

## Installation

Run the installation script:

```bash
./scripts/install-hooks.sh
```

## Valid Commit Types

| Type       | Emoji | Description      |
| ---------- | ----- | ---------------- |
| `feat`     | ✨    | New features     |
| `fix`      | 🐛    | Bug fixes        |
| `docs`     | 📚    | Documentation    |
| `refactor` | ♻️    | Code refactoring |
| `perf`     | ⚡    | Performance      |
| `test`     | 🧪    | Tests            |
| `build`    | 🛠️    | Build system     |
| `ci`       | 🔄    | CI/CD            |
| `chore`    | 🧹    | Maintenance      |
| `style`    | 🎨    | Code style       |
| `revert`   | ⏪    | Revert           |
| `release`  | 🚀    | Release          |

## Examples

### ✅ Valid Commits

```bash
git commit -m "feat: add new feature"     # → ✨ feat: add new feature
git commit -m "fix: resolve bug"          # → 🐛 fix: resolve bug
git commit -m "test: add unit tests"      # → 🧪 test: add unit tests
```

### ❌ Invalid Commits

```bash
git commit -m "add new feature"           # → REJECTED (no tag)
git commit -m "invalid: test something"   # → REJECTED (invalid tag)
git commit -m "FEAT: add feature"         # → REJECTED (uppercase)
```

## Manual Installation

If you prefer to install manually:

```bash
# Copy hooks to git directory
cp .githooks/commit-msg .git/hooks/
cp .githooks/prepare-commit-msg .git/hooks/

# Make them executable
chmod +x .git/hooks/commit-msg
chmod +x .git/hooks/prepare-commit-msg
```

## Troubleshooting

- **Hooks not working**: Ensure they're executable (`chmod +x`)
- **Permission denied**: Check file permissions
- **Hooks not found**: Run `./scripts/install-hooks.sh`
