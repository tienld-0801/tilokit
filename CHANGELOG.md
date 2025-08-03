# Changelog

All notable changes to TiLoKit will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.1.0-dev] - 2025-08-03

### 🚀 New Features

#### Core CLI Framework
- ✨ **Project scaffolding system** with multiple framework support
- ✨ **Template engine** for React, Vue, and custom templates
- ✨ **Interactive CLI** with beautiful project creation wizard
- ✨ **Auto-update functionality** - Self-updating from GitHub releases
- ✨ **Global `--update` flag** for quick updates
- ✨ **Cross-platform support** - macOS, Linux, and Windows
- ✨ **Smart binary detection** with automatic platform identification
- ✨ **Release notes display** during updates

#### Template System
- 📦 **React templates** with modern build tools (Vite, Create React App)
- 📦 **Vue templates** with complete component structure
- 📦 **Custom template support** with flexible configuration
- 📦 **Template validation** and structure verification
- 📦 **Package management** with npm/yarn support
- 📦 **TypeScript configurations** for modern development

#### Development Tools
- 🔧 **Professional release system** with GitFlow integration
- 🔧 **Automated changelog generation** from conventional commits
- 🔧 **Multi-platform builds** with Apple Silicon optimizations
- 🔧 **GitHub Actions integration** for CI/CD
- 🔧 **Hotfix workflow** for emergency releases
- 🔧 **Release branch management** with automatic cleanup

#### GitHub Integration
- 🐙 **GitHub Pages deployment** with modern responsive interface
- 🐙 **Automated PR labeling** with conventional commit parsing
- 🐙 **Release notes extraction** with proper categorization
- 🐙 **Multi-platform asset publishing** with checksums
- 🐙 **Install script generation** with dependency checking

### 🧹 Updates

#### Dependencies & Compatibility
- ⬆️ **Go 1.24+ compatibility** with modern language features
- ⬆️ **Cobra CLI framework** updated to v1.9.1
- ⬆️ **Viper configuration** updated to v1.20.1
- ⬆️ **Color output support** with fatih/color v1.18.0
- ⬆️ **Git integration** with go-git v5.16.2
- ⬆️ **Cross-platform compatibility** improvements

#### Build & Release
- 🏗️ **Build optimization** with size reduction flags (-w -s)
- 🏗️ **Apple Silicon support** with darwin/arm64 optimizations
- 🏗️ **Windows executable** handling with proper .exe extensions
- 🏗️ **Checksum generation** with cross-platform sha256 support
- 🏗️ **Install script** with Git Bash/WSL detection

### 🐛 Bug Fixes

#### Release System
- 🐛 **ANSI color codes** removed from git tag messages
- 🐛 **macOS sed compatibility** with empty string extensions
- 🐛 **Regex syntax errors** in changelog generation
- 🐛 **Associative array issues** replaced with simple variables
- 🐛 **Awk newline handling** for proper changelog formatting
- 🐛 **stderr redirection** for clean command output

#### Template & Structure
- 🐛 **Missing Vue components** and icons directory
- 🐛 **TypeScript configuration** for Vue templates
- 🐛 **Base CSS files** inclusion in templates
- 🐛 **Package.json validation** and dependency checking
- 🐛 **Template structure** completeness verification

#### CI/CD & Workflows
- 🐛 **GitHub Actions permissions** resolved
- 🐛 **PR auto-labeling** workflow conditions
- 🐛 **Markdown linting** with proper header checking
- 🐛 **Security scanning** and testing improvements
- 🐛 **Deploy process** automation fixes

### 🛠️ Maintenance

#### Code Quality
- 🧹 **Code organization** with proper module structure
- 🧹 **Error handling** improvements across all components
- 🧹 **Documentation** updates with comprehensive examples
- 🧹 **License management** and compliance
- 🧹 **README formatting** and content updates

#### Testing & Security
- 🔒 **CodeQL integration** for security scanning
- 🔒 **Dependabot configuration** for automated updates
- 🔒 **Security policy** establishment
- 🔒 **Code review** automation with proper workflows

### 📚 Documentation

#### User Guides
- 📖 **Installation instructions** for all platforms
- 📖 **Usage examples** with command-line references
- 📖 **Template creation** guides and best practices
- 📖 **Release process** documentation
- 📖 **Contributing guidelines** and development setup

#### Technical Documentation
- 📝 **API documentation** with docstrings
- 📝 **Architecture diagrams** and system overview
- 📝 **Configuration options** reference
- 📝 **Troubleshooting guides** for common issues

Full Changelog: [Initial Release](https://github.com/tienld-0801/tilokit/releases/tag/v0.1.0-dev)

## [0.1.6-dev] - 2025-08-03

### Fixed
- 🐛 improve GitHub release notes extraction with working awk pattern

### Maintenance
- 🏠 bump version to v0.1.6-dev for development

### Changed
- Development continues...

## [0.1.5-dev] - 2025-08-03

### Added
- ✨ simplify release workflow and remove main branch dependency

### Maintenance
- 🏠 bump version to v0.1.5-dev for development

### Changed
- Development continues...

## [0.1.4-dev] - 2025-08-03

### Fixed
- 🐛 use actual release notes content for GitHub releases

### Maintenance
- 🏠 bump version to v0.1.4-dev for development

### Changed
- Development continues...

## [0.1.3-dev] - 2025-08-03

### Added
- ✨ enhance release notes and changelog generation
- ✨ add auto-update PR title with emoji functionality
- ✨ refactor pull request template
- ✨ plus heading pull request template
- ✨ plus heading pull request template with space remove
- ✨ check clean header for file .md
- ✨ enhance release workflow with multi-platform optimizations
- ✨ enhance GitHub workflows with improved labels and PR automation
- ✨ enhance GitHub workflows with improved labels and automation

### Fixed
- 🐛 resolve permissions issues in GitHub Actions
- 🐛 correct PR auto-label workflow condition
- 🐛 resolve install.sh issues identified by CodeRabbit
- 🐛 resolve regex syntax errors in generate-changelog.sh
- 🐛 replace regex with sed for conventional commit parsing
- 🐛 remove dangling else statement in generate-changelog.sh
- 🐛 replace associative arrays with simple variables
- 🐛 use literal newlines instead of \n escapes
- 🐛 use temporary file for awk changelog content
- 🐛 add empty string extension for macOS sed -i compatibility
- 🐛 use temporary file for git tag message to avoid ANSI codes
- 🐛 redirect print messages to stderr in create_release_branch

### Changed
- ♻️ replace emoji workflows with auto-labeling system

### Documentation
- 📚 update CHANGELOG.md with latest improvements
- 📚 fix CHANGELOG.md duplicate entries from testing

### CI/CD
- 🔄 add automated workflows for PR formatting and changelog generation

### Changed
- Development continues...

## [0.1.2-dev] - 2025-08-03

### Added
- ✨ enhance release notes and changelog generation
- ✨ add auto-update PR title with emoji functionality
- ✨ refactor pull request template
- ✨ plus heading pull request template
- ✨ plus heading pull request template with space remove
- ✨ check clean header for file .md
- ✨ enhance release workflow with multi-platform optimizations
- ✨ enhance GitHub workflows with improved labels and PR automation
- ✨ enhance GitHub workflows with improved labels and automation

### Fixed
- 🐛 resolve permissions issues in GitHub Actions
- 🐛 correct PR auto-label workflow condition
- 🐛 resolve install.sh issues identified by CodeRabbit
- 🐛 resolve regex syntax errors in generate-changelog.sh
- 🐛 replace regex with sed for conventional commit parsing
- 🐛 remove dangling else statement in generate-changelog.sh
- 🐛 replace associative arrays with simple variables
- 🐛 use literal newlines instead of \n escapes
- 🐛 use temporary file for awk changelog content

### Changed
- ♻️ replace emoji workflows with auto-labeling system

### Documentation
- 📚 update CHANGELOG.md with latest improvements

### CI/CD
- 🔄 add automated workflows for PR formatting and changelog generation

### Changed
- Development continues...

## [0.1.2-dev] - 2025-08-03

### Added
- ✨ enhance release notes and changelog generation
- ✨ add auto-update PR title with emoji functionality
- ✨ refactor pull request template
- ✨ plus heading pull request template
- ✨ plus heading pull request template with space remove
- ✨ check clean header for file .md
- ✨ enhance release workflow with multi-platform optimizations
- ✨ enhance GitHub workflows with improved labels and PR automation
- ✨ enhance GitHub workflows with improved labels and automation

### Fixed
- 🐛 resolve permissions issues in GitHub Actions
- 🐛 correct PR auto-label workflow condition
- 🐛 resolve install.sh issues identified by CodeRabbit
- 🐛 resolve regex syntax errors in generate-changelog.sh
- 🐛 replace regex with sed for conventional commit parsing
- 🐛 remove dangling else statement in generate-changelog.sh
- 🐛 replace associative arrays with simple variables
- 🐛 use literal newlines instead of \n escapes
- 🐛 use temporary file for awk changelog content

### Changed
- ♻️ replace emoji workflows with auto-labeling system

### CI/CD
- 🔄 add automated workflows for PR formatting and changelog generation

### Added
- ✨ Enhanced multi-platform build system with Apple Silicon optimizations
- 🚀 Comprehensive GitHub Pages deployment with interactive download interface
- 📱 Universal install script supporting Linux, macOS, and Windows (Git Bash/WSL)
- 🔍 Advanced platform detection with WSL and Apple Silicon recognition
- 📊 Release statistics tracking for features, bugs, and improvements
- 🔄 Retry logic for downloads with timeout and connection handling
- 🎨 Modern web interface with copy-to-clipboard functionality
- 💾 Platform-specific download cards with command-line examples

### Fixed
- 🐛 Corrected Go version in workflow to match go.mod (1.24.4)
- 🚀 Fixed Apple Silicon binary builds with proper optimization flags
- 🔧 Enhanced checksum generation compatibility across platforms
- 📝 Improved release notes extraction with proper feature/bugfix categorization
- 🔐 Fixed Windows binary installation and PATH detection
- ⚙️ Better error handling in installation script with detailed feedback
- 📱 Resolved download URL consistency across all platforms

### Changed
- ♾️ Enhanced build flags with size optimization (-w -s)
- 📦 Improved binary naming convention and file extension handling
- 🎨 Modern gradient design for GitHub Pages with responsive layout
- 📄 Better structured release notes with emoji categorization
- 🔍 Enhanced logging system with colored output and progress indicators
- 🚀 Streamlined installation process with better user feedback

## [0.1.1-dev] - 2025-08-01

### Added
- ✨ Improved hotfix script with conventional commits support
- 🔧 Enhanced release automation with better validation
- 📋 Comprehensive pre-flight checks for releases
- 🎯 Interactive prompts with confirmation dialogs
- 🛡️ Tag existence validation to prevent duplicates

### Fixed
- 🐛 Fixed commit message formatting in hotfix script
- 🔧 Corrected merge process for hotfix branches
- 📝 Improved CHANGELOG.md update mechanism

### Changed
- ♻️ Refactored release scripts to follow conventional commits
- 🔄 Enhanced GitFlow integration with better branch management
- 📚 Updated documentation and inline help messages

