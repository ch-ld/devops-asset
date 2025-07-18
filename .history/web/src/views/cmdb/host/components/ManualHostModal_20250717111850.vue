<template>
  <a-modal
    title="添加自建主机"
    :open="visible"
    :width="700"
    :confirm-loading="loading"
    @ok="handleSubmit"
    @cancel="handleCancel"
  >
    <a-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      :label-col="{ span: 6 }"
      :wrapper-col="{ span: 18 }"
    >
      <a-form-item label="主机名称" name="name">
        <a-input v-model:value="formData.name" placeholder="请输入主机名称" />
      </a-form-item>

      <a-form-item label="实例ID" name="instance_id">
        <a-input v-model:value="formData.instance_id" placeholder="请输入实例ID或自定义标识" />
      </a-form-item>

      <a-form-item label="主机组" name="group_id">
        <a-tree-select
          v-model:value="formData.group_id"
          :tree-data="hostGroupOptions"
          placeholder="选择主机组"
          allow-clear
          tree-default-expand-all
        />
      </a-form-item>

      <a-form-item label="公网IP" name="public_ip">
        <a-input v-model:value="formData.public_ip" placeholder="请输入公网IP地址" />
      </a-form-item>

      <a-form-item label="私网IP" name="private_ip">
        <a-input v-model:value="formData.private_ip" placeholder="请输入私网IP地址" />
      </a-form-item>

      <a-form-item label="用户名" name="username">
        <a-input v-model:value="formData.username" placeholder="请输入SSH用户名" />
      </a-form-item>

      <a-form-item label="密码" name="password">
        <a-input-password
          v-model:value="formData.password"
          placeholder="请输入SSH密码"
          :autocomplete="false"
        />
      </a-form-item>

      <a-form-item label="操作系统" name="os">
        <a-select v-model:value="formData.os" placeholder="选择操作系统">
          <a-select-option value="CentOS 7">CentOS 7</a-select-option>
          <a-select-option value="CentOS 8">CentOS 8</a-select-option>
          <a-select-option value="Ubuntu 18.04">Ubuntu 18.04</a-select-option>
          <a-select-option value="Ubuntu 20.04">Ubuntu 20.04</a-select-option>
          <a-select-option value="Ubuntu 22.04">Ubuntu 22.04</a-select-option>
          <a-select-option value="Windows Server 2019">Windows Server 2019</a-select-option>
          <a-select-option value="Windows Server 2022">Windows Server 2022</a-select-option>
        </a-select>
      </a-form-item>

      <a-form-item label="状态" name="status">
        <a-select v-model:value="formData.status" placeholder="选择状态">
          <a-select-option value="running">运行中</a-select-option>
          <a-select-option value="stopped">已停止</a-select-option>
          <a-select-option value="unknown">未知</a-select-option>
        </a-select>
      </a-form-item>

      <a-form-item label="区域" name="region">
        <a-input v-model:value="formData.region" placeholder="请输入区域或机房信息" />
      </a-form-item>

      <a-form-item label="备注" name="remark">
        <a-textarea v-model:value="formData.remark" placeholder="请输入备注信息" :rows="3" />
      </a-form-item>

      <a-alert
        message="提示"
        description="自建主机将标记为手动管理，不会进行云端同步。密码将被加密存储。"
        type="info"
        show-icon
        style="margin-bottom: 16px"
      />
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
  import { ref, reactive, computed } from 'vue'
  import { ElMessage } from 'element-plus'
  import type { FormInstance } from 'ant-design-vue'
  import type { HostGroup } from '@/types/api/host'

  interface Props {
    hostGroups: HostGroup[]
  }

  const props = defineProps<Props>()

  const emit = defineEmits<{
    success: []
  }>()

  const visible = ref(false)
  const loading = ref(false)
  const formRef = ref<FormInstance>()

  const formData = reactive({
    name: '',
    instance_id: '',
    group_id: undefined as number | undefined,
    public_ip: '',
    private_ip: '',
    username: '',
    password: '',
    os: '',
    status: 'running',
    region: '',
    remark: ''
  })

  const rules = {
    name: [{ required: true, message: '请输入主机名称', trigger: 'blur' }],
    instance_id: [{ required: true, message: '请输入实例ID', trigger: 'blur' }],
    public_ip: [
      { required: true, message: '请输入公网IP', trigger: 'blur' },
      {
        pattern:
          /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/,
        message: '请输入正确的IP地址格式',
        trigger: 'blur'
      }
    ],
    username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
    password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
    os: [{ required: true, message: '请选择操作系统', trigger: 'change' }]
  }

  const hostGroupOptions = computed(() => {
    const buildTreeData = (groups: HostGroup[], parentId: number | null = null): any[] => {
      if (!groups || !Array.isArray(groups)) {
        return []
      }
      return groups
        .filter((group) => group.parent_id === parentId)
        .map((group) => ({
          title: group.name,
          value: group.id,
          key: group.id,
          children: buildTreeData(groups, group.id)
        }))
    }
    return buildTreeData(props.hostGroups || [])
  })

  const resetForm = () => {
    Object.assign(formData, {
      name: '',
      instance_id: '',
      group_id: undefined,
      public_ip: '',
      private_ip: '',
      username: '',
      password: '',
      os: '',
      status: 'running',
      region: '',
      remark: ''
    })
  }

  const handleSubmit = async () => {
    try {
      await formRef.value?.validate()
      loading.value = true

      const submitData = {
        ...formData,
        provider_type: 'manual',
        resource_type: 'manual',
        public_ip: formData.public_ip ? [formData.public_ip] : [],
        private_ip: formData.private_ip ? [formData.private_ip] : []
      }

      // 这里应该调用API
      console.log('Submit manual host data:', submitData)

      // 模拟API调用
      await new Promise((resolve) => setTimeout(resolve, 1000))

      ElMessage.success('自建主机创建成功')
      visible.value = false
      emit('success')
    } catch (error) {
      console.error('Submit error:', error)
      ElMessage.error('创建失败')
    } finally {
      loading.value = false
    }
  }

  const handleCancel = () => {
    visible.value = false
    resetForm()
  }

  const open = () => {
    visible.value = true
    resetForm()
  }

  defineExpose({
    open
  })
</script>

<style scoped>
  .ant-form-item {
    margin-bottom: 16px;
  }
</style>
