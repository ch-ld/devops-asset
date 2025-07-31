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
      <!-- 注册表单 -->
      <div class="card-body">
        <ElForm
          ref="formRef"
          :model="formData"
          :rules="rules"
          @keyup.enter="handleSubmit"
          class="auth-form"
          size="large"
        >
          <div class="form-row">
            <ElFormItem prop="username">
              <div class="input-group">
                <div class="input-icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                    <path d="M12 12C14.7614 12 17 9.76142 17 7C17 4.23858 14.7614 2 12 2C9.23858 2 7 4.23858 7 7C7 9.76142 9.23858 12 12 12Z" fill="currentColor"/>
                    <path d="M12 14C7.58172 14 4 17.5817 4 22H20C20 17.5817 16.4183 14 12 14Z" fill="currentColor"/>
                  </svg>
                </div>
                <ElInput
                  placeholder="请输入用户名"
                  v-model.trim="formData.username"
                  class="modern-input"
                />
              </div>
            </ElFormItem>

            <ElFormItem prop="name">
              <div class="input-group">
                <div class="input-icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                    <path d="M20 21V19C20 17.9391 19.5786 16.9217 18.8284 16.1716C18.0783 15.4214 17.0609 15 16 15H8C6.93913 15 5.92172 15.4214 5.17157 16.1716C4.42143 16.9217 4 17.9391 4 19V21" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    <circle cx="12" cy="7" r="4" stroke="currentColor" stroke-width="2"/>
                  </svg>
                </div>
                <ElInput
                  placeholder="请输入真实姓名"
                  size="large"
                  v-model.trim="formData.name"
                  class="modern-input"
                />
              </div>
            </ElFormItem>
          </div>

          <div class="form-row">
            <ElFormItem prop="password">
              <div class="input-group">
                <div class="input-icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                    <path d="M6 10V8C6 5.79086 7.79086 4 10 4H14C16.2091 4 18 5.79086 18 8V10" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                    <rect x="4" y="10" width="16" height="10" rx="2" fill="currentColor"/>
                  </svg>
                </div>
                <ElInput
                  placeholder="请输入密码"
                  size="large"
                  type="password"
                  show-password
                  v-model.trim="formData.password"
                  class="modern-input"
                />
              </div>
            </ElFormItem>

            <ElFormItem prop="confirmPassword">
              <div class="input-group">
                <div class="input-icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                    <path d="M6 10V8C6 5.79086 7.79086 4 10 4H14C16.2091 4 18 5.79086 18 8V10" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                    <rect x="4" y="10" width="16" height="10" rx="2" fill="currentColor"/>
                    <path d="M9 14L11 16L15 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                </div>
                <ElInput
                  placeholder="请确认密码"
                  size="large"
                  type="password"
                  show-password
                  v-model.trim="formData.confirmPassword"
                  class="modern-input"
                />
              </div>
            </ElFormItem>
          </div>

          <div class="form-row">
            <ElFormItem prop="phone">
              <div class="input-group">
                <div class="input-icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                    <path d="M22 16.92V19.92C22.0011 20.1985 21.9441 20.4742 21.8325 20.7293C21.7209 20.9845 21.5573 21.2136 21.3521 21.4019C21.1468 21.5901 20.9046 21.7335 20.6407 21.8227C20.3769 21.9119 20.0974 21.9451 19.82 21.92C16.7428 21.5856 13.787 20.5341 11.19 18.85C8.77382 17.3147 6.72533 15.2662 5.18999 12.85C3.49997 10.2412 2.44824 7.27099 2.11999 4.18C2.095 3.90347 2.12787 3.62476 2.21649 3.36162C2.30512 3.09849 2.44756 2.85669 2.63476 2.65162C2.82196 2.44655 3.0498 2.28271 3.30379 2.17052C3.55777 2.05833 3.83233 2.00026 4.10999 2H7.10999C7.59531 1.99522 8.06579 2.16708 8.43376 2.48353C8.80173 2.79999 9.04207 3.23945 9.10999 3.72C9.23662 4.68007 9.47144 5.62273 9.80999 6.53C9.94454 6.88792 9.97366 7.27691 9.8939 7.65088C9.81415 8.02485 9.62886 8.36811 9.35999 8.64L8.08999 9.91C9.51355 12.4135 11.5865 14.4864 14.09 15.91L15.36 14.64C15.6319 14.3711 15.9751 14.1858 16.3491 14.1061C16.7231 14.0263 17.1121 14.0555 17.47 14.19C18.3773 14.5286 19.3199 14.7634 20.28 14.89C20.7658 14.9585 21.2094 15.2032 21.5265 15.5775C21.8437 15.9518 22.0122 16.4296 21.1 16.92H22Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                </div>
                <ElInput
                  placeholder="请输入手机号"
                  size="large"
                  v-model.trim="formData.phone"
                  class="modern-input"
                />
              </div>
            </ElFormItem>

            <ElFormItem prop="gender">
              <div class="input-group">
                <div class="input-icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                    <circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="2"/>
                    <path d="M12 1V3" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                    <path d="M12 21V23" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                    <path d="M4.22 4.22L5.64 5.64" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                    <path d="M18.36 18.36L19.78 19.78" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                  </svg>
                </div>
                <ElSelect
                  placeholder="请选择性别"
                  size="large"
                  v-model="formData.gender"
                  class="modern-select"
                >
                  <ElOption label="男" :value="1" />
                  <ElOption label="女" :value="2" />
                </ElSelect>
              </div>
            </ElFormItem>
          </div>

          <ElFormItem>
            <ElButton
              type="primary"
              size="large"
              :loading="loading"
              @click="handleSubmit"
              class="auth-button"
            >
              {{ loading ? '注册中...' : '立即注册' }}
            </ElButton>
          </ElFormItem>
        </ElForm>
      </div>

      <!-- 卡片底部 -->
      <div class="card-footer">
        <div class="register-link">
          <span>已有账号？</span>
          <RouterLink to="/auth/login" class="link">立即登录</RouterLink>
        </div>

        <div class="footer-links">
          <a href="https://your-docs.com/devops-asset" target="_blank">文档</a>
          <span>·</span>
          <a href="https://github.com/your-repo/devops-asset" target="_blank">Github</a>
          <span>·</span>
          <a href="https://your-homepage.com" target="_blank">官网</a>
        </div>

        <div class="copyright">
          Copyright © 2024 DevOps Asset
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage, type FormRules } from 'element-plus'
import { useRouter } from 'vue-router'
import { registerAPI } from '@/api/system/api'

const router = useRouter()
const formRef = ref()
const loading = ref(false)

// 粒子动画样式生成
const getParticleStyle = (index: number) => {
  const size = Math.random() * 4 + 2
  const duration = Math.random() * 20 + 10
  const delay = Math.random() * 5
  const x = Math.random() * 100
  const y = Math.random() * 100

  return {
    width: `${size}px`,
    height: `${size}px`,
    left: `${x}%`,
    top: `${y}%`,
    animationDuration: `${duration}s`,
    animationDelay: `${delay}s`
  }
}

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

<style lang="scss" scoped>
.auth-container {
  position: relative;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.background-wrapper {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1;
}

.gradient-bg {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg,
    #0f0f23 0%,
    #1a1a2e 25%,
    #16213e 50%,
    #0f3460 75%,
    #533483 100%
  );
  background-size: 400% 400%;
  animation: gradientShift 20s ease infinite;
}

@keyframes gradientShift {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

.particles {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.particle {
  position: absolute;
  background: rgba(255, 255, 255, 0.6);
  border-radius: 50%;
  animation: float linear infinite;
}

@keyframes float {
  0% {
    transform: translateY(100vh) rotate(0deg);
    opacity: 0;
  }
  10% {
    opacity: 1;
  }
  90% {
    opacity: 1;
  }
  100% {
    transform: translateY(-100px) rotate(360deg);
    opacity: 0;
  }
}

.geometric-shapes {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.shape {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  animation: shapeFloat 20s ease-in-out infinite;
}

.shape-1 {
  width: 200px;
  height: 200px;
  top: 10%;
  left: 10%;
  animation-delay: 0s;
}

.shape-2 {
  width: 150px;
  height: 150px;
  top: 60%;
  right: 15%;
  animation-delay: 5s;
}

.shape-3 {
  width: 100px;
  height: 100px;
  bottom: 20%;
  left: 20%;
  animation-delay: 10s;
}

.shape-4 {
  width: 120px;
  height: 120px;
  top: 30%;
  right: 30%;
  animation-delay: 15s;
}

@keyframes shapeFloat {
  0%, 100% { transform: translateY(0px) rotate(0deg); }
  50% { transform: translateY(-20px) rotate(180deg); }
}

.tech-grid {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  opacity: 0.4;
  background-image:
    linear-gradient(rgba(99, 102, 241, 0.1) 1px, transparent 1px),
    linear-gradient(90deg, rgba(99, 102, 241, 0.1) 1px, transparent 1px);
  background-size: 50px 50px;
  animation: gridMove 30s linear infinite;
}

@keyframes gridMove {
  0% { transform: translate(0, 0); }
  100% { transform: translate(50px, 50px); }
}

.auth-card {
  position: relative;
  z-index: 10;
  width: 100%;
  max-width: 520px;
  margin: 0 20px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  box-shadow:
    0 8px 32px rgba(0, 0, 0, 0.1),
    0 0 0 1px rgba(255, 255, 255, 0.2);
  overflow: hidden;
  animation: cardSlideIn 0.8s ease-out;
  max-height: 90vh;
  overflow-y: auto;
}

@keyframes cardSlideIn {
  0% {
    opacity: 0;
    transform: translateY(30px) scale(0.95);
  }
  100% {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.card-header {
  padding: 32px 32px 16px;
  text-align: center;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.brand-section {
  margin-bottom: 24px;
}

.brand-icon {
  margin-bottom: 12px;
  display: flex;
  justify-content: center;
}

.brand-text {
  .brand-title {
    font-size: 24px;
    font-weight: 700;
    color: #1a1a1a;
    margin: 0 0 6px 0;
    background: linear-gradient(135deg, #6366f1, #8b5cf6);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .brand-subtitle {
    font-size: 13px;
    color: #6b7280;
    margin: 0;
  }
}

.welcome-section {
  .welcome-title {
    font-size: 20px;
    font-weight: 600;
    color: #1f2937;
    margin: 0 0 6px 0;
  }

  .welcome-subtitle {
    font-size: 13px;
    color: #6b7280;
    margin: 0;
  }
}

.card-body {
  padding: 16px 32px 32px;
}

.auth-form {
  .el-form-item {
    margin-bottom: 20px;

    &:last-child {
      margin-bottom: 0;
    }
  }
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;

  @media (max-width: 640px) {
    grid-template-columns: 1fr;
    gap: 0;
  }
}

.input-group {
  position: relative;
  display: flex;
  align-items: center;

  .input-icon {
    position: absolute;
    left: 16px;
    z-index: 2;
    color: #9ca3af;
    transition: color 0.3s ease;
  }

  .modern-input {
    :deep(.el-input__wrapper) {
      padding-left: 48px;
      padding-right: 16px;
      border-radius: 12px;
      border: 2px solid #e5e7eb;
      background: #ffffff;
      box-shadow: none;
      transition: all 0.3s ease;

      &:hover {
        border-color: #d1d5db;
      }

      &.is-focus {
        border-color: #6366f1;
        box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
      }
    }

    :deep(.el-input__inner) {
      font-size: 15px;
      color: #1f2937;

      &::placeholder {
        color: #9ca3af;
      }
    }
  }

  .modern-select {
    width: 100%;

    :deep(.el-select__wrapper) {
      padding-left: 48px;
      padding-right: 16px;
      border-radius: 12px;
      border: 2px solid #e5e7eb;
      background: #ffffff;
      box-shadow: none;
      transition: all 0.3s ease;

      &:hover {
        border-color: #d1d5db;
      }

      &.is-focus {
        border-color: #6366f1;
        box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
      }
    }

    :deep(.el-select__selected-item) {
      font-size: 15px;
      color: #1f2937;
    }

    :deep(.el-select__placeholder) {
      color: #9ca3af;
      font-size: 15px;
    }

    :deep(.el-select__caret) {
      color: #9ca3af;
    }
  }

  &:focus-within .input-icon {
    color: #6366f1;
  }
}

.auth-button {
  width: 100%;
  height: 48px;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  border: none;
  transition: all 0.3s ease;
  margin-top: 8px;

  &:hover {
    background: linear-gradient(135deg, #5b5bd6, #7c3aed);
    transform: translateY(-2px);
    box-shadow: 0 8px 25px rgba(99, 102, 241, 0.3);
  }

  &:active {
    transform: translateY(0);
  }

  &.is-loading {
    background: linear-gradient(135deg, #9ca3af, #6b7280);
  }
}

.card-footer {
  padding: 16px 32px 32px;
  text-align: center;
  border-top: 1px solid rgba(0, 0, 0, 0.05);
}

.register-link {
  margin-bottom: 20px;
  font-size: 14px;
  color: #6b7280;

  .link {
    color: #6366f1;
    text-decoration: none;
    font-weight: 500;
    margin-left: 4px;
    transition: color 0.3s ease;

    &:hover {
      color: #4f46e5;
      text-decoration: underline;
    }
  }
}

.footer-links {
  margin-bottom: 12px;
  font-size: 13px;
  color: #9ca3af;

  a {
    color: #6b7280;
    text-decoration: none;
    transition: color 0.3s ease;

    &:hover {
      color: #6366f1;
    }
  }

  span {
    margin: 0 8px;
  }
}

.copyright {
  font-size: 12px;
  color: #9ca3af;
}

// 响应式设计
@media (max-width: 768px) {
  .auth-card {
    margin: 0 16px;
    border-radius: 16px;
    max-width: none;
  }

  .card-header {
    padding: 24px 20px 12px;
  }

  .card-body {
    padding: 12px 20px 24px;
  }

  .card-footer {
    padding: 12px 20px 24px;
  }

  .brand-text .brand-title {
    font-size: 20px;
  }

  .welcome-section .welcome-title {
    font-size: 18px;
  }
}

// 暗色主题支持
@media (prefers-color-scheme: dark) {
  .auth-card {
    background: rgba(17, 24, 39, 0.95);
    box-shadow:
      0 8px 32px rgba(0, 0, 0, 0.3),
      0 0 0 1px rgba(255, 255, 255, 0.1);
  }

  .brand-text .brand-title {
    color: #f9fafb;
  }

  .welcome-section .welcome-title {
    color: #f9fafb;
  }

  .input-group .modern-input :deep(.el-input__wrapper),
  .input-group .modern-select :deep(.el-select__wrapper) {
    background: rgba(31, 41, 55, 0.8);
    border-color: #374151;

    &:hover {
      border-color: #4b5563;
    }

    &.is-focus {
      border-color: #6366f1;
    }
  }

  .input-group .modern-input :deep(.el-input__inner),
  .input-group .modern-select :deep(.el-select__selected-item) {
    color: #f9fafb;

    &::placeholder {
      color: #9ca3af;
    }
  }
}
</style>
