<template>
  <div class="modern-sidebar">
    <!-- 头部区域 -->
    <div class="sidebar-header">
      <div class="header-content">
        <div class="title-section">
          <div class="icon-wrapper">
            <el-icon class="title-icon"><FolderOpened /></el-icon>
          </div>
          <h3 class="title">主机分组</h3>
        </div>
        <div class="header-actions">
          <el-tooltip content="新建分组" placement="bottom">
            <el-button 
              class="action-btn" 
              circle 
              size="small" 
              @click="handleAddGroup"
              :icon="Plus"
            />
          </el-tooltip>
          <el-tooltip content="刷新" placement="bottom">
            <el-button 
              class="action-btn" 
              circle 
              size="small" 
              @click="handleRefresh"
              :loading="loading"
              :icon="Refresh"
            />
          </el-tooltip>
        </div>
      </div>
    </div>

    <!-- 搜索区域 -->
    <div class="search-section">
      <el-input
        v-model="searchQuery"
        placeholder="搜索分组..."
        size="default"
        clearable
        class="search-input"
      >
        <template #prefix>
          <el-icon class="search-icon"><Search /></el-icon>
        </template>
      </el-input>
    </div>

    <!-- 统计信息 -->
    <div class="stats-section">
      <div class="stat-item">
        <div class="stat-icon">📊</div>
        <div class="stat-content">
          <div class="stat-label">总分组数</div>
          <div class="stat-value">{{ totalGroups }}</div>
        </div>
      </div>
      <div class="stat-item">
        <div class="stat-icon">🖥️</div>
        <div class="stat-content">
          <div class="stat-label">总主机数</div>
          <div class="stat-value">{{ totalHosts }}</div>
        </div>
      </div>
    </div>

    <!-- 分组树 -->
    <div class="group-tree-container">
      <el-scrollbar class="tree-scrollbar">
        <el-tree
          ref="treeRef"
          :data="treeData"
          :props="treeProps"
          :expand-on-click-node="false"
          :highlight-current="true"
          node-key="id"
          class="group-tree"
          @node-click="handleNodeClick"
          @node-contextmenu="handleContextMenu"
        >
          <template #default="{ node, data }">
            <div class="tree-node" :class="{ 'is-selected': selectedGroupId === data.id }">
              <div class="node-content">
                <div class="node-info">
                  <el-icon class="node-icon" :class="getNodeIconClass(data)">
                    <component :is="getNodeIcon(data)" />
                  </el-icon>
                  <span class="node-label">{{ data.name }}</span>
                </div>
                <div class="node-meta">
                  <el-tag 
                    class="host-count-tag" 
                    size="small" 
                    :type="getCountTagType(data.host_count)"
                    round
                  >
                    {{ data.host_count || 0 }}
                  </el-tag>
                </div>
              </div>
              <div class="node-actions" v-if="data.id !== null">
                <el-tooltip content="编辑" placement="top">
                  <el-button 
                    class="node-action-btn" 
                    size="small" 
                    text 
                    @click.stop="handleEditGroup(data)"
                    :icon="Edit"
                  />
                </el-tooltip>
                <el-tooltip content="删除" placement="top">
                  <el-button 
                    class="node-action-btn delete-btn" 
                    size="small" 
                    text 
                    @click.stop="handleDeleteGroup(data)"
                    :icon="Delete"
                  />
                </el-tooltip>
              </div>
            </div>
          </template>
        </el-tree>
      </el-scrollbar>
    </div>

    <!-- 右键菜单 -->
    <el-dropdown
      ref="contextMenuRef"
      trigger="contextmenu"
      :teleported="false"
      class="context-menu"
    >
      <span></span>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item @click="handleAddSubGroup" :icon="Plus">
            新建子分组
          </el-dropdown-item>
          <el-dropdown-item @click="handleEditContextGroup" :icon="Edit">
            编辑分组
          </el-dropdown-item>
          <el-dropdown-item 
            @click="handleDeleteContextGroup" 
            :icon="Delete"
            class="danger-item"
          >
            删除分组
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { 
  FolderOpened, 
  Plus, 
  Refresh, 
  Search, 
  Edit, 
  Delete,
  Folder,
  House
} from '@element-plus/icons-vue'
import type { HostGroup } from '@/types/api/host'

interface Props {
  groups: HostGroup[]
  selectedGroupId: number | null
  loading?: boolean
}

interface Emits {
  (e: 'group-select', group: HostGroup): void
  (e: 'add-group'): void
  (e: 'edit-group', group: HostGroup): void
  (e: 'delete-group', group: HostGroup): void
  (e: 'refresh'): void
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
})

const emit = defineEmits<Emits>()

// 响应式数据
const searchQuery = ref('')
const treeRef = ref()
const contextMenuRef = ref()
const contextMenuGroup = ref<HostGroup | null>(null)

// 树形组件配置
const treeProps = {
  children: 'children',
  label: 'name'
}

// 计算属性
const totalGroups = computed(() => {
  const countGroups = (groups: HostGroup[]): number => {
    let count = 0
    groups.forEach(group => {
      if (group.id !== null) count++
      if (group.children) count += countGroups(group.children)
    })
    return count
  }
  return countGroups(props.groups)
})

const totalHosts = computed(() => {
  const countHosts = (groups: HostGroup[]): number => {
    let count = 0
    groups.forEach(group => {
      count += group.host_count || 0
      if (group.children) count += countHosts(group.children)
    })
    return count
  }
  return countHosts(props.groups)
})

const treeData = computed(() => {
  // 添加"全部主机"根节点
  const allHostsNode: HostGroup = {
    id: null,
    name: '全部主机',
    host_count: totalHosts.value,
    children: props.groups.filter(group => group.id !== null)
  }

  const filterTree = (nodes: HostGroup[]): HostGroup[] => {
    if (!searchQuery.value) return nodes

    return nodes.filter(node => {
      const matchesSearch = node.name.toLowerCase().includes(searchQuery.value.toLowerCase())
      const hasMatchingChildren = node.children && filterTree(node.children).length > 0
      
      if (matchesSearch || hasMatchingChildren) {
        return {
          ...node,
          children: node.children ? filterTree(node.children) : []
        }
      }
      return false
    }).filter(Boolean)
  }

  return filterTree([allHostsNode])
})

// 方法
const getNodeIcon = (data: HostGroup) => {
  if (data.id === null) return House
  return Folder
}

const getNodeIconClass = (data: HostGroup) => {
  return {
    'all-hosts-icon': data.id === null,
    'group-icon': data.id !== null
  }
}

const getCountTagType = (count: number) => {
  if (count === 0) return 'info'
  if (count < 5) return 'success'
  if (count < 20) return 'warning'
  return 'danger'
}

const handleNodeClick = (data: HostGroup) => {
  emit('group-select', data)
}

const handleAddGroup = () => {
  emit('add-group')
}

const handleEditGroup = (group: HostGroup) => {
  emit('edit-group', group)
}

const handleDeleteGroup = (group: HostGroup) => {
  emit('delete-group', group)
}

const handleRefresh = () => {
  emit('refresh')
}

const handleContextMenu = (event: MouseEvent, data: HostGroup) => {
  if (data.id === null) return // 不允许对"全部主机"右键
  
  contextMenuGroup.value = data
  nextTick(() => {
    contextMenuRef.value?.handleOpen()
  })
}

const handleAddSubGroup = () => {
  // TODO: 实现添加子分组
  emit('add-group')
}

const handleEditContextGroup = () => {
  if (contextMenuGroup.value) {
    emit('edit-group', contextMenuGroup.value)
  }
}

const handleDeleteContextGroup = () => {
  if (contextMenuGroup.value) {
    emit('delete-group', contextMenuGroup.value)
  }
}

// 监听选中状态变化
watch(() => props.selectedGroupId, (newId) => {
  if (newId !== null && treeRef.value) {
    nextTick(() => {
      treeRef.value.setCurrentKey(newId)
    })
  }
})
</script>

<style scoped>
.modern-sidebar {
  width: 280px;
  height: 100%;
  background: linear-gradient(180deg, #ffffff 0%, #fafbfc 100%);
  border-right: 1px solid #e1e5e9;
  display: flex;
  flex-direction: column;
  box-shadow: 2px 0 12px rgba(0, 0, 0, 0.08);
  position: relative;
}

/* 头部区域 */
.sidebar-header {
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  position: relative;
  overflow: hidden;
}

.sidebar-header::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, rgba(255,255,255,0.1) 0%, transparent 70%);
  animation: shimmer 3s ease-in-out infinite;
}

@keyframes shimmer {
  0%, 100% { transform: rotate(0deg); }
  50% { transform: rotate(180deg); }
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  position: relative;
  z-index: 1;
}

.title-section {
  display: flex;
  align-items: center;
  gap: 12px;
}

.icon-wrapper {
  width: 36px;
  height: 36px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  backdrop-filter: blur(10px);
}

.title-icon {
  font-size: 18px;
  color: white;
}

.title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  background: rgba(255, 255, 255, 0.15) !important;
  border: 1px solid rgba(255, 255, 255, 0.2) !important;
  color: white !important;
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
}

.action-btn:hover {
  background: rgba(255, 255, 255, 0.25) !important;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 搜索区域 */
.search-section {
  padding: 16px 20px;
  border-bottom: 1px solid #f0f2f5;
}

.search-input {
  --el-input-border-radius: 12px;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  transition: all 0.3s ease;
}

.search-input :deep(.el-input__wrapper:hover) {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.search-icon {
  color: #9ca3af;
}

/* 统计信息 */
.stats-section {
  padding: 16px 20px;
  display: flex;
  gap: 12px;
  border-bottom: 1px solid #f0f2f5;
}

.stat-item {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-radius: 12px;
  border: 1px solid #e2e8f0;
}

.stat-icon {
  font-size: 20px;
  filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.1));
}

.stat-content {
  flex: 1;
}

.stat-label {
  font-size: 11px;
  color: #6b7280;
  margin-bottom: 2px;
  font-weight: 500;
}

.stat-value {
  font-size: 16px;
  font-weight: 700;
  color: #1f2937;
}

/* 分组树容器 */
.group-tree-container {
  flex: 1;
  padding: 8px 0;
  overflow: hidden;
}

.tree-scrollbar {
  height: 100%;
}

.group-tree {
  padding: 0 12px;
}

/* 树节点样式 */
.tree-node {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 8px 12px;
  border-radius: 10px;
  transition: all 0.3s ease;
  position: relative;
  margin: 2px 0;
}

.tree-node:hover {
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  transform: translateX(2px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.tree-node.is-selected {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.3);
  transform: translateX(4px);
}

.node-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex: 1;
  min-width: 0;
}

.node-info {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  min-width: 0;
}

.node-icon {
  font-size: 16px;
  flex-shrink: 0;
}

.node-icon.all-hosts-icon {
  color: #f59e0b;
}

.node-icon.group-icon {
  color: #6366f1;
}

.tree-node.is-selected .node-icon {
  color: white;
}

.node-label {
  font-weight: 500;
  font-size: 13px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
}

.node-meta {
  margin-left: 8px;
}

.host-count-tag {
  font-weight: 600;
  font-size: 10px;
  min-width: 20px;
}

.tree-node.is-selected .host-count-tag {
  background: rgba(255, 255, 255, 0.2) !important;
  color: white !important;
  border-color: rgba(255, 255, 255, 0.3) !important;
}

.node-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: all 0.3s ease;
  transform: translateX(8px);
}

.tree-node:hover .node-actions {
  opacity: 1;
  transform: translateX(0);
}

.node-action-btn {
  padding: 4px !important;
  border-radius: 6px !important;
  background: rgba(255, 255, 255, 0.8) !important;
  border: 1px solid rgba(0, 0, 0, 0.1) !important;
  color: #6b7280 !important;
}

.node-action-btn:hover {
  background: white !important;
  color: #374151 !important;
  transform: scale(1.1) !important;
}

.node-action-btn.delete-btn:hover {
  color: #ef4444 !important;
}

/* 右键菜单 */
.context-menu {
  position: absolute;
  z-index: 1000;
}

.danger-item {
  color: #ef4444 !important;
}

/* 树组件样式覆盖 */
.group-tree :deep(.el-tree-node__content) {
  padding: 0 !important;
  height: auto !important;
  background: transparent !important;
}

.group-tree :deep(.el-tree-node__expand-icon) {
  color: #9ca3af;
  font-size: 14px;
  margin-right: 8px;
}

.group-tree :deep(.el-tree-node__expand-icon.expanded) {
  transform: rotate(90deg);
}

.group-tree :deep(.el-tree-node__expand-icon:hover) {
  color: #6366f1;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .modern-sidebar {
    width: 240px;
  }
  
  .stats-section {
    flex-direction: column;
    gap: 8px;
  }
}

@media (max-width: 992px) {
  .modern-sidebar {
    width: 100%;
    max-height: 300px;
    border-right: none;
    border-bottom: 1px solid #e1e5e9;
  }
  
  .group-tree-container {
    max-height: 200px;
  }
}
</style>
