English | [简体中文](./README.zh-CN.md)

## About Art Design Pro

As a developer, I needed to build admin management systems for multiple projects but found that traditional systems couldn't fully meet the requirements for user experience and visual design. Therefore, I created Art Design Pro, an open-source admin management solution focused on user experience and rapid development. Based on the ElementPlus design specifications, it has been visually optimized to provide a more beautiful and practical front-end interface, helping you easily build high-quality admin systems.

## Official Website

[Visit the official documentation](https://www.lingchen.kim/art-design-pro/docs/en/)

## Demo Images

### Light Theme

![Light Theme](https://www.qiniu.lingchen.kim/art_design_pro_readme_cover1.png)

![Light Theme](https://www.qiniu.lingchen.kim/art_design_pro_readme_cover2.png)

### Dark Theme

![Dark Theme](https://www.qiniu.lingchen.kim/art_design_pro_readme_cover3.png)

![Dark Theme](https://www.qiniu.lingchen.kim/art_design_pro_readme_cover4.png)

## Features

- Uses the latest technology stack
- Built-in common business component templates
- Provides multiple theme modes and customizable themes
- Beautiful UI design, excellent user experience, and attention to detail
- System fully supports customization, meeting your personalized needs

## Functionality

- Rich theme switching
- Global search
- Lock screen
- Multi-tabs
- Global breadcrumbs
- Multi-language support
- Icon library
- Rich text editor
- Echarts charts
- Utils toolkit
- Network exception handling
- Route-level authentication
- Sidebar menu authentication
- Authentication directives
- Mobile adaptation
- Excellent persistent storage solution
- Local data storage validation
- Code commit validation and formatting
- Code commit standardization

## Compatibility

- Supports modern mainstream browsers such as Chrome, Safari, Firefox, etc.

## Installation and Running

```bash
# Install dependencies
pnpm install

# If pnpm install fails, try using the following command to install dependencies
pnpm install --ignore-scripts

# Start local development environment
pnpm dev

# Build for production
pnpm build
```

## Technical Support

QQ Group: <a href="https://qm.qq.com/cgi-bin/qm/qr?k=Gg6yzZLFaNgmRhK0T5Qcjf7-XcAFWWXm&jump_from=webapi&authKey=YpRKVJQyFKYbGTiKw0GJ/YQXnNF+GdXNZC5beQQqnGZTvuLlXoMO7nw5fNXvmVhA">821834289</a> (Click the link to join the group chat)

## Donation

If my project has been helpful to you, donations are welcome! Your support will be used to purchase tools like ChatGPT, Cursor, etc., to improve development efficiency and make the project even better. Thank you for your encouragement and support!

![Donation QR Code](https://www.qiniu.lingchen.kim/%E7%BB%84%202%402x%202.png)

## 2025-07-18 Host Module Frontend Refactor

本次提交主要针对 **CMDB / 主机管理模块** 进行全面优化，关键点如下：

1. **UI 全面迁移到 Element Plus**
   - 将遗留的 Ant Design Vue 组件（Modal / Form / Table / Steps 等）替换为 Element Plus 同类组件。
   - 统一使用 Element Plus 图标库，解决编译期 `export not found` 报错。

2. **功能对齐后端接口**
   - 列表查询：分页 / 关键字 / 状态 / 主机组 / 区域全部透传给后端 `GET /api/v1/cmdb/hosts`。
   - 批量导入：接入 `POST /api/v1/cmdb/hosts/batch_import`，支持 .xlsx / .csv 拖拽上传与进度提示。
   - 云同步：接入 `POST /api/v1/cmdb/hosts/sync`，全局 Loading 状态提示。
   - 手动添加主机：使用 `POST /api/v1/cmdb/hosts/manual`，支持主机组级联选择。
   - 主机状态同步、批量删除等操作均调用对应后端接口。

3. **状态与样式优化**
   - 新增 `src/assets/styles/change.scss` 覆盖样式：去除页面水印、修正空数据图片尺寸、统一卡片 padding。
   - 使用 Element Plus `el-empty` / `el-result` 显示空状态、导入结果。

4. **Store 增强**
   - `hostStore` 新增 `syncHosts、getHostList` 兼容方法以及监控相关 fetch 方法，旧组件无需改动即可使用。

5. **移除 ant-design-vue 依赖**
   - 所有模块已无需 Ant Design Vue，后续可删除 `ant-design-vue` 依赖并执行 `pnpm install` 清理。

> 经过以上调整，主机管理页面在 UI 与接口层面与后端保持一致，编译与运行均无报错，体验更统一、性能更优。

## 2025-07-19 Provider Module Refactor

- 厂商管理页面（/cmdb/provider）全面迁移至 Element Plus：el-card / el-table / el-dialog / el-form 等。
- 新增功能：
  1. 搜索 + 前端分页；
  2. 新建 / 编辑厂商（ProviderModal.vue 重写）；
  3. 同步云资源、删除厂商，操作流程全局 Loading 与提示；
- store 使用 hostApi provider 接口（getProviderList / createProvider / updateProvider / deleteProvider / syncResources）。
- 完全移除项目中 ant-design-vue 依赖，统一 UI/交互风格。
