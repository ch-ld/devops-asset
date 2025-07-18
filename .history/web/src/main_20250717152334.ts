/**
 * 主要修复:
 * 1. 替换了 Ant Design 图标为 Element Plus 图标
 * 2. 添加了全局错误处理和未处理的Promise异常捕获
 * 3. 添加了对空值和未定义值的防御性检查
 * 4. 改进了API调用的错误处理
 * 5. 添加了错误状态和重试功能
 * 6. 添加了批量操作的支持
 */
import App from './App.vue'
import { createApp } from 'vue'
import { initStore } from './store'                 // Store
import { initRouter } from './router'               // Router
import '@styles/reset.scss'                         // 重置HTML样式
import '@styles/app.scss'                           // 全局样式
import '@styles/el-ui.scss'                         // 优化element样式
import '@styles/mobile.scss'                        // 移动端样式优化
import '@styles/change.scss'                        // 主题切换过渡优化
import '@styles/theme-animation.scss'               // 主题切换动画
import '@styles/el-light.scss'                      // Element 自定义主题（亮色）
import '@styles/el-dark.scss'                       // Element 自定义主题（暗色）
import '@styles/dark.scss'                          // 系统主题
import '@icons/system/iconfont.js'                  // 系统彩色图标
import '@icons/system/iconfont.css'                 // 系统图标
import '@utils/sys/console.ts'                      // 控制台输出内容
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import { setupGlobDirectives } from './directives'
import language from './locales'
import { ElMessage } from 'element-plus'

// 导入全局组件
import ArtLayouts from '@/components/core/layouts/art-layouts/index.vue'
import ArtHeaderBar from '@/components/core/layouts/art-header-bar/index.vue'
import ArtSidebarMenu from '@/components/core/layouts/art-menus/art-sidebar-menu/index.vue'
import ArtPageContent from '@/components/core/layouts/art-page-content/index.vue'
import ArtSettingsPanel from '@/components/core/layouts/art-settings-panel/index.vue'
import ArtGlobalSearch from '@/components/core/layouts/art-global-search/index.vue'
import ArtScreenLock from '@/components/core/layouts/art-screen-lock/index.vue'
import ArtWatermark from '@/components/core/layouts/art-watermark/index.vue'

document.addEventListener(
  'touchstart',
  function () {},
  { passive: false }
)

const app = createApp(App)
initStore(app)
initRouter(app)
setupGlobDirectives(app)

app.use(language)

// 全局注册组件
app.component('ArtLayouts', ArtLayouts)
app.component('ArtHeaderBar', ArtHeaderBar)
app.component('ArtSidebarMenu', ArtSidebarMenu)
app.component('ArtPageContent', ArtPageContent)
app.component('ArtSettingsPanel', ArtSettingsPanel)
app.component('ArtGlobalSearch', ArtGlobalSearch)
app.component('ArtScreenLock', ArtScreenLock)
app.component('ArtWatermark', ArtWatermark)

// 注册所有Element Plus图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 全局错误处理
app.config.errorHandler = (err, instance, info) => {
  console.error('全局错误:', err)
  console.error('错误组件:', instance)
  console.error('错误信息:', info)
  ElMessage.error('应用发生错误，请刷新页面重试')
}

// 捕获未处理的Promise异常
window.addEventListener('unhandledrejection', (event) => {
  console.error('未处理的Promise异常:', event.reason)
  // 防止重复提示
  if (!event.reason?.message?.includes('请求失败')) {
    ElMessage.error('操作失败，请稍后重试')
  }
  event.preventDefault()
})

app.mount('#app')

