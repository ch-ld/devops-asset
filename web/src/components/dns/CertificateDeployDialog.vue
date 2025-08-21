<template>
  <el-dialog
    v-model="dialogVisible"
    title="部署证书"
    width="600px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <div class="deploy-dialog-content">
      <!-- 证书信息 -->
      <div class="cert-info">
        <h4>证书信息</h4>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="域名">
            {{ certificate?.common_name || certificate?.domain_name || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType(certificate?.status)">
              {{ getStatusText(certificate?.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="过期时间">
            {{ formatDate(certificate?.expires_at) }}
          </el-descriptions-item>
          <el-descriptions-item label="证书类型">
            {{ getCATypeName(certificate?.ca_type) }}
          </el-descriptions-item>
        </el-descriptions>
      </div>

      <!-- 部署配置 -->
      <div class="deploy-config">
        <h4>部署配置</h4>
        <el-form ref="formRef" :model="deployForm" :rules="deployRules" label-width="100px">
          <el-form-item label="目标主机" prop="hostIds" required>
            <el-select
              v-model="deployForm.hostIds"
              multiple
              placeholder="请选择要部署的主机"
              style="width: 100%"
              filterable
              remote
              reserve-keyword
              :remote-method="searchHosts"
              :loading="hostsLoading"
            >
              <el-option
                v-for="host in filteredHosts.length > 0 ? filteredHosts : availableHosts"
                :key="host.id"
                :label="`${host.name} (${host.ip})`"
                :value="host.id"
              >
                <div class="host-option">
                  <div class="host-info">
                    <span class="host-name">{{ host.name }}</span>
                    <span class="host-ip">{{ host.ip }}</span>
                  </div>
                  <el-tag
                    :type="getHostStatusType(host.status)"
                    size="small"
                  >
                    {{ getHostStatusText(host.status) }}
                  </el-tag>
                </div>
              </el-option>
            </el-select>
          </el-form-item>

          <el-form-item label="部署路径" prop="deployPath">
            <el-input
              v-model="deployForm.deployPath"
              placeholder="证书部署路径，如：/etc/ssl/certs/"
            />
            <div class="form-tip">
              证书文件将部署到指定路径，请确保路径存在且有写入权限
            </div>
          </el-form-item>

          <el-form-item label="服务名称">
            <el-input
              v-model="deployForm.serviceName"
              placeholder="相关服务名称，如：nginx, apache"
            />
          </el-form-item>

          <el-form-item label="重启命令">
            <el-input
              v-model="deployForm.restartCommand"
              placeholder="部署后执行的重启命令，如：systemctl reload nginx"
            />
            <div class="form-tip">
              部署完成后执行的命令，用于重启相关服务使证书生效
            </div>
          </el-form-item>

          <el-form-item label="部署选项">
            <el-checkbox-group v-model="deployForm.options">
              <el-checkbox label="backup">备份现有证书</el-checkbox>
              <el-checkbox label="verify">部署前验证证书</el-checkbox>
              <el-checkbox label="notification">发送部署通知</el-checkbox>
            </el-checkbox-group>
          </el-form-item>
        </el-form>
      </div>

      <!-- 部署预览 -->
      <div class="deploy-preview" v-if="deployForm.hostIds.length > 0">
        <h4>部署预览</h4>
        <el-alert
          title="部署信息确认"
          type="info"
          :closable="false"
          show-icon
        >
          <template #default>
            <p>将把证书部署到 <strong>{{ deployForm.hostIds.length }}</strong> 台主机：</p>
            <ul class="host-list">
              <li v-for="hostId in deployForm.hostIds" :key="hostId">
                {{ getHostName(hostId) }}
              </li>
            </ul>
            <p>部署路径：<code>{{ deployForm.deployPath }}</code></p>
            <p v-if="deployForm.restartCommand">重启命令：<code>{{ deployForm.restartCommand }}</code></p>
          </template>
        </el-alert>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button
          type="primary"
          :loading="deploying"
          @click="handleDeploy"
          :disabled="!canDeploy"
        >
          {{ deploying ? '部署中...' : '开始部署' }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { certificateApi } from '@/api/dns/certificate'
import { getHostListWithCount } from '@/api/system/host'
import type { Certificate } from '@/types/dns'

interface Props {
  visible: boolean
  certificate: Certificate | null
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'success'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 状态管理
const deploying = ref(false)
const hostsLoading = ref(false)
const availableHosts = ref<any[]>([])
const filteredHosts = ref<any[]>([])
const formRef = ref()

// 表单数据
const deployForm = reactive({
  hostIds: [] as number[],
  deployPath: '/etc/ssl/certs/',
  serviceName: 'nginx',
  restartCommand: 'systemctl reload nginx',
  options: ['backup', 'verify'] as string[]
})

// 验证规则
const deployRules = {
  hostIds: [
    { required: true, message: '请选择要部署的主机', trigger: 'change' }
  ],
  deployPath: [
    { required: true, message: '请输入部署路径', trigger: 'blur' }
  ]
}

// 计算属性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const canDeploy = computed(() => {
  return deployForm.hostIds.length > 0 && 
         deployForm.deployPath.trim() && 
         !deploying.value &&
         props.certificate?.status === 'issued'
})

// 工具方法
const getStatusType = (status?: string) => {
  const statusMap: Record<string, any> = {
    pending: 'warning',
    issued: 'success',
    expired: 'danger',
    revoked: 'info'
  }
  return statusMap[status || ''] || 'info'
}

const getStatusText = (status?: string) => {
  const statusMap: Record<string, string> = {
    pending: '申请中',
    issued: '已签发',
    expired: '已过期',
    revoked: '已吊销'
  }
  return statusMap[status || ''] || status || '-'
}

const getCATypeName = (caType?: string) => {
  const typeMap: Record<string, string> = {
    letsencrypt: "Let's Encrypt",
    zerossl: 'ZeroSSL',
    custom: '自定义'
  }
  return typeMap[caType || ''] || caType || '-'
}

const formatDate = (dateValue?: string | number) => {
  if (!dateValue) return '-'
  
  try {
    let date: Date
    if (typeof dateValue === 'number') {
      date = new Date(dateValue < 10000000000 ? dateValue * 1000 : dateValue)
    } else {
      date = new Date(dateValue)
    }
    
    if (isNaN(date.getTime()) || date.getFullYear() === 1970) {
      return '-'
    }
    
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch (error) {
    return '-'
  }
}

const getHostName = (hostId: number) => {
  const host = availableHosts.value.find(h => h.id === hostId)
  return host ? `${host.name} (${host.ip})` : `主机 ID: ${hostId}`
}

const getHostStatusType = (status: string) => {
  const statusMap: Record<string, string> = {
    online: 'success',
    running: 'success',
    active: 'success',
    offline: 'danger',
    stopped: 'danger',
    error: 'danger',
    unknown: 'info'
  }
  return statusMap[status] || 'info'
}

const getHostStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    online: '在线',
    running: '运行中',
    active: '活跃',
    offline: '离线',
    stopped: '已停止',
    error: '错误',
    unknown: '未知'
  }
  return statusMap[status] || status
}

const searchHosts = async (query: string) => {
  if (!query) {
    filteredHosts.value = availableHosts.value
    return
  }
  
  // 本地过滤
  filteredHosts.value = availableHosts.value.filter(host => 
    host.name.toLowerCase().includes(query.toLowerCase()) ||
    host.ip.includes(query) ||
    host.provider.toLowerCase().includes(query.toLowerCase())
  )
}

// 方法
const loadHosts = async () => {
  try {
    hostsLoading.value = true
    
    // 使用正确的CMDB主机查询API
    const response = await getHostListWithCount({
      page: 1,
      page_size: 100
      // 移除status过滤，显示所有主机状态
    })
    
    console.log('证书部署 - CMDB主机API响应:', response)
    
    // 根据API响应结构解析数据
    // getHostListWithCount返回格式: {code: 200, data: [...], count: 10}
    let hostList = []
    
    if (response && response.data && Array.isArray(response.data)) {
      // 标准API响应格式
      hostList = response.data
    } else if (response && Array.isArray(response)) {
      // 直接数组格式
      hostList = response
    } else {
      console.warn('证书部署 - 未识别的主机API响应格式:', response)
      hostList = []
    }
    
    // 映射主机数据，确保字段完整
    availableHosts.value = hostList.map((host: any) => ({
      id: host.id,
      name: host.name || host.hostname || `主机-${host.id}`,
      ip: host.ip || host.private_ip || host.public_ip || '未知IP', 
      status: host.status || 'unknown',
      provider: host.provider || host.provider_name || '未知',
      region: host.region || '未知'
    }))
    
    console.log(`证书部署 - 成功加载 ${availableHosts.value.length} 台主机:`, availableHosts.value)
    filteredHosts.value = availableHosts.value // 初始化过滤列表
    
  } catch (error) {
    console.error('证书部署 - 加载主机列表失败:', error)
    ElMessage.error('加载主机列表失败，请确保CMDB主机管理模块正常运行')
  } finally {
    hostsLoading.value = false
  }
}

const handleDeploy = async () => {
  if (!props.certificate || !formRef.value) return
  
  try {
    await formRef.value.validate()
    deploying.value = true
    
    await certificateApi.deploy(props.certificate.id, deployForm.hostIds)
    
    ElMessage.success('证书部署任务已启动，请稍后查看部署状态')
    emit('success')
    handleClose()
  } catch (error: any) {
    console.error('证书部署失败:', error)
    ElMessage.error(error.message || '证书部署失败')
  } finally {
    deploying.value = false
  }
}

const handleClose = () => {
  emit('update:visible', false)
  // 重置表单
  deployForm.hostIds = []
  deployForm.deployPath = '/etc/ssl/certs/'
  deployForm.serviceName = 'nginx'
  deployForm.restartCommand = 'systemctl reload nginx'
  deployForm.options = ['backup', 'verify']
}

// 监听对话框显示
watch(() => props.visible, (visible) => {
  if (visible) {
    loadHosts()
  }
})
</script>

<style scoped lang="scss">
.deploy-dialog-content {
  .cert-info,
  .deploy-config,
  .deploy-preview {
    margin-bottom: 24px;
    
    h4 {
      margin: 0 0 16px 0;
      font-size: 16px;
      font-weight: 600;
      color: #303133;
    }
  }
  
  .host-option {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    
    .host-info {
      display: flex;
      flex-direction: column;
      align-items: flex-start;
      
      .host-name {
        font-weight: 500;
        color: #303133;
      }
      
      .host-ip {
        font-size: 12px;
        color: #909399;
        margin-top: 2px;
      }
    }
  }
  
  .form-tip {
    margin-top: 4px;
    font-size: 12px;
    color: #909399;
    line-height: 1.4;
  }
  
  .host-list {
    margin: 8px 0;
    padding-left: 20px;
    
    li {
      margin: 4px 0;
      color: #606266;
    }
  }
  
  code {
    background: #f0f2f5;
    padding: 2px 6px;
    border-radius: 4px;
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
    font-size: 12px;
    color: #e6a23c;
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>