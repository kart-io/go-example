# File Logging Demo

这个示例展示了如何使用kart-io/logger库将日志输出到文件。包含多种文件日志配置场景和最佳实践。

## 功能特性

### 📁 支持的输出目标
- **单文件输出**: 将所有日志写入单个文件
- **多输出路径**: 同时输出到控制台和文件
- **分级文件**: 不同级别的日志写入不同文件
- **时间戳文件**: 支持文件轮转的时间戳命名

### 🌟 实际应用场景
- **Web服务器日志**: 访问日志和应用日志分离
- **业务操作日志**: 记录关键业务流程
- **错误追踪**: 专门的错误日志文件
- **开发调试**: 控制台+文件双输出

## 运行示例

### 基础运行
```bash
cd /home/hellotalk/code/go/src/github.com/kart-io/go-example/file-logging-demo
go run main.go
```

### 查看生成的日志文件
```bash
# 列出所有日志文件
ls -la logs/

# 查看特定日志文件
cat logs/single.log
cat logs/access.log
cat logs/application.log

# 实时监控日志文件
tail -f logs/application.log
```

## 示例说明

### Demo 1: 单文件日志
```go
logOption := &option.LogOption{
    Engine:      "slog",
    Level:       "info", 
    Format:      "json",
    OutputPaths: []string{"logs/single.log"}, // 指定日志文件路径
}
```
- ✅ 所有日志写入单个文件
- ✅ JSON格式便于后续解析
- ✅ 生产环境常用配置

### Demo 2: 多输出路径
```go
logOption := &option.LogOption{
    Engine:      "zap",
    Level:       "debug",
    Format:      "console", 
    OutputPaths: []string{"stdout", "logs/multiple.log"}, // 多输出
}
```
- ✅ 同时输出到控制台和文件
- ✅ 开发环境理想配置
- ✅ 便于调试和持久化存储

### Demo 3: 分级日志文件
```go
// Info级别日志
infoOption := &option.LogOption{
    Level:       "info",
    OutputPaths: []string{"logs/info.log"},
}

// Error级别日志  
errorOption := &option.LogOption{
    Level:       "error",
    OutputPaths: []string{"logs/error.log"},
}
```
- ✅ 不同级别日志分别存储
- ✅ 便于问题定位和监控
- ✅ 支持不同的处理策略

### Demo 4: 文件轮转模拟
```go
timestamp := time.Now().Format("20060102-150405")
logFile := fmt.Sprintf("logs/rotated-%s.log", timestamp)
```
- ✅ 时间戳文件命名
- ✅ 支持按时间轮转
- ✅ 避免单文件过大

### Demo 5: Web服务器日志
```go
// 访问日志
accessLogOption := &option.LogOption{
    OutputPaths: []string{"logs/access.log"},
}

// 应用日志
appLogOption := &option.LogOption{
    OutputPaths: []string{"stdout", "logs/application.log"},
}
```
- ✅ 访问日志和应用日志分离
- ✅ 自定义Gin中间件记录请求
- ✅ 结构化日志便于分析

## 日志文件格式

### JSON格式示例
```json
{
    "time": "2025-09-01T15:30:45.123456789+08:00",
    "level": "info",
    "msg": "User login",
    "service.name": "go-example-api",
    "service.version": "39b038f", 
    "pod": "hellotalk",
    "user_id": "12345",
    "ip": "192.168.1.100"
}
```

### Console格式示例
```
2025-09-01T15:30:45.123+08:00 INFO User login user_id=12345 ip=192.168.1.100
```

## 配置最佳实践

### 1. 生产环境配置
```go
logOption := &option.LogOption{
    Engine:         "zap",           // 高性能
    Level:          "info",          // 适中的日志级别
    Format:         "json",          // 结构化格式
    OutputPaths:    []string{"logs/app.log"},
    ServiceName:    "my-service",
    ServiceVersion: "v1.0.0",
}
```

### 2. 开发环境配置
```go
logOption := &option.LogOption{
    Engine:      "slog",             // 标准库
    Level:       "debug",            // 详细调试信息
    Format:      "console",          // 人类可读
    OutputPaths: []string{"stdout", "logs/dev.log"},
}
```

### 3. 高并发环境
```go
logOption := &option.LogOption{
    Engine:      "zap",              // 最高性能
    Level:       "warn",             // 减少日志量
    Format:      "json",             // 高效解析
    OutputPaths: []string{"logs/high-perf.log"},
}
```

## 文件管理建议

### 目录结构
```
logs/
├── access.log          # HTTP访问日志
├── application.log     # 应用程序日志
├── error.log           # 错误日志
├── debug.log           # 调试日志
└── rotated-20250901-153045.log  # 轮转日志
```

### 日志轮转
- 使用logrotate工具管理日志文件大小
- 按时间或大小进行轮转
- 保留适当数量的历史日志
- 定期清理过期日志

### 监控和告警
- 监控日志文件大小和磁盘空间
- 对ERROR级别日志设置告警
- 使用ELK Stack或类似工具分析日志
- 定期备份重要日志文件

## 性能考虑

### 文件I/O优化
- 使用缓冲写入减少系统调用
- 考虑异步日志写入
- 避免频繁的文件打开关闭
- 选择适当的日志级别

### 磁盘空间管理
- 设置合理的日志保留策略
- 实施自动日志清理
- 监控磁盘空间使用情况
- 考虑日志压缩存储

## 故障排查

### 常见问题
1. **权限问题**: 确保应用有写入日志目录的权限
2. **磁盘空间不足**: 监控并清理日志文件
3. **文件锁定**: 避免多进程同时写入同一文件
4. **路径不存在**: 确保日志目录存在

### 调试命令
```bash
# 检查文件权限
ls -la logs/

# 监控磁盘使用
df -h

# 实时查看日志
tail -f logs/application.log

# 搜索错误日志
grep -i error logs/*.log

# 统计日志数量
wc -l logs/*.log
```

## 扩展功能

这个示例可以进一步扩展：
- 集成日志轮转库（如lumberjack）
- 添加日志采样以减少高频日志
- 实现结构化字段验证
- 添加日志格式化模板
- 集成分布式追踪ID
- 支持日志加密存储