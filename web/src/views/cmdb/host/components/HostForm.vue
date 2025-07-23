<template>
  <div class="host-form-wrapper">
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="100px"
      @submit.prevent
      class="host-form"
    >
      <!-- 表单标题 -->
      <div class="form-title">
        <el-icon class="title-icon">
          <Plus />
        </el-icon>
        <span>{{ isEdit ? '编辑主机' : '添加主机' }}</span>
      </div>

      <!-- 基本信息 -->
      <div class="form-section">
        <h3 class="section-title">基本信息</h3>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="主机名称" prop="name">
              <el-input
                v-model="formData.name"
                placeholder="请输入主机名称"
              >
                <template #prefix>
                  <el-icon><Monitor /></el-icon>
                </template>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="实例ID" prop="instance_id">
              <el-input
                v-model="formData.instance_id"
                placeholder="自动生成"
                readonly
              >
                <template #prefix>
                  <el-icon><Cpu /></el-icon>
                </template>
                <template #append>
                  <el-button @click="handleGenerateInstanceId" size="small">
                    <el-icon><Refresh /></el-icon>
                    重新生成
                  </el-button>
                </template>
              </el-input>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="IP地址" prop="public_ip">
              <el-input
                v-model="formData.public_ip"
                placeholder="请输入IP地址"
              >
                <template #prefix>
                  <el-icon><Connection /></el-icon>
                </template>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="操作系统" prop="os">
              <div class="custom-select-wrapper">
                <select
                  v-model="formData.os"
                  class="custom-select"
                >
                  <option value="">请选择操作系统</option>
                  <option value="CentOS 7">CentOS 7</option>
                  <option value="CentOS 8">CentOS 8</option>
                  <option value="Ubuntu 18.04">Ubuntu 18.04</option>
                  <option value="Ubuntu 20.04">Ubuntu 20.04</option>
                  <option value="Ubuntu 22.04">Ubuntu 22.04</option>
                  <option value="Windows Server 2019">Windows Server 2019</option>
                  <option value="Windows Server 2022">Windows Server 2022</option>
                  <option value="其他">其他</option>
                </select>
                <div class="select-arrow">
                  <el-icon><ArrowDown /></el-icon>
                </div>
              </div>
            </el-form-item>
          </el-col>
        </el-row>
      </div>

      <!-- 认证信息 -->
      <div class="form-section">
        <h3 class="section-title">认证信息</h3>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="用户名" prop="username">
              <el-input
                v-model="formData.username"
                placeholder="请输入登录用户名"
              >
                <template #prefix>
                  <el-icon><User /></el-icon>
                </template>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="密码" prop="password">
              <el-input
                v-model="formData.password"
                type="password"
                placeholder="请输入登录密码"
                show-password
              >
                <template #prefix>
                  <el-icon><Lock /></el-icon>
                </template>
              </el-input>
            </el-form-item>
          </el-col>
        </el-row>
      </div>

      <!-- 分组信息 -->
      <div class="form-section">
        <h3 class="section-title">分组信息</h3>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="主机组" prop="group_id">
              <div class="custom-select-wrapper">
                <select
                  v-model="formData.group_id"
                  class="custom-select"
                  @change="handleGroupChange"
                >
                  <option value="">请选择主机组</option>
                  <option
                    v-for="group in groupOptions"
                    :key="group.value"
                    :value="group.value"
                  >
                    {{ group.label }}
                  </option>
                </select>
                <div class="select-arrow">
                  <el-icon><ArrowDown /></el-icon>
                </div>
              </div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="标签">
              <div class="tag-container">
                <el-tag
                  v-for="tag in formData.tags"
                  :key="tag"
                  closable
                  @close="removeTag(tag)"
                  class="tag-item"
                >
                  {{ tag }}
                </el-tag>
                <el-input
                  v-if="inputVisible"
                  ref="inputRef"
                  v-model="inputValue"
                  size="small"
                  @keyup.enter="handleInputConfirm"
                  @blur="handleInputConfirm"
                  class="tag-input"
                  placeholder="输入标签"
                />
                <el-button
                  v-else
                  size="small"
                  @click="showInput"
                  type="primary"
                  plain
                >
                  <el-icon><Plus /></el-icon>
                  添加标签
                </el-button>
              </div>
            </el-form-item>
          </el-col>
        </el-row>
      </div>

      <!-- 备注信息 -->
      <div class="form-section">
        <el-form-item label="备注" prop="remark">
          <el-input
            v-model="formData.remark"
            type="textarea"
            :rows="3"
            placeholder="请输入备注信息（可选）"
          />
        </el-form-item>
      </div>

      <!-- 操作按钮 -->
      <div class="form-actions">
        <el-button @click="handleCancel" size="large">
          取消
        </el-button>
        <el-button
          type="primary"
          @click="handleSubmit"
          :loading="submitting"
          size="large"
        >
          {{ isEdit ? '更新主机' : '创建主机' }}
        </el-button>
      </div>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, nextTick, watch } from 'vue'
import { ArrowDown } from '@element-plus/icons-vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { Plus, Monitor, Cpu, Connection, User, Lock, Refresh } from '@element-plus/icons-vue'
import { 
  createHost, 
  updateHost, 
  createManualHost,
  getProviderList,
  getHostGroupTree
} from '@/api/system/host'
import type { Host, Provider, HostGroup } from '@/types/api/host'

interface Props {
  host?: Host | null
}

const props = defineProps<Props>()
const emit = defineEmits<{
  success: []
  cancel: []
}>()

const formRef = ref<FormInstance>()
const submitting = ref(false)
const providerList = ref<Provider[]>([])
const groupTree = ref<any[]>([])
const groupOptions = ref<any[]>([])
const selectKey = ref(0)
const groupSelectRef = ref()

// 标签相关
const inputVisible = ref(false)
const inputValue = ref('')
const inputRef = ref()

const isEdit = computed(() => !!props.host)



// 生成随机实例ID
const generateInstanceId = () => {
  const prefix = 'manual-'
  const date = new Date()
  const dateStr = date.getFullYear().toString().slice(-2) +
                  (date.getMonth() + 1).toString().padStart(2, '0') +
                  date.getDate().toString().padStart(2, '0')
  const random = Math.random().toString(36).substr(2, 6).toUpperCase()
  return `${prefix}${dateStr}-${random}`
}

// 处理生成实例ID
const handleGenerateInstanceId = () => {
  formData.instance_id = generateInstanceId()
  ElMessage.success('实例ID已生成')
}

// 表单数据
const formData = reactive({
  name: '',
  instance_id: '',
  public_ip: '',
  os: '',
  username: '',
  password: '',
  group_id: undefined as number | undefined,
  remark: '',
  tags: [] as string[]
})

// 表单验证规则
const formRules: FormRules = {
  name: [{ required: true, message: '请输入主机名称', trigger: 'blur' }],
  public_ip: [{ required: true, message: '请输入IP地址', trigger: 'blur' }],
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }]
}

// 初始化表单数据
const initFormData = () => {
  if (props.host) {
    // 编辑模式：使用现有数据
    Object.assign(formData, {
      name: props.host.name,
      instance_id: props.host.instance_id,
      public_ip: Array.isArray(props.host.public_ip)
        ? props.host.public_ip.join(',')
        : props.host.public_ip,
      os: props.host.os,
      username: props.host.username,
      password: props.host.password,
      group_id: props.host.group_id,
      remark: props.host.remark,
      tags: props.host.tags || []
    })
  } else {
    // 新建模式：自动生成实例ID
    formData.instance_id = generateInstanceId()
  }
}

// 获取云厂商列表
const fetchProviders = async () => {
  try {
    const response = await getProviderList()
    // 修复数据结构问题，response直接就是数组
    if (Array.isArray(response)) {
      providerList.value = response
    } else if (response && response.data) {
      providerList.value = response.data
    } else {
      providerList.value = []
    }

    // 如果没有数据，添加一些默认选项
    if (providerList.value.length === 0) {
      providerList.value = [
        { id: 1, name: '阿里云', code: 'aliyun' },
        { id: 2, name: '腾讯云', code: 'tencent' },
        { id: 3, name: '华为云', code: 'huawei' },
        { id: 4, name: 'AWS', code: 'aws' },
        { id: 5, name: '自建机房', code: 'self' }
      ]
    }
  } catch (error) {
    console.error('获取云厂商列表失败:', error)
    // 失败时使用默认数据
    providerList.value = [
      { id: 1, name: '阿里云', code: 'aliyun' },
      { id: 2, name: '腾讯云', code: 'tencent' },
      { id: 3, name: '华为云', code: 'huawei' },
      { id: 4, name: 'AWS', code: 'aws' },
      { id: 5, name: '自建机房', code: 'self' }
    ]
  }
}

// 获取主机组树
const fetchGroupTree = async () => {
  try {
    const response = await getHostGroupTree()
    console.log('🌳 HostForm获取主机组数据:', response)

    let groupData = []

    // 根据实际的API响应结构处理数据
    if (Array.isArray(response)) {
      groupData = response
    } else if (response && Array.isArray(response.data)) {
      groupData = response.data
    } else {
      console.warn('⚠️ HostForm主机组数据格式异常:', response)
      groupData = []
    }

    // 过滤掉"全部主机"选项（id为0）
    const filteredGroups = groupData.filter(item => item.id !== 0)

    console.log('🌳 HostForm过滤后的主机组:', filteredGroups)

    // 构建树形结构的函数
    const buildTreeData = (items: any[], parentId: number | null = null): any[] => {
      return items
        .filter(item => item.parent_id === parentId)
        .map(item => ({
          value: item.id,
          label: item.name,
          children: buildTreeData(items, item.id)
        }))
    }

    const treeData = buildTreeData(filteredGroups)

    // 添加"无分组"选项
    groupTree.value = [
      { value: 0, label: '无分组', children: [] },
      ...treeData
    ]

    console.log('🌳 HostForm最终的主机组树:', groupTree.value)

    // 直接设置扁平化选项
    const flattenOptions = (options: any[], prefix = ''): any[] => {
      const result: any[] = []

      options.forEach(option => {
        const label = prefix ? `${prefix} / ${option.label}` : option.label
        result.push({
          value: option.value,
          label: label
        })

        if (option.children && option.children.length > 0) {
          result.push(...flattenOptions(option.children, label))
        }
      })

      return result
    }

    groupOptions.value = flattenOptions(groupTree.value)
    console.log('🔍 HostForm扁平化选项:', groupOptions.value)
    console.log('🔍 HostForm扁平化选项详情:', JSON.stringify(groupOptions.value, null, 2))

    // 验证每个选项
    groupOptions.value.forEach((option, index) => {
      console.log(`🔍 选项 ${index}:`, option)
    })

    // 强制重新渲染Select组件
    selectKey.value++
    await nextTick()
    console.log('✅ HostForm数据更新完成，强制重新渲染Select')
  } catch (error) {
    console.error('❌ HostForm获取主机组失败:', error)
    groupTree.value = [
      { value: 0, label: '无分组', children: [] }
    ]
    groupOptions.value = [
      { value: 0, label: '无分组' }
    ]
  }
}

// 标签操作
const removeTag = (tag: string) => {
  formData.tags.splice(formData.tags.indexOf(tag), 1)
}

const showInput = () => {
  inputVisible.value = true
  nextTick(() => {
    inputRef.value?.focus()
  })
}

const handleInputConfirm = () => {
  if (inputValue.value && !formData.tags.includes(inputValue.value)) {
    formData.tags.push(inputValue.value)
  }
  inputVisible.value = false
  inputValue.value = ''
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    submitting.value = true

    // 处理IP地址
    const publicIPs = formData.public_ip ? formData.public_ip.split(',').map(ip => ip.trim()).filter(Boolean) : []

    const submitData = {
      name: formData.name,
      instance_id: formData.instance_id,
      public_ip: publicIPs,
      private_ip: [], // 自建主机暂不需要私网IP
      os: formData.os || '',
      region: '', // 自建主机暂不需要地域
      username: formData.username || '',
      password: formData.password || '',
      status: 'running',
      provider_type: 'manual',
      resource_type: 'manual',
      provider_id: null,
      group_id: formData.group_id === 0 ? null : formData.group_id || null,
      expired_at: null, // 自建主机暂不设置过期时间
      remark: formData.remark || '',
      tags: formData.tags || [],
      configuration: {
        cpu_cores: null,
        memory_size: null,
        instance_type: '',
        zone_id: '',
        vpc_id: ''
      }
    }

    if (isEdit.value && props.host) {
      await updateHost(props.host.id, submitData)
      ElMessage.success('更新成功')
    } else {
      await createManualHost(submitData)
      ElMessage.success('创建成功')
    }

    emit('success')
  } catch (error) {
    ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
  } finally {
    submitting.value = false
  }
}

// 取消
const handleCancel = () => {
  emit('cancel')
}

// 处理Select焦点事件
const handleSelectFocus = () => {
  console.log('🎯 Select获得焦点，当前选项:', groupOptions.value)
  console.log('🎯 Select获得焦点，选项数量:', groupOptions.value.length)
}

// 处理Select点击事件
const handleSelectClick = () => {
  console.log('🖱️ Select被点击，当前选项:', groupOptions.value)
  console.log('🖱️ Select被点击，选项数量:', groupOptions.value.length)
}

// 处理主机组选择变化
const handleGroupChange = (event: Event) => {
  const target = event.target as HTMLSelectElement
  const value = target.value
  console.log('🔄 主机组选择变化:', value)

  // 确保数值类型正确
  if (value === '') {
    formData.group_id = null
  } else {
    formData.group_id = parseInt(value)
  }
}

// 监听groupTree变化
watch(groupTree, (newVal) => {
  console.log('🔄 HostForm groupTree变化:', newVal)
}, { deep: true })

// 监听groupOptions变化
watch(groupOptions, (newVal) => {
  console.log('🔄 HostForm groupOptions变化:', newVal)
  console.log('🔄 HostForm groupOptions数量:', newVal.length)
}, { deep: true })

// 初始化
onMounted(() => {
  initFormData()
  fetchProviders()
  fetchGroupTree()
})
</script>

<style scoped>
.host-form-wrapper {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.host-form {
  background: white;
  border-radius: 12px;
  padding: 32px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

/* 表单标题 */
.form-title {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 32px;
  padding-bottom: 16px;
  border-bottom: 2px solid #f0f2f5;
}

.title-icon {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 18px;
}

.form-title span {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
}

/* 表单分组 */
.form-section {
  margin-bottom: 28px;
}

.form-section:last-child {
  margin-bottom: 0;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #374151;
  margin: 0 0 16px 0;
  padding-left: 12px;
  border-left: 3px solid #4f46e5;
}

/* 表单项样式 */
:deep(.el-form-item) {
  margin-bottom: 20px;
}

:deep(.el-form-item__label) {
  font-weight: 500;
  color: #374151;
  margin-bottom: 6px;
}

:deep(.el-input__wrapper) {
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

:deep(.el-input__wrapper:hover) {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
}

:deep(.el-select .el-input__wrapper) {
  cursor: pointer;
}

:deep(.el-textarea__inner) {
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

:deep(.el-textarea__inner:focus) {
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
}

/* 标签容器 */
.tag-container {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
  min-height: 32px;
}

.tag-item {
  margin: 0;
  border-radius: 6px;
  background: #f3f4f6;
  border-color: #d1d5db;
}

.tag-input {
  width: 120px;
}

/* 修复Select下拉框样式 */
:deep(.el-select-dropdown) {
  z-index: 9999 !important;
}

:deep(.host-form-select-dropdown) {
  z-index: 9999 !important;
  max-height: 300px !important;
}

/* 确保下拉选项可见 */
:deep(.el-select-dropdown__item) {
  padding: 8px 12px !important;
  font-size: 14px !important;
}

/* 自定义Select样式 */
.custom-select-wrapper {
  position: relative;
  width: 100%;
}

.custom-select {
  width: 100%;
  height: 40px;
  padding: 8px 32px 8px 12px;
  border: 1px solid #dcdfe6;
  border-radius: 6px;
  background-color: #ffffff;
  font-size: 14px;
  color: #606266;
  appearance: none;
  -webkit-appearance: none;
  -moz-appearance: none;
  cursor: pointer;
  transition: all 0.3s ease;
  outline: none;
}

.custom-select:hover {
  border-color: #c0c4cc;
}

.custom-select:focus {
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
}

.custom-select option {
  padding: 8px 12px;
  color: #606266;
  background-color: #ffffff;
}

.custom-select option:hover {
  background-color: #f5f7fa;
}

.select-arrow {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  pointer-events: none;
  color: #c0c4cc;
  transition: transform 0.3s ease;
}

.custom-select:focus + .select-arrow {
  color: #409eff;
  transform: translateY(-50%) rotate(180deg);
}

/* 操作按钮 */
.form-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid #f0f2f5;
}

.form-actions .el-button {
  min-width: 120px;
  height: 44px;
  border-radius: 8px;
  font-weight: 500;
}

.form-actions .el-button--primary {
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  border: none;
}

.form-actions .el-button--primary:hover {
  background: linear-gradient(135deg, #4338ca 0%, #6d28d9 100%);
  transform: translateY(-1px);
  box-shadow: 0 8px 25px rgba(79, 70, 229, 0.3);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .host-form-wrapper {
    padding: 12px;
  }

  .host-form {
    padding: 20px;
  }

  .form-title {
    flex-direction: column;
    text-align: center;
    gap: 8px;
  }

  .form-actions {
    flex-direction: column;
  }

  .form-actions .el-button {
    width: 100%;
  }
}

/* 动画效果 */
.host-form {
  animation: fadeInUp 0.6s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
