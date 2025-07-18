# DevOps Asset Management System

企业级DevOps资产管理系统 - 一个现代化的前后端分离管理平台

## 🎯 项目简介

DevOps Asset Management System 是一个基于Vue 3 + Go的企业级资产管理系统，提供完整的用户管理、角色权限、菜单管理等核心功能。系统采用现代化的技术栈，具备高性能、易扩展、易维护的特点。

### ✨ 主要特点

- 🏗️ **现代化架构**: 前后端分离，微服务架构
- 🔐 **安全可靠**: JWT认证 + RBAC权限控制
- 🚀 **高性能**: Redis缓存 + 数据库连接池优化
- 🐳 **容器化**: Docker支持，开箱即用
- 📊 **实时监控**: 健康检查 + 日志管理
- 🎨 **现代UI**: 基于Element Plus的现代化界面
- 📝 **规范化**: 统一的错误处理和日志管理

## 🛠️ 技术栈

### 后端技术栈
- **Go 1.21+** - 主要开发语言
- **Gin** - Web框架
- **GORM** - ORM框架
- **MySQL 8.0** - 主数据库
- **Redis 7.0** - 缓存数据库
- **JWT** - 身份认证
- **Swagger** - API文档

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

### 🔧 系统管理
- **健康检查**: 系统状态监控
- **日志管理**: 系统操作日志
- **配置管理**: 系统参数配置
- **API文档**: Swagger API文档

## 🚀 快速开始

### 使用Docker Compose（推荐）

1. **克隆项目**
```bash
git clone <your-repo-url>
cd devops-asset
```

2. **启动所有服务**
```bash
# 生产环境部署
docker-compose --profile production up -d

# 开发环境部署
docker-compose --profile development up -d
```

3. **初始化数据库**
```bash
# 执行数据库迁移
docker exec -it devops-asset-server ./server --migrate
```

4. **访问系统**
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
go run cmd/server/main.go --migrate

# 启动开发服务
go run cmd/server/main.go --dev
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
- 至少 2GB 内存
- 至少 10GB 磁盘空间

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

### API设计
- RESTful API设计风格
- 统一的响应格式
- 完整的错误处理
- Swagger文档自动生成

### 数据库设计
- 规范化的数据库设计
- 外键约束保证数据一致性
- 软删除支持
- 自动时间戳管理

## 🔒 安全配置

### 生产环境安全检查清单

- [ ] 修改默认数据库密码
- [ ] 配置强密码策略
- [ ] 启用HTTPS
- [ ] 配置防火墙规则
- [ ] 定期备份数据
- [ ] 监控异常登录
- [ ] 更新安全补丁

### 权限控制
- 基于角色的访问控制(RBAC)
- 菜单级权限控制
- API接口权限验证
- 数据行级权限控制

## 📞 技术支持

### 常见问题
- [FAQ文档](docs/FAQ.md)
- [故障排除指南](docs/TROUBLESHOOTING.md)
- [性能优化指南](docs/PERFORMANCE.md)

### 获取帮助
- 提交Issue：[GitHub Issues](issues/)
- 查看文档：[项目文档](docs/)
- 社区讨论：[Discussions](discussions/)

## 🤝 贡献指南

我们欢迎社区贡献！请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解如何参与项目开发。

### 贡献类型
- 🐛 Bug修复
- ✨ 新功能开发
- 📚 文档改进
- 🎨 UI/UX优化
- ⚡ 性能优化

## 📄 许可证

本项目基于 [MIT License](LICENSE) 开源协议。

## 🙏 致谢

感谢以下开源项目的支持：
- [Vue.js](https://vuejs.org/) - 渐进式JavaScript框架
- [Gin](https://gin-gonic.com/) - Go Web框架
- [Element Plus](https://element-plus.org/) - Vue 3组件库
- [GORM](https://gorm.io/) - Go ORM库

---

⭐ 如果这个项目对您有帮助，请给我们一个星标！ 