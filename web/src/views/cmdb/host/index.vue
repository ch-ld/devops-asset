<template>
  <div class="p-4">
    <a-card :bordered="false">
      <div class="flex justify-between mb-4">
        <a-button type="primary" @click="openHostModal()">添加主机</a-button>
        <div>
          <a-input-search placeholder="搜索主机..." class="mr-2" style="width: 240px" @search="onSearch" />
          <a-button @click="syncHosts">同步主机</a-button>
        </div>
      </div>
      <a-table :columns="columns" :data-source="hostList" :loading="loading" row-key="id">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'action'">
            <a-button type="link" size="small" @click="openSftpWindow(record)">SFTP</a-button>
            <a-button type="link" size="small" @click="openWebSsh(record)">WebSSH</a-button>
            <a-button type="link" size="small" @click="openHostModal(record)">编辑</a-button>
            <a-popconfirm
              title="确定要删除此主机吗？"
              ok-text="确定"
              cancel-text="取消"
              @confirm="handleDelete(record.id)"
            >
              <a-button type="link" size="small" danger>删除</a-button>
            </a-popconfirm>
          </template>
        </template>
      </a-table>
    </a-card>

    <HostModal ref="hostModalRef" @success="fetchHosts" />
    <SftpWindow ref="sftpWindowRef" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useHostStore } from '@/store/modules/host';
import { storeToRefs } from 'pinia';

import HostModal from '@/views/system/host/components/HostModal.vue';
import SftpWindow from '@/views/system/host/components/SftpWindow.vue';

const router = useRouter();
const hostStore = useHostStore();
const { hostList } = storeToRefs(hostStore);

const loading = ref(false);
const hostModalRef = ref<InstanceType<typeof HostModal> | null>(null);
const sftpWindowRef = ref<InstanceType<typeof SftpWindow> | null>(null);

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id' },
  { title: '实例ID', dataIndex: 'instance_id', key: 'instance_id' },
  { title: '主机名', dataIndex: 'name', key: 'name' },
  { title: '公网IP', dataIndex: 'public_ip', key: 'public_ip' },
  { title: '私网IP', dataIndex: 'private_ip', key: 'private_ip' },
  { title: '状态', dataIndex: 'status', key: 'status' },
  { title: '操作', key: 'action' },
];

const fetchHosts = async () => {
  loading.value = true;
  try {
    await hostStore.fetchHosts();
  } catch (error) {
    message.error('获取主机列表失败');
  } finally {
    loading.value = false;
  }
};

const onSearch = (value: string) => {
  // Implement search logic if needed
};

const syncHosts = () => {
  message.info('同步功能待实现');
};

const openHostModal = (record?: any) => {
  hostModalRef.value?.open(record);
};

const openSftpWindow = (record: any) => {
  sftpWindowRef.value?.open(record.id);
};

const openWebSsh = (record: any) => {
  const routeUrl = router.resolve({
    path: `/webssh`, // A dedicated route for WebSSH might be better
    query: { host_id: record.id },
  });
  window.open(routeUrl.href, '_blank');
};

const handleDelete = async (id: number) => {
  try {
    await hostStore.deleteHost(id);
    message.success('删除成功');
    fetchHosts();
  } catch (error) {
    message.error('删除失败');
  }
};

onMounted(() => {
  fetchHosts();
});
</script>

<style scoped></style> 