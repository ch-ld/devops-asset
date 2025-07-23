<template>
  <div class="host-detail-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <el-button @click="goBack" text>
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button>
        <div class="host-title">
          <h1>{{ host.name }}</h1>
          <p class="instance-id">{{ host.instance_id }}</p>
        </div>
      </div>
      <div class="header-right">
        <div class="status-tags">
          <el-tag :type="getStatusTagType(host.status)" size="large" class="status-tag">
            <el-icon class="status-icon">
              <component :is="getStatusIcon(host.status)" />
            </el-icon>
            {{ getStatusText(host.status) }}
          </el-tag>
          <el-tag v-if="host.provider_type" :type="getProviderTagType(host.provider_type)" size="large">
            {{ getProviderDisplayName(host.provider_type) }}
          </el-tag>
        </div>
        <div class="action-buttons">
          <el-button @click="refreshHostData" :loading="loading">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
          <el-button type="primary" @click="editHost">
            <el-icon><Edit /></el-icon>
            编辑
          </el-button>
          <el-dropdown @command="handleMenuClick" trigger="click">
            <el-button>
              更多操作
              <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="ssh">
                  <el-icon><Monitor /></el-icon>
                  打开终端
                </el-dropdown-item>
                <el-dropdown-item command="sftp">
                  <el-icon><FolderOpened /></el-icon>
                  文件管理
                </el-dropdown-item>
                <el-dropdown-item command="sync">
                  <el-icon><Refresh /></el-icon>
                  同步状态
                </el-dropdown-item>
                <el-dropdown-item command="restart" divided>
                  <el-icon><SwitchButton /></el-icon>
                  重启主机
                </el-dropdown-item>
                <el-dropdown-item command="delete">
                  <el-icon><Delete /></el-icon>
                  删除
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="main-content" v-loading="loading">
      <el-row :gutter="24">
        <!-- 左侧信息面板 -->
        <el-col :lg="16" :md="24">
          <!-- 基本信息卡片 -->
          <el-card class="info-card" shadow="never">
            <template #header>
              <div class="card-header">
                <h3>基本信息</h3>
                <el-button size="small" text @click="editHost">
                  <el-icon><Edit /></el-icon>
                  编辑
                </el-button>
              </div>
            </template>
            
            <div class="info-grid">
              <div class="info-item">
                <label>主机名称</label>
                <span>{{ host.name }}</span>
              </div>
              <div class="info-item">
                <label>实例ID</label>
                <span class="instance-id">{{ host.instance_id }}</span>
              </div>
              <div class="info-item">
                <label>状态</label>
                <el-tag :type="getStatusTagType(host.status)">
                  <el-icon><component :is="getStatusIcon(host.status)" /></el-icon>
                  {{ getStatusText(host.status) }}
                </el-tag>
              </div>
              <div class="info-item">
                <label>云提供商</label>
                <el-tag :type="getProviderTagType(host.provider_type)">
                  {{ getProviderDisplayName(host.provider_type) }}
                </el-tag>
              </div>
              <div class="info-item">
                <label>资源类型</label>
                <span>{{ host.resource_type?.toUpperCase() || '-' }}</span>
              </div>
              <div class="info-item">
                <label>地区</label>
                <span>{{ host.region || '-' }}</span>
              </div>
              <div class="info-item">
                <label>操作系统</label>
                <span>{{ host.os || '-' }}</span>
              </div>
              <div class="info-item">
                <label>主机组</label>
                <span>{{ host.group?.name || '未分组' }}</span>
              </div>
              <div class="info-item">
                <label>创建时间</label>
                <span>{{ formatDateTime(host.created_at) }}</span>
              </div>
              <div class="info-item">
                <label>更新时间</label>
                <span>{{ formatDateTime(host.updated_at) }}</span>
              </div>
              <div v-if="host.expired_at" class="info-item">
                <label>过期时间</label>
                <span :class="{ 'expired': isExpired(host.expired_at), 'expiring': isExpiring(host.expired_at) }">
                  {{ formatDateTime(host.expired_at) }}
                  <el-tag v-if="isExpiring(host.expired_at)" type="warning" size="small">即将过期</el-tag>
                  <el-tag v-if="isExpired(host.expired_at)" type="danger" size="small">已过期</el-tag>
                </span>
              </div>
            </div>
          </el-card>

          <!-- 网络信息卡片 -->
          <el-card class="info-card" shadow="never">
            <template #header>
              <div class="card-header">
                <h3>网络信息</h3>
              </div>
            </template>
            
            <div class="network-info">
              <div v-if="host.public_ip && host.public_ip.length > 0" class="ip-section">
                <h4>公网IP地址</h4>
                <div class="ip-list">
                  <div v-for="ip in host.public_ip" :key="ip" class="ip-item">
                    <el-tag type="success" size="small">{{ ip }}</el-tag>
                    <el-button size="small" text @click="copyToClipboard(ip)">
                      <el-icon><CopyDocument /></el-icon>
                    </el-button>
                  </div>
                </div>
              </div>
              
              <div v-if="host.private_ip && host.private_ip.length > 0" class="ip-section">
                <h4>私网IP地址</h4>
                <div class="ip-list">
                  <div v-for="ip in host.private_ip" :key="ip" class="ip-item">
                    <el-tag type="info" size="small">{{ ip }}</el-tag>
                    <el-button size="small" text @click="copyToClipboard(ip)">
                      <el-icon><CopyDocument /></el-icon>
                    </el-button>
                  </div>
                </div>
              </div>
            </div>
          </el-card>

          <!-- 配置信息卡片 -->
          <el-card v-if="host.configuration" class="info-card" shadow="never">
            <template #header>
              <div class="card-header">
                <h3>配置信息</h3>
              </div>
            </template>
            
            <div class="config-info">
              <div class="config-grid">
                <div v-if="host.configuration.cpu_cores" class="config-item">
                  <div class="config-icon">
                    <el-icon><Cpu /></el-icon>
                  </div>
                  <div class="config-details">
                    <span class="config-label">CPU核心</span>
                    <span class="config-value">{{ host.configuration.cpu_cores }} 核</span>
                  </div>
                </div>
                
                <div v-if="host.configuration.memory_size" class="config-item">
                  <div class="config-icon">
                    <el-icon><MemoryCard /></el-icon>
                  </div>
                  <div class="config-details">
                    <span class="config-label">内存大小</span>
                    <span class="config-value">{{ formatMemory(host.configuration.memory_size) }}</span>
                  </div>
                </div>
                
                <div v-if="host.configuration.disk_size" class="config-item">
                  <div class="config-icon">
                    <el-icon><Coin /></el-icon>
                  </div>
                  <div class="config-details">
                    <span class="config-label">磁盘大小</span>
                    <span class="config-value">{{ host.configuration.disk_size }} GB</span>
                  </div>
                </div>
                
                <div v-if="host.configuration.instance_type" class="config-item">
                  <div class="config-icon">
                    <el-icon><Monitor /></el-icon>
                  </div>
                  <div class="config-details">
                    <span class="config-label">实例类型</span>
                    <span class="config-value">{{ host.configuration.instance_type }}</span>
                  </div>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>

        <!-- 右侧操作面板 -->
        <el-col :lg="8" :md="24">
          <!-- 快速操作卡片 -->
          <el-card class="action-card" shadow="never">
            <template #header>
              <h3>快速操作</h3>
            </template>
            
            <div class="quick-actions">
              <el-button 
                type="primary" 
                size="large" 
                class="action-btn"
                @click="handleMenuClick('ssh')"
                :disabled="!canConnect"
              >
                <el-icon><Monitor /></el-icon>
                SSH终端
              </el-button>
              
              <el-button 
                type="success" 
                size="large" 
                class="action-btn"
                @click="handleMenuClick('sftp')"
                :disabled="!canConnect"
              >
                <el-icon><FolderOpened /></el-icon>
                SFTP文件
              </el-button>
              
              <el-button 
                size="large" 
                class="action-btn"
                @click="handleMenuClick('sync')"
              >
                <el-icon><Refresh /></el-icon>
                同步状态
              </el-button>
              
              <el-button 
                type="warning" 
                size="large" 
                class="action-btn"
                @click="handleMenuClick('restart')"
                :disabled="host.status !== 'running'"
              >
                <el-icon><SwitchButton /></el-icon>
                重启主机
              </el-button>
            </div>
          </el-card>

          <!-- 监控信息卡片 -->
          <el-card class="monitor-card" shadow="never">
            <template #header>
              <h3>监控信息</h3>
            </template>
            
            <div class="monitor-info">
              <div class="monitor-item">
                <div class="monitor-label">CPU使用率</div>
                <el-progress :percentage="cpuUsage" :color="getProgressColor(cpuUsage)" />
              </div>
              
              <div class="monitor-item">
                <div class="monitor-label">内存使用率</div>
                <el-progress :percentage="memoryUsage" :color="getProgressColor(memoryUsage)" />
              </div>
              
              <div class="monitor-item">
                <div class="monitor-label">磁盘使用率</div>
                <el-progress :percentage="diskUsage" :color="getProgressColor(diskUsage)" />
              </div>
              
              <div class="monitor-item">
                <div class="monitor-label">网络状态</div>
                <el-tag :type="networkStatus === 'normal' ? 'success' : 'danger'">
                  {{ networkStatus === 'normal' ? '正常' : '异常' }}
                </el-tag>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 模态框组件 -->
    <terminal-window v-model:visible="terminalVisible" :host="host" />
    <sftp-window v-model:visible="sftpVisible" :host="host" />
  </div>
</template>

<script setup>
  import { ref, computed, onMounted } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import {
    ArrowLeft, Refresh, Edit, ArrowDown, Monitor, FolderOpened,
    SwitchButton, Delete, CircleCheck, CircleClose, Warning,
    CopyDocument, Cpu, MemoryCard, Coin
  } from '@element-plus/icons-vue'
  import dayjs from 'dayjs'

  import { getHost, syncHostStatus as apiSyncHostStatus, deleteHost as apiDeleteHost } from '@/api/system/host'
  import TerminalWindow from './components/TerminalWindow.vue'
  import SftpWindow from './components/SftpWindow.vue'

  const route = useRoute()
  const router = useRouter()

  // 响应式数据
  const loading = ref(false)
  const host = ref({})
  const terminalVisible = ref(false)
  const sftpVisible = ref(false)

  // 模拟监控数据
  const cpuUsage = ref(45)
  const memoryUsage = ref(68)
  const diskUsage = ref(32)
  const networkStatus = ref('normal')

  // 计算属性
  const canConnect = computed(() => {
    return host.value.status === 'running' &&
           (host.value.public_ip?.length > 0 || host.value.private_ip?.length > 0)
  })

  // 生命周期
  onMounted(() => {
    loadHostDetail()
  })

  // 方法
  const loadHostDetail = async () => {
    try {
      loading.value = true
      const hostId = route.params.id
      const response = await getHost(hostId)
      host.value = response.data
    } catch (error) {
      ElMessage.error('加载主机详情失败: ' + error.message)
    } finally {
      loading.value = false
    }
  }

  const refreshHostData = async () => {
    await loadHostDetail()
    ElMessage.success('数据刷新成功')
  }

  const goBack = () => {
    router.back()
  }

  const editHost = () => {
    router.push(`/cmdb/host/edit/${host.value.id}`)
  }

  const handleMenuClick = async (command) => {
    switch (command) {
      case 'ssh':
        terminalVisible.value = true
        break
      case 'sftp':
        sftpVisible.value = true
        break
      case 'sync':
        await syncHostStatus()
        break
      case 'restart':
        await restartHost()
        break
      case 'delete':
        await deleteHost()
        break
    }
  }

  const syncHostStatus = async () => {
    try {
      loading.value = true
      await apiSyncHostStatus(host.value.id)
      await loadHostDetail()
      ElMessage.success('状态同步成功')
    } catch (error) {
      ElMessage.error('状态同步失败: ' + error.message)
    } finally {
      loading.value = false
    }
  }

  const restartHost = async () => {
    ElMessageBox.confirm('确定要重启这台主机吗？', '重启主机', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      try {
        loading.value = true
        // TODO: 实现重启主机功能
        ElMessage.info('重启功能开发中...')
        // await hostApi.restartHost(host.value.id)
        // ElMessage.success('重启命令已发送')
        await loadHostDetail()
      } catch (error) {
        ElMessage.error('重启失败: ' + error.message)
      } finally {
        loading.value = false
      }
    }).catch(() => {
      // 用户取消
    })
  }

  const deleteHost = async () => {
    ElMessageBox.confirm(`确定要删除主机 "${host.value.name}" 吗？`, '删除主机', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      try {
        loading.value = true
        await apiDeleteHost(host.value.id)
        ElMessage.success('主机删除成功')
        router.push('/cmdb/host')
      } catch (error) {
        ElMessage.error('删除失败: ' + error.message)
      } finally {
        loading.value = false
      }
    }).catch(() => {
      // 用户取消
    })
  }

  const copyToClipboard = async (text) => {
    try {
      await navigator.clipboard.writeText(text)
      ElMessage.success('已复制到剪贴板')
    } catch (error) {
      ElMessage.error('复制失败')
    }
  }

  // 工具函数
  const getStatusTagType = (status) => {
    const statusMap = {
      'running': 'success',
      'stopped': 'info',
      'error': 'danger',
      'expired': 'warning'
    }
    return statusMap[status] || 'info'
  }

  const getStatusIcon = (status) => {
    const iconMap = {
      'running': 'CircleCheck',
      'stopped': 'CircleClose',
      'error': 'Warning',
      'expired': 'Warning'
    }
    return iconMap[status] || 'CircleClose'
  }

  const getStatusText = (status) => {
    const textMap = {
      'running': '运行中',
      'stopped': '已停止',
      'error': '错误',
      'expired': '已过期'
    }
    return textMap[status] || status
  }

  const getProviderTagType = (providerType) => {
    const typeMap = {
      'aliyun': 'warning',
      'tencent': 'primary',
      'aws': 'success',
      'manual': 'info'
    }
    return typeMap[providerType] || 'info'
  }

  const getProviderDisplayName = (providerType) => {
    const nameMap = {
      'aliyun': '阿里云',
      'tencent': '腾讯云',
      'aws': 'AWS',
      'manual': '自建'
    }
    return nameMap[providerType] || providerType
  }

  const formatDateTime = (dateTime) => {
    if (!dateTime) return '-'
    return dayjs(dateTime).format('YYYY-MM-DD HH:mm:ss')
  }

  const formatMemory = (memory) => {
    if (!memory) return '-'
    if (memory >= 1024) {
      return `${(memory / 1024).toFixed(1)} GB`
    }
    return `${memory} MB`
  }

  const isExpired = (expiredAt) => {
    if (!expiredAt) return false
    return dayjs(expiredAt).isBefore(dayjs())
  }

  const isExpiring = (expiredAt) => {
    if (!expiredAt) return false
    const diffDays = dayjs(expiredAt).diff(dayjs(), 'day')
    return diffDays <= 30 && diffDays >= 0
  }

  const getProgressColor = (percentage) => {
    if (percentage >= 80) return '#f56c6c'
    if (percentage >= 60) return '#e6a23c'
    return '#67c23a'
  }
</script>

<style lang="scss" scoped>
.host-detail-page {
  padding: 16px;
  background-color: #f5f5f5;
  min-height: 100vh;

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
    padding: 20px 24px;
    background: white;
    border-radius: 8px;
    box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);

    .header-left {
      display: flex;
      align-items: center;
      gap: 16px;

      .host-title {
        h1 {
          margin: 0;
          font-size: 24px;
          font-weight: 600;
          color: #1f2937;
        }

        .instance-id {
          margin: 4px 0 0 0;
          color: #6b7280;
          font-size: 14px;
          font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
        }
      }
    }

    .header-right {
      display: flex;
      align-items: center;
      gap: 16px;

      .status-tags {
        display: flex;
        gap: 8px;

        .status-tag {
          display: flex;
          align-items: center;
          gap: 4px;

          .status-icon {
            font-size: 14px;
          }
        }
      }

      .action-buttons {
        display: flex;
        gap: 8px;
      }
    }
  }

  .main-content {
    .info-card {
      margin-bottom: 24px;
      border-radius: 8px;
      border: none;

      .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;

        h3 {
          margin: 0;
          font-size: 18px;
          font-weight: 600;
          color: #1f2937;
        }
      }

      .info-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
        gap: 16px;

        .info-item {
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: 12px 0;
          border-bottom: 1px solid #f3f4f6;

          &:last-child {
            border-bottom: none;
          }

          label {
            font-weight: 500;
            color: #374151;
            min-width: 100px;
          }

          span {
            color: #6b7280;
            text-align: right;

            &.instance-id {
              font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
            }

            &.expired {
              color: #dc2626;
            }

            &.expiring {
              color: #d97706;
            }
          }
        }
      }

      .network-info {
        .ip-section {
          margin-bottom: 24px;

          &:last-child {
            margin-bottom: 0;
          }

          h4 {
            margin: 0 0 12px 0;
            font-size: 16px;
            font-weight: 600;
            color: #374151;
          }

          .ip-list {
            display: flex;
            flex-wrap: wrap;
            gap: 8px;

            .ip-item {
              display: flex;
              align-items: center;
              gap: 4px;
            }
          }
        }
      }

      .config-info {
        .config-grid {
          display: grid;
          grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
          gap: 16px;

          .config-item {
            display: flex;
            align-items: center;
            gap: 12px;
            padding: 16px;
            background: #f9fafb;
            border-radius: 8px;

            .config-icon {
              width: 40px;
              height: 40px;
              border-radius: 8px;
              background: #e5e7eb;
              display: flex;
              align-items: center;
              justify-content: center;
              color: #6b7280;
              font-size: 18px;
            }

            .config-details {
              flex: 1;

              .config-label {
                display: block;
                font-size: 12px;
                color: #6b7280;
                margin-bottom: 4px;
              }

              .config-value {
                display: block;
                font-size: 16px;
                font-weight: 600;
                color: #1f2937;
              }
            }
          }
        }
      }
    }
  }

  .action-card {
    margin-bottom: 24px;
    border-radius: 8px;
    border: none;

    .quick-actions {
      display: flex;
      flex-direction: column;
      gap: 12px;

      .action-btn {
        width: 100%;
        height: 48px;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 8px;
        font-size: 14px;
        font-weight: 500;
      }
    }
  }

  .monitor-card {
    border-radius: 8px;
    border: none;

    .monitor-info {
      .monitor-item {
        margin-bottom: 20px;

        &:last-child {
          margin-bottom: 0;
        }

        .monitor-label {
          margin-bottom: 8px;
          font-size: 14px;
          color: #374151;
          font-weight: 500;
        }
      }
    }
  }
}

// 全局样式覆盖
:deep(.el-card) {
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06);
}

:deep(.el-progress-bar__outer) {
  border-radius: 4px;
}

:deep(.el-progress-bar__inner) {
  border-radius: 4px;
}

@media (max-width: 768px) {
  .host-detail-page {
    padding: 12px;

    .page-header {
      flex-direction: column;
      align-items: flex-start;
      gap: 16px;

      .header-right {
        width: 100%;
        justify-content: space-between;

        .action-buttons {
          flex-wrap: wrap;
        }
      }
    }

    .main-content {
      .info-grid {
        grid-template-columns: 1fr;
      }

      .config-grid {
        grid-template-columns: 1fr;
      }
    }
  }
}
</style>
