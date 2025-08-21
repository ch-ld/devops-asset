<template>
  <el-dialog
    v-model="visible"
    title="选择连接IP"
    width="500px"
    :before-close="handleClose"
    append-to-body
  >
    <div class="ip-selection-content">
      <div class="host-info">
        <div class="host-name">
          <el-icon><Monitor /></el-icon>
          <span>{{ host?.name || '未知主机' }}</span>
        </div>
        <div class="connection-type">
          <el-icon><Connection /></el-icon>
          <span>{{ connectionTypeText }}</span>
        </div>
      </div>

      <div class="ip-options">
        <div class="section-title">请选择要使用的IP地址：</div>
        
        <!-- 公网IP选项 -->
        <div v-if="publicIPs.length > 0" class="ip-group">
          <div class="ip-group-title">
            <el-icon><Connection /></el-icon>
            公网IP
          </div>
          <el-radio-group v-model="selectedIP" class="ip-radio-group">
            <el-radio
              v-for="(ip, index) in publicIPs"
              :key="`public-${index}`"
              :label="`public:${ip}`"
              class="ip-radio"
            >
              <div class="ip-option">
                <span class="ip-address">{{ ip }}</span>
                <el-tag type="success" size="small">公网</el-tag>
              </div>
            </el-radio>
          </el-radio-group>
        </div>

        <!-- 私网IP选项 -->
        <div v-if="privateIPs.length > 0" class="ip-group">
          <div class="ip-group-title">
            <i class="el-icon-office-building"></i>
            私网IP
          </div>
          <el-radio-group v-model="selectedIP" class="ip-radio-group">
            <el-radio
              v-for="(ip, index) in privateIPs"
              :key="`private-${index}`"
              :label="`private:${ip}`"
              class="ip-radio"
            >
              <div class="ip-option">
                <span class="ip-address">{{ ip }}</span>
                <el-tag type="info" size="small">私网</el-tag>
              </div>
            </el-radio>
          </el-radio-group>
        </div>

        <!-- 无可用IP提示 -->
        <div v-if="publicIPs.length === 0 && privateIPs.length === 0" class="no-ip-warning">
          <el-alert
            title="无可用IP地址"
            description="该主机没有配置任何可用的IP地址，请先在主机管理中配置IP地址。"
            type="warning"
            :closable="false"
          />
        </div>
      </div>

      <!-- 连接提示 -->
      <div v-if="selectedIP" class="connection-hint">
        <el-alert
          :title="getConnectionHint()"
          type="info"
          :closable="false"
        />
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button
          type="primary"
          @click="handleConfirm"
          :disabled="!selectedIP"
        >
          连接
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { Connection, Monitor, OfficeBuilding } from '@element-plus/icons-vue'

// Props
interface Props {
  host?: any
  connectionType?: 'ssh' | 'sftp'
}

const props = withDefaults(defineProps<Props>(), {
  connectionType: 'ssh'
})

// Emits
const emit = defineEmits<{
  confirm: [ipType: 'public' | 'private', ip: string]
  cancel: []
}>()

// 响应式数据
const visible = ref(false)
const selectedIP = ref('')

// 计算属性
const connectionTypeText = computed(() => {
  return props.connectionType === 'ssh' ? 'SSH终端连接' : 'SFTP文件管理'
})

const publicIPs = computed(() => {
  if (!props.host) return []
  
  try {
    const ips = Array.isArray(props.host.public_ip) 
      ? props.host.public_ip 
      : JSON.parse(props.host.public_ip || '[]')
    return ips.filter((ip: string) => ip && ip.trim())
  } catch {
    return []
  }
})

const privateIPs = computed(() => {
  if (!props.host) return []
  
  try {
    const ips = Array.isArray(props.host.private_ip) 
      ? props.host.private_ip 
      : JSON.parse(props.host.private_ip || '[]')
    return ips.filter((ip: string) => ip && ip.trim())
  } catch {
    return []
  }
})

// 方法
const open = () => {
  visible.value = true
  selectedIP.value = ''
  
  // 自动选择第一个可用IP
  if (publicIPs.value.length > 0) {
    selectedIP.value = `public:${publicIPs.value[0]}`
  } else if (privateIPs.value.length > 0) {
    selectedIP.value = `private:${privateIPs.value[0]}`
  }
}

const handleClose = () => {
  visible.value = false
  emit('cancel')
}

const handleConfirm = () => {
  if (!selectedIP.value) return
  
  const [type, ip] = selectedIP.value.split(':')
  emit('confirm', type as 'public' | 'private', ip)
  visible.value = false
}

const getConnectionHint = () => {
  if (!selectedIP.value) return ''
  
  const [type] = selectedIP.value.split(':')
  if (type === 'public') {
    return '将通过公网IP连接，请确保网络可达且防火墙已开放相应端口'
  } else {
    return '将通过私网IP连接，请确保在同一网络环境中且网络可达'
  }
}

// 监听主机变化，自动选择IP
watch(() => props.host, () => {
  if (visible.value) {
    selectedIP.value = ''
    if (publicIPs.value.length > 0) {
      selectedIP.value = `public:${publicIPs.value[0]}`
    } else if (privateIPs.value.length > 0) {
      selectedIP.value = `private:${privateIPs.value[0]}`
    }
  }
})

// 暴露方法
defineExpose({
  open
})
</script>

<style lang="scss" scoped>
.ip-selection-content {
  .host-info {
    background: #f8f9fa;
    border-radius: 8px;
    padding: 16px;
    margin-bottom: 20px;
    
    .host-name {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 16px;
      font-weight: 600;
      color: #2c3e50;
      margin-bottom: 8px;
      
      i {
        color: #6366f1;
      }
    }
    
    .connection-type {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 14px;
      color: #6b7280;
      
      i {
        color: #10b981;
      }
    }
  }
  
  .ip-options {
    .section-title {
      font-size: 14px;
      font-weight: 600;
      color: #374151;
      margin-bottom: 16px;
    }
    
    .ip-group {
      margin-bottom: 20px;
      
      .ip-group-title {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 14px;
        font-weight: 600;
        color: #4b5563;
        margin-bottom: 12px;
        
        i {
          color: #6366f1;
        }
      }
      
      .ip-radio-group {
        display: flex;
        flex-direction: column;
        gap: 8px;
        
        .ip-radio {
          margin: 0;
          padding: 12px;
          border: 1px solid #e5e7eb;
          border-radius: 8px;
          transition: all 0.2s ease;
          
          &:hover {
            border-color: #6366f1;
            background: #f8faff;
          }
          
          &.is-checked {
            border-color: #6366f1;
            background: #f0f4ff;
          }
          
          .ip-option {
            display: flex;
            align-items: center;
            justify-content: space-between;
            width: 100%;
            
            .ip-address {
              font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
              font-size: 14px;
              font-weight: 500;
              color: #1f2937;
            }
          }
        }
      }
    }
    
    .no-ip-warning {
      margin: 20px 0;
    }
  }
  
  .connection-hint {
    margin-top: 16px;
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
