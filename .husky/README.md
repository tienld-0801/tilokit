# 🐕 TiLoKit Git Hooks (.husky style)

This directory contains Git hooks for TiLoKit project, organized in a `.husky`-style structure (inspired by the popular JavaScript tool).

## 📋 Available Hooks

### `commit-msg` 
- **Purpose**: Validates commit messages follow conventional commit format
- **Format**: `type(scope): description`
- **Valid Types**: `feat`, `fix`, `docs`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `style`, `revert`

### `pre-commit`
- **Purpose**: Runs pre-commit checks before allowing commits
- **Checks**: 
  - Go syntax validation
  - Merge conflict markers
  - Large files detection (>10MB)
  - Sensitive files detection

## 🚀 Quick Setup

```bash
# Install hooks using Makefile
make install-hooks

# Or run installer directly
./.husky/hooks/install-hooks.sh
```

## 📝 Commit Message Examples

### ✅ Valid Commits
```bash
git commit -m "feat: add user authentication"
git commit -m "fix(core): resolve memory leak issue"
git commit -m "docs: update API documentation"
git commit -m "chore(deps): upgrade dependencies"
git commit -m "refactor(auth): simplify login logic"
```

### ❌ Invalid Commits
```bash
git commit -m "Add new feature"          # Missing type
git commit -m "FEAT: add feature"        # Uppercase type
git commit -m "feat:add feature"         # Missing space after colon
git commit -m "feature: add new thing"   # Invalid type
```

## 🔧 Manual Installation

If you prefer to install hooks manually:

```bash
# Copy hooks to .git/hooks/
cp .husky/hooks/commit-msg .git/hooks/
cp .husky/hooks/pre-commit .git/hooks/

# Make them executable
chmod +x .git/hooks/commit-msg
chmod +x .git/hooks/pre-commit
```

## 🗑️ Uninstall Hooks

```bash
# Using Makefile
make uninstall-hooks

# Or manually
rm .git/hooks/commit-msg
rm .git/hooks/pre-commit
```

## 🔍 How It Works

1. **Installation**: Hooks are copied from `.husky/hooks/` to `.git/hooks/`
2. **Execution**: Git automatically runs these hooks during commit process
3. **Validation**: If hooks fail, the commit is rejected with helpful error messages

## 📚 Conventional Commit Types

| Type | Emoji | Description | Example |
|------|-------|-------------|---------|
| `feat` | ✨ | New features | `feat: add OAuth integration` |
| `fix` | 🐛 | Bug fixes | `fix: resolve login timeout` |
| `docs` | 📚 | Documentation | `docs: update README` |
| `refactor` | ♻️ | Code refactoring | `refactor: optimize database queries` |
| `perf` | ⚡ | Performance improvements | `perf: cache user sessions` |
| `test` | 🧪 | Tests | `test: add unit tests for auth` |
| `build` | 🛠️ | Build system | `build: update Dockerfile` |
| `ci` | 🔄 | CI/CD | `ci: add GitHub Actions` |
| `chore` | 🧹 | Maintenance | `chore: update dependencies` |
| `style` | 🎨 | Code style | `style: format Go code` |
| `revert` | ⏪ | Reverts | `revert: undo previous commit` |

## 🎯 Integration with GitHub

These commit message formats integrate with:
- **.github/workflows/pr-auto-label.yml**: Auto-adds labels to PRs
- **Release notes generation**: Groups commits by type
- **Changelog automation**: Creates organized changelogs

## 💡 Tips

- Use `git commit --amend` to fix commit messages
- Hooks apply to all contributors when they install them
- Merge commits and reverts are automatically skipped
- You can temporarily bypass hooks with `git commit --no-verify` (not recommended)

## 🔗 Related Files

- `.github/workflows/pr-auto-label.yml` - Auto PR labeling
- `scripts/generate-changelog.sh` - Changelog generation
- `Makefile` - Build and hook management commands
