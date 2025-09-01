# Viper Configuration Demo Makefile

# Build variables
BINARY_NAME=viper-config-demo
SERVICE_NAME=viper-config-api
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT ?= $(shell git rev-parse HEAD)
BRANCH ?= $(shell git branch --show-current)
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# Directories
BUILD_DIR=./bin
CONFIG_DIR=./config
LOGS_DIR=./logs

# Available config files
CONFIGS=app.yaml production.yaml testing.yaml

# Build flags with version injection
LDFLAGS=-ldflags "\
	-X 'github.com/kart-io/version.serviceName=$(SERVICE_NAME)' \
	-X 'github.com/kart-io/version.gitVersion=$(VERSION)' \
	-X 'github.com/kart-io/version.gitCommit=$(COMMIT)' \
	-X 'github.com/kart-io/version.gitBranch=$(BRANCH)' \
	-X 'github.com/kart-io/version.buildDate=$(BUILD_DATE)'"

.PHONY: help
help: ## Show available commands
	@echo "Viper Configuration Demo Commands:"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "Configuration Files:"
	@for config in $(CONFIGS); do \
		echo "  📋 $$config"; \
	done

.PHONY: setup
setup: ## Setup project directories and dependencies
	@echo "📁 Creating directories..."
	@mkdir -p $(BUILD_DIR) $(LOGS_DIR)
	@echo "📦 Installing dependencies..."
	@go mod tidy
	@echo "✅ Setup completed"

.PHONY: build
build: setup ## Build with version injection
	@echo "🔨 Building $(BINARY_NAME) with version injection..."
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✅ Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

.PHONY: version
version: build ## Show version information
	@echo "📋 Version Information:"
	@$(BUILD_DIR)/$(BINARY_NAME) --version 2>/dev/null || \
		echo "Version info will be displayed when service starts"

# Run commands for different configurations
.PHONY: run-dev
run-dev: build ## Run with development config (app.yaml)
	@echo "🚀 Starting development server with app.yaml..."
	@$(BUILD_DIR)/$(BINARY_NAME) app.yaml

.PHONY: run-prod
run-prod: build ## Run with production config (production.yaml)
	@echo "🚀 Starting production server with production.yaml..."
	@$(BUILD_DIR)/$(BINARY_NAME) production.yaml

.PHONY: run-test
run-test: build ## Run with testing config (testing.yaml)
	@echo "🚀 Starting testing server with testing.yaml..."
	@$(BUILD_DIR)/$(BINARY_NAME) testing.yaml

.PHONY: run
run: run-dev ## Default run (development mode)

# Environment-based runs
.PHONY: run-env-dev
run-env-dev: build ## Run using APP_ENV=app environment variable
	@echo "🚀 Starting with APP_ENV=app..."
	@APP_ENV=app $(BUILD_DIR)/$(BINARY_NAME)

.PHONY: run-env-prod
run-env-prod: build ## Run using APP_ENV=production environment variable
	@echo "🚀 Starting with APP_ENV=production..."
	@APP_ENV=production $(BUILD_DIR)/$(BINARY_NAME)

.PHONY: run-env-test
run-env-test: build ## Run using APP_ENV=testing environment variable
	@echo "🚀 Starting with APP_ENV=testing..."
	@APP_ENV=testing $(BUILD_DIR)/$(BINARY_NAME)

# Test endpoints (requires service to be running)
.PHONY: test-endpoints
test-endpoints: ## Test all API endpoints (service must be running)
	@echo "🧪 Testing API endpoints..."
	@echo ""
	@echo "📍 Testing root endpoint:"
	@curl -s http://localhost:8083/ | python3 -m json.tool || echo "❌ Failed to connect to port 8083"
	@echo ""
	@echo "📍 Testing health endpoint:"
	@curl -s http://localhost:8083/health | python3 -m json.tool || echo "❌ Failed to connect"
	@echo ""
	@echo "📍 Testing version endpoint:"
	@curl -s http://localhost:8083/version | python3 -m json.tool || echo "❌ Failed to connect"
	@echo ""
	@echo "📍 Testing config endpoint:"
	@curl -s http://localhost:8083/config | python3 -m json.tool || echo "❌ Failed to connect"
	@echo ""
	@echo "📍 Testing logger test endpoint:"
	@curl -s http://localhost:8083/logger/test | python3 -m json.tool || echo "❌ Failed to connect"
	@echo ""

.PHONY: test-prod-endpoints
test-prod-endpoints: ## Test production endpoints (port 8080)
	@echo "🧪 Testing production API endpoints..."
	@echo ""
	@echo "📍 Testing root endpoint:"
	@curl -s http://localhost:8080/ | python3 -m json.tool || echo "❌ Failed to connect to port 8080"
	@echo ""
	@echo "📍 Testing config endpoint:"
	@curl -s http://localhost:8080/config | python3 -m json.tool || echo "❌ Failed to connect"

.PHONY: test-testing-endpoints  
test-testing-endpoints: ## Test testing endpoints (port 8084)
	@echo "🧪 Testing testing API endpoints..."
	@echo ""
	@echo "📍 Testing root endpoint:"
	@curl -s http://localhost:8084/ | python3 -m json.tool || echo "❌ Failed to connect to port 8084"

# Show configuration files
.PHONY: show-configs
show-configs: ## Display all configuration files
	@echo "📋 Available Configuration Files:"
	@echo ""
	@for config in $(CONFIGS); do \
		echo "=== $$config ==="; \
		cat $(CONFIG_DIR)/$$config; \
		echo ""; \
	done

.PHONY: validate-configs
validate-configs: ## Validate YAML configuration files
	@echo "✅ Validating configuration files..."
	@for config in $(CONFIGS); do \
		echo "Checking $$config..."; \
		python3 -c "import yaml; yaml.safe_load(open('$(CONFIG_DIR)/$$config'))" && \
			echo "  ✅ $$config is valid" || \
			echo "  ❌ $$config has syntax errors"; \
	done

# Log management
.PHONY: logs
logs: ## Show recent logs
	@echo "📋 Recent logs:"
	@if [ -d "$(LOGS_DIR)" ]; then \
		find $(LOGS_DIR) -name "*.log" -exec echo "=== {} ===" \; -exec tail -10 {} \; 2>/dev/null || \
		echo "No log files found"; \
	else \
		echo "Logs directory not found"; \
	fi

.PHONY: clean-logs
clean-logs: ## Clean log files
	@echo "🧹 Cleaning log files..."
	@rm -rf $(LOGS_DIR)/*.log 2>/dev/null || true
	@echo "✅ Log files cleaned"

# Development utilities
.PHONY: clean
clean: clean-logs ## Clean build artifacts and logs
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean completed"

.PHONY: dev
dev: clean setup run-dev ## Full development setup and run

# Docker commands
.PHONY: docker-build
docker-build: ## Build Docker image with version info
	@echo "🐳 Building Docker image..."
	@docker build --build-arg VERSION="$(VERSION)" \
		--build-arg COMMIT="$(COMMIT)" \
		--build-arg BUILD_DATE="$(BUILD_DATE)" \
		-t $(SERVICE_NAME):$(VERSION) .
	@echo "✅ Docker image built: $(SERVICE_NAME):$(VERSION)"

# Configuration testing
.PHONY: test-config-loading
test-config-loading: build ## Test configuration loading for all configs
	@echo "🧪 Testing configuration loading..."
	@for config in $(CONFIGS); do \
		echo "Testing $$config..."; \
		timeout 3s $(BUILD_DIR)/$(BINARY_NAME) $$config >/dev/null 2>&1 && \
			echo "  ✅ $$config loads successfully" || \
			echo "  ❌ $$config failed to load"; \
	done

.PHONY: info
info: ## Show project information
	@echo "📋 Project Information:"
	@echo "  Service Name: $(SERVICE_NAME)"
	@echo "  Version: $(VERSION)"
	@echo "  Commit: $(COMMIT)"
	@echo "  Branch: $(BRANCH)"
	@echo "  Build Date: $(BUILD_DATE)"
	@echo "  Binary: $(BUILD_DIR)/$(BINARY_NAME)"
	@echo ""
	@echo "📁 Directory Structure:"
	@echo "  Config: $(CONFIG_DIR)"
	@echo "  Logs: $(LOGS_DIR)"
	@echo "  Build: $(BUILD_DIR)"
	@echo ""
	@echo "🌐 Default Endpoints:"
	@echo "  Development: http://localhost:8083"
	@echo "  Production: http://localhost:8080"  
	@echo "  Testing: http://localhost:8084"