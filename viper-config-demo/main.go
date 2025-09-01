package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kart-io/logger"
	"github.com/kart-io/logger/core"
	"github.com/kart-io/logger/option"
	"github.com/kart-io/version"

	"github.com/kart-io/go-example/viper-config-demo/config"
)

func main() {
	fmt.Println("=== Viper Configuration Demo ===")
	fmt.Printf("Starting application with configuration-driven logging\n\n")

	// Load configuration based on environment or command line argument
	var configFile string
	if len(os.Args) > 1 {
		configFile = os.Args[1]
		fmt.Printf("📁 Loading config from argument: %s\n", configFile)
	} else {
		env := os.Getenv("APP_ENV")
		if env == "" {
			env = "app"
		}
		configFile = env + ".yaml"
		fmt.Printf("📁 Loading config from environment (%s): %s\n", env, configFile)
	}

	// Load configuration and create logger option
	appConfig, logOption, err := config.LoadConfigFromFile(configFile)
	if err != nil {
		fmt.Printf("❌ Failed to load configuration: %v\n", err)
		fmt.Println("Available config files:")
		fmt.Println("  - app.yaml (development)")
		fmt.Println("  - production.yaml")
		fmt.Println("  - testing.yaml")
		fmt.Println("\nUsage: go run main.go [config-file]")
		fmt.Println("   or: APP_ENV=production go run main.go")
		os.Exit(1)
	}

	// Show loaded configuration
	fmt.Printf("✅ Configuration loaded successfully\n")
	fmt.Printf("   Engine: %s\n", logOption.Engine)
	fmt.Printf("   Level: %s\n", logOption.Level)
	fmt.Printf("   Format: %s\n", logOption.Format)
	fmt.Printf("   Output Paths: %v\n", logOption.OutputPaths)
	// Extract service information for display (from config file and version package)
	serviceName := appConfig.Service.Name
	serviceVersion := appConfig.Service.Version

	// Get version info for display (build-time injected values take precedence)
	versionInfo := version.Get()
	if versionInfo.ServiceName != "" {
		serviceName = versionInfo.ServiceName
	}
	if versionInfo.GitVersion != "" {
		serviceVersion = versionInfo.GitVersion
	}

	fmt.Printf("   Service: %s v%s\n", serviceName, serviceVersion)
	if logOption.IsOTLPEnabled() {
		fmt.Printf("   OTLP: %s (%s)\n", logOption.OTLPEndpoint, logOption.OTLP.Protocol)
	} else {
		fmt.Printf("   OTLP: disabled\n")
	}
	fmt.Println()

	// Service info is now handled automatically via version package and -ldflags injection
	// No need to override OTLP service fields manually

	// Add service info and additional context as initial fields using the new methods
	logOption.WithInitialFields(map[string]interface{}{
		"service.name":    versionInfo.ServiceName,
		"service.version": versionInfo.GitVersion,
		"config_file":     configFile,
		"environment":     appConfig.Server.Environment,
	}).AddInitialField("commit", getShortCommit(versionInfo.GitCommit)).
		AddInitialField("build_date", versionInfo.BuildDate)

	// Create logger with all initial fields
	serviceLogger, err := logger.New(logOption)
	if err != nil {
		fmt.Printf("❌ Failed to initialize logger with initial fields: %v\n", err)
		os.Exit(1)
	}

	// Log startup information
	serviceLogger.Infow("Application starting",
		"config_loaded", true,
		"server_port", appConfig.Server.Port,
		"logger_engine", logOption.Engine,
		"otlp_enabled", logOption.IsOTLPEnabled(),
	)

	// Setup Gin with appropriate mode
	if appConfig.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else if appConfig.Server.Environment == "testing" {
		gin.SetMode(gin.TestMode)
	}

	r := gin.Default()

	// Add middleware for request logging
	r.Use(loggingMiddleware(serviceLogger))

	// Routes
	r.GET("/", func(c *gin.Context) {
		serviceLogger.Infow("Handling root request", "endpoint", "/", "method", "GET")
		c.JSON(http.StatusOK, gin.H{
			"message":     "Viper Configuration Demo API",
			"service":     appConfig.Service.Name,
			"version":     appConfig.Service.Version,
			"environment": appConfig.Server.Environment,
			"config_file": configFile,
			"description": appConfig.Service.Description,
		})
	})

	r.GET("/health", func(c *gin.Context) {
		serviceLogger.Debugw("Health check requested", "endpoint", "/health")
		c.JSON(http.StatusOK, gin.H{
			"status":      "healthy",
			"service":     appConfig.Service.Name,
			"version":     appConfig.Service.Version,
			"environment": appConfig.Server.Environment,
			"uptime":      "running",
		})
	})

	r.GET("/version", func(c *gin.Context) {
		serviceLogger.Infow("Version info requested", "endpoint", "/version", "method", "GET")
		c.JSON(http.StatusOK, gin.H{
			"build_info":  versionInfo,
			"config_info": appConfig.Service,
			"server_info": appConfig.Server,
		})
	})

	r.GET("/config", func(c *gin.Context) {
		serviceLogger.Infow("Configuration info requested", "endpoint", "/config")

		// Return sanitized configuration (without sensitive data)
		sanitizedConfig := sanitizeConfig(appConfig)
		c.JSON(http.StatusOK, gin.H{
			"config":    sanitizedConfig,
			"loaded_from": configFile,
		})
	})

	r.GET("/logger/test", func(c *gin.Context) {
		serviceLogger.Infow("Logger test endpoint accessed", "endpoint", "/logger/test")

		// Test all log levels
		serviceLogger.Debug("This is a debug message")
		serviceLogger.Info("This is an info message")
		serviceLogger.Warn("This is a warning message")
		serviceLogger.Error("This is an error message (simulated)")

		// Test structured logging
		serviceLogger.Infow("Structured logging test",
			"user_id", "12345",
			"action", "test_logging",
			"timestamp", "2025-09-01T15:00:00Z",
			"success", true,
		)

		c.JSON(http.StatusOK, gin.H{
			"message": "Logger test completed",
			"levels_tested": []string{"debug", "info", "warn", "error"},
			"check": "See logs for output",
		})
	})

	// Environment-specific routes
	if appConfig.Server.Environment == "development" {
		r.GET("/debug/config", func(c *gin.Context) {
			serviceLogger.Debugw("Debug config endpoint accessed", "endpoint", "/debug/config")
			c.JSON(http.StatusOK, gin.H{
				"raw_config": appConfig,
				"log_option": sanitizeLogOption(logOption),
				"env_vars":   getRelevantEnvVars(),
			})
		})
	}

	// Start server
	port := ":" + strconv.Itoa(appConfig.Server.Port)

	serviceLogger.Infow("Starting server",
		"port", port,
		"environment", appConfig.Server.Environment,
		"endpoints", []string{"/", "/health", "/version", "/config", "/logger/test"},
		"logger_config", fmt.Sprintf("%s/%s/%s", logOption.Engine, logOption.Level, logOption.Format),
	)

	if err := r.Run(port); err != nil {
		serviceLogger.Fatalw("Failed to start server", "error", err.Error(), "port", port)
	}
}

// loggingMiddleware creates a Gin middleware for request logging
func loggingMiddleware(logger core.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Log request details
		logger.Infow("HTTP request processed",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		)
	}
}

// sanitizeConfig removes sensitive information from configuration
func sanitizeConfig(cfg *config.Config) *config.Config {
	sanitized := *cfg

	// Remove sensitive OTLP headers from logger config
	if sanitized.Logger.OTLP != nil && sanitized.Logger.OTLP.Headers != nil {
		sanitizedHeaders := make(map[string]string)
		for k, v := range sanitized.Logger.OTLP.Headers {
			if k == "x-api-key" || k == "authorization" {
				sanitizedHeaders[k] = "***REDACTED***"
			} else {
				sanitizedHeaders[k] = v
			}
		}
		// Create a copy of the OTLP config to avoid modifying the original
		otlpCopy := *sanitized.Logger.OTLP
		otlpCopy.Headers = sanitizedHeaders
		sanitized.Logger.OTLP = &otlpCopy
	}

	return &sanitized
}

// sanitizeLogOption removes sensitive information from log option
func sanitizeLogOption(opt *option.LogOption) map[string]interface{} {
	result := map[string]interface{}{
		"engine":             opt.Engine,
		"level":              opt.Level,
		"format":             opt.Format,
		"development":        opt.Development,
		"disable_caller":     opt.DisableCaller,
		"disable_stacktrace": opt.DisableStacktrace,
		"output_paths":       opt.OutputPaths,
		"otlp_enabled":       opt.IsOTLPEnabled(),
		"otlp_endpoint":      opt.OTLPEndpoint,
	}

	// Service info is handled via version package and -ldflags injection
	// Add version info for API response
	versionInfo := version.Get()
	result["service_name"] = versionInfo.ServiceName
	result["service_version"] = versionInfo.GitVersion

	return result
}

// getRelevantEnvVars returns relevant environment variables
func getRelevantEnvVars() map[string]string {
	envVars := map[string]string{}
	relevantVars := []string{
		"APP_ENV", "APP_SERVER_PORT", "APP_LOGGER_LEVEL",
		"APP_LOGGER_ENGINE", "APP_OTLP_ENABLED", "APP_OTLP_ENDPOINT",
	}

	for _, varName := range relevantVars {
		if value := os.Getenv(varName); value != "" {
			envVars[varName] = value
		}
	}

	return envVars
}

// isDefaultVersion checks if version is a default/placeholder value
func isDefaultVersion(version string) bool {
	defaultVersions := []string{
		"v0.0.0-master+$Format:%h$",
		"unknown",
		"dev",
		"",
	}

	for _, defaultVer := range defaultVersions {
		if version == defaultVer {
			return true
		}
	}
	return false
}

// getShortCommit returns the first 8 characters of a commit hash
func getShortCommit(commit string) string {
	if len(commit) >= 8 {
		return commit[:8]
	}
	return commit
}