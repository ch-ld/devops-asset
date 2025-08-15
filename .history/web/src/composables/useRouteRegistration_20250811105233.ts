import { ref, readonly } from 'vue'

// 路由注册状态
const isRouteRegistered = ref(false)
const routeRegistrationPromise = ref<Promise<void> | null>(null)

/**
 * 路由注册状态管理
 */
export function useRouteRegistration() {
  // 设置路由注册状态
  const setRouteRegistered = (status: boolean) => {
    isRouteRegistered.value = status
    console.log('路由注册状态已更新:', status)
  }

  // 等待路由注册完成
  const waitForRouteRegistration = (): Promise<void> => {
    if (isRouteRegistered.value) {
      return Promise.resolve()
    }

    // 如果已经有等待的Promise，复用它
    if (routeRegistrationPromise.value) {
      return routeRegistrationPromise.value
    }

    // 创建新的等待Promise
    routeRegistrationPromise.value = new Promise((resolve) => {
      const maxWait = 15000 // 最多等待15秒
      const checkInterval = 100 // 每100ms检查一次
      let waitTime = 0
      
      const checkRegistration = () => {
        if (isRouteRegistered.value) {
          console.log('路由注册已完成')
          routeRegistrationPromise.value = null
          resolve()
        } else if (waitTime >= maxWait) {
          console.warn('等待路由注册超时')
          routeRegistrationPromise.value = null
          resolve()
        } else {
          waitTime += checkInterval
          setTimeout(checkRegistration, checkInterval)
        }
      }
      
      checkRegistration()
    })

    return routeRegistrationPromise.value
  }

  // 重置路由注册状态
  const resetRouteRegistration = () => {
    isRouteRegistered.value = false
    routeRegistrationPromise.value = null
    console.log('路由注册状态已重置')
  }

  return {
    isRouteRegistered: readonly(isRouteRegistered),
    setRouteRegistered,
    waitForRouteRegistration,
    resetRouteRegistration
  }
}
