.PHONY: help build test clean install dev lint security-check security-report docker hooks install-hooks uninstall-hooks

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -X main.Version=$(VERSION) -X main.BuildDate=$(BUILD_DATE) -X main.GitCommit=$(GIT_COMMIT)

BINARY_NAME := tilokit

default: help build

help:
	@echo "TiLoKit - Modern Multi-Framework Project Generator"
	@echo ""
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development
dev:
	@./.husky/check-hooks.sh 2>/dev/null || true
	go run . --help

run: test lint markdown-lint
	@./.husky/check-hooks.sh 2>/dev/null || true
	@echo "🚀 Starting TiLoKit in interactive mode..."
	go run . -i

demo:
	@./.husky/check-hooks.sh 2>/dev/null || true
	@echo "🎯 Creating demo React project..."
	go run . -n demo-react -f react -b vite -L js -q

build:
	@./.husky/check-hooks.sh 2>/dev/null || true
	@echo "Building $(BINARY_NAME)..."
	go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) .
	@echo "✅ Build completed: ./$(BINARY_NAME)"

lint:
	@echo "Running linter..."
	golangci-lint run

security-check:
	@echo "🔒 Running security analysis with gosec..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec -fmt=colored -stdout -verbose=text ./...; \
	else \
		echo "❌ gosec not found. Install with: go install github.com/securego/gosec/v2/cmd/gosec@latest"; \
		exit 1; \
	fi

##@ Security
security-report:
	@echo "📊 Generating detailed security report..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec -fmt=json -out=security-report.json ./... && \
		echo "✅ Security report saved to security-report.json"; \
	else \
		echo "❌ gosec not found. Install with: go install github.com/securego/gosec/v2/cmd/gosec@latest"; \
		exit 1; \
	fi

##@ Installation
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "✅ $(BINARY_NAME) installed successfully"

uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✅ $(BINARY_NAME) uninstalled"

##@ Testing & Examples
test:
	@echo "Running tests..."
	go test -v ./...

markdown-lint:
	@echo "Running markdownlint..."
	@if command -v markdownlint-cli2 >/dev/null 2>&1; then \
		markdownlint-cli2 "**/*.md" "#vendor"; \
	elif command -v npx >/dev/null 2>&1; then \
		npx markdownlint-cli2 "**/*.md" "#vendor"; \
	else \
		echo "❌ markdownlint-cli2 not found. Install with: npm install -g markdownlint-cli2"; \
		exit 1; \
	fi

test-react: build
	@echo "Testing React project generation..."
	./$(BINARY_NAME) --name example-react --framework react --build-tool vite --output ./examples --force --quiet
	@echo "✅ React project generated in ./examples/example-react"

test-vue: build
	@echo "Testing Vue project generation..."
	./$(BINARY_NAME) --name example-vue --framework vue --build-tool vite --output ./examples --force --quiet
	@echo "✅ Vue project generated in ./examples/example-vue"

test-all: test-react test-vue
	@echo "✅ All framework tests completed"

##@ Cleanup
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	rm -rf dist/
	rm -rf examples/
	rm -rf test-*-app/
	@echo "✅ Cleanup complete"

clean-examples:
	@echo "Cleaning example projects..."
	rm -rf examples/
	rm -rf test-*-app/
	@echo "✅ Example projects cleaned"

##@ Docker
docker:
	@echo "Building Docker image..."
	docker build -t tilokit:$(VERSION) .
	@echo "✅ Docker image built: tilokit:$(VERSION)"

##@ Release Management
init-branches:
	@echo "Initializing branch structure..."
	@chmod +x scripts/init-branches.sh
	./scripts/init-branches.sh

release:
	@if [ -z "$(VERSION)" ]; then \
		echo "❌ VERSION is required. Usage: make release VERSION=v0.1.0"; \
		exit 1; \
	fi
	@echo "🚀 Starting release process for $(VERSION)..."
	@echo "🔧 Ensuring Git hooks are installed..."
	@./.husky/check-hooks.sh
	@echo "🔍 Checking release readiness..."
	@$(MAKE) check-release
	@echo "🏠 Creating release..."
	@chmod +x .github/scripts/release.sh
	@./.github/scripts/release.sh $(VERSION)
	@echo "✨ Release $(VERSION) completed!"

quick-release:
	@echo "⚡ Quick release process..."
	@NEXT_VERSION=$$(date +"v0.1.%s-dev"); \
	echo "🚀 Creating quick release: $$NEXT_VERSION"; \
	$(MAKE) release VERSION=$$NEXT_VERSION

dev-release:
	@NEXT_VERSION="v0.1.$$(date +%s)-dev"; \
	echo "🔧 Creating dev release: $$NEXT_VERSION"; \
	$(MAKE) release VERSION=$$NEXT_VERSION

hotfix:
	@if [ -z "$(VERSION)" ]; then \
		echo "❌ VERSION is required. Usage: make hotfix VERSION=v0.1.1"; \
		exit 1; \
	fi
	@chmod +x scripts/hotfix.sh
	./scripts/hotfix.sh $(VERSION)

revert-release:
	@echo "🔄 Reverting to previous version..."
	@PREV_VERSION=$$(git tag --sort=-version:refname | head -1 2>/dev/null || echo "v0.0.0-dev"); \
	if [ "$$PREV_VERSION" = "v0.0.0-dev" ]; then \
		echo "❌ No previous tags found"; \
		exit 1; \
	fi; \
	echo "📦 Previous version: $$PREV_VERSION"; \
	if [[ "$$OSTYPE" == "darwin"* ]]; then \
		sed -i '' "s/Version   = \".*\"/Version   = \"$$PREV_VERSION\"/" pkg/constants/constants.go; \
	else \
		sed -i "s/Version   = \".*\"/Version   = \"$$PREV_VERSION\"/" pkg/constants/constants.go; \
	fi; \
	echo "✅ Version reverted to $$PREV_VERSION in pkg/constants/constants.go"; \
	echo "💡 You can now commit this change: git add pkg/constants/constants.go && git commit -m '⏪ revert: version to $$PREV_VERSION'"

check-release:
	@echo "🔍 Checking release readiness..."
	@echo "Current branch: $$(git branch --show-current)"
	@echo "Working directory status:"
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "❌ Working directory is dirty:"; \
		git status --short; \
		echo "Please commit or stash changes before release"; \
		exit 1; \
	else \
		echo "✅ Working directory is clean"; \
	fi
	@echo "Latest tags:"
	@git tag --sort=-version:refname | head -5
	@echo "Unreleased commits:"
	@git log --oneline $$(git describe --tags --abbrev=0 2>/dev/null || echo "HEAD~10")..HEAD | head -10

release-status:
	@echo "📋 Release Status"
	@echo "==============="
	@echo "Current branch: $$(git branch --show-current)"
	@echo "Latest tag: $$(git describe --tags --abbrev=0 2>/dev/null || echo 'No tags found')"
	@echo "Commits since last tag: $$(git rev-list $$(git describe --tags --abbrev=0 2>/dev/null || echo HEAD)..HEAD --count 2>/dev/null || echo '0')"
	@echo "Version in code: $$(grep 'Version =' cmd/version.go | cut -d'"' -f2)"
	@echo "Working directory: $$(if [ -n "$$(git status --porcelain)" ]; then echo 'dirty'; else echo 'clean'; fi)"

validate-release:
	@echo "⚙️ Validating release readiness..."
	@$(MAKE) lint
	@$(MAKE) test
	@$(MAKE) build
	@echo "✅ Release validation passed!"

version-info:
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

##@ Git Hooks
hooks: install-hooks

install-hooks:
	@echo "🔧 Installing Git hooks from .husky..."
	@chmod +x .husky/hooks/install-hooks.sh
	@./.husky/hooks/install-hooks.sh

uninstall-hooks:
	@echo "🗑️ Uninstalling Git hooks..."
	@if [ -f ".git/hooks/commit-msg" ]; then \
		rm .git/hooks/commit-msg && echo "✅ Removed commit-msg hook"; \
	else \
		echo "ℹ️ commit-msg hook not found"; \
	fi
	@if [ -f ".git/hooks/pre-commit" ]; then \
		rm .git/hooks/pre-commit && echo "✅ Removed pre-commit hook"; \
	else \
		echo "ℹ️ pre-commit hook not found"; \
	fi
	@echo "🎉 Git hooks uninstalled!"

check-hooks:
	@./.husky/check-hooks.sh

validate-commits:
	@echo "🔍 Validating commit messages..."
	@chmod +x .husky/ci-check-commits.sh
	@./.husky/ci-check-commits.sh

test-emoji:
	@echo "🧪 Testing emoji validation system..."
	@chmod +x scripts/test-emoji-validation.sh
	@./scripts/test-emoji-validation.sh
