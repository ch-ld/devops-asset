<template>
  <div class="user-center">
    <div class="user-center-grid">
      <!-- 左侧用户信息卡片 -->
      <div class="user-profile-card">
        <div class="profile-header">
          <div class="profile-bg">
            <img src="@imgs/user/bg.webp" alt="背景图" />
            <div class="bg-overlay"></div>
          </div>
          
          <div class="profile-avatar-section">
            <div class="avatar-wrapper">
              <div class="avatar-container">
                <img 
                  class="avatar-image" 
                  :src="userInfo.avatar || defaultAvatar" 
                  :alt="userInfo.name || '用户头像'"
                />
                <el-upload
                  class="avatar-uploader"
                  action="#"
                  :show-file-list="false"
                  :before-upload="handleAvatarChange"
                  accept="image/*"
                >
                  <div class="uploader-content">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"/>
                        <circle cx="12" cy="13" r="4"/>
                    </svg>
                    <span>更换头像</span>
                  </div>
                </el-upload>
              </div>
            </div>
            
            <div class="profile-info">
              <h2 class="profile-name">{{ userInfo.name || userInfo.userName || '用户' }}</h2>
              <p class="profile-description">{{ userInfo.description || 'DevOps Asset 是一款现代化的运维资产管理平台' }}</p>
            </div>
          </div>
        </div>

        <div class="profile-details">
          <div class="detail-section">
            <h3 class="section-title">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                <circle cx="12" cy="7" r="4"/>
              </svg>
              联系信息
            </h3>
            <div class="detail-items">
              <div class="detail-item" v-if="userInfo.email">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/>
                  <polyline points="22,6 12,13 2,6"/>
                </svg>
                <span>{{ userInfo.email }}</span>
              </div>
              <div class="detail-item" v-if="userInfo.phone">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"/>
                </svg>
                <span>{{ userInfo.phone }}</span>
              </div>
              <div class="detail-item" v-if="userInfo.roleName">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M9 12l2 2 4-4"/>
                  <path d="M21 4H8l-7 8 7 8h13a2 2 0 0 0 2-2V6a2 2 0 0 0-2-2z"/>
                </svg>
                <span>{{ userInfo.roleName }}</span>
              </div>
              <div class="detail-item" v-if="userInfo.departmentName">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>
                  <polyline points="9,22 9,12 15,12 15,22"/>
                </svg>
                <span>{{ userInfo.departmentName }}</span>
              </div>
            </div>
          </div>

          <div class="detail-section">
            <h3 class="section-title">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"/>
                <line x1="7" y1="7" x2="7.01" y2="7"/>
              </svg>
              标签
            </h3>
            <div class="tags-container user-tags">
              <span 
                v-for="(tag, index) in userTags" 
                :key="tag" 
                :class="[
                  'tag',
                  { 'role-tag': tag === userInfo.roleName },
                  { 'department-tag': tag === userInfo.departmentName }
                ]"
              >
                {{ tag }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧表单区域 -->
      <div class="settings-panel">
        <!-- 基本设置 -->
        <div class="settings-card">
          <div class="card-header">
            <div class="header-title">
              <span>基本设置</span>
            </div>
            <div class="header-description">管理您的基本信息</div>
          </div>

          <div class="card-content">
            <ElForm
              :model="form"
              ref="ruleFormRef"
              :rules="rules"
              class="settings-form"
              label-position="top"
            >
              <ElRow :gutter="24">
                <ElCol :span="12">
                  <ElFormItem label="姓名" prop="name">
                    <ElInput 
                      v-model="form.name" 
                      :disabled="!isEdit"
                      placeholder="请输入姓名"
                    />
                  </ElFormItem>
                </ElCol>
                <ElCol :span="12">
                  <ElFormItem label="性别" prop="gender">
                    <ElSelect 
                      v-model="form.gender" 
                      placeholder="请选择性别" 
                      :disabled="!isEdit"
                    >
                      <ElOption label="男" :value="1" />
                      <ElOption label="女" :value="2" />
                    </ElSelect>
                  </ElFormItem>
                </ElCol>
              </ElRow>
              <ElRow :gutter="24">
                <ElCol :span="12">
                   <ElFormItem label="用户名" prop="username">
                    <ElInput 
                      v-model="form.username" 
                      disabled 
                      placeholder="用户名"
                    />
                  </ElFormItem>
                </ElCol>
                <ElCol :span="12">
                  <ElFormItem label="邮箱" prop="email">
                    <ElInput 
                      v-model="form.email" 
                      :disabled="!isEdit"
                      placeholder="请输入邮箱地址"
                    />
                  </ElFormItem>
                </ElCol>
              </ElRow>
              <ElRow :gutter="24">
                <ElCol :span="12">
                  <ElFormItem label="手机" prop="phone">
                    <ElInput 
                      v-model="form.phone" 
                      :disabled="!isEdit"
                      placeholder="请输入手机号码"
                    />
                  </ElFormItem>
                </ElCol>
                <ElCol :span="12">
                  <ElFormItem label="地址" prop="address">
                    <ElInput 
                      v-model="form.address" 
                      :disabled="!isEdit"
                      placeholder="请输入地址"
                    />
                  </ElFormItem>
                </ElCol>
              </ElRow>
              <ElFormItem label="个人介绍" prop="description">
                <ElInput
                  type="textarea"
                  :rows="4"
                  v-model="form.description"
                  :disabled="!isEdit"
                  placeholder="请输入个人介绍，让大家更好地了解您..."
                  maxlength="200"
                  show-word-limit
                />
              </ElFormItem>
              <ElFormItem>
                <ElButton 
                  type="primary"
                  @click="handleEdit"
                  :loading="loading"
                  class="form-btn"
                >
                  {{ isEdit ? '保存修改' : '编辑信息' }}
                </ElButton>
              </ElFormItem>
            </ElForm>
          </div>
        </div>

        <!-- 更改密码 -->
        <div class="settings-card">
          <div class="card-header">
            <div class="header-title">
              <span>更改密码</span>
            </div>
            <div class="header-description">确保您的账户安全</div>
          </div>

          <div class="card-content">
            <ElForm 
              :model="pwdForm" 
              ref="pwdFormRef"
              :rules="pwdRules"
              class="settings-form"
              label-position="top"
            >
              <ElFormItem label="当前密码" prop="oldPassword">
                <ElInput
                  v-model="pwdForm.oldPassword"
                  type="password"
                  :disabled="!isEditPwd"
                  show-password
                  placeholder="请输入当前密码"
                />
              </ElFormItem>
              <ElFormItem label="新密码" prop="newPassword">
                <ElInput
                  v-model="pwdForm.newPassword"
                  type="password"
                  :disabled="!isEditPwd"
                  show-password
                  placeholder="请输入新密码"
                />
              </ElFormItem>
              <ElFormItem label="确认新密码" prop="confirmPassword">
                <ElInput
                  v-model="pwdForm.confirmPassword"
                  type="password"
                  :disabled="!isEditPwd"
                  show-password
                  placeholder="请再次输入新密码"
                />
              </ElFormItem>
              <ElFormItem>
                <ElButton 
                  type="primary"
                  @click="handleEditPwd"
                  :loading="pwdLoading"
                  class="form-btn"
                >
                  {{ isEditPwd ? '保存密码' : '修改密码' }}
                </ElButton>
              </ElFormItem>
            </ElForm>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useUserStore } from '@/store/modules/user'
import { ElForm, ElMessage, FormInstance, FormRules } from 'element-plus'
import { getUserInfo, updateUserInfo, changePassword } from '@/api/system/api'
import { ApiStatus } from '@/utils/http/status'

defineOptions({ name: 'UserCenter' })

const userStore = useUserStore()
const defaultAvatar = new URL('@/assets/img/user/avatar.webp', import.meta.url).href

const isEdit = ref(false)
const isEditPwd = ref(false)
const loading = ref(false)
const pwdLoading = ref(false)
const userInfo = ref<any>({})

const form = reactive({
  name: '',
  username: '',
  email: '',
  phone: '',
  gender: 1,
  address: '',
  description: ''
})

const pwdForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const ruleFormRef = ref<FormInstance>()
const pwdFormRef = ref<FormInstance>()

const rules = reactive<FormRules>({
  name: [
    { required: true, message: '请输入姓名', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ]
})

const pwdRules = reactive<FormRules>({
  oldPassword: [
    { required: true, message: '请输入当前密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== pwdForm.newPassword) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
})

const userTags = computed(() => {
  const tags = ['DevOps工程师', '云原生架构', '容器化部署', 'CI/CD专家', '基础设施即代码', '微服务架构']
  if (userInfo.value.roleName) {
    tags.push(userInfo.value.roleName)
  }
  if (userInfo.value.departmentName) {
    tags.push(userInfo.value.departmentName)
  }
  return tags
})

// 获取用户信息
const fetchUserInfo = async () => {
  try {
    const response = await getUserInfo()
    if (response.code === ApiStatus.success && response.data) {
      userInfo.value = response.data
      
      // 填充表单
      form.name = userInfo.value.name || ''
      form.username = userInfo.value.userName || ''
      form.email = userInfo.value.email || ''
      form.phone = userInfo.value.phone || ''
      form.gender = userInfo.value.gender || 1
      form.address = userInfo.value.address || ''
      form.description = userInfo.value.description || ''
    }
  } catch (error) {
    console.error('获取用户信息失败:', error)
    ElMessage.error('获取用户信息失败')
  }
}

// 处理头像上传
const handleAvatarChange = (file: File) => {
  // 验证文件类型
  const isImage = file.type.startsWith('image/')
  if (!isImage) {
    ElMessage.error('只能上传图片文件')
    return false
  }
  
  // 验证文件大小
  const isLt2M = file.size / 1024 / 1024 < 2
  if (!isLt2M) {
    ElMessage.error('头像大小不能超过 2MB')
    return false
  }
  
  // 创建预览
  const reader = new FileReader()
  reader.onload = e => {
    userInfo.value.avatar = e.target?.result as string
  }
  reader.readAsDataURL(file)
  
  // 这里可以添加上传到服务器的逻辑
  ElMessage.success('头像上传成功')
  return false // 阻止自动上传
}

// 处理编辑/保存
const handleEdit = async () => {
  if (!isEdit.value) {
    isEdit.value = true
    return
  }
  
  if (!ruleFormRef.value) return
  
  await ruleFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const data = {
          name: form.name,
          phone: form.phone,
          gender: form.gender,
          email: form.email,
          address: form.address,
          description: form.description
        }
        
        const response = await updateUserInfo(data)
        if (response.code === ApiStatus.success) {
          ElMessage.success('保存成功')
          isEdit.value = false
          await fetchUserInfo() // 重新获取用户信息
        } else {
          ElMessage.error(response.message || '保存失败')
        }
      } catch (error) {
        console.error('保存失败:', error)
        ElMessage.error('保存失败')
      } finally {
        loading.value = false
      }
    }
  })
}

// 处理密码编辑/保存
const handleEditPwd = async () => {
  if (!isEditPwd.value) {
    isEditPwd.value = true
    return
  }
  
  if (!pwdFormRef.value) return
  
  await pwdFormRef.value.validate(async (valid) => {
    if (valid) {
      pwdLoading.value = true
      try {
        const data = {
          old_password: pwdForm.oldPassword,
          new_password: pwdForm.newPassword
        }
        
        const response = await changePassword(data)
        if (response.code === ApiStatus.success) {
          ElMessage.success('密码修改成功')
          isEditPwd.value = false
          pwdForm.oldPassword = ''
          pwdForm.newPassword = ''
          pwdForm.confirmPassword = ''
        } else {
          ElMessage.error(response.message || '密码修改失败')
        }
      } catch (error) {
        console.error('密码修改失败:', error)
        ElMessage.error('密码修改失败')
      } finally {
        pwdLoading.value = false
      }
    }
  })
}

onMounted(() => {
  fetchUserInfo()
})
</script>

<style lang="scss" scoped>
.user-center {
  padding: 20px;
  background: #f0f2f5;
}

.user-center-grid {
  display: grid;
  grid-template-columns: 350px 1fr;
  gap: 20px;
  align-items: start;
}

/* Left user profile card */
.user-profile-card {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;

  &:hover {
    transform: translateY(-5px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
  }

  .profile-header {
    position: relative;
    .profile-bg {
      height: 150px;
      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
      .bg-overlay {
        position: absolute;
        inset: 0;
        background: linear-gradient(135deg, rgba(102, 126, 234, 0.5), rgba(118, 75, 162, 0.5));
      }
    }
    .profile-avatar-section {
      position: relative;
      margin-top: -50px;
      padding: 0 24px 24px;
      text-align: center;
      
      .avatar-wrapper {
        .avatar-container {
          position: relative;
          margin: 0 auto 16px;

          .avatar-image {
            width: 100px;
            height: 100px;
            border-radius: 50%;
            border: 4px solid #fff;
            box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
            background-color: #f0f2f5;
            display: block;
          }

          .avatar-uploader {
            position: absolute;
            inset: 0;
            border-radius: 50%;
            background-color: rgba(0, 0, 0, 0.4);
            opacity: 0;
            transition: opacity 0.3s ease;
            cursor: pointer;

            .uploader-content {
              display: flex;
              flex-direction: column;
              align-items: center;
              justify-content: center;
              color: #fff;
              height: 100%;
              gap: 4px;

              svg {
                width: 24px;
                height: 24px;
              }
              span {
                font-size: 13px;
              }
            }

            // 让 el-upload 充满整个区域
            :deep(.el-upload) {
              width: 100%;
              height: 100%;
            }
          }

          &:hover .avatar-uploader {
            opacity: 1;
          }
        }
      }
      .profile-info {
        margin-top: 16px;
        .profile-name {
          font-size: 20px;
          font-weight: 600;
          color: #1a1a1a;
        }
        .profile-description {
          font-size: 14px;
          color: #666;
          margin-top: 4px;
        }
      }
    }
  }

  .profile-details {
    padding: 0 24px 24px;
    .detail-section {
      margin-top: 24px;
      .section-title {
        font-size: 16px;
        font-weight: 600;
        color: #333;
        margin-bottom: 16px;
        display: flex;
        align-items: center;
        gap: 8px;
        svg { width: 18px; height: 18px; color: #1677ff; }
      }
      .detail-items .detail-item {
        display: flex;
        align-items: center;
        gap: 12px;
        font-size: 14px;
        color: #555;
        padding: 8px 0;
        svg { width: 16px; height: 16px; color: #888; flex-shrink: 0; }
        span { word-break: break-all; }
      }
      .tags-container {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;
        .tag {
          padding: 4px 10px;
          border-radius: 12px;
          font-size: 12px;
          font-weight: 500;
          background-color: #f0f5ff;
          color: #1677ff;
          border: 1px solid #d6e4ff;
        }
      }
    }
  }
}

/* Right settings panel */
.settings-panel {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.settings-card {
  background: #fff;
  padding: 24px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);

  .card-header {
    display: flex;
    align-items: baseline;
    gap: 12px;
    padding-bottom: 16px;
    margin-bottom: 24px;
    border-bottom: 1px solid #f0f0f0;

    .header-title span {
      font-size: 18px;
      font-weight: 600;
      color: #1d2129;
    }
    .header-description {
      font-size: 14px;
      color: #86909c;
    }
  }

  .settings-form {
    max-width: 100%;
    
    :deep(.el-form-item__label) {
      font-weight: 500;
      color: #4e5969;
      padding-bottom: 8px;
    }
    :deep(.el-input__wrapper),
    :deep(.el-textarea__inner),
    :deep(.el-select .el-input__wrapper) {
      background-color: #f7f8fa;
      box-shadow: none;
      border: 1px solid transparent;
      transition: all 0.2s ease;
      min-height: 40px;
      border-radius: 6px;
      &:hover {
        border-color: #a9a9a9;
      }
      &.is-focus {
        background-color: #fff;
        border-color: #1677ff;
        box-shadow: 0 0 0 2px rgba(22, 119, 255, 0.2);
      }
    }
    .form-btn {
      min-width: 100px;
      height: 40px;
    }
  }
}

/* Responsive */
@media (max-width: 1200px) {
  .user-center-grid {
    grid-template-columns: 1fr;
  }
}
</style>
