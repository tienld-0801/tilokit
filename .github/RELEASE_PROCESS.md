# TiLoKit Release Process

This document outlines the release process for TiLoKit, following GitFlow methodology with automated CI/CD.

## 🌊 Branch Strategy

### Main Branches
- **`main`**: Production-ready code, always stable
- **`develop`**: Integration branch for features, active development

### Supporting Branches
- **`feature/*`**: New features (branch from `develop`, merge back to `develop`)
- **`release/*`**: Prepare releases (branch from `develop`, merge to `main` and `develop`)
- **`hotfix/*`**: Critical fixes (branch from `main`, merge to `main` and `develop`)

## 🚀 Release Types

### 1. Regular Release
For new features and improvements:

```bash
# Ensure you're on develop branch
git checkout develop
git pull origin develop

# Start release process
./scripts/release.sh v0.2.0
```

### 2. Hotfix Release
For critical bug fixes:

```bash
# Ensure you're on main branch
git checkout main
git pull origin main

# Start hotfix process
./scripts/hotfix.sh v0.1.1
```

### 3. Development Releases
For testing and pre-release versions:

```bash
# Manual tagging for alpha/beta/rc versions
git tag v0.2.0-alpha.1
git push origin v0.2.0-alpha.1
```

## 📋 Version Numbering

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR** (v1.0.0): Breaking changes
- **MINOR** (v0.1.0): New features, backwards compatible
- **PATCH** (v0.0.1): Bug fixes, backwards compatible

### Development Versions
- **`v0.1.0-dev`**: Development version
- **`v0.2.0-alpha.1`**: Alpha release
- **`v0.2.0-beta.1`**: Beta release
- **`v0.2.0-rc.1`**: Release candidate

## 🔄 Release Workflow

### Automated Process
1. **Tag Creation**: Push a version tag (e.g., `v0.1.0`)
2. **CI/CD Trigger**: GitHub Actions automatically starts
3. **Build**: Multi-platform binaries are built
4. **Test**: Full test suite runs
5. **Deploy**: Binaries deployed to GitHub Pages
6. **Release**: GitHub Release created with assets

### Manual Steps (via Script)
1. **Preparation**: Run `./scripts/release.sh v0.1.0`
2. **Branch Management**: Script handles branch creation/merging
3. **Changelog**: Automatically updated
4. **Version Bump**: Code version updated
5. **Tag Creation**: Annotated tag created and pushed

## 📝 Changelog Management

### Format
We follow [Keep a Changelog](https://keepachangelog.com/) format:

```markdown
## [1.0.0] - 2024-01-15

### Added
- New feature descriptions

### Changed
- Modified functionality

### Deprecated
- Features to be removed

### Removed
- Deleted features

### Fixed
- Bug fixes

### Security
- Security improvements
```

### Automation
- Release script automatically updates CHANGELOG.md
- Extracts release notes from changelog for GitHub releases
- Maintains proper versioning and dates

## 🛠️ Scripts

### Release Script (`scripts/release.sh`)
- Creates release branch from develop
- Updates version and changelog
- Creates and pushes annotated tag
- Merges to main and back to develop
- Cleans up release branch

### Hotfix Script (`scripts/hotfix.sh`)
- Creates hotfix branch from main
- Allows manual fixes
- Updates version and changelog
- Creates and pushes tag
- Merges to main and develop

## ✅ Pre-release Checklist

- [ ] All tests passing
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Version bumped appropriately
- [ ] Breaking changes documented
- [ ] Migration guide (if needed)

## 🔐 Security

- Signed tags for releases
- Checksums for all binaries
- Security scan in CI/CD
- Dependency vulnerability checks

## 📊 Monitoring

- GitHub Actions for build status
- Release metrics tracking
- Download statistics
- Issue/PR linking to releases

## 🆘 Rollback Process

### If Release Fails
1. **Stop**: Cancel ongoing deployments
2. **Revert**: Remove problematic tag
3. **Fix**: Address issues on release branch
4. **Retry**: Create new patch version

### Emergency Rollback
1. **Hotfix**: Create immediate hotfix
2. **Communication**: Notify users of issues
3. **Documentation**: Update known issues

## 🔄 Regular Maintenance

- **Weekly**: Review and merge approved PRs
- **Monthly**: Dependency updates
- **Quarterly**: Security audit
- **Annually**: Process review and improvements

## 📞 Support

For questions about the release process:
- Create issue with `release` label
- Contact maintainers
- Check workflow logs for CI/CD issues

---
*Last updated: 2024-01-31*
