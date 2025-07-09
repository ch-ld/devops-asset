<template>
  <a-modal
    v-model:open="visible"
    :title="isEdit ? '编辑主机' : '添加主机'"
    @ok="handleOk"
    :confirm-loading="confirmLoading"
    @cancel="handleCancel"
    width="600px"
  >
    <a-form ref="formRef" :model="formState" :rules="rules" layout="vertical">
      <a-form-item label="主机名称" name="name">
        <a-input v-model:value="formState.name" placeholder="请输入主机名称" />
      </a-form-item>
      <a-form-item label="实例ID" name="instance_id">
        <a-input v-model:value="formState.instance_id" placeholder="请输入云厂商实例ID" />
      </a-form-item>
      <a-form-item label="归属厂商" name="provider_id">
        <a-select v-model:value="formState.provider_id" placeholder="请选择归属厂商" :loading="providerLoading">
          <a-select-option v-for="p in providers" :key="p.id" :value="p.id">
            {{ p.name }}
          </a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="SSH 用户名" name="username">
        <a-input v-model:value="formState.username" placeholder="请输入SSH登录用户名" />
      </a-form-item>
      <a-form-item label="SSH 密码" name="password">
        <a-input-password
          v-model:value="formState.password"
          :placeholder="isEdit ? '如需修改请输入新密码' : '请输入SSH登录密码'"
        />
      </a-form-item>
      <a-form-item label="公网IP" name="public_ip">
        <a-select
          v-model:value="formState.public_ip"
          mode="tags"
          style="width: 100%"
          placeholder="请输入并按回车确认"
          :token-separators="[',']"
        ></a-select>
      </a-form-item>
      <a-form-item label="私网IP" name="private_ip">
        <a-select
          v-model:value="formState.private_ip"
          mode="tags"
          style="width: 100%"
          placeholder="请输入并按回车确认"
          :token-separators="[',']"
        ></a-select>
      </a-form-item>
      <a-form-item label="备注" name="remark">
        <a-textarea v-model:value="formState.remark" placeholder="请输入备注信息" />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { useHostStore } from '@/store/modules/host';
import { listProviders } from '@/api/cmdb/provider';
import type { Provider } from '@/types/cmdb';

const emit = defineEmits(['success']);
const hostStore = useHostStore();

const visible = ref(false);
const confirmLoading = ref(false);
const isEdit = ref(false);
const formRef = ref();
const recordId = ref<number | null>(null);

const providers = ref<Provider[]>([]);
const providerLoading = ref(false);

const createInitialFormState = () => ({
  name: '',
  instance_id: '',
  provider_id: null,
  username: 'root',
  password: '',
  public_ip: [],
  private_ip: [],
  remark: '',
  // Default non-form values
  resource_type: 'ecs',
  region: '',
  configuration: {},
  os: '',
  status: 'running',
  expired_at: null,
});

const formState = reactive(createInitialFormState());

const rules = {
  name: [{ required: true, message: '请输入主机名称' }],
  instance_id: [{ required: true, message: '请输入实例ID' }],
  provider_id: [{ required: true, message: '请选择归属厂商' }],
  username: [{ required: true, message: '请输入SSH用户名' }],
};

const fetchAllProviders = async () => {
  providerLoading.value = true;
  try {
    // Assuming listProviders can fetch all without pagination if no params are given
    const { data } = await listProviders();
    providers.value = data;
  } catch (error) {
    message.error('获取厂商列表失败');
  } finally {
    providerLoading.value = false;
  }
};

const open = (record?: any) => {
  formRef.value?.resetFields();
  visible.value = true;
  isEdit.value = !!record;

  if (record) {
    recordId.value = record.id;
    const public_ip = typeof record.public_ip === 'string' ? JSON.parse(record.public_ip) : record.public_ip || [];
    const private_ip = typeof record.private_ip === 'string' ? JSON.parse(record.private_ip) : record.private_ip || [];
    Object.assign(formState, { ...record, password: '', public_ip, private_ip });
  } else {
    recordId.value = null;
    Object.assign(formState, createInitialFormState());
  }
};

const handleOk = async () => {
  try {
    await formRef.value.validate();
    confirmLoading.value = true;
    
    // The backend expects JSON strings for IP arrays
    const payload = {
      ...formState,
      public_ip: JSON.stringify(formState.public_ip),
      private_ip: JSON.stringify(formState.private_ip),
    };

    if (isEdit.value && recordId.value) {
      await hostStore.updateHost(recordId.value, payload);
      message.success('更新成功');
    } else {
      await hostStore.addHost(payload);
      message.success('创建成功');
    }
    visible.value = false;
    emit('success');
  } catch (error) {
    // Error message is handled in the store
    console.error(error);
  } finally {
    confirmLoading.value = false;
  }
};

const handleCancel = () => {
  visible.value = false;
};

onMounted(() => {
  fetchAllProviders();
});

defineExpose({
  open,
});
</script>

<style scoped>
.host-modal {
  width: 600px;
}
</style> 