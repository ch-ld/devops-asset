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

<script setup>
  import { computed } from 'vue'
  import { formatMenuTitle } from '@/router/utils/utils'
  import { handleMenuJump } from '@/utils/navigation'
  import { ElMessage } from 'element-plus'
  import { useRouter } from 'vue-router'

  // Props定义
  const props = defineProps({
    title: {
      type: String,
      default: ''
    },
    list: {
      type: Array,
      default: () => []
    },
    theme: {
      type: Object,
      default: () => ({})
    },
    isMobile: {
      type: Boolean,
      default: false
    },
    level: {
      type: Number,
      default: 0
    }
  })

  // Emits定义
  const emit = defineEmits(['close'])
  const router = useRouter()

  // 计算属性
  const safeList = computed(() => (Array.isArray(props.list) ? props.list : []))
  const filteredMenuItems = computed(() => filterRoutes(safeList.value))

  // 处理菜单点击
  const handleMenuClick = async (item) => {
    try {
      closeMenu()

      // 检查路由是否存在
      const routePath = item.path
      const routeName = item.name

      if (!routePath) {
        console.error('菜单路径为空:', item)
        ElMessage.warning('菜单配置错误，请联系管理员')
        return
      }

      // 尝试跳转路由
      try {
        await handleMenuJump(item)
      } catch (error) {
        console.error('菜单跳转失败:', error)
        ElMessage.warning('页面加载中，请稍后再试')
      }
    } catch (error) {
      console.error('菜单处理错误:', error)
    }
  }

  // 关闭菜单
  const closeMenu = () => emit('close')

  // 判断是否有子菜单
  const hasChildren = (item) => {
    return Boolean(item.children?.length)
  }

  // 过滤菜单项
  const filterRoutes = (items) => {
    if (!Array.isArray(items)) return []
    return items
      .filter((item) => item && !item.meta?.isHide)
      .map((item) => ({
        ...item,
        children: Array.isArray(item.children) ? filterRoutes(item.children) : []
      }))
  }
</script>

<script>
  import { defineComponent, h } from 'vue'

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
