import App from './App.vue'
import { createApp, defineComponent, h } from 'vue'
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

document.addEventListener(
  'touchstart',
  function () {},
  { passive: false }
)

// 全局错误边界组件
const ErrorBoundary = defineComponent({
  name: 'ErrorBoundary',
  setup(_, { slots }) {
    const error = ref(false)
    const errorMsg = ref('')

    onErrorCaptured((err) => {
      error.value = true
      errorMsg.value = err.message || '未知错误'
      console.error('全局捕获到渲染错误:', err)
      return false // 阻止错误继续传播
    })

    return () => {
      if (error.value) {
        return h('div', { class: 'global-error-container' }, [
          h('h3', { class: 'error-title' }, '页面渲染错误'),
          h('p', { class: 'error-message' }, errorMsg.value),
          h('button', {
            class: 'error-retry-btn',
            onClick: () => window.location.reload()
          }, '刷新页面')
        ])
      }

      return slots.default?.()
    }
  }
})

// 全局未捕获异常处理
window.addEventListener('error', (event) => {
  console.error('全局错误:', event.error)
  ElMessage.error(`系统错误: ${event.error?.message || '未知错误'}`)
  event.preventDefault()
})

// 全局Promise拒绝错误处理
window.addEventListener('unhandledrejection', (event) => {
  console.error('未处理的Promise拒绝:', event.reason)
  ElMessage.error(`系统错误: ${event.reason?.message || '未知错误'}`)
  event.preventDefault()
})

const app = createApp(App)
// 注册全局错误边界组件
app.component('ErrorBoundary', ErrorBoundary)

initStore(app)
initRouter(app)
setupGlobDirectives(app)

app.use(language)

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 全局错误处理器
app.config.errorHandler = (err, vm, info) => {
  console.error('Vue错误处理器:', err, info)
  ElMessage.error(`系统错误: ${err.message || '未知错误'}`)
}

// 将整个应用包裹在错误边界中
const AppWithErrorBoundary = defineComponent({
  render() {
    return h(ErrorBoundary, {}, {
      default: () => [h(App)]
    })
  }
})

app.mount('#app', AppWithErrorBoundary)

