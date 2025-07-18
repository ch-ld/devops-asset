<template>
  <a-modal
    :title="isEdit ? '编辑主机' : '添加主机'"
    :open="visible"
    :width="800"
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
      <a-tabs v-model:activeKey="activeTab">
        <a-tab-pane key="basic" tab="基本信息">
          <a-form-item label="云账号" name="provider_id">
            <a-select
              v-model:value="formData.provider_id"
              placeholder="选择云账号"
              allow-clear
              @change="handleProviderChange"
            >
              <a-select-option
                v-for="provider in providers"
                :key="provider.id"
                :value="provider.id"
              >
                {{ provider.name }} ({{ provider.type }})
              </a-select-option>
            </a-select>
          </a-form-item>

          <a-form-item label="实例ID" name="instance_id">
            <a-input v-model:value="formData.instance_id" placeholder="请输入实例ID" />
          </a-form-item>

          <a-form-item label="主机名称" name="name">
            <a-input v-model:value="formData.name" placeholder="请输入主机名称" />
          </a-form-item>

          <a-form-item label="资源类型" name="resource_type">
            <a-select v-model:value="formData.resource_type" placeholder="选择资源类型">
              <a-select-option value="ecs">阿里云ECS</a-select-option>
              <a-select-option value="ec2">AWS EC2</a-select-option>
              <a-select-option value="cvm">腾讯云CVM</a-select-option>
              <a-select-option value="manual">自建</a-select-option>
            </a-select>
          </a-form-item>

          <a-form-item label="区域" name="region">
            <a-input v-model:value="formData.region" placeholder="请输入区域" />
          </a-form-item>

          <a-form-item label="状态" name="status">
            <a-select v-model:value="formData.status" placeholder="选择状态">
              <a-select-option value="running">运行中</a-select-option>
              <a-select-option value="stopped">已停止</a-select-option>
              <a-select-option value="starting">启动中</a-select-option>
              <a-select-option value="stopping">停止中</a-select-option>
              <a-select-option value="unknown">未知</a-select-option>
            </a-select>
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

          <a-form-item label="操作系统" name="os">
            <a-input v-model:value="formData.os" placeholder="请输入操作系统" />
          </a-form-item>

          <a-form-item label="过期时间" name="expired_at">
            <a-date-picker
              v-model:value="formData.expired_at"
              placeholder="选择过期时间"
              style="width: 100%"
              show-time
              format="YYYY-MM-DD HH:mm:ss"
            />
          </a-form-item>
        </a-tab-pane>

        <a-tab-pane key="network" tab="网络配置">
          <a-form-item label="公网IP">
            <a-space direction="vertical" style="width: 100%">
              <div v-for="(ip, index) in formData.public_ip" :key="index">
                <a-input-group compact>
                  <a-input
                    v-model:value="formData.public_ip[index]"
                    placeholder="请输入公网IP"
                    style="width: calc(100% - 32px)"
                  />
                  <a-button type="link" danger @click="removeIP('public_ip', index)">
                    <template #icon><DeleteOutlined /></template>
                  </a-button>
                </a-input-group>
              </div>
              <a-button type="dashed" block @click="addIP('public_ip')">
                <template #icon><PlusOutlined /></template>
                添加公网IP
              </a-button>
            </a-space>
          </a-form-item>

          <a-form-item label="私网IP">
            <a-space direction="vertical" style="width: 100%">
              <div v-for="(ip, index) in formData.private_ip" :key="index">
                <a-input-group compact>
                  <a-input
                    v-model:value="formData.private_ip[index]"
                    placeholder="请输入私网IP"
                    style="width: calc(100% - 32px)"
                  />
                  <a-button type="link" danger @click="removeIP('private_ip', index)">
                    <template #icon><DeleteOutlined /></template>
                  </a-button>
                </a-input-group>
              </div>
              <a-button type="dashed" block @click="addIP('private_ip')">
                <template #icon><PlusOutlined /></template>
                添加私网IP
              </a-button>
            </a-space>
          </a-form-item>
        </a-tab-pane>

        <a-tab-pane key="auth" tab="认证信息">
          <a-form-item label="用户名" name="username">
            <a-input v-model:value="formData.username" placeholder="请输入用户名" />
          </a-form-item>

          <a-form-item label="密码" name="password">
            <a-input-password
              v-model:value="formData.password"
              placeholder="请输入密码"
              :autocomplete="false"
            />
          </a-form-item>

          <a-alert message="密码将被加密存储" type="info" show-icon style="margin-bottom: 16px" />
        </a-tab-pane>

        <a-tab-pane key="config" tab="配置信息">
          <a-form-item label="CPU核数">
            <a-input-number
              v-model:value="formData.configuration.cpu_cores"
              placeholder="CPU核数"
              :min="1"
              style="width: 100%"
            />
          </a-form-item>

          <a-form-item label="内存大小(GB)">
            <a-input-number
              v-model:value="formData.configuration.memory_size"
              placeholder="内存大小"
              :min="1"
              style="width: 100%"
            />
          </a-form-item>

          <a-form-item label="实例类型">
            <a-input v-model:value="formData.configuration.instance_type" placeholder="实例类型" />
          </a-form-item>

          <a-form-item label="可用区">
            <a-input v-model:value="formData.configuration.zone_id" placeholder="可用区ID" />
          </a-form-item>

          <a-form-item label="VPC ID">
            <a-input v-model:value="formData.configuration.vpc_id" placeholder="VPC ID" />
          </a-form-item>
        </a-tab-pane>

        <a-tab-pane key="other" tab="其他信息">
          <a-form-item label="标签">
            <a-select
              v-model:value="formData.tags"
              mode="tags"
              placeholder="请输入标签"
              style="width: 100%"
            />
          </a-form-item>

          <a-form-item label="备注" name="remark">
            <a-textarea v-model:value="formData.remark" placeholder="请输入备注信息" :rows="4" />
          </a-form-item>
        </a-tab-pane>
      </a-tabs>
    </a-form>
  </a-modal>
</template>

<script setup>
  import { ref, reactive, computed } from 'vue'
  import { ElMessage } from 'element-plus'
  // 使用 require 而非 import 以避免类型错误
  const { PlusOutlined, DeleteOutlined } = require('@ant-design/icons-vue')

  const props = defineProps({
    providers: Array,
    hostGroups: Array,
    visible: {
      type: Boolean,
      default: false
    },
    host: {
      type: Object,
      default: null
    },
    isEdit: {
      type: Boolean,
      default: false
    },
    onSuccess: Function
  })

  const emit = defineEmits(['success', 'update:visible'])

  const visible = ref(props.visible)
  const loading = ref(false)
  const isEdit = ref(props.isEdit)
  const activeTab = ref('basic')
  const formRef = ref<any>() // Changed from FormInstance to any

  const formData = reactive({
    id: undefined as number | undefined,
    provider_id: undefined as number | undefined,
    instance_id: '',
    name: '',
    resource_type: '',
    region: '',
    username: '',
    password: '',
    public_ip: [] as string[],
    private_ip: [] as string[],
    configuration: {
      cpu_cores: undefined as number | undefined,
      memory_size: undefined as number | undefined,
      instance_type: '',
      zone_id: '',
      vpc_id: ''
    },
    os: '',
    status: 'unknown',
    expired_at: undefined as any,
    group_id: undefined as number | undefined,
    tags: [] as string[],
    remark: ''
  })

  const rules = {
    instance_id: [{ required: true, message: '请输入实例ID', trigger: 'blur' }],
    name: [{ required: true, message: '请输入主机名称', trigger: 'blur' }],
    resource_type: [{ required: true, message: '请选择资源类型', trigger: 'change' }],
    username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
    password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
  }

  const hostGroupOptions = computed(() => {
    const buildTreeData = (groups: HostGroup[], parentId: number | null = null): any[] => {
      return groups
        .filter((group) => group.parent_id === parentId)
        .map((group) => ({
          title: group.name,
          value: group.id,
          key: group.id,
          children: buildTreeData(groups, group.id)
        }))
    }
    return buildTreeData(props.hostGroups)
  })

  const resetForm = () => {
    Object.assign(formData, {
      id: undefined,
      provider_id: undefined,
      instance_id: '',
      name: '',
      resource_type: '',
      region: '',
      username: '',
      password: '',
      public_ip: [],
      private_ip: [],
      configuration: {
        cpu_cores: undefined,
        memory_size: undefined,
        instance_type: '',
        zone_id: '',
        vpc_id: ''
      },
      os: '',
      status: 'unknown',
      expired_at: undefined,
      group_id: undefined,
      tags: [],
      remark: ''
    })
    activeTab.value = 'basic'
  }

  const handleProviderChange = (providerId: number) => {
    const provider = props.providers.find((p) => p.id === providerId)
    if (provider) {
      formData.resource_type =
        provider.type === 'aliyun'
          ? 'ecs'
          : provider.type === 'aws'
            ? 'ec2'
            : provider.type === 'tencent'
              ? 'cvm'
              : 'manual'
      formData.region = provider.region
    }
  }

  const addIP = (type: 'public_ip' | 'private_ip') => {
    formData[type].push('')
  }

  const removeIP = (type: 'public_ip' | 'private_ip', index: number) => {
    formData[type].splice(index, 1)
  }

  const handleSubmit = async () => {
    try {
      await formRef.value?.validate()
      loading.value = true

      const submitData = {
        ...formData,
        expired_at: formData.expired_at
          ? new Date(formData.expired_at).toISOString().slice(0, 19).replace('T', ' ')
          : undefined,
        public_ip: formData.public_ip.filter((ip) => ip.trim()),
        private_ip: formData.private_ip.filter((ip) => ip.trim())
      }

      // 这里应该调用API
      console.log('Submit data:', submitData)

      // 模拟API调用
      await new Promise((resolve) => setTimeout(resolve, 1000))

      ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
      visible.value = false
      emit('success')
    } catch (error) {
      console.error('Submit error:', error)
      ElMessage.error('操作失败')
    } finally {
      loading.value = false
    }
  }

  const handleCancel = () => {
    visible.value = false
    resetForm()
  }

  const open = (record?: Host) => {
    visible.value = true
    isEdit.value = !!record

    if (record) {
      Object.assign(formData, {
        ...record,
        expired_at: record.expired_at ? new Date(record.expired_at) : undefined,
        public_ip: Array.isArray(record.public_ip) ? record.public_ip : [],
        private_ip: Array.isArray(record.private_ip) ? record.private_ip : [],
        configuration: record.configuration || {},
        tags: Array.isArray(record.tags) ? record.tags : []
      })
    } else {
      resetForm()
    }
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
