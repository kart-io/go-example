# Go Example - kart-io Logger & Version Integration

这个项目展示了如何在实际的Go应用中集成和使用kart-io生态系统的核心组件：

- **`kart-io/logger`** - 高性能双引擎（Zap/Slog）日志库
- **`kart-io/version`** - 版本管理和构建信息包

## 项目结构

```
go-example/
├── gin-demo/              # Gin web服务示例
│   └── main.go           # 基础的web API，展示OTLP集成
├── file-logging-demo/     # 文件日志示例
│   ├── main.go           # 完整的文件日志演示
│   ├── config-examples.go # 配置示例参考
│   ├── README.md         # 详细的文件日志文档
│   └── Makefile          # 便捷的命令工具
├── Dockerfile            # Docker容器化配置
├── Makefile              # 项目级别命令
└── go.mod                # Go模块定义
```

## 快速开始

### 运行Gin Web服务示例
```bash
cd gin-demo
make run

# 或者指定端口
PORT=8080 make run
```

访问端点：
- `http://localhost:8082/` - 主页
- `http://localhost:8082/health` - 健康检查
- `http://localhost:8082/version` - 版本信息

### 运行文件日志示例
```bash
cd file-logging-demo
make run       # 运行完整演示
make logs      # 查看生成的日志文件
make examples  # 查看配置示例
```

## 功能特性

### 🌐 Web服务集成 (gin-demo)
- **Gin框架集成**: 展示在web服务中使用logger
- **OTLP导出**: 自动将日志发送到OpenTelemetry Collector
- **版本信息**: 通过API端点暴露构建信息
- **结构化日志**: 使用统一的字段格式

### 📁 文件日志系统 (file-logging-demo)
- **多种输出模式**: 单文件、多文件、控制台+文件
- **分级日志**: 不同级别的日志分别存储
- **文件轮转**: 时间戳命名支持日志轮转
- **Web访问日志**: HTTP请求和应用日志分离
- **配置示例**: 生产和开发环境的最佳实践

## InitialFields 详解

`InitialFields` 是一个强大的功能，允许你在创建 logger 时定义一组字段，这些字段会自动包含在每个日志条目中。

### 💡 **核心概念**
- **一次定义，处处使用**: 在 `InitialFields` 中定义的字段会出现在每个日志条目中
- **自动默认值**: `service.name` 和 `service.version` 如果未提供，会自动设为 "unknown"
- **类型支持**: 支持所有 Go 基础类型和复杂类型
- **优先级**: 运行时通过 `Infow()` 等方法添加的字段会覆盖同名的 InitialFields

### 📊 **常见使用场景**

#### **服务标识**
```go
InitialFields: map[string]interface{}{
    "service.name":    "user-service",
    "service.version": "v2.1.0",
}
```

#### **环境信息**
```go
InitialFields: map[string]interface{}{
    "environment": "production",
    "region":      "us-west-2",
    "datacenter":  "dc-1",
    "cluster":     "prod-cluster",
}
```

#### **团队和所有权**
```go
InitialFields: map[string]interface{}{
    "team":        "platform",
    "squad":       "api-team", 
    "owner":       "john.doe@company.com",
    "cost_center": "engineering",
}
```

#### **技术上下文**
```go
InitialFields: map[string]interface{}{
    "language":    "go",
    "framework":   "gin", 
    "port":        8080,
    "go_version":  "1.21.0",
    "debug_mode":  false,
}
```

#### **业务上下文**
```go
InitialFields: map[string]interface{}{
    "project":          "customer-portal",
    "business_unit":    "sales",
    "compliance_scope": "pci-dss",
    "data_classification": "confidential",
}
```

#### **监控和告警**
```go
InitialFields: map[string]interface{}{
    "monitoring.team":      "sre",
    "monitoring.runbook":   "https://runbooks.company.com/user-service",
    "monitoring.dashboard": "https://grafana.company.com/d/user-service",
    "alert_level":          "critical",
}
```

### 🚀 **完整示例**
```go
// 真实生产环境的初始字段配置
versionInfo := version.Get()
logOption := &option.LogOption{
    Engine:      "zap",
    Level:       "info",
    Format:      "json", 
    OutputPaths: []string{"stdout", "logs/app.log"},
    InitialFields: map[string]interface{}{
        // === 必需的服务标识 ===
        "service.name":    versionInfo.ServiceName,
        "service.version": versionInfo.GitVersion,
        
        // === 环境上下文 ===
        "environment": os.Getenv("ENVIRONMENT"),
        "region":      os.Getenv("AWS_REGION"),
        "az":          os.Getenv("AWS_AZ"),
        
        // === 容器/K8s 上下文 ===
        "pod_name":   os.Getenv("POD_NAME"),
        "node_name":  os.Getenv("NODE_NAME"),
        "namespace":  os.Getenv("POD_NAMESPACE"),
        
        // === 应用上下文 ===
        "build_date": versionInfo.BuildDate,
        "commit":     versionInfo.GitCommit[:8],
        "port":       8080,
        
        // === 团队信息 ===
        "team":  "platform",
        "owner": "platform-team@company.com",
        
        // === 自定义字段 ===
        "feature_flags": map[string]bool{
            "new_auth":      true,
            "rate_limiting": false,
        },
    },
}

logger, _ := logger.New(logOption)

// 每个日志条目都会包含上述所有字段
logger.Infow("User login", "user_id", "12345")
logger.Errorw("Database error", "error", "connection timeout")
```

## 配置示例

### 基础文件日志配置
```go
// 获取版本信息
versionInfo := version.Get()

// 创建带初始字段的日志器
logOption := &option.LogOption{
    Engine:         "slog",
    Level:          "info", 
    Format:         "json",
    OutputPaths:    []string{"logs/app.log"},
    // 初始字段会添加到每个日志条目中
    // 如果不提供 service.name 和 service.version，将默认为 "unknown"
    InitialFields: map[string]interface{}{
        "service.name":    versionInfo.ServiceName,     // 构建时注入
        "service.version": versionInfo.GitVersion,      // 构建时注入
    },
}

serviceLogger, err := logger.New(logOption)
if err != nil {
    log.Fatal(err)
}
```

#### 默认值行为
```go
// 如果不设置 InitialFields，会自动添加默认值
logOption := &option.LogOption{
    Engine:      "slog",
    Level:       "info",
    Format:      "json",
    OutputPaths: []string{"stdout"},
    // 不设置 InitialFields
}

logger, _ := logger.New(logOption)
logger.Info("Hello")
// 输出会包含: "service.name":"unknown", "service.version":"unknown"
```

#### 自定义初始字段
`InitialFields` 中的所有字段都会出现在每个日志条目中：
```go
logOption := &option.LogOption{
    Engine:      "slog",
    Level:       "info",
    Format:      "json",
    OutputPaths: []string{"stdout"},
    InitialFields: map[string]interface{}{
        "service.name":    "my-api",
        "service.version": "v1.2.0",
        "environment":     "production",
        "region":          "us-west-2",
        "team":           "platform",
        "datacenter":     "dc-1",
        "cost_center":    "engineering",
        "debug_mode":     false,
        "port":           8080,
    },
}

logger, _ := logger.New(logOption)
logger.Infow("User login", "user_id", "12345")

// 输出将包含所有 InitialFields + 运行时字段:
// {
//   "service.name": "my-api",
//   "service.version": "v1.2.0", 
//   "environment": "production",
//   "region": "us-west-2",
//   "team": "platform",
//   "datacenter": "dc-1",
//   "cost_center": "engineering",
//   "debug_mode": false,
//   "port": 8080,
//   "user_id": "12345",
//   "msg": "User login",
//   ...
// }
```

### 生产环境配置
```go
// 获取版本信息
versionInfo := version.Get()

// 生产环境日志器配置
logOption := &option.LogOption{
    Engine:            "zap",              // 高性能
    Level:             "info",             // 适中日志级别
    Format:            "json",             // 结构化格式
    OutputPaths:       []string{"logs/prod.log"},
    Development:       false,              // 生产模式
    DisableCaller:     true,               // 提升性能
    DisableStacktrace: true,               // 减少日志大小
    OTLPEndpoint:      "otel-collector:4317",  // gRPC端点
    // 生产环境初始字段
    InitialFields: map[string]interface{}{
        "service.name":    versionInfo.ServiceName,
        "service.version": versionInfo.GitVersion,
        "commit":          versionInfo.GitCommit[:8],
        "build_date":      versionInfo.BuildDate,
        "environment":     "production",
    },
}

serviceLogger, err := logger.New(logOption)
if err != nil {
    log.Fatal(err)
}

// 构建命令示例:
// go build -ldflags "-X 'github.com/kart-io/version.serviceName=my-service' -X 'github.com/kart-io/version.gitVersion=v1.0.0'"
```

### 开发环境配置
```go
logOption := &option.LogOption{
    Engine:      "slog",                   // 标准库
    Level:       "debug",                  // 详细调试信息
    Format:      "console",                // 人类可读
    OutputPaths: []string{"stdout", "logs/dev.log"},
    Development: true,                     // 开发模式特性
}
```

## 版本信息集成

项目使用`kart-io/version`包提供完整的构建信息：

```go
versionInfo := version.Get()

// 在logger创建时使用版本信息作为初始字段
logOption := &option.LogOption{
    // ... 其他配置 ...
    InitialFields: map[string]interface{}{
        "service.name":    versionInfo.ServiceName,      // 标准字段名
        "service.version": versionInfo.GitVersion,       // 标准字段名
        "commit":          versionInfo.GitCommit[:8],
        "build_date":      versionInfo.BuildDate,
    },
}

serviceLogger, err := logger.New(logOption)

// 在API中暴露版本信息
r.GET("/version", func(c *gin.Context) {
    c.JSON(http.StatusOK, versionInfo)
})
```

## 构建和部署

### 本地构建
```bash
make build    # 构建应用（带版本信息注入）
make run      # 运行应用
make version  # 显示版本信息
```

### Docker构建
```bash
make docker-build    # 构建Docker镜像
make docker-run      # 运行Docker容器
```

### 版本信息注入
构建时自动注入Git版本信息：
```bash
go build -ldflags "
  -X 'github.com/kart-io/version.serviceName=go-example-api'
  -X 'github.com/kart-io/version.gitVersion=$(git describe --tags)'
  -X 'github.com/kart-io/version.gitCommit=$(git rev-parse HEAD)'
  -X 'github.com/kart-io/version.buildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)'
" ./gin-demo
```

## OTLP集成

项目支持将日志导出到OpenTelemetry生态系统：

### 启动OTLP测试环境
```bash
# 在logger项目中启动OTLP栈
cd ../logger/otlp-docker
./deploy.sh

# 访问VictoriaLogs查看日志
curl 'http://127.0.0.1:9428/select/logsql/query?query=*&limit=10'
```

### OTLP配置
```go
// 获取版本信息
versionInfo := version.Get()

// OTLP配置 - 服务信息通过 InitialFields 添加
logOption := &option.LogOption{
    // ... 其他配置 ...
    OTLPEndpoint: "localhost:4317",  // gRPC端点
    InitialFields: map[string]interface{}{
        "service.name":    versionInfo.ServiceName,
        "service.version": versionInfo.GitVersion,
    },
    OTLP: &option.OTLPOption{
        Enabled: &[]bool{true}[0],  // 显式启用
    },
}

serviceLogger, _ := logger.New(logOption)
```

系统会自动：
- 检测OTLP端点是否可用
- 发送结构化日志到collector
- 添加服务标识和环境信息
- 支持Kubernetes环境检测

## 监控和可观测性

### 日志查询
生产环境中的日志查询示例：
```bash
# 查看服务日志
curl 'http://victorialogs:9428/select/logsql/query?query=service.name:go-example-api'

# 错误日志过滤
curl 'http://victorialogs:9428/select/logsql/query?query=level:error'

# 特定时间范围
curl 'http://victorialogs:9428/select/logsql/query?query=_time:2025-09-01'
```

### 指标和追踪
- 日志自动包含请求ID和追踪信息
- 支持分布式追踪上下文传递
- 与Prometheus和Jaeger集成

## 最佳实践

### 日志配置
1. **生产环境**: 使用Zap引擎，JSON格式，适中的日志级别
2. **开发环境**: 使用Slog引擎，Console格式，详细的调试信息
3. **高并发**: 减少日志级别，禁用caller和stacktrace
4. **调试**: 启用所有调试特性，使用多输出路径

### 版本管理
1. 使用Git标签进行版本控制
2. 构建时自动注入版本信息
3. 在日志和API中包含版本信息
4. 监控不同版本的性能差异

### 文件管理
1. 使用logrotate管理日志文件大小
2. 设置适当的日志保留策略
3. 监控磁盘空间使用情况
4. 定期备份重要日志文件

## 故障排查

### 常见问题
- **OTLP连接失败**: 检查collector是否运行，端点配置是否正确
- **文件权限问题**: 确保应用有写入日志目录的权限
- **版本信息为空**: 检查构建时的ldflags参数
- **日志格式不一致**: 验证引擎配置和字段映射

### 调试命令
```bash
# 检查日志文件
tail -f logs/app.log

# 验证OTLP连接
curl http://localhost:13133/  # Agent健康检查

# 测试版本注入
./app --version
```

## 扩展示例

这些示例展示了kart-io生态系统的核心功能，您可以：
- 扩展更多web框架集成（Echo、Fiber等）
- 添加数据库日志集成（GORM等）
- 实现自定义日志中间件
- 集成更多可观测性工具

## 相关链接

- [kart-io/logger](../logger/) - 日志库核心功能
- [kart-io/version](../version/) - 版本管理包
- [OTLP文档](../logger/otlp/) - OpenTelemetry集成
- [配置文档](../logger/option/) - 详细配置选项