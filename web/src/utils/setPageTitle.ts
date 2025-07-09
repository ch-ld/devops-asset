import { RouteLocationNormalized } from 'vue-router'
import { formatMenuTitle } from '@/router/utils/utils'
import AppConfig from '@/config'

/**
 * 设置页面标题，根据路由元信息和系统信息拼接标题
 * @param to 当前路由对象
 */
export const setPageTitle = (): void => {
  setTimeout(() => {
    document.title = AppConfig.systemInfo.name
  }, 150)
}
