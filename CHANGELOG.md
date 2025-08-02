# Changelog

All notable changes to TiLoKit will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
- 🐛 use literal newlines instead of 
 escapes
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

## [0.1.0-dev] - Development Phase

### Added
- 🚀 **Multi-Framework Support**: Complete plugin architecture supporting 25+ frameworks
- 🔧 **Universal Build Tools**: Support for language-specific build tools and package managers
- 🎯 **Interactive CLI**: Beautiful command-line interface with prompts and validation
- 📦 **Plugin System**: Extensible architecture for easy framework additions
- 🐳 **Docker Support**: Container configurations for all supported frameworks
- 🔄 **CI/CD Integration**: GitHub Actions workflows for testing and releases
- 📋 **Comprehensive Testing**: Unit tests and integration testing framework
- 🛠️ **Development Tools**: Build scripts, linting, and formatting automation

### Framework Support
- **JavaScript/TypeScript**: React, Vue, Angular, Svelte, Next.js, Nuxt.js
- **Python**: Django, Flask, FastAPI
- **PHP**: Laravel, Symfony
- **Java**: Spring Boot, Quarkus
- **Go**: Gin, Echo, Fiber
- **Rust**: Actix, Rocket, Axum
- **C#**: ASP.NET Core, Blazor
- **Ruby**: Rails, Sinatra
- **Node.js**: Express, NestJS, Fastify
- **Mobile**: React Native, Flutter, Ionic
- **Desktop**: Electron, Tauri, Wails

### Build Tools
- **JavaScript**: Vite, Webpack, Rollup, Parcel
- **Package Managers**: npm, yarn, pnpm
- **Language-Specific**: pip, poetry, composer, maven, gradle, cargo, dotnet, bundler

### Infrastructure
- Modern Go architecture with clean separation of concerns
- Plugin registry with lifecycle management
- Rich execution context for project generation
- Configuration management with YAML support
- Template engine for dynamic file generation
- Git integration with smart .gitignore generation

### Development Features
- Initial project structure
- Core engine implementation
- Basic CLI framework setup
- Basic React and Vue support
- Vite build tool integration
- CLI foundation with Cobra

---

> **Note**: TiLoKit is currently in active development. This changelog will be updated as features are implemented and released.
