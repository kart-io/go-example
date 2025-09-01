# Viper Configuration Demo

A comprehensive demonstration of using **Viper** to read YAML configuration files and convert them to logger options. This package showcases configuration-driven logging with multiple environment support.

## Features

- **YAML Configuration Management**: Load logger settings from YAML files using Viper
- **Multi-Environment Support**: Development, production, and testing configurations
- **Environment Variable Override**: Configuration values can be overridden via environment variables
- **Dynamic Configuration**: Load different configs via command line or environment variables
- **Logger Integration**: Seamless conversion from YAML config to logger options
- **OTLP Support**: OpenTelemetry configuration through YAML files
- **Configuration Validation**: Built-in validation for all configuration values
- **API Endpoints**: RESTful endpoints to inspect configuration and test logging

## Project Structure

```
viper-config-demo/
├── config/
│   ├── app.yaml         # Development configuration
│   ├── production.yaml  # Production configuration
│   ├── testing.yaml     # Testing configuration
│   └── config.go        # Configuration structs and management
├── main.go              # Main application with Gin web server
├── Makefile            # Build and run commands
└── README.md           # This file
```

## Configuration Files

### app.yaml (Development)
- **Port**: 8083
- **Logger**: Slog engine with debug level and console/file output
- **OTLP**: Enabled with development settings

### production.yaml (Production)
- **Port**: 8080
- **Logger**: Zap engine with info level and structured logging
- **OTLP**: Enabled with production collector endpoint

### testing.yaml (Testing)
- **Port**: 8084
- **Logger**: Slog engine with warn level and minimal output
- **OTLP**: Disabled for testing

## Quick Start

### 1. Setup and Build

```bash
make setup    # Create directories and install dependencies
make build    # Build with version injection
```

### 2. Run with Different Configurations

```bash
# Development mode (app.yaml)
make run-dev

# Production mode (production.yaml)
make run-prod

# Testing mode (testing.yaml)
make run-test
```

### 3. Run with Environment Variables

```bash
# Use APP_ENV environment variable
make run-env-dev    # APP_ENV=app
make run-env-prod   # APP_ENV=production
make run-env-test   # APP_ENV=testing
```

## Configuration Loading Methods

### Method 1: Command Line Argument

```bash
./bin/viper-config-demo app.yaml
./bin/viper-config-demo production.yaml
./bin/viper-config-demo testing.yaml
```

### Method 2: Environment Variable

```bash
APP_ENV=production ./bin/viper-config-demo
APP_ENV=testing ./bin/viper-config-demo
```

### Method 3: Environment Variable Override

Override specific configuration values:

```bash
# Override server port
APP_SERVER_PORT=9000 ./bin/viper-config-demo

# Override logger level
APP_LOGGER_LEVEL=debug ./bin/viper-config-demo production.yaml

# Override OTLP endpoint
APP_LOGGER_OTLP_ENDPOINT=custom-collector:4317 ./bin/viper-config-demo
```

## API Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /` | Service information and configuration summary |
| `GET /health` | Health check with service status |
| `GET /version` | Build and version information |
| `GET /config` | Current configuration (sanitized) |
| `GET /logger/test` | Test all log levels and structured logging |
| `GET /debug/config` | Raw configuration (development only) |

### Test Endpoints

```bash
# Test development endpoints (port 8083)
make test-endpoints

# Test production endpoints (port 8080)
make test-prod-endpoints

# Test testing endpoints (port 8084)  
make test-testing-endpoints
```

## Configuration Structure

### Complete YAML Structure

```yaml
# Server configuration
server:
  port: 8083
  name: "viper-config-demo"
  environment: "development"

# Service information
service:
  name: "viper-config-api"
  version: "v1.0.0"
  description: "Viper configuration demo service"

# Logger configuration
logger:
  engine: "slog"              # "zap" or "slog"
  level: "debug"              # "debug", "info", "warn", "error", "fatal"
  format: "json"              # "json" or "console"
  development: true
  disable_caller: false
  disable_stacktrace: false
  output_paths:
    - "stdout"
    - "logs/app.log"
  otlp_endpoint: "localhost:4317"  # OTLP endpoint for OpenTelemetry export
  otlp:
    enabled: true
    endpoint: "localhost:4317"
    protocol: "grpc"            # "grpc" or "http"
    timeout: "10s"
    insecure: true
    # service_name and service_version handled via -ldflags injection
    headers:
      x-api-key: "demo-key"
      x-environment: "development"
```

### Environment Variable Mapping

Viper automatically maps environment variables with `APP_` prefix:

| Environment Variable | Configuration Path |
|---------------------|-------------------|
| `APP_SERVER_PORT` | `server.port` |
| `APP_LOGGER_ENGINE` | `logger.engine` |
| `APP_LOGGER_LEVEL` | `logger.level` |
| `APP_LOGGER_OTLP_ENABLED` | `logger.otlp.enabled` |
| `APP_LOGGER_OTLP_ENDPOINT` | `logger.otlp.endpoint` |

## Logger Integration

### Configuration to Logger Option Conversion

The `ConfigManager.ToLoggerOption()` method converts YAML configuration to logger options:

```go
// Load configuration
config, logOption, err := config.LoadConfigFromFile("app.yaml")
if err != nil {
    log.Fatal(err)
}

// Create logger with configuration
logger, err := logger.New(logOption)
if err != nil {
    log.Fatal(err)
}
```

### Features

- **Engine Selection**: Choose between Zap and Slog engines
- **Level Configuration**: Set log level via configuration
- **Output Paths**: Multiple output destinations (console, files)
- **OTLP Integration**: OpenTelemetry configuration
- **Service Context**: Automatic service name and version injection
- **Development Mode**: Enhanced debugging features

## Configuration Validation

### Built-in Validation Rules

- **Server Port**: Must be between 1-65535
- **Logger Engine**: Must be "zap" or "slog"
- **Logger Level**: Must be valid log level
- **Logger Format**: Must be "json" or "console"
- **OTLP Protocol**: Must be "grpc" or "http"
- **OTLP Timeout**: Must be valid duration format

### Validation Commands

```bash
make validate-configs   # Validate YAML syntax
make test-config-loading # Test configuration loading
```

## Development Commands

### Build and Run

```bash
make help              # Show all available commands
make build            # Build with version injection
make run              # Run development mode
make dev              # Clean, setup, and run development
```

### Configuration Management

```bash
make show-configs     # Display all config files
make validate-configs # Validate YAML files
make info            # Show project information
```

### Testing and Debugging

```bash
make test-endpoints           # Test API endpoints
make test-config-loading     # Test config loading
make logs                   # Show recent logs
make clean-logs            # Clean log files
```

### Docker

```bash
make docker-build     # Build Docker image with version info
```

## Environment Variables Reference

### Configuration Override Variables

```bash
# Server configuration
export APP_SERVER_PORT=8083
export APP_SERVER_NAME="viper-config-demo"
export APP_SERVER_ENVIRONMENT="development"

# Logger configuration  
export APP_LOGGER_ENGINE="slog"
export APP_LOGGER_LEVEL="debug"
export APP_LOGGER_FORMAT="json"
export APP_LOGGER_DEVELOPMENT="true"

# OTLP configuration
export APP_LOGGER_OTLP_ENABLED="true"
export APP_LOGGER_OTLP_ENDPOINT="localhost:4317"
export APP_LOGGER_OTLP_PROTOCOL="grpc"
```

### Runtime Environment Variables

```bash
# Environment selection
export APP_ENV="production"

# Kubernetes environment (auto-detected)
export POD_NAME="my-pod-123"
export POD_NAMESPACE="default"
export KUBERNETES_SERVICE_HOST="10.96.0.1"
```

## Integration with Logger Package

### Logger Option Creation

```go
// Load configuration and create logger option
appConfig, logOption, err := config.LoadConfigFromFile(configFile)
if err != nil {
    return err
}

// Create logger with loaded configuration
logger, err := logger.New(logOption)
if err != nil {
    return err
}
```

### Version Integration

```go
// Get version information
versionInfo := version.Get()

// Service info is automatically handled via version package
// No manual override needed - service information is injected at build time
// using -ldflags and retrieved from version.Get() in the logger initialization
```

## Best Practices

### Configuration Management

1. **Environment-Specific Configs**: Use separate YAML files for each environment
2. **Sensitive Data**: Use environment variables for secrets, not YAML files
3. **Validation**: Always validate configuration before using
4. **Defaults**: Provide reasonable defaults for all configuration values

### Logger Configuration

1. **Production Settings**: Use structured JSON format and appropriate log levels
2. **Development Settings**: Enable debug mode and console output
3. **File Output**: Configure log rotation for file outputs
4. **OTLP Integration**: Enable OpenTelemetry for observability

### Security Considerations

1. **Sanitized Output**: Remove sensitive data from API responses
2. **Environment Variables**: Use environment variables for API keys and secrets
3. **Access Control**: Restrict access to debug endpoints in production

## Troubleshooting

### Common Issues

1. **Configuration Not Found**: Ensure YAML files are in the correct directory
2. **Port Already in Use**: Check if another service is using the same port
3. **OTLP Connection Failed**: Verify OTLP collector endpoint and network connectivity
4. **Permission Denied**: Check file permissions for log directories

### Debug Commands

```bash
# Check configuration files
make show-configs

# Validate configuration syntax
make validate-configs

# Test configuration loading
make test-config-loading

# Check logs
make logs
```

## Examples

### Custom Configuration

Create your own configuration file:

```yaml
# custom.yaml
server:
  port: 9000
  environment: "custom"

logger:
  engine: "zap"
  level: "info"
  output_paths: ["stdout", "logs/custom.log"]

  otlp:
    enabled: false
```

```bash
./bin/viper-config-demo custom.yaml
```

### Environment Override Example

```bash
# Start with production config but override port and log level
APP_SERVER_PORT=9090 APP_LOGGER_LEVEL=debug ./bin/viper-config-demo production.yaml
```

This demonstrates the power of Viper's configuration management with environment variable override capabilities.