<template>
  <el-dialog
    v-model="visible"
    title=""
    width="600px"
    :close-on-click-modal="false"
    class="batch-delete-dialog"
    :show-close="false"
  >
    <!-- 自定义头部 -->
    <template #header>
      <div class="dialog-header">
        <div class="header-left">
          <div class="header-icon">
            <el-icon size="24"><Delete /></el-icon>
          </div>
          <div class="header-content">
            <h2 class="dialog-title">批量删除主机</h2>
            <p class="dialog-subtitle">此操作将永久删除选中的主机，请谨慎操作</p>
          </div>
        </div>
        <el-button 
                          type="link" 
          size="large" 
          @click="handleCancel"
          class="close-btn"
        >
          <el-icon size="20"><Close /></el-icon>
        </el-button>
      </div>
    </template>

    <div class="delete-container">
      <!-- 警告提示 -->
      <div class="warning-section">
        <div class="warning-card">
          <div class="warning-icon">
            <el-icon size="32" color="#F56C6C"><WarningFilled /></el-icon>
          </div>
          <div class="warning-content">
            <h3>危险操作警告</h3>
            <p>您即将删除 <strong>{{ hosts.length }}</strong> 台主机，此操作不可恢复！</p>
            <ul class="warning-list">
              <li>主机的所有配置信息将被永久删除</li>
              <li>相关的监控数据和日志记录将被清除</li>
              <li>删除后无法通过系统恢复</li>
            </ul>
          </div>
        </div>
      </div>

      <!-- 主机列表 -->
      <div class="hosts-section">
        <div class="section-header">
          <h4>将要删除的主机列表</h4>
          <div class="host-count">
            <el-tag type="danger" effect="light">{{ hosts.length }} 台主机</el-tag>
          </div>
        </div>
        
        <div class="hosts-table-wrapper">
          <el-table 
            :data="hosts" 
            size="default" 
            max-height="300"
            class="hosts-table"
            :show-header="true"
          >
            <el-table-column type="index" label="#" width="50" />
            <el-table-column prop="name" label="主机名称" min-width="150">
              <template #default="{ row }">
                <div class="host-name">
                  <el-icon class="host-icon"><Monitor /></el-icon>
                  <span>{{ row.name }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="public_ip" label="公网IP" width="140">
              <template #default="{ row }">
                <span class="ip-address">
                  {{ Array.isArray(row.public_ip) ? row.public_ip[0] : row.public_ip }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag 
                  :type="getStatusType(row.status)" 
                  size="small"
                  effect="light"
                >
                  {{ getStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="provider.name" label="云厂商" width="120">
              <template #default="{ row }">
                <div class="provider-info">
                  <span>{{ row.provider?.name || '-' }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="group.name" label="主机组" width="120">
              <template #default="{ row }">
                <span class="group-name">{{ row.group?.name || '-' }}</span>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>

      <!-- 确认输入 -->
      <div class="confirm-section">
        <div class="confirm-card">
          <div class="confirm-header">
            <el-icon size="20" color="#E6A23C"><Key /></el-icon>
            <span>安全确认</span>
          </div>
          <div class="confirm-content">
            <p>为了防止误操作，请在下方输入框中输入 <code>DELETE</code> 来确认删除操作：</p>
            <el-input
              v-model="confirmText"
              placeholder="请输入 DELETE 确认删除"
              size="large"
              class="confirm-input"
              :class="{ 'is-valid': isConfirmValid }"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- 底部操作栏 -->
    <template #footer>
      <div class="dialog-footer">
        <div class="footer-left">
          <div class="delete-info">
            <el-icon><InfoFilled /></el-icon>
            <span>删除操作不可恢复，请谨慎操作</span>
          </div>
        </div>
        <div class="footer-right">
          <el-button @click="handleCancel" size="large">
            取消
          </el-button>
          <el-button 
            type="danger" 
            @click="handleConfirm" 
            :loading="deleting"
            :disabled="!isConfirmValid"
            size="large"
            :icon="Delete"
          >
            {{ deleting ? '删除中...' : '确认删除' }}
          </el-button>
        </div>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Delete, 
  Close, 
  WarningFilled, 
  Monitor, 
  Key, 
  InfoFilled 
} from '@element-plus/icons-vue'
import { batchDeleteHosts } from '@/api/system/host'
import type { Host } from '@/types/api/host'

interface Props {
  modelValue: boolean
  hosts: Host[]
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
  (e: 'cancel'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const deleting = ref(false)
const confirmText = ref('')

// 确认文本验证
const isConfirmValid = computed(() => {
  return confirmText.value.trim().toUpperCase() === 'DELETE'
})

// 状态相关方法
const getStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    running: 'success',
    stopped: 'warning',
    error: 'danger'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    running: '运行中',
    stopped: '已停止',
    error: '异常'
  }
  return statusMap[status] || status
}

// 确认删除
const handleConfirm = async () => {
  if (!isConfirmValid.value) {
    ElMessage.error('请输入 DELETE 确认删除操作')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除这 ${props.hosts.length} 台主机吗？此操作不可恢复！`,
      '最终确认',
      { 
        type: 'error',
        confirmButtonText: '确认删除',
        cancelButtonText: '取消',
        confirmButtonClass: 'el-button--danger'
      }
    )

    deleting.value = true
    const hostIds = props.hosts.map(host => host.id)
    await batchDeleteHosts(hostIds)
    
    ElMessage.success(`成功删除 ${props.hosts.length} 台主机`)
    emit('success')
    handleCancel()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败：' + (error as Error).message)
    }
  } finally {
    deleting.value = false
  }
}

// 取消操作
const handleCancel = () => {
  visible.value = false
  confirmText.value = ''
  deleting.value = false
  emit('cancel')
}
</script>

<style scoped>
.batch-delete-dialog {
  --el-dialog-border-radius: 12px;
}

.batch-delete-dialog :deep(.el-dialog__header) {
  padding: 0;
  margin: 0;
}

.batch-delete-dialog :deep(.el-dialog__body) {
  padding: 0 24px 24px;
}

.batch-delete-dialog :deep(.el-dialog__footer) {
  padding: 0 24px 24px;
}

/* 对话框头部 */
.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px 24px 0;
  background: linear-gradient(135deg, #ff6b6b 0%, #ee5a24 100%);
  color: white;
  border-radius: 12px 12px 0 0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-icon {
  width: 48px;
  height: 48px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.dialog-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}

.dialog-subtitle {
  margin: 4px 0 0;
  font-size: 14px;
  opacity: 0.9;
}

.close-btn {
  color: white !important;
  background: rgba(255, 255, 255, 0.1) !important;
  border: none !important;
  border-radius: 8px !important;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.2) !important;
}

/* 删除容器 */
.delete-container {
  padding: 32px 0;
}

/* 警告区域 */
.warning-section {
  margin-bottom: 32px;
}

.warning-card {
  display: flex;
  gap: 16px;
  padding: 24px;
  background: #fef0f0;
  border: 1px solid #fbc4c4;
  border-radius: 12px;
}

.warning-icon {
  flex-shrink: 0;
}

.warning-content h3 {
  margin: 0 0 12px;
  font-size: 16px;
  font-weight: 600;
  color: #f56c6c;
}

.warning-content p {
  margin: 0 0 16px;
  color: #606266;
  line-height: 1.5;
}

.warning-content strong {
  color: #f56c6c;
  font-weight: 600;
}

.warning-list {
  margin: 0;
  padding-left: 20px;
  color: #909399;
}

.warning-list li {
  margin-bottom: 8px;
  line-height: 1.4;
}

/* 主机区域 */
.hosts-section {
  margin-bottom: 32px;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.section-header h4 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.hosts-table-wrapper {
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #ebeef5;
}

.hosts-table :deep(.el-table__header) {
  background: #f8f9fa;
}

.hosts-table :deep(.el-table__header th) {
  background: #f8f9fa;
  color: #303133;
  font-weight: 600;
}

.host-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.host-icon {
  color: #909399;
}

.ip-address {
  font-family: monospace;
  color: #606266;
}

.provider-info,
.group-name {
  color: #606266;
}

/* 确认区域 */
.confirm-section {
  margin-bottom: 24px;
}

.confirm-card {
  padding: 24px;
  background: #fff7e6;
  border: 1px solid #ffd591;
  border-radius: 12px;
}

.confirm-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  font-size: 16px;
  font-weight: 600;
  color: #e6a23c;
}

.confirm-content p {
  margin: 0 0 16px;
  color: #606266;
  line-height: 1.5;
}

.confirm-content code {
  background: #f4f4f5;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: monospace;
  color: #e6a23c;
  font-weight: 600;
}

.confirm-input {
  width: 100%;
}

.confirm-input.is-valid :deep(.el-input__wrapper) {
  border-color: #67c23a;
  box-shadow: 0 0 0 1px rgba(103, 194, 58, 0.2);
}

/* 底部操作栏 */
.dialog-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 24px;
  border-top: 1px solid #ebeef5;
}

.footer-left {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #f56c6c;
  font-size: 14px;
}

.footer-right {
  display: flex;
  gap: 12px;
}

.delete-info {
  display: flex;
  align-items: center;
  gap: 6px;
}
</style>
