<template>
  <el-dialog
    v-model="dialogVisible"
    title="导出证书配置"
    width="700px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <div class="export-dialog-content">
      <!-- 证书信息 -->
      <div class="cert-info">
        <h4>证书信息</h4>
        <el-descriptions :column="2" border size="small">
          <el-descriptions-item label="域名">
            {{ certificate?.common_name || certificate?.domain_name || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType(certificate?.status)" size="small">
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

      <!-- 导出选项 -->
      <div class="export-options">
        <h4>导出选项</h4>
        <el-form :model="exportForm" label-width="120px">
          <el-form-item label="服务器类型">
            <el-radio-group v-model="exportForm.serverType" @change="handleServerTypeChange">
              <el-radio-button value="nginx">Nginx</el-radio-button>
              <el-radio-button value="apache">Apache</el-radio-button>
              <el-radio-button value="iis">IIS</el-radio-button>
              <el-radio-button value="tomcat">Tomcat</el-radio-button>
              <el-radio-button value="custom">自定义</el-radio-button>
            </el-radio-group>
          </el-form-item>

          <el-form-item label="配置格式">
            <el-checkbox-group v-model="exportForm.formats">
              <el-checkbox label="config">服务器配置文件</el-checkbox>
              <el-checkbox label="script">部署脚本</el-checkbox>
              <el-checkbox label="docker">Docker配置</el-checkbox>
              <el-checkbox label="docs">说明文档</el-checkbox>
            </el-checkbox-group>
          </el-form-item>

          <el-form-item label="域名配置">
            <el-input
              v-model="exportForm.domainConfig"
              type="textarea"
              :rows="3"
              placeholder="请输入域名相关配置，如虚拟主机、端口等"
            />
          </el-form-item>

          <el-form-item label="证书路径">
            <el-input
              v-model="exportForm.certPath"
              placeholder="证书文件在服务器上的路径"
            />
          </el-form-item>

          <el-form-item label="私钥路径">
            <el-input
              v-model="exportForm.keyPath"
              placeholder="私钥文件在服务器上的路径"
            />
          </el-form-item>

          <el-form-item label="证书链路径" v-if="exportForm.serverType !== 'iis'">
            <el-input
              v-model="exportForm.chainPath"
              placeholder="证书链文件在服务器上的路径（可选）"
            />
          </el-form-item>
        </el-form>
      </div>

      <!-- 配置预览 -->
      <div class="config-preview" v-if="configPreview">
        <h4>配置预览</h4>
        <el-tabs v-model="activeTab" type="border-card">
          <el-tab-pane
            v-for="(content, format) in configPreview"
            :key="format"
            :label="getFormatLabel(format)"
            :name="format"
          >
            <div class="config-content">
              <el-input
                :model-value="content"
                type="textarea"
                :rows="15"
                readonly
                class="config-textarea"
              />
              <div class="config-actions">
                <el-button size="small" @click="copyConfig(content)">
                  <el-icon><CopyDocument /></el-icon>
                  复制配置
                </el-button>
                <el-button size="small" type="primary" @click="downloadConfig(format, content)">
                  <el-icon><Download /></el-icon>
                  下载文件
                </el-button>
              </div>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">关闭</el-button>
        <el-button type="primary" @click="generateConfig" :loading="generating">
          {{ generating ? '生成中...' : '生成配置' }}
        </el-button>
        <el-button type="success" @click="exportAll" :disabled="!configPreview">
          <el-icon><Download /></el-icon>
          导出全部
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { CopyDocument, Download } from '@element-plus/icons-vue'
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

// 状态管理
const generating = ref(false)
const activeTab = ref('config')
const configPreview = ref<Record<string, string> | null>(null)

// 表单数据
const exportForm = reactive({
  serverType: 'nginx',
  formats: ['config'],
  domainConfig: '',
  certPath: '/etc/ssl/certs/cert.pem',
  keyPath: '/etc/ssl/private/cert.key',
  chainPath: '/etc/ssl/certs/chain.pem'
})

// 计算属性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
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
      day: '2-digit'
    })
  } catch (error) {
    return '-'
  }
}

const getFormatLabel = (format: string) => {
  const labels: Record<string, string> = {
    config: '配置文件',
    script: '部署脚本',
    docker: 'Docker配置',
    docs: '说明文档'
  }
  return labels[format] || format
}

// 方法
const handleServerTypeChange = (serverType: string) => {
  // 根据服务器类型更新默认路径
  const pathMap: Record<string, any> = {
    nginx: {
      certPath: '/etc/ssl/certs/cert.pem',
      keyPath: '/etc/ssl/private/cert.key',
      chainPath: '/etc/ssl/certs/chain.pem'
    },
    apache: {
      certPath: '/etc/ssl/certs/cert.crt',
      keyPath: '/etc/ssl/private/cert.key',
      chainPath: '/etc/ssl/certs/chain.crt'
    },
    iis: {
      certPath: 'C:\\SSL\\cert.pfx',
      keyPath: '',
      chainPath: ''
    },
    tomcat: {
      certPath: '/opt/tomcat/ssl/keystore.jks',
      keyPath: '',
      chainPath: ''
    }
  }
  
  if (pathMap[serverType]) {
    Object.assign(exportForm, pathMap[serverType])
  }
}

const generateConfig = async () => {
  if (!props.certificate) return
  
  try {
    generating.value = true
    
    // 模拟生成配置
    const configs: Record<string, string> = {}
    
    if (exportForm.formats.includes('config')) {
      configs.config = generateServerConfig()
    }
    
    if (exportForm.formats.includes('script')) {
      configs.script = generateDeployScript()
    }
    
    if (exportForm.formats.includes('docker')) {
      configs.docker = generateDockerConfig()
    }
    
    if (exportForm.formats.includes('docs')) {
      configs.docs = generateDocumentation()
    }
    
    configPreview.value = configs
    activeTab.value = Object.keys(configs)[0]
    
    ElMessage.success('配置生成成功')
  } catch (error) {
    ElMessage.error('配置生成失败')
  } finally {
    generating.value = false
  }
}

const generateServerConfig = () => {
  const domain = props.certificate?.common_name || 'example.com'
  
  switch (exportForm.serverType) {
    case 'nginx':
      return `server {
    listen 443 ssl;
    server_name ${domain};
    
    ssl_certificate ${exportForm.certPath};
    ssl_certificate_key ${exportForm.keyPath};
    ${exportForm.chainPath ? `ssl_trusted_certificate ${exportForm.chainPath};` : ''}
    
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    
    location / {
        root /var/www/html;
        index index.html index.htm;
    }
}

# HTTP重定向到HTTPS
server {
    listen 80;
    server_name ${domain};
    return 301 https://$server_name$request_uri;
}`

    case 'apache':
      return `<VirtualHost *:443>
    ServerName ${domain}
    DocumentRoot /var/www/html
    
    SSLEngine on
    SSLCertificateFile ${exportForm.certPath}
    SSLCertificateKeyFile ${exportForm.keyPath}
    ${exportForm.chainPath ? `SSLCertificateChainFile ${exportForm.chainPath}` : ''}
    
    SSLProtocol all -SSLv3 -TLSv1 -TLSv1.1
    SSLCipherSuite ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256
    SSLHonorCipherOrder off
</VirtualHost>

<VirtualHost *:80>
    ServerName ${domain}
    Redirect permanent / https://${domain}/
</VirtualHost>`

    case 'iis':
      return `# IIS证书配置说明
1. 将证书文件转换为PFX格式
2. 在IIS管理器中导入证书
3. 在站点绑定中选择证书

# PowerShell命令示例
Import-PfxCertificate -FilePath "${exportForm.certPath}" -CertStoreLocation Cert:\\LocalMachine\\My
New-WebBinding -Name "Default Web Site" -IP "*" -Port 443 -Protocol https`

    default:
      return `# ${exportForm.serverType} 配置
# 请根据您的服务器类型自定义配置

证书文件: ${exportForm.certPath}
私钥文件: ${exportForm.keyPath}
证书链文件: ${exportForm.chainPath}
域名: ${domain}`
  }
}

const generateDeployScript = () => {
  const domain = props.certificate?.common_name || 'example.com'
  
  return `#!/bin/bash
# ${domain} 证书部署脚本

set -e

echo "开始部署证书..."

# 创建证书目录
sudo mkdir -p $(dirname ${exportForm.certPath})
sudo mkdir -p $(dirname ${exportForm.keyPath})
${exportForm.chainPath ? `sudo mkdir -p $(dirname ${exportForm.chainPath})` : ''}

# 设置文件权限
sudo chmod 644 ${exportForm.certPath}
sudo chmod 600 ${exportForm.keyPath}
${exportForm.chainPath ? `sudo chmod 644 ${exportForm.chainPath}` : ''}

# 重启服务
${exportForm.serverType === 'nginx' ? 'sudo systemctl reload nginx' : ''}
${exportForm.serverType === 'apache' ? 'sudo systemctl reload apache2' : ''}

echo "证书部署完成！"

# 验证证书
openssl x509 -in ${exportForm.certPath} -text -noout
`
}

const generateDockerConfig = () => {
  const domain = props.certificate?.common_name || 'example.com'
  
  if (exportForm.serverType === 'nginx') {
    return `# Docker Compose配置
version: '3.8'

services:
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ${exportForm.certPath}:/etc/ssl/certs/cert.pem:ro
      - ${exportForm.keyPath}:/etc/ssl/private/cert.key:ro
      ${exportForm.chainPath ? `- ${exportForm.chainPath}:/etc/ssl/certs/chain.pem:ro` : ''}
    restart: unless-stopped

# Dockerfile示例
FROM nginx:alpine
COPY nginx.conf /etc/nginx/nginx.conf
COPY cert.pem /etc/ssl/certs/
COPY cert.key /etc/ssl/private/
${exportForm.chainPath ? 'COPY chain.pem /etc/ssl/certs/' : ''}
EXPOSE 80 443`
  }
  
  return `# Docker配置示例
# 请根据您的应用类型自定义Docker配置`
}

const generateDocumentation = () => {
  const domain = props.certificate?.common_name || 'example.com'
  
  return `# ${domain} SSL证书部署文档

## 证书信息
- 域名: ${domain}
- 证书类型: ${getCATypeName(props.certificate?.ca_type)}
- 过期时间: ${formatDate(props.certificate?.expires_at)}
- 服务器类型: ${exportForm.serverType}

## 文件路径
- 证书文件: ${exportForm.certPath}
- 私钥文件: ${exportForm.keyPath}
${exportForm.chainPath ? `- 证书链文件: ${exportForm.chainPath}` : ''}

## 部署步骤
1. 将证书文件上传到服务器指定路径
2. 设置正确的文件权限
3. 更新服务器配置文件
4. 重启或重载服务器
5. 验证证书是否生效

## 注意事项
- 确保证书文件权限正确（644）
- 确保私钥文件权限正确（600）
- 定期检查证书过期时间
- 配置自动续期机制

## 验证命令
\`\`\`bash
# 检查证书有效性
openssl x509 -in ${exportForm.certPath} -text -noout

# 测试HTTPS连接
curl -I https://${domain}

# 检查证书过期时间
openssl x509 -in ${exportForm.certPath} -noout -dates
\`\`\`
`
}

const copyConfig = async (content: string) => {
  try {
    await navigator.clipboard.writeText(content)
    ElMessage.success('配置已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

const downloadConfig = (format: string, content: string) => {
  const domain = props.certificate?.common_name || 'certificate'
  const extensions: Record<string, string> = {
    config: exportForm.serverType === 'nginx' ? '.conf' : '.txt',
    script: '.sh',
    docker: '.yml',
    docs: '.md'
  }
  
  const filename = `${domain}_${format}${extensions[format] || '.txt'}`
  
  const blob = new Blob([content], { type: 'text/plain' })
  const url = window.URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  link.click()
  window.URL.revokeObjectURL(url)
  
  ElMessage.success(`已下载 ${filename}`)
}

const exportAll = () => {
  if (!configPreview.value) return
  
  const domain = props.certificate?.common_name || 'certificate'
  
  Object.entries(configPreview.value).forEach(([format, content]) => {
    downloadConfig(format, content)
  })
  
  ElMessage.success('所有配置文件已下载')
}

const handleClose = () => {
  emit('update:visible', false)
  configPreview.value = null
  exportForm.formats = ['config']
  exportForm.domainConfig = ''
}

// 监听证书变化，自动填充域名配置
watch(() => props.certificate, (cert) => {
  if (cert) {
    exportForm.domainConfig = `域名: ${cert.common_name || cert.domain_name || ''}`
  }
}, { immediate: true })
</script>

<style scoped lang="scss">
.export-dialog-content {
  .cert-info,
  .export-options,
  .config-preview {
    margin-bottom: 24px;
    
    h4 {
      margin: 0 0 16px 0;
      font-size: 16px;
      font-weight: 600;
      color: #303133;
    }
  }
  
  .config-content {
    .config-textarea {
      margin-bottom: 12px;
      
      :deep(.el-textarea__inner) {
        font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
        font-size: 13px;
        line-height: 1.6;
      }
    }
    
    .config-actions {
      display: flex;
      gap: 8px;
    }
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>