# 证书管理功能测试指南

## 修复的问题

### 1. 删除证书显示"undefined"问题
**问题原因**: 缺少单个删除证书的API接口
**修复方案**: 
- 添加了 `DELETE /api/v1/dns/certificates/:id` 接口
- 在handler中添加了 `DeleteCertificate` 方法
- 在service中添加了 `DeleteCertificate` 方法
- 改善了前端错误处理，显示详细错误信息

### 2. 证书下载格式问题
**问题原因**: 下载弹窗UI不够友好
**修复方案**:
- 简化了下载格式选择弹窗
- 参考阿里云证书下载页面，提供清晰的格式说明
- 支持的格式：
  - PEM格式 (通用) - 标准格式，适用于大多数服务器
  - Nginx配置 - 包含Nginx配置示例
  - Apache配置 - 包含Apache配置示例
  - IIS配置 - Windows IIS服务器配置
  - Tomcat配置 - 包含JKS格式证书
  - 私钥文件 - 仅下载私钥(.key)
  - 证书文件 - 仅下载证书(.crt)
  - 证书链 - 完整证书链文件

### 3. 通配符域名申请无输出问题
**修复方案**:
- 增加了详细的日志记录
- 改善了错误处理和错误信息显示
- 在申请过程中添加了更多调试信息

### 4. 申请证书页面跳转错误
**问题原因**: 路由配置错误，申请证书路由指向了上传证书页面
**修复方案**:
- 修改路由配置：
  - `/dns/certs/create` -> 指向申请证书页面 (`create-new.vue`)
  - `/dns/certs/upload` -> 指向上传证书页面 (`create.vue`)
- 更新页面跳转逻辑

## 证书状态说明

证书在数据库中的状态包括：
- `pending` - 待处理（申请中）
- `issued` - 已签发（正常可用）
- `expired` - 已过期
- `revoked` - 已吊销（不可用）

## 测试步骤

### 1. 测试删除功能
1. 进入证书管理页面
2. 选择一个证书，点击删除
3. 确认删除弹窗显示正确的证书名称
4. 点击确定，应该显示"删除成功"而不是"undefined"

### 2. 测试下载功能
1. 选择一个已签发的证书
2. 点击下载按钮
3. 应该弹出格式选择对话框，包含8种格式选项
4. 选择任意格式，应该能正常下载

### 3. 测试申请证书功能
1. 点击"申请证书"按钮
2. 应该跳转到申请证书页面，而不是上传证书页面
3. 尝试申请通配符域名证书（如 *.example.com）
4. 应该能看到详细的申请进度和错误信息

### 4. 测试上传证书功能
1. 在下拉菜单中选择"上传证书"
2. 应该跳转到上传证书页面
3. 可以上传证书文件或粘贴证书内容

### 5. 测试吊销功能
1. 选择一个已签发的证书
2. 点击吊销按钮
3. 确认吊销后，证书状态应该变为"已吊销"
4. 在数据库中状态字段应该为"revoked"

## API接口测试

可以使用以下curl命令测试API接口：

```bash
# 删除证书
curl -X DELETE "http://localhost:8080/api/v1/dns/certificates/1" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 下载证书
curl -X GET "http://localhost:8080/api/v1/dns/certificates/1/download?format=pem" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  --output certificate.zip

# 吊销证书
curl -X POST "http://localhost:8080/api/v1/dns/certificates/1/revoke" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 注意事项

1. 确保后端服务正在运行
2. 确保数据库连接正常
3. 确保有有效的JWT token
4. 测试时使用测试环境，避免影响生产数据
