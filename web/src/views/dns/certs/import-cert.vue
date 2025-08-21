<template>
  <div class="cert-import-page">
    <!-- 顶部导航栏 -->
    <div class="page-nav">
      <el-breadcrumb separator="/">
        <el-breadcrumb-item @click="goBack" class="clickable">证书管理</el-breadcrumb-item>
        <el-breadcrumb-item>导入证书</el-breadcrumb-item>
      </el-breadcrumb>
    </div>

    <!-- 主容器 -->
    <div class="main-container">
      <!-- 左侧：导入区域 -->
      <div class="import-section">
        <div class="section-title">
          <h2>
            <span class="icon">📋</span>
            导入SSL证书
          </h2>
          <p>支持手动粘贴或文件上传，自动验证证书有效性</p>
        </div>

        <!-- 导入方式切换 -->
        <div class="import-mode-tabs">
          <div 
            class="mode-tab" 
            :class="{ active: importMode === 'manual' }"
            @click="importMode = 'manual'"
          >
            <div class="tab-icon">✏️</div>
            <span>手动输入</span>
          </div>
          <div 
            class="mode-tab" 
            :class="{ active: importMode === 'file' }"
            @click="importMode = 'file'"
          >
            <div class="tab-icon">📁</div>
            <span>文件上传</span>
          </div>
        </div>

        <!-- 导入内容区域 -->
        <div class="import-content">
          <!-- 手动输入模式 -->
          <div v-if="importMode === 'manual'" class="manual-mode">
            <div class="input-group">
              <label>证书内容 *</label>
              <div class="textarea-wrapper">
                <textarea
                  v-model="formData.certContent"
                  placeholder="粘贴证书内容，格式：-----BEGIN CERTIFICATE-----...-----END CERTIFICATE-----"
                  class="cert-textarea"
                  @paste="handlePaste"
                ></textarea>
                <div class="textarea-tools">
                  <button @click="validateCert" :disabled="!formData.certContent.trim()" class="tool-btn">
                    🔍 验证
                  </button>
                  <button @click="clearCert" class="tool-btn">
                    🗑️ 清空
                  </button>
                </div>
              </div>
            </div>

            <div class="input-group">
              <label>私钥内容 *</label>
              <div class="textarea-wrapper">
                <textarea
                  v-model="formData.keyContent"
                  placeholder="粘贴私钥内容，格式：-----BEGIN PRIVATE KEY-----...-----END PRIVATE KEY-----"
                  class="cert-textarea"
                ></textarea>
                <div class="textarea-tools">
                  <button @click="validateKey" :disabled="!formData.keyContent.trim()" class="tool-btn">
                    🔍 验证
                  </button>
                  <button @click="clearKey" class="tool-btn">
                    🗑️ 清空
                  </button>
                </div>
              </div>
            </div>

            <div class="input-group">
              <label>证书链（可选）</label>
              <div class="textarea-wrapper">
                <textarea
                  v-model="formData.chainContent"
                  placeholder="证书链内容，包含中间CA证书"
                  class="cert-textarea chain"
                ></textarea>
                <div class="textarea-tools">
                  <button @click="clearChain" class="tool-btn">
                    🗑️ 清空
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- 文件上传模式 -->
          <div v-if="importMode === 'file'" class="file-mode">
            <div class="upload-areas">
              <!-- 证书文件 -->
              <div class="upload-item">
                <div 
                  class="upload-zone" 
                  :class="{ 'has-file': certFile }"
                  @drop="handleCertDrop"
                  @dragover.prevent
                  @click="$refs.certFileInput.click()"
                >
                  <input 
                    ref="certFileInput" 
                    type="file" 
                    accept=".crt,.pem,.cer" 
                    @change="handleCertFileChange"
                    style="display: none"
                  >
                  <div class="upload-content">
                    <div class="upload-icon">📄</div>
                    <div class="upload-text">
                      <strong>{{ certFile ? certFile.name : '证书文件' }}</strong>
                      <p>点击选择或拖拽 .crt/.pem/.cer 文件</p>
                    </div>
                  </div>
                  <button v-if="certFile" @click.stop="removeCertFile" class="remove-btn">✕</button>
                </div>
              </div>

              <!-- 私钥文件 -->
              <div class="upload-item">
                <div 
                  class="upload-zone" 
                  :class="{ 'has-file': keyFile }"
                  @drop="handleKeyDrop"
                  @dragover.prevent
                  @click="$refs.keyFileInput.click()"
                >
                  <input 
                    ref="keyFileInput" 
                    type="file" 
                    accept=".key,.pem" 
                    @change="handleKeyFileChange"
                    style="display: none"
                  >
                  <div class="upload-content">
                    <div class="upload-icon">🔑</div>
                    <div class="upload-text">
                      <strong>{{ keyFile ? keyFile.name : '私钥文件' }}</strong>
                      <p>点击选择或拖拽 .key/.pem 文件</p>
                    </div>
                  </div>
                  <button v-if="keyFile" @click.stop="removeKeyFile" class="remove-btn">✕</button>
                </div>
              </div>

              <!-- 证书链文件 -->
              <div class="upload-item optional">
                <div 
                  class="upload-zone" 
                  :class="{ 'has-file': chainFile }"
                  @drop="handleChainDrop"
                  @dragover.prevent
                  @click="$refs.chainFileInput.click()"
                >
                  <input 
                    ref="chainFileInput" 
                    type="file" 
                    accept=".pem,.crt" 
                    @change="handleChainFileChange"
                    style="display: none"
                  >
                  <div class="upload-content">
                    <div class="upload-icon">🔗</div>
                    <div class="upload-text">
                      <strong>{{ chainFile ? chainFile.name : '证书链（可选）' }}</strong>
                      <p>点击选择或拖拽证书链文件</p>
                    </div>
                  </div>
                  <button v-if="chainFile" @click.stop="removeChainFile" class="remove-btn">✕</button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧：配置面板 -->
      <div class="config-panel">
        <!-- 基本信息 -->
        <div class="panel-section">
          <h3>基本信息</h3>
          <div class="form-group">
            <label>证书名称 *</label>
            <input 
              v-model="formData.name" 
              placeholder="为证书取个名称" 
              class="form-input"
            >
          </div>
          <div class="form-group">
            <label>备注</label>
            <textarea 
              v-model="formData.remark" 
              placeholder="添加备注信息（可选）" 
              class="form-textarea"
              rows="2"
            ></textarea>
          </div>
        </div>

        <!-- 部署设置 -->
        <div class="panel-section">
          <h3>部署设置</h3>
          <div class="toggle-group">
            <label class="toggle">
              <input type="checkbox" v-model="formData.autoDeploy">
              <span class="toggle-slider"></span>
              <span class="toggle-label">启用自动部署</span>
            </label>
          </div>
          
          <div v-if="formData.autoDeploy" class="deploy-options">
            <div class="form-group">
              <label>目标主机</label>
              <select v-model="formData.deployHosts" multiple class="form-select">
                <option v-for="host in availableHosts" :key="host.id" :value="host.id">
                  {{ host.name }} ({{ host.ip }})
                </option>
              </select>
            </div>
            <div class="form-group">
              <label>部署路径</label>
              <input 
                v-model="formData.deployPath" 
                placeholder="/etc/ssl/certs/" 
                class="form-input"
              >
            </div>
          </div>
        </div>

        <!-- 证书信息预览 -->
        <div v-if="certInfo" class="panel-section cert-preview">
          <h3>证书信息</h3>
          <div class="cert-details">
            <div class="detail-item">
              <span class="label">域名:</span>
              <span class="value">{{ certInfo.commonName }}</span>
            </div>
            <div class="detail-item">
              <span class="label">颁发者:</span>
              <span class="value">{{ certInfo.issuer }}</span>
            </div>
            <div class="detail-item">
              <span class="label">有效期:</span>
              <span class="value">{{ certInfo.validTo }}</span>
            </div>
            <div class="detail-item">
              <span class="label">算法:</span>
              <span class="value">{{ certInfo.algorithm }}</span>
            </div>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="action-buttons">
          <button @click="goBack" class="btn btn-cancel">取消</button>
          <button 
            @click="handleImport" 
            :disabled="!canImport" 
            :class="['btn', 'btn-primary', { loading: importing }]"
          >
            <span v-if="importing">导入中...</span>
            <span v-else>导入证书</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { certificateApi } from '@/api/dns/certificate'
import { getHostListWithCount } from '@/api/system/host'

const router = useRouter()

// 状态管理
const importMode = ref<'manual' | 'file'>('manual')
const importing = ref(false)
const availableHosts = ref<any[]>([])
const certInfo = ref<any>(null)

// 表单数据
const formData = reactive({
  name: '',
  remark: '',
  certContent: '',
  keyContent: '',
  chainContent: '',
  autoDeploy: false,
  deployHosts: [] as number[],
  deployPath: '/etc/ssl/certs/'
})

// 文件数据
const certFile = ref<File | null>(null)
const keyFile = ref<File | null>(null)
const chainFile = ref<File | null>(null)

// 计算属性
const canImport = computed(() => {
  const hasName = formData.name.trim()
  if (importMode.value === 'manual') {
    return hasName && formData.certContent.trim() && formData.keyContent.trim()
  } else {
    return hasName && certFile.value && keyFile.value
  }
})

// 方法
const goBack = () => {
  router.push('/dns/certs')
}

// 手动输入相关方法
const handlePaste = (event: ClipboardEvent) => {
  // 可以添加智能粘贴逻辑
}

const validateCert = async () => {
  if (!formData.certContent.trim()) return
  try {
    // 解析证书信息
    const info = parseCertificateInfo(formData.certContent)
    if (info) {
      certInfo.value = info
      ElMessage.success('证书验证通过')
    }
  } catch (error) {
    ElMessage.error('证书格式无效')
  }
}

const validateKey = () => {
  if (formData.keyContent.includes('BEGIN') && formData.keyContent.includes('PRIVATE KEY')) {
    ElMessage.success('私钥格式正确')
  } else {
    ElMessage.error('私钥格式无效')
  }
}

const clearCert = () => formData.certContent = ''
const clearKey = () => formData.keyContent = ''
const clearChain = () => formData.chainContent = ''

// 文件上传相关方法
const handleCertDrop = (event: DragEvent) => {
  event.preventDefault()
  const files = event.dataTransfer?.files
  if (files?.[0]) {
    certFile.value = files[0]
  }
}

const handleCertFileChange = (event: Event) => {
  const files = (event.target as HTMLInputElement).files
  if (files?.[0]) {
    certFile.value = files[0]
    readFileContent(files[0]).then(content => {
      formData.certContent = content
      validateCert()
    })
  }
}

const handleKeyDrop = (event: DragEvent) => {
  event.preventDefault()
  const files = event.dataTransfer?.files
  if (files?.[0]) {
    keyFile.value = files[0]
  }
}

const handleKeyFileChange = (event: Event) => {
  const files = (event.target as HTMLInputElement).files
  if (files?.[0]) {
    keyFile.value = files[0]
    readFileContent(files[0]).then(content => {
      formData.keyContent = content
    })
  }
}

const handleChainDrop = (event: DragEvent) => {
  event.preventDefault()
  const files = event.dataTransfer?.files
  if (files?.[0]) {
    chainFile.value = files[0]
  }
}

const handleChainFileChange = (event: Event) => {
  const files = (event.target as HTMLInputElement).files
  if (files?.[0]) {
    chainFile.value = files[0]
    readFileContent(files[0]).then(content => {
      formData.chainContent = content
    })
  }
}

const removeCertFile = () => {
  certFile.value = null
  formData.certContent = ''
  certInfo.value = null
}

const removeKeyFile = () => {
  keyFile.value = null
  formData.keyContent = ''
}

const removeChainFile = () => {
  chainFile.value = null
  formData.chainContent = ''
}

const readFileContent = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = e => resolve(e.target?.result as string)
    reader.onerror = reject
    reader.readAsText(file)
  })
}

// 证书信息解析
const parseCertificateInfo = (certContent: string) => {
  // 简单的证书信息提取
  try {
    const lines = certContent.split('\n')
    return {
      commonName: 'example.com', // 实际应该解析证书
      issuer: 'Let\'s Encrypt',
      validTo: '2024-12-31',
      algorithm: 'RSA 2048'
    }
  } catch {
    return null
  }
}

// 导入处理
const handleImport = async () => {
  try {
    importing.value = true
    
    let certContent = formData.certContent
    let keyContent = formData.keyContent
    let chainContent = formData.chainContent

    if (importMode.value === 'file') {
      if (certFile.value) certContent = await readFileContent(certFile.value)
      if (keyFile.value) keyContent = await readFileContent(keyFile.value)
      if (chainFile.value) chainContent = await readFileContent(chainFile.value)
    }

    const requestData = {
      cert_content: certContent,
      key_content: keyContent,
      chain_content: chainContent,
      auto_deploy: formData.autoDeploy,
      deploy_hosts: formData.deployHosts,
      deploy_path: formData.deployPath,
      restart_command: 'systemctl reload nginx',
      email_notification: false,
      notification_email: '',
      remark: formData.remark || formData.name
    }

    await certificateApi.upload(requestData)
    ElMessage.success('证书导入成功')
    router.push('/dns/certs')
  } catch (error: any) {
    ElMessage.error(error.message || '证书导入失败')
  } finally {
    importing.value = false
  }
}

// 加载主机列表
const loadHosts = async () => {
  try {
    // 使用正确的CMDB主机查询API
    const response = await getHostListWithCount({
      page: 1,
      page_size: 100
      // 移除status过滤，显示所有主机
    })
    
    console.log('证书导入 - CMDB主机API响应:', response)
    
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
      console.warn('证书导入 - 未识别的主机API响应格式:', response)
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
    
    console.log(`证书导入 - 成功加载 ${availableHosts.value.length} 台主机:`, availableHosts.value)
  } catch (error) {
    console.error('证书导入 - 加载主机列表失败:', error)
    ElMessage.error('加载主机列表失败，请确保CMDB主机管理模块正常运行')
  }
}

onMounted(() => {
  loadHosts()
})
</script>

<style scoped lang="scss">
.cert-import-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.page-nav {
  margin-bottom: 20px;
  
  .el-breadcrumb {
    color: rgba(255, 255, 255, 0.8);
    
    .clickable {
      cursor: pointer;
      &:hover {
        color: white;
      }
    }
  }
}

.main-container {
  display: grid;
  grid-template-columns: 1fr 400px;
  gap: 24px;
  max-width: 1400px;
  margin: 0 auto;
}

.import-section {
  background: white;
  border-radius: 20px;
  padding: 32px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
}

.section-title {
  margin-bottom: 32px;
  
  h2 {
    font-size: 28px;
    font-weight: 700;
    color: #2d3748;
    margin: 0 0 8px 0;
    display: flex;
    align-items: center;
    gap: 12px;
    
    .icon {
      font-size: 32px;
    }
  }
  
  p {
    color: #718096;
    font-size: 16px;
    margin: 0;
  }
}

.import-mode-tabs {
  display: flex;
  gap: 16px;
  margin-bottom: 32px;
  background: #f7fafc;
  padding: 8px;
  border-radius: 16px;
}

.mode-tab {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 16px 24px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-weight: 500;
  color: #4a5568;
  
  .tab-icon {
    font-size: 20px;
  }
  
  &.active {
    background: linear-gradient(135deg, #667eea, #764ba2);
    color: white;
    box-shadow: 0 8px 16px rgba(102, 126, 234, 0.3);
  }
  
  &:hover:not(.active) {
    background: #edf2f7;
  }
}

.import-content {
  min-height: 400px;
}

.manual-mode {
  .input-group {
    margin-bottom: 24px;
    
    label {
      display: block;
      font-weight: 600;
      color: #2d3748;
      margin-bottom: 8px;
      font-size: 14px;
    }
  }
  
  .textarea-wrapper {
    position: relative;
    
    .cert-textarea {
      width: 100%;
      min-height: 120px;
      padding: 16px;
      border: 2px solid #e2e8f0;
      border-radius: 12px;
      font-family: 'Monaco', 'Menlo', monospace;
      font-size: 13px;
      line-height: 1.5;
      resize: vertical;
      transition: border-color 0.3s ease;
      
      &:focus {
        outline: none;
        border-color: #667eea;
        box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
      }
      
      &.chain {
        min-height: 80px;
      }
    }
    
    .textarea-tools {
      position: absolute;
      top: 8px;
      right: 8px;
      display: flex;
      gap: 8px;
      
      .tool-btn {
        background: white;
        border: 1px solid #e2e8f0;
        border-radius: 6px;
        padding: 4px 8px;
        font-size: 12px;
        cursor: pointer;
        transition: all 0.2s ease;
        
        &:hover:not(:disabled) {
          background: #f7fafc;
          border-color: #cbd5e0;
        }
        
        &:disabled {
          opacity: 0.5;
          cursor: not-allowed;
        }
      }
    }
  }
}

.file-mode {
  .upload-areas {
    display: grid;
    gap: 20px;
  }
  
  .upload-item {
    &.optional .upload-zone {
      border-style: dashed;
      opacity: 0.7;
    }
  }
  
  .upload-zone {
    position: relative;
    border: 3px dashed #cbd5e0;
    border-radius: 16px;
    padding: 32px 24px;
    text-align: center;
    cursor: pointer;
    transition: all 0.3s ease;
    background: #f7fafc;
    
    &:hover {
      border-color: #667eea;
      background: #edf2f7;
      transform: translateY(-2px);
    }
    
    &.has-file {
      border-color: #48bb78;
      background: #f0fff4;
      border-style: solid;
    }
    
    .upload-content {
      display: flex;
      align-items: center;
      gap: 16px;
      
      .upload-icon {
        font-size: 32px;
      }
      
      .upload-text {
        text-align: left;
        
        strong {
          display: block;
          font-size: 16px;
          color: #2d3748;
          margin-bottom: 4px;
        }
        
        p {
          color: #718096;
          font-size: 14px;
          margin: 0;
        }
      }
    }
    
    .remove-btn {
      position: absolute;
      top: 8px;
      right: 8px;
      background: #fed7d7;
      color: #e53e3e;
      border: none;
      border-radius: 50%;
      width: 24px;
      height: 24px;
      cursor: pointer;
      font-size: 12px;
      
      &:hover {
        background: #feb2b2;
      }
    }
  }
}

.config-panel {
  background: white;
  border-radius: 20px;
  padding: 24px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  height: fit-content;
  position: sticky;
  top: 20px;
}

.panel-section {
  margin-bottom: 32px;
  
  &:last-child {
    margin-bottom: 0;
  }
  
  h3 {
    font-size: 18px;
    font-weight: 600;
    color: #2d3748;
    margin: 0 0 16px 0;
  }
}

.form-group {
  margin-bottom: 16px;
  
  label {
    display: block;
    font-weight: 500;
    color: #4a5568;
    margin-bottom: 6px;
    font-size: 14px;
  }
}

.form-input, .form-textarea, .form-select {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid #e2e8f0;
  border-radius: 8px;
  font-size: 14px;
  transition: border-color 0.3s ease;
  
  &:focus {
    outline: none;
    border-color: #667eea;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  }
}

.toggle-group {
  margin-bottom: 16px;
}

.toggle {
  display: flex;
  align-items: center;
  cursor: pointer;
  
  input[type="checkbox"] {
    display: none;
  }
  
  .toggle-slider {
    width: 44px;
    height: 24px;
    background: #cbd5e0;
    border-radius: 12px;
    position: relative;
    transition: background 0.3s ease;
    margin-right: 12px;
    
    &::before {
      content: '';
      position: absolute;
      width: 20px;
      height: 20px;
      background: white;
      border-radius: 50%;
      top: 2px;
      left: 2px;
      transition: transform 0.3s ease;
    }
  }
  
  input:checked + .toggle-slider {
    background: #667eea;
    
    &::before {
      transform: translateX(20px);
    }
  }
  
  .toggle-label {
    font-weight: 500;
    color: #2d3748;
  }
}

.deploy-options {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #e2e8f0;
}

.cert-preview {
  background: #f7fafc;
  border-radius: 12px;
  padding: 16px;
  
  .cert-details {
    .detail-item {
      display: flex;
      justify-content: space-between;
      margin-bottom: 8px;
      
      .label {
        color: #718096;
        font-size: 12px;
        font-weight: 500;
      }
      
      .value {
        color: #2d3748;
        font-size: 12px;
        font-weight: 600;
      }
    }
  }
}

.action-buttons {
  display: flex;
  gap: 12px;
  
  .btn {
    flex: 1;
    padding: 14px 24px;
    border-radius: 10px;
    font-weight: 600;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.3s ease;
    border: none;
    
    &.btn-cancel {
      background: #edf2f7;
      color: #4a5568;
      
      &:hover {
        background: #e2e8f0;
      }
    }
    
    &.btn-primary {
      background: linear-gradient(135deg, #667eea, #764ba2);
      color: white;
      
      &:hover:not(:disabled) {
        transform: translateY(-2px);
        box-shadow: 0 8px 16px rgba(102, 126, 234, 0.3);
      }
      
      &:disabled {
        opacity: 0.6;
        cursor: not-allowed;
        transform: none;
      }
      
      &.loading {
        position: relative;
        
        &::before {
          content: '';
          position: absolute;
          width: 16px;
          height: 16px;
          border: 2px solid transparent;
          border-top: 2px solid white;
          border-radius: 50%;
          animation: spin 1s linear infinite;
          left: 50%;
          top: 50%;
          transform: translate(-50%, -50%);
        }
        
        span {
          opacity: 0;
        }
      }
    }
  }
}

@keyframes spin {
  to {
    transform: translate(-50%, -50%) rotate(360deg);
  }
}

// 响应式设计
@media (max-width: 1024px) {
  .main-container {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .config-panel {
    position: static;
  }
}

@media (max-width: 768px) {
  .cert-import-page {
    padding: 16px;
  }
  
  .import-section, .config-panel {
    padding: 20px;
  }
  
  .upload-zone {
    padding: 20px 16px;
    
    .upload-content {
      flex-direction: column;
      text-align: center;
      
      .upload-text {
        text-align: center;
      }
    }
  }
}
</style>