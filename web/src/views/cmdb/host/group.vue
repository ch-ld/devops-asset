<template>
  <div class="host-group-management">
    <a-card :bordered="false">
      <template #title>
        <div class="header">
          <span>主机组管理</span>
          <a-button type="primary" @click="openGroupModal()">
            <template #icon><PlusOutlined /></template>
            添加主机组
          </a-button>
        </div>
      </template>

      <div class="content">
        <div class="tree-container">
          <a-tree
            v-model:expandedKeys="expandedKeys"
            :tree-data="treeData"
            :load-data="loadData"
            draggable
            block-node
            @drop="onDrop"
            @select="onSelect"
          >
            <template #title="{ title, key, dataRef }">
              <div class="tree-node">
                <span class="node-title">{{ title }}</span>
                <div class="node-actions">
                  <a-tooltip title="添加子组">
                    <a-button
                      type="text"
                      size="small"
                      @click.stop="openGroupModal(dataRef)"
                    >
                      <template #icon><PlusOutlined /></template>
                    </a-button>
                  </a-tooltip>
                  <a-tooltip title="编辑">
                    <a-button
                      type="text"
                      size="small"
                      @click.stop="editGroup(dataRef)"
                    >
                      <template #icon><EditOutlined /></template>
                    </a-button>
                  </a-tooltip>
                  <a-tooltip title="删除">
                    <a-button
                      type="text"
                      size="small"
                      danger
                      @click.stop="deleteGroup(dataRef)"
                    >
                      <template #icon><DeleteOutlined /></template>
                    </a-button>
                  </a-tooltip>
                </div>
              </div>
            </template>
          </a-tree>
        </div>

        <div class="detail-container" v-if="selectedGroup">
          <a-card size="small" title="主机组详情">
            <a-descriptions :column="1" size="small">
              <a-descriptions-item label="名称">
                {{ selectedGroup.name }}
              </a-descriptions-item>
              <a-descriptions-item label="描述">
                {{ selectedGroup.description || '-' }}
              </a-descriptions-item>
              <a-descriptions-item label="路径">
                {{ selectedGroup.path }}
              </a-descriptions-item>
              <a-descriptions-item label="排序">
                {{ selectedGroup.sort }}
              </a-descriptions-item>
              <a-descriptions-item label="创建时间">
                {{ formatDate(selectedGroup.created_at) }}
              </a-descriptions-item>
            </a-descriptions>
          </a-card>

          <a-card size="small" title="主机列表" style="margin-top: 16px">
            <a-table
              :columns="hostColumns"
              :data-source="groupHosts"
              :loading="hostLoading"
              :pagination="hostPagination"
              size="small"
              row-key="id"
              @change="handleHostTableChange"
            >
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'status'">
                  <a-tag :color="getStatusColor(record.status)">
                    {{ getStatusText(record.status) }}
                  </a-tag>
                </template>
                <template v-if="column.key === 'action'">
                  <a-button
                    type="link"
                    size="small"
                    @click="removeHostFromGroup(record.id)"
                  >
                    移出分组
                  </a-button>
                </template>
              </template>
            </a-table>
          </a-card>
        </div>
      </div>
    </a-card>

    <!-- 主机组表单弹窗 -->
    <HostGroupModal
      ref="groupModalRef"
      :host-groups="hostGroupTree"
      @success="handleModalSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { useHostStore } from '@/store/modules/host'
import { storeToRefs } from 'pinia'
import type { HostGroup, Host } from '@/types/api/host'
import {
  PlusOutlined,
  EditOutlined,
  DeleteOutlined
} from '@ant-design/icons-vue'
import HostGroupModal from './components/HostGroupModal.vue'
import * as hostApi from '@/api/system/host'

const hostStore = useHostStore()
const { hostGroupTree, isLoading } = storeToRefs(hostStore)

const expandedKeys = ref<string[]>([])
const selectedGroup = ref<HostGroup | null>(null)
const groupHosts = ref<Host[]>([])
const hostLoading = ref(false)
const hostPagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const groupModalRef = ref()

const treeData = computed(() => {
  const buildTreeData = (groups: HostGroup[], parentId: number | null = null): any[] => {
    return groups
      .filter(group => group.parent_id === parentId)
      .map(group => ({
        title: group.name,
        key: group.id.toString(),
        dataRef: group,
        children: buildTreeData(groups, group.id)
      }))
  }
  return buildTreeData(hostGroupTree.value)
})

const hostColumns = [
  {
    title: '主机名',
    dataIndex: 'name',
    key: 'name',
    width: 200
  },
  {
    title: '实例ID',
    dataIndex: 'instance_id',
    key: 'instance_id',
    width: 150
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    width: 100
  },
  {
    title: '公网IP',
    dataIndex: 'public_ip',
    key: 'public_ip',
    width: 150,
    customRender: ({ record }: { record: Host }) => {
      return record.public_ip && record.public_ip.length > 0 ? record.public_ip[0] : '-'
    }
  },
  {
    title: '操作',
    key: 'action',
    width: 100
  }
]

const getStatusColor = (status: string) => {
  const colorMap: Record<string, string> = {
    running: 'green',
    stopped: 'red',
    starting: 'blue',
    stopping: 'orange',
    unknown: 'default'
  }
  return colorMap[status] || 'default'
}

const getStatusText = (status: string) => {
  const textMap: Record<string, string> = {
    running: '运行中',
    stopped: '已停止',
    starting: '启动中',
    stopping: '停止中',
    unknown: '未知'
  }
  return textMap[status] || status
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleString()
}

const loadData = (treeNode: any) => {
  // 这里可以实现懒加载逻辑
  return Promise.resolve()
}

const onDrop = async (info: any) => {
  const dragNode = info.dragNode
  const dropNode = info.node
  const dropPosition = info.dropPosition
  
  try {
    // 确定新的父节点ID
    let newParentId: number | undefined
    
    if (dropPosition === 0) {
      // 拖拽到节点内部，作为子节点
      newParentId = parseInt(dropNode.key)
    } else {
      // 拖拽到节点前后，与该节点同级
      newParentId = dropNode.dataRef.parent_id
    }
    
    await hostStore.moveHostGroup(parseInt(dragNode.key), newParentId)
    message.success('移动成功')
  } catch (error) {
    message.error('移动失败')
  }
}

const onSelect = async (selectedKeys: string[]) => {
  if (selectedKeys.length === 0) {
    selectedGroup.value = null
    groupHosts.value = []
    return
  }
  
  const groupId = parseInt(selectedKeys[0])
  const group = findGroupById(hostGroupTree.value, groupId)
  
  if (group) {
    selectedGroup.value = group
    await loadGroupHosts(groupId)
  }
}

const findGroupById = (groups: HostGroup[], id: number): HostGroup | null => {
  for (const group of groups) {
    if (group.id === id) {
      return group
    }
    const found = findGroupById(groups, id)
    if (found) return found
  }
  return null
}

const loadGroupHosts = async (groupId: number) => {
  hostLoading.value = true
  try {
    const { data } = await hostApi.getGroupHosts(
      groupId,
      hostPagination.current,
      hostPagination.pageSize
    )
    groupHosts.value = data.data || []
    hostPagination.total = data.total || 0
  } catch (error) {
    message.error('加载主机列表失败')
  } finally {
    hostLoading.value = false
  }
}

const handleHostTableChange = (pagination: any) => {
  hostPagination.current = pagination.current
  hostPagination.pageSize = pagination.pageSize
  
  if (selectedGroup.value) {
    loadGroupHosts(selectedGroup.value.id)
  }
}

const openGroupModal = (parentGroup?: HostGroup) => {
  groupModalRef.value?.open(null, parentGroup)
}

const editGroup = (group: HostGroup) => {
  groupModalRef.value?.open(group)
}

const deleteGroup = (group: HostGroup) => {
  Modal.confirm({
    title: '确认删除',
    content: `确定要删除主机组"${group.name}"吗？`,
    onOk: async () => {
      try {
        await hostStore.deleteHostGroup(group.id)
        message.success('删除成功')
        if (selectedGroup.value?.id === group.id) {
          selectedGroup.value = null
          groupHosts.value = []
        }
      } catch (error) {
        message.error('删除失败')
      }
    }
  })
}

const removeHostFromGroup = async (hostId: number) => {
  try {
    await hostStore.moveHost(hostId, undefined)
    message.success('移出成功')
    if (selectedGroup.value) {
      loadGroupHosts(selectedGroup.value.id)
    }
  } catch (error) {
    message.error('移出失败')
  }
}

const handleModalSuccess = () => {
  // 刷新主机组树
  hostStore.fetchHostGroupTree()
}

onMounted(() => {
  hostStore.fetchHostGroupTree()
})
</script>

<style scoped lang="scss">
.host-group-management {
  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .content {
    display: flex;
    gap: 16px;
    height: 600px;

    .tree-container {
      flex: 1;
      overflow: auto;
      border: 1px solid #f0f0f0;
      border-radius: 6px;
      padding: 16px;

      .tree-node {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 4px 0;

        .node-title {
          flex: 1;
        }

        .node-actions {
          display: none;
          gap: 4px;
        }

        &:hover .node-actions {
          display: flex;
        }
      }
    }

    .detail-container {
      flex: 1;
      overflow: auto;
    }
  }
}
</style> 