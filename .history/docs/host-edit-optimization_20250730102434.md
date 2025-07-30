# 主机管理编辑功能优化

## 问题描述

在主机编辑功能中存在以下问题：
1. 编辑主机时会显示明文密码，存在安全风险
2. 即使用户没有修改密码，系统也会重新提交密码，导致不必要的加密操作
3. 用户体验不友好，无法区分是否需要修改认证信息

## 优化方案

### 1. 前端优化 (web/src/views/cmdb/host/index.vue)

#### 新增响应式变量
```javascript
const isPasswordModified = ref(false)
const isPrivateKeyModified = ref(false)
```

#### 编辑时数据绑定优化
- 不再直接显示明文密码和私钥
- 添加输入监听，检测用户是否修改了认证信息
- 提供友好的提示信息

#### 保存逻辑优化
- 只有在新增模式或用户真的修改了密码时才提交password字段
- 只有在新增模式或用户真的修改了私钥时才提交private_key字段

```javascript
// 处理密码字段：只有在新增模式或用户修改了密码时才提交
if (!isEditMode.value || isPasswordModified.value) {
  saveData.password = hostFormData.ssh_config.password || ''
}

// 处理私钥字段：只有在新增模式或用户修改了私钥时才提交
if (!isEditMode.value || isPrivateKeyModified.value) {
  saveData.private_key = hostFormData.ssh_config.private_key || ''
}
```

### 2. 编辑页面优化 (web/src/views/cmdb/host/edit.vue)

#### 类似的优化措施
- 添加密码修改标记
- 编辑时不显示明文密码
- 保存时只有修改了密码才提交密码字段

### 3. 后端逻辑 (server/internal/service/cmdb/host_service.go)

后端已经实现了正确的逻辑：
- 如果密码为空，保留原密码
- 如果私钥为空，保留原私钥
- 自动处理加密和解密

## 用户体验改进

### 1. 界面提示
- 在编辑模式下，密码和私钥输入框显示提示："如需修改密码请重新输入"
- 添加信息提示："当前已设置密码，如不修改请保持为空"

### 2. 安全性提升
- 编辑时不再显示明文密码
- 减少不必要的密码重新加密操作
- 保持原有的"不填则不修改"逻辑

### 3. 样式优化
添加了`.password-hint`样式类，提供美观的提示信息展示。

## 使用说明

### 添加主机
- 正常填写所有字段，包括密码和私钥

### 编辑主机
- 如果不需要修改密码，保持密码字段为空
- 如果需要修改密码，重新输入新密码
- 系统会自动检测并只在必要时更新认证信息

## 技术细节

### 修改标记机制
- `isPasswordModified`: 检测密码是否被用户修改
- `isPrivateKeyModified`: 检测私钥是否被用户修改
- 通过`@input`事件监听用户输入

### 条件提交逻辑
- 新增模式：提交所有字段
- 编辑模式：只提交用户修改的认证字段

这种优化既保证了安全性，又提供了良好的用户体验，符合"不填则不修改"的业务逻辑。 
