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

### 🖥️ 主机管理
- **主机列表**: 支持筛选、排序和分页
- **主机详情**: 查看主机详细信息
- **主机操作**: 添加、编辑、删除主机
- **主机分组**: 主机组织和管理
- **远程操作**: SSH终端和SFTP文件传输
- **批量操作**: 批量管理主机

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

### UI组件库说明
本项目前端使用 **Element Plus** 作为UI组件库，设计风格参考 [art-design-pro](https://github.com/Daymychen/art-design-pro)。所有自定义组件和页面都应当遵循Element Plus的设计规范，保持界面风格的一致性。

## 🛠️ 模块修复记录

### 🖥️ 主机模块修复记录

#### 2023-07-18 修复

1. **UI组件库迁移**
   - 将主机模块从Ant Design Vue完全迁移到Element Plus
   - 优化了主机列表页面的样式和交互
   - 修复了主机模态框组件的兼容性问题

2. **样式优化**
   - 统一使用Element Plus的设计风格
   - 优化了表格、表单、按钮等组件的样式
   - 添加了加载状态、错误状态和空数据状态的处理

3. **功能完善**
   - 完善了主机状态显示
   - 优化了主机操作菜单
   - 添加了主机组显示
   - 改进了批量操作功能

4. **代码优化**
   - 重构了组件结构，提高代码可维护性
   - 优化了API调用方式
   - 统一了状态管理方式

### ☁️ 云账号管理模块修复记录

#### 2025-07-23 修复

**问题描述**：云账号管理页面无法正常加载，前端动态路由解析失败

**根本原因**：
- 数据库中云账号管理菜单的 `component` 字段包含前导斜杠 `/cmdb/provider/index`
- 前端动态路由处理期望的路径格式为 `cmdb/provider/index`（不带前导斜杠）
- 路由解析失败导致组件无法正确加载

**修复内容**：

1. **后端路由修复**
   - 文件：`server/internal/db/mysql/system/host_menu.go`
   - 函数：`addProviderManagementMenuIfNotExists`
   - 修复：添加自动检测和修复逻辑，当菜单已存在时检查 `component` 字段格式
   - 结果：确保数据库中的路径格式为 `cmdb/provider/index`

2. **前端组件优化**
   - 文件：`web/src/views/cmdb/provider/index.vue`
   - 优化：现代化的云账号管理界面设计
   - 功能：支持阿里云、腾讯云、AWS等多云厂商管理
   - 特性：卡片式布局、实时状态显示、批量操作

3. **前端模态框重构**
   - 文件：`web/src/views/cmdb/provider/ProviderModal.vue`
   - 设计：三步骤向导式添加流程
   - 功能：云厂商选择 → 认证配置 → 完成设置
   - 验证：实时凭证验证和区域获取

4. **API接口完善**
   - 路径：`/api/v1/cmdb/providers/*`
   - 功能：CRUD操作、凭证验证、资源同步
   - 状态：所有接口正常注册并可用

**技术细节**：

```go
// 修复逻辑示例
if menu.Component != correctComponent {
    zap.L().Info("fixing provider management menu component",
        zap.String("old", menu.Component),
        zap.String("new", correctComponent))

    err = db.Model(&menu).Update("component", correctComponent).Error
    if err != nil {
        return err
    }
}
```

**验证结果**：
- ✅ 后端服务正常启动
- ✅ 所有API路由正确注册
- ✅ 数据库菜单路径自动修复
- ✅ 前端动态路由正确解析
- ✅ 云账号管理页面正常加载

## 📋 后续优化计划

### 🔄 短期优化（1-2周）

1. **云账号管理功能完善**
   - [ ] 添加更多云厂商支持（华为云、百度云等）
   - [ ] 实现云账号配额监控和告警
   - [ ] 优化凭证验证的错误提示
   - [ ] 添加云账号使用统计和成本分析

2. **前端体验优化**
   - [ ] 添加页面加载骨架屏
   - [ ] 优化移动端响应式布局
   - [ ] 实现暗色主题支持
   - [ ] 添加操作引导和帮助文档

3. **性能优化**
   - [ ] 实现前端组件懒加载
   - [ ] 优化API接口响应时间
   - [ ] 添加Redis缓存策略
   - [ ] 实现数据库查询优化

### 🚀 中期规划（1-2月）

1. **监控告警系统**
   - [ ] 实现资源监控大盘
   - [ ] 添加自定义告警规则
   - [ ] 集成邮件/短信通知
   - [ ] 实现告警历史和统计

2. **自动化运维**
   - [ ] 实现资源自动发现
   - [ ] 添加定时任务管理
   - [ ] 实现批量运维脚本执行
   - [ ] 添加变更审批流程

3. **安全增强**
   - [ ] 实现操作审计日志
   - [ ] 添加敏感数据加密
   - [ ] 实现IP白名单控制
   - [ ] 添加双因子认证

### 🎯 长期目标（3-6月）

1. **微服务架构**
   - [ ] 服务拆分和治理
   - [ ] 实现服务注册发现
   - [ ] 添加分布式链路追踪
   - [ ] 实现配置中心

2. **多租户支持**
   - [ ] 实现租户隔离
   - [ ] 添加资源配额管理
   - [ ] 实现计费和结算
   - [ ] 添加租户自定义配置

## ⚠️ 已知问题和限制

### 🐛 当前已知问题

1. **前端问题**
   - [ ] 部分页面在IE浏览器下兼容性问题
   - [ ] 大数据量表格渲染性能待优化
   - [ ] 文件上传组件在某些情况下进度显示异常

2. **后端问题**
   - [ ] Redis连接失败时的降级策略需要完善
   - [ ] 并发场景下的数据一致性问题
   - [ ] 长时间运行后的内存泄漏风险

3. **部署问题**
   - [ ] Docker容器重启后某些配置丢失
   - [ ] 数据库迁移在某些版本下可能失败
   - [ ] SSL证书自动续期功能待实现

### 🔧 技术债务

1. **代码质量**
   - [ ] 部分历史代码需要重构
   - [ ] 单元测试覆盖率需要提升
   - [ ] API文档需要完善和更新

2. **架构优化**
   - [ ] 数据库连接池配置需要调优
   - [ ] 缓存策略需要统一规范
   - [ ] 错误处理机制需要标准化

3. **运维工具**
   - [ ] 缺少自动化部署脚本
   - [ ] 监控指标需要更加完善
   - [ ] 日志收集和分析工具待集成

### 💡 改进建议

1. **开发流程**
   - 建议引入代码审查机制
   - 完善CI/CD流水线
   - 建立更完善的测试环境

2. **文档维护**
   - 定期更新API文档
   - 完善开发者指南
   - 建立问题排查手册

3. **社区建设**
   - 建立问题反馈渠道
   - 定期发布版本更新
   - 收集用户使用反馈

## 📞 技术支持

如果在使用过程中遇到问题，请按以下方式获取帮助：

1. **查看文档**：首先查阅本README和相关文档
2. **检查日志**：查看应用和系统日志获取错误信息
3. **问题反馈**：通过Issue提交问题报告
4. **社区讨论**：参与项目讨论和经验分享

### 注意事项

1. 本项目UI使用Element Plus，请不要混用其他UI库
2. 参考art-design-pro的设计风格进行开发
3. 保持代码风格一致性，遵循项目既定的开发规范
4. 提交代码前请确保通过所有测试用例
