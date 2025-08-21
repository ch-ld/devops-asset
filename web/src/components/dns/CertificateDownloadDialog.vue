<template>
  <el-dialog
    v-model="dialogVisible"
    title="下载证书"
    width="600px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <div v-if="certificate" class="download-dialog">
      <!-- 证书信息 -->
      <div class="cert-info">
        <h4>{{ certificate.common_name }}</h4>
        <p class="cert-meta">
          <el-tag :type="getStatusType(certificate.status)">{{ getStatusText(certificate.status) }}</el-tag>
          <span class="divider">|</span>
          <span>{{ getCATypeName(certificate.ca_type) }}</span>
          <span class="divider">|</span>
          <span>到期时间：{{ formatDate(certificate.valid_to) }}</span>
        </p>
      </div>

      <!-- 下载格式选择 -->
      <div class="format-selection">
        <h5>选择下载格式</h5>
        <el-row :gutter="16">
          <el-col :span="12" v-for="format in downloadFormats" :key="format.value">
            <div 
              class="format-card" 
              :class="{ active: selectedFormat === format.value }"
              @click="selectedFormat = format.value"
            >
              <div class="format-header">
                <el-icon class="format-icon">
                  <component :is="format.icon" />
                </el-icon>
                <span class="format-name">{{ format.name }}</span>
              </div>
              <p class="format-desc">{{ format.description }}</p>
              <div class="format-files">
                <el-tag 
                  v-for="file in format.files" 
                  :key="file" 
                  size="small" 
                  type="info"
                >
                  {{ file }}
                </el-tag>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>

      <!-- 使用说明 -->
      <div class="usage-guide" v-if="selectedFormatInfo">
        <h5>使用说明</h5>
        <div class="guide-content">
          <p>{{ selectedFormatInfo.guide }}</p>
          <div v-if="selectedFormatInfo.example" class="example-code">
            <h6>配置示例：</h6>
            <pre><code>{{ selectedFormatInfo.example }}</code></pre>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button 
          type="primary" 
          @click="handleDownload"
          :loading="downloading"
          :disabled="!selectedFormat"
        >
          <el-icon><Download /></el-icon>
          下载证书
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Download, Document, Setting, Monitor, Coffee } from '@element-plus/icons-vue'
import { certificateApi } from '@/api/dns/certificate'
import type { Certificate } from '@/types/dns'

interface Props {
  visible: boolean
  certificate: Certificate | null
}

interface Emits {
  (e: 'update:visible', value: boolean): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const selectedFormat = ref('nginx')
const downloading = ref(false)

// 下载格式配置
const downloadFormats = [
  {
    value: 'nginx',
    name: 'Nginx',
    icon: Monitor,
    description: '适用于 Nginx 服务器',
    files: ['cert.crt', 'private.key'],
    guide: '将证书文件上传到服务器，在 Nginx 配置中指定证书和私钥路径。',
    example: `server {
    listen 443 ssl;
    server_name example.com;
    ssl_certificate /path/to/cert.crt;
    ssl_certificate_key /path/to/private.key;
}`
  },
  {
    value: 'apache',
    name: 'Apache',
    icon: Document,
    description: '适用于 Apache 服务器',
    files: ['cert.crt', 'private.key', 'chain.crt'],
    guide: '将证书文件上传到服务器，在 Apache 虚拟主机配置中指定证书路径。',
    example: `<VirtualHost *:443>
    ServerName example.com
    SSLEngine on
    SSLCertificateFile /path/to/cert.crt
    SSLCertificateKeyFile /path/to/private.key
    SSLCertificateChainFile /path/to/chain.crt
</VirtualHost>`
  },
  {
    value: 'iis',
    name: 'IIS',
    icon: Setting,
    description: '适用于 Windows IIS 服务器',
    files: ['cert.pfx'],
    guide: '导入 PFX 证书到 Windows 证书存储，然后在 IIS 管理器中绑定到网站。',
    example: `1. 双击 .pfx 文件导入证书
2. 打开 IIS 管理器
3. 选择网站 -> 绑定 -> 添加 HTTPS 绑定
4. 选择导入的证书`
  },
  {
    value: 'tomcat',
    name: 'Tomcat',
    icon: Coffee,
    description: '适用于 Java Tomcat 服务器',
    files: ['keystore.jks'],
    guide: '将 JKS 密钥库文件放置在 Tomcat 配置目录，在 server.xml 中配置 HTTPS 连接器。',
    example: `<Connector port="8443" protocol="HTTP/1.1"
    SSLEnabled="true" scheme="https" secure="true"
    keystoreFile="conf/keystore.jks"
    keystorePass="password" />`
  }
]

const selectedFormatInfo = computed(() => {
  return downloadFormats.find(f => f.value === selectedFormat.value)
})

const handleClose = () => {
  dialogVisible.value = false
  selectedFormat.value = 'nginx'
}

const handleDownload = async () => {
  if (!props.certificate || !selectedFormat.value) return

  try {
    downloading.value = true
    ElMessage.info('正在准备下载，请稍候...')
    
    const response = await certificateApi.download(props.certificate.id, selectedFormat.value)
    
    // 创建下载链接
    const url = window.URL.createObjectURL(response as Blob)
    const link = document.createElement('a')
    link.href = url
    
    const domainName = props.certificate.common_name || `cert_${props.certificate.id}`
    const formatInfo = selectedFormatInfo.value
    const extension = getFileExtension(selectedFormat.value)
    
    link.download = `${domainName}_${selectedFormat.value}.${extension}`
    link.click()
    window.URL.revokeObjectURL(url)
    
    ElMessage.success('证书下载成功')
    handleClose()
  } catch (error: any) {
    console.error('证书下载失败:', error)
    const errorMessage = error?.response?.data?.message || error?.message || '证书下载失败，请重试'
    ElMessage.error(errorMessage)
  } finally {
    downloading.value = false
  }
}

const getFileExtension = (format: string) => {
  const extensionMap: Record<string, string> = {
    'nginx': 'crt',
    'apache': 'crt',
    'iis': 'pfx',
    'tomcat': 'jks',
    'pem': 'pem',
    'crt': 'crt',
    'key': 'key',
    'chain': 'pem'
  }
  return extensionMap[format] || 'txt'
}

// 工具方法
const getStatusType = (status: string) => {
  const statusMap: Record<string, any> = {
    pending: 'warning',
    issued: 'success',
    expired: 'danger',
    revoked: 'info'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: '申请中',
    issued: '已签发',
    expired: '已过期',
    revoked: '已吊销'
  }
  return statusMap[status] || status
}

const getCATypeName = (caType: string) => {
  const typeMap: Record<string, string> = {
    letsencrypt: "Let's Encrypt",
    zerossl: 'ZeroSSL',
    custom: '自定义'
  }
  return typeMap[caType] || caType
}

const formatDate = (dateValue: string | number) => {
  if (!dateValue) return '-'
  
  try {
    let date: Date
    
    if (typeof dateValue === 'number') {
      if (dateValue < 10000000000) {
        date = new Date(dateValue * 1000)
      } else {
        date = new Date(dateValue)
      }
    } else {
      date = new Date(dateValue)
    }
    
    if (isNaN(date.getTime())) {
      return dateValue.toString()
    }
    
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit'
    })
  } catch {
    return dateValue.toString()
  }
}
</script>

<style scoped>
.download-dialog {
  padding: 0;
}

.cert-info {
  margin-bottom: 24px;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
}

.cert-info h4 {
  margin: 0 0 8px 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.cert-meta {
  margin: 0;
  font-size: 14px;
  color: #606266;
  display: flex;
  align-items: center;
  gap: 8px;
}

.divider {
  color: #dcdfe6;
}

.format-selection {
  margin-bottom: 24px;
}

.format-selection h5 {
  margin: 0 0 16px 0;
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}

.format-card {
  padding: 16px;
  border: 2px solid #e4e7ed;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  margin-bottom: 16px;
  height: 140px;
  display: flex;
  flex-direction: column;
}

.format-card:hover {
  border-color: #409eff;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.1);
}

.format-card.active {
  border-color: #409eff;
  background: #f0f8ff;
}

.format-header {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.format-icon {
  font-size: 20px;
  color: #409eff;
  margin-right: 8px;
}

.format-name {
  font-weight: 600;
  color: #303133;
}

.format-desc {
  margin: 0 0 12px 0;
  font-size: 13px;
  color: #606266;
  flex: 1;
}

.format-files {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.usage-guide {
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
}

.usage-guide h5 {
  margin: 0 0 12px 0;
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}

.guide-content p {
  margin: 0 0 12px 0;
  font-size: 14px;
  color: #606266;
  line-height: 1.5;
}

.example-code {
  margin-top: 12px;
}

.example-code h6 {
  margin: 0 0 8px 0;
  font-size: 13px;
  font-weight: 600;
  color: #303133;
}

.example-code pre {
  margin: 0;
  padding: 12px;
  background: #f5f5f5;
  border-radius: 4px;
  font-size: 12px;
  line-height: 1.4;
  overflow-x: auto;
}

.example-code code {
  color: #e96900;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
