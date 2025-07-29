# Element Plus Popper 白块问题修复记录

> 创建日期：2025-07-29  
> 适用版本：Art Design Pro + Element Plus `^2.x`

## 1. 问题现象

页面中 **Tooltip / Dropdown / Select 等弹层** 出现不可点击的白色矩形（包含黑色小箭头），如下图示意：

- 弹层内容未渲染
- 位置随机漂浮
- 滚动页面仍然固定在视窗中

## 2. 根因分析

1. 为解决「厂商管理-区域选择」下拉框被遮挡问题，`web/src/assets/styles/app.scss` **169 ~ 214 行** 写了“一锅端”式覆盖：

   ```scss
   .el-select-dropdown,
   .el-select__popper,
   .el-popper.el-select__popper,
   .el-popper[data-popper-placement],
   div[data-popper-placement],
   .el-popper,
   [class*="el-popper"],
   [id*="el-popper"] {
     position: fixed !important;
     opacity: 1 !important;
     background-color: #fff !important;
     z-index: 99999999 !important;
     /* … */
   }
   ```

2. Element Plus 在 **弹层隐藏** 时仅设置 `display:none`、`opacity:0`。上面的覆盖又把 `opacity` 强制成 `1`，`display:none` 仍然生效 ⇒ 结果就是 _空白弹层_ 显示出来。
3. `position:fixed` 让这些空白层漂浮在视口任意位置，导致“白块”现象。

## 3. 修复方案

### 3.1 精准提升层级

仅对 **Select 下拉框等易被遮挡组件** 提升 `z-index`，_不再_ 修改 `position / opacity / background-color`。

```scss
/* web/src/assets/styles/app.scss */
/* 全局修复 Element Plus 下拉框层级 */
.el-select-dropdown,
.el-select__popper,
.el-popper.el-select__popper {
  z-index: 3000 !important;   // 默认 2000，足够压过 Dialog(2001)
}
```

### 3.2 通过框架配置统一 zIndex

在入口 `web/src/main.ts`：

```ts
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

app.use(ElementPlus, { zIndex: 3000 })
```

这样所有基于 Popper 的组件（Tooltip / Dropdown / Cascader / DatePicker …）都会继承 3000 层级，不必再写额外 CSS。

### 3.3 删除有副作用的覆盖

移除下列选择器中的 `position / opacity / background-color` 等强制样式：

```scss
.el-popper[data-popper-placement]
div[data-popper-placement]
.el-popper
[class*="el-popper"]
[id*="el-popper"]
[role="listbox"], [role="option"]
```

> 仅保留必要的 **z-index** 覆盖（如需）。

## 4. 回归测试 Checklist

- [ ] Tooltip / Popover 悬停显示、消失正常
- [ ] Select / Cascader / DatePicker 在 Dialog / Drawer 内不被遮罩遮挡
- [ ] Dropdown 菜单在表格或复杂容器中仍可点击
- [ ] 无任何空白矩形浮层残留

## 5. 后续建议

1. 尽量使用 Element Plus 的 `teleported` 默认行为，避免被父级 `overflow` 裁剪。
2. 若特定场景仍被裁剪，可局部关闭 `teleported` 并为该组件单独写 `position:fixed` 覆盖，_不要全局覆盖_。
3. 统一层级请优先使用 `app.use(ElementPlus,{zIndex})`，避免重复 CSS Hack。

---

如再次出现类似问题，可先排查：

- 弹层元素是否含有 `.el-popper` 且 `display:none`；
- 是否被新的全局样式覆盖了 `opacity`；
- 父级容器是否设置了 `overflow:hidden`。 
