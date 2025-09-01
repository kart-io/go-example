package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kart-io/logger"
	"github.com/kart-io/logger/option"
	"github.com/kart-io/version"
)

func main() {
	// Get version information
	versionInfo := version.Get()

	// Initialize logger with service information and OTLP export
	logOption := &option.LogOption{
		Engine:      "slog",
		Level:       "info",
		Format:      "json",
		OutputPaths: []string{"stdout"},
		// Add initial fields that should be in every log entry
		InitialFields: map[string]interface{}{
			"service.name":    versionInfo.ServiceName,   // Service name from build-time injection
			"service.version": versionInfo.GitVersion,    // Version from build-time injection  
			"commit":          versionInfo.GitCommit[:8], // Short commit hash
			"build_date":      versionInfo.BuildDate,
		},
		// Smart OTLP configuration - will auto-enable if endpoint is available
		OTLPEndpoint: "localhost:4317", // Jaeger default gRPC endpoint (no http:// prefix for gRPC)
		OTLP: &option.OTLPOption{
			// Basic OTLP configuration
		},
	}

	// Create logger with initial fields already included
	serviceLogger, err := logger.New(logOption)
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	// Log OTLP configuration status
	if logOption.OTLPEndpoint != "" {
		fmt.Printf("OTLP configured for endpoint: %s (connection may fail if collector is not running)\n", logOption.OTLPEndpoint)
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		serviceLogger.Infow("Handling root request", "endpoint", "/", "method", "GET")
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Go Example API",
			"version": versionInfo.GitVersion,
		})
	})

	r.GET("/health", func(c *gin.Context) {
		serviceLogger.Infow("Health check requested", "endpoint", "/health", "method", "GET")
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"version": versionInfo.GitVersion,
		})
	})

	r.GET("/version", func(c *gin.Context) {
		serviceLogger.Infow("Version info requested", "endpoint", "/version", "method", "GET")
		c.JSON(http.StatusOK, versionInfo)
	})

	// Log startup with all service information
	port := ":8082" // Default port
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = ":" + envPort
	}
	serviceLogger.Infow("Starting server",
		"port", port,
		"endpoints", []string{"/", "/health", "/version"},
		"go_version", versionInfo.GoVersion,
		"platform", versionInfo.Platform,
	)

	if err := r.Run(port); err != nil {
		serviceLogger.Fatalw("Failed to start server", "error", err.Error())
	}
}
