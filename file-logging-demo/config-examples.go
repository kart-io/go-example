// +build ignore

package main

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/kart-io/logger"
	"github.com/kart-io/logger/option"
)

// ConfigExamples demonstrates various logger configuration patterns
// This file is meant for reference and learning purposes

func main() {
	fmt.Println("=== Logger Configuration Examples ===")
	fmt.Println("This file demonstrates different ways to configure file logging")
	fmt.Println()

	// Example 1: Basic file logging
	basicFileLogging()

	// Example 2: Advanced configuration
	advancedConfiguration()

	// Example 3: Production-ready setup
	productionSetup()

	// Example 4: Development environment
	developmentSetup()
}

// Example 1: Basic file logging
func basicFileLogging() {
	fmt.Println("1. Basic File Logging:")
	fmt.Println("```go")
	fmt.Println(`logOption := &option.LogOption{
    Engine:      "slog",
    Level:       "info",
    Format:      "json", 
    OutputPaths: []string{"logs/app.log"},
}

logger, err := logger.New(logOption)
if err != nil {
    panic(err)
}

logger.Info("Application started")`)
	fmt.Println("```")
	fmt.Println()
}

// Example 2: Advanced configuration
func advancedConfiguration() {
	fmt.Println("2. Advanced Configuration:")
	fmt.Println("```go")
	fmt.Println(`logOption := &option.LogOption{
    Engine:            "zap",
    Level:             "debug",
    Format:            "json",
    OutputPaths:       []string{"stdout", "logs/app.log", "logs/debug.log"},
    Development:       false,
    DisableCaller:     false,
    DisableStacktrace: true,
    OTLP: &option.OTLPOption{
        // ServiceName and ServiceVersion handled via -ldflags injection
    },
}`)
	fmt.Println("```")
	fmt.Println()
}

// Example 3: Production-ready setup
func productionSetup() {
	fmt.Println("3. Production Setup:")
	fmt.Println("```go")
	fmt.Println(`// Create timestamped log file
timestamp := time.Now().Format("20060102")
logFile := fmt.Sprintf("logs/prod-%s.log", timestamp)

logOption := &option.LogOption{
    Engine:            "zap",           // High performance
    Level:             "info",          // Appropriate level for prod
    Format:            "json",          // Structured logging
    OutputPaths:       []string{logFile},
    Development:       false,          // Production mode
    DisableCaller:     true,           // Better performance
    DisableStacktrace: true,           // Reduce log size
    
    // OTLP configuration for centralized logging
    OTLPEndpoint: "http://otel-collector:4317",
    OTLP: &option.OTLPOption{
        // ServiceName and ServiceVersion handled via -ldflags injection
    },
}`)
	fmt.Println("```")
	fmt.Println()
}

// Example 4: Development environment
func developmentSetup() {
	fmt.Println("4. Development Environment:")
	fmt.Println("```go")
	fmt.Println(`logOption := &option.LogOption{
    Engine:            "slog",          // Standard library
    Level:             "debug",         // Verbose for debugging
    Format:            "console",       // Human-readable
    OutputPaths:       []string{"stdout", "logs/dev.log"},
    Development:       true,            // Enable development features
    DisableCaller:     false,           // Show caller info
    DisableStacktrace: false,           // Show stack traces
    OTLP: &option.OTLPOption{
        ServiceName:       "dev-service",
        ServiceVersion:    "dev-build",
    },
}`)
	fmt.Println("```")
	fmt.Println()
}

// ShowRealExamples creates actual logger instances to demonstrate functionality
func ShowRealExamples() {
	// Ensure logs directory exists
	// os.MkdirAll("logs", 0755)

	fmt.Println("=== Live Configuration Examples ===")

	// Example 1: Single file JSON logging
	singleFileExample()

	// Example 2: Multi-output console + file
	multiOutputExample()

	// Example 3: Structured logging with context
	structuredLoggingExample()
}

func singleFileExample() {
	logOption := &option.LogOption{
		Engine:      "slog",
		Level:       "info",
		Format:      "json",
		OutputPaths: []string{filepath.Join("logs", "example-single.log")},
	}

	logger, err := logger.New(logOption)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}

	logger.Infow("Single file example", "example", 1, "type", "basic")
}

func multiOutputExample() {
	logOption := &option.LogOption{
		Engine:      "zap",
		Level:       "debug", 
		Format:      "console",
		OutputPaths: []string{"stdout", filepath.Join("logs", "example-multi.log")},
	}

	logger, err := logger.New(logOption)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}

	logger.Debugw("Multi-output example", "example", 2, "outputs", []string{"console", "file"})
}

func structuredLoggingExample() {
	logOption := &option.LogOption{
		Engine:      "slog",
		Level:       "info",
		Format:      "json",
		OutputPaths: []string{filepath.Join("logs", "example-structured.log")},
		OTLP: &option.OTLPOption{
			ServiceName:    "example-service",
			ServiceVersion: "v1.0.0",
		},
	}

	logger, err := logger.New(logOption)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		return
	}

	// Create structured logger with context
	structuredLogger := logger.With(
		"component", "user-service",
		"environment", "staging",
		"request_id", "req-12345",
	)

	structuredLogger.Infow("User operation completed",
		"user_id", "user-67890",
		"operation", "profile_update",
		"duration_ms", 150,
		"success", true,
	)
}

// FileRotationExample demonstrates file rotation patterns
func FileRotationExample() {
	fmt.Println("=== File Rotation Patterns ===")
	
	// Pattern 1: Daily rotation
	fmt.Println("Daily rotation:")
	fmt.Println("logs/app-20250901.log")
	fmt.Println("logs/app-20250902.log")
	fmt.Println()
	
	// Pattern 2: Hourly rotation
	fmt.Println("Hourly rotation:")
	fmt.Println("logs/app-2025090108.log")
	fmt.Println("logs/app-2025090109.log")
	fmt.Println()
	
	// Pattern 3: Size-based rotation
	fmt.Println("Size-based rotation:")
	fmt.Println("logs/app.log")
	fmt.Println("logs/app.log.1")
	fmt.Println("logs/app.log.2")
	fmt.Println()

	// Example implementation
	dailyLogFile := fmt.Sprintf("logs/app-%s.log", time.Now().Format("20060102"))
	hourlyLogFile := fmt.Sprintf("logs/app-%s.log", time.Now().Format("2006010215"))
	
	fmt.Printf("Today's log file would be: %s\n", dailyLogFile)
	fmt.Printf("This hour's log file would be: %s\n", hourlyLogFile)
}

// LogLevelExamples shows different log level configurations
func LogLevelExamples() {
	fmt.Println("=== Log Level Examples ===")
	
	levels := map[string]string{
		"debug": "Development debugging, very verbose",
		"info":  "General information, production default",
		"warn":  "Warning conditions, potential issues", 
		"error": "Error conditions, requires attention",
		"fatal": "Fatal errors, application will exit",
	}
	
	for level, description := range levels {
		fmt.Printf("%-8s: %s\n", level, description)
	}
	fmt.Println()
	
	fmt.Println("Usage guidelines:")
	fmt.Println("- Development: debug level for comprehensive logging")
	fmt.Println("- Staging: info level for detailed operational info")
	fmt.Println("- Production: info or warn level for performance")
	fmt.Println("- High-traffic: warn or error level to reduce volume")
}