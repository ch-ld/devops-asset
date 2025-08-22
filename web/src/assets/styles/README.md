# Art Design Pro 现代化设计系统

本项目采用了基于 Art Design Pro 的现代化设计系统，提供统一的视觉语言和组件样式。

## 设计原则

### 1. 简洁明了
- 采用简洁的视觉设计，减少不必要的装饰元素
- 突出内容本身，提升用户体验

### 2. 一致性
- 统一的颜色系统、字体规范、间距标准
- 保持组件样式的一致性

### 3. 现代化
- 使用现代化的设计语言
- 支持响应式设计
- 优雅的动画和交互效果

## 设计令牌 (Design Tokens)

### 颜色系统
```scss
// 主色调
--primary-color: #1890ff;
--primary-hover: #40a9ff;
--primary-active: #096dd9;
--primary-light: #e6f7ff;

// 功能色
--success-color: #52c41a;
--warning-color: #faad14;
--error-color: #ff4d4f;

// 中性色
--text-primary: #262626;
--text-secondary: #595959;
--text-tertiary: #8c8c8c;
--text-quaternary: #bfbfbf;

// 背景色
--bg-primary: #ffffff;
--bg-secondary: #fafafa;
--bg-tertiary: #f5f5f5;
--bg-quaternary: #f0f0f0;
```

### 间距系统
```scss
--spacing-xs: 4px;
--spacing-sm: 8px;
--spacing-md: 16px;
--spacing-lg: 24px;
--spacing-xl: 32px;
--spacing-2xl: 48px;
```

### 圆角系统
```scss
--radius-sm: 4px;
--radius-md: 6px;
--radius-lg: 8px;
--radius-xl: 12px;
```

### 阴影系统
```scss
--shadow-sm: 0 1px 2px 0 rgba(0, 0, 0, 0.03), 0 1px 6px -1px rgba(0, 0, 0, 0.02);
--shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
--shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
```

## 组件样式

### 1. 页面容器 (.modern-page-container)
现代化的页面容器，提供统一的页面布局和背景。

**特性：**
- 浅灰色背景 (#fafafa)
- 统一的内边距
- 最小高度为视窗高度

### 2. 页面头部 (.page-header)
现代化的页面头部组件，包含标题、描述和操作按钮。

**结构：**
- `.header-content`: 头部内容容器
- `.title-section`: 标题区域
- `.page-title`: 页面标题（带图标）
- `.page-description`: 页面描述
- `.header-actions`: 操作按钮区域

### 3. 统计卡片 (.modern-stats-grid)
现代化的统计卡片网格，用于展示关键数据指标。

**特性：**
- 响应式网格布局
- 悬停效果和动画
- 图标和趋势指示器
- 装饰性背景图案

**卡片结构：**
- `.stat-card`: 统计卡片容器
- `.stat-header`: 卡片头部（图标和趋势）
- `.stat-content`: 卡片内容（数字、标签、描述）

### 4. 内容卡片 (.modern-content-card)
现代化的内容卡片，用于包装表格、表单等内容。

**结构：**
- `.card-header`: 卡片头部
- `.header-content`: 头部内容
- `.header-left`: 左侧标题区域
- `.header-actions`: 右侧操作区域
- `.card-content`: 卡片内容区域

### 5. 搜索区域 (.modern-search-section)
现代化的搜索和筛选区域。

**特性：**
- 灵活的布局
- 统一的输入框和按钮样式
- 响应式设计

### 6. 按钮系统 (.modern-btn)
统一的按钮样式系统，支持多种类型和状态。

**类型：**
- `.primary`: 主要按钮
- `.secondary`: 次要按钮
- `.success`: 成功按钮
- `.warning`: 警告按钮
- `.danger`: 危险按钮
- `.ghost`: 幽灵按钮
- `.text`: 文本按钮

### 7. 表格样式 (.modern-table)
现代化的表格样式，提供清晰的数据展示。

**特性：**
- 圆角边框
- 统一的头部样式
- 悬停效果
- 清晰的分割线

## 使用指南

### 1. 引入样式
在 `app.scss` 中引入现代化设计系统：
```scss
@import './modern-design.scss';
```

### 2. 页面结构
```html
<div class="modern-page-container">
  <!-- 页面头部 -->
  <div class="page-header">
    <div class="header-content">
      <div class="title-section">
        <h1 class="page-title">
          <div class="title-icon">
            <el-icon><Icon /></el-icon>
          </div>
          页面标题
        </h1>
        <p class="page-description">页面描述</p>
      </div>
      <div class="header-actions">
        <el-button class="modern-btn primary">操作按钮</el-button>
      </div>
    </div>
  </div>

  <!-- 统计卡片 -->
  <div class="modern-stats-grid">
    <div class="stat-card">
      <div class="stat-header">
        <div class="stat-icon primary">
          <el-icon><Icon /></el-icon>
        </div>
        <div class="stat-trend up">+12%</div>
      </div>
      <div class="stat-content">
        <div class="stat-number">123</div>
        <div class="stat-label">标签</div>
        <div class="stat-description">描述</div>
      </div>
    </div>
  </div>

  <!-- 内容卡片 -->
  <div class="modern-content-card">
    <div class="card-header">
      <div class="header-content">
        <div class="header-left">
          <h3 class="card-title">卡片标题</h3>
          <p class="card-subtitle">卡片副标题</p>
        </div>
        <div class="header-actions">
          <el-button class="modern-btn secondary">操作</el-button>
        </div>
      </div>
    </div>
    <div class="card-content">
      <!-- 内容区域 -->
    </div>
  </div>
</div>
```

### 3. 按钮使用
```html
<!-- 主要按钮 -->
<el-button class="modern-btn primary">主要操作</el-button>

<!-- 次要按钮 -->
<el-button class="modern-btn secondary">次要操作</el-button>

<!-- 带图标的按钮 -->
<el-button class="modern-btn primary" :icon="Plus">添加</el-button>
```

## 最佳实践

1. **保持一致性**：始终使用设计系统中定义的颜色、间距和组件样式
2. **响应式设计**：确保在不同屏幕尺寸下都有良好的显示效果
3. **可访问性**：注意颜色对比度和键盘导航支持
4. **性能优化**：合理使用动画和过渡效果，避免影响性能

## 更新日志

### v1.0.0 (2024-01-22)
- 初始版本发布
- 完成基础设计系统搭建
- 实现页面容器、统计卡片、内容卡片等核心组件
- 完成 DNS 证书管理和域名解析页面的现代化改造
