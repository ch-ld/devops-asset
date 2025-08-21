<template>
  <div class="domain-edit">
    <el-page-header @back="goBack" :content="isEdit ? '编辑域名' : '添加域名'">
      <template #extra>
        <el-space>
          <el-button @click="goBack">取消</el-button>
          <el-button type="primary" :loading="saving" @click="handleSave">
            {{ isEdit ? '保存' : '添加' }}
          </el-button>
        </el-space>
      </template>
    </el-page-header>

    <el-card class="mt-4">
      <el-form 
        ref="formRef" 
        :model="formData" 
        :rules="formRules" 
        label-width="120px"
        class="domain-form"
      >
        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="域名" prop="name">
              <el-input 
                v-model="formData.name" 
                placeholder="请输入域名，如：example.com"
                :disabled="isEdit"
              />
              <div class="form-tip">
                域名格式：example.com 或 sub.example.com
              </div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-select v-model="formData.status" placeholder="请选择域名状态">
                <el-option label="正常" value="active" />
                <el-option label="暂停" value="suspended" />
                <el-option label="过期" value="expired" />
                <el-option label="待处理" value="pending" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="分组" prop="groupId">
              <el-cascader
                v-model="formData.groupId"
                :options="groupOptions"
                :props="cascaderProps"
                placeholder="请选择域名分组"
                clearable
                filterable
                @change="handleGroupChange"
              />
              <div class="form-actions">
                <el-button text @click="showGroupModal = true">
                  <el-icon><Plus /></el-icon>
                  新建分组
                </el-button>
              </div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="DNS提供商" prop="providerId">
              <el-select v-model="formData.providerId" placeholder="请选择DNS提供商">
                <el-option 
                  v-for="provider in providerOptions" 
                  :key="provider.id"
                  :label="provider.name" 
                  :value="provider.id" 
                />
              </el-select>
              <div class="form-actions">
                <el-button text @click="showProviderModal = true">
                  <el-icon><Plus /></el-icon>
                  新建提供商
                </el-button>
              </div>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="注册商">
              <el-input v-model="formData.registrar" placeholder="请输入注册商名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="到期时间">
              <el-date-picker
                v-model="formData.expiresAt"
                type="date"
                placeholder="请选择到期时间"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="标签">
          <div class="tags-container">
            <el-tag
              v-for="tag in formData.tags"
              :key="tag.id"
              :color="tag.color"
              closable
              @close="removeTag(tag)"
              class="tag-item"
            >
              {{ tag.name }}
            </el-tag>
            <el-dropdown @command="addTag" v-if="availableTags.length > 0">
              <el-button text>
                <el-icon><Plus /></el-icon>
                添加标签
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item 
                    v-for="tag in availableTags" 
                    :key="tag.id"
                    :command="tag"
                  >
                    <el-tag :color="tag.color" size="small">{{ tag.name }}</el-tag>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <el-button text @click="showTagModal = true">
              <el-icon><Plus /></el-icon>
              新建标签
            </el-button>
          </div>
        </el-form-item>

        <el-row :gutter="24">
          <el-col :span="24">
            <el-form-item label="备注">
              <el-input
                v-model="formData.remark"
                type="textarea"
                :rows="3"
                placeholder="请输入备注信息"
                maxlength="500"
                show-word-limit
              />
            </el-form-item>
          </el-col>
        </el-row>

        <!-- 高级设置 -->
        <el-form-item>
          <el-divider content-position="left">高级设置</el-divider>
        </el-form-item>

        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="自动续费">
              <el-switch 
                v-model="formData.autoRenew" 
                active-text="开启"
                inactive-text="关闭"
              />
              <div class="form-tip">
                开启后系统将在到期前自动尝试续费
              </div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="到期提醒">
              <el-select v-model="formData.expiryReminderDays" multiple placeholder="选择提醒时间">
                <el-option label="30天前" :value="30" />
                <el-option label="15天前" :value="15" />
                <el-option label="7天前" :value="7" />
                <el-option label="3天前" :value="3" />
                <el-option label="1天前" :value="1" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="24">
          <el-col :span="12">
            <el-form-item label="HTTPS监控">
              <el-switch 
                v-model="formData.enableHttpsMonitor" 
                active-text="开启"
                inactive-text="关闭"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12" v-if="formData.enableHttpsMonitor">
            <el-form-item label="监控URL">
              <el-input 
                v-model="formData.monitorUrl" 
                placeholder="https://example.com"
              />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-card>

    <!-- 分组管理弹窗 -->
    <GroupModal 
      v-model:visible="showGroupModal"
      @success="handleGroupCreated"
    />

    <!-- 提供商管理弹窗 -->
    <ProviderModal 
      v-model:visible="showProviderModal"
      @success="handleProviderCreated"
    />

    <!-- 标签管理弹窗 -->
    <TagModal 
      v-model:visible="showTagModal"
      @success="handleTagCreated"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElForm } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { domainApi } from '@/api/dns/domain'
import { domainGroupApi } from '@/api/dns/domainGroup'
import { dnsProviderApi } from '@/api/dns/provider'
import { tagApi } from '@/api/dns/tag'
import GroupModal from './components/GroupModal.vue'
import ProviderModal from '../providers/components/ProviderModal.vue'
import TagModal from './components/TagModal.vue'

const route = useRoute()
const router = useRouter()
const formRef = ref<InstanceType<typeof ElForm>>()

const domainId = computed(() => route.params.id as string)
const isEdit = computed(() => !!domainId.value && domainId.value !== 'create')

// 表单数据
const formData = reactive({
  name: '',
  status: 'active',
  groupId: null as number | null,
  providerId: null as number | null,
  registrar: '',
  expiresAt: null as Date | null,
  tags: [] as any[],
  remark: '',
  autoRenew: false,
  expiryReminderDays: [30, 7, 1] as number[],
  enableHttpsMonitor: false,
  monitorUrl: ''
})

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入域名', trigger: 'blur' },
    { 
      pattern: /^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)*[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?$/,
      message: '请输入有效的域名格式',
      trigger: 'blur'
    }
  ],
  status: [
    { required: true, message: '请选择域名状态', trigger: 'change' }
  ]
}

// 选项数据
const groupOptions = ref([])
const providerOptions = ref([])
const tagOptions = ref([])

// 级联选择器配置
const cascaderProps = {
  value: 'id',
  label: 'name',
  children: 'children',
  checkStrictly: true,
  emitPath: false
}

// 可用标签（排除已选择的）
const availableTags = computed(() => {
  const selectedTagIds = formData.tags.map(tag => tag.id)
  return tagOptions.value.filter((tag: any) => !selectedTagIds.includes(tag.id))
})

// 弹窗控制
const showGroupModal = ref(false)
const showProviderModal = ref(false)
const showTagModal = ref(false)
const saving = ref(false)

// 获取域名详情
const fetchDomainDetail = async () => {
  if (!isEdit.value) return
  
  try {
    const response = await domainApi.getDetail(domainId.value)
    const domain = response.data
    
    Object.assign(formData, {
      name: domain.name,
      status: domain.status,
      groupId: domain.groupId,
      providerId: domain.providerId,
      registrar: domain.registrar,
      expiresAt: domain.expiresAt ? new Date(domain.expiresAt) : null,
      tags: domain.tags || [],
      remark: domain.remark,
      autoRenew: domain.autoRenew,
      expiryReminderDays: domain.expiryReminderDays || [30, 7, 1],
      enableHttpsMonitor: domain.enableHttpsMonitor,
      monitorUrl: domain.monitorUrl
    })
  } catch (error) {
    ElMessage.error('获取域名详情失败')
  }
}

// 获取分组选项
const fetchGroupOptions = async () => {
  try {
    const response = await domainGroupApi.list()
    groupOptions.value = buildTreeOptions(response.data)
  } catch (error) {
    console.error('获取分组列表失败:', error)
  }
}

// 获取提供商选项
const fetchProviderOptions = async () => {
  try {
    const response = await providerApi.list()
    providerOptions.value = response.data
  } catch (error) {
    console.error('获取提供商列表失败:', error)
  }
}

// 获取标签选项
const fetchTagOptions = async () => {
  try {
    const response = await tagApi.list()
    tagOptions.value = response.data
  } catch (error) {
    console.error('获取标签列表失败:', error)
  }
}

// 构建树形选项
const buildTreeOptions = (items: any[], parentId = 0): any[] => {
  const result: any[] = []
  items.forEach(item => {
    if (item.parentId === parentId) {
      const children = buildTreeOptions(items, item.id)
      const option: any = {
        id: item.id,
        name: item.name,
        value: item.id,
        label: item.name
      }
      if (children.length > 0) {
        option.children = children
      }
      result.push(option)
    }
  })
  return result
}

// 分组变化处理
const handleGroupChange = (value: any) => {
  formData.groupId = value
}

// 添加标签
const addTag = (tag: any) => {
  if (!formData.tags.find(t => t.id === tag.id)) {
    formData.tags.push({ ...tag })
  }
}

// 移除标签
const removeTag = (tag: any) => {
  const index = formData.tags.findIndex(t => t.id === tag.id)
  if (index > -1) {
    formData.tags.splice(index, 1)
  }
}

// 处理分组创建成功
const handleGroupCreated = (group: any) => {
  fetchGroupOptions()
  formData.groupId = group.id
  ElMessage.success('分组创建成功')
}

// 处理提供商创建成功
const handleProviderCreated = (provider: any) => {
  fetchProviderOptions()
  formData.providerId = provider.id
  ElMessage.success('提供商创建成功')
}

// 处理标签创建成功
const handleTagCreated = (tag: any) => {
  fetchTagOptions()
  addTag(tag)
  ElMessage.success('标签创建成功')
}

// 保存处理
const handleSave = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    saving.value = true
    
    const data = {
      ...formData,
      tagIds: formData.tags.map(tag => tag.id),
      expiresAt: formData.expiresAt ? formData.expiresAt.toISOString() : null
    }
    
    if (isEdit.value) {
      await domainApi.update(domainId.value, data)
      ElMessage.success('域名更新成功')
    } else {
      await domainApi.create(data)
      ElMessage.success('域名创建成功')
    }
    
    router.push('/dns/domains')
  } catch (error: any) {
    if (error.fields) {
      // 表单验证错误
      return
    }
    ElMessage.error(isEdit.value ? '域名更新失败' : '域名创建失败')
  } finally {
    saving.value = false
  }
}

// 返回
const goBack = () => {
  router.back()
}

onMounted(() => {
  fetchGroupOptions()
  fetchProviderOptions()
  fetchTagOptions()
  fetchDomainDetail()
})
</script>

<style scoped>
.domain-edit {
  padding: 20px;
}

.domain-form {
  max-width: 800px;
}

.form-tip {
  font-size: 12px;
  color: var(--el-color-info);
  margin-top: 4px;
}

.form-actions {
  margin-top: 4px;
}

.tags-container {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.tag-item {
  margin: 0;
}

.mt-4 {
  margin-top: 1rem;
}
</style> 
