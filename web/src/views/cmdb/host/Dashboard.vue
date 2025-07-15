<template>
  <div class="host-dashboard">
    <a-row :gutter="[16, 16]">
      <!-- 统计卡片 -->
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card class="stats-card">
          <div class="stats-value">{{ stats.totalHosts }}</div>
          <div class="stats-label">
            <CloudServerOutlined /> 主机总数
          </div>
          <div class="stats-trend" :class="stats.hostsTrend >= 0 ? 'up' : 'down'">
            <span v-if="stats.hostsTrend > 0">
              <ArrowUpOutlined /> 增长 {{ stats.hostsTrend }}%
            </span>
            <span v-else-if="stats.hostsTrend < 0">
              <ArrowDownOutlined /> 下降 {{ Math.abs(stats.hostsTrend) }}%
            </span>
            <span v-else>持平</span>
            <span class="trend-period">（近30天）</span>
          </div>
        </a-card>
      </a-col>
      
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card class="stats-card">
          <div class="stats-value">{{ stats.runningHosts }}</div>
          <div class="stats-label">
            <PlayCircleOutlined /> 运行中主机
          </div>
          <div class="stats-trend">
            <span>运行率 {{ stats.runningPercent }}%</span>
          </div>
        </a-card>
      </a-col>
      
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card class="stats-card">
          <div class="stats-value">{{ stats.expiringHosts }}</div>
          <div class="stats-label">
            <ClockCircleOutlined /> 30天内到期
          </div>
          <div class="stats-trend">
            <span>占比 {{ stats.expiringPercent }}%</span>
          </div>
        </a-card>
      </a-col>
      
      <a-col :xs="24" :sm="12" :lg="6">
        <a-card class="stats-card">
          <div class="stats-value">{{ stats.alertCount }}</div>
          <div class="stats-label">
            <WarningOutlined /> 告警总数
          </div>
          <div class="stats-trend" :class="stats.alertTrend >= 0 ? 'up' : 'down'">
            <span v-if="stats.alertTrend > 0">
              <ArrowUpOutlined /> 增长 {{ stats.alertTrend }}%
            </span>
            <span v-else-if="stats.alertTrend < 0">
              <ArrowDownOutlined /> 下降 {{ Math.abs(stats.alertTrend) }}%
            </span>
            <span v-else>持平</span>
            <span class="trend-period">（近7天）</span>
          </div>
        </a-card>
      </a-col>
      
      <!-- 分布图表 -->
      <a-col :xs="24" :md="12">
        <a-card title="主机分布" :bordered="false">
          <div class="chart-container">
            <div ref="providerChartRef" class="provider-chart"></div>
          </div>
        </a-card>
      </a-col>
      
      <a-col :xs="24" :md="12">
        <a-card title="地区分布" :bordered="false">
          <div class="chart-container">
            <div ref="regionChartRef" class="region-chart"></div>
          </div>
        </a-card>
      </a-col>
      
      <!-- 状态图表 -->
      <a-col :xs="24" :md="12">
        <a-card title="主机状态" :bordered="false">
          <div class="chart-container">
            <div ref="statusChartRef" class="status-chart"></div>
          </div>
        </a-card>
      </a-col>
      
      <a-col :xs="24" :md="12">
        <a-card title="系统分布" :bordered="false">
          <div class="chart-container">
            <div ref="osChartRef" class="os-chart"></div>
          </div>
        </a-card>
      </a-col>
      
      <!-- 告警和主机组 -->
      <a-col :xs="24" :lg="12">
        <alerts-widget />
      </a-col>
      
      <a-col :xs="24" :lg="12">
        <a-card title="主机组" :bordered="false">
          <a-tree
            v-if="hostGroups.length > 0"
            :tree-data="formattedHostGroups"
            :default-expanded-keys="['all']"
          >
            <template #title="{ title, key, hosts }">
              <span>{{ title }}</span>
              <a-tag color="blue" style="margin-left: 8px">
                {{ hosts || '0' }} 台主机
              </a-tag>
            </template>
          </a-tree>
          <a-empty v-else description="暂无主机组" />
        </a-card>
      </a-col>
      
      <!-- 资源使用情况 -->
      <a-col :span="24">
        <a-card title="资源使用情况" :bordered="false">
          <a-table
            :columns="resourceColumns"
            :data-source="resourceData"
            :pagination="false"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'cpu'">
                <a-progress 
                  :percent="record.cpu_usage" 
                  :status="record.cpu_usage > 90 ? 'exception' : 'normal'" 
                  size="small"
                />
              </template>
              
              <template v-if="column.key === 'memory'">
                <a-progress 
                  :percent="record.memory_usage" 
                  :status="record.memory_usage > 90 ? 'exception' : 'normal'" 
                  size="small"
                />
              </template>
              
              <template v-if="column.key === 'disk'">
                <a-progress 
                  :percent="record.disk_usage" 
                  :status="record.disk_usage > 90 ? 'exception' : 'normal'" 
                  size="small"
                />
              </template>
              
              <template v-if="column.key === 'action'">
                <a-button type="link" size="small" @click="viewHostDetail(record.id)">
                  查看
                </a-button>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, onBeforeUnmount, watch } from 'vue';
import { useRouter } from 'vue-router';
import * as echarts from 'echarts/core';
import { PieChart } from 'echarts/charts';
import { TooltipComponent, LegendComponent, TitleComponent } from 'echarts/components';
import { CanvasRenderer } from 'echarts/renderers';
import AlertsWidget from './components/AlertsWidget.vue';
import { useHostStore } from '@/store/modules/host';
import { message } from 'ant-design-vue';
import { storeToRefs } from 'pinia';
import { 
  CloudServerOutlined, 
  PlayCircleOutlined, 
  ClockCircleOutlined,
  WarningOutlined,
  ArrowUpOutlined,
  ArrowDownOutlined
} from '@ant-design/icons-vue';
import { getHostList, getHostGroupTree, getHostsOverallMetrics, getHostAlerts } from '@/api/system/host';
import type { Host, HostGroup, HostAlert } from '@/types/api/host';

echarts.use([PieChart, TooltipComponent, LegendComponent, TitleComponent, CanvasRenderer]);

const router = useRouter();
const hostStore = useHostStore();
const { hostGroupTree } = storeToRefs(hostStore);

// 图表引用
const providerChartRef = ref(null);
const regionChartRef = ref(null);
const statusChartRef = ref(null);
const osChartRef = ref(null);

// 图表实例
const providerChart = ref(null);
const regionChart = ref(null);
const statusChart = ref(null);
const osChart = ref(null);

// 统计数据
const stats = reactive({
  totalHosts: 0,
  hostsTrend: 0,
  runningHosts: 0,
  runningPercent: 0,
  expiringHosts: 0,
  expiringPercent: 0,
  alertCount: 0,
  alertTrend: 0,
});

// 主机组数据
const hostGroups = ref<HostGroup[]>([]);

// 资源使用表格
const resourceColumns = [
  { title: '主机名称', dataIndex: 'name', key: 'name' },
  { title: '状态', dataIndex: 'status', key: 'status' },
  { title: 'IP地址', dataIndex: 'ip', key: 'ip' },
  { title: 'CPU使用率', dataIndex: 'cpu_usage', key: 'cpu' },
  { title: '内存使用率', dataIndex: 'memory_usage', key: 'memory' },
  { title: '磁盘使用率', dataIndex: 'disk_usage', key: 'disk' },
  { title: '操作', key: 'action' },
];

const resourceData = ref<any[]>([]);

// 格式化主机组数据为树形结构
const formattedHostGroups = computed(() => {
  const buildTree = (groups: any[], parentId = null) => {
    return groups
      .filter(group => group.parent_id === parentId)
      .map(group => ({
        title: group.name,
        key: group.id,
        hosts: group.host_count || 0,
        children: buildTree(groups, group.id)
      }));
  };
  
  const result = buildTree(hostGroupTree.value);
  // 添加"全部"根节点
  return [{
    title: '全部主机组',
    key: 'all',
    hosts: hostGroupTree.value.reduce((total, group) => total + (group.host_count || 0), 0),
    children: result
  }];
});

// 查看主机详情
const viewHostDetail = (id: number) => {
  router.push(`/cmdb/host/detail/${id}`);
};

// 加载主机统计数据
const loadHostStats = async () => {
  try {
    // 获取主机列表
    const response = await getHostList({
      page: 1,
      page_size: 1000 // 获取足够的数据进行统计分析
    });
    
    const hosts: Host[] = response.data.data;
    const total = hosts.length;
    
    if (total === 0) {
      return;
    }
    
    // 计算运行中的主机数量
    const running = hosts.filter(host => host.status === 'running').length;
    
    // 计算即将过期的主机
    const now = new Date();
    const thirtyDaysLater = new Date();
    thirtyDaysLater.setDate(now.getDate() + 30);
    
    const expiring = hosts.filter(host => {
      if (!host.expired_at) return false;
      const expireDate = new Date(host.expired_at);
      return expireDate > now && expireDate <= thirtyDaysLater;
    }).length;
    
    // 统计数据计算
    stats.totalHosts = total;
    stats.runningHosts = running;
    stats.runningPercent = Math.round((running / total) * 100);
    stats.expiringHosts = expiring;
    stats.expiringPercent = Math.round((expiring / total) * 100);
    
    // 这里是模拟的趋势数据，实际项目中可能需要从API获取
    stats.hostsTrend = Math.floor(Math.random() * 20) - 5; // -5 到 15 的随机数
    
    // 绘制分布图表
    drawProviderChart(hosts);
    drawRegionChart(hosts);
    drawStatusChart(hosts);
    drawOSChart(hosts);
    
    // 获取资源使用率较高的主机
    loadResourceUsage();
    
  } catch (error) {
    console.error('Failed to load host statistics:', error);
    message.error('加载主机统计数据失败');
  }
};

// 加载告警统计数据
const loadAlertStats = async () => {
  try {
    const { data } = await getHostAlerts(30);
    const alerts: HostAlert[] = data;
    
    stats.alertCount = alerts.length;
    
    // 模拟告警趋势
    stats.alertTrend = Math.floor(Math.random() * 30) - 15; // -15 到 15 的随机数
  } catch (error) {
    console.error('Failed to load alert statistics:', error);
    message.error('加载告警统计数据失败');
  }
};

// 加载主机组数据
const loadHostGroups = async () => {
  if (hostGroupTree.value.length > 0) {
    hostGroups.value = hostGroupTree.value;
    return;
  }
  
  try {
    const { data } = await getHostGroupTree();
    hostGroups.value = data;
    hostStore.setHostGroupTree(data);
  } catch (error) {
    console.error('Failed to load host groups:', error);
    message.error('加载主机组数据失败');
  }
};

// 加载资源使用情况
const loadResourceUsage = async () => {
  try {
    // 实际项目中应该从API获取资源使用数据
    // 这里使用模拟数据
    const { data } = await getHostsOverallMetrics();
    
    if (data && Array.isArray(data)) {
      resourceData.value = data.map((item: any) => ({
        id: item.host_id,
        name: item.host_name,
        status: item.status,
        ip: item.ip,
        cpu_usage: Math.round(item.cpu_usage * 100),
        memory_usage: Math.round(item.memory_usage * 100),
        disk_usage: Math.round(item.disk_usage * 100)
      }));
    } else {
      // 如果API未提供数据，使用模拟数据
      generateMockResourceData();
    }
  } catch (error) {
    console.error('Failed to load resource usage data:', error);
    // 出错时使用模拟数据
    generateMockResourceData();
  }
};

// 生成模拟资源使用数据
const generateMockResourceData = () => {
  const mockData = [];
  
  // 生成10条模拟数据
  for (let i = 1; i <= 10; i++) {
    const cpuUsage = Math.floor(Math.random() * 100);
    const memoryUsage = Math.floor(Math.random() * 100);
    const diskUsage = Math.floor(Math.random() * 100);
    
    mockData.push({
      id: i,
      name: `主机-${i}`,
      status: i % 5 === 0 ? '异常' : '运行中',
      ip: `192.168.1.${i + 100}`,
      cpu_usage: cpuUsage,
      memory_usage: memoryUsage,
      disk_usage: diskUsage
    });
  }
  
  resourceData.value = mockData;
};

// 绘制云厂商分布图表
const drawProviderChart = (hosts: Host[]) => {
  if (!providerChartRef.value) return;
  
  // 按云厂商统计
  const providerStats: Record<string, number> = {};
  
  hosts.forEach(host => {
    const provider = host.provider_type || '其他';
    providerStats[provider] = (providerStats[provider] || 0) + 1;
  });
  
  // 转换为饼图数据格式
  const providerData = Object.entries(providerStats).map(([name, value]) => ({
    name,
    value
  }));
  
  // 初始化图表
  if (!providerChart.value) {
    providerChart.value = echarts.init(providerChartRef.value);
  }
  
  // 设置图表选项
  providerChart.value.setOption({
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'horizontal',
      bottom: 10,
      data: Object.keys(providerStats)
    },
    series: [
      {
        name: '云厂商分布',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: '18',
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: providerData
      }
    ]
  });
};

// 绘制地区分布图表
const drawRegionChart = (hosts: Host[]) => {
  if (!regionChartRef.value) return;
  
  // 按地区统计
  const regionStats: Record<string, number> = {};
  
  hosts.forEach(host => {
    const region = host.region || '未知地区';
    regionStats[region] = (regionStats[region] || 0) + 1;
  });
  
  // 转换为饼图数据格式
  const regionData = Object.entries(regionStats).map(([name, value]) => ({
    name,
    value
  }));
  
  // 初始化图表
  if (!regionChart.value) {
    regionChart.value = echarts.init(regionChartRef.value);
  }
  
  // 设置图表选项
  regionChart.value.setOption({
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'horizontal',
      bottom: 10,
      data: Object.keys(regionStats)
    },
    series: [
      {
        name: '地区分布',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: '18',
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: regionData
      }
    ]
  });
};

// 绘制主机状态图表
const drawStatusChart = (hosts: Host[]) => {
  if (!statusChartRef.value) return;
  
  // 按状态统计
  const statusStats: Record<string, number> = {};
  
  hosts.forEach(host => {
    const status = host.status || '未知状态';
    statusStats[status] = (statusStats[status] || 0) + 1;
  });
  
  // 转换为饼图数据格式
  const statusData = Object.entries(statusStats).map(([name, value]) => ({
    name: translateStatus(name),
    value
  }));
  
  // 初始化图表
  if (!statusChart.value) {
    statusChart.value = echarts.init(statusChartRef.value);
  }
  
  // 设置图表选项
  statusChart.value.setOption({
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'horizontal',
      bottom: 10,
      data: statusData.map(item => item.name)
    },
    series: [
      {
        name: '主机状态',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: '18',
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: statusData
      }
    ]
  });
};

// 绘制操作系统分布图表
const drawOSChart = (hosts: Host[]) => {
  if (!osChartRef.value) return;
  
  // 按操作系统统计
  const osStats: Record<string, number> = {};
  
  hosts.forEach(host => {
    let os = host.os || '未知系统';
    // 提取操作系统名称，忽略版本号
    if (os.toLowerCase().includes('windows')) {
      os = 'Windows';
    } else if (os.toLowerCase().includes('centos')) {
      os = 'CentOS';
    } else if (os.toLowerCase().includes('ubuntu')) {
      os = 'Ubuntu';
    } else if (os.toLowerCase().includes('debian')) {
      os = 'Debian';
    } else if (os.toLowerCase().includes('red hat') || os.toLowerCase().includes('redhat')) {
      os = 'Red Hat';
    } else if (os.toLowerCase().includes('linux')) {
      os = 'Linux';
    }
    
    osStats[os] = (osStats[os] || 0) + 1;
  });
  
  // 转换为饼图数据格式
  const osData = Object.entries(osStats).map(([name, value]) => ({
    name,
    value
  }));
  
  // 初始化图表
  if (!osChart.value) {
    osChart.value = echarts.init(osChartRef.value);
  }
  
  // 设置图表选项
  osChart.value.setOption({
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'horizontal',
      bottom: 10,
      data: Object.keys(osStats)
    },
    series: [
      {
        name: '操作系统分布',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: '18',
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: osData
      }
    ]
  });
};

// 状态翻译
const translateStatus = (status: string) => {
  const statusMap: Record<string, string> = {
    'running': '运行中',
    'stopped': '已停止',
    'error': '错误',
    'starting': '启动中',
    'stopping': '停止中',
    'rebooting': '重启中',
    'expired': '已过期',
    'recycled': '已回收'
  };
  
  return statusMap[status] || status;
};

// 处理窗口大小变化，重新绘制图表
const handleResize = () => {
  providerChart.value?.resize();
  regionChart.value?.resize();
  statusChart.value?.resize();
  osChart.value?.resize();
};

// 加载所有数据
const loadAllData = async () => {
  await Promise.all([
    loadHostStats(),
    loadAlertStats(),
    loadHostGroups()
  ]);
};

// 监听窗口大小变化
window.addEventListener('resize', handleResize);

// 组件挂载时加载数据
onMounted(() => {
  loadAllData();
});

// 组件卸载前移除事件监听
onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize);
  
  // 销毁图表实例
  providerChart.value?.dispose();
  regionChart.value?.dispose();
  statusChart.value?.dispose();
  osChart.value?.dispose();
});
</script>

<style lang="scss" scoped>
.host-dashboard {
  .stats-card {
    height: 100%;
    
    .stats-value {
      font-size: 28px;
      font-weight: bold;
      margin-bottom: 8px;
    }
    
    .stats-label {
      font-size: 14px;
      color: rgba(0, 0, 0, 0.65);
      margin-bottom: 16px;
      
      .anticon {
        margin-right: 4px;
      }
    }
    
    .stats-trend {
      font-size: 13px;
      
      &.up {
        color: #52c41a;
      }
      
      &.down {
        color: #f5222d;
      }
      
      .trend-period {
        color: rgba(0, 0, 0, 0.45);
        margin-left: 4px;
      }
    }
  }
  
  .chart-container {
    height: 300px;
    
    .provider-chart,
    .region-chart,
    .status-chart,
    .os-chart {
      width: 100%;
      height: 100%;
    }
  }
}
</style> 