package main

import (
	"fmt"

	"github.com/kart-io/logger"
	"github.com/kart-io/logger/option"
	"github.com/kart-io/version"
)

func main() {
	fmt.Println("=== Default Fields Demo ===")
	fmt.Println("Demonstrating logger behavior with and without InitialFields\n")

	// Demo 1: Logger without InitialFields - should show "unknown"
	fmt.Println("1. Logger without InitialFields (should show 'unknown' values):")
	basicOption := &option.LogOption{
		Engine:      "slog",
		Level:       "info",
		Format:      "json",
		OutputPaths: []string{"stdout"},
		// No InitialFields specified
	}

	basicLogger, err := logger.New(basicOption)
	if err != nil {
		panic(err)
	}
	basicLogger.Infow("Basic logger message", "test", "value1")
	fmt.Println()

	// Demo 2: Logger with empty InitialFields - should still show "unknown"
	fmt.Println("2. Logger with empty InitialFields (should show 'unknown' values):")
	emptyOption := &option.LogOption{
		Engine:        "slog",
		Level:         "info",
		Format:        "json",
		OutputPaths:   []string{"stdout"},
		InitialFields: map[string]interface{}{}, // Empty map
	}

	emptyLogger, err := logger.New(emptyOption)
	if err != nil {
		panic(err)
	}
	emptyLogger.Infow("Empty initial fields message", "test", "value2")
	fmt.Println()

	// Demo 3: Logger with partial InitialFields - should show mix of provided and "unknown"
	fmt.Println("3. Logger with partial InitialFields (service.name provided, service.version unknown):")
	partialOption := &option.LogOption{
		Engine:      "slog",
		Level:       "info",
		Format:      "json",
		OutputPaths: []string{"stdout"},
		InitialFields: map[string]interface{}{
			"service.name": "partial-demo-service",
			// service.version not provided - should be "unknown"
		},
	}

	partialLogger, err := logger.New(partialOption)
	if err != nil {
		panic(err)
	}
	partialLogger.Infow("Partial fields message", "test", "value3")
	fmt.Println()

	// Demo 4: Logger with complete InitialFields - should show all provided values
	fmt.Println("4. Logger with complete InitialFields (all values provided):")
	versionInfo := version.Get()
	completeOption := &option.LogOption{
		Engine:      "slog",
		Level:       "info",
		Format:      "json",
		OutputPaths: []string{"stdout"},
		InitialFields: map[string]interface{}{
			"service.name":    versionInfo.ServiceName,
			"service.version": versionInfo.GitVersion,
			"environment":     "demo",
		},
	}

	completeLogger, err := logger.New(completeOption)
	if err != nil {
		panic(err)
	}
	completeLogger.Infow("Complete fields message", "test", "value4")
	fmt.Println()

	// Demo 5: Logger with overridden default values
	fmt.Println("5. Logger with custom values overriding defaults:")
	customOption := &option.LogOption{
		Engine:      "slog",
		Level:       "info",
		Format:      "json",
		OutputPaths: []string{"stdout"},
		InitialFields: map[string]interface{}{
			"service.name":    "custom-service",
			"service.version": "custom-version",
			"team":           "platform",
		},
	}

	customLogger, err := logger.New(customOption)
	if err != nil {
		panic(err)
	}
	customLogger.Infow("Custom fields message", "test", "value5")
	fmt.Println()

	fmt.Println("=== Demo Complete ===")
	fmt.Println("Notice how 'service.name' and 'service.version' are always present,")
	fmt.Println("with 'unknown' as default when not explicitly provided.")
}