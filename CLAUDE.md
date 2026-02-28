# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在此代码仓库中工作时提供指导。

## 构建命令

```bash
# 构建当前平台（linux-amd64 最常用）
make linux-amd64

# 构建所有主要平台
make all

# 构建特定平台
make darwin-arm64    # macOS Apple Silicon
make linux-arm64     # Linux ARM64
make windows-amd64   # Windows 64位
```

## 运行服务器

```bash
# 使用默认配置运行 (etc/ppanel.yaml)
./bin/ppanel-server-linux-amd64 run

# 使用自定义配置运行
./bin/ppanel-server-linux-amd64 run --config /path/to/config.yaml

# 显示版本
./bin/ppanel-server-linux-amd64 version
```

## 代码生成

本项目使用自定义的 go-zero 风格 API 代码生成器。修改 `.api` 文件后：

```bash
# 从 API 定义生成代码
./script/generate.sh
```

生成器读取 `ppanel.api`（该文件导入所有子 API）并生成：
- `internal/handler/` - HTTP 处理器和路由
- `internal/types/` - 请求/响应类型

## 测试

```bash
# 运行所有测试
go test ./...

# 运行测试并显示详细输出
go test -v ./...

# 运行特定包的测试
go test -v ./pkg/cache/...

# 运行特定测试
go test -v -run TestName ./path/to/package
```

## 架构概览

PPanel 是一个使用 Go 构建的代理面板管理后端，采用分层架构：

### 服务层
应用程序运行三个并发服务（参见 `cmd/run.go`）：
- **HTTP 服务器** (`internal/`) - REST API 处理
- **队列服务** (`queue/`) - 通过 Asynq 进行异步任务处理
- **调度器** (`scheduler/`) - 基于 Cron 的定时任务

### 请求流程
```
API 定义 (.api 文件)
    ↓ generate.sh
Handler (internal/handler/) → Middleware → Logic (internal/logic/) → Model (internal/model/)
```

### 主要目录

| 目录 | 用途 |
|------|------|
| `apis/` | go-zero 格式的 API 定义。主入口：`ppanel.api` |
| `internal/handler/` | HTTP 处理器（从 .api 文件自动生成） |
| `internal/logic/` | 按 API 区域组织的业务逻辑（admin、public、auth 等） |
| `internal/model/` | GORM 数据模型，用于数据库操作 |
| `internal/svc/` | ServiceContext，持有 DB、Redis 和所有模型 |
| `internal/middleware/` | HTTP 中间件（认证、日志、CORS 等） |
| `pkg/` | 可复用工具（缓存、日志、支付、邮件、短信等） |
| `queue/` | Asynq 任务处理器，用于异步操作 |
| `scheduler/` | 定时任务定义 |
| `initialize/` | 系统初始化（配置、迁移、默认数据） |
| `adapter/` | 外部服务适配器，用于节点通信 |

### API 组织

API 分为三个主要区域（参见 `apis/`）：

- **admin/**: 管理面板 API（用户管理、订单、工单、系统配置）
- **public/**: 终端用户 API（订阅、订单、公告）
- **auth/**: 认证 API（登录、注册、OAuth）

每个 API 文件定义路由、请求/响应类型。`ppanel.api` 文件导入所有子 API。

### ServiceContext 模式

`internal/svc/serviceContext.go` 定义注入到所有处理器的中心上下文：

```go
type ServiceContext struct {
    DB           *gorm.DB
    Redis        *redis.Client
    Config       config.Config
    Queue        *asynq.Client
    // 领域模型
    UserModel    user.Model
    OrderModel   order.Model
    // ... 其他模型
}
```

添加需要数据库访问的新功能时，请将模型添加到 ServiceContext。

### 添加新 API 端点

1. 在 `apis/` 中相应的 `.api` 文件中添加路由定义
2. 运行 `./script/generate.sh` 生成处理器和类型
3. 在 `internal/logic/<area>/<feature>/` 中实现业务逻辑
4. 将任何新模型添加到 `internal/model/` 和 ServiceContext

### 配置

默认配置位置：`etc/ppanel.yaml`

必需设置：
- MySQL 连接（Addr、Username、Password、Dbname）
- Redis 连接（Host、Pass、DB）
- JWT 密钥（AccessSecret）

环境变量可覆盖配置：
- `PPANEL_DB`: MySQL DSN（例如：`user:pass@tcp(localhost:3306)/dbname`）
- `PPANEL_REDIS`: Redis URI（例如：`redis://localhost:6379`）

### 数据库

使用 GORM 配合 MySQL。模型位于 `internal/model/`。项目通过 `gorm.io/plugin/soft_delete` 使用软删除。

迁移在 `initialize/migrate/` 中处理。

### 主要依赖

- **Gin**: HTTP 框架
- **GORM**: ORM
- **Asynq**: 任务队列（基于 Redis）
- **go-resty**: 外部 API 的 HTTP 客户端
- **Stripe/Alipay**: 支付集成
- **Telegram Bot API**: 通知