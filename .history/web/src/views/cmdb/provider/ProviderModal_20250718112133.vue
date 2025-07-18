<template>
  <el-dialog
    v-model="dialogVisible"
    :title="isEdit ? '编辑厂商' : '添加厂商'"
    width="520px"
    :close-on-click-modal="false"
    @close="handleCancel"
  >
    <el-form ref="formRef" :model="formState" :rules="rules" label-width="110px">
      <el-form-item label="厂商名称" prop="name">
        <el-input v-model="formState.name" placeholder="请输入厂商名称" />
      </el-form-item>

      <el-form-item label="厂商类型" prop="type">
        <el-select v-model="formState.type" placeholder="请选择厂商类型">
          <el-option label="阿里云" value="aliyun" />
          <el-option label="腾讯云" value="tencent" />
          <el-option label="AWS" value="aws" />
        </el-select>
      </el-form-item>

      <el-form-item label="AccessKey ID" prop="access_key_id">
        <el-input v-model="formState.access_key_id" placeholder="请输入 AccessKey ID" />
      </el-form-item>

      <el-form-item label="AccessKey Secret" prop="access_key_secret">
        <el-input
          v-model="formState.access_key_secret"
          type="password"
          show-password
          placeholder="请输入 AccessKey Secret"
        />
      </el-form-item>

      <el-form-item label="状态" prop="status">
        <el-switch v-model="formState.status" active-text="启用" inactive-text="禁用" />
      </el-form-item>

      <el-form-item label="备注">
        <el-input v-model="formState.remark" type="textarea" placeholder="请输入备注" />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="handleCancel">取消</el-button>
      <el-button type="primary" :loading="loading" @click="handleOk">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
  import { ref, reactive, watch, type PropType } from 'vue'
  import { ElMessage } from 'element-plus'
  import * as hostApi from '@/api/system/host'
  import type { Provider } from '@/types/api/host'
  import type { FormInstance } from 'element-plus'

  // expose events
  const emit = defineEmits(['success'])

  // dialog visibility control
  const props = defineProps({
    visible: Boolean,
    record: Object as PropType<Provider | null>
  })

  const dialogVisible = ref(props.visible)
  watch(
    () => props.visible,
    (val) => (dialogVisible.value = val)
  )
  watch(dialogVisible, (val) => {
    if (!val) resetForm()
  })

  const isEdit = ref(false)
  const loading = ref(false)
  const formRef = ref<FormInstance>()

  const formState = reactive<Omit<Provider, 'id' | 'created_at' | 'updated_at'>>({
    name: '',
    type: 'aliyun',
    access_key_id: '',
    access_key_secret: '',
    status: true,
    remark: ''
  })

  const rules = {
    name: [{ required: true, message: '请输入厂商名称', trigger: 'blur' }],
    type: [{ required: true, message: '请选择厂商类型', trigger: 'change' }],
    access_key_id: [{ required: true, message: '请输入Access Key ID', trigger: 'blur' }],
    access_key_secret: [{ required: true, message: '请输入Access Key Secret', trigger: 'blur' }]
  }

  function open(record?: Provider) {
    dialogVisible.value = true
    isEdit.value = !!record
    if (record) {
      Object.assign(formState, {
        name: record.name,
        type: record.type,
        access_key_id: record.access_key_id,
        access_key_secret: '', // 安全原因不回填
        status: record.status,
        remark: record.remark
      })(formState as any).id = record.id
    }
  }

  defineExpose({ open })

  function resetForm() {
    formRef.value?.resetFields()
    Object.assign(formState, {
      name: '',
      type: 'aliyun',
      access_key_id: '',
      access_key_secret: '',
      status: true,
      remark: ''
    })
    isEdit.value = false
  }

  async function handleOk() {
    await formRef.value?.validate()
    loading.value = true
    try {
      if (isEdit.value && (formState as any).id) {
        await hostApi.updateProvider((formState as any).id, formState)
        ElMessage.success('更新成功')
      } else {
        await hostApi.createProvider(formState)
        ElMessage.success('创建成功')
      }
      emit('success')
      dialogVisible.value = false
    } catch (err: any) {
      console.error(err)
      ElMessage.error(err.message || '操作失败')
    } finally {
      loading.value = false
    }
  }

  function handleCancel() {
    dialogVisible.value = false
  }
</script>

<style scoped>
  .form-help-text {
    font-size: 12px;
    color: var(--el-text-color-secondary);
  }
</style>
