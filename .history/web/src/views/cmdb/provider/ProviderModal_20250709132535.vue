<template>
  <a-modal
    v-model:open="visible"
    :title="isEdit ? '编辑厂商' : '添加厂商'"
    @ok="handleOk"
    :confirm-loading="confirmLoading"
    @cancel="handleCancel"
  >
    <a-form ref="formRef" :model="formState" :rules="rules" layout="vertical">
      <a-form-item label="厂商名称" name="name">
        <a-input v-model:value="formState.name" placeholder="请输入厂商名称" />
      </a-form-item>
      <a-form-item label="厂商类型" name="type">
        <a-select v-model:value="formState.type" placeholder="请选择厂商类型">
          <a-select-option value="aliyun">阿里云</a-select-option>
          <a-select-option value="tencent">腾讯云</a-select-option>
          <a-select-option value="aws">AWS</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="Access Key ID" name="access_key_id">
        <a-input v-model:value="formState.access_key_id" placeholder="请输入 Access Key ID" />
      </a-form-item>
      <a-form-item label="Access Key Secret" name="access_key_secret">
        <a-input-password v-model:value="formState.access_key_secret" placeholder="请输入 Access Key Secret" />
      </a-form-item>
      <a-form-item label="状态" name="status">
        <a-switch v-model:checked="formState.status" checked-children="启用" un-checked-children="禁用" />
      </a-form-item>
      <a-form-item label="备注" name="remark">
        <a-textarea v-model:value="formState.remark" placeholder="请输入备注" />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { createProvider, updateProvider } from '@/api/cmdb/provider';
import type { Provider } from '@/types/cmdb';

const emit = defineEmits(['success']);

const visible = ref(false);
const confirmLoading = ref(false);
const isEdit = ref(false);
const formRef = ref();
const recordId = ref<number | null>(null);

const formState = reactive<Omit<Provider, 'id' | 'created_at' | 'updated_at'>>({
  name: '',
  type: 'aliyun',
  access_key_id: '',
  access_key_secret: '',
  status: true,
  remark: '',
});

const rules = {
  name: [{ required: true, message: '请输入厂商名称' }],
  type: [{ required: true, message: '请选择厂商类型' }],
  access_key_id: [{ required: true, message: '请输入Access Key ID' }],
  access_key_secret: [{ required: true, message: '请输入Access Key Secret' }],
};

const open = (record?: Provider) => {
  visible.value = true;
  isEdit.value = !!record;
  if (record) {
    recordId.value = record.id;
    Object.assign(formState, record);
  } else {
    recordId.value = null;
    formRef.value?.resetFields();
  }
};

const handleOk = async () => {
  try {
    await formRef.value.validate();
    confirmLoading.value = true;
    if (isEdit.value && recordId.value) {
      await updateProvider(recordId.value, formState);
      message.success('更新成功');
    } else {
      await createProvider(formState);
      message.success('创建成功');
    }
    visible.value = false;
    emit('success');
  } catch (error) {
    console.error(error);
  } finally {
    confirmLoading.value = false;
  }
};

const handleCancel = () => {
  formRef.value.resetFields();
  visible.value = false;
};

defineExpose({
  open,
});
</script> 