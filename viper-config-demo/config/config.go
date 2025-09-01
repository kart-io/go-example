package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/kart-io/logger/option"
)

// Config represents the complete application configuration
type Config struct {
	Server ServerConfig `mapstructure:"server" yaml:"server" json:"server"`
	Service ServiceConfig `mapstructure:"service" yaml:"service" json:"service"`
	Logger option.LogOption `mapstructure:"logger" yaml:"logger" json:"logger"`
}

// ServerConfig contains server-specific settings
type ServerConfig struct {
	Port        int    `mapstructure:"port" yaml:"port" json:"port"`
	Name        string `mapstructure:"name" yaml:"name" json:"name"`
	Environment string `mapstructure:"environment" yaml:"environment" json:"environment"`
}

// ServiceConfig contains service identification information
type ServiceConfig struct {
	Name        string `mapstructure:"name" yaml:"name" json:"name"`
	Version     string `mapstructure:"version" yaml:"version" json:"version"`
	Description string `mapstructure:"description" yaml:"description" json:"description"`
}



// ConfigManager manages configuration loading and conversion
type ConfigManager struct {
	viper  *viper.Viper
	config *Config
}

// NewConfigManager creates a new configuration manager
func NewConfigManager() *ConfigManager {
	v := viper.New()
	
	// Set configuration defaults
	setDefaults(v)
	
	return &ConfigManager{
		viper:  v,
		config: &Config{},
	}
}

// LoadConfig loads configuration from file and environment
func (cm *ConfigManager) LoadConfig(configPath string, configName string) (*Config, error) {
	v := cm.viper
	
	// Set configuration file details
	if configPath != "" {
		v.AddConfigPath(configPath)
	}
	v.AddConfigPath("./config")      // Default config directory
	v.AddConfigPath(".")             // Current directory
	v.AddConfigPath("/etc/app")      // System config directory
	
	if configName != "" {
		v.SetConfigName(configName)
	} else {
		// Try to load based on environment
		env := os.Getenv("APP_ENV")
		if env == "" {
			env = "app" // Default config file
		}
		v.SetConfigName(env)
	}
	
	v.SetConfigType("yaml")
	
	// Enable environment variable override
	v.AutomaticEnv()
	v.SetEnvPrefix("APP")           // Prefix for environment variables
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	
	// Try to read configuration file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	// Unmarshal configuration into struct
	if err := v.Unmarshal(cm.config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	// Validate configuration
	if err := cm.validateConfig(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	
	return cm.config, nil
}

// ToLoggerOption converts the configuration to logger.Option
func (cm *ConfigManager) ToLoggerOption() (*option.LogOption, error) {
	if cm.config == nil {
		return nil, fmt.Errorf("configuration not loaded")
	}
	
	loggerConfig := &cm.config.Logger
	
	// Service info is handled via version package and -ldflags injection
	// No need to set OTLP service fields from config
	
	return loggerConfig, nil
}

// GetConfig returns the loaded configuration
func (cm *ConfigManager) GetConfig() *Config {
	return cm.config
}

// GetViper returns the underlying viper instance for advanced usage
func (cm *ConfigManager) GetViper() *viper.Viper {
	return cm.viper
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.name", "viper-config-demo")
	v.SetDefault("server.environment", "development")
	
	// Service defaults
	v.SetDefault("service.name", "viper-config-api")
	v.SetDefault("service.version", "v1.0.0")
	v.SetDefault("service.description", "Viper configuration demo service")
	
	// Logger defaults (including OTLP)
	v.SetDefault("logger.engine", "slog")
	v.SetDefault("logger.level", "info")
	v.SetDefault("logger.format", "json")
	v.SetDefault("logger.development", false)
	v.SetDefault("logger.disable_caller", false)
	v.SetDefault("logger.disable_stacktrace", false)
	v.SetDefault("logger.output_paths", []string{"stdout"})
}

// validateConfig validates the loaded configuration
func (cm *ConfigManager) validateConfig() error {
	config := cm.config
	
	// Validate server config
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}
	
	// Validate logger config
	validEngines := map[string]bool{"zap": true, "slog": true}
	if !validEngines[config.Logger.Engine] {
		return fmt.Errorf("invalid logger engine: %s (must be 'zap' or 'slog')", config.Logger.Engine)
	}
	
	validLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true, "fatal": true,
	}
	if !validLevels[config.Logger.Level] {
		return fmt.Errorf("invalid logger level: %s", config.Logger.Level)
	}
	
	validFormats := map[string]bool{"json": true, "console": true}
	if !validFormats[config.Logger.Format] {
		return fmt.Errorf("invalid logger format: %s (must be 'json' or 'console')", config.Logger.Format)
	}
	
	// OTLP validation is handled by the logger package
	
	return nil
}

// LoadConfigFromFile is a convenience function to load config from a specific file
func LoadConfigFromFile(filePath string) (*Config, *option.LogOption, error) {
	cm := NewConfigManager()
	
	// Parse file path
	var configPath, configName string
	if strings.Contains(filePath, "/") {
		parts := strings.Split(filePath, "/")
		configName = strings.TrimSuffix(parts[len(parts)-1], ".yaml")
		configPath = strings.Join(parts[:len(parts)-1], "/")
	} else {
		configName = strings.TrimSuffix(filePath, ".yaml")
		configPath = "./config"
	}
	
	// Load configuration
	config, err := cm.LoadConfig(configPath, configName)
	if err != nil {
		return nil, nil, err
	}
	
	// Convert to logger option
	logOption, err := cm.ToLoggerOption()
	if err != nil {
		return nil, nil, err
	}
	
	return config, logOption, nil
}

// LoadConfigFromEnv loads configuration using environment-based file selection
func LoadConfigFromEnv() (*Config, *option.LogOption, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "app" // Default to app.yaml
	}
	
	return LoadConfigFromFile(env + ".yaml")
}