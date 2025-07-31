<template>
  <div class="auth-container">
    <!-- 动态背景 -->
    <div class="background-wrapper">
      <!-- 渐变背景 -->
      <div class="gradient-bg"></div>

      <!-- 浮动粒子效果 -->
      <div class="particles">
        <div class="particle" v-for="i in 50" :key="i" :style="getParticleStyle(i)"></div>
      </div>

      <!-- 几何图形装饰 -->
      <div class="geometric-shapes">
        <div class="shape shape-1"></div>
        <div class="shape shape-2"></div>
        <div class="shape shape-3"></div>
        <div class="shape shape-4"></div>
      </div>

      <!-- 科技网格 -->
      <div class="tech-grid">
        <svg width="100%" height="100%" class="grid-svg">
          <defs>
            <pattern id="grid" width="60" height="60" patternUnits="userSpaceOnUse">
              <path d="M 60 0 L 0 0 0 60" fill="none" stroke="rgba(99, 102, 241, 0.1)" stroke-width="1"/>
            </pattern>
          </defs>
          <rect width="100%" height="100%" fill="url(#grid)" />
        </svg>
      </div>
    </div>

    <!-- 主注册卡片 -->
    <div class="auth-card">
      <!-- 卡片头部 -->
      <div class="card-header">
        <div class="brand-section">
          <div class="brand-icon">
            <svg width="48" height="48" viewBox="0 0 48 48" fill="none">
              <rect width="48" height="48" rx="12" fill="url(#brandGradient)"/>
              <path d="M16 20L24 14L32 20V32L24 38L16 32V20Z" stroke="white" stroke-width="2" fill="none"/>
              <circle cx="24" cy="24" r="3" fill="white"/>
              <defs>
                <linearGradient id="brandGradient" x1="0%" y1="0%" x2="100%" y2="100%">
                  <stop offset="0%" style="stop-color:#6366f1"/>
                  <stop offset="100%" style="stop-color:#8b5cf6"/>
                </linearGradient>
              </defs>
            </svg>
          </div>
          <div class="brand-text">
            <h1 class="brand-title">DevOps Asset</h1>
            <p class="brand-subtitle">统一的 DevOps 资产管理平台</p>
          </div>
        </div>

        <div class="welcome-section">
          <h2 class="welcome-title">创建新账户</h2>
          <p class="welcome-subtitle">请填写以下信息完成注册</p>
        </div>
      </div>
      <ElForm
        ref="formRef"
        :model="formData"
        :rules="rules"
        @keyup.enter="handleSubmit"
        class="login-form register-form"
      >
        <ElFormItem prop="username">
          <ElInput placeholder="请输入用户名" size="large" v-model.trim="formData.username" />
        </ElFormItem>
        <ElFormItem prop="password">
          <ElInput
            placeholder="请输入密码"
            size="large"
            type="password"
            show-password
            v-model.trim="formData.password"
          />
        </ElFormItem>
        <ElFormItem prop="confirmPassword">
          <ElInput
            placeholder="请确认密码"
            size="large"
            type="password"
            show-password
            v-model.trim="formData.confirmPassword"
          />
        </ElFormItem>
        <ElFormItem prop="name">
          <ElInput placeholder="请输入真实姓名" size="large" v-model.trim="formData.name" />
        </ElFormItem>
        <ElFormItem prop="phone">
          <ElInput placeholder="请输入手机号" size="large" v-model.trim="formData.phone" />
        </ElFormItem>
        <ElFormItem prop="gender">
          <ElSelect placeholder="请选择性别" size="large" v-model="formData.gender" style="width: 100%">
            <ElOption label="男" :value="1" />
            <ElOption label="女" :value="2" />
          </ElSelect>
        </ElFormItem>
        <ElFormItem>
          <ElButton
            type="primary"
            size="large"
            :loading="loading"
            @click="handleSubmit"
            class="login-btn"
          >
            注册
          </ElButton>
        </ElFormItem>
      </ElForm>
      <div class="footer-links">
        <RouterLink to="/auth/login" class="link">已有账号？立即登录</RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage, type FormRules } from 'element-plus'
import { useRouter } from 'vue-router'
import ArtLogo from '@/components/core/base/art-logo/index.vue'
import { registerAPI } from '@/api/system/api'

const router = useRouter()
const formRef = ref()
const loading = ref(false)

// 表单数据
const formData = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  name: '',
  phone: '',
  gender: 1
})

// 验证规则
const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 50, message: '用户名长度在 3 到 50 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 50, message: '密码长度在 6 到 50 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== formData.password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  name: [
    { required: true, message: '请输入真实姓名', trigger: 'blur' },
    { min: 2, max: 50, message: '姓名长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号格式', trigger: 'blur' }
  ],
  gender: [{ required: true, message: '请选择性别', trigger: 'change' }]
}

// 提交注册
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      loading.value = true
      try {
        await registerAPI({
          username: formData.username,
          password: formData.password,
          name: formData.name,
          phone: formData.phone,
          gender: formData.gender
        })
        
        ElMessage.success('注册成功！请登录')
        router.push('/auth/login')
      } catch (error: any) {
        ElMessage.error(error.message || '注册失败')
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped lang="scss">
@import '../login/index.scss';

.register-card {
  min-height: 600px;
  width: 420px;
}

.register-form {
  .el-form-item {
    margin-bottom: 20px;
  }
}

.footer-links {
  text-align: center;
  margin-top: 20px;
  
  .link {
    color: var(--el-color-primary);
    text-decoration: none;
    font-size: 14px;
    
    &:hover {
      text-decoration: underline;
    }
  }
}
</style>
