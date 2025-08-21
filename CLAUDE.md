# CLAUDE.md

此文件为 Claude Code（claude.ai/code）在处理本仓库代码时提供指导。

- 技术问题优先用中文回复
- 代码注释要详细，新人能看懂
- 优先使用成熟开源方案，避免重复造轮子
- 性能优化必须有benchmark数据
- 数据库操作必须考虑并发安全

## 项目概述
这是一个 **DevOps 资产管理系统** —— 一个基于 Vue 3 + Go 架构的现代化云原生企业级资产管理平台。

**主要功能：**
* 云账号管理（多云支持：阿里云、腾讯云、AWS）
* 主机管理，包含监控与告警
* DNS 与证书管理（集成 Let's Encrypt）
* CMDB（配置管理数据库）
* 用户与角色管理（基于 RBAC）
* 实时监控与告警
* SSH/SFTP 终端访问

## 技术栈
### 后端 (Go)

* **框架**: Gin + GORM + MySQL + Redis
* **语言**: Go 1.24.1
* **认证**: JWT
* **数据库**: MySQL 8.0, Redis 7.0
* **文档**: Swagger/OpenAPI
* **日志**: Zap 结构化日志
* **架构**: 领域驱动设计 + 分层架构

### 前端 (Vue)

* **框架**: Vue 3 + TypeScript + Vite
* **UI 库**: Element Plus（主要），少量 Ant Design Vue 组件
* **状态管理**: Pinia
* **路由**: Vue Router（支持动态路由）
* **构建工具**: Vite
* **模板基础**: Art-Design-Pro

## 开发命令

### 后端开发

```bash
cd server

# 安装依赖
go mod download

# 数据库迁移
go run cmd/main.go --dev --migrate

# 启动开发服务
go run cmd/main.go --dev

# 构建生产环境
go build -o main cmd/main.go
```

### 前端开发

```bash
cd web

# 安装依赖
pnpm install

# 启动开发服务
pnpm dev

# 构建生产环境
pnpm build

# 代码检查与修复
pnpm lint
pnpm fix

# 运行测试
pnpm test
```

### Docker 部署

```bash
# 生产环境部署
docker-compose --profile production up -d

# 开发环境部署
docker-compose --profile development up -d

# 使用部署脚本
./deploy.sh production          # 生产
./deploy.sh development         # 开发
./deploy.sh status              # 查看状态
./deploy.sh health              # 健康检查
```

## 项目结构

### 后端结构

```
server/
├── cmd/main.go                 # 应用入口
├── internal/
│   ├── api/handler/            # HTTP 控制器（Gin handler）
│   │   ├── cmdb/              # CMDB 模块
│   │   ├── dns/               # DNS 模块
│   │   └── system/            # 系统模块
│   ├── service/               # 业务逻辑层
│   ├── repository/            # 数据访问层 (GORM)
│   ├── model/                 # 数据库模型
│   ├── router/                # 路由定义
│   ├── middleware/            # HTTP 中间件
│   ├── response/              # 统一响应处理
│   ├── config/                # 配置管理
│   └── db/                    # 数据库客户端
├── configs/                   # 配置文件
├── pkg/                       # 公共包
└── docs/                      # API 文档
```

### 前端结构

```
web/
├── src/
│   ├── api/                   # API 客户端
│   ├── components/            # 可复用组件
│   ├── views/                 # 页面组件
│   │   ├── cmdb/             # CMDB 页面
│   │   ├── dns/              # DNS 页面
│   │   └── system/           # 系统页面
│   ├── router/               # 路由配置
│   ├── store/                # Pinia 状态管理
│   ├── utils/                # 工具函数
│   └── assets/               # 静态资源
├── public/                   # 公共资源
└── docs/                     # 前端文档
```

## 架构原则

### 后端设计模式

* **分层架构**: Handler → Service → Repository → Model
* **统一响应**: 所有接口返回 `200`，并使用标准 JSON 格式
* **错误处理**: Service 层返回 `ServiceError`，由中间件统一处理
* **认证**: JWT 中间件，开发模式支持跳过认证
* **多租户**: 数据库层面隔离

### 响应格式

```go
{
    "code": 200,            // 业务状态码
    "status": "success",    // 状态
    "message": "OK",        // 可选提示
    "data": {},             // 返回数据
    "count": 0,             // 分页可用
    "timestamp": 1640995200 // Unix 时间戳
}
```

### Service 层示例

```go
// 业务错误使用 ServiceError
func (s *SomeService) DoSomething() error {
    if err := validation; err != nil {
        return response.NewServiceError(response.INVALID_PARAMS, "参数校验失败")
    }
    return nil
}

// Handler 转换 ServiceError → HTTP 响应
func (h *SomeHandler) HandleRequest(c *gin.Context) {
    if err := h.service.DoSomething(); err != nil {
        middleware.HandleServiceError(c, err)
        return
    }
    response.ReturnSuccess(c)
}
```

### 前端约定

* **API 客户端**: 统一在 `src/api/client.ts`
* **错误处理**: 拦截器统一处理 API 错误
* **UI 库**: 优先使用 Element Plus (`el-*`)，少量 AntD Vue
* **路由守卫**: 认证与权限校验在 `router/guards/`
* **路由注册**: 动态路由在后端配置
* **状态管理**: Pinia 按模块组织

## 关键模块

### 1. DNS 与证书管理

* **位置**: `server/internal/*/dns/`, `web/src/views/dns/`
* **功能**: 域名管理、DNS 记录、Let’s Encrypt 证书、HTTPS 监控
* **支持厂商**: 阿里云、AWS Route53、腾讯云、GoDaddy
* **关键文件**:

  * 后端: `internal/service/dns/`, `internal/provider/dns/`
  * 前端: `views/dns/domains/`, `views/dns/certs/`

### 2. CMDB 主机管理

* **位置**: `server/internal/*/cmdb/`, `web/src/views/cmdb/`
* **功能**: 主机清单、监控、SSH 终端、文件传输
* **关键文件**:

  * 后端: `internal/service/cmdb/`, `internal/adapter/`
  * 前端: `views/cmdb/host/`, `components/ssh/`

### 3. 用户与系统管理

* **位置**: `server/internal/*/system/`, `web/src/views/system/`
* **功能**: 用户、角色、权限、菜单
* **认证**: 基于 JWT + RBAC
* **关键文件**:

  * 后端: `internal/service/user_service.go`
  * 前端: `views/system/user/`, `store/modules/user.ts`

## 开发规范

### 代码规范

* **Go**: 遵循官方规范，使用 `gofmt` + `golangci-lint`
* **Vue**: 使用 Composition API，开启 TS 严格模式
* **命名**: 清晰、易懂，避免缩写
* **注释**: 必要时描述业务逻辑与复杂算法

### API 开发流程

1. 在 `internal/model/` 定义模型
2. 在 `internal/repository/` 创建仓储层
3. 在 `internal/service/` 实现业务逻辑
4. 在 `internal/api/handler/` 编写控制器
5. 在 `internal/router/` 注册路由
6. 在前端 `src/api/` 添加 API 客户端

### 数据库规范

* **迁移**: `server/internal/database/migrations/`
* **模型**: 遵循 GORM 规范，启用软删除 `deleted_at`
* **关系**: 明确外键
* **索引**: 高频查询字段加索引

### 安全最佳实践

* 不记录密钥、token、密码
* 敏感数据使用 AES 加密
* 校验所有输入参数
* 使用预处理语句（GORM 已支持）
* 严格执行 RBAC

## 测试

### 后端测试

```bash
cd server
go test ./...                  # 运行所有测试
go test -cover ./...           # 带覆盖率
go test -v ./internal/service/ # 指定包
```

### 前端测试

```bash
cd web
pnpm test                      # 运行测试
pnpm test:coverage             # 带覆盖率
```

## 配置

### 环境变量

* `APP_MODE`: `dev` / `prod`
* `MYSQL_*`: MySQL 配置
* `REDIS_*`: Redis 配置
* `JWT_SECRET`: JWT 签名秘钥
* `AES_KEY`: AES 加密秘钥（32 位）

### 配置文件

* 后端: `server/configs/config.yaml`
* 前端: `web/.env.development`, `web/.env.production`

## 常见开发任务

### 新增功能模块

1. 建立数据库模型与迁移
2. 实现仓储层 CRUD
3. 编写 Service 层逻辑
4. 创建 HTTP handler + Swagger 文档
5. 注册路由 + 中间件
6. 前端页面 + 组件
7. 添加 API 客户端方法
8. 更新导航菜单

### 调试问题

* 后端日志: `server/log/`
* 前端: 浏览器开发者工具
* API 调试: Swagger UI `/swagger/index.html`
* 健康检查: GET `/health`

## 注意事项

* **UI 框架**: 项目以 Element Plus 为主，避免混用其他框架
* **认证**: 开发模式（`--dev`）可跳过 JWT 认证
* **多租户**: 所有数据查询需包含租户隔离
* **错误码**: 统一在 `response/code.go` 定义
* **文档**: 如有架构调整，需更新 CLAUDE.md

## 常用脚本与命令

### 数据库操作

```bash
# 执行迁移
go run cmd/main.go --migrate

# 重置数据库（仅开发环境）
go run cmd/main.go --migrate --reset
```

### 开发流程

```bash
# 启动后端（开发模式）
cd server && go run cmd/main.go --dev

# 启动前端
cd web && pnpm dev

# 全栈启动（Docker）
./deploy.sh development
```

### 生产部署

```bash
# 构建并部署
./deploy.sh production

# 健康检查
./deploy.sh health

# 查看日志
docker-compose logs -f
```

---

此项目遵循现代 DevOps 实践，支持容器化部署、全面监控，并以安全优先的设计理念为核心。
