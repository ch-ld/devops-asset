<template>
  <template v-for="item in filteredMenuItems" :key="item.path">
    <!-- 包含子菜单的项目 -->
    <el-sub-menu v-if="hasChildren(item)" :index="item.path || item.meta.title" :level="level">
      <template #title>
        <MenuItemIcon :icon="item.meta.icon" :color="theme?.iconColor" />
        <span class="menu-name">{{ formatMenuTitle(item.meta.title) }}</span>
        <div v-if="item.meta.showBadge" class="badge" style="right: 35px" />
      </template>
      <SidebarSubmenu
        :list="item.children"
        :is-mobile="isMobile"
        :level="level + 1"
        :theme="theme"
        @close="closeMenu"
      />
    </el-sub-menu>

    <!-- 普通菜单项 -->
    <el-menu-item
      v-else
      :index="item.path || item.meta.title"
      :level-item="level + 1"
      @click="handleMenuClick(item)"
    >
      <MenuItemIcon :icon="item.meta.icon" :color="theme?.iconColor" />
      <template #title>
        <span class="menu-name">{{ formatMenuTitle(item.meta.title) }}</span>
        <div v-if="item.meta.showBadge" class="badge" />
        <div v-if="item.meta.showTextBadge" class="text-badge">
          {{ item.meta.showTextBadge }}
        </div>
      </template>
    </el-menu-item>
  </template>
</template>

<script setup lang="ts">
  import { computed } from 'vue'
  import type { AppRouteRecord } from '@/types/router'
  import { formatMenuTitle } from '@/router/utils/utils'
  import { handleMenuJump } from '@/utils/navigation'
  import { router } from '@/router'
  import { ElMessage } from 'element-plus'

  // 类型定义
  interface Props {
    title?: string
    list?: AppRouteRecord[]
    theme?: {
      iconColor?: string
    }
    isMobile?: boolean
    level?: number
  }

  // Props定义
  const props = withDefaults(defineProps<Props>(), {
    title: '',
    list: () => [],
    theme: () => ({}),
    isMobile: false,
    level: 0
  })

  // Emits定义
  const emit = defineEmits<{
    (e: 'close'): void
  }>()

  // 计算属性
  const safeList = Array.isArray(props.list) ? props.list : []
  const filteredMenuItems = computed(() => filterRoutes(safeList))

  // 处理菜单点击
  const handleMenuClick = async (item: AppRouteRecord) => {
    try {
      closeMenu()

      // 检查路由是否已注册 - 优先使用 name 判断
      const routeName = item.name || item.path
      if (!router.hasRoute(routeName)) {
        console.warn(`路由 ${routeName} 未注册，尝试重新注册...`)
        // 等待一小段时间让路由注册完成
        await new Promise((resolve) => setTimeout(resolve, 500))

        // 再次检查路由
        if (!router.hasRoute(routeName)) {
          ElMessage.warning('页面正在加载中，请稍后再试')
          return
        }
      }

      await handleMenuJump(item)
    } catch (error) {
      console.error('菜单跳转失败:', error)
      ElMessage.error('页面跳转失败，请刷新页面重试')
    }
  }

  // 关闭菜单
  const closeMenu = () => emit('close')

  // 判断是否有子菜单
  const hasChildren = (item: AppRouteRecord): boolean => {
    return Boolean(item.children?.length)
  }

  // 过滤菜单项
  const filterRoutes = (items: AppRouteRecord[]): AppRouteRecord[] => {
    if (!Array.isArray(items)) return []
    return items
      .filter((item) => !item.meta.isHide)
      .map((item) => ({
        ...item,
        children: Array.isArray(item.children) ? filterRoutes(item.children) : []
      }))
  }
</script>

<script lang="ts">
  // 抽取图标组件
  const MenuItemIcon = defineComponent({
    name: 'MenuItemIcon',
    props: {
      icon: String,
      color: String
    },
    setup(props) {
      return () =>
        h('i', {
          class: 'menu-icon iconfont-sys',
          style: props.color ? { color: props.color } : undefined,
          innerHTML: props.icon
        })
    }
  })
</script>
