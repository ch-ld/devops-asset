<template>
  <el-dialog
    title="添加自建主机"
    v-model="dialogVisible"
    :width="700"
    :close-on-click-modal="false"
    @close="handleCancel"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-width="100px"
      label-position="right"
      status-icon
    >
      <el-form-item label="主机名称" prop="name">
        <el-input v-model="formData.name" placeholder="请输入主机名称" />
      </el-form-item>

      <el-form-item label="实例ID" prop="instance_id">
        <el-input v-model="formData.instance_id" placeholder="请输入实例ID或自定义标识" />
      </el-form-item>

      <el-form-item label="主机组" prop="group_id">
        <el-cascader
          v-model="formData.group_id"
          :options="hostGroupOptions"
          :props="{ checkStrictly: true, emitPath: false }"
          placeholder="选择主机组"
          clearable
          style="width: 100%"
        />
      </el-form-item>

      <el-form-item label="公网IP" prop="public_ip">
        <el-input v-model="formData.public_ip" placeholder="请输入公网IP地址" />
      </el-form-item>

      <el-form-item label="私网IP" prop="private_ip">
        <el-input v-model="formData.private_ip" placeholder="请输入私网IP地址" />
      </el-form-item>

      <el-form-item label="用户名" prop="username">
        <el-input v-model="formData.username" placeholder="请输入SSH用户名" />
      </el-form-item>

      <el-form-item label="密码" prop="password">
        <el-input
          v-model="formData.password"
          type="password"
          placeholder="请输入SSH密码"
          show-password
          autocomplete="off"
        />
      </el-form-item>

      <el-form-item label="操作系统" prop="os">
        <el-select v-model="formData.os" placeholder="选择操作系统" style="width: 100%">
          <el-option value="CentOS 7" label="CentOS 7" />
          <el-option value="CentOS 8" label="CentOS 8" />
          <el-option value="Ubuntu 18.04" label="Ubuntu 18.04" />
          <el-option value="Ubuntu 20.04" label="Ubuntu 20.04" />
          <el-option value="Ubuntu 22.04" label="Ubuntu 22.04" />
          <el-option value="Windows Server 2019" label="Windows Server 2019" />
          <el-option value="Windows Server 2022" label="Windows Server 2022" />
        </el-select>
      </el-form-item>

      <el-form-item label="状态" prop="status">
        <el-select v-model="formData.status" placeholder="选择状态" style="width: 100%">
          <el-option value="running" label="运行中" />
          <el-option value="stopped" label="已停止" />
          <el-option value="unknown" label="未知" />
        </el-select>
      </el-form-item>

      <el-form-item label="区域" prop="region">
        <el-input v-model="formData.region" placeholder="请输入区域或机房信息" />
      </el-form-item>

      <el-form-item label="备注" prop="remark">
        <el-input
          v-model="formData.remark"
          type="textarea"
          placeholder="请输入备注信息"
          :rows="3"
        />
      </el-form-item>

      <el-alert
        title="提示"
        type="info"
        description="自建主机将标记为手动管理，不会进行云端同步。密码将被加密存储。"
        show-icon
        :closable="false"
        style="margin-bottom: 16px"
      />
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleCancel">取消</el-button>
        <el-button type="primary" :loading="loading" @click="handleSubmit"> 确认 </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
  import { ref, reactive, computed, onMounted, watch } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { createManualHost, updateHost, getHostGroupTree } from '@/api/system/host'
  import { useHostStore } from '@/store/modules/host'

  const props = defineProps({
    visible: {
      type: Boolean,
      default: false
    },
    host: {
      type: Object,
      default: () => ({})
    },
    isEdit: {
      type: Boolean,
      default: false
    }
  })

  const emit = defineEmits(['update:visible', 'success'])

  const hostStore = useHostStore()
  const dialogVisible = ref(props.visible)
  const loading = ref(false)
  const formRef = ref(null)
  const hostGroups = ref([])

  const formData = reactive({
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

  // 监听visible属性变化
  watch(
    () => props.visible,
    (val) => {
      dialogVisible.value = val
      if (val && props.isEdit && props.host) {
        // 编辑模式下，填充表单数据
        fillFormData()
      } else if (val) {
        // 添加模式下，重置表单
        resetForm()
      }
    }
  )

  // 监听dialogVisible变化，同步更新父组件的visible属性
  watch(dialogVisible, (val) => {
    emit('update:visible', val)
  })

  // 获取主机组数据
  const fetchHostGroups = async () => {
    try {
      loading.value = true
      const response = await getHostGroupTree()
      hostGroups.value = response || []
    } catch (error) {
      console.error('Failed to fetch host groups:', error)
      ElMessage.error('获取主机组失败，请刷新重试')
      hostGroups.value = []
    } finally {
      loading.value = false
    }
  }

  onMounted(() => {
    fetchHostGroups()
    if (props.isEdit && props.host) {
      fillFormData()
    }
  })

  // 将主机组数据转换为级联选择器所需的格式
  const hostGroupOptions = computed(() => {
    if (!hostGroups.value || !Array.isArray(hostGroups.value)) {
      return []
    }

    const buildCascaderOptions = (groups, parentId = null) => {
      return groups
        .filter((group) => group.parent_id === parentId)
        .map((group) => ({
          value: group.id,
          label: group.name,
          children: buildCascaderOptions(groups, group.id)
        }))
    }
    return buildCascaderOptions(hostGroups.value)
  })

  // 重置表单
  const resetForm = () => {
    if (formRef.value) {
      formRef.value.resetFields()
    }
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

  // 填充表单数据（编辑模式）
  const fillFormData = () => {
    if (!props.host) return

    // 处理IP地址数组转字符串
    const publicIp = Array.isArray(props.host.public_ip)
      ? props.host.public_ip.join(',')
      : props.host.public_ip || ''

    const privateIp = Array.isArray(props.host.private_ip)
      ? props.host.private_ip.join(',')
      : props.host.private_ip || ''

    Object.assign(formData, {
      name: props.host.name || '',
      instance_id: props.host.instance_id || '',
      group_id: props.host.group_id,
      public_ip: publicIp,
      private_ip: privateIp,
      username: props.host.username || '',
      password: '', // 出于安全考虑，不回填密码
      os: props.host.os || '',
      status: props.host.status || 'running',
      region: props.host.region || '',
      remark: props.host.remark || ''
    })
  }

  // 提交表单
  const handleSubmit = async () => {
    if (!formRef.value) return

    try {
      await formRef.value.validate()

      loading.value = true

      // 处理IP地址，将逗号分隔的字符串转为数组
      const formattedData = {
        ...formData,
        public_ip: formData.public_ip ? formData.public_ip.split(',').map((ip) => ip.trim()) : [],
        private_ip: formData.private_ip ? formData.private_ip.split(',').map((ip) => ip.trim()) : []
      }

      if (props.isEdit && props.host.id) {
        // 编辑模式
        await updateHost(props.host.id, formattedData)
        ElMessage.success('更新主机成功')
      } else {
        // 添加模式
        await createManualHost(formattedData)
        ElMessage.success('添加主机成功')
      }

      dialogVisible.value = false
      emit('success')
    } catch (error) {
      console.error('表单提交失败:', error)
      ElMessage.error(props.isEdit ? '更新主机失败' : '添加主机失败')
    } finally {
      loading.value = false
    }
  }

  // 取消操作
  const handleCancel = () => {
    dialogVisible.value = false
    resetForm()
  }
</script>

<style lang="scss" scoped>
  .dialog-footer {
    display: flex;
    justify-content: flex-end;
  }
</style>
