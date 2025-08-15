<template>
  <div class="dns-group-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="page-header-content">
        <div class="page-title">
          <h1>域名分组管理</h1>
          <p>管理域名分组，支持层级结构和拖拽排序</p>
        </div>
        <div class="page-actions">
          <el-button type="primary" @click="handleAdd" class="page-action-btn primary-btn">
            <el-icon><Plus /></el-icon>
            <span>新建分组</span>
          </el-button>
        </div>
      </div>
    </div>

    <!-- 分组树形结构 -->
    <el-card class="tree-card">
      <template #header>
        <div class="card-header">
          <h3>分组结构</h3>
          <div class="header-actions">
            <el-button @click="expandAll" size="small">
              <el-icon><Expand /></el-icon>
              展开全部
            </el-button>
            <el-button @click="collapseAll" size="small">
              <el-icon><Fold /></el-icon>
              折叠全部
            </el-button>
          </div>
        </div>
      </template>

      <el-tree
        ref="treeRef"
        :data="groupTree"
        :props="treeProps"
        node-key="id"
        :expand-on-click-node="false"
        :default-expand-all="false"
        draggable
        @node-drop="handleNodeDrop"
        class="group-tree"
      >
        <template #default="{ node, data }">
          <div class="tree-node">
            <div class="node-content">
              <el-icon class="node-icon">
                <Folder />
              </el-icon>
              <span class="node-label">{{ data.name }}</span>
              <el-tag v-if="data.domain_count > 0" size="small" type="info">
                {{ data.domain_count }}个域名
              </el-tag>
            </div>
            <div class="node-actions">
              <el-button 
                size="small" 
                type="primary" 
                text 
                @click="handleEdit(data)"
              >
                编辑
              </el-button>
              <el-button 
                size="small" 
                type="success" 
                text 
                @click="handleAddChild(data)"
              >
                添加子分组
              </el-button>
              <el-button 
                size="small" 
                type="danger" 
                text 
                @click="handleDelete(data)"
                :disabled="data.domain_count > 0"
              >
                删除
              </el-button>
            </div>
          </div>
        </template>
      </el-tree>
    </el-card>

    <!-- 分组表单弹窗 -->
    <GroupModal
      v-model:visible="modalVisible"
      :group="currentGroup"
      :parent-options="parentOptions"
      @success="handleModalSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Plus,
  Folder,
  Expand,
  Fold
} from '@element-plus/icons-vue'
import GroupModal from './components/GroupModal.vue'
import { domainGroupApi } from '@/api/dns/group'
import type { DomainGroup } from '@/types/dns'

// 响应式数据
const treeRef = ref()
const modalVisible = ref(false)
const currentGroup = ref<DomainGroup | null>(null)
const groupTree = ref<DomainGroup[]>([])

// 树形组件配置
const treeProps = {
  children: 'children',
  label: 'name'
}

// 计算父级选项（用于弹窗）
const parentOptions = computed(() => {
  const options: any[] = [{ id: null, name: '根级分组' }]
  
  const flattenOptions = (groups: DomainGroup[], level = 0) => {
    groups.forEach(group => {
      options.push({
        id: group.id,
        name: '  '.repeat(level) + group.name
      })
      if (group.children && group.children.length > 0) {
        flattenOptions(group.children, level + 1)
      }
    })
  }
  
  flattenOptions(groupTree.value)
  return options
})

// 事件处理
const handleAdd = () => {
  currentGroup.value = null
  modalVisible.value = true
}

const handleEdit = (group: DomainGroup) => {
  currentGroup.value = { ...group }
  modalVisible.value = true
}

const handleAddChild = (parent: DomainGroup) => {
  currentGroup.value = {
    id: 0,
    name: '',
    description: '',
    parent_id: parent.id,
    sort: 0,
    domain_count: 0,
    children: []
  } as DomainGroup
  modalVisible.value = true
}

const handleDelete = async (group: DomainGroup) => {
  if (group.domain_count > 0) {
    ElMessage.warning('该分组下还有域名，无法删除')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除分组 "${group.name}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await domainGroupApi.delete(group.id)
    ElMessage.success('删除成功')
    await fetchGroupTree()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleNodeDrop = async (draggingNode: any, dropNode: any, dropType: string) => {
  try {
    const dragData = draggingNode.data
    const dropData = dropNode.data
    
    let newParentId = null
    let newSort = 0
    
    if (dropType === 'inner') {
      // 拖拽到节点内部，成为子节点
      newParentId = dropData.id
    } else {
      // 拖拽到节点前后，成为同级节点
      newParentId = dropData.parent_id
      newSort = dropData.sort + (dropType === 'after' ? 1 : 0)
    }
    
    await domainGroupApi.update(dragData.id, {
      parent_id: newParentId,
      sort: newSort
    })
    
    ElMessage.success('分组移动成功')
    await fetchGroupTree()
  } catch (error) {
    ElMessage.error('分组移动失败')
    await fetchGroupTree() // 重新加载数据恢复原状
  }
}

const expandAll = () => {
  const nodes = treeRef.value?.store._getAllNodes()
  nodes?.forEach((node: any) => {
    node.expanded = true
  })
}

const collapseAll = () => {
  const nodes = treeRef.value?.store._getAllNodes()
  nodes?.forEach((node: any) => {
    node.expanded = false
  })
}

const handleModalSuccess = () => {
  modalVisible.value = false
  fetchGroupTree()
}

// 数据获取
const fetchGroupTree = async () => {
  try {
    const response = await domainGroupApi.tree()
    groupTree.value = response || []
  } catch (error) {
    ElMessage.error('获取分组树失败')
  }
}

// 生命周期
onMounted(() => {
  fetchGroupTree()
})
</script>

<style scoped lang="scss">
.dns-group-container {
  padding: 24px;
  background: linear-gradient(135deg, #f0f8ff 0%, #e6f3ff 50%, #f8f9fa 100%);
  min-height: calc(100vh - 64px);
  position: relative;
}

.dns-group-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 200px;
  background: linear-gradient(135deg, rgba(24, 144, 255, 0.05) 0%, rgba(64, 169, 255, 0.03) 100%);
  pointer-events: none;
  z-index: 0;
}

.dns-group-container > * {
  position: relative;
  z-index: 1;
}

.page-header {
  margin-bottom: 24px;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.95) 0%, rgba(255, 255, 255, 0.9) 100%);
  padding: 32px;
  border-radius: 20px;
  box-shadow: 
    0 8px 32px rgba(24, 144, 255, 0.12),
    0 2px 8px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(24, 144, 255, 0.1);
  backdrop-filter: blur(20px);
}

.page-header-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.page-title h1 {
  margin: 0 0 12px 0;
  font-size: 32px;
  font-weight: 700;
  background: linear-gradient(135deg, #1890ff 0%, #722ed1 50%, #eb2f96 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: -0.5px;
}

.page-title p {
  margin: 0;
  color: #64748b;
  font-size: 15px;
}

.page-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.page-action-btn {
  border-radius: 12px;
  font-weight: 600;
  font-size: 14px;
  padding: 12px 20px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  gap: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border: none;
}

.page-action-btn:hover {
  transform: translateY(-2px) scale(1.02);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.15);
}

.page-action-btn.primary-btn {
  background: linear-gradient(135deg, #1890ff, #40a9ff);
  color: white;
}

.page-action-btn.primary-btn:hover {
  background: linear-gradient(135deg, #096dd9, #1890ff);
  box-shadow: 0 8px 30px rgba(24, 144, 255, 0.4);
}

.tree-card {
  border-radius: 16px;
  border: 1px solid rgba(24, 144, 255, 0.1);
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.9) 0%, rgba(255, 255, 255, 0.95) 100%);
  backdrop-filter: blur(10px);
  box-shadow: 0 6px 30px rgba(24, 144, 255, 0.12);
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.header-actions {
  display: flex;
  gap: 8px;
}

.group-tree {
  :deep(.el-tree-node) {
    .el-tree-node__content {
      height: 48px;
      border-radius: 8px;
      margin: 2px 0;
      transition: all 0.3s ease;
      
      &:hover {
        background: linear-gradient(135deg, rgba(24, 144, 255, 0.08) 0%, rgba(64, 169, 255, 0.05) 100%);
      }
    }
    
    .el-tree-node__expand-icon {
      color: var(--el-color-primary);
      font-size: 16px;
    }
  }
}

.tree-node {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 0 12px 0 0;
}

.node-content {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.node-icon {
  color: var(--el-color-primary);
  font-size: 16px;
}

.node-label {
  font-weight: 500;
  color: var(--el-text-color-primary);
}

.node-actions {
  display: flex;
  gap: 8px;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.tree-node:hover .node-actions {
  opacity: 1;
}
</style>