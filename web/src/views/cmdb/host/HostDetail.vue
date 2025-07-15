<template>
  <div class="host-detail-page">
    <a-row :gutter="16" class="page-header">
      <a-col :span="12">
        <a-breadcrumb>
          <a-breadcrumb-item>
            <router-link to="/cmdb/host">主机管理</router-link>
          </a-breadcrumb-item>
          <a-breadcrumb-item>主机详情</a-breadcrumb-item>
        </a-breadcrumb>
        <div class="page-title">
          <span v-if="loading">加载中...</span>
          <template v-else>
            {{ host?.name }}
            <a-tag :color="getStatusColor(host?.status)">{{ getStatusText(host?.status) }}</a-tag>
            <a-tag color="blue">{{ host?.provider_type }}</a-tag>
          </template>
        </div>
      </a-col>
      <a-col :span="12" style="text-align: right">
        <a-space>
          <a-button @click="refreshHostData">
            <template #icon><ReloadOutlined /></template>
            刷新
          </a-button>
          <a-button @click="backToList">
            <template #icon><ArrowLeftOutlined /></template>
            返回列表
          </a-button>
          <a-button type="primary" @click="openTerminal">
            <template #icon><CodeOutlined /></template>
            SSH 终端
          </a-button>
          <a-dropdown>
            <a-button type="primary">
              <template #icon><SettingOutlined /></template>
              操作
              <template #icon><DownOutlined /></template>
            </a-button>
            <template #overlay>
              <a-menu @click="handleHostAction">
                <a-menu-item key="edit">
                  <EditOutlined /> 编辑主机
                </a-menu-item>
                <a-menu-item key="sftp">
                  <FolderOpenOutlined /> SFTP 文件管理
                </a-menu-item>
                <a-menu-item key="sync">
                  <SyncOutlined /> 同步状态
                </a-menu-item>
                <a-menu-item key="move">
                  <SwapOutlined /> 移动分组
                </a-menu-item>
                <a-menu-item key="restart" :disabled="host?.status !== 'running'">
                  <ReloadOutlined /> 重启主机
                </a-menu-item>
                <a-menu-item key="stop" :disabled="host?.status !== 'running'">
                  <PauseCircleOutlined /> 停止主机
                </a-menu-item>
                <a-menu-item key="start" :disabled="host?.status !== 'stopped'">
                  <PlayCircleOutlined /> 启动主机
                </a-menu-item>
                <a-menu-item key="delete" danger>
                  <DeleteOutlined /> 删除主机
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </a-space>
      </a-col>
    </a-row>

    <!-- 加载中状态 -->
    <a-spin :spinning="loading" tip="加载中...">
      <div v-if="host">
        <!-- 信息概览 -->
        <a-row :gutter="16">
          <a-col :span="24">
            <a-card :bordered="false" class="info-card">
              <a-descriptions :column="{ xs: 1, sm: 2, md: 3, lg: 4 }" bordered>
                <a-descriptions-item label="主机ID">{{ host.id }}</a-descriptions-item>
                <a-descriptions-item label="实例ID">{{ host.instance_id || '--' }}</a-descriptions-item>
                <a-descriptions-item label="公网IP">
                  {{ host.public_ip && host.public_ip.length > 0 ? host.public_ip.join(', ') : '--' }}
                </a-descriptions-item>
                <a-descriptions-item label="内网IP">
                  {{ host.private_ip && host.private_ip.length > 0 ? host.private_ip.join(', ') : '--' }}
                </a-descriptions-item>
                <a-descriptions-item label="地区">{{ host.region || '--' }}</a-descriptions-item>
                <a-descriptions-item label="类型">{{ host.resource_type || '--' }}</a-descriptions-item>
                <a-descriptions-item label="操作系统">{{ host.os || '--' }}</a-descriptions-item>
                <a-descriptions-item label="状态">
                  <a-tag :color="getStatusColor(host.status)">{{ getStatusText(host.status) }}</a-tag>
                </a-descriptions-item>
                <a-descriptions-item label="到期时间">
                  <a-tag :color="getExpiryColor(host.expired_at)">
                    {{ formatExpiryTime(host.expired_at) }}
                  </a-tag>
                </a-descriptions-item>
                <a-descriptions-item label="所属分组">{{ host.group?.name || '--' }}</a-descriptions-item>
                <a-descriptions-item label="云提供商">{{ host.provider?.name || '--' }}</a-descriptions-item>
                <a-descriptions-item label="创建时间">{{ formatDate(host.created_at) }}</a-descriptions-item>
              </a-descriptions>
            </a-card>
          </a-col>
        </a-row>

        <!-- 详细信息和监控 -->
        <a-row :gutter="16" style="margin-top: 16px">
          <!-- 左侧：配置信息 -->
          <a-col :xs="24" :lg="12">
            <a-card title="配置信息" :bordered="false" class="config-card">
              <a-tabs>
                <a-tab-pane key="basic" tab="基本配置">
                  <a-descriptions :column="1" bordered>
                    <a-descriptions-item label="CPU">
                      {{ host.configuration?.cpu_cores ? `${host.configuration.cpu_cores} 核` : '--' }}
                    </a-descriptions-item>
                    <a-descriptions-item label="内存">
                      {{ host.configuration?.memory_size ? formatMemorySize(host.configuration.memory_size) : '--' }}
                    </a-descriptions-item>
                    <a-descriptions-item label="实例类型">
                      {{ host.configuration?.instance_type || '--' }}
                    </a-descriptions-item>
                    <a-descriptions-item label="可用区">
                      {{ host.configuration?.zone_id || '--' }}
                    </a-descriptions-item>
                    <a-descriptions-item label="VPC ID">
                      {{ host.configuration?.vpc_id || '--' }}
                    </a-descriptions-item>
                  </a-descriptions>
                </a-tab-pane>
                <a-tab-pane key="custom" tab="自定义字段">
                  <a-empty v-if="!host.extra_fields || Object.keys(host.extra_fields).length === 0" description="暂无自定义字段" />
                  <a-descriptions v-else :column="1" bordered>
                    <template v-for="(value, key) in host.extra_fields" :key="key">
                      <a-descriptions-item :label="key">{{ value }}</a-descriptions-item>
                    </template>
                  </a-descriptions>
                </a-tab-pane>
                <a-tab-pane key="tags" tab="标签">
                  <div v-if="!host.tags || host.tags.length === 0" class="empty-tags">
                    <a-empty description="暂无标签" />
                  </div>
                  <div v-else class="tag-list">
                    <a-tag v-for="tag in host.tags" :key="tag" color="blue">{{ tag }}</a-tag>
                  </div>
                </a-tab-pane>
              </a-tabs>
            </a-card>

            <!-- 主机备注 -->
            <a-card title="备注信息" :bordered="false" style="margin-top: 16px">
              <div class="remark-content">
                <div v-if="host.remark" class="remark-text">{{ host.remark }}</div>
                <a-empty v-else description="暂无备注信息" />
              </div>
              <template #extra>
                <a-button type="link" @click="editRemark">
                  <template #icon><EditOutlined /></template>
                  编辑
                </a-button>
              </template>
            </a-card>
          </a-col>

          <!-- 右侧：监控数据 -->
          <a-col :xs="24" :lg="12">
            <a-card title="资源监控" :bordered="false" class="monitor-card">
              <a-empty v-if="!hasMonitorData" description="暂无监控数据">
                <template #extra>
                  <a-button @click="fetchMonitorData">
                    <template #icon><ReloadOutlined /></template>
                    刷新
                  </a-button>
                </template>
              </a-empty>
              <div v-else>
                <div class="chart-container">
                  <div class="chart-title">CPU 使用率</div>
                  <div ref="cpuChartRef" class="chart"></div>
                </div>
                <div class="chart-container">
                  <div class="chart-title">内存使用率</div>
                  <div ref="memoryChartRef" class="chart"></div>
                </div>
                <div class="chart-container">
                  <div class="chart-title">磁盘使用率</div>
                  <div ref="diskChartRef" class="chart"></div>
                </div>
                <div class="chart-container">
                  <div class="chart-title">网络流量</div>
                  <div ref="networkChartRef" class="chart"></div>
                </div>
              </div>
            </a-card>
          </a-col>
        </a-row>

        <!-- 操作记录 -->
        <a-row :gutter="16" style="margin-top: 16px">
          <a-col :span="24">
            <a-card title="操作记录" :bordered="false">
              <a-table
                :columns="operationColumns"
                :data-source="operationLogs"
                :pagination="{ pageSize: 10 }"
                :loading="logsLoading"
              >
                <template #bodyCell="{ column, record }">
                  <template v-if="column.key === 'operation_type'">
                    <a-tag :color="getOperationTypeColor(record.operation_type)">
                      {{ getOperationTypeText(record.operation_type) }}
                    </a-tag>
                  </template>
                  <template v-else-if="column.key === 'status'">
                    <a-tag :color="record.status ? 'green' : 'red'">
                      {{ record.status ? '成功' : '失败' }}
                    </a-tag>
                  </template>
                </template>
              </a-table>
            </a-card>
          </a-col>
        </a-row>
      </div>
      <a-empty v-else-if="!loading" description="未找到主机信息">
        <template #extra>
          <a-button type="primary" @click="backToList">返回主机列表</a-button>
        </template>
      </a-empty>
    </a-spin>

    <!-- 编辑备注对话框 -->
    <a-modal
      v-model:visible="remarkModalVisible"
      title="编辑备注信息"
      @ok="saveRemark"
    >
      <a-form :model="remarkForm">
        <a-form-item label="备注信息">
          <a-textarea v-model:value="remarkForm.remark" :rows="5" placeholder="请输入备注信息..." />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- SSH 终端窗口 -->
    <terminal-window
      v-model:visible="terminalVisible"
      :host="host"
    />

    <!-- SFTP 文件管理窗口 -->
    <sftp-window
      v-model:visible="sftpVisible"
      :host="host"
    />

    <!-- 主机编辑模态框 -->
    <host-modal
      v-model:visible="hostModalVisible"
      :host="host"
      :is-edit="true"
      @success="refreshHostData"
    />

    <!-- 移动分组模态框 -->
    <batch-move-modal
      v-model:visible="moveModalVisible"
      :host-ids="host ? [host.id] : []"
      @success="refreshHostData"
    />
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted, nextTick, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { message, Modal } from 'ant-design-vue';
import * as hostApi from '@/api/system/host';
import { useHostStore } from '@/store/modules/host';
import type { Host } from '@/types/api/host';
import { formatDistanceToNow } from 'date-fns';
import { zhCN } from 'date-fns/locale';
import dayjs from 'dayjs';
import * as echarts from 'echarts/core';
import { LineChart } from 'echarts/charts';
import { TooltipComponent, LegendComponent, TitleComponent, GridComponent } from 'echarts/components';
import { CanvasRenderer } from 'echarts/renderers';
import TerminalWindow from './components/TerminalWindow.vue';
import SftpWindow from './components/SftpWindow.vue';
import HostModal from './components/HostModal.vue';
import BatchMoveModal from './components/BatchMoveModal.vue';

import {
  ReloadOutlined,
  ArrowLeftOutlined,
  CodeOutlined,
  SettingOutlined,
  EditOutlined,
  DeleteOutlined,
  FolderOpenOutlined,
  SyncOutlined,
  SwapOutlined,
  DownOutlined,
  PauseCircleOutlined,
  PlayCircleOutlined
} from '@ant-design/icons-vue';

// 初始化 ECharts
echarts.use([LineChart, TooltipComponent, LegendComponent, TitleComponent, GridComponent, CanvasRenderer]);

const route = useRoute();
const router = useRouter();
const hostStore = useHostStore();

// 状态变量
const host = ref<Host | null>(null);
const loading = ref(true);
const logsLoading = ref(false);
const operationLogs = ref<any[]>([]);
const hasMonitorData = ref(false);

// 图表引用
const cpuChartRef = ref(null);
const memoryChartRef = ref(null);
const diskChartRef = ref(null);
const networkChartRef = ref(null);
const cpuChart = ref<echarts.ECharts | null>(null);
const memoryChart = ref<echarts.ECharts | null>(null);
const diskChart = ref<echarts.ECharts | null>(null);
const networkChart = ref<echarts.ECharts | null>(null);

// 操作记录列
const operationColumns = [
  {
    title: '操作类型',
    dataIndex: 'operation_type',
    key: 'operation_type',
  },
  {
    title: '操作人',
    dataIndex: 'operator',
    key: 'operator',
  },
  {
    title: '操作时间',
    dataIndex: 'created_at',
    key: 'created_at',
    sorter: true,
    defaultSortOrder: 'descend',
    render: (text: string) => formatDate(text),
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
  },
  {
    title: '详情',
    dataIndex: 'detail',
    key: 'detail',
    ellipsis: true,
  }
];

// 模态框控制
const terminalVisible = ref(false);
const sftpVisible = ref(false);
const hostModalVisible = ref(false);
const moveModalVisible = ref(false);
const remarkModalVisible = ref(false);
const remarkForm = reactive({
  remark: '',
});

// 生命周期钩子
onMounted(async () => {
  const hostId = parseInt(route.params.id as string);
  if (!hostId || isNaN(hostId)) {
    message.error('无效的主机ID');
    backToList();
    return;
  }
  
  await fetchHostData(hostId);
  fetchOperationLogs(hostId);
});

// 获取主机详情
const fetchHostData = async (hostId: number) => {
  loading.value = true;
  try {
    const { data } = await hostApi.getHost(hostId);
    host.value = data;
    remarkForm.remark = data.remark || '';
    
    // 初始化图表
    nextTick(() => {
      initCharts();
      fetchMonitorData();
    });
  } catch (error) {
    message.error('获取主机详情失败');
    console.error(error);
  } finally {
    loading.value = false;
  }
};

// 刷新主机数据
const refreshHostData = async () => {
  if (host.value) {
    await fetchHostData(host.value.id);
  }
};

// 返回列表
const backToList = () => {
  router.push('/cmdb/host');
};

// 获取操作记录
const fetchOperationLogs = async (hostId: number) => {
  logsLoading.value = true;
  try {
    // 假设有获取主机操作日志的API
    // const { data } = await hostApi.getHostHistory(hostId);
    // operationLogs.value = data;
    
    // 模拟数据
    operationLogs.value = [
      {
        id: 1,
        operation_type: 'create',
        operator: 'admin',
        created_at: '2023-07-01T10:00:00Z',
        status: true,
        detail: '创建主机'
      },
      {
        id: 2,
        operation_type: 'update',
        operator: 'admin',
        created_at: '2023-07-02T11:30:00Z',
        status: true,
        detail: '修改主机配置'
      },
      {
        id: 3,
        operation_type: 'restart',
        operator: 'operator',
        created_at: '2023-07-03T14:20:00Z',
        status: true,
        detail: '重启主机'
      }
    ];
  } catch (error) {
    console.error('获取操作记录失败:', error);
  } finally {
    logsLoading.value = false;
  }
};

// 初始化图表
const initCharts = () => {
  // CPU 图表
  if (cpuChartRef.value) {
    cpuChart.value = echarts.init(cpuChartRef.value);
    const cpuOption = getChartOption('CPU 使用率', '%');
    cpuChart.value.setOption(cpuOption);
  }
  
  // 内存图表
  if (memoryChartRef.value) {
    memoryChart.value = echarts.init(memoryChartRef.value);
    const memoryOption = getChartOption('内存使用率', '%');
    memoryChart.value.setOption(memoryOption);
  }
  
  // 磁盘图表
  if (diskChartRef.value) {
    diskChart.value = echarts.init(diskChartRef.value);
    const diskOption = getChartOption('磁盘使用率', '%');
    diskChart.value.setOption(diskOption);
  }
  
  // 网络图表
  if (networkChartRef.value) {
    networkChart.value = echarts.init(networkChartRef.value);
    const networkOption = getNetworkChartOption();
    networkChart.value.setOption(networkOption);
  }
  
  // 窗口大小变化时重新调整图表大小
  window.addEventListener('resize', () => {
    cpuChart.value?.resize();
    memoryChart.value?.resize();
    diskChart.value?.resize();
    networkChart.value?.resize();
  });
};

// 修改 fetchMonitorData 方法，使用真实的API获取监控数据
const fetchMonitorData = async () => {
  try {
    loading.value = true;
    
    if (!host.value) return;
    
    // 调用获取最新指标的API
    const { data: latestMetrics } = await hostApi.getHostMetrics(host.value.id);
    
    // 调用获取历史指标的API
    const { data: historyMetrics } = await hostApi.getHostMetricsHistory(
      host.value.id, 
      { period: 'last_day' }
    );
    
    // 设置有监控数据的标志
    hasMonitorData.value = true;
    
    // 处理历史数据，用于图表展示
    processMetricsHistory(historyMetrics);
    
  } catch (error) {
    console.error('获取监控数据失败:', error);
    message.error('获取监控数据失败');
    
    // 如果API失败，使用模拟数据
    generateMockData();
  } finally {
    loading.value = false;
  }
};

// 处理历史监控数据
const processMetricsHistory = (metricsHistory: any[]) => {
  if (!metricsHistory || metricsHistory.length === 0) {
    generateMockData();
    return;
  }
  
  // 提取时间轴
  const times = metricsHistory.map(m => {
    const date = new Date(m.timestamp);
    return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' });
  });
  
  // 提取各指标数据
  const cpuData = metricsHistory.map(m => m.cpu_usage);
  const memoryData = metricsHistory.map(m => m.memory_usage);
  const diskData = metricsHistory.map(m => m.disk_usage);
  const rxData = metricsHistory.map(m => m.net_io?.rx_bytes / (1024 * 1024) || 0); // 转为MB
  const txData = metricsHistory.map(m => m.net_io?.tx_bytes / (1024 * 1024) || 0); // 转为MB
  
  // 更新图表
  updateChart(cpuChart.value, times, [{ name: 'CPU', data: cpuData }]);
  updateChart(memoryChart.value, times, [{ name: '内存', data: memoryData }]);
  updateChart(diskChart.value, times, [{ name: '磁盘', data: diskData }]);
  updateNetworkChart(networkChart.value, times, rxData, txData);
};

// 生成模拟数据（仅在API调用失败时使用）
const generateMockData = () => {
  hasMonitorData.value = true;
  
  // 生成模拟数据
  const now = new Date();
  const times = Array.from({ length: 24 }, (_, i) => {
    const time = new Date(now);
    time.setHours(now.getHours() - 24 + i);
    return time.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' });
  });
  
  // CPU 数据
  const cpuData = Array.from({ length: 24 }, () => Math.floor(Math.random() * 100));
  updateChart(cpuChart.value, times, [{ name: 'CPU', data: cpuData }]);
  
  // 内存数据
  const memoryData = Array.from({ length: 24 }, () => Math.floor(Math.random() * 100));
  updateChart(memoryChart.value, times, [{ name: '内存', data: memoryData }]);
  
  // 磁盘数据
  const diskData = Array.from({ length: 24 }, () => Math.floor(50 + Math.random() * 30));
  updateChart(diskChart.value, times, [{ name: '磁盘', data: diskData }]);
  
  // 网络数据
  const inboundData = Array.from({ length: 24 }, () => Math.floor(Math.random() * 100));
  const outboundData = Array.from({ length: 24 }, () => Math.floor(Math.random() * 80));
  updateNetworkChart(networkChart.value, times, inboundData, outboundData);
};

// 更新图表数据
const updateChart = (chart: echarts.ECharts | null, xAxisData: string[], seriesData: { name: string; data: number[] }[]) => {
  if (!chart) return;
  
  chart.setOption({
    xAxis: { data: xAxisData },
    series: seriesData.map(item => ({
      name: item.name,
      data: item.data
    }))
  });
};

// 更新网络图表
const updateNetworkChart = (chart: echarts.ECharts | null, xAxisData: string[], inboundData: number[], outboundData: number[]) => {
  if (!chart) return;
  
  chart.setOption({
    xAxis: { data: xAxisData },
    series: [
      {
        name: '入站',
        data: inboundData
      },
      {
        name: '出站',
        data: outboundData
      }
    ]
  });
};

// 获取图表通用配置
const getChartOption = (name: string, unit: string) => {
  return {
    tooltip: {
      trigger: 'axis',
      formatter: function(params: any) {
        const param = params[0];
        return `${param.name}<br/>${param.seriesName}: ${param.value}${unit}`;
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      top: '10%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: [],
      axisLabel: {
        formatter: '{value}'
      }
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        formatter: `{value}${unit}`
      }
    },
    series: [
      {
        name: name,
        type: 'line',
        smooth: true,
        data: [],
        areaStyle: {},
        markPoint: {
          data: [
            { type: 'max', name: '最大值' },
            { type: 'min', name: '最小值' }
          ]
        }
      }
    ]
  };
};

// 获取网络图表配置
const getNetworkChartOption = () => {
  return {
    tooltip: {
      trigger: 'axis'
    },
    legend: {
      data: ['入站', '出站'],
      top: 0
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      top: '40px',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: []
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        formatter: '{value} Mbps'
      }
    },
    series: [
      {
        name: '入站',
        type: 'line',
        stack: 'Total',
        data: [],
        areaStyle: {},
        emphasis: {
          focus: 'series'
        }
      },
      {
        name: '出站',
        type: 'line',
        stack: 'Total',
        data: [],
        areaStyle: {},
        emphasis: {
          focus: 'series'
        }
      }
    ]
  };
};

// 主机操作
const handleHostAction = ({ key }: { key: string }) => {
  switch (key) {
    case 'edit':
      hostModalVisible.value = true;
      break;
    case 'sftp':
      sftpVisible.value = true;
      break;
    case 'sync':
      syncHostStatus();
      break;
    case 'move':
      moveModalVisible.value = true;
      break;
    case 'restart':
      confirmRestart();
      break;
    case 'stop':
      confirmStop();
      break;
    case 'start':
      confirmStart();
      break;
    case 'delete':
      confirmDelete();
      break;
  }
};

// 同步主机状态
const syncHostStatus = async () => {
  if (!host.value) return;
  
  try {
    message.loading('正在同步主机状态...');
    await hostStore.syncHostStatus(host.value.id);
    message.success('同步主机状态成功');
    refreshHostData();
  } catch (error) {
    message.error('同步主机状态失败');
    console.error(error);
  }
};

// 确认重启主机
const confirmRestart = () => {
  if (!host.value) return;
  
  Modal.confirm({
    title: '重启主机',
    content: `确定要重启主机 ${host.value.name} 吗？`,
    onOk: async () => {
      try {
        // 假设有重启主机的API
        message.loading('正在重启主机...');
        // await hostApi.restartHost(host.value.id);
        message.success('已发送重启命令');
        refreshHostData();
      } catch (error) {
        message.error('重启主机失败');
        console.error(error);
      }
    }
  });
};

// 确认停止主机
const confirmStop = () => {
  if (!host.value) return;
  
  Modal.confirm({
    title: '停止主机',
    content: `确定要停止主机 ${host.value.name} 吗？`,
    onOk: async () => {
      try {
        // 假设有停止主机的API
        message.loading('正在停止主机...');
        // await hostApi.stopHost(host.value.id);
        message.success('已发送停止命令');
        refreshHostData();
      } catch (error) {
        message.error('停止主机失败');
        console.error(error);
      }
    }
  });
};

// 确认启动主机
const confirmStart = () => {
  if (!host.value) return;
  
  Modal.confirm({
    title: '启动主机',
    content: `确定要启动主机 ${host.value.name} 吗？`,
    onOk: async () => {
      try {
        // 假设有启动主机的API
        message.loading('正在启动主机...');
        // await hostApi.startHost(host.value.id);
        message.success('已发送启动命令');
        refreshHostData();
      } catch (error) {
        message.error('启动主机失败');
        console.error(error);
      }
    }
  });
};

// 确认删除主机
const confirmDelete = () => {
  if (!host.value) return;
  
  Modal.confirm({
    title: '删除主机',
    content: `确定要删除主机 ${host.value.name} 吗？此操作不可恢复。`,
    onOk: async () => {
      try {
        await hostStore.deleteHost(host.value!.id);
        message.success('删除主机成功');
        backToList();
      } catch (error) {
        message.error('删除主机失败');
        console.error(error);
      }
    }
  });
};

// 打开终端
const openTerminal = () => {
  terminalVisible.value = true;
};

// 编辑备注
const editRemark = () => {
  remarkModalVisible.value = true;
};

// 保存备注
const saveRemark = async () => {
  if (!host.value) return;
  
  try {
    message.loading('正在保存备注...');
    await hostApi.updateHost(host.value.id, { remark: remarkForm.remark });
    message.success('保存备注成功');
    remarkModalVisible.value = false;
    refreshHostData();
  } catch (error) {
    message.error('保存备注失败');
    console.error(error);
  }
};

// 工具函数
const formatMemorySize = (size?: number) => {
  if (!size) return '--';
  return size >= 1024 ? `${(size / 1024).toFixed(1)} GB` : `${size} MB`;
};

const getStatusColor = (status?: string) => {
  if (!status) return 'default';
  const colorMap: Record<string, string> = {
    'running': 'green',
    'stopped': 'orange',
    'error': 'red',
    'expired': 'red',
    'starting': 'blue',
    'stopping': 'orange',
    'rebooting': 'blue',
  };
  return colorMap[status.toLowerCase()] || 'default';
};

const getStatusText = (status?: string) => {
  if (!status) return '--';
  const textMap: Record<string, string> = {
    'running': '运行中',
    'stopped': '已停止',
    'error': '错误',
    'expired': '已过期',
    'starting': '启动中',
    'stopping': '停止中',
    'rebooting': '重启中',
  };
  return textMap[status.toLowerCase()] || status;
};

const getExpiryColor = (expiryDate?: string) => {
  if (!expiryDate) return 'default';
  
  const now = new Date();
  const expiry = new Date(expiryDate);
  const diffDays = Math.floor((expiry.getTime() - now.getTime()) / (1000 * 60 * 60 * 24));
  
  if (diffDays < 0) return 'red';
  if (diffDays < 7) return 'orange';
  if (diffDays < 30) return 'gold';
  return 'green';
};

const formatExpiryTime = (expiryDate?: string) => {
  if (!expiryDate) return '--';
  try {
    const date = new Date(expiryDate);
    const now = new Date();
    
    if (date < now) {
      return '已过期';
    }
    
    return formatDistanceToNow(date, { addSuffix: true, locale: zhCN });
  } catch (error) {
    return expiryDate;
  }
};

const formatDate = (dateStr?: string) => {
  if (!dateStr) return '--';
  try {
    return dayjs(dateStr).format('YYYY-MM-DD HH:mm:ss');
  } catch (error) {
    return dateStr;
  }
};

const getOperationTypeColor = (type?: string) => {
  if (!type) return 'default';
  const colorMap: Record<string, string> = {
    'create': 'green',
    'update': 'blue',
    'delete': 'red',
    'start': 'green',
    'stop': 'orange',
    'restart': 'purple',
    'sync': 'cyan',
  };
  return colorMap[type.toLowerCase()] || 'default';
};

const getOperationTypeText = (type?: string) => {
  if (!type) return '--';
  const textMap: Record<string, string> = {
    'create': '创建',
    'update': '更新',
    'delete': '删除',
    'start': '启动',
    'stop': '停止',
    'restart': '重启',
    'sync': '同步',
  };
  return textMap[type.toLowerCase()] || type;
};
</script>

<style lang="scss" scoped>
.host-detail-page {
  .page-header {
    margin-bottom: 16px;
    
    .page-title {
      font-size: 20px;
      font-weight: 500;
      margin-top: 8px;
      
      .ant-tag {
        margin-left: 8px;
      }
    }
  }
  
  .info-card {
    .ant-descriptions {
      margin-bottom: 0;
    }
  }
  
  .config-card {
    height: 100%;
    
    .empty-tags {
      padding: 20px 0;
    }
    
    .tag-list {
      padding: 8px 0;
      
      .ant-tag {
        margin-bottom: 8px;
      }
    }
  }
  
  .monitor-card {
    height: 100%;
  }
  
  .chart-container {
    margin-bottom: 20px;
    
    .chart-title {
      font-size: 14px;
      font-weight: 500;
      margin-bottom: 8px;
    }
    
    .chart {
      height: 180px;
    }
  }
  
  .remark-content {
    min-height: 60px;
    
    .remark-text {
      white-space: pre-line;
    }
  }
}

// 暗色主题适配
html.dark {
  .host-detail-page {
    .page-header .page-title {
      color: rgba(255, 255, 255, 0.85);
    }
    
    .chart-container .chart-title {
      color: rgba(255, 255, 255, 0.85);
    }
  }
}
</style> 