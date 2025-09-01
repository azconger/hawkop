# HawkOp CLI Makefile
# Go CLI companion tool for StackHawk scanner and platform

.PHONY: help build test clean install format lint check release deps

# Default target
.DEFAULT_GOAL := help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=hawkop
BINARY_UNIX=$(BINARY_NAME)_unix

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS=-ldflags "-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Install and verify dependencies
	$(GOMOD) download
	$(GOMOD) verify
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && $(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)

format: deps ## Format code and organize imports
	$(GOCMD) fmt ./...
	$(GOMOD) tidy
	@which gci > /dev/null || $(GOCMD) install github.com/daixiang0/gci@latest
	@which gci > /dev/null && gci write --skip-generated -s standard -s default -s "prefix(hawkop)" . || echo "gci not available, skipping import organization"

lint: format ## Run linters and static analysis
	$(GOCMD) vet ./...
	golangci-lint run

test: deps ## Run all tests (unit, integration, end-to-end)
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-unit: deps ## Run unit tests only
	$(GOTEST) -v -short ./...

test-integration: deps ## Run integration tests (with external dependencies)
	$(GOTEST) -v -run Integration ./...

bench: deps ## Run benchmarks
	$(GOTEST) -bench=. -benchmem ./...

check: lint test ## Run all checks (format, lint, test)

build: deps ## Build the binary
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v .

build-all: deps ## Build for all platforms
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)_linux_amd64 .
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)_linux_arm64 .
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)_darwin_amd64 .
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)_darwin_arm64 .
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)_windows_amd64.exe .

install: build ## Install the binary to GOPATH/bin
	$(GOCMD) install $(LDFLAGS) .

clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)_*
	rm -f coverage.out coverage.html

release: check ## Create and push a new release tag
	@echo "Current version: $(VERSION)"
	@read -p "Enter new version (e.g., v1.2.3): " new_version; \
	if [ -z "$$new_version" ]; then \
		echo "Version cannot be empty"; \
		exit 1; \
	fi; \
	if ! echo "$$new_version" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+$$'; then \
		echo "Version must be in format v1.2.3"; \
		exit 1; \
	fi; \
	echo "Creating release $$new_version..."; \
	git tag -a "$$new_version" -m "Release $$new_version"; \
	git push origin "$$new_version"; \
	echo "Release $$new_version created and pushed. GitHub Actions will handle the build and release."

pre-commit: format lint test ## Run all pre-commit checks
	@echo "✅ All pre-commit checks passed"

dev: ## Install development tools
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOCMD) install github.com/goreleaser/goreleaser/v2@latest
	$(GOCMD) install github.com/daixiang0/gci@latest

# Quick development workflow targets
quick-test: ## Run tests without coverage
	$(GOTEST) ./...

quick-build: ## Build without full checks
	$(GOBUILD) -o $(BINARY_NAME) .

watch: ## Watch for changes and run tests (requires entr: brew install entr)
	find . -name "*.go" | entr -c make quick-test