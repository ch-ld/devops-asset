<template>
  <div class="host-edit-container">
    <a-card :bordered="false">
      <template #title>
        <div class="card-title">
          <span>{{ isEdit ? '编辑主机' : '添加主机' }}</span>
          <a-tag v-if="isEdit" color="blue">ID: {{ hostId }}</a-tag>
        </div>
      </template>
      <template #extra>
        <a-space>
          <a-button @click="goBack">
            <template #icon><ArrowLeftOutlined /></template>
            返回
          </a-button>
        </a-space>
      </template>

      <a-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        :label-col="{ span: 6 }"
        :wrapper-col="{ span: 16 }"
        @finish="handleSubmit"
      >
        <a-divider>基本信息</a-divider>
        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="主机名称" name="name">
              <a-input v-model:value="formData.name" placeholder="请输入主机名称" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="实例ID" name="instance_id">
              <a-input v-model:value="formData.instance_id" placeholder="请输入实例ID" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="主机状态" name="status">
              <a-select v-model:value="formData.status" placeholder="请选择主机状态">
                <a-select-option value="running">运行中</a-select-option>
                <a-select-option value="stopped">已停止</a-select-option>
                <a-select-option value="error">错误</a-select-option>
                <a-select-option value="expired">已过期</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="所属组" name="group_id">
              <a-select
                v-model:value="formData.group_id"
                placeholder="请选择所属主机组"
                :loading="groupLoading"
                :options="groupOptions"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="地区" name="region">
              <a-input v-model:value="formData.region" placeholder="请输入地区" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="云提供商" name="provider_id">
              <a-select
                v-model:value="formData.provider_id"
                placeholder="请选择云提供商"
                :loading="providerLoading"
                :options="providerOptions"
                allow-clear
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="到期时间" name="expired_at">
              <a-date-picker
                v-model:value="formData.expired_at"
                :show-time="{ format: 'HH:mm' }"
                format="YYYY-MM-DD HH:mm"
                placeholder="选择到期时间"
                style="width: 100%"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="资源类型" name="resource_type">
              <a-input v-model:value="formData.resource_type" placeholder="例如：ecs.g6.large" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-divider>连接信息</a-divider>
        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="公网IP" name="public_ip">
              <a-input v-model:value="publicIpInput" placeholder="多个IP用逗号分隔" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="内网IP" name="private_ip">
              <a-input v-model:value="privateIpInput" placeholder="多个IP用逗号分隔" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="用户名" name="username">
              <a-input v-model:value="formData.username" placeholder="请输入用户名" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="密码" name="password">
              <a-input-password 
                v-model:value="formData.password" 
                :placeholder="isEdit ? '如需修改密码请重新输入' : '请输入密码'"
                @input="isPasswordModified = true"
              />
              <div v-if="isEdit && !isPasswordModified" class="password-hint">
                <a-typography-text type="secondary" :style="{ fontSize: '12px' }">
                  <ExclamationCircleOutlined style="margin-right: 4px;" />
                  当前已设置密码，如不修改请保持为空
                </a-typography-text>
              </div>
            </a-form-item>
          </a-col>
        </a-row>

        <a-divider>配置信息</a-divider>
        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="CPU核心数" name="cpu_cores">
              <a-input-number
                v-model:value="formData.configuration.cpu_cores"
                :min="1"
                style="width: 100%"
                placeholder="CPU核心数"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="内存大小(GB)" name="memory_size">
              <a-input-number
                v-model:value="formData.configuration.memory_size"
                :min="1"
                style="width: 100%"
                placeholder="内存大小(GB)"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="操作系统" name="os">
              <a-input v-model:value="formData.os" placeholder="例如：CentOS 7.9" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="实例类型" name="instance_type">
              <a-input
                v-model:value="formData.configuration.instance_type"
                placeholder="例如：ecs.g6.large"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-divider>其他信息</a-divider>
        <a-row :gutter="24">
          <a-col :span="24">
            <a-form-item label="备注" name="remark">
              <a-textarea v-model:value="formData.remark" placeholder="请输入备注信息" :rows="4" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row>
          <a-col :span="24" style="text-align: center">
            <a-form-item :wrapper-col="{ offset: 6, span: 16 }">
              <a-space>
                <a-button type="primary" html-type="submit" :loading="loading"> 保存 </a-button>
                <a-button @click="goBack">取消</a-button>
              </a-space>
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-card>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, onMounted, computed } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { message } from 'ant-design-vue'
  import { ArrowLeftOutlined, ExclamationCircleOutlined } from '@ant-design/icons-vue'
  import dayjs from 'dayjs'
  import { useHostStore } from '@/store/modules/host'
  import type { Host } from '@/types/api/host'
  import { formatJsonData } from '@/utils/dataprocess/format'

  const route = useRoute()
  const router = useRouter()
  const hostStore = useHostStore()

  // 表单引用
  const formRef = ref()
  const loading = ref(false)
  const groupLoading = ref(false)
  const providerLoading = ref(false)

  // 编辑模式判断
  const hostId = computed(() => (route.params.id ? Number(route.params.id) : 0))
  const isEdit = computed(() => hostId.value > 0)

  // 添加修改标记
  const isPasswordModified = ref(false)

  // 主机组和云提供商选项
  const groupOptions = ref<any[]>([])
  const providerOptions = ref<any[]>([])

  // IP输入处理
  const publicIpInput = ref('')
  const privateIpInput = ref('')

  // 表单数据
  const formData = reactive<Host>({
    id: 0,
    instance_id: '',
    name: '',
    resource_type: '',
    region: '',
    username: '',
    password: '',
    public_ip: [],
    private_ip: [],
    configuration: {
      cpu_cores: 1,
      memory_size: 1,
      instance_type: ''
    },
    os: '',
    status: 'running',
    expired_at: undefined,
    provider_id: undefined,
    provider_type: 'manual',
    group_id: undefined,
    remark: '',
    created_at: '',
    updated_at: '',
    tags: []
  })

  // 表单校验规则
  const rules = {
    name: [{ required: true, message: '请输入主机名称', trigger: 'blur' }],
    instance_id: [{ required: true, message: '请输入实例ID', trigger: 'blur' }],
    status: [{ required: true, message: '请选择主机状态', trigger: 'change' }]
  }

  // 加载主机组选项
  const loadHostGroups = async () => {
    try {
      groupLoading.value = true
      await hostStore.fetchHostGroupTree()
      const flattenGroups = (groups: any[], result: any[] = []) => {
        groups.forEach((group) => {
          result.push({
            label: group.name,
            value: group.id
          })
          if (group.children && group.children.length) {
            flattenGroups(group.children, result)
          }
        })
        return result
      }
      groupOptions.value = flattenGroups(hostStore.hostGroupTree)
    } catch (error) {
      console.error('加载主机组失败:', error)
      message.error('加载主机组失败')
    } finally {
      groupLoading.value = false
    }
  }

  // 加载云提供商选项
  const loadProviders = async () => {
    try {
      providerLoading.value = true
      await hostStore.fetchProviders()
      providerOptions.value = hostStore.providerList.map((provider) => ({
        label: provider.name,
        value: provider.id
      }))
    } catch (error) {
      console.error('加载云提供商失败:', error)
      message.error('加载云提供商失败')
    } finally {
      providerLoading.value = false
    }
  }

  // 加载主机详情
  const loadHostDetail = async () => {
    if (!isEdit.value) return

    try {
      loading.value = true
      const hostDetail = await hostStore.getHost(hostId.value)

      // 重置修改标记
      isPasswordModified.value = false

      // 填充表单数据
      Object.assign(formData, hostDetail)

      // 在编辑模式下不显示明文密码
      if (isEdit.value) {
        formData.password = ''
      }

      // 处理日期时间
      if (formData.expired_at) {
        formData.expired_at = dayjs(formData.expired_at)
      }

      // 处理IP列表
      if (Array.isArray(formData.public_ip)) {
        publicIpInput.value = formData.public_ip.join(',')
      } else if (typeof formData.public_ip === 'string') {
        publicIpInput.value = formData.public_ip
        formData.public_ip = [formData.public_ip]
      }

      if (Array.isArray(formData.private_ip)) {
        privateIpInput.value = formData.private_ip.join(',')
      } else if (typeof formData.private_ip === 'string') {
        privateIpInput.value = formData.private_ip
        formData.private_ip = [formData.private_ip]
      }

      // 确保配置对象存在
      if (!formData.configuration) {
        formData.configuration = {
          cpu_cores: 1,
          memory_size: 1,
          instance_type: ''
        }
      }
    } catch (error) {
      console.error('加载主机详情失败:', error)
      message.error('加载主机详情失败')
      goBack()
    } finally {
      loading.value = false
    }
  }

  // 表单提交
  const handleSubmit = async () => {
    try {
      loading.value = true

      // 处理IP数据
      formData.public_ip = publicIpInput.value
        ? publicIpInput.value.split(',').map((ip) => ip.trim())
        : []
      formData.private_ip = privateIpInput.value
        ? privateIpInput.value.split(',').map((ip) => ip.trim())
        : []

      // JSON格式化数据
      const submitData = { ...formData }
      submitData.configuration = formatJsonData(formData.configuration)

      // 处理密码字段：在编辑模式下，只有用户修改了密码时才提交密码字段
      if (isEdit.value && !isPasswordModified.value) {
        delete submitData.password
      }

      // 处理日期
      if (submitData.expired_at && typeof submitData.expired_at !== 'string') {
        submitData.expired_at = (submitData.expired_at as any).format('YYYY-MM-DD HH:mm:ss')
      }

      if (isEdit.value) {
        await hostStore.updateHost(hostId.value, submitData)
        message.success('更新主机成功')
      } else {
        await hostStore.addHost(submitData)
        message.success('添加主机成功')
      }

      // 返回列表页
      goBack()
    } catch (error) {
      console.error('保存主机失败:', error)
      message.error('保存主机失败')
    } finally {
      loading.value = false
    }
  }

  // 返回上一页
  const goBack = () => {
    router.push('/cmdb/hosts')
  }

  // 生命周期钩子
  onMounted(async () => {
    await loadHostGroups()
    await loadProviders()
    await loadHostDetail()
  })
</script>

<style lang="scss" scoped>
  .host-edit-container {
    padding: 0;

    .card-title {
      display: flex;
      align-items: center;
      gap: 8px;
    }

    .ant-form-item {
      margin-bottom: 12px;
    }

    .password-hint {
      margin-top: 4px;
      padding-left: 24px;
    }
  }
</style>
