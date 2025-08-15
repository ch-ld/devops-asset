<template>
  <el-dialog
    v-model="dialogVisible"
    :title="modalTitle"
    width="600px"
    :close-on-click-modal="false"
    class="ali-record-modal"
  >
    <template #header>
      <div class="modal-header">
        <span class="modal-title">{{ modalTitle }}</span>
        <div class="modal-tabs">
          <el-radio-group v-model="inputMode" size="small">
            <el-radio-button label="form">表单导向模式</el-radio-button>
            <el-radio-button label="expert">图形编辑模式</el-radio-button>
          </el-radio-group>
        </div>
      </div>
    </template>

    <div class="modal-content">
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="120px"
        class="record-form"
      >
        <!-- 记录类型 -->
        <el-form-item label="记录类型" prop="type">
          <el-select 
            v-model="formData.type" 
            placeholder="请选择记录类型"
            @change="handleTypeChange"
            style="width: 100%"
          >
            <el-option 
              v-for="type in recordTypes" 
              :key="type.value" 
              :label="type.label" 
              :value="type.value"
            >
              <div class="type-option">
                <span class="type-name">{{ type.label }}</span>
                <span class="type-desc">{{ type.description }}</span>
              </div>
            </el-option>
          </el-select>
        </el-form-item>

        <!-- 主机记录 -->
        <el-form-item label="主机记录" prop="name">
          <el-input 
            v-model="formData.name" 
            placeholder="请输入主机记录"
            style="width: 100%"
          >
            <template #suffix>
              <span class="domain-suffix">.{{ domainName }}</span>
            </template>
          </el-input>
          <div class="form-help">
            <el-icon><InfoFilled /></el-icon>
            <span>完整域名：{{ fullDomainName }}</span>
          </div>
        </el-form-item>

        <!-- 解析请求来源 -->
        <el-form-item label="解析请求来源" prop="line">
          <el-select 
            v-model="formData.line" 
            placeholder="请选择解析线路"
            style="width: 100%"
          >
            <el-option label="默认" value="default">
              <div class="line-option">
                <span>默认</span>
                <el-tag size="small" type="success">推荐</el-tag>
              </div>
            </el-option>
            <el-option label="电信" value="telecom" />
            <el-option label="联通" value="unicom" />
            <el-option label="移动" value="mobile" />
            <el-option label="海外" value="overseas" />
          </el-select>
          <div class="form-help">
            <el-link 
              type="primary" 
              :underline="false"
              @click="handleUpgradeToEnterprise"
            >
              升级至企业版使用自定义线路
            </el-link>
          </div>
        </el-form-item>

        <!-- TTL时间 -->
        <el-form-item label="TTL时间" prop="ttl">
          <el-select 
            v-model="formData.ttl" 
            placeholder="请选择TTL"
            style="width: 200px"
          >
            <el-option label="10 分钟" :value="600" />
            <el-option label="1 小时" :value="3600" />
            <el-option label="12 小时" :value="43200" />
            <el-option label="1 天" :value="86400" />
          </el-select>
          <span class="ttl-unit">分钟</span>
        </el-form-item>

        <!-- 记录值设置 -->
        <el-form-item label="记录值设置" prop="value">
          <div class="record-value-section">
            <!-- 单个记录值 -->
            <div v-if="!isMultiValue" class="single-value">
              <el-input
                v-model="formData.value"
                :placeholder="getValuePlaceholder()"
                type="textarea"
                :rows="3"
                style="width: 100%"
              />
            </div>

            <!-- 多个记录值 -->
            <div v-else class="multi-value">
              <div class="value-header">
                <span>记录值集合</span>
                <el-button 
                  type="primary" 
                  link 
                  @click="handleAddValue"
                >
                  <el-icon><Plus /></el-icon>
                  添加条目
                </el-button>
              </div>
              
              <div class="value-list">
                <div 
                  v-for="(item, index) in recordValues" 
                  :key="index"
                  class="value-item"
                >
                  <div class="value-inputs">
                    <el-input
                      v-model="item.value"
                      placeholder="请输入记录值"
                      style="flex: 1"
                    />
                    <el-switch
                      v-model="item.enabled"
                      active-text="启用"
                      inactive-text="暂停"
                      style="margin-left: 12px"
                    />
                    <el-input
                      v-model="item.remark"
                      placeholder="备注"
                      style="width: 120px; margin-left: 12px"
                    />
                  </div>
                  <div class="value-actions">
                    <el-button 
                      type="danger" 
                      link 
                      @click="handleRemoveValue(index)"
                      :disabled="recordValues.length <= 1"
                    >
                      删除
                    </el-button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </el-form-item>

        <!-- MX优先级 (仅MX记录显示) -->
        <el-form-item 
          v-if="formData.type === 'MX'" 
          label="MX优先级" 
          prop="mx_priority"
        >
          <el-input-number
            v-model="formData.mx_priority"
            :min="0"
            :max="65535"
            placeholder="10"
            style="width: 200px"
          />
        </el-form-item>

        <!-- 权重 (企业版功能) -->
        <el-form-item label="权重" prop="weight">
          <el-input-number
            v-model="formData.weight"
            :min="0"
            :max="255"
            placeholder="1"
            style="width: 200px"
            :disabled="!isEnterpriseUser"
          />
          <div class="form-help" v-if="!isEnterpriseUser">
            <el-link 
              type="primary" 
              :underline="false"
              @click="handleUpgradeToEnterprise"
            >
              升级至企业版使用权重功能
            </el-link>
          </div>
        </el-form-item>
      </el-form>
    </div>

    <template #footer>
      <div class="modal-footer">
        <div class="footer-left">
          <!-- 云解析个人版优惠广告 -->
          <div class="upgrade-banner">
            <el-icon><Star /></el-icon>
            <span>云解析个人版优惠</span>
            <el-tag type="danger" size="small">全球部署100+节点</el-tag>
            <el-tag type="warning" size="small">域名问询量更稳定</el-tag>
            <span class="price">￥19.9/年</span>
            <span class="original-price">原价48元/年</span>
            <el-button type="warning" size="small" @click="handleUpgrade">
              立即购买
            </el-button>
          </div>
        </div>
        <div class="footer-right">
          <el-button @click="handleCancel">取消</el-button>
          <el-button type="primary" @click="handleConfirm" :loading="loading">
            确定
          </el-button>
        </div>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import {
  InfoFilled,
  Plus,
  Star
} from '@element-plus/icons-vue'

// Props
interface Props {
  visible: boolean
  domainId?: number
  domainName?: string
  record?: any
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  domainName: '',
  record: null
})

// Emits
const emit = defineEmits(['update:visible', 'success'])

// 响应式数据
const dialogVisible = ref(false)
const loading = ref(false)
const inputMode = ref('form')
const formRef = ref()

// 表单数据
const formData = reactive({
  type: 'A',
  name: '',
  line: 'default',
  ttl: 600,
  value: '',
  mx_priority: 10,
  weight: 1
})

// 记录值集合（用于多值记录）
const recordValues = ref([
  { value: '', enabled: true, remark: '' }
])

// 记录类型选项
const recordTypes = [
  { value: 'A', label: 'A', description: '将域名指向一个IPv4地址' },
  { value: 'CNAME', label: 'CNAME', description: '将域名指向另一个域名' },
  { value: 'MX', label: 'MX', description: '将域名指向邮件服务器地址' },
  { value: 'TXT', label: 'TXT', description: '可任意填写，长度限制255' },
  { value: 'NS', label: 'NS', description: '域名服务器记录' },
  { value: 'AAAA', label: 'AAAA', description: '将域名指向一个IPv6地址' }
]

// 表单验证规则
const formRules = {
  type: [{ required: true, message: '请选择记录类型', trigger: 'change' }],
  name: [{ required: true, message: '请输入主机记录', trigger: 'blur' }],
  line: [{ required: true, message: '请选择解析线路', trigger: 'change' }],
  value: [{ required: true, message: '请输入记录值', trigger: 'blur' }]
}

// 计算属性
const modalTitle = computed(() => {
  return props.record ? '修改记录' : '添加记录'
})

const fullDomainName = computed(() => {
  const name = formData.name || '@'
  return name === '@' ? props.domainName : `${name}.${props.domainName}`
})

const isMultiValue = computed(() => {
  return false // 暂时不支持多值记录
})

const isEnterpriseUser = computed(() => {
  return false // 模拟非企业版用户
})

// 监听visible变化
watch(() => props.visible, (val) => {
  dialogVisible.value = val
  if (val && props.record) {
    // 编辑模式，填充数据
    Object.assign(formData, props.record)
  } else if (val) {
    // 新增模式，重置表单
    resetForm()
  }
})

watch(dialogVisible, (val) => {
  emit('update:visible', val)
})

// 方法
const resetForm = () => {
  Object.assign(formData, {
    type: 'A',
    name: '',
    line: 'default',
    ttl: 600,
    value: '',
    mx_priority: 10,
    weight: 1
  })
  recordValues.value = [{ value: '', enabled: true, remark: '' }]
}

const getValuePlaceholder = () => {
  const placeholders = {
    A: '请输入IPv4地址，如：192.168.1.1',
    CNAME: '请输入域名，如：www.example.com',
    MX: '请输入邮件服务器地址，如：mail.example.com',
    TXT: '请输入文本内容',
    NS: '请输入域名服务器，如：ns1.example.com',
    AAAA: '请输入IPv6地址'
  }
  return placeholders[formData.type] || '请输入记录值'
}

const handleTypeChange = () => {
  formData.value = ''
  if (formData.type === 'MX') {
    formData.mx_priority = 10
  }
}

const handleAddValue = () => {
  recordValues.value.push({ value: '', enabled: true, remark: '' })
}

const handleRemoveValue = (index) => {
  recordValues.value.splice(index, 1)
}

const handleUpgradeToEnterprise = () => {
  ElMessage.info('升级功能开发中...')
}

const handleUpgrade = () => {
  ElMessage.info('购买功能开发中...')
}

const handleCancel = () => {
  dialogVisible.value = false
}

const handleConfirm = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    ElMessage.success(props.record ? '修改成功' : '添加成功')
    emit('success')
    dialogVisible.value = false
  } catch (error) {
    console.error('表单验证失败:', error)
  } finally {
    loading.value = false
  }
}
</script>
