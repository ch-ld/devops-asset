<template>
  <div class="quick-nav-ant">
    <div class="nav-group-bar">
      <div class="nav-group-tabs">
        <button
          v-for="(group, _idx) in groupTabs"
          :key="group"
          :class="[
            'ant-btn',
            'ant-btn-round',
            currentGroup === group ? 'ant-btn-primary' : 'ant-btn-default'
          ]"
          @click="currentGroup = group"
        >
          {{ group }}
        </button>
      </div>
      <button class="ant-btn ant-btn-default ant-btn-round manage-nav-btn" @click="onManageNav">
        <i class="anticon anticon-setting"></i> 管理导航
      </button>
    </div>
    <div class="nav-group-list">
      <div class="ant-row nav-row">
        <div class="ant-col nav-col" v-for="nav in filteredNavs" :key="nav.id">
          <div class="nav-card-ant">
            <div class="nav-card-header">
              <img
                class="nav-logo"
                :src="getNavLogo(nav)"
                :alt="nav.title"
                @error="onLogoError($event, nav)"
                loading="lazy"
              />
              <div class="nav-title">
                <a :href="nav.url" target="_blank">{{ nav.title }}</a>
                <div class="nav-desc">{{ nav.description }}</div>
              </div>
              <div class="nav-more ant-dropdown" @click.stop="toggleMenu(nav.id)">
                <i class="anticon anticon-ellipsis"></i>
                <div class="dropdown-menu" v-show="menuOpenId === nav.id">
                  <div class="dropdown-item" @click.stop="onEditNav(nav)"
                    ><i class="anticon anticon-edit"></i> 编辑</div
                  >
                  <div class="dropdown-item danger" @click.stop="onDeleteNav(nav.id)"
                    ><i class="anticon anticon-delete"></i> 删除</div
                  >
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <!-- 导航编辑弹窗 -->
    <NavEditModal
      ref="navEditModalRef"
      :group-options="groupTabs.filter((g) => g !== '便捷导航')"
      @refresh="loadNavData"
    />
    <NavManageModal
      ref="navManageModalRef"
      :nav-list="navList"
      :get-nav-logo="getNavLogo"
      @edit="onManageEdit"
      @delete="onManageDelete"
      @add="onManageAdd"
      @move="onManageMove"
      @batchDelete="onManageBatchDelete"
    />
  </div>
</template>

<script setup lang="ts">
  import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { getNavList, deleteNav } from '@/api/system/nav'
  import { ApiStatus } from '@/utils/http/status'
  import NavEditModal from './modal/NavEditModal.vue'
  import NavManageModal from './NavManageModal.vue'
  import defaultLogo from '@/assets/img/common/logo.png'

  interface Nav {
    id: number | string
    title: string
    url: string
    group?: string
    icon?: string
    description?: string
    [key: string]: any
  }

  const navList = ref<Nav[]>([])
  const navEditModalRef = ref()
  const navManageModalRef = ref()
  const currentGroup = ref('便捷导航')
  const logoErrorMap = ref<Record<string, boolean>>({})
  // favicon 缓存，避免重复请求
  const faviconCache = ref<Record<string, string>>({})
  const menuOpenId = ref<number | string | null>(null)

  // 计算所有分组Tab
  const groupTabs = computed(() => {
    const groups = navList.value.map((nav: Nav) => nav.group || '未分组')
    const unique = Array.from(new Set(groups))
    return ['便捷导航', ...unique]
  })

  // 当前分组下的导航
  const filteredNavs = computed<Nav[]>(() => {
    if (currentGroup.value === '便捷导航') return navList.value
    return navList.value.filter((nav) => (nav.group || '未分组') === currentGroup.value)
  })

  // 加载导航数据
  const loadNavData = async (): Promise<void> => {
    try {
      const res: any = await getNavList()
      if (res.code === ApiStatus.success) {
        // 修复：获取 res.data.list 而不是 res.data
        const rawNavs = res.data?.list || []

        // 修复：映射后端字段到前端期望的字段
        navList.value = rawNavs.map((nav) => ({
          id: nav.id,
          title: nav.name, // 后端的 name -> 前端的 title
          url: nav.links, // 后端的 links -> 前端的 url
          group: nav.group_name, // 后端的 group_name -> 前端的 group
          icon: nav.icon_url, // 后端的 icon_url -> 前端的 icon
          description: nav.description,
          open_in_new_tab: nav.open_in_new_tab,
          status: nav.status,
          isEnable: nav.status === 1, // 后端的 status -> 前端的 isEnable (1:启用, 2:禁用)
          sort: nav.order_num, // 后端的 order_num -> 前端的 sort
          openInNewTab: nav.open_in_new_tab,
          created_at: nav.created_at,
          updated_at: nav.updated_at
        }))

        console.log('导航数据加载成功:', navList.value)
      } else {
        console.error('导航数据加载失败:', res.message)
        ElMessage.error(`加载导航数据失败: ${res.message}`)
      }
    } catch (error) {
      console.error('导航数据加载异常:', error as any)
      ElMessage.error('加载导航数据失败')
    }
  }

  // 智能获取 logo，带 favicon 缓存
  function getNavLogo(nav: Nav): string {
    if (logoErrorMap.value[nav.id]) return defaultLogo
    if (nav.icon) return nav.icon
    // favicon 缓存命中
    if (faviconCache.value[nav.id]) return faviconCache.value[nav.id]
    // 自动抓取 favicon
    try {
      const url = new URL(nav.url)
      const favicon = url.origin + '/favicon.ico'
      faviconCache.value[nav.id] = favicon
      return favicon
    } catch {
      faviconCache.value[nav.id] = defaultLogo
      return defaultLogo
    }
  }
  function onLogoError(e: Event, nav: Nav): void {
    logoErrorMap.value[nav.id] = true
    faviconCache.value[nav.id] = defaultLogo
    const img = e.target as HTMLImageElement | null
    if (img) img.src = defaultLogo
  }

  // 编辑导航
  const onEditNav = (nav: Nav) => {
    navEditModalRef.value?.showModal('edit', nav)
  }

  // 删除导航
  const onDeleteNav = async (id: number | string) => {
    try {
      await ElMessageBox.confirm('确定要删除这个导航吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })

      const res = await deleteNav(id)
      if (res.code === ApiStatus.success) {
        ElMessage.success('删除成功')
        loadNavData()
      } else {
        ElMessage.error(`删除失败: ${res.message}`)
      }
    } catch (error) {
      if (error !== 'cancel') {
        ElMessage.error('删除失败')
      }
    }
  }

  const onManageNav = () => {
    navManageModalRef.value?.show()
  }

  function onManageEdit(nav: Nav) {
    onEditNav(nav)
  }
  function onManageDelete(id: number | string) {
    onDeleteNav(id)
  }
  function onManageAdd() {
    navEditModalRef.value?.showModal('add')
  }
  function onManageMove({ from, to }: { from: number; to: number }) {
    // 参数验证
    if (from === to) return

    if (!navList.value || !Array.isArray(navList.value)) {
      ElMessage.error('数据异常，请刷新页面重试')
      return
    }

    if (from < 0 || from >= navList.value.length || to < 0 || to >= navList.value.length) {
      ElMessage.error('移动位置无效')
      return
    }

    try {
      // 执行移动操作
      const arr = navList.value.slice()
      const moved = arr.splice(Number(from), 1)[0]
      arr.splice(Number(to), 0, moved)
      navList.value = arr

      // 显示成功消息
      ElMessage.success(`已将 "${moved.title}" 移动到新位置`)

      // TODO: 调用后端保存排序
      // 这里可以添加后端保存逻辑
    } catch (err) {
      ElMessage.error('移动失败，请重试')
    }
  }

  function onManageBatchDelete(ids: (number | string)[]): void {
    Promise.all(ids.map((id) => onDeleteNav(id))).then(() => {
      ElMessage.success('批量删除成功')
      loadNavData()
    })
  }

  function toggleMenu(navId: number | string): void {
    menuOpenId.value = menuOpenId.value === navId ? null : navId
  }
  function handleClickOutside(e: MouseEvent): void {
    const target = e.target as HTMLElement | null
    if (!target?.closest('.nav-more')) menuOpenId.value = null
  }

  onMounted(() => {
    document.addEventListener('click', handleClickOutside)
    loadNavData()
  })
  onBeforeUnmount(() => {
    document.removeEventListener('click', handleClickOutside)
  })
</script>

<style lang="scss" scoped>
  .quick-nav-ant {
    .nav-group-bar {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 18px;
    }
    .nav-group-tabs {
      display: flex;
      align-items: center;
      flex-wrap: wrap;
      gap: 10px;
    }
    .manage-nav-btn {
      display: flex;
      align-items: center;
      font-size: 15px;
      gap: 6px;
      padding: 0 18px;
      height: 36px;
      border-radius: 20px;
      background: #fff;
      border: 1px solid #e5e6eb;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.03);
      transition: all 0.2s;
      &:hover {
        border-color: #1677ff;
        color: #1677ff;
        background: #f0f7ff;
      }
    }
    .nav-group-list {
      .nav-row {
        display: flex;
        flex-wrap: wrap;
        gap: 24px 0;
      }
      .nav-col {
        flex: 0 0 25%;
        max-width: 25%;
        padding: 0 12px;
        box-sizing: border-box;
      }
      .nav-card-ant {
        background: #fff;
        border-radius: 18px;
        box-shadow: 0 4px 24px rgba(22, 119, 255, 0.08);
        padding: 0;
        margin-bottom: 24px;
        transition:
          box-shadow 0.2s,
          transform 0.2s;
        &:hover {
          box-shadow: 0 8px 32px rgba(22, 119, 255, 0.16);
          transform: translateY(-2px) scale(1.03);
        }
        .nav-card-header {
          display: flex;
          align-items: center;
          padding: 22px 24px 22px 22px;
          position: relative;
          .nav-logo {
            width: 54px;
            height: 54px;
            border-radius: 12px;
            margin-right: 20px;
            object-fit: cover;
            background: #f5f6fa;
            border: 1px solid #e5e6eb;
          }
          .nav-title {
            flex: 1;
            display: flex;
            flex-direction: column;
            a {
              color: #222;
              font-size: 18px;
              font-weight: 600;
              text-decoration: none;
              margin-bottom: 2px;
              &:hover {
                color: #1677ff;
              }
            }
            .nav-desc {
              color: #888;
              font-size: 12px;
              margin-top: 2px;
              line-height: 1.5;
              word-break: break-all;
              display: -webkit-box;
              -webkit-line-clamp: 3;
              -webkit-box-orient: vertical;
              overflow: hidden;
              text-overflow: ellipsis;
            }
          }
          .nav-more {
            margin-left: 10px;
            color: #bbb;
            font-size: 22px;
            cursor: pointer;
            position: relative;
            &:hover .dropdown-menu {
              display: block;
            }
            .dropdown-menu {
              display: none;
              position: absolute;
              right: 0;
              top: 32px;
              min-width: 110px;
              background: #fff;
              border-radius: 8px;
              box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
              z-index: 10;
              .dropdown-item {
                padding: 10px 18px;
                font-size: 15px;
                color: #222;
                cursor: pointer;
                transition: background 0.15s;
                display: flex;
                align-items: center;
                gap: 6px;
                &:hover {
                  background: #f0f7ff;
                  color: #1677ff;
                }
                &.danger {
                  color: #ff4d4f;
                  &:hover {
                    background: #fff1f0;
                    color: #ff4d4f;
                  }
                }
              }
            }
          }
        }
      }
    }
    // Ant Design 按钮样式
    .ant-btn {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      height: 36px;
      padding: 0 18px;
      font-size: 15px;
      border-radius: 20px;
      border: 1px solid #e5e6eb;
      background: #fff;
      color: #222;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.2s;
      outline: none;
      box-shadow: 0 2px 8px rgba(22, 119, 255, 0.03);
      &:hover {
        border-color: #1677ff;
        color: #1677ff;
        background: #f0f7ff;
      }
      &.ant-btn-primary {
        background: linear-gradient(90deg, #1677ff 0%, #69b1ff 100%);
        color: #fff;
        border: none;
        box-shadow: 0 2px 8px rgba(22, 119, 255, 0.1);
        &:hover {
          background: linear-gradient(90deg, #4096ff 0%, #69b1ff 100%);
          color: #fff;
        }
      }
      &.ant-btn-round {
        border-radius: 20px;
      }
      &.ant-btn-default {
        background: #fff;
        color: #222;
        border: 1px solid #e5e6eb;
      }
    }
    // Ant Design 图标
    .anticon {
      display: inline-block;
      font-style: normal;
      font-size: 18px;
      line-height: 1;
      vertical-align: middle;
      &.anticon-setting:before {
        content: '\2699';
      }
      &.anticon-edit:before {
        content: '\270E';
      }
      &.anticon-delete:before {
        content: '\1F5D1';
      }
      &.anticon-ellipsis:before {
        content: '\2026';
      }
    }
  }
</style>
