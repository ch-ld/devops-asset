import { getTabConfig } from '@/utils/ui'
import { useSettingStore } from '@/store/modules/setting'
import { useMenuStore } from '@/store/modules/menu'

// 通用函数
export function useCommon() {
  const settingStore = useSettingStore()
  const { showWorkTab, tabStyle } = storeToRefs(settingStore)

  // 是否是前端控制模式
  const isFrontendMode = computed(() => {
    // 暂时强制使用后端模式获取菜单
    return false
    // return import.meta.env.VITE_ACCESS_MODE === 'frontend'
  })

  // 首页路径
  const homePath = computed(() => useMenuStore().getHomePath())

  // 刷新页面
  const refresh = () => {
    settingStore.reload()
  }

  // 回到顶部
  const scrollToTop = () => {
    window.scrollTo({ top: 0 })
  }

  // 页面最小高度
  const containerMinHeight = computed(() => {
    const { openHeight, closeHeight } = getTabConfig(tabStyle.value)
    return `calc(100vh - ${showWorkTab.value ? openHeight : closeHeight}px)`
  })

  return {
    isFrontendMode,
    homePath,
    refresh,
    scrollToTop,
    containerMinHeight
  }
}
