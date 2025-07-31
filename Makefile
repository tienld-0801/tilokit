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
