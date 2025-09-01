# Makefile for go-example project

# Project information
PROJECT_NAME := go-example
SERVICE_NAME := go-example-api
MODULE_NAME := github.com/kart-io/go-example
VERSION_PKG := github.com/kart-io/version

# Git information
GIT_VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.1.0-dev")
GIT_COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH := $(shell git branch --show-current 2>/dev/null || echo "unknown")
GIT_TREE_STATE := $(shell if [ -n "$$(git status --porcelain 2>/dev/null)" ]; then echo "dirty"; else echo "clean"; fi)
BUILD_DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

# Build flags
LDFLAGS := -w -s \
	-X '$(VERSION_PKG).serviceName=$(SERVICE_NAME)' \
	-X '$(VERSION_PKG).gitVersion=$(GIT_VERSION)' \
	-X '$(VERSION_PKG).gitCommit=$(GIT_COMMIT)' \
	-X '$(VERSION_PKG).gitBranch=$(GIT_BRANCH)' \
	-X '$(VERSION_PKG).gitTreeState=$(GIT_TREE_STATE)' \
	-X '$(VERSION_PKG).buildDate=$(BUILD_DATE)'

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: help
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(BLUE)%-15s$(NC) %s\n", $$1, $$2}'

.PHONY: deps
deps: ## Download and tidy dependencies
	@echo "$(GREEN)[INFO]$(NC) Downloading dependencies..."
	go mod download
	go mod tidy

.PHONY: build
build: ## Build the gin demo application
	@echo "$(GREEN)[INFO]$(NC) Building $(SERVICE_NAME)..."
	@echo "Version: $(GIT_VERSION)"
	@echo "Commit: $(GIT_COMMIT)"
	@echo "Branch: $(GIT_BRANCH)"
	@echo "Date: $(BUILD_DATE)"
	go build -ldflags "$(LDFLAGS)" -o bin/gin-demo ./gin-demo

.PHONY: run
run: ## Run the gin demo application with version injection
	@echo "$(GREEN)[INFO]$(NC) Starting $(SERVICE_NAME)..."
	@echo "$(YELLOW)[INFO]$(NC) Server will start on :8082"
	@echo "$(YELLOW)[INFO]$(NC) Available endpoints:"
	@echo "  - http://localhost:8082/        (root)"
	@echo "  - http://localhost:8082/health  (health check)"
	@echo "  - http://localhost:8082/version (version info)"
	@echo "$(YELLOW)[INFO]$(NC) Press Ctrl+C to stop"
	go run -ldflags "$(LDFLAGS)" ./gin-demo

.PHONY: dev
dev: deps run ## Setup dependencies and run in development mode

.PHONY: test
test: ## Run tests
	@echo "$(GREEN)[INFO]$(NC) Running tests..."
	go test -v ./...

.PHONY: fmt
fmt: ## Format Go code
	@echo "$(GREEN)[INFO]$(NC) Formatting Go code..."
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	@echo "$(GREEN)[INFO]$(NC) Running go vet..."
	go vet ./...

.PHONY: check
check: fmt vet test ## Run all checks (format, vet, test)
	@echo "$(GREEN)[INFO]$(NC) All checks completed successfully"

.PHONY: clean
clean: ## Clean build artifacts
	@echo "$(GREEN)[INFO]$(NC) Cleaning build artifacts..."
	rm -rf bin/

.PHONY: version
version: ## Show version information
	@echo "Project: $(PROJECT_NAME)"
	@echo "Service: $(SERVICE_NAME)"
	@echo "Module: $(MODULE_NAME)"
	@echo "Version: $(GIT_VERSION)"
	@echo "Commit: $(GIT_COMMIT)"
	@echo "Branch: $(GIT_BRANCH)"
	@echo "Tree State: $(GIT_TREE_STATE)"
	@echo "Build Date: $(BUILD_DATE)"

.PHONY: file-logging-demo
file-logging-demo: ## Run file logging demonstration
	@echo "$(GREEN)[INFO]$(NC) Running file logging demo..."
	@echo "$(YELLOW)[INFO]$(NC) This demo shows various file logging configurations"
	@cd file-logging-demo && make run

.PHONY: file-logging-examples
file-logging-examples: ## Show file logging configuration examples
	@echo "$(GREEN)[INFO]$(NC) Showing file logging configuration examples..."
	@cd file-logging-demo && make examples

.PHONY: file-logging-logs
file-logging-logs: ## View generated file logs
	@echo "$(GREEN)[INFO]$(NC) Generated log files:"
	@cd file-logging-demo && make logs

.PHONY: viper-config-demo
viper-config-demo: ## Run Viper configuration demo
	@echo "$(GREEN)[INFO]$(NC) Running Viper configuration demo..."
	@echo "$(YELLOW)[INFO]$(NC) This demo shows YAML configuration with Viper"
	@cd viper-config-demo && make run

.PHONY: viper-config-dev
viper-config-dev: ## Run Viper demo in development mode
	@echo "$(GREEN)[INFO]$(NC) Running Viper demo in development mode..."
	@cd viper-config-demo && make run-dev

.PHONY: viper-config-prod
viper-config-prod: ## Run Viper demo in production mode
	@echo "$(GREEN)[INFO]$(NC) Running Viper demo in production mode..."
	@cd viper-config-demo && make run-prod

.PHONY: viper-config-test
viper-config-test: ## Run Viper demo in testing mode
	@echo "$(GREEN)[INFO]$(NC) Running Viper demo in testing mode..."
	@cd viper-config-demo && make run-test

.PHONY: demos
demos: ## Run all available demos
	@echo "$(GREEN)[INFO]$(NC) Running all available demos..."
	@echo ""
	@echo "$(YELLOW)[DEMO 1]$(NC) File Logging Demo"
	@make file-logging-demo
	@echo ""
	@echo "$(YELLOW)[DEMO 2]$(NC) Viper Configuration Demo"
	@echo "$(YELLOW)[INFO]$(NC) Starting Viper demo briefly to show configuration loading"
	@timeout 5s make viper-config-demo || echo "Demo completed"
	@echo ""
	@echo "$(YELLOW)[DEMO 3]$(NC) Gin Web Server Demo"
	@echo "$(YELLOW)[INFO]$(NC) Starting Gin demo on port 8085 (to avoid conflict)"
	@PORT=8085 make run &
	@PID=$$!; sleep 3; \
	echo "$(GREEN)[INFO]$(NC) Testing endpoints..."; \
	curl -s http://localhost:8085/ | head -c 100; echo "..."; \
	curl -s http://localhost:8085/health | head -c 100; echo "..."; \
	curl -s http://localhost:8085/version | head -c 100; echo "..."; \
	kill $$PID 2>/dev/null || true
	@echo ""
	@echo "$(GREEN)[INFO]$(NC) All demos completed!"

.PHONY: clean-logs
clean-logs: ## Clean generated log files
	@echo "$(GREEN)[INFO]$(NC) Cleaning generated log files..."
	@cd file-logging-demo && make clean || true
	@echo "$(GREEN)[INFO]$(NC) Log files cleaned!"

.PHONY: test-endpoints
test-endpoints: ## Test API endpoints (requires service to be running)
	@echo "$(GREEN)[INFO]$(NC) Testing API endpoints..."
	@echo "Testing root endpoint:"
	@curl -s http://localhost:8082/ | jq . || echo "Failed to connect - is the service running?"
	@echo "\nTesting health endpoint:"
	@curl -s http://localhost:8082/health | jq . || echo "Failed to connect - is the service running?"
	@echo "\nTesting version endpoint:"
	@curl -s http://localhost:8082/version | jq . || echo "Failed to connect - is the service running?"

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "$(GREEN)[INFO]$(NC) Building Docker image..."
	docker build \
		--build-arg SERVICE_NAME="$(SERVICE_NAME)" \
		--build-arg VERSION="$(GIT_VERSION)" \
		--build-arg COMMIT="$(GIT_COMMIT)" \
		--build-arg BRANCH="$(GIT_BRANCH)" \
		--build-arg BUILD_DATE="$(BUILD_DATE)" \
		-t $(SERVICE_NAME):$(GIT_VERSION) \
		-t $(SERVICE_NAME):latest \
		.

.PHONY: install
install: build ## Install binary to GOPATH/bin
	@echo "$(GREEN)[INFO]$(NC) Installing $(SERVICE_NAME) to $(shell go env GOPATH)/bin"
	cp bin/gin-demo $(shell go env GOPATH)/bin/gin-demo