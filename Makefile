# Netwatch - Network Traffic Monitor Makefile
# Build configuration
APP_NAME := netwatch
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date +%Y-%m-%dT%H:%M:%S%z)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go configuration
GO := go
GO_VERSION := 1.21
BINARY_NAME := $(APP_NAME)
MAIN_PATH := ./cmd/netwatch
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT) -s -w"

# Build directories
BUILD_DIR := build
DIST_DIR := dist

# Testing configuration
TEST_TIMEOUT := 30s
COVERAGE_FILE := coverage.out

# Default target
.DEFAULT_GOAL := build

# Help target
.PHONY: help
help: ## Display this help message
	@echo "Netwatch - Network Traffic Monitor"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development targets
.PHONY: build
build: ## Build the application binary
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

.PHONY: run
run: ## Run the application in development mode
	@echo "Running $(BINARY_NAME) in development mode..."
	$(GO) run $(MAIN_PATH) --dev-mode --log-level=debug

.PHONY: install
install: ## Install the application binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install $(LDFLAGS) $(MAIN_PATH)

# Testing targets
.PHONY: test
test: ## Run unit tests
	@echo "Running unit tests..."
	$(GO) test -v -timeout=$(TEST_TIMEOUT) ./...

.PHONY: test-race
test-race: ## Run unit tests with race detection
	@echo "Running unit tests with race detection..."
	$(GO) test -v -race -timeout=$(TEST_TIMEOUT) ./...

.PHONY: test-coverage
test-coverage: ## Run unit tests with coverage report
	@echo "Running unit tests with coverage..."
	$(GO) test -v -coverprofile=$(COVERAGE_FILE) ./...
	$(GO) tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-integration
test-integration: ## Run integration tests
	@echo "Running integration tests..."
	@if [ -n "$$(find ./tests/integration -name '*_test.go' -print -quit 2>/dev/null)" ]; then \
		$(GO) test -v -tags=integration -timeout=60s ./tests/integration/...; \
	else \
		echo "⚠️ No integration tests found in ./tests/integration/ - skipping"; \
		echo "✅ Integration test check completed (no tests to run)"; \
	fi

.PHONY: benchmark
benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem ./...

# Code quality targets
.PHONY: lint
lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

.PHONY: fmt
fmt: ## Format code
	@echo "Formatting code..."
	$(GO) fmt ./...

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	$(GO) vet ./...

.PHONY: tidy
tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	$(GO) mod tidy

.PHONY: verify
verify: fmt vet tidy test ## Run all verification steps
	@echo "All verification steps completed"

# Build targets
.PHONY: build-all
build-all: ## Build for all supported platforms
	@echo "Building for all platforms..."
	@mkdir -p $(DIST_DIR)
	@for os in linux darwin windows; do \
		for arch in amd64 arm64; do \
			if [ "$$os" = "windows" ]; then \
				ext=".exe"; \
			else \
				ext=""; \
			fi; \
			echo "Building $$os/$$arch..."; \
			GOOS=$$os GOARCH=$$arch $(GO) build $(LDFLAGS) \
				-o $(DIST_DIR)/$(BINARY_NAME)-$$os-$$arch$$ext $(MAIN_PATH); \
		done; \
	done
	@echo "Cross-platform builds complete"

.PHONY: release
release: verify build-all ## Create release builds
	@echo "Creating release packages..."
	@cd $(DIST_DIR) && for f in $(BINARY_NAME)-*; do \
		if [ -f "$$f" ]; then \
			echo "Creating archive for $$f..."; \
			tar -czf "$$f.tar.gz" "$$f"; \
		fi; \
	done
	@echo "Release packages created in $(DIST_DIR)/"

# Cleanup targets
.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(DIST_DIR)
	@rm -f $(COVERAGE_FILE)
	@rm -f coverage.html
	@echo "Clean complete"

.PHONY: clean-cache
clean-cache: ## Clean Go build cache
	@echo "Cleaning Go build cache..."
	$(GO) clean -cache -testcache -modcache

# Development utilities
.PHONY: deps
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GO) mod download

.PHONY: deps-update
deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	$(GO) get -u ./...
	$(GO) mod tidy

.PHONY: generate
generate: ## Run go generate
	@echo "Running go generate..."
	$(GO) generate ./...

# Version information
.PHONY: version
version: ## Display version information
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Go Version: $(shell $(GO) version)"

# Check Go version
.PHONY: check-go-version
check-go-version: ## Check Go version meets requirements
	@echo "Checking Go version..."
	@$(GO) version | grep -q "go$(GO_VERSION)" || (echo "Error: Go $(GO_VERSION)+ required" && exit 1)
	@echo "Go version check passed"

# Development server with auto-reload (requires air)
.PHONY: dev
dev: ## Run development server with auto-reload (requires air)
	@if command -v air > /dev/null; then \
		echo "Starting development server with auto-reload..."; \
		air -c .air.toml; \
	else \
		echo "air not found. Install it with: go install github.com/cosmtrek/air@latest"; \
		echo "Falling back to normal run..."; \
		$(MAKE) run; \
	fi

# Docker targets (optional)
.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):$(VERSION) .

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run --rm -p 8080:8080 $(APP_NAME):$(VERSION)

# Security scan (requires gosec)
.PHONY: security
security: ## Run security scan (requires gosec)
	@if command -v gosec > /dev/null; then \
		echo "Running security scan..."; \
		gosec ./...; \
	else \
		echo "gosec not found. Install it with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi