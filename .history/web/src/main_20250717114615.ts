import { createApp } from 'vue'
import App from './App.vue'
import { router } from './router'
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import { ElMessage } from 'element-plus'

// 导入样式
import './assets/styles/app.scss'
import 'element-plus/dist/index.css'

// 创建应用
const app = createApp(App)

// 全局错误处理
app.config.errorHandler = (err, instance, info) => {
  console.error('全局错误:', err)
  console.error('错误信息:', info)
  ElMessage.error('应用发生错误，请刷新页面重试')
}

// 捕获未处理的Promise错误
window.addEventListener('unhandledrejection', (event) => {
  console.error('未处理的Promise错误:', event.reason)
  // 阻止默认处理（如控制台日志）
  event.preventDefault()
})

// 初始化 Pinia
const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

// 挂载应用
app.use(pinia)
app.use(router)
app.mount('#app')

