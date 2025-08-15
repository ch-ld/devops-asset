<template>
  <el-dialog
    v-model="dialogVisible"
    title="添加域名"
    width="500px"
    :close-on-click-modal="false"
    class="add-domain-modal"
  >
    <div class="modal-content">
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
        class="domain-form"
      >
        <el-form-item label="域名" prop="name">
          <el-input
            v-model="formData.name"
            placeholder="请输入域名，如：example.com"
            style="width: 100%"
          />
          <div class="form-help">
            <el-icon><InfoFilled /></el-icon>
            <span>请输入完整的域名，不包含协议前缀</span>
          </div>
        </el-form-item>

        <el-form-item label="DNS提供商" prop="provider_id">
          <el-select 
            v-model="formData.provider_id" 
            placeholder="请选择DNS提供商"
            style="width: 100%"
          >
            <el-option 
              v-for="provider in providers" 
              :key="provider.id" 
              :label="provider.name" 
              :value="provider.id"
            >
              <div class="provider-option">
                <span class="provider-name">{{ provider.name }}</span>
                <span class="provider-type">{{ provider.type }}</span>
              </div>
            </el-option>
          </el-select>
          <div class="form-help">
            <span>选择管理此域名的DNS提供商</span>
          </div>
        </el-form-item>

        <el-form-item label="描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入域名描述（可选）"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="自动同步" prop="auto_sync">
          <el-switch
            v-model="formData.auto_sync"
            active-text="启用"
            inactive-text="禁用"
          />
          <div class="form-help">
            <span>启用后将自动同步DNS记录变更</span>
          </div>
        </el-form-item>
      </el-form>
    </div>

    <template #footer>
      <div class="modal-footer">
        <el-button @click="handleCancel">取消</el-button>
        <el-button type="primary" @click="handleConfirm" :loading="loading">
          确定
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { InfoFilled } from '@element-plus/icons-vue'
import { dnsProviderApi as providerApi } from '@/api/dns/provider'

// Props
interface Props {
  visible: boolean
}

const props = withDefaults(defineProps<Props>(), {
  visible: false
})

// Emits
const emit = defineEmits(['update:visible', 'success'])

// 响应式数据
const dialogVisible = ref(false)
const loading = ref(false)
const formRef = ref()
const providers = ref([])

// 表单数据
const formData = reactive({
  name: '',
  provider_id: null,
  description: '',
  auto_sync: true
})

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入域名', trigger: 'blur' },
    { 
      pattern: /^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$/, 
      message: '请输入有效的域名格式', 
      trigger: 'blur' 
    }
  ],
  provider_id: [
    { required: true, message: '请选择DNS提供商', trigger: 'change' }
  ]
}

// 监听visible变化
watch(() => props.visible, (val) => {
  dialogVisible.value = val
  if (val) {
    resetForm()
  }
})

watch(dialogVisible, (val) => {
  emit('update:visible', val)
})

// 获取DNS提供商列表
const fetchProviders = async () => {
  try {
    const response = await providerApi.list({})
    providers.value = response.items || []
  } catch (error) {
    console.error('获取DNS提供商失败:', error)
    // 使用模拟数据
    providers.value = [
      { id: 1, name: 'Cloudflare', type: 'cloudflare' },
      { id: 2, name: '阿里云DNS', type: 'aliyun' },
      { id: 3, name: '腾讯云DNS', type: 'tencent' },
      { id: 4, name: 'DNSPod', type: 'dnspod' }
    ]
  }
}

// 方法
const resetForm = () => {
  Object.assign(formData, {
    name: '',
    provider_id: null,
    description: '',
    auto_sync: true
  })
  formRef.value?.clearValidate()
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
    
    ElMessage.success('添加域名成功')
    emit('success')
    dialogVisible.value = false
  } catch (error) {
    console.error('表单验证失败:', error)
  } finally {
    loading.value = false
  }
}

// 初始化
onMounted(() => {
  fetchProviders()
})
</script>

<style scoped lang="scss">
.add-domain-modal {
  .modal-content {
    padding: 20px 0;

    .domain-form {
      .el-form-item {
        margin-bottom: 24px;

        .form-help {
          margin-top: 8px;
          font-size: 12px;
          color: #8c8c8c;
          display: flex;
          align-items: center;
          gap: 4px;

          .el-icon {
            font-size: 14px;
            color: #1890ff;
          }
        }
      }

      .provider-option {
        display: flex;
        justify-content: space-between;
        align-items: center;
        width: 100%;

        .provider-name {
          font-weight: 500;
        }

        .provider-type {
          font-size: 12px;
          color: #8c8c8c;
        }
      }
    }
  }

  .modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
  }
}
</style>
