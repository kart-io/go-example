package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kart-io/logger"
	"github.com/kart-io/logger/option"
	"github.com/kart-io/version"
)

func main() {
	fmt.Println("=== Real-World Initial Fields Demo ===")
	fmt.Println("Web service with comprehensive initial fields\n")

	// Get version and environment info
	versionInfo := version.Get()
	
	// Create logger with comprehensive initial fields
	// These fields will appear in EVERY log entry
	logOption := &option.LogOption{
		Engine:      "slog",
		Level:       "info",
		Format:      "json",
		OutputPaths: []string{"stdout", "logs/app.log"},
		InitialFields: map[string]interface{}{
			// === Required service identification ===
			"service.name":    getEnvOrDefault("SERVICE_NAME", versionInfo.ServiceName),
			"service.version": getEnvOrDefault("SERVICE_VERSION", versionInfo.GitVersion),
			
			// === Environment context ===
			"environment": getEnvOrDefault("ENVIRONMENT", "development"),
			"region":      getEnvOrDefault("AWS_REGION", "us-west-2"),
			"az":          getEnvOrDefault("AWS_AZ", "us-west-2a"),
			
			// === Kubernetes/Container context ===
			"pod_name":      getEnvOrDefault("POD_NAME", "local-pod"),
			"node_name":     getEnvOrDefault("NODE_NAME", "local-node"),
			"namespace":     getEnvOrDefault("POD_NAMESPACE", "default"),
			"cluster":       getEnvOrDefault("CLUSTER_NAME", "local-cluster"),
			
			// === Application context ===
			"app_name":    "customer-api",
			"app_version": "v2.1.0",
			"go_version":  versionInfo.GoVersion,
			"build_date":  versionInfo.BuildDate,
			"commit":      versionInfo.GitCommit[:8],
			
			// === Team and ownership ===
			"team":         "platform",
			"squad":        "api-team",
			"owner":        "platform-team@company.com",
			"on_call":      getEnvOrDefault("ONCALL_CONTACT", "platform-oncall@company.com"),
			
			// === Business context ===
			"business_unit": "customer-success",
			"cost_center":   "engineering",
			"project":       "customer-portal-v2",
			
			// === Technical configuration ===
			"server_port":    getEnvOrDefault("PORT", "8080"),
			"log_level":      "info",
			"metrics_port":   "9090",
			"health_port":    "8081",
			
			// === Compliance and governance ===
			"data_classification": "confidential",
			"compliance_scope":    "pci-dss",
			"retention_policy":    "90-days",
			
			// === Monitoring tags ===
			"monitoring.team":        "platform",
			"monitoring.runbook":     "https://wiki.company.com/runbooks/customer-api",
			"monitoring.dashboard":   "https://grafana.company.com/d/customer-api",
			"monitoring.alert_level": "critical",
			
			// === Feature flags context ===
			"feature.new_auth":      true,
			"feature.rate_limiting": true,
			"feature.caching":       false,
			
			// === Performance context ===
			"max_connections": 1000,
			"timeout_seconds": 30,
			"workers":         4,
		},
	}

	// Create logger - all fields above will be in every log entry
	appLogger, err := logger.New(logOption)
	if err != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", err))
	}

	// Create Gin router
	r := gin.New()
	
	// Use our logger for Gin middleware
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Log HTTP requests with our structured logger
		appLogger.Infow("HTTP request",
			"method", param.Method,
			"path", param.Path,
			"status", param.StatusCode,
			"latency_ms", param.Latency.Milliseconds(),
			"client_ip", param.ClientIP,
			"user_agent", param.Request.UserAgent(),
		)
		return ""
	}))

	// Routes with different log scenarios
	r.GET("/", func(c *gin.Context) {
		// Business logic log - all InitialFields will be included
		appLogger.Infow("Homepage accessed",
			"user_type", "anonymous",
			"referrer", c.Request.Header.Get("Referer"),
		)
		
		c.JSON(http.StatusOK, gin.H{
			"message": "Customer API",
			"version": versionInfo.GitVersion,
			"status":  "healthy",
		})
	})

	r.GET("/users/:id", func(c *gin.Context) {
		userID := c.Param("id")
		
		// Simulate user lookup with detailed logging
		appLogger.Infow("User lookup started",
			"user_id", userID,
			"operation", "get_user",
			"cache_enabled", true,
		)
		
		// Simulate some business logic
		time.Sleep(10 * time.Millisecond)
		
		if userID == "123" {
			appLogger.Infow("User found",
				"user_id", userID,
				"user_status", "active",
				"last_login", "2025-09-01T10:30:00Z",
				"permission_level", "standard",
			)
			
			c.JSON(http.StatusOK, gin.H{
				"user_id": userID,
				"name":    "John Doe",
				"status":  "active",
			})
		} else {
			// Error case - still includes all InitialFields
			appLogger.Warnw("User not found",
				"user_id", userID,
				"lookup_duration_ms", 10,
				"searched_indexes", []string{"primary", "email", "username"},
			)
			
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
				"user_id": userID,
			})
		}
	})

	r.POST("/users", func(c *gin.Context) {
		// Simulate user creation with error handling
		appLogger.Infow("User creation started",
			"operation", "create_user",
			"request_size_bytes", c.Request.ContentLength,
		)
		
		// Simulate validation error
		appLogger.Errorw("User creation failed",
			"error", "email already exists",
			"validation_errors", []string{"email", "username"},
			"retry_recommended", true,
		)
		
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email already exists",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		// Health check with system status
		appLogger.Debugw("Health check performed",
			"check_type", "http",
			"response_time_ms", 1,
			"dependencies", map[string]string{
				"database": "healthy",
				"redis":    "healthy",
				"queue":    "healthy",
			},
		)
		
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Start the server
	port := getEnvOrDefault("PORT", "8080")
	
	appLogger.Infow("Server starting",
		"startup_time", time.Now().Format(time.RFC3339),
		"pid", os.Getpid(),
		"available_endpoints", []string{"/", "/users/:id", "/users", "/health"},
	)

	fmt.Printf("Starting server on port %s\n", port)
	fmt.Println("Try these endpoints:")
	fmt.Printf("  curl http://localhost:%s/\n", port)
	fmt.Printf("  curl http://localhost:%s/users/123\n", port)
	fmt.Printf("  curl http://localhost:%s/users/999\n", port)
	fmt.Printf("  curl -X POST http://localhost:%s/users\n", port)
	fmt.Printf("  curl http://localhost:%s/health\n", port)
	fmt.Println("\nNotice how EVERY log entry contains all the InitialFields!")

	if err := r.Run(":" + port); err != nil {
		appLogger.Fatalw("Server failed to start",
			"error", err.Error(),
			"port", port,
		)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}