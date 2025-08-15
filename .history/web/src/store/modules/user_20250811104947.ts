import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { router } from '@/router'
import { useSettingStore } from './setting'
import { useWorktabStore } from './worktab'
import { AppRouteRecord } from '@/types/router'
import { setPageTitle } from '@/router/utils/utils'
import { resetRouterState } from '@/router/guards/beforeEach'
import { useRouteRegistration } from '@/composables/useRouteRegistration'
import { RoutesAlias } from '@/router/routesAlias'

// 用户
export const useUserStore = defineStore(
  'userStore',
  () => {
    const isLogin = ref(false)
    const isLock = ref(false)
    const lockPassword = ref('')
    const info = ref<Partial<Api.User.UserInfo>>({})
    const searchHistory = ref<AppRouteRecord[]>([])
    const accessToken = ref('')
    const refreshToken = ref('')

    const getUserInfo = computed(() => info.value)
    const getSettingState = computed(() => useSettingStore().$state)
    const getWorktabState = computed(() => useWorktabStore().$state)

    const setUserInfo = (newInfo: Api.User.UserInfo) => {
      info.value = newInfo
    }

    const setLoginStatus = (status: boolean) => {
      isLogin.value = status
      // 如果登录状态变为true，重置路由注册状态，确保重新注册动态路由
      if (status) {
        console.log('用户登录状态已设置，准备注册动态路由')
      }
    }

    const setSearchHistory = (list: AppRouteRecord[]) => {
      searchHistory.value = list
    }

    const setLockStatus = (status: boolean) => {
      isLock.value = status
    }

    const setLockPassword = (password: string) => {
      lockPassword.value = password
    }

    const setToken = (newAccessToken: string, newRefreshToken?: string) => {
      accessToken.value = newAccessToken
      if (newRefreshToken) {
        refreshToken.value = newRefreshToken
      }
    }

    const logOut = () => {
      info.value = {}
      isLogin.value = false
      isLock.value = false
      lockPassword.value = ''
      accessToken.value = ''
      refreshToken.value = ''
      useWorktabStore().opened = []
      sessionStorage.removeItem('iframeRoutes')
      resetRouterState()
      // 重置路由注册状态
      const { resetRouteRegistration } = useRouteRegistration()
      resetRouteRegistration()
      router.push(RoutesAlias.Login)
    }

    return {
      isLogin,
      isLock,
      lockPassword,
      info,
      searchHistory,
      accessToken,
      refreshToken,
      getUserInfo,
      getSettingState,
      getWorktabState,
      setUserInfo,
      setLoginStatus,
      setSearchHistory,
      setLockStatus,
      setLockPassword,
      setToken,
      logOut
    }
  },
  {
    persist: {
      key: 'user',
      storage: localStorage
    }
  }
)
