<template>
  <el-dialog
    v-model="visible"
    :title="isEdit ? '编辑主机' : '添加主机'"
    width="800px"
    :close-on-click-modal="false"
    @close="handleCancel"
  >
    <el-form ref="formRef" :model="formData" :rules="rules" label-width="110px">
      <el-tabs v-model="activeTab" class="mb-2">
        <el-tab-pane label="基本信息" name="basic">
          <!-- 基本信息 Tab 已迁移 -->
          <el-form-item label="云账号" prop="provider_id">
            <el-select v-model="formData.provider_id" placeholder="选择云账号" clearable>
              <el-option
                v-for="p in providers"
                :key="p.id"
                :label="`${p.name} (${p.type})`"
                :value="p.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="实例ID" prop="instance_id">
            <el-input v-model="formData.instance_id" placeholder="请输入实例ID" />
          </el-form-item>

          <el-form-item label="主机名称" prop="name">
            <el-input v-model="formData.name" placeholder="请输入主机名称" />
          </el-form-item>

          <el-form-item label="资源类型" prop="resource_type">
            <el-select v-model="formData.resource_type" placeholder="选择资源类型">
              <el-option label="阿里云ECS" value="ecs" />
              <el-option label="AWS EC2" value="ec2" />
              <el-option label="腾讯云CVM" value="cvm" />
              <el-option label="自建" value="manual" />
            </el-select>
          </el-form-item>

          <el-form-item label="区域" prop="region">
            <el-input v-model="formData.region" placeholder="请输入区域" />
          </el-form-item>

          <el-form-item label="状态" prop="status">
            <el-select v-model="formData.status" placeholder="选择状态">
              <el-option label="运行中" value="running" />
              <el-option label="已停止" value="stopped" />
              <el-option label="启动中" value="starting" />
              <el-option label="停止中" value="stopping" />
              <el-option label="未知" value="unknown" />
            </el-select>
          </el-form-item>

          <el-form-item label="主机组" prop="group_id">
            <el-tree-select
              v-model="formData.group_id"
              :data="hostGroupOptions"
              placeholder="选择主机组"
              clearable
              node-key="value"
              :props="{ label: 'title', children: 'children' }"
            />
          </el-form-item>

          <el-form-item label="操作系统" prop="os">
            <el-input v-model="formData.os" placeholder="请输入操作系统" />
          </el-form-item>

          <el-form-item label="过期时间" prop="expired_at">
            <el-date-picker
              v-model="formData.expired_at"
              type="datetime"
              placeholder="选择过期时间"
              style="width: 100%"
              format="YYYY-MM-DD HH:mm:ss"
            />
          </el-form-item>
        </el-tab-pane>
        <el-tab-pane label="网络配置" name="network">
          <!-- 公网 IP 列表 -->
          <el-form-item label="公网IP">
            <el-space direction="vertical" style="width: 100%">
              <div v-for="(ip, idx) in formData.public_ip" :key="`pub-${idx}`" class="w-100">
                <el-input
                  v-model="formData.public_ip[idx]"
                  placeholder="请输入公网IP"
                  class="mr-1"
                  style="width: calc(100% - 40px)"
                />
                <el-button type="danger" plain @click="removeIP('public_ip', idx)">
                  <el-icon><Delete /></el-icon>
                </el-button>
              </div>
              <el-button type="primary" plain @click="addIP('public_ip')">
                <el-icon><Plus /></el-icon> 添加公网IP
              </el-button>
            </el-space>
          </el-form-item>

          <!-- 私网 IP 列表 -->
          <el-form-item label="私网IP">
            <el-space direction="vertical" style="width: 100%">
              <div v-for="(ip, idx) in formData.private_ip" :key="`pri-${idx}`" class="w-100">
                <el-input
                  v-model="formData.private_ip[idx]"
                  placeholder="请输入私网IP"
                  class="mr-1"
                  style="width: calc(100% - 40px)"
                />
                <el-button type="danger" plain @click="removeIP('private_ip', idx)">
                  <el-icon><Delete /></el-icon>
                </el-button>
              </div>
              <el-button type="primary" plain @click="addIP('private_ip')">
                <el-icon><Plus /></el-icon> 添加私网IP
              </el-button>
            </el-space>
          </el-form-item>
        </el-tab-pane>

        <el-tab-pane label="认证信息" name="auth">
          <el-form-item label="用户名" prop="username">
            <el-input v-model="formData.username" placeholder="请输入用户名" />
          </el-form-item>

          <el-form-item label="密码" prop="password">
            <el-input
              v-model="formData.password"
              type="password"
              show-password
              placeholder="请输入密码"
            />
          </el-form-item>

          <el-alert type="info" title="密码将被加密存储" show-icon class="mb-2" />
        </el-tab-pane>

        <el-tab-pane label="配置信息" name="config">
          <el-form-item label="CPU核数">
            <el-input-number
              v-model="formData.configuration.cpu_cores"
              :min="1"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item label="内存(GB)">
            <el-input-number
              v-model="formData.configuration.memory_size"
              :min="1"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item label="实例类型">
            <el-input v-model="formData.configuration.instance_type" placeholder="实例类型" />
          </el-form-item>
          <el-form-item label="可用区">
            <el-input v-model="formData.configuration.zone_id" placeholder="可用区ID" />
          </el-form-item>
          <el-form-item label="VPC ID">
            <el-input v-model="formData.configuration.vpc_id" placeholder="VPC ID" />
          </el-form-item>
        </el-tab-pane>

        <el-tab-pane label="其他信息" name="other">
          <el-form-item label="标签">
            <el-select
              v-model="formData.tags"
              multiple
              filterable
              allow-create
              placeholder="请输入标签"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item label="备注" prop="remark">
            <el-input
              type="textarea"
              v-model="formData.remark"
              :rows="4"
              placeholder="请输入备注信息"
            />
          </el-form-item>
        </el-tab-pane>
      </el-tabs>
    </el-form>

    <template #footer>
      <el-button @click="handleCancel">取 消</el-button>
      <el-button type="primary" :loading="loading" @click="handleSubmit">确 定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
  import { ref, reactive, computed } from 'vue'
  import { ElMessage } from 'element-plus'
  import type { FormInstance } from 'element-plus'
  // 使用 Element Plus 图标替代 Ant Design 图标
  import { Plus, Delete } from '@element-plus/icons-vue'

  interface Provider {
    id: number
    name: string
    type: string
    region?: string
  }

  interface HostGroup {
    id: number
    name: string
    parent_id?: number | null
  }

  const props = defineProps<{
    providers: Provider[]
    hostGroups: HostGroup[]
    visible: boolean
    host: Record<string, any> | null
    isEdit: boolean
    onSuccess?: () => void
  }>({
    providers: Array as unknown as () => Provider[],
    hostGroups: Array as unknown as () => HostGroup[],
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
  const formRef = ref(null)

  const formData = reactive({
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

  const rules = {
    instance_id: [{ required: true, message: '请输入实例ID', trigger: 'blur' }],
    name: [{ required: true, message: '请输入主机名称', trigger: 'blur' }],
    resource_type: [{ required: true, message: '请选择资源类型', trigger: 'change' }],
    username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
    password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
  }

  const hostGroupOptions = computed(() => {
    if (!props.hostGroups || !Array.isArray(props.hostGroups)) {
      return []
    }

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

  const handleProviderChange = (providerId: number | undefined) => {
    const provider =
      props.providers && Array.isArray(props.providers)
        ? props.providers.find((p: Provider) => p.id === providerId)
        : null

    if (provider) {
      formData.resource_type =
        provider.type === 'aliyun'
          ? 'ecs'
          : provider.type === 'aws'
            ? 'ec2'
            : provider.type === 'tencent'
              ? 'cvm'
              : 'manual'
      formData.region = provider.region || ''
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

      try {
        // 调用真实的API
        if (isEdit.value) {
          // 更新主机
          await import('@/api/system/host').then(({ updateHost }) => {
            return updateHost(formData.id, submitData)
          })
          ElMessage.success('主机信息已更新')
        } else {
          // 创建主机
          await import('@/api/system/host').then(({ createHost }) => {
            return createHost(submitData)
          })
          ElMessage.success('主机创建成功')
        }

        visible.value = false
        emit('success')
        emit('update:visible', false)
        resetForm()
      } catch (error) {
        console.error('API call error:', error)
        ElMessage.error('操作失败: ' + (error.message || '未知错误'))
      }
    } catch (error) {
      console.error('Submit error:', error)
      ElMessage.error('表单验证失败，请检查输入')
    } finally {
      loading.value = false
    }
  }

  const handleCancel = () => {
    visible.value = false
    resetForm()
  }

  const open = (record) => {
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
