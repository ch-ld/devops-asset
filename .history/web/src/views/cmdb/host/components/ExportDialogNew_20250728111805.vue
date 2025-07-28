<template>
  <el-dialog
    v-model="visible"
    title=""
    width="700px"
    :close-on-click-modal="false"
    class="export-dialog"
    :show-close="false"
  >
    <!-- 自定义头部 -->
    <template #header>
      <div class="dialog-header">
        <div class="header-left">
          <div class="header-icon">
            <el-icon size="24"><Download /></el-icon>
          </div>
          <div class="header-content">
            <h2 class="dialog-title">导出主机数据</h2>
            <p class="dialog-subtitle">选择导出范围和格式，快速导出主机信息</p>
          </div>
        </div>
        <el-button 
          type="text" 
          size="large" 
          @click="handleClose"
          class="close-btn"
        >
          <el-icon size="20"><Close /></el-icon>
        </el-button>
      </div>
    </template>

    <div class="export-container">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="0">
        <!-- 导出范围 -->
        <div class="form-section">
          <div class="section-header">
            <div class="section-icon">
              <el-icon size="20" color="#409EFF"><Filter /></el-icon>
            </div>
            <div class="section-title">
              <h3>导出范围</h3>
              <p>选择要导出的主机数据范围</p>
            </div>
          </div>
          
          <div class="scope-options">
            <div 
              v-for="option in scopeOptions" 
              :key="option.value"
              :class="['scope-card', { active: formData.scope === option.value }]"
              @click="formData.scope = option.value"
            >
              <div class="scope-icon">
                <el-icon size="24" :color="formData.scope === option.value ? '#409EFF' : '#909399'">
                  <component :is="option.icon" />
                </el-icon>
              </div>
              <div class="scope-content">
                <div class="scope-title">{{ option.title }}</div>
                <div class="scope-desc">{{ option.desc }}</div>
              </div>
              <div class="scope-radio">
                <el-radio 
                  v-model="formData.scope" 
                  :value="option.value"
                  @click.stop
                />
              </div>
            </div>
          </div>

          <!-- 主机组选择 -->
          <div v-if="formData.scope === 'group'" class="group-selector">
            <el-form-item prop="group_id">
              <el-tree-select
                v-model="formData.group_id"
                :data="hostGroupOptions"
                placeholder="请选择要导出的主机组"
                clearable
                :default-expand-all="true"
                :render-after-expand="false"
                node-key="id"
                :props="{ label: 'name', value: 'id' }"
                class="group-select"
              />
            </el-form-item>
          </div>
        </div>

        <!-- 导出格式 -->
        <div class="form-section">
          <div class="section-header">
            <div class="section-icon">
              <el-icon size="20" color="#67C23A"><Document /></el-icon>
            </div>
            <div class="section-title">
              <h3>导出格式</h3>
              <p>选择导出文件的格式类型</p>
            </div>
          </div>
          
          <div class="format-options">
            <div 
              v-for="format in formatOptions" 
              :key="format.value"
              :class="['format-card', { active: formData.format === format.value }]"
              @click="formData.format = format.value"
            >
              <div class="format-icon">
                <el-icon size="32" :color="formData.format === format.value ? format.color : '#909399'">
                  <component :is="format.icon" />
                </el-icon>
              </div>
              <div class="format-content">
                <div class="format-title">{{ format.title }}</div>
                <div class="format-desc">{{ format.desc }}</div>
              </div>
              <div class="format-radio">
                <el-radio 
                  v-model="formData.format" 
                  :value="format.value"
                  @click.stop
                />
              </div>
            </div>
          </div>
        </div>

        <!-- 导出字段 -->
        <div class="form-section">
          <div class="section-header">
            <div class="section-icon">
              <el-icon size="20" color="#E6A23C"><Grid /></el-icon>
            </div>
            <div class="section-title">
              <h3>导出字段</h3>
              <p>选择要导出的字段信息</p>
            </div>
          </div>
          
          <div class="fields-container">
            <div class="field-groups">
              <div 
                v-for="group in fieldGroups" 
                :key="group.key"
                class="field-group"
              >
                <div class="field-group-header">
                  <el-checkbox 
                    :model-value="group.allSelected" 
                    @change="handleSelectAllGroup(group.key, $event)"
                    :indeterminate="group.indeterminate"
                    class="group-checkbox"
                  >
                    <div class="group-info">
                      <div class="group-title">{{ group.title }}</div>
                      <div class="group-count">{{ group.selectedCount }}/{{ group.fields.length }}</div>
                    </div>
                  </el-checkbox>
                </div>
                <div class="field-list">
                  <el-checkbox-group v-model="formData.fields">
                    <div class="field-items">
                      <el-checkbox 
                        v-for="field in group.fields" 
                        :key="field.value"
                        :value="field.value"
                        class="field-item"
                      >
                        {{ field.label }}
                      </el-checkbox>
                    </div>
                  </el-checkbox-group>
                </div>
              </div>
            </div>
          </div>
        </div>
      </el-form>
    </div>

    <!-- 底部操作栏 -->
    <template #footer>
      <div class="dialog-footer">
        <div class="footer-left">
          <div class="export-info">
            <el-icon><InfoFilled /></el-icon>
            <span>已选择 {{ formData.fields.length }} 个字段</span>
          </div>
        </div>
        <div class="footer-right">
          <el-button @click="handleClose">取消</el-button>
          <el-button 
            type="primary" 
            @click="handleExport" 
            :loading="exporting"
            :icon="Download"
          >
            {{ exporting ? '导出中...' : '确认导出' }}
          </el-button>
        </div>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { 
  Download, 
  Close, 
  Filter, 
  Document, 
  Grid, 
  InfoFilled,
  Files,
  Folder,
  Monitor,
  DocumentCopy
} from '@element-plus/icons-vue'
import { getHostGroupTree } from '@/api/cmdb/hostGroup'
import type { HostGroup } from '@/api/cmdb/hostGroup'

interface ExportFormData {
  scope: 'all' | 'group' | 'current'
  group_id?: number
  format: 'excel' | 'csv'
  fields: string[]
}

interface Props {
  modelValue: boolean
  currentFilters?: Record<string, any>
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'export', data: ExportFormData): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const formRef = ref<FormInstance>()
const exporting = ref(false)
const hostGroupOptions = ref<HostGroup[]>([])

// 表单数据
const formData = reactive<ExportFormData>({
  scope: 'all',
  format: 'excel',
  fields: ['name', 'instance_id', 'status', 'public_ip', 'private_ip', 'os', 'region']
})

// 导出范围选项
const scopeOptions = [
  {
    value: 'all',
    title: '全部主机',
    desc: '导出系统中的所有主机',
    icon: Monitor
  },
  {
    value: 'group',
    title: '指定主机组',
    desc: '导出特定主机组的主机',
    icon: Folder
  },
  {
    value: 'current',
    title: '当前筛选结果',
    desc: '导出当前页面筛选的主机',
    icon: Filter
  }
]

// 格式选项
const formatOptions = [
  {
    value: 'excel',
    title: 'Excel 格式',
    desc: '适合数据分析和编辑',
    icon: Files,
    color: '#67C23A'
  },
  {
    value: 'csv',
    title: 'CSV 格式',
    desc: '通用格式，兼容性好',
    icon: DocumentCopy,
    color: '#409EFF'
  }
]

// 字段分组
const fieldGroups = computed(() => {
  const basicFields = ['name', 'instance_id', 'status', 'public_ip', 'private_ip', 'os', 'region']
  const configFields = ['configuration', 'username', 'provider_type', 'resource_type', 'group_name', 'provider_name']
  const extraFields = ['tags', 'expired_at', 'remark', 'created_at', 'updated_at']
  
  return [
    {
      key: 'basic',
      title: '基础信息',
      fields: [
        { value: 'name', label: '主机名称' },
        { value: 'instance_id', label: '实例ID' },
        { value: 'status', label: '状态' },
        { value: 'public_ip', label: '公网IP' },
        { value: 'private_ip', label: '私网IP' },
        { value: 'os', label: '操作系统' },
        { value: 'region', label: '区域' }
      ],
      get selectedCount() {
        return basicFields.filter(field => formData.fields.includes(field)).length
      },
      get allSelected() {
        return basicFields.every(field => formData.fields.includes(field))
      },
      get indeterminate() {
        const selected = basicFields.filter(field => formData.fields.includes(field))
        return selected.length > 0 && selected.length < basicFields.length
      }
    },
    {
      key: 'config',
      title: '配置信息',
      fields: [
        { value: 'configuration', label: '配置规格' },
        { value: 'username', label: '用户名' },
        { value: 'provider_type', label: '提供商类型' },
        { value: 'resource_type', label: '资源类型' },
        { value: 'group_name', label: '主机组' },
        { value: 'provider_name', label: '云账号' }
      ],
      get selectedCount() {
        return configFields.filter(field => formData.fields.includes(field)).length
      },
      get allSelected() {
        return configFields.every(field => formData.fields.includes(field))
      },
      get indeterminate() {
        const selected = configFields.filter(field => formData.fields.includes(field))
        return selected.length > 0 && selected.length < configFields.length
      }
    },
    {
      key: 'extra',
      title: '扩展信息',
      fields: [
        { value: 'tags', label: '标签' },
        { value: 'expired_at', label: '过期时间' },
        { value: 'remark', label: '备注' },
        { value: 'created_at', label: '创建时间' },
        { value: 'updated_at', label: '更新时间' }
      ],
      get selectedCount() {
        return extraFields.filter(field => formData.fields.includes(field)).length
      },
      get allSelected() {
        return extraFields.every(field => formData.fields.includes(field))
      },
      get indeterminate() {
        const selected = extraFields.filter(field => formData.fields.includes(field))
        return selected.length > 0 && selected.length < extraFields.length
      }
    }
  ]
})

// 表单验证规则
const rules: FormRules = {
  group_id: [
    {
      required: true,
      message: '请选择主机组',
      trigger: 'change',
      validator: (rule, value, callback) => {
        if (formData.scope === 'group' && !value) {
          callback(new Error('请选择主机组'))
        } else {
          callback()
        }
      }
    }
  ],
  fields: [
    {
      required: true,
      message: '请至少选择一个导出字段',
      trigger: 'change',
      validator: (rule, value, callback) => {
        if (!value || value.length === 0) {
          callback(new Error('请至少选择一个导出字段'))
        } else {
          callback()
        }
      }
    }
  ]
}

// 处理分组全选
const handleSelectAllGroup = (groupKey: string, checked: boolean) => {
  const group = fieldGroups.value.find(g => g.key === groupKey)
  if (!group) return

  const groupFields = group.fields.map(f => f.value)

  if (checked) {
    // 添加该分组的所有字段
    groupFields.forEach(field => {
      if (!formData.fields.includes(field)) {
        formData.fields.push(field)
      }
    })
  } else {
    // 移除该分组的所有字段
    formData.fields = formData.fields.filter(field => !groupFields.includes(field))
  }
}

// 获取主机组列表
const fetchHostGroups = async () => {
  try {
    const response = await getHostGroupTree()
    hostGroupOptions.value = response.data || []
  } catch (error) {
    console.error('获取主机组列表失败:', error)
    ElMessage.error('获取主机组列表失败')
  }
}

// 处理导出
const handleExport = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()

    exporting.value = true
    emit('export', { ...formData })

  } catch (error) {
    console.error('表单验证失败:', error)
  }
}

// 处理关闭
const handleClose = () => {
  visible.value = false
  exporting.value = false
}

// 监听对话框打开
watch(visible, (newVal) => {
  if (newVal) {
    fetchHostGroups()
  }
})

// 暴露方法给父组件
defineExpose({
  setExporting: (value: boolean) => {
    exporting.value = value
  }
})
</script>
