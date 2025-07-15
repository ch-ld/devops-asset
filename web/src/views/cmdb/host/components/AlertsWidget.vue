<template>
  <a-card title="主机告警" :bordered="false" class="host-alerts-widget">
    <template #extra>
      <a-space>
        <a-select 
          v-model:value="alertType" 
          style="width: 120px"
          placeholder="告警类型"
        >
          <a-select-option value="all">全部类型</a-select-option>
          <a-select-option value="expired">已过期</a-select-option>
          <a-select-option value="expiring">即将过期</a-select-option>
          <a-select-option value="error">错误</a-select-option>
          <a-select-option value="abnormal">异常状态</a-select-option>
          <a-select-option value="unreachable">无法连接</a-select-option>
        </a-select>
        <a-button type="link" size="small" @click="loadAlerts">
          <template #icon><ReloadOutlined /></template>
        </a-button>
      </a-space>
    </template>

    <a-spin :spinning="loading">
      <a-empty v-if="filteredAlerts.length === 0" description="暂无告警" />
      
      <a-list
        v-else
        :data-source="filteredAlerts"
        :pagination="{ pageSize: 5 }"
      >
        <template #renderItem="{ item }">
          <a-list-item>
            <a-list-item-meta>
              <template #title>
                <a-space>
                  <a-tag :color="getAlertTypeColor(item.alert_type)">
                    {{ getAlertTypeText(item.alert_type) }}
                  </a-tag>
                  <a @click="viewHostDetail(item.host)">{{ item.host.name }}</a>
                </a-space>
              </template>
              <template #description>
                <div class="alert-info">
                  <div class="alert-message">{{ item.message }}</div>
                  <div class="alert-time">{{ formatTime(item.time) }}</div>
                </div>
              </template>
            </a-list-item-meta>
            
            <template #actions>
              <a-button-group size="small">
                <a-button @click="viewHostDetail(item.host)">
                  <template #icon><EyeOutlined /></template>
                  查看
                </a-button>
                <a-button @click="handleAction(item)" type="primary">
                  <template #icon><ToolOutlined /></template>
                  处理
                </a-button>
              </a-button-group>
            </template>
          </a-list-item>
        </template>
      </a-list>
    </a-spin>
    
    <a-modal
      v-model:open="actionModalVisible"
      title="处理告警"
      @ok="confirmAction"
      width="640px"
    >
      <a-form :model="actionForm" layout="vertical">
        <a-form-item label="告警信息">
          <a-alert
            :message="selectedAlert?.message || ''"
            :type="getAlertTypeForAntd(selectedAlert?.alert_type)"
            show-icon
          />
        </a-form-item>
        
        <a-form-item label="处理方案">
          <a-select v-model:value="actionForm.action" style="width: 100%">
            <a-select-option v-if="selectedAlert?.alert_type === 'expired' || selectedAlert?.alert_type === 'expiring'" value="extend">
              延长过期时间
            </a-select-option>
            <a-select-option v-if="selectedAlert?.alert_type === 'abnormal' || selectedAlert?.alert_type === 'error'" value="restart">
              尝试重启主机
            </a-select-option>
            <a-select-option value="ignore">忽略此告警</a-select-option>
            <a-select-option value="manual">手动处理</a-select-option>
          </a-select>
        </a-form-item>
        
        <a-form-item v-if="actionForm.action === 'extend'" label="延长到">
          <a-date-picker v-model:value="actionForm.extendDate" style="width: 100%" />
        </a-form-item>
        
        <a-form-item label="处理备注">
          <a-textarea v-model:value="actionForm.remark" :rows="3" placeholder="请输入处理备注..." />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import type { HostAlert } from '@/types/api/host';
import { getHostAlerts, batchLifecycleHosts, batchChangeStatus } from '@/api/system/host';
import { message, notification } from 'ant-design-vue';
import { 
  ReloadOutlined,
  EyeOutlined,
  ToolOutlined
} from '@ant-design/icons-vue';
import { formatDistanceToNow } from 'date-fns';
import { zhCN } from 'date-fns/locale';
import dayjs from 'dayjs';

const router = useRouter();

// 基础状态
const loading = ref(false);
const alerts = ref<HostAlert[]>([]);
const alertType = ref('all');

// 处理告警相关
const actionModalVisible = ref(false);
const selectedAlert = ref<HostAlert | null>(null);
const actionForm = ref({
  action: 'manual',
  extendDate: null,
  remark: '',
});

// 计算属性
const filteredAlerts = computed(() => {
  if (alertType.value === 'all') {
    return alerts.value;
  }
  return alerts.value.filter(alert => alert.alert_type === alertType.value);
});

// 获取主机告警列表
const loadAlerts = async () => {
  loading.value = true;
  try {
    const { data } = await getHostAlerts(30); // 获取30天内的告警
    alerts.value = data;
  } catch (error) {
    console.error('Failed to load alerts:', error);
    message.error('获取告警信息失败');
  } finally {
    loading.value = false;
  }
};

// 告警类型颜色映射
const getAlertTypeColor = (type: string) => {
  const colorMap: Record<string, string> = {
    'expired': 'red',
    'expiring': 'orange',
    'error': 'red',
    'abnormal': 'gold',
    'unreachable': 'purple'
  };
  return colorMap[type] || 'blue';
};

// 告警类型文本映射
const getAlertTypeText = (type: string) => {
  const textMap: Record<string, string> = {
    'expired': '已过期',
    'expiring': '即将过期',
    'error': '错误',
    'abnormal': '状态异常',
    'unreachable': '无法连接'
  };
  return textMap[type] || type;
};

// 告警类型对应Ant Design Alert组件类型
const getAlertTypeForAntd = (type?: string) => {
  const typeMap: Record<string, string> = {
    'expired': 'error',
    'expiring': 'warning',
    'error': 'error',
    'abnormal': 'warning',
    'unreachable': 'error'
  };
  return typeMap[type || ''] || 'info';
};

// 格式化时间为"多久前"的形式
const formatTime = (timeStr: string) => {
  try {
    const date = new Date(timeStr);
    return formatDistanceToNow(date, { addSuffix: true, locale: zhCN });
  } catch (error) {
    return timeStr;
  }
};

// 查看主机详情
const viewHostDetail = (host: any) => {
  if (host && host.id) {
    router.push({ path: `/cmdb/host/detail/${host.id}` });
  }
};

// 处理告警
const handleAction = (alert: HostAlert) => {
  selectedAlert.value = alert;
  
  // 根据告警类型设置默认处理方式
  if (alert.alert_type === 'expired' || alert.alert_type === 'expiring') {
    actionForm.value.action = 'extend';
    // 默认延长30天
    const date = alert.host?.expired_at 
      ? new Date(alert.host.expired_at) 
      : new Date();
    date.setDate(date.getDate() + 30);
    actionForm.value.extendDate = dayjs(date);
  } else if (alert.alert_type === 'abnormal' || alert.alert_type === 'error') {
    actionForm.value.action = 'restart';
  } else {
    actionForm.value.action = 'manual';
  }
  
  actionForm.value.remark = '';
  actionModalVisible.value = true;
};

// 确认处理告警
const confirmAction = async () => {
  if (!selectedAlert.value || !selectedAlert.value.host) return;
  
  const host = selectedAlert.value.host;
  const alertType = selectedAlert.value.alert_type;
  const hostId = host.id;
  
  try {
    // 处理中状态
    actionModalVisible.value = false;
    loading.value = true;
    
    switch (actionForm.value.action) {
      case 'extend':
        // 延长过期时间
        if (!actionForm.value.extendDate) {
          message.error('请选择延长的过期时间');
          actionModalVisible.value = true;
          loading.value = false;
          return;
        }
        
        const newExpireDate = actionForm.value.extendDate.format('YYYY-MM-DD');
        await batchLifecycleHosts({
          ids: [hostId],
          expired_at: newExpireDate
        });
        
        notification.success({
          message: '处理成功',
          description: `已将主机 ${host.name} 的过期时间延长至 ${newExpireDate}`
        });
        break;
        
      case 'restart':
        // 重启主机
        await batchChangeStatus({
          ids: [hostId],
          status: 'restarting'
        });
        
        notification.success({
          message: '处理成功',
          description: `已下发重启命令至主机 ${host.name}`
        });
        break;
        
      case 'ignore':
        // 忽略告警
        notification.success({
          message: '处理成功',
          description: `已忽略主机 ${host.name} 的告警`
        });
        break;
        
      case 'manual':
        // 手动处理 - 仅记录
        notification.success({
          message: '记录已保存',
          description: `已记录对主机 ${host.name} 的手动处理`
        });
        break;
    }
    
    // 刷新告警列表
    await loadAlerts();
    
  } catch (error) {
    console.error('Failed to process alert:', error);
    notification.error({
      message: '处理失败',
      description: '处理告警时发生错误，请重试或联系系统管理员'
    });
  } finally {
    loading.value = false;
  }
};

// 组件挂载时加载数据
onMounted(loadAlerts);
</script>

<style lang="scss" scoped>
.host-alerts-widget {
  height: 100%;
  
  .alert-info {
    display: flex;
    flex-direction: column;
    
    .alert-message {
      font-size: 13px;
      margin-bottom: 4px;
      color: rgba(0, 0, 0, 0.85);
    }
    
    .alert-time {
      font-size: 12px;
      color: rgba(0, 0, 0, 0.45);
    }
  }
  
  :deep(.ant-list-item) {
    padding: 12px 0;
  }
  
  :deep(.ant-list-item-meta-title) {
    margin-bottom: 4px;
  }
  
  :deep(.ant-list-item-action) {
    margin-left: 16px;
  }
  
  :deep(.ant-empty) {
    margin: 32px 0;
  }
}
</style> 