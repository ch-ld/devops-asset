# DevOps Asset Management System

企业级DevOps资产管理系统 - 一个现代化的云原生资产管理平台

## 🎯 项目简介

DevOps Asset Management System 是一个基于Vue 3 + Go的企业级云资产管理系统，提供完整的云账号管理、主机管理、监控告警等核心功能。系统采用现代化的技术栈，具备高性能、易扩展、易维护的特点。

### ✨ 主要特点

- 🏗️ **现代化架构**: 前后端分离，容器化部署
- 🔐 **安全可靠**: JWT认证 + RBAC权限控制 + 数据加密
- ☁️ **多云支持**: 支持阿里云、腾讯云、AWS等主流云服务商
- 🚀 **高性能**: Redis缓存 + 数据库连接池优化
- 🐳 **容器化**: Docker支持，开箱即用
- 📊 **实时监控**: 健康检查 + 日志管理 + 指标监控
- 🎨 **现代UI**: 基于Element Plus的现代化界面
- 📝 **规范化**: 统一的错误处理和日志管理

## 🛠️ 技术栈

### 后端技术栈
- **Go 1.24+** - 主要开发语言
- **Gin** - Web框架
- **GORM** - ORM框架
- **MySQL 8.0** - 主数据库
- **Redis 7.0** - 缓存数据库
- **JWT** - 身份认证
- **Swagger** - API文档
- **Zap** - 结构化日志

### 前端技术栈
- **Vue 3** - 前端框架
- **TypeScript** - 类型支持
- **Vite** - 构建工具
- **Element Plus** - UI组件库
- **Pinia** - 状态管理
- **Vue Router** - 路由管理

### 部署技术栈
- **Docker & Docker Compose** - 容器化部署
- **Nginx** - 反向代理和静态文件服务
- **MySQL** - 生产数据库
- **Redis** - 生产缓存

## 📋 功能模块

### 🔐 认证与授权
- **用户登录/登出**: 支持图形验证码
- **JWT认证**: 无状态身份验证
- **会话管理**: 登录状态持久化
- **密码管理**: 安全的密码修改机制

### 👥 用户管理
- **用户信息管理**: 创建、编辑、删除用户
- **用户状态控制**: 启用/禁用用户账号
- **用户详情查看**: 完整的用户信息展示
- **登录日志**: 用户登录历史记录

### 🎭 角色管理
- **角色定义**: 创建和管理系统角色
- **权限分配**: 为角色分配菜单权限
- **角色状态**: 启用/禁用角色
- **角色关联**: 用户与角色的关联管理

### 📂 菜单管理
- **菜单结构**: 树形菜单结构管理
- **权限控制**: 细粒度的菜单权限控制
- **动态菜单**: 基于权限的动态菜单生成
- **菜单配置**: 图标、路径、组件配置

### 🏢 部门管理
- **组织架构**: 树形部门结构
- **部门信息**: 部门基本信息管理
- **人员分配**: 用户与部门的关联

### ☁️ 云账号管理
- **多云支持**: 支持阿里云、腾讯云、AWS等
- **凭证管理**: 加密存储访问密钥
- **凭证验证**: 实时验证云账号凭证有效性
- **区域管理**: 支持多区域资源管理
- **资源同步**: 自动同步云资源到本地

### 🖥️ 主机管理
- **主机列表**: 支持筛选、排序和分页
- **主机详情**: 查看主机详细信息
- **主机操作**: 添加、编辑、删除主机
- **主机分组**: 主机组织和管理
- **远程操作**: SSH终端和SFTP文件传输
- **批量操作**: 批量管理主机
- **状态监控**: 实时监控主机状态
- **到期提醒**: 主机到期时间提醒

### 📊 监控告警
- **指标收集**: CPU、内存、磁盘、网络监控
- **阈值告警**: 可自定义告警阈值
- **告警通知**: 支持多种通知方式
- **历史数据**: 监控数据历史查询

### 🔧 系统管理
- **健康检查**: 系统状态监控
- **日志管理**: 系统操作日志
- **配置管理**: 系统参数配置
- **API文档**: Swagger API文档

## 🚀 快速开始

### 一键部署（推荐）

1. **克隆项目**
```bash
git clone <your-repo-url>
cd devops-asset
```

2. **配置环境变量**
```bash
# 复制环境配置模板
cp .env.example .env

# 编辑配置文件
vim .env
```

3. **一键启动**
```bash
# 生产环境部署
docker-compose --profile production up -d

# 开发环境部署
docker-compose --profile development up -d
```

4. **初始化数据库**
```bash
# 执行数据库迁移
docker exec -it devops-asset-server ./server --migrate
```

5. **访问系统**
- 前端界面：http://localhost (生产) 或 http://localhost:3000 (开发)
- API文档：http://localhost:8080/swagger/index.html
- 健康检查：http://localhost:8080/health

### 本地开发

#### 后端开发
```bash
cd server

# 安装依赖
go mod download

# 配置数据库
cp configs/config.dev.yaml configs/config.dev.local.yaml
# 编辑配置文件...

# 数据库迁移
go run cmd/main.go --migrate

# 启动开发服务
go run cmd/main.go --dev
```

#### 前端开发
```bash
cd web

# 安装依赖
pnpm install

# 配置后端地址
cp .env.example .env.development
# 编辑环境变量...

# 启动开发服务
pnpm run dev
```

## 🐳 部署指南

### 环境要求
- Docker 20.10+
- Docker Compose 2.0+
- 至少 4GB 内存
- 至少 20GB 磁盘空间

### 生产部署

1. **准备环境文件**
```bash
# 复制环境配置模板
cp .env.example .env

# 编辑生产环境配置
vim .env
```

2. **配置数据库**
```bash
# 编辑数据库配置
vim docker-compose.yml
```

3. **启动生产服务**
```bash
# 启动所有生产服务
docker-compose --profile production up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

4. **SSL证书配置（可选）**
```bash
# 将SSL证书放置到docker/ssl目录
mkdir -p docker/ssl
cp your-cert.pem docker/ssl/
cp your-key.pem docker/ssl/

# 修改nginx配置启用HTTPS
vim docker/nginx.conf
```

### 数据备份与恢复

**数据备份**
```bash
# 备份MySQL数据
docker exec devops-asset-mysql mysqldump -u devops -p devops_asset > backup.sql

# 备份Redis数据
docker exec devops-asset-redis redis-cli --rdb /data/dump.rdb
```

**数据恢复**
```bash
# 恢复MySQL数据
docker exec -i devops-asset-mysql mysql -u devops -p devops_asset < backup.sql

# 恢复Redis数据
docker cp dump.rdb devops-asset-redis:/data/
docker restart devops-asset-redis
```

## 📊 系统监控

### 健康检查
```bash
# 检查后端服务健康状态
curl http://localhost:8080/health

# 检查数据库连接
curl http://localhost:8080/api/v1/system/health
```

### 日志管理
```bash
# 查看应用日志
docker-compose logs -f server

# 查看数据库日志
docker-compose logs -f mysql

# 查看nginx日志
docker-compose logs -f nginx
```

### 性能监控
- CPU使用率监控
- 内存使用监控
- 数据库连接池监控
- Redis缓存命中率

## 🔧 开发指南

### 代码规范
- 后端遵循Go官方代码规范
- 前端遵循Vue 3 + TypeScript最佳实践
- 统一使用Prettier进行代码格式化
- 提交信息遵循Conventional Commits规范

### API文档
本项目使用Swagger生成API文档，启动服务后可访问：
- 开发环境：http://localhost:8080/swagger/index.html
- 生产环境：http://your-domain/swagger/index.html

### 数据库设计
- 使用GORM进行数据库操作
- 统一的数据库迁移管理
- 支持多环境配置

### 前端架构
- 基于Vue 3 Composition API
- TypeScript类型安全
- Element Plus组件库
- Pinia状态管理

## 🔒 安全特性

### 数据安全
- **密码加密**: 使用AES加密存储敏感数据
- **JWT认证**: 无状态身份验证
- **权限控制**: 基于角色的访问控制(RBAC)
- **数据传输**: HTTPS加密传输

### 系统安全
- **输入验证**: 严格的参数校验
- **SQL注入防护**: 使用GORM ORM防护
- **XSS防护**: 前端输出转义
- **CSRF防护**: 跨站请求伪造防护

## 📝 配置说明

### 环境变量
```bash
# 应用配置
APP_MODE=prod
APP_PORT=8080

# 数据库配置
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=devops
MYSQL_PASSWORD=your_password
MYSQL_DB=devops_asset

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password

# JWT配置
JWT_SECRET=your_jwt_secret
JWT_EXPIRATION=24h
```

### 应用配置
```yaml
app:
  name: "DevOps Asset Management System"
  version: "1.0.0"
  mode: "prod"
  port: 8080
  aes_key: "your_aes_key_32_characters_long"

database:
  type: "mysql"
  mysql:
    host: "localhost"
    port: 3306
    username: "devops"
    password: "your_password"
    database: "devops_asset"
    charset: "utf8mb4"
    parseTime: true
    loc: "Local"

redis:
  host: "localhost"
  port: 6379
  password: "your_redis_password"
  db: 0

jwt:
  secret: "your_jwt_secret"
  expiration: "24h"

log:
  level: "info"
  maxSize: 100
  maxBackups: 3
  maxAge: 7
  compress: true
```

## 🚨 故障排除

### 常见问题

1. **数据库连接失败**
```bash
# 检查数据库服务状态
docker-compose logs mysql

# 检查网络连接
docker exec -it devops-asset-server ping devops-asset-mysql
```

2. **Redis连接失败**
```bash
# 检查Redis服务状态
docker-compose logs redis

# 测试Redis连接
docker exec -it devops-asset-redis redis-cli ping
```

3. **前端无法访问后端**
```bash
# 检查nginx配置
docker-compose logs nginx

# 检查后端服务状态
curl http://localhost:8080/health
```

### 日志位置
- 应用日志：`server/log/`
- Nginx访问日志：`/var/log/nginx/access.log`
- Nginx错误日志：`/var/log/nginx/error.log`

## 🤝 贡献指南

### 开发流程
1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

### 代码提交规范
```bash
# 功能开发
git commit -m "feat: 添加云账号管理功能"

# 问题修复
git commit -m "fix: 修复主机连接超时问题"

# 文档更新
git commit -m "docs: 更新API文档"

# 性能优化
git commit -m "perf: 优化主机列表查询性能"
```

## 📞 技术支持

如果在使用过程中遇到问题，请按以下方式获取帮助：

1. **查看文档**：首先查阅本README和相关文档
2. **检查日志**：查看应用和系统日志获取错误信息
3. **问题反馈**：通过Issue提交问题报告
4. **社区讨论**：参与项目讨论和经验分享

### 联系方式
- 项目仓库：[GitHub Repository]
- 文档网站：[Documentation Site]
- 问题反馈：[Issues Page]

## 📄 许可证

本项目采用 MIT 许可证 - 详情请查看 [LICENSE](LICENSE) 文件

## 🎉 致谢

感谢所有为本项目做出贡献的开发者和用户！

---

**注意**: 本项目UI使用Element Plus，请不要混用其他UI库。开发时请遵循项目既定的开发规范，保持代码风格一致性。
