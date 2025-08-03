.PHONY: help build test clean install dev lint docker

# Build configurations  
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -X main.Version=$(VERSION) -X main.BuildDate=$(BUILD_DATE) -X main.GitCommit=$(GIT_COMMIT)

# Binary name
BINARY_NAME := tilokit

# Default target
default: help build

help: ## Show this help message
	@echo "TiLoKit - Modern Multi-Framework Project Generator"
	@echo ""
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

dev: ## Run in development mode
	go run . --help

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) .
	@echo "✅ Build complete: ./$(BINARY_NAME)"

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

##@ Installation

install: build ## Install the binary to /usr/local/bin
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "✅ $(BINARY_NAME) installed successfully"

uninstall: ## Uninstall the binary
	@echo "Uninstalling $(BINARY_NAME)..."
	sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✅ $(BINARY_NAME) uninstalled"

##@ Testing & Examples

test-react: build ## Test React project generation
	@echo "Testing React project generation..."
	./$(BINARY_NAME) --name example-react --framework react --build-tool vite --output ./examples --force --quiet
	@echo "✅ React project generated in ./examples/example-react"

test-vue: build ## Test Vue project generation
	@echo "Testing Vue project generation..."
	./$(BINARY_NAME) --name example-vue --framework vue --build-tool vite --output ./examples --force --quiet
	@echo "✅ Vue project generated in ./examples/example-vue"

test-all: test-react test-vue ## Test all framework generations
	@echo "✅ All framework tests completed"

##@ Cleanup

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	rm -rf dist/
	rm -rf examples/
	rm -rf test-*-app/
	@echo "✅ Cleanup complete"

clean-examples: ## Clean only example projects
	@echo "Cleaning example projects..."
	rm -rf examples/
	rm -rf test-*-app/
	@echo "✅ Example projects cleaned"

##@ Docker

docker: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t tilokit:$(VERSION) .
	@echo "✅ Docker image built: tilokit:$(VERSION)"

##@ Release Management

init-branches: ## Initialize Git branch structure
	@echo "Initializing branch structure..."
	@chmod +x scripts/init-branches.sh
	./scripts/init-branches.sh

release: ## Create a new release (usage: make release VERSION=v0.1.0)
	@if [ -z "$(VERSION)" ]; then \
		echo "❌ VERSION is required. Usage: make release VERSION=v0.1.0"; \
		exit 1; \
	fi
	@echo "Starting release process for $(VERSION)..."
	@chmod +x scripts/release.sh
	./scripts/release.sh $(VERSION)

hotfix: ## Create a hotfix release (usage: make hotfix VERSION=v0.1.1)
	@if [ -z "$(VERSION)" ]; then \
		echo "❌ VERSION is required. Usage: make hotfix VERSION=v0.1.1"; \
		exit 1; \
	fi
	@echo "Starting hotfix process for $(VERSION)..."
	@chmod +x scripts/hotfix.sh
	./scripts/hotfix.sh $(VERSION)

generate-changelog: ## Generate changelog from conventional commits (usage: make generate-changelog VERSION=v0.1.1)
	@if [ -z "$(VERSION)" ]; then \
		echo "❌ VERSION is required. Usage: make generate-changelog VERSION=v0.1.1"; \
		exit 1; \
	fi
	@echo "Generating changelog for $(VERSION)..."
	@chmod +x scripts/generate-changelog.sh
	./scripts/generate-changelog.sh $(VERSION)

check-release: ## Check if ready for release
	@echo "🔍 Checking release readiness..."
	@echo "Current branch: $$(git branch --show-current)"
	@echo "Working directory status:"
	@git status --porcelain || echo "✅ Working directory is clean"
	@echo "Latest tags:"
	@git tag --sort=-version:refname | head -5
	@echo "Unreleased commits:"
	@git log --oneline $$(git describe --tags --abbrev=0 2>/dev/null || echo "HEAD~10")..HEAD | head -10

version-info: ## Show current version information
	@echo "📋 Version Information"
	@echo "Current Version: $(VERSION)"
	@echo "Build Date: $(BUILD_DATE)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Binary: ./$(BINARY_NAME)"
	@if [ -f "./$(BINARY_NAME)" ]; then \
		echo "Binary version:"; \
		./$(BINARY_NAME) version; \
	else \
		echo "Binary not built (run 'make build' first)"; \
	fi
