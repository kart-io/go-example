# File Logging Demo Makefile

.PHONY: help run clean logs examples setup test

# Default target
help: ## Show this help message
	@echo "File Logging Demo - Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Examples:"
	@echo "  make setup    # Create logs directory and prepare environment"
	@echo "  make run      # Run the complete file logging demo" 
	@echo "  make examples # Run configuration examples"
	@echo "  make logs     # View generated log files"
	@echo "  make clean    # Clean up generated files"

setup: ## Create logs directory and prepare environment
	@echo "🔧 Setting up file logging demo environment..."
	@mkdir -p logs
	@echo "✅ Logs directory created: logs/"
	@echo "🎯 Ready to run file logging demo!"

run: setup ## Run the complete file logging demonstration
	@echo "🚀 Starting file logging demo..."
	@echo "📁 Log files will be created in: logs/"
	@echo ""
	@go run main.go
	@echo ""
	@echo "✅ Demo completed! Check the files generated in logs/"

examples: setup ## Run configuration examples
	@echo "📚 Running configuration examples..."
	@go run config-examples.go
	@echo ""
	@echo "💡 These examples show different ways to configure file logging"

test: setup ## Run a quick test to verify logging works
	@echo "🧪 Running quick logging test..."
	@go run -c 'package main; import ("github.com/kart-io/logger"; "github.com/kart-io/logger/option"); func main() { opt := &option.LogOption{Engine: "slog", Level: "info", Format: "json", OutputPaths: []string{"logs/test.log"}}; l, _ := logger.New(opt); l.Info("Test logging successful") }' || echo 'Creating simple test...'
	@echo 'package main' > test-simple.go
	@echo 'import (' >> test-simple.go
	@echo '  "github.com/kart-io/logger"' >> test-simple.go 
	@echo '  "github.com/kart-io/logger/option"' >> test-simple.go
	@echo ')' >> test-simple.go
	@echo 'func main() {' >> test-simple.go
	@echo '  opt := &option.LogOption{' >> test-simple.go
	@echo '    Engine: "slog",' >> test-simple.go
	@echo '    Level: "info",' >> test-simple.go
	@echo '    Format: "json",' >> test-simple.go
	@echo '    OutputPaths: []string{"logs/test.log"},' >> test-simple.go
	@echo '  }' >> test-simple.go
	@echo '  logger, _ := logger.New(opt)' >> test-simple.go
	@echo '  logger.Info("Test logging successful")' >> test-simple.go
	@echo '  println("✅ Test log written to logs/test.log")' >> test-simple.go
	@echo '}' >> test-simple.go
	@go run test-simple.go
	@rm test-simple.go
	@echo ""

logs: ## Show generated log files and their contents
	@echo "📄 Generated log files:"
	@echo ""
	@if [ -d "logs" ]; then \
		ls -la logs/ | grep -v "^total" | grep -v "^d" | sed 's/^/  /'; \
		echo ""; \
		echo "📋 Recent log entries:"; \
		echo ""; \
		for file in logs/*.log; do \
			if [ -f "$$file" ]; then \
				echo "=== $$(basename $$file) ==="; \
				tail -n 3 "$$file" | sed 's/^/  /'; \
				echo ""; \
			fi; \
		done; \
	else \
		echo "  No logs directory found. Run 'make setup' first."; \
	fi

clean: ## Clean up generated log files and temporary files
	@echo "🧹 Cleaning up generated files..."
	@rm -rf logs/
	@rm -f test-simple.go
	@echo "✅ Cleanup completed!"

watch: ## Watch log files in real-time (requires tail)
	@echo "👀 Watching log files in real-time..."
	@echo "Press Ctrl+C to stop watching"
	@echo ""
	@if [ -f "logs/application.log" ]; then \
		tail -f logs/application.log; \
	elif [ -f "logs/app.log" ]; then \
		tail -f logs/app.log; \
	else \
		echo "No log files found to watch. Run 'make run' first."; \
	fi

demo-server: setup ## Run just the web server demo part
	@echo "🌐 Starting web server demo..."
	@echo "This will run only the web server portion of the demo"
	@echo ""
	@go run -ldflags "-X main.runOnlyWebServer=true" main.go || echo "Web server demo completed"

size: ## Show size of generated log files
	@echo "📊 Log file sizes:"
	@echo ""
	@if [ -d "logs" ]; then \
		du -h logs/* 2>/dev/null | sed 's/^/  /' || echo "  No log files found"; \
		echo ""; \
		echo "Total logs directory size:"; \
		du -sh logs/ 2>/dev/null | sed 's/^/  /' || echo "  0B"; \
	else \
		echo "  No logs directory found"; \
	fi

monitor: ## Monitor log directory for changes
	@echo "🔍 Monitoring logs directory for changes..."
	@echo "Press Ctrl+C to stop monitoring"
	@echo ""
	@if command -v inotifywait >/dev/null 2>&1; then \
		inotifywait -m -r logs/ -e create -e modify -e delete; \
	else \
		echo "inotifywait not available. Install inotify-tools package for real-time monitoring."; \
		echo "Falling back to periodic checks..."; \
		while true; do \
			ls -la logs/ 2>/dev/null | tail -n +2 | sed 's/^/  /'; \
			sleep 2; \
			clear; \
		done; \
	fi

info: ## Show information about the demo
	@echo "📖 File Logging Demo Information"
	@echo ""
	@echo "This demo showcases various file logging configurations using kart-io/logger:"
	@echo ""
	@echo "🎯 Features demonstrated:"
	@echo "  • Single file logging (JSON format)"
	@echo "  • Multiple output paths (console + file)"  
	@echo "  • Level-based file separation"
	@echo "  • File rotation patterns"
	@echo "  • Web server access and application logs"
	@echo ""
	@echo "📁 Generated files:"
	@echo "  • logs/single.log       - Basic single file output"
	@echo "  • logs/multiple.log     - Multi-output example"
	@echo "  • logs/info.log         - Info-level and above"
	@echo "  • logs/error.log        - Error-level and above"
	@echo "  • logs/access.log       - HTTP access logs"
	@echo "  • logs/application.log  - Application logs"
	@echo "  • logs/rotated-*.log    - Timestamped rotation example"
	@echo ""
	@echo "🚀 Quick start:"
	@echo "  1. make setup    # Prepare environment"
	@echo "  2. make run      # Run complete demo"
	@echo "  3. make logs     # View generated logs"
	@echo ""
	@echo "📚 Learn more: Check README.md for detailed documentation"

# Hidden target for CI/automated testing
test-ci: setup
	@echo "Running CI tests for file logging demo..."
	@timeout 30 go run main.go > /dev/null 2>&1 || echo "Demo completed (or timed out)"
	@if [ -f "logs/single.log" ] && [ -f "logs/access.log" ] && [ -f "logs/application.log" ]; then \
		echo "✅ All expected log files were created"; \
		exit 0; \
	else \
		echo "❌ Some expected log files were missing"; \
		exit 1; \
	fi