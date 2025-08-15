<template>
  <div class="domain-group-sidebar">
    <!-- 头部区域 -->
    <div class="sidebar-header">
      <div class="header-content">
        <div class="title-section">
          <div class="icon-wrapper">
            <el-icon class="title-icon"><FolderOpened /></el-icon>
          </div>
          <h3 class="title">域名分组</h3>
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
        <div class="stat-icon">🌐</div>
        <div class="stat-content">
          <div class="stat-label">总域名数</div>
          <div class="stat-value">{{ totalDomains }}</div>
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
          :draggable="false"
          :allow-drop="allowDrop"
          :allow-drag="allowDrag"
          :default-expand-all="false"
          :default-expanded-keys="defaultExpandedKeys"
          node-key="id"
          class="group-tree"
          @node-click="handleNodeClick"
          @node-contextmenu="handleContextMenu"
          @node-drop="handleNodeDrop"
        >
          <template #default="{ node, data }">
            <div class="tree-node" :class="{ 'is-selected': selectedGroupId === data.id }">
              <div class="node-content">
                <div class="node-info">
                  <el-icon class="node-icon" :class="getNodeIconClass(data)">
                    <component :is="getNodeIcon(data)" />
                  </el-icon>
                  <span class="node-label" :title="data.name">{{ data.name }}</span>
                </div>
                <div class="node-meta">
                  <div class="domain-count-display">
                    <!-- 显示直接域名数量 -->
                    <el-tag 
                      class="domain-count-tag direct-count" 
                      size="small" 
                      :type="getCountTagType(data.domain_count || 0)"
                      round
                      :title="`直接归属域名: ${data.domain_count || 0}`"
                    >
                      {{ data.domain_count || 0 }}
                    </el-tag>
                    <!-- 如果有子分组且总数不等于直接数量，显示总数 -->
                    <el-tag 
                      v-if="data.children && data.children.length > 0 && (data.total_domain_count || 0) !== (data.domain_count || 0)"
                      class="domain-count-tag total-count" 
                      size="small" 
                      type="info"
                      round
                      :title="`包含子分组总数: ${data.total_domain_count || 0}`"
                    >
                      总{{ data.total_domain_count || 0 }}
                    </el-tag>
                  </div>
                </div>
              </div>
              <div class="node-actions" v-if="data.id !== null">
                <el-tooltip content="新建子分组" placement="top">
                  <el-button 
                    class="node-action-btn" 
                    size="small" 
                    text 
                    @click.stop="handleAddSubGroup(data)"
                    :icon="Plus"
                  />
                </el-tooltip>
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
import type { DomainGroup } from '@/types/dns'

interface Props {
  groups: DomainGroup[]
  selectedGroupId: number | null
  loading?: boolean
  totalDomainCount?: number
}

interface Emits {
  (e: 'group-select', group: DomainGroup): void
  (e: 'add-group'): void
  (e: 'edit-group', group: DomainGroup): void
  (e: 'delete-group', group: DomainGroup): void
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
const contextMenuGroup = ref<DomainGroup | null>(null)

// 树形组件配置
const treeProps = {
  children: 'children',
  label: 'name'
}

// 计算属性
const totalGroups = computed(() => {
  const countGroups = (groups: DomainGroup[]): number => {
    let count = 0
    groups.forEach(group => {
      if (group.id !== null) count++
      if (group.children) count += countGroups(group.children)
    })
    return count
  }
  return countGroups(props.groups)
})

const totalDomains = computed(() => {
  if (props.totalDomainCount !== undefined) {
    return props.totalDomainCount
  }

  let count = 0
  props.groups.forEach(group => {
    if (group.id !== null) {
      count += group.total_domain_count || group.domain_count || 0
    }
  })
  return count
})

const treeData = computed(() => {
  // 添加"全部域名"根节点
  const allDomainsNode: DomainGroup = {
    id: null,
    name: '全部域名',
    domain_count: totalDomains.value,
    total_domain_count: totalDomains.value,
    children: props.groups.filter(group => group.id !== null)
  }

  const filterTree = (nodes: DomainGroup[]): DomainGroup[] => {
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

  return filterTree([allDomainsNode])
})

// 方法
const getNodeIcon = (data: DomainGroup) => {
  if (data.id === null) return House
  return Folder
}

const getNodeIconClass = (data: DomainGroup) => {
  return {
    'all-domains-icon': data.id === null,
    'group-icon': data.id !== null
  }
}

const getCountTagType = (count: number) => {
  if (count === 0) return 'info'
  if (count < 5) return 'success'
  if (count < 20) return 'warning'
  return 'danger'
}

const handleNodeClick = (data: DomainGroup) => {
  emit('group-select', data)
}

const handleAddGroup = () => {
  emit('add-group')
}

const handleEditGroup = (group: DomainGroup) => {
  emit('edit-group', group)
}

const handleDeleteGroup = (group: DomainGroup) => {
  emit('delete-group', group)
}

const handleRefresh = () => {
  emit('refresh')
}

const handleContextMenu = (event: MouseEvent, data: DomainGroup) => {
  if (data.id === null) return // 不允许对"全部域名"右键
  
  contextMenuGroup.value = data
  nextTick(() => {
    contextMenuRef.value?.handleOpen()
  })
}

const handleAddSubGroup = (group: DomainGroup) => {
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

// 拖拽相关方法
const allowDrag = (draggingNode: any) => {
  return draggingNode.data.id !== null
}

const allowDrop = (draggingNode: any, dropNode: any, type: string) => {
  if (dropNode.data.id === null && type === 'inner') {
    return false
  }
  if (type === 'inner') {
    return !isDescendant(draggingNode.data, dropNode.data)
  }
  return true
}

const isDescendant = (ancestor: DomainGroup, node: DomainGroup): boolean => {
  if (!ancestor.children) return false
  return ancestor.children.some(child => 
    child.id === node.id || isDescendant(child, node)
  )
}

const handleNodeDrop = (draggingNode: any, dropNode: any, dropType: string) => {
  console.log('Node dropped:', { draggingNode, dropNode, dropType })
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
.domain-group-sidebar {
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
  padding: 10px 12px;
  border-radius: 12px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  margin: 2px 0;
  border: 1px solid transparent;
}

.tree-node::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background: transparent;
  border-radius: 0 2px 2px 0;
  transition: all 0.3s ease;
}

.tree-node:hover {
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  transform: translateX(3px);
  box-shadow: 
    0 4px 12px rgba(0, 0, 0, 0.08),
    inset 0 1px 0 rgba(255, 255, 255, 0.6);
  border-color: #e2e8f0;
}

.tree-node:hover::before {
  background: linear-gradient(180deg, #667eea 0%, #764ba2 100%);
}

.tree-node.is-selected {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 
    0 6px 20px rgba(102, 126, 234, 0.3),
    0 2px 8px rgba(102, 126, 234, 0.2),
    inset 0 1px 0 rgba(255, 255, 255, 0.2);
  transform: translateX(5px);
  border-color: rgba(255, 255, 255, 0.2);
}

.tree-node.is-selected::before {
  background: rgba(255, 255, 255, 0.8);
  width: 4px;
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
  gap: 10px;
  flex: 1;
  min-width: 0;
}

.node-icon {
  font-size: 16px;
  flex-shrink: 0;
  transition: all 0.3s ease;
}

.node-icon.all-domains-icon {
  color: #f59e0b;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.node-icon.group-icon {
  color: #6366f1;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.tree-node.is-selected .node-icon {
  color: white;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
  transform: scale(1.1);
}

.node-label {
  font-weight: 500;
  font-size: 13px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
  transition: all 0.3s ease;
}

.tree-node:hover .node-label {
  font-weight: 600;
}

.tree-node.is-selected .node-label {
  font-weight: 600;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.node-meta {
  margin-left: 8px;
  flex-shrink: 0;
}

.domain-count-display {
  display: flex;
  gap: 4px;
  align-items: center;
}

.domain-count-tag {
  font-weight: 600;
  font-size: 10px;
  min-width: 22px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.domain-count-tag.direct-count {
  background: linear-gradient(135deg, #ffffff 0%, #f8fafc 100%);
  border: 1px solid #d1d5db;
  color: #374151;
}

.domain-count-tag.total-count {
  background: linear-gradient(135deg, #eff6ff 0%, #dbeafe 100%);
  border: 1px solid #93c5fd;
  color: #1e40af;
  font-size: 9px;
}

.tree-node:hover .domain-count-tag {
  transform: scale(1.05);
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.15);
}

.tree-node.is-selected .domain-count-tag {
  background: rgba(255, 255, 255, 0.25) !important;
  color: white !important;
  border-color: rgba(255, 255, 255, 0.4) !important;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.node-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  transform: translateX(12px);
}

.tree-node:hover .node-actions {
  opacity: 1;
  transform: translateX(0);
}

.tree-node.is-selected .node-actions {
  opacity: 1;
  transform: translateX(0);
}

.node-action-btn {
  padding: 6px !important;
  border-radius: 8px !important;
  background: rgba(255, 255, 255, 0.9) !important;
  border: 1px solid rgba(0, 0, 0, 0.08) !important;
  color: #6b7280 !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
  backdrop-filter: blur(8px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.node-action-btn:hover {
  background: white !important;
  color: #374151 !important;
  transform: scale(1.15) translateY(-1px) !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15) !important;
}

.node-action-btn.delete-btn:hover {
  color: #ef4444 !important;
  background: #fef2f2 !important;
  border-color: #fecaca !important;
}

.tree-node.is-selected .node-action-btn {
  background: rgba(255, 255, 255, 0.2) !important;
  border-color: rgba(255, 255, 255, 0.3) !important;
  color: rgba(255, 255, 255, 0.9) !important;
}

.tree-node.is-selected .node-action-btn:hover {
  background: rgba(255, 255, 255, 0.3) !important;
  color: white !important;
  transform: scale(1.15) translateY(-1px) !important;
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
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border-radius: 4px;
  padding: 2px;
  background: rgba(255, 255, 255, 0.6);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.group-tree :deep(.el-tree-node__expand-icon.expanded) {
  transform: rotate(90deg);
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
}

.group-tree :deep(.el-tree-node__expand-icon:hover) {
  color: #667eea;
  background: rgba(102, 126, 234, 0.15);
  transform: scale(1.1);
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.15);
}

/* 选中状态下的展开图标 */
.group-tree :deep(.el-tree-node.is-current > .el-tree-node__content .el-tree-node__expand-icon) {
  color: rgba(255, 255, 255, 0.9);
  background: rgba(255, 255, 255, 0.2);
}

.group-tree :deep(.el-tree-node.is-current > .el-tree-node__content .el-tree-node__expand-icon:hover) {
  color: white;
  background: rgba(255, 255, 255, 0.3);
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .domain-group-sidebar {
    width: 240px;
  }
  
  .stats-section {
    flex-direction: column;
    gap: 8px;
  }
}

@media (max-width: 992px) {
  .domain-group-sidebar {
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
