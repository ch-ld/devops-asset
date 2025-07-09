<template>
  <div class="p-4">
    <a-card :bordered="false" title="厂商管理">
      <div class="flex justify-between mb-4">
        <a-button type="primary" @click="handleAdd">添加厂商</a-button>
        <a-input-search
          v-model:value="searchKeyword"
          placeholder="搜索厂商名称..."
          style="width: 240px"
          @search="onSearch"
        />
      </div>
      <a-table
        :columns="columns"
        :data-source="providerList"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="record.status ? 'success' : 'error'">
              {{ record.status ? '有效' : '无效' }}
            </a-tag>
          </template>
          <template v-if="column.key === 'type'">
            <a-tag>
              {{ record.type.toUpperCase() }}
            </a-tag>
          </template>
          <template v-if="column.key === 'created_at'">
            {{ new Date(record.created_at).toLocaleString() }}
          </template>
          <template v-if="column.key === 'action'">
            <a-button type="link" size="small" @click="handleSync(record.id)" :loading="syncing[record.id]">同步</a-button>
            <a-button type="link" size="small" @click="handleEdit(record)">编辑</a-button>
            <a-popconfirm
              title="确定要删除此厂商吗？"
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
    <ProviderModal ref="providerModalRef" @success="fetchProviders" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue';
import { listProviders, deleteProvider, syncProvider } from '@/api/cmdb/provider';
import type { Provider } from '@/types/cmdb';
import ProviderModal from './ProviderModal.vue';

const providerList = ref<Provider[]>([]);
const loading = ref(false);
const searchKeyword = ref('');
const providerModalRef = ref<InstanceType<typeof ProviderModal> | null>(null);
const syncing = reactive<Record<number, boolean>>({});

const pagination = ref({
  current: 1,
  pageSize: 10,
  total: 0,
});

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
  { title: '厂商名称', dataIndex: 'name', key: 'name' },
  { title: '类型', dataIndex: 'type', key: 'type', align: 'center' },
  { title: '状态', dataIndex: 'status', key: 'status', align: 'center' },
  { title: '备注', dataIndex: 'remark', key: 'remark', ellipsis: true },
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at' },
  { title: '操作', key: 'action', width: 220, align: 'center' },
];

const fetchProviders = async () => {
  loading.value = true;
  try {
    const { data } = await listProviders({
      page: pagination.value.current,
      pageSize: pagination.value.pageSize,
      keyword: searchKeyword.value,
    });
    providerList.value = data;
    // Note: The backend should ideally return pagination info.
    // Here we'll just use the length for a simple demo.
    pagination.value.total = data.length; 
  } catch (error) {
    message.error('获取厂商列表失败');
  } finally {
    loading.value = false;
  }
};

const handleTableChange = (pag: any) => {
  pagination.value.current = pag.current;
  fetchProviders();
};

const onSearch = () => {
  pagination.value.current = 1;
  fetchProviders();
};

const handleAdd = () => {
  providerModalRef.value?.open();
};

const handleEdit = (record: Provider) => {
  providerModalRef.value?.open(record);
};

const handleDelete = async (id: number) => {
  try {
    await deleteProvider(id);
    message.success('删除成功');
    fetchProviders();
  } catch (error) {
    message.error('删除失败');
  }
};

const handleSync = async (id: number) => {
  syncing[id] = true;
  try {
    await syncProvider(id);
    message.success('同步指令已发送');
  } catch (error) {
    message.error('同步失败');
  } finally {
    syncing[id] = false;
  }
};

onMounted(() => {
  fetchProviders();
});
</script>

<style scoped>
</style> 