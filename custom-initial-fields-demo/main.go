package main

import (
	"fmt"
	"os"

	"github.com/kart-io/logger"
	"github.com/kart-io/logger/option"
	"github.com/kart-io/version"
)

func main() {
	fmt.Println("=== Custom Initial Fields Demo ===")
	fmt.Println("Demonstrating how all InitialFields are included in every log entry\n")

	// Get version info for some fields
	versionInfo := version.Get()

	// Demo: Logger with many different types of initial fields
	fmt.Println("Logger with various types of initial fields:")
	logOption := &option.LogOption{
		Engine:      "slog",
		Level:       "info",
		Format:      "json",
		OutputPaths: []string{"stdout"},
		InitialFields: map[string]interface{}{
			// Service information
			"service.name":    versionInfo.ServiceName,
			"service.version": versionInfo.GitVersion,
			
			// Environment information
			"environment": "production",
			"region":      "us-west-2",
			"datacenter":  "dc-1",
			
			// Team and ownership
			"team":       "platform",
			"squad":      "infrastructure",
			"maintainer": "john.doe@company.com",
			
			// Technical details
			"language":    "go",
			"framework":   "gin",
			"port":        8080,
			
			// Custom application fields
			"app_type":     "web-api",
			"health_check": true,
			"debug_mode":   false,
			
			// Business context
			"cost_center": "engineering",
			"project":     "customer-portal",
			
			// Infrastructure
			"container_id": getContainerID(),
			"node_name":    getNodeName(),
			
			// Compliance and governance
			"data_classification": "internal",
			"retention_days":      30,
		},
	}

	logger, err := logger.New(logOption)
	if err != nil {
		panic(err)
	}

	// All these log entries will include ALL the initial fields above
	fmt.Println("\n1. Simple info log:")
	logger.Info("Application started successfully")

	fmt.Println("\n2. Structured log with additional fields:")
	logger.Infow("User login", 
		"user_id", "user-12345",
		"ip_address", "192.168.1.100",
		"login_method", "oauth2",
	)

	fmt.Println("\n3. Error log:")
	logger.Errorw("Database connection failed",
		"error", "connection timeout",
		"retry_count", 3,
		"duration_ms", 5000,
	)

	fmt.Println("\n4. Debug log with nested data:")
	logger.Debugw("Processing request",
		"request_id", "req-789",
		"user_agent", "Mozilla/5.0...",
		"headers", map[string]string{
			"authorization": "Bearer ***",
			"content-type":  "application/json",
		},
	)

	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("\nNotice how EVERY log entry includes all the InitialFields:")
	fmt.Println("- Service information (name, version)")
	fmt.Println("- Environment details (region, datacenter)")
	fmt.Println("- Team ownership (team, squad, maintainer)")
	fmt.Println("- Technical context (language, framework, port)")
	fmt.Println("- Business context (cost_center, project)")
	fmt.Println("- Infrastructure details (container_id, node_name)")
	fmt.Println("- Compliance info (data_classification, retention_days)")
	fmt.Println("- Plus any runtime fields added via Infow(), Errorw(), etc.")
}

// Simulated functions to get infrastructure details
func getContainerID() string {
	if id := os.Getenv("CONTAINER_ID"); id != "" {
		return id
	}
	return "container-abc123"
}

func getNodeName() string {
	if name := os.Getenv("NODE_NAME"); name != "" {
		return name
	}
	return "node-worker-01"
}