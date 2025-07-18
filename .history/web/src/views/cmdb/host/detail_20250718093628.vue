<template>
  <div class="host-detail-container">
    <a-spin :spinning="loading">
      <a-page-header :title="host.name" :subtitle="host.instance_id" @back="goBack">
        <template #tags>
          <a-tag :color="getStatusColor(host.status)">{{ getStatusText(host.status) }}</a-tag>
          <a-tag v-if="host.provider_type" color="blue">{{ host.provider_type }}</a-tag>
        </template>
        <template #extra>
          <a-space>
            <a-button @click="refreshHostData">
              <template #icon><ReloadOutlined /></template>
              刷新
            </a-button>
            <a-button type="primary" @click="editHost">
              <template #icon><EditOutlined /></template>
              编辑
            </a-button>
            <a-dropdown>
              <a-button>
                更多操作
                <template #icon><DownOutlined /></template>
              </a-button>
              <template #overlay>
                <a-menu @click="handleMenuClick">
                  <a-menu-item key="ssh">
                    <template #icon><CodeOutlined /></template>
                    打开终端
                  </a-menu-item>
                  <a-menu-item key="sftp">
                    <template #icon><FolderOpenOutlined /></template>
                    文件管理
                  </a-menu-item>
                  <a-menu-item key="sync">
                    <template #icon><SyncOutlined /></template>
                    同步状态
                  </a-menu-item>
                  <a-menu-item key="restart">
                    <template #icon><PoweroffOutlined /></template>
                    重启主机
                  </a-menu-item>
                  <a-menu-item key="delete">
                    <template #icon><DeleteOutlined /></template>
                    删除
                  </a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </a-space>
        </template>
      </a-page-header>

      <a-row :gutter="16" class="detail-content">
        <!-- 主机基本信息 -->
        <a-col :span="16">
          <a-card title="基本信息" :bordered="false">
            <a-descriptions :column="2" bordered>
              <a-descriptions-item label="主机名称">{{ host.name }}</a-descriptions-item>
              <a-descriptions-item label="实例ID">{{ host.instance_id }}</a-descriptions-item>
              <a-descriptions-item label="状态">
                <a-tag :color="getStatusColor(host.status)">{{ getStatusText(host.status) }}</a-tag>
              </a-descriptions-item>
              <a-descriptions-item label="地区">{{ host.region || '未知' }}</a-descriptions-item>
              <a-descriptions-item label="操作系统">{{ host.os || '未知' }}</a-descriptions-item>
              <a-descriptions-item label="资源类型">{{
                host.resource_type || '未知'
              }}</a-descriptions-item>
              <a-descriptions-item label="CPU">
                {{ host.configuration?.cpu_cores || '未知' }} 核
              </a-descriptions-item>
              <a-descriptions-item label="内存">
                {{ formatMemorySize(host.configuration?.memory_size) }}
              </a-descriptions-item>
              <a-descriptions-item label="公网IP">
                <div v-if="host.public_ip && host.public_ip.length">
                  <a-tag v-for="(ip, index) in host.public_ip" :key="index" color="green">
                    {{ ip }}
                  </a-tag>
                </div>
                <span v-else>无</span>
              </a-descriptions-item>
              <a-descriptions-item label="内网IP">
                <div v-if="host.private_ip && host.private_ip.length">
                  <a-tag v-for="(ip, index) in host.private_ip" :key="index" color="blue">
                    {{ ip }}
                  </a-tag>
                </div>
                <span v-else>无</span>
              </a-descriptions-item>
              <a-descriptions-item label="到期时间" :span="2">
                <a-tag v-if="host.expired_at" :color="getExpiryColor(host.expired_at)">
                  {{ formatExpiryTime(host.expired_at) }}
                </a-tag>
                <span v-else>无</span>
              </a-descriptions-item>
              <a-descriptions-item label="备注" :span="2">
                {{ host.remark || '无备注' }}
              </a-descriptions-item>
            </a-descriptions>
          </a-card>

          <!-- 监控信息 -->
          <a-card title="监控信息" :bordered="false" class="metrics-card">
            <a-space direction="vertical" style="width: 100%">
              <!-- CPU使用率 -->
              <div class="metric-item">
                <div class="metric-header">
                  <span class="metric-title">CPU使用率</span>
                  <span class="metric-value">{{
                    metrics.cpu_usage ? (metrics.cpu_usage * 100).toFixed(2) + '%' : '未知'
                  }}</span>
                </div>
                <a-progress
                  :percent="metrics.cpu_usage ? metrics.cpu_usage * 100 : 0"
                  :status="metrics.cpu_usage && metrics.cpu_usage > 0.9 ? 'exception' : 'normal'"
                />
              </div>

              <!-- 内存使用率 -->
              <div class="metric-item">
                <div class="metric-header">
                  <span class="metric-title">内存使用率</span>
                  <span class="metric-value">{{
                    metrics.memory_usage ? (metrics.memory_usage * 100).toFixed(2) + '%' : '未知'
                  }}</span>
                </div>
                <a-progress
                  :percent="metrics.memory_usage ? metrics.memory_usage * 100 : 0"
                  :status="
                    metrics.memory_usage && metrics.memory_usage > 0.9 ? 'exception' : 'normal'
                  "
                />
              </div>

              <!-- 磁盘使用率 -->
              <div class="metric-item">
                <div class="metric-header">
                  <span class="metric-title">磁盘使用率</span>
                  <span class="metric-value">{{
                    metrics.disk_usage ? (metrics.disk_usage * 100).toFixed(2) + '%' : '未知'
                  }}</span>
                </div>
                <a-progress
                  :percent="metrics.disk_usage ? metrics.disk_usage * 100 : 0"
                  :status="metrics.disk_usage && metrics.disk_usage > 0.9 ? 'exception' : 'normal'"
                />
              </div>
            </a-space>
          </a-card>
        </a-col>

        <!-- 侧边信息 -->
        <a-col :span="8">
          <!-- 快速操作 -->
          <a-card title="快速操作" :bordered="false" class="action-card">
            <div class="quick-actions">
              <a-button type="primary" block @click="openTerminal">
                <template #icon><CodeOutlined /></template>
                打开终端
              </a-button>
              <a-button block @click="openSftp">
                <template #icon><FolderOpenOutlined /></template>
                文件管理
              </a-button>
              <a-button block @click="syncHostStatus">
                <template #icon><SyncOutlined /></template>
                同步状态
              </a-button>
              <a-button block @click="showMonitoringModal">
                <template #icon><DashboardOutlined /></template>
                监控详情
              </a-button>
            </div>
          </a-card>

          <!-- 云服务商信息 -->
          <a-card v-if="host.provider" title="云服务商信息" :bordered="false" class="provider-card">
            <a-descriptions :column="1" size="small">
              <a-descriptions-item label="云服务商">{{ host.provider.name }}</a-descriptions-item>
              <a-descriptions-item label="账号类型">{{ host.provider.type }}</a-descriptions-item>
              <a-descriptions-item label="地区">{{ host.provider.region }}</a-descriptions-item>
            </a-descriptions>
          </a-card>

          <!-- 标签信息 -->
          <a-card title="标签信息" :bordered="false" class="tags-card">
            <div v-if="host.tags && host.tags.length" class="host-tags">
              <a-tag v-for="(tag, index) in host.tags" :key="index" color="blue">{{ tag }}</a-tag>
            </div>
            <a-empty v-else description="暂无标签" />
            <div class="tag-actions">
              <a-button type="dashed" size="small" @click="showTagsModal"> 管理标签 </a-button>
            </div>
          </a-card>

          <!-- 主机组信息 -->
          <a-card title="主机组" :bordered="false" class="group-card">
            <div v-if="host.group">
              <div class="group-info">
                <span class="group-name">{{ host.group.name }}</span>
              </div>
              <div class="group-actions">
                <a-button type="dashed" size="small" @click="showMoveGroupModal">
                  变更分组
                </a-button>
              </div>
            </div>
            <a-empty v-else description="未分配主机组" />
          </a-card>
        </a-col>
      </a-row>
    </a-spin>

    <!-- 监控详情弹窗 -->
    <a-modal v-model:visible="monitoringModalVisible" title="监控详情" width="800px" :footer="null">
      <div class="monitoring-charts">
        <!-- 图表展示区域 -->
        <a-spin :spinning="metricsLoading">
          <div ref="cpuChartRef" class="metrics-chart"></div>
          <div ref="memoryChartRef" class="metrics-chart"></div>
          <div ref="diskChartRef" class="metrics-chart"></div>
        </a-spin>
      </div>
    </a-modal>

    <!-- 标签管理弹窗 -->
    <a-modal
      v-model:visible="tagsModalVisible"
      title="标签管理"
      @ok="saveHostTags"
      :confirm-loading="tagsSaving"
    >
      <div class="tags-editor">
        <a-tag v-for="(tag, index) in editTags" :key="index" closable @close="removeTag(tag)">
          {{ tag }}
        </a-tag>
        <a-input
          v-if="tagInputVisible"
          ref="tagInputRef"
          v-model:value="tagInputValue"
          type="text"
          size="small"
          style="width: 78px"
          @blur="handleTagInputConfirm"
          @pressEnter="handleTagInputConfirm"
        />
        <a-tag v-else style="background: #fff; border-style: dashed" @click="showTagInput">
          <plus-outlined /> 新标签
        </a-tag>
      </div>
    </a-modal>

    <!-- 移动分组弹窗 -->
    <a-modal
      v-model:visible="moveGroupModalVisible"
      title="变更分组"
      @ok="saveHostGroup"
      :confirm-loading="groupSaving"
    >
      <div class="group-selector">
        <a-tree-select
          v-model:value="selectedGroupId"
          style="width: 100%"
          :dropdown-style="{ maxHeight: '400px', overflow: 'auto' }"
          :tree-data="groupTreeData"
          placeholder="请选择主机组"
          tree-default-expand-all
          :replaceFields="{ children: 'children', key: 'id', value: 'id', title: 'name' }"
          allow-clear
        />
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
  // @ts-nocheck
  import { ref, reactive, onMounted, nextTick } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { useHostStore } from '@/store/modules/host'
  import { useECharts } from '@/utils/echarts/useECharts'
  import { formatTimestamp } from '@/utils/dataprocess/format'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import dayjs from 'dayjs'
  import {
    Edit as EditOutlined,
    Delete as DeleteOutlined,
    Refresh as ReloadOutlined,
    ArrowDown as DownOutlined,
    Monitor as CodeOutlined,
    FolderOpened as FolderOpenOutlined,
    RefreshRight as SyncOutlined,
    SwitchButton as PoweroffOutlined,
    DataAnalysis as DashboardOutlined,
    Plus as PlusOutlined
  } from '@element-plus/icons-vue'
  import type { Host, HostMetrics } from '@/types/api/host'

  const route = useRoute()
  const router = useRouter()
  const hostStore = useHostStore()

  // 获取主机ID
  const hostId = ref<number>(Number(route.params.id))
  const host = reactive<Host>({} as Host)
  const loading = ref<boolean>(true)

  // 监控指标数据
  const metrics = reactive<HostMetrics>({
    host_id: 0,
    host_name: '',
    timestamp: '',
    cpu_usage: 0,
    memory_usage: 0,
    disk_usage: 0,
    network_in: 0,
    network_out: 0,
    load_1: 0,
    load_5: 0,
    load_15: 0,
    process_count: 0,
    uptime: 0,
    status: '',
    ip: ''
  })

  // 监控详情弹窗
  const monitoringModalVisible = ref<boolean>(false)
  const metricsLoading = ref<boolean>(false)
  const cpuChartRef = ref<HTMLElement | null>(null)
  const memoryChartRef = ref<HTMLElement | null>(null)
  const diskChartRef = ref<HTMLElement | null>(null)
  let cpuChart: any = null
  let memoryChart: any = null
  let diskChart: any = null

  // 标签管理弹窗
  const tagsModalVisible = ref<boolean>(false)
  const tagsSaving = ref<boolean>(false)
  const editTags = ref<string[]>([])
  const tagInputRef = ref<HTMLElement | null>(null)
  const tagInputVisible = ref<boolean>(false)
  const tagInputValue = ref<string>('')

  // 移动分组弹窗
  const moveGroupModalVisible = ref<boolean>(false)
  const groupSaving = ref<boolean>(false)
  const selectedGroupId = ref<number | null>(null)
  const groupTreeData = ref<any[]>([])

  // 加载主机详情
  const loadHostDetail = async () => {
    loading.value = true
    try {
      const response = await hostStore.fetchHost(hostId.value)
      Object.assign(host, response)

      if (host.tags && Array.isArray(host.tags)) {
        editTags.value = [...host.tags]
      } else {
        editTags.value = []
      }

      // 加载监控指标
      await loadHostMetrics()
    } catch (error) {
      console.error('加载主机详情失败:', error)
      ElMessage.error('加载主机详情失败')
    } finally {
      loading.value = false
    }
  }

  // 加载主机监控指标
  const loadHostMetrics = async () => {
    try {
      const response = await hostStore.fetchHostMetrics(hostId.value)
      if (response) {
        Object.assign(metrics, response)
      }
    } catch (error) {
      console.error('加载监控指标失败:', error)
      // 静默处理，不影响主体功能
    }
  }

  // 状态文本和颜色
  const getStatusText = (status: string) => {
    const statusMap: Record<string, string> = {
      running: '运行中',
      stopped: '已停止',
      error: '错误',
      expired: '已过期'
    }
    return statusMap[status] || status
  }

  const getStatusColor = (status: string) => {
    const colorMap: Record<string, string> = {
      running: 'green',
      stopped: 'orange',
      error: 'red',
      expired: 'gray'
    }
    return colorMap[status] || 'blue'
  }

  // 格式化内存大小
  const formatMemorySize = (size?: number) => {
    if (!size) return '未知'
    return `${size} GB`
  }

  // 格式化到期时间和颜色
  const formatExpiryTime = (time: string) => {
    return formatTimestamp(time)
  }

  const getExpiryColor = (time: string) => {
    const now = dayjs()
    const expiry = dayjs(time)
    const days = expiry.diff(now, 'day')

    if (days < 0) return 'red'
    if (days < 30) return 'orange'
    return 'green'
  }

  // 处理菜单点击事件
  const handleMenuClick = ({ key }: { key: string }) => {
    switch (key) {
      case 'ssh':
        openTerminal()
        break
      case 'sftp':
        openSftp()
        break
      case 'sync':
        syncHostStatus()
        break
      case 'restart':
        confirmRestart()
        break
      case 'delete':
        confirmDelete()
        break
    }
  }

  // 编辑主机
  const editHost = () => {
    router.push(`/cmdb/host-edit/${hostId.value}`)
  }

  // 返回上一页
  const goBack = () => {
    router.back()
  }

  // 刷新数据
  const refreshHostData = async () => {
    await loadHostDetail()
    ElMessage.success('数据已刷新')
  }

  // 打开终端
  const openTerminal = () => {
    // 实现打开终端的逻辑，可以是新窗口或者弹窗
    ElMessage.info('打开终端功能正在开发中')
  }

  // 打开SFTP
  const openSftp = () => {
    // 实现打开SFTP的逻辑
    ElMessage.info('打开文件管理功能正在开发中')
  }

  // 同步主机状态
  const syncHostStatus = async () => {
    try {
      await hostStore.syncHostStatus(hostId.value)
      ElMessage.success('同步状态成功')
      await loadHostDetail()
    } catch (error) {
      console.error('同步状态失败:', error)
      ElMessage.error('同步状态失败')
    }
  }

  // 确认重启主机
  const confirmRestart = () => {
    ElMessageBox.confirm('确定要重启该主机吗？这可能会导致正在运行的服务中断。', '确认重启主机', {
      confirmButtonText: '确认重启',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      try {
        ElMessage.success('重启指令已发送')
        // 实际项目中应调用API重启主机
      } catch (error) {
        console.error('重启主机失败:', error)
        ElMessage.error('重启主机失败')
      }
    })
  }

  // 确认删除主机
  const confirmDelete = () => {
    ElMessageBox.confirm('确定要删除该主机吗？此操作不可恢复！', '确认删除主机', {
      confirmButtonText: '确认删除',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      try {
        await hostStore.deleteHost(hostId.value)
        ElMessage.success('删除主机成功')
        router.push('/cmdb/hosts')
      } catch (error) {
        console.error('删除主机失败:', error)
        ElMessage.error('删除主机失败')
      }
    })
  }

  // 显示监控详情弹窗
  const showMonitoringModal = async () => {
    monitoringModalVisible.value = true
    metricsLoading.value = true

    try {
      const metricsHistory = await hostStore.fetchHostMetricsHistory(hostId.value, {
        period: 'last_day'
      })

      nextTick(() => {
        initMetricsCharts(metricsHistory)
      })
    } catch (error) {
      console.error('加载监控历史数据失败:', error)
      ElMessage.error('加载监控历史数据失败')
    } finally {
      metricsLoading.value = false
    }
  }

  // 初始化监控图表
  const initMetricsCharts = (data: any) => {
    const { initChart } = useECharts()

    // CPU使用率图表
    if (cpuChartRef.value) {
      cpuChart = initChart(cpuChartRef.value)
      cpuChart.setOption({
        title: {
          text: 'CPU使用率'
        },
        tooltip: {
          trigger: 'axis'
        },
        xAxis: {
          type: 'category',
          data: data.timestamps
        },
        yAxis: {
          type: 'value',
          axisLabel: {
            formatter: '{value} %'
          },
          max: 100
        },
        series: [
          {
            name: 'CPU使用率',
            type: 'line',
            data: data.cpu.map((v: number) => (v * 100).toFixed(2)),
            areaStyle: {}
          }
        ]
      })
    }

    // 内存使用率图表
    if (memoryChartRef.value) {
      memoryChart = initChart(memoryChartRef.value)
      memoryChart.setOption({
        title: {
          text: '内存使用率'
        },
        tooltip: {
          trigger: 'axis'
        },
        xAxis: {
          type: 'category',
          data: data.timestamps
        },
        yAxis: {
          type: 'value',
          axisLabel: {
            formatter: '{value} %'
          },
          max: 100
        },
        series: [
          {
            name: '内存使用率',
            type: 'line',
            data: data.memory.map((v: number) => (v * 100).toFixed(2)),
            areaStyle: {}
          }
        ]
      })
    }

    // 磁盘使用率图表
    if (diskChartRef.value) {
      diskChart = initChart(diskChartRef.value)
      diskChart.setOption({
        title: {
          text: '磁盘使用率'
        },
        tooltip: {
          trigger: 'axis'
        },
        xAxis: {
          type: 'category',
          data: data.timestamps
        },
        yAxis: {
          type: 'value',
          axisLabel: {
            formatter: '{value} %'
          },
          max: 100
        },
        series: [
          {
            name: '磁盘使用率',
            type: 'line',
            data: data.disk.map((v: number) => (v * 100).toFixed(2)),
            areaStyle: {}
          }
        ]
      })
    }
  }

  // 显示标签管理弹窗
  const showTagsModal = () => {
    tagsModalVisible.value = true
  }

  // 显示标签输入框
  const showTagInput = () => {
    tagInputVisible.value = true
    nextTick(() => {
      tagInputRef.value?.focus()
    })
  }

  // 处理标签输入确认
  const handleTagInputConfirm = () => {
    if (tagInputValue.value && !editTags.value.includes(tagInputValue.value)) {
      editTags.value.push(tagInputValue.value)
    }
    tagInputVisible.value = false
    tagInputValue.value = ''
  }

  // 移除标签
  const removeTag = (tag: string) => {
    editTags.value = editTags.value.filter((t) => t !== tag)
  }

  // 保存标签
  const saveHostTags = async () => {
    tagsSaving.value = true
    try {
      await hostStore.batchUpdateTags([hostId.value], editTags.value, 'replace')
      host.tags = [...editTags.value]
      ElMessage.success('标签更新成功')
      tagsModalVisible.value = false
    } catch (error) {
      console.error('更新标签失败:', error)
      ElMessage.error('更新标签失败')
    } finally {
      tagsSaving.value = false
    }
  }

  // 显示移动分组弹窗
  const showMoveGroupModal = async () => {
    moveGroupModalVisible.value = true

    // 加载主机组树
    if (groupTreeData.value.length === 0) {
      try {
        await hostStore.fetchHostGroupTree()
        groupTreeData.value = hostStore.hostGroupTree
        selectedGroupId.value = host.group_id || null
      } catch (error) {
        console.error('加载主机组失败:', error)
        ElMessage.error('加载主机组失败')
      }
    }
  }

  // 保存主机组
  const saveHostGroup = async () => {
    groupSaving.value = true
    try {
      await hostStore.moveHost(hostId.value, selectedGroupId.value || undefined)
      ElMessage.success('分组移动成功')
      moveGroupModalVisible.value = false
      await loadHostDetail()
    } catch (error) {
      console.error('更新主机分组失败:', error)
      ElMessage.error('更新主机分组失败')
    } finally {
      groupSaving.value = false
    }
  }

  // 初始化
  onMounted(() => {
    if (!hostStore.fetchHost) {
      // 添加fetchHost方法到hostStore（如果不存在）
      hostStore.fetchHost = async (id: number) => {
        try {
          const response = await hostStore.getHost(id)
          return response
        } catch (error) {
          throw error
        }
      }
    }

    loadHostDetail()
  })
</script>

<style lang="scss" scoped>
  .host-detail-container {
    padding: 0;

    .detail-content {
      margin-top: 20px;

      .ant-card {
        margin-bottom: 16px;
      }
    }

    .metrics-card {
      margin-top: 16px;

      .metric-item {
        margin-bottom: 16px;

        .metric-header {
          display: flex;
          justify-content: space-between;
          margin-bottom: 8px;

          .metric-title {
            font-weight: 500;
          }

          .metric-value {
            font-weight: 600;
          }
        }
      }
    }

    .action-card {
      .quick-actions {
        display: flex;
        flex-direction: column;
        gap: 12px;
      }
    }

    .provider-card {
      margin-top: 16px;
    }

    .tags-card {
      margin-top: 16px;

      .host-tags {
        margin-bottom: 16px;
        display: flex;
        flex-wrap: wrap;
        gap: 8px;
      }

      .tag-actions {
        margin-top: 16px;
      }
    }

    .group-card {
      margin-top: 16px;

      .group-info {
        margin-bottom: 16px;

        .group-name {
          font-size: 16px;
          font-weight: 500;
        }
      }
    }

    .monitoring-charts {
      .metrics-chart {
        height: 300px;
        margin-bottom: 16px;
      }
    }

    .tags-editor {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
    }

    .group-selector {
      width: 100%;
    }
  }
</style>
