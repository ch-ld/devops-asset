<template>
  <a-modal
    :visible="visible"
    :title="title"
    @ok="handleSubmit"
    :confirm-loading="confirmLoading"
    :width="500"
    @cancel="handleCancel"
  >
    <a-form :model="formState" layout="vertical">
      <a-alert
        v-if="hostIds.length > 0"
        :message="`已选择 ${hostIds.length} 台主机`"
        type="info"
        style="margin-bottom: 16px"
      />

      <a-form-item label="操作类型">
        <a-radio-group v-model:value="formState.actionType">
          <a-radio value="add">添加标签</a-radio>
          <a-radio value="remove">删除标签</a-radio>
          <a-radio value="replace">替换标签</a-radio>
        </a-radio-group>
      </a-form-item>

      <a-form-item label="标签" name="tags">
        <a-select
          v-model:value="formState.selectedTags"
          mode="tags"
          placeholder="请输入或选择标签"
          style="width: 100%"
          :disabled="loading"
          :options="tagOptions"
          :max-tag-count="10"
        >
        </a-select>
        <div class="form-help-text"> 输入标签后按回车添加，或从下拉列表中选择 </div>
      </a-form-item>

      <a-form-item v-if="commonTags.length > 0" label="常用标签">
        <div class="common-tags">
          <a-tag
            v-for="tag in commonTags"
            :key="tag"
            class="clickable-tag"
            @click="addCommonTag(tag)"
          >
            {{ tag }}
          </a-tag>
        </div>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
  import { ref, reactive, computed, onMounted } from 'vue'
  import { ElMessage } from 'element-plus'
  import * as hostApi from '@/api/system/host'

  // 属性定义
  const props = defineProps({
    visible: {
      type: Boolean,
      default: false
    },
    hostIds: {
      type: Array as () => number[],
      default: () => []
    }
  })

  // 事件
  const emit = defineEmits(['update:visible', 'success'])

  // 状态
  const loading = ref(false)
  const confirmLoading = ref(false)

  // 表单状态
  const formState = reactive({
    actionType: 'add' as 'add' | 'remove' | 'replace',
    selectedTags: [] as string[]
  })

  // 处理取消
  const handleCancel = () => {
    emit('update:visible', false)
  }

  // 已有标签选项
  const tagOptions = ref<{ value: string; label: string }[]>([])
  // 常用标签
  const commonTags = ref<string[]>([
    '生产环境',
    '测试环境',
    '开发环境',
    'Linux',
    'Windows',
    'Web服务器',
    '数据库',
    '高优先级',
    '低优先级'
  ])

  // 计算属性
  const title = computed(() => {
    switch (formState.actionType) {
      case 'add':
        return '批量添加标签'
      case 'remove':
        return '批量删除标签'
      case 'replace':
        return '批量替换标签'
      default:
        return '标签操作'
    }
  })

  // 生命周期钩子
  onMounted(async () => {
    await loadAllTags()
  })

  // 加载所有已有标签
  const loadAllTags = async () => {
    loading.value = true
    try {
      // 获取所有主机的标签
      const response = await hostApi.getHostList()
      const allTags = new Set<string>()

      // 从response中正确提取hosts数据
      const hosts = response.data || []

      hosts.forEach((host) => {
        if (host.tags && Array.isArray(host.tags)) {
          host.tags.forEach((tag: string) => allTags.add(tag))
        } else if (typeof host.tags === 'string') {
          try {
            const parsedTags = JSON.parse(host.tags as string)
            if (Array.isArray(parsedTags)) {
              parsedTags.forEach((tag: string) => allTags.add(tag))
            }
          } catch (e) {
            // 忽略解析错误
          }
        }
      })

      // 转换为选项格式
      tagOptions.value = Array.from(allTags).map((tag) => ({
        value: tag,
        label: tag
      }))
    } catch (error) {
      console.error('加载标签失败:', error)
    } finally {
      loading.value = false
    }
  }

  // 添加常用标签
  const addCommonTag = (tag: string) => {
    if (!formState.selectedTags.includes(tag)) {
      formState.selectedTags.push(tag)
    }
  }

  // 提交表单
  const handleSubmit = async () => {
    if (formState.selectedTags.length === 0) {
      ElMessage.warning('请至少选择或输入一个标签')
      return
    }

    confirmLoading.value = true
    try {
      // 直接调用API而不是通过store
      await hostApi.batchUpdateTags({
        ids: props.hostIds,
        tags: formState.selectedTags,
        action: formState.actionType
      })

      ElMessage.success('标签操作成功')
      emit('success')
      emit('update:visible', false)

      // 重置表单
      formState.selectedTags = []
      formState.actionType = 'add'
    } catch (error) {
      console.error('标签操作失败:', error)
      ElMessage.error('标签操作失败，请重试')
    } finally {
      confirmLoading.value = false
    }
  }
</script>

<style lang="scss" scoped>
  .form-help-text {
    color: rgba(0, 0, 0, 0.45);
    font-size: 12px;
    margin-top: 4px;
  }

  .common-tags {
    margin-top: 8px;

    .clickable-tag {
      cursor: pointer;
      margin-bottom: 8px;

      &:hover {
        opacity: 0.8;
      }
    }
  }

  // 暗色主题适配
  html.dark {
    .form-help-text {
      color: rgba(255, 255, 255, 0.45);
    }
  }
</style>
