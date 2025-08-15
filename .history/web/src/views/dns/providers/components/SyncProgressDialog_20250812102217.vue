<template>
  <el-dialog
    v-model="visible"
    title="域名同步进度"
    width="600px"
    :close-on-click-modal="false"
    :show-close="!syncing"
  >
    <div class="sync-progress">
      <!-- 总体进度 -->
      <div class="progress-header">
        <h4>同步进度：{{ currentIndex }}/{{ totalProviders }}</h4>
        <el-progress 
          :percentage="progressPercent" 
          :status="syncing ? 'success' : 'success'"
        />
      </div>

      <!-- 当前正在同步的提供商 -->
      <div v-if="currentProvider" class="current-sync">
        <el-alert
          :title="`正在同步：${currentProvider.name} (${currentProvider.type})`"
          type="info"
          :closable="false"
          show-icon
        />
      </div>

      <!-- 同步结果列表 -->
      <div class="sync-results">
        <el-table :data="results" style="width: 100%" max-height="300">
          <el-table-column prop="name" label="提供商" width="150" />
          <el-table-column prop="type" label="类型" width="100" />
          <el-table-column prop="status" label="状态" width="80">
            <template #default="{ row }">
              <el-tag 
                :type="row.status === 'success' ? 'success' : row.status === 'error' ? 'danger' : 'info'"
                size="small"
              >
                {{ row.status === 'success' ? '成功' : row.status === 'error' ? '失败' : '同步中' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="domains_count" label="域名数" width="80" />
          <el-table-column prop="message" label="说明" />
        </el-table>
      </div>

      <!-- 汇总统计 -->
      <div v-if="!syncing" class="sync-summary">
        <el-divider />
        <el-row :gutter="20">
          <el-col :span="6">
            <div class="summary-item">
              <div class="summary-value">{{ totalDomains }}</div>
              <div class="summary-label">总域名数</div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="summary-item">
              <div class="summary-value">{{ successCount }}</div>
              <div class="summary-label">成功</div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="summary-item">
              <div class="summary-value">{{ errorCount }}</div>
              <div class="summary-label">失败</div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="summary-item">
              <div class="summary-value">{{ Math.round((successCount / totalProviders) * 100) }}%</div>
              <div class="summary-label">成功率</div>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>

    <template #footer>
      <el-button v-if="!syncing" @click="handleClose">关闭</el-button>
      <el-button v-if="syncing" @click="handleCancel" type="danger">取消同步</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'

// Props
interface Props {
  visible: boolean
  providers: any[]
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  providers: () => []
})

// Emits
const emit = defineEmits<{
  'update:visible': [value: boolean]
  cancel: []
}>()

// 响应式数据
const syncing = ref(false)
const currentIndex = ref(0)
const currentProvider = ref<any>(null)
const results = ref<any[]>([])

// 计算属性
const totalProviders = computed(() => props.providers.length)

const progressPercent = computed(() => {
  if (totalProviders.value === 0) return 0
  return Math.round((currentIndex.value / totalProviders.value) * 100)
})

const totalDomains = computed(() => {
  return results.value.reduce((sum, result) => sum + (result.domains_count || 0), 0)
})

const successCount = computed(() => {
  return results.value.filter(result => result.status === 'success').length
})

const errorCount = computed(() => {
  return results.value.filter(result => result.status === 'error').length
})

// 监听可见性变化
watch(() => props.visible, (newVal) => {
  if (newVal) {
    resetProgress()
    startSync()
  }
})

// 重置进度
const resetProgress = () => {
  syncing.value = true
  currentIndex.value = 0
  currentProvider.value = null
  results.value = []
}

// 开始同步
const startSync = async () => {
  for (let i = 0; i < props.providers.length; i++) {
    if (!syncing.value) break // 如果取消了同步，则退出

    const provider = props.providers[i]
    currentIndex.value = i + 1
    currentProvider.value = provider

    // 添加到结果列表（初始状态为同步中）
    const resultItem = {
      id: provider.id,
      name: provider.name,
      type: provider.type,
      status: 'syncing',
      domains_count: 0,
      message: '同步中...'
    }
    results.value.push(resultItem)

    try {
      // 这里应该调用真实的同步API
      // const response = await dnsProviderApi.syncDomains(provider.id)
      
      // 模拟API调用
      await new Promise(resolve => setTimeout(resolve, 1000 + Math.random() * 2000))
      
      // 模拟成功结果
      const domainsCount = Math.floor(Math.random() * 20) + 1
      resultItem.status = 'success'
      resultItem.domains_count = domainsCount
      resultItem.message = `成功同步 ${domainsCount} 个域名`
      
    } catch (error) {
      resultItem.status = 'error'
      resultItem.message = error instanceof Error ? error.message : '同步失败'
    }
  }

  syncing.value = false
  currentProvider.value = null
}

// 处理关闭
const handleClose = () => {
  emit('update:visible', false)
}

// 处理取消
const handleCancel = () => {
  syncing.value = false
  emit('cancel')
}
</script>

<style scoped>
.sync-progress {
  padding: 20px 0;
}

.progress-header {
  margin-bottom: 20px;
}

.progress-header h4 {
  margin: 0 0 10px 0;
  color: #303133;
}

.current-sync {
  margin-bottom: 20px;
}

.sync-results {
  margin: 20px 0;
}

.sync-summary {
  margin-top: 20px;
}

.summary-item {
  text-align: center;
  padding: 10px;
  background: #f8f9fa;
  border-radius: 6px;
}

.summary-value {
  font-size: 24px;
  font-weight: bold;
  color: #409eff;
  margin-bottom: 5px;
}

.summary-label {
  font-size: 14px;
  color: #909399;
}
</style>
