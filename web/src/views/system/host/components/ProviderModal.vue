<template>
  <a-modal
    v-model:open="isModalVisible"
    :title="isEditMode ? '编辑云厂商账号' : '新建云厂商账号'"
    @ok="handleOk"
    @cancel="handleCancel"
    :confirm-loading="isLoading"
  >
    <a-form ref="formRef" :model="formState" :rules="rules" layout="vertical">
      <a-form-item label="名称" name="name">
        <a-input v-model:value="formState.name" placeholder="请输入账号别名" />
      </a-form-item>
      <a-form-item label="厂商" name="type">
        <a-select v-model:value="formState.type" placeholder="请选择云厂商">
          <a-select-option value="aliyun">阿里云</a-select-option>
          <a-select-option value="tencent">腾讯云</a-select-option>
          <a-select-option value="huawei">华为云</a-select-option>
          <a-select-option value="aws">AWS</a-select-option>
          <a-select-option value="volcengine">火山引擎</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="AccessKey ID" name="access_key_id">
        <a-input v-model:value="formState.access_key_id" placeholder="请输入 AccessKey ID" />
      </a-form-item>
      <a-form-item label="AccessKey Secret" name="access_key_secret">
        <a-input-password v-model:value="formState.access_key_secret" placeholder="请输入 AccessKey Secret" />
      </a-form-item>
      <a-form-item label="备注" name="remark">
        <a-textarea v-model:value="formState.remark" :rows="3" />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, watch, reactive, toRaw } from 'vue';
import type { PropType } from 'vue';
import type { FormInstance } from 'ant-design-vue';
import type { Provider } from '@/types/api/host';

const props = defineProps({
  visible: { type: Boolean, required: true },
  provider: { type: Object as PropType<Provider | null>, default: null },
  isLoading: { type: Boolean, default: false },
});

const emit = defineEmits(['update:visible', 'submit']);

const isModalVisible = ref(false);
const formRef = ref<FormInstance>();

const defaultFormState = {
  name: '',
  type: undefined,
  access_key_id: '',
  access_key_secret: '',
  remark: '',
};

const formState = reactive({ ...defaultFormState });
const isEditMode = ref(false);

watch(() => props.visible, (newValue) => {
  isModalVisible.value = newValue;
  if (newValue) {
    if (props.provider) {
      isEditMode.value = true;
      // Populate form for editing
      Object.assign(formState, {
        ...toRaw(props.provider),
        access_key_secret: '', // Do not show existing secret
      });
    } else {
      isEditMode.value = false;
      // Reset form for creating
      Object.assign(formState, defaultFormState);
    }
  }
});

const rules = {
  name: [{ required: true, message: '请输入账号名称' }],
  type: [{ required: true, message: '请选择厂商类型' }],
  access_key_id: [{ required: true, message: '请输入AccessKey ID' }],
  access_key_secret: [{ required: !isEditMode.value, message: '请输入AccessKey Secret' }],
};

const handleOk = async () => {
  try {
    const values = await formRef.value?.validate();
    if (values) {
      emit('submit', { ...values, id: props.provider?.id });
    }
  } catch (errorInfo) {
    console.log('Failed:', errorInfo);
  }
};

const handleCancel = () => {
  emit('update:visible', false);
};
</script> 