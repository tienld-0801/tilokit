# Contributing to TiLoKit

> ⚠️ **Development Status**: TiLoKit is currently under active development. This contributing guide will be updated as the project matures.

Thank you for your interest in contributing to TiLoKit! This document provides preliminary guidelines for contributors.

## 🚀 Getting Started

> **Note**: As the project is in development, these requirements may change.

### Prerequisites

- Go 1.25.0 or later
- Git
- Make (recommended for development)
- Node.js and npm (for testing JavaScript projects)
- Additional language runtimes as needed (Python, PHP, Java, etc.)

### Development Setup

1. **Fork and clone the repository**
   ```bash
   git clone https://github.com/your-username/tilokit.git
   cd tilokit
   ```

2. **Install dependencies and setup hooks**
   ```bash
   go mod download && go mod tidy
   make install-hooks  # Install Git hooks for commit validation
   ```

3. **Build the project**
   ```bash
   make build
   ```

4. **Run development mode**
   ```bash
   make run  # Run with full validation (test, lint, markdown-lint)
   # or
   make dev  # Run help mode only
   ```

5. **Run security checks**
   ```bash
   make security-check      # Quick security scan
   make security-report     # Detailed JSON report
   ```

## 🏗️ Project Structure

```
tilokit/
├── cmd/                    # CLI commands and entry points
│   └── root.go            # Main command configuration
├── internal/
│   ├── cli/               # CLI logic and handlers
│   │   ├── generation/    # Project generation orchestration
│   │   │   ├── orchestrator.go    # Main workflow coordinator
│   │   │   ├── prompts/           # Interactive user input handling
│   │   │   ├── validator/         # Input validation logic
│   │   │   ├── registry/          # Plugin registry management
│   │   │   └── progress/          # Animated progress UI
│   │   ├── handlers.go    # CLI command handlers
│   │   ├── manager.go     # Project management logic
│   │   ├── update.go      # Self-update functionality
│   │   └── validation.go  # Input validation
│   ├── core/              # Core engine and registry
│   │   ├── engine/        # Template execution engine
│   │   ├── registry/      # Plugin registry system
│   │   └── context/       # Execution context management
│   ├── plugins/           # Plugin implementations
│   │   ├── frameworks/    # Framework plugins (40+ frameworks)
│   │   ├── builders/      # Build tool plugins (Vite, Webpack, etc.)
│   │   └── templates/     # Template processing plugins
│   ├── templates/         # Template definitions for frameworks
│   ├── ui/                # User interface components (Bubble Tea)
│   ├── utils/             # Utility functions
│   └── config/            # Configuration management
├── .github/               # GitHub workflows and templates
├── .husky/                # Git hooks for commit validation
├── assets/                # Static assets (logos, icons)
├── config/                # Configuration files
├── docs/                  # Documentation
├── scripts/               # Build and deployment scripts
├── Dockerfile            # Container configuration
├── Makefile              # Build automation with security checks
└── main.go               # Application entry point
```

## 🔌 Plugin Development

> ⚠️ **Coming Soon**: Plugin development documentation is being prepared.

TiLoKit uses a plugin-based architecture. The plugin system supports multiple languages and frameworks:

### Framework Plugin

```go
package frameworks

import (
    "github.com/ti-lo/tilokit/internal/core/context"
    "github.com/ti-lo/tilokit/internal/core/registry"
)

type MyFrameworkPlugin struct{}

func NewMyFrameworkPlugin() *MyFrameworkPlugin {
    return &MyFrameworkPlugin{}
}

func (p *MyFrameworkPlugin) Name() string {
    return "my-framework"
}

func (p *MyFrameworkPlugin) SupportedFrameworks() []string {
    return []string{"myframework"}
}

func (p *MyFrameworkPlugin) SupportedBuildTools() []string {
    return []string{"vite", "webpack"}
}

func (p *MyFrameworkPlugin) PreGenerate(ctx *tilocontext.ExecutionContext) error {
    // Pre-generation setup
    return nil
}

func (p *MyFrameworkPlugin) Generate(ctx *tilocontext.ExecutionContext) error {
    // Main generation logic
    return nil
}

func (p *MyFrameworkPlugin) PostGenerate(ctx *tilocontext.ExecutionContext) error {
    // Post-generation cleanup
    return nil
}
```

### Build Tool Plugin

```go
package builders

type MyBuilderPlugin struct{}

func (p *MyBuilderPlugin) Name() string {
    return "my-builder"
}

func (p *MyBuilderPlugin) SupportedBuildTools() []string {
    return []string{"mybuildtool"}
}

// Implement other required methods...
```

## 📝 Coding Standards

> **Note**: These standards are being refined during development.

### Go Code Style

- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions small and focused
- Handle errors appropriately

### Commit Messages

**IMPORTANT**: TiLoKit requires emojis in commit messages!

Required format:
```
🎯 type: description

[optional body]

[optional footer]
```

**Valid emoji mappings:**
- ✨ `feat` - new features
- 🐛 `fix` - bug fixes
- 📚 `docs` - documentation changes
- ♻️ `refactor` - code refactoring
- ⚡ `perf` - performance improvements
- 🧪 `test` - adding or updating tests
- 🛠️ `build` - build system changes
- 🔄 `ci` - CI/CD changes
- 🧹 `chore` - maintenance tasks
- 🎨 `style` - code style changes
- ⏪ `revert` - reverting previous commits
- 🚀 `release` - version releases

Examples:
```
✨ feat: add Svelte framework support
🐛 fix: resolve Vite configuration path issue
📚 docs: update installation instructions
🔒 feat: add gosec security analysis
```

**Git hooks will automatically validate commit messages!**

## 🧪 Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...

# Run integration tests
make test-all
```

### Writing Tests

- Write unit tests for all new functions
- Use table-driven tests when appropriate
- Mock external dependencies
- Test both success and error cases

Example test:
```go
func TestMyFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    "test",
            expected: "test-result",
            wantErr:  false,
        },
        // Add more test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := MyFunction(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("MyFunction() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if result != tt.expected {
                t.Errorf("MyFunction() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

## 📋 Pull Request Process

> **Note**: As the project is in active development, the PR process may evolve.

1. **Create a feature branch**
   ```bash
   git checkout -b feature/my-new-feature
   ```

2. **Make your changes**
   - Write code following the style guidelines
   - Add tests for new functionality
   - Update documentation if needed

3. **Test your changes**
   ```bash
   make test
   make security-check
   make lint
   ```

4. **Commit your changes**
   ```bash
   git add .
   git commit -m "✨ feat: add my new feature"
   ```

5. **Push to your fork**
   ```bash
   git push origin feature/my-new-feature
   ```

6. **Create a Pull Request**
   - Use a clear title and description
   - Reference any related issues
   - Include screenshots for UI changes
   - Ensure all CI checks pass

### PR Review Checklist

- [ ] Code follows project style guidelines
- [ ] Tests are included and passing
- [ ] Documentation is updated
- [ ] No breaking changes (or properly documented)
- [ ] CI/CD pipeline passes
- [ ] Code is reviewed by at least one maintainer

## 🐛 Bug Reports

> **Note**: During development phase, bugs are expected. Please report them to help improve the project.

When reporting bugs, please include:

1. **Environment information**
   - Operating system
   - Go version
   - TiLoKit version

2. **Steps to reproduce**
   - Exact commands used
   - Expected behavior
   - Actual behavior

3. **Additional context**
   - Error messages
   - Log output
   - Screenshots (if applicable)

## 💡 Feature Requests

> **Note**: Feature requests are welcome as we shape the project's direction.

When requesting features:

1. **Describe the problem** you're trying to solve
2. **Explain the proposed solution**
3. **Consider alternatives** you've thought about
4. **Provide use cases** and examples

## 📄 License

By contributing to TiLoKit, you agree that your contributions will be licensed under the MIT License.

## 🤝 Community

> **Note**: Community features will be expanded as the project grows.

- 💬 **Join our Discord**: [Discord](https://discord.gg/BqTqZ46uT9) - Get help, share projects, and connect with other developers
- Join discussions in GitHub Issues
- Follow the project for development updates
- Star the repository to show support
- Share feedback and suggestions

## 📞 Getting Help

> **Note**: Support channels are being established.

If you need help:

1. Check the [README](README.md) for current status
2. Search existing issues for similar problems
3. Create a new issue with detailed information
4. Be patient as the project is in development

---

## 🚧 Development Roadmap

### Phase 1: Core Architecture *(Completed)*
- ✅ Plugin system foundation
- ✅ CLI structure với Bubble Tea UI
- ✅ Configuration system
- ✅ Template engine với custom delimiters
- ✅ Modular generation orchestration

### Phase 2: Multi-Framework Support *(Completed)*
- ✅ JavaScript/TypeScript: React, Vue, Svelte, Angular, Next.js, Nuxt
- ✅ Node.js: Express, Fastify, NestJS
- ✅ Python: Django, Flask, FastAPI
- ✅ PHP: Laravel, Symfony, CakePHP, CodeIgniter
- ✅ Java/Kotlin: Spring Boot
- ✅ Go: Gin, Echo, Fiber
- ✅ Ruby: Rails
- ✅ C#: ASP.NET Core

### Phase 3: Mobile & Desktop *(Completed)*
- ✅ Mobile: React Native (Expo), Flutter
- ✅ Desktop: Electron
- 📋 iOS: Swift, SwiftUI
- 📋 Android: Kotlin, Jetpack Compose

### Phase 4: Advanced Features *(In Progress)*
- ✅ Security scanning với gosec
- ✅ Discord CI notifications
- ✅ Animated progress bars
- 📋 Plugin marketplace
- 📋 Custom templates
- 📋 IDE integrations

Thank you for your interest in contributing to TiLoKit! 🚀
