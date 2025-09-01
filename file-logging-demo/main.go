package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kart-io/logger"
	"github.com/kart-io/logger/option"
	"github.com/kart-io/version"
)

func main() {
	// Get version information
	versionInfo := version.Get()

	fmt.Println("=== File Logging Demo ===")
	fmt.Printf("Service: %s\n", versionInfo.ServiceName)
	fmt.Printf("Version: %s\n", versionInfo.GitVersion)
	fmt.Printf("Build Date: %s\n", versionInfo.BuildDate)
	fmt.Println()

	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create logs directory: %v", err))
	}

	// Demo 1: Single file logging
	fmt.Println("=== Demo 1: Single File Logging ===")
	singleFileDemo(versionInfo)

	// Demo 2: Multiple output paths (console + file)
	fmt.Println("\n=== Demo 2: Multiple Output Paths ===")
	multipleOutputDemo(versionInfo)

	// Demo 3: Different log levels to different files
	fmt.Println("\n=== Demo 3: Level-based File Logging ===")
	levelBasedDemo(versionInfo)

	// Demo 4: File rotation simulation
	fmt.Println("\n=== Demo 4: File Rotation Simulation ===")
	fileRotationDemo(versionInfo)

	// Demo 5: Web server with file logging
	fmt.Println("\n=== Demo 5: Web Server with File Logging ===")
	webServerDemo(versionInfo)
}

// Demo 1: Log to a single file
func singleFileDemo(versionInfo version.Info) {
	logFile := filepath.Join("logs", "single.log")
	
	logOption := &option.LogOption{
		Engine:      "slog",
		Level:       "info",
		Format:      "json",
		OutputPaths: []string{logFile},
		// Add service info as initial fields
		InitialFields: map[string]interface{}{
			"service.name":    versionInfo.ServiceName,
			"service.version": versionInfo.GitVersion,
		},
		OTLP: &option.OTLPOption{
			// Basic OTLP configuration
		},
	}

	logger, err := logger.New(logOption)
	if err != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", err))
	}

	// Log some messages
	logger.Info("Single file logging demo started")
	logger.Infow("User login", "user_id", "12345", "ip", "192.168.1.100")
	logger.Warnw("High memory usage", "usage", "85%", "threshold", "80%")
	logger.Errorw("Database connection failed", "error", "connection timeout", "retry_count", 3)

	fmt.Printf("✅ Logs written to: %s\n", logFile)
	
	// Show file contents
	if content, err := os.ReadFile(logFile); err == nil {
		fmt.Printf("📄 File contents (last 200 chars):\n")
		if len(content) > 200 {
			fmt.Printf("...%s", content[len(content)-200:])
		} else {
			fmt.Printf("%s", content)
		}
	}
}

// Demo 2: Log to both console and file
func multipleOutputDemo(versionInfo version.Info) {
	logFile := filepath.Join("logs", "multiple.log")
	
	logOption := &option.LogOption{
		Engine:      "zap",
		Level:       "debug",
		Format:      "console",
		OutputPaths: []string{"stdout", logFile},
		OTLP: &option.OTLPOption{
			// ServiceName and ServiceVersion removed - handled via -ldflags injection
		},
	}

	coreLogger, err := logger.New(logOption)
	if err != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", err))
	}

	// Create a service-specific logger with service info
	logger := coreLogger.With(
		"service.name", versionInfo.ServiceName,
		"service.version", versionInfo.GitVersion,
	)
	serviceLogger := logger.With(
		"component", "payment-service",
		"environment", "production",
	)

	// Log messages that will appear both in console and file
	fmt.Println("📺 Watch the console output while logs are also written to file:")
	serviceLogger.Debug("Payment processing started")
	serviceLogger.Info("Payment validation successful")
	serviceLogger.Warn("Payment amount exceeds daily limit")
	serviceLogger.Error("Payment gateway error")

	fmt.Printf("✅ Logs written to both console and: %s\n", logFile)
}

// Demo 3: Different log levels to different files
func levelBasedDemo(versionInfo version.Info) {
	// Create separate loggers for different levels
	infoLogFile := filepath.Join("logs", "info.log")
	errorLogFile := filepath.Join("logs", "error.log")

	// Info level logger (info and above)
	infoOption := &option.LogOption{
		Engine:      "slog",
		Level:       "info",
		Format:      "json",
		OutputPaths: []string{infoLogFile},
		OTLP: &option.OTLPOption{
			// ServiceName and ServiceVersion removed - handled via -ldflags injection
		},
	}

	coreInfoLogger, err := logger.New(infoOption)
	if err != nil {
		panic(fmt.Sprintf("Failed to create info logger: %v", err))
	}

	// Add service info
	infoLogger := coreInfoLogger.With(
		"service.name", versionInfo.ServiceName,
		"service.version", versionInfo.GitVersion,
	)

	// Error level logger (error and above)
	errorOption := &option.LogOption{
		Engine:      "slog",
		Level:       "error",
		Format:      "json",
		OutputPaths: []string{errorLogFile},
		OTLP: &option.OTLPOption{
			// ServiceName and ServiceVersion removed - handled via -ldflags injection
		},
	}

	coreErrorLogger, err := logger.New(errorOption)
	if err != nil {
		panic(fmt.Sprintf("Failed to create error logger: %v", err))
	}

	// Add service info
	errorLogger := coreErrorLogger.With(
		"service.name", versionInfo.ServiceName,
		"service.version", versionInfo.GitVersion,
	)

	// Log different levels
	infoLogger.Info("Application started successfully")
	infoLogger.Warn("Configuration file not found, using defaults")
	infoLogger.Error("Failed to connect to database")

	errorLogger.Error("Critical system error")
	// Note: Fatal() would exit the program, so we use Error() instead for demo
	errorLogger.Error("System shutdown due to critical error (simulated fatal)")

	fmt.Printf("✅ Info logs written to: %s\n", infoLogFile)
	fmt.Printf("✅ Error logs written to: %s\n", errorLogFile)
}

// Demo 4: Simulate file rotation by creating timestamped files
func fileRotationDemo(versionInfo version.Info) {
	timestamp := time.Now().Format("20060102-150405")
	logFile := filepath.Join("logs", fmt.Sprintf("rotated-%s.log", timestamp))

	logOption := &option.LogOption{
		Engine:      "zap",
		Level:       "info",
		Format:      "json",
		OutputPaths: []string{logFile},
		OTLP: &option.OTLPOption{
			// ServiceName and ServiceVersion removed - handled via -ldflags injection
		},
	}

	coreLogger, err := logger.New(logOption)
	if err != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", err))
	}

	// Add service info
	logger := coreLogger.With(
		"service.name", versionInfo.ServiceName,
		"service.version", versionInfo.GitVersion,
	)

	// Simulate some business operations
	operations := []string{
		"user_registration",
		"order_creation", 
		"payment_processing",
		"inventory_update",
		"email_notification",
	}

	for i, op := range operations {
		logger.Infow("Business operation", 
			"operation", op,
			"step", i+1,
			"timestamp", time.Now().Unix(),
		)
		time.Sleep(100 * time.Millisecond) // Simulate processing time
	}

	fmt.Printf("✅ Timestamped logs written to: %s\n", logFile)
}

// Demo 5: Web server with comprehensive file logging
func webServerDemo(versionInfo version.Info) {
	// Create logs for different components
	accessLogFile := filepath.Join("logs", "access.log")
	appLogFile := filepath.Join("logs", "application.log")

	// Access logger (all requests)
	accessLogOption := &option.LogOption{
		Engine:      "slog",
		Level:       "info",
		Format:      "json",
		OutputPaths: []string{accessLogFile},
		OTLP: &option.OTLPOption{
			// ServiceName and ServiceVersion removed - handled via -ldflags injection
		},
	}

	coreAccessLogger, err := logger.New(accessLogOption)
	if err != nil {
		panic(fmt.Sprintf("Failed to create access logger: %v", err))
	}

	// Add service info
	accessLogger := coreAccessLogger.With(
		"service.name", versionInfo.ServiceName,
		"service.version", versionInfo.GitVersion,
	)

	// Application logger (console + file for development)
	appLogOption := &option.LogOption{
		Engine:      "slog",
		Level:       "debug",
		Format:      "console",
		OutputPaths: []string{"stdout", appLogFile},
		OTLP: &option.OTLPOption{
			// ServiceName and ServiceVersion removed - handled via -ldflags injection
		},
	}

	coreAppLogger, err := logger.New(appLogOption)
	if err != nil {
		panic(fmt.Sprintf("Failed to create app logger: %v", err))
	}

	// Add service info
	appLogger := coreAppLogger.With(
		"service.name", versionInfo.ServiceName,
		"service.version", versionInfo.GitVersion,
	)

	// Create loggers with context
	accessLoggerWithContext := accessLogger.With("component", "http-access")
	appLoggerWithContext := appLogger.With("component", "application")

	// Set up Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Custom logging middleware
	r.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		// Log access information
		accessLoggerWithContext.Infow("HTTP request",
			"method", method,
			"path", path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		)
	})

	// Routes
	r.GET("/", func(c *gin.Context) {
		appLoggerWithContext.Debug("Handling root request")
		c.JSON(http.StatusOK, gin.H{
			"message": "File Logging Demo API",
			"version": versionInfo.GitVersion,
			"logs": map[string]string{
				"access": accessLogFile,
				"app":    appLogFile,
			},
		})
	})

	r.GET("/health", func(c *gin.Context) {
		appLoggerWithContext.Debug("Health check requested")
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	r.GET("/error", func(c *gin.Context) {
		appLoggerWithContext.Error("Simulated error endpoint accessed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Simulated error"})
	})

	r.GET("/logs", func(c *gin.Context) {
		appLoggerWithContext.Debug("Log files listing requested")
		
		// List all log files
		logFiles := []string{}
		files, err := filepath.Glob(filepath.Join("logs", "*.log"))
		if err == nil {
			for _, file := range files {
				logFiles = append(logFiles, filepath.Base(file))
			}
		}
		
		c.JSON(http.StatusOK, gin.H{
			"log_files": logFiles,
			"logs_dir": "logs/",
		})
	})

	// Start server in background
	srv := &http.Server{
		Addr:    ":8084",
		Handler: r,
	}

	go func() {
		appLoggerWithContext.Infow("Starting web server",
			"port", 8084,
			"access_log", accessLogFile,
			"app_log", appLogFile,
		)
		
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLoggerWithContext.Errorw("Server failed to start", "error", err)
		}
	}()

	// Wait a moment for server to start
	time.Sleep(500 * time.Millisecond)

	fmt.Printf("🚀 Web server started on http://localhost:8084\n")
	fmt.Printf("📝 Access logs: %s\n", accessLogFile)
	fmt.Printf("📱 App logs: %s\n", appLogFile)
	fmt.Println()
	fmt.Println("Available endpoints:")
	fmt.Println("  GET http://localhost:8084/        - Main page")
	fmt.Println("  GET http://localhost:8084/health  - Health check")
	fmt.Println("  GET http://localhost:8084/error   - Simulate error")
	fmt.Println("  GET http://localhost:8084/logs    - List log files")
	fmt.Println()
	fmt.Println("Making some test requests...")

	// Make some test requests
	testEndpoints := []string{
		"http://localhost:8084/",
		"http://localhost:8084/health",
		"http://localhost:8084/logs",
		"http://localhost:8084/error",
	}

	client := &http.Client{Timeout: 2 * time.Second}
	for i, endpoint := range testEndpoints {
		appLoggerWithContext.Debugw("Making test request", "endpoint", endpoint, "request", i+1)
		
		resp, err := client.Get(endpoint)
		if err != nil {
			appLoggerWithContext.Warnw("Test request failed", "endpoint", endpoint, "error", err)
		} else {
			resp.Body.Close()
			appLoggerWithContext.Infow("Test request completed", "endpoint", endpoint, "status", resp.StatusCode)
		}
		time.Sleep(200 * time.Millisecond)
	}

	// Graceful shutdown
	fmt.Println("\nShutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		appLoggerWithContext.Errorw("Server shutdown failed", "error", err)
	} else {
		appLoggerWithContext.Info("Server shutdown completed")
	}

	fmt.Printf("✅ Demo completed. Check log files in the 'logs/' directory\n")
}