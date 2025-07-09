<template>
  <div class="host-management-container">
    <!-- Left Panel: Provider List -->
    <div class="left-panel">
      <div class="panel-header">
        <h3>分组列表</h3>
        <a-switch v-model:checked="isTreeMode" checked-children="树" un-checked-children="列表" />
      </div>
      <div class="provider-list">
        <a-menu mode="inline" :selected-keys="[selectedProviderId]" @select="onProviderSelect">
          <a-menu-item key="all" class="provider-item">
            <span>全部主机</span>
          </a-menu-item>
          <a-sub-menu v-for="provider in providerList" :key="provider.id" :title="provider.name">
            <template #title>
              <div class="provider-item">
                <span>{{ provider.name }}</span>
                <span class="actions">
                  <a-tooltip title="同步">
                    <sync-outlined :spin="isLoading && syncingProviderId === provider.id" @click.stop="handleSyncProvider(provider.id)" />
                  </a-tooltip>
                  <a-tooltip title="编辑">
                    <edit-outlined @click.stop="handleEditProvider(provider)" />
                  </a-tooltip>
                  <a-popconfirm
                    title="确定删除此账号吗? 所有关联主机将一并删除！"
                    @confirm.stop="handleDeleteProvider(provider.id)"
                  >
                    <a-tooltip title="删除">
                      <delete-outlined />
                    </a-tooltip>
                  </a-popconfirm>
                </span>
              </div>
            </template>
            <a-menu-item :key="provider.id">
              {{ provider.name }} ({{ getHostCountForProvider(provider.id) }})
            </a-menu-item>
          </a-sub-menu>
        </a-menu>
      </div>
    </div>

    <!-- Right Panel: Host Table -->
    <div class="right-panel">
      <div class="toolbar">
        <a-input-search placeholder="输入名称/IP搜索" style="width: 240px" />
        <a-space>
          <a-button type="primary" @click="handleCreateHost">新建主机</a-button>
        </a-space>
      </div>

      <a-table
        :columns="hostColumns"
        :data-source="filteredHosts"
        :loading="isLoading"
        row-key="id"
        class="host-table"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <a>{{ record.name }}</a>
          </template>
          <template v-else-if="column.key === 'ip_address'">
            <div v-for="ip in record.public_ip" :key="ip">{{ ip }} (公)</div>
            <div v-for="ip in record.private_ip" :key="ip">{{ ip }} (内)</div>
          </template>
          <template v-else-if="column.key === 'configuration'">
            <span>{{ record.configuration?.cpu }}核 {{ record.configuration?.memory }}</span>
          </template>
          <template v-else-if="column.key === 'status'">
            <a-tag :color="record.status === 'Running' ? 'green' : 'volcano'">{{ record.status }}</a-tag>
          </template>
          <template v-else-if="column.key === 'action'">
            <a-space>
              <a @click="handleEditHost(record)">编辑</a>
              <a @click="handleSsh(record)">终端</a>
              <a @click="openSftpWindow(record.id)">SFTP</a>
              <a-popconfirm title="确定删除此主机吗?" @confirm="handleDeleteHost(record.id)">
                <a style="color: red">删除</a>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </div>
  </div>

  <ProviderModal
    v-model:visible="isProviderModalVisible"
    :provider="editingProvider"
    :is-loading="isSubmitting"
    @submit="handleProviderSubmit"
  />

  <HostModal ref="hostModalRef" :providers="providerList" @submit="handleHostSubmit" />

  <TerminalWindow
    v-if="isTerminalVisible"
    v-model:visible="isTerminalVisible"
    :host="connectingHost"
  />

  <SftpWindow ref="sftpWindowRef" />
</template>

<script setup lang="ts">
import { ref, onMounted, computed, h } from 'vue';
import { useHostStore } from '@/store/modules/host';
import { storeToRefs } from 'pinia';
import type { Host, Provider } from '@/types/api/host';
import ProviderModal from './components/ProviderModal.vue';
import HostModal from './components/HostModal.vue'; // Import HostModal
import TerminalWindow from './components/TerminalWindow.vue';
import SftpWindow from './components/SftpWindow.vue';
import { EditOutlined, DeleteOutlined, SyncOutlined } from '@ant-design/icons-vue';

const isTreeMode = ref(false);
const isProviderModalVisible = ref(false);
const editingProvider = ref<Provider | null>(null);
const syncingProviderId = ref<number | null>(null);

const hostModalRef = ref();
const sftpWindowRef = ref();

const isTerminalVisible = ref(false);
const connectingHost = ref<Host | null>(null);

const hostStore = useHostStore();
const { providerList, hostList, isLoading, isSubmitting } = storeToRefs(hostStore);

const selectedProviderId = ref<string | number>('all');

onMounted(() => {
  hostStore.fetchProviders();
  hostStore.fetchHosts();
});

const getHostCountForProvider = (providerId: number) => {
  return hostList.value.filter(h => h.provider_id === providerId).length;
}

const onProviderSelect = ({ key }: { key: string | number }) => {
  selectedProviderId.value = key;
}

const filteredHosts = computed(() => {
  if (selectedProviderId.value === 'all') {
    return hostList.value;
  }
  return hostList.value.filter(h => h.provider_id === selectedProviderId.value);
})

const hostColumns = [
  { title: '主机名称', dataIndex: 'name', key: 'name' },
  { title: 'IP地址', key: 'ip_address' },
  { title: '配置信息', key: 'configuration' },
  { title: '状态', key: 'status' },
  { title: '操作', key: 'action', width: 220 },
];

const handleCreate = () => {
  editingProvider.value = null;
  isProviderModalVisible.value = true;
};

const handleEditProvider = (provider: Provider) => {
  editingProvider.value = provider;
  isProviderModalVisible.value = true;
};

const handleDeleteProvider = async (id: number) => {
  // TODO: Add a confirmation dialog
  await hostStore.deleteProvider(id);
};

const handleProviderSubmit = async (providerData: Partial<Provider>) => {
  isSubmitting.value = true;
  try {
    if (editingProvider.value) {
      // Update
      await hostStore.updateProvider(editingProvider.value.id, providerData);
    } else {
      // Create
      await hostStore.addProvider(providerData);
    }
    isProviderModalVisible.value = false;
  } catch (error) {
    console.error("Submission failed:", error);
    // Optionally show an error message to the user
  } finally {
    isSubmitting.value = false;
  }
};


const handleDelete = (id: number) => {
  console.log('Delete:', id);
  // TODO: Implement delete logic for hosts
};

const handleSsh = (record: Host) => {
  connectingHost.value = record;
  isTerminalVisible.value = true;
};

const handleSyncProvider = async (id: number) => {
  syncingProviderId.value = id;
  await hostStore.syncProviderResources(id);
  syncingProviderId.value = null;
}

function openSftpWindow(hostId: number) {
  sftpWindowRef.value.open(hostId);
}

// --- Host CRUD Handlers ---
const handleCreateHost = () => {
  hostModalRef.value.open();
};

const handleEditHost = (host: Host) => {
  hostModalRef.value.open(host);
};

const handleDeleteHost = async (id: number) => {
  await hostStore.deleteHost(id);
};

const handleHostSubmit = async (hostData: Host) => {
  if (hostData.id) {
    await hostStore.updateHost(hostData.id, hostData);
  } else {
    await hostStore.addHost(hostData);
  }
};

</script>

<style lang="scss" scoped>
.host-management-container {
  display: flex;
  height: 100%;
  background-color: #f0f2f5;
}

.left-panel {
  width: 240px;
  background-color: #fff;
  border-right: 1px solid #e8e8e8;
  display: flex;
  flex-direction: column;

  .panel-header {
    padding: 12px;
    border-bottom: 1px solid #e8e8e8;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .provider-list {
    flex-grow: 1;
    overflow-y: auto;
  }
}

.right-panel {
  flex-grow: 1;
  padding: 16px;
  background-color: #fff;
  margin-left: 8px;
  display: flex;
  flex-direction: column;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  margin-bottom: 16px;
}

.host-table {
  flex-grow: 1;
}

.provider-item {
  display: flex;
  justify-content: space-between;
  align-items: center;

  .actions {
    display: none; // Hidden by default
    gap: 8px;
    
    .anticon {
      color: #999;
      &:hover {
        color: #1890ff;
      }
    }
  }

  &:hover .actions {
    display: inline-flex; // Show on hover
  }
}
</style> 