<template>
  <div class="auth-container">
    <!-- 左侧 DevOps 主题插画区域 -->
    <div class="hero-section">
      <!-- 背景渐变 -->
      <div class="hero-background"></div>
      
      <!-- DevOps 插画内容 -->
      <div class="hero-content">
        <div class="hero-title">
          <h1>加入 DevOps Asset 平台</h1>
          <p>开启智能化运维管理，构建高效协作的开发运维生态</p>
        </div>
        
        <!-- 静态装饰图案 -->
        <div class="hero-decoration">
          <div class="decoration-grid">
            <div class="grid-item">
              <div class="icon-circle">
                <svg width="32" height="32" viewBox="0 0 32 32" fill="none">
                  <path d="M16 2L20 12H30L22 18L26 28L16 22L6 28L10 18L2 12H12L16 2Z" fill="rgba(255,255,255,0.8)"/>
                </svg>
              </div>
              <span class="item-text">团队协作</span>
            </div>
            
            <div class="grid-item">
              <div class="icon-circle">
                <svg width="32" height="32" viewBox="0 0 32 32" fill="none">
                  <rect x="6" y="8" width="20" height="16" rx="2" fill="none" stroke="rgba(255,255,255,0.8)" stroke-width="2"/>
                  <path d="M10 12h12M10 16h8M10 20h10" stroke="rgba(255,255,255,0.8)" stroke-width="2"/>
                </svg>
              </div>
              <span class="item-text">代码管理</span>
            </div>
            
            <div class="grid-item">
              <div class="icon-circle">
                <svg width="32" height="32" viewBox="0 0 32 32" fill="none">
                  <circle cx="16" cy="16" r="10" fill="none" stroke="rgba(255,255,255,0.8)" stroke-width="2"/>
                  <path d="M16 6v20M6 16h20" stroke="rgba(255,255,255,0.8)" stroke-width="2"/>
                  <circle cx="16" cy="16" r="3" fill="rgba(255,255,255,0.8)"/>
                </svg>
              </div>
              <span class="item-text">自动化</span>
            </div>
            
            <div class="grid-item">
              <div class="icon-circle">
                <svg width="32" height="32" viewBox="0 0 32 32" fill="none">
                  <rect x="4" y="10" width="24" height="14" rx="2" fill="none" stroke="rgba(255,255,255,0.8)" stroke-width="2"/>
                  <circle cx="8" cy="14" r="1" fill="#52c41a"/>
                  <circle cx="12" cy="14" r="1" fill="#52c41a"/>
                  <circle cx="16" cy="14" r="1" fill="#faad14"/>
                  <path d="M6 18h20" stroke="rgba(255,255,255,0.8)" stroke-width="2"/>
                </svg>
              </div>
              <span class="item-text">监控部署</span>
            </div>
          </div>
          
          <!-- 装饰性几何图形 -->
          <div class="geometric-shapes">
            <div class="shape circle-1"></div>
            <div class="shape circle-2"></div>
            <div class="shape triangle-1"></div>
            <div class="shape square-1"></div>
          </div>
        </div>
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
                  <stop offset="0%" style="stop-color:#0f52ba"/>
                  <stop offset="100%" style="stop-color:#00bcd4"/>
                </linearGradient>
              </defs>
            </svg>
          </div>
          <div class="brand-text">
            <h1 class="brand-title">DevOps Asset</h1>
          </div>
        </div>

        <div class="welcome-section">
          <h2 class="welcome-title">创建新账户</h2>
          <p class="welcome-subtitle">请填写以下信息完成注册</p>
        </div>
        
        <!-- 简洁装饰线条 -->
        <div class="header-decoration" aria-hidden="true">
          <div class="decoration-line">
            <div class="line-segment"></div>
            <div class="line-segment"></div>
            <div class="line-segment"></div>
          </div>
          <div class="decoration-dots">
            <span class="dot active"></span>
            <span class="dot"></span>
            <span class="dot"></span>
            <span class="dot"></span>
          </div>
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
                    <rect x="3" y="11" width="18" height="10" rx="2" ry="2" stroke="currentColor" stroke-width="2"/>
                    <circle cx="12" cy="16" r="1" fill="currentColor"/>
                    <path d="M7 11V7a5 5 0 0 1 10 0v4" stroke="currentColor" stroke-width="2"/>
                  </svg>
                </div>
                <ElInput
                  type="password"
                  placeholder="请输入密码"
                  v-model.trim="formData.password"
                  class="modern-input"
                />
              </div>
            </ElFormItem>

            <ElFormItem prop="confirmPassword">
              <div class="input-group">
                <div class="input-icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                    <rect x="3" y="11" width="18" height="10" rx="2" ry="2" stroke="currentColor" stroke-width="2"/>
                    <path d="M7 11V7a5 5 0 0 1 10 0v4" stroke="currentColor" stroke-width="2"/>
                    <path d="M9 16l2 2 4-4" stroke="currentColor" stroke-width="2"/>
                  </svg>
                </div>
                <ElInput
                  type="password"
                  placeholder="请确认密码"
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
                    <path d="M22 16.92V19.92C22 20.92 21.11 21.92 20 21.92C8.95 21.92 0 12.97 0 1.92C0 0.92 1 0.02 2 0.02H5C6.11 0.02 7 0.92 7 1.92V4.92C7 5.92 6.11 6.92 5 6.92H3C3 12.47 7.53 16.92 13 16.92V14.92C13 13.81 13.89 12.92 15 12.92H18C19.11 12.92 20 13.81 20 14.92V16.92H22Z" fill="currentColor"/>
                  </svg>
                </div>
                <ElInput
                  placeholder="请输入手机号"
                  v-model.trim="formData.phone"
                  class="modern-input"
                />
              </div>
            </ElFormItem>

            <ElFormItem prop="gender">
              <div class="input-group">
                <div class="input-icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                    <circle cx="12" cy="8" r="5" stroke="currentColor" stroke-width="2"/>
                    <path d="M20 21a8 8 0 1 0-16 0" stroke="currentColor" stroke-width="2"/>
                  </svg>
                </div>
                <ElSelect
                  v-model="formData.gender"
                  placeholder="男"
                  class="modern-select"
                  style="width: 100%"
                >
                  <ElOption label="男" :value="1" />
                  <ElOption label="女" :value="2" />
                </ElSelect>
              </div>
            </ElFormItem>
          </div>

          <ElButton
            type="primary"
            size="large"
            :loading="loading"
            @click="handleSubmit"
            class="auth-button"
          >
            立即注册
          </ElButton>
        </ElForm>
      </div>

      <!-- 卡片底部 -->
      <div class="card-footer">
        <div class="register-link">
          <span>已有账号？</span>
          <RouterLink to="/auth/login" class="link">立即登录</RouterLink>
        </div>

        <div class="sso-section">
          <div class="sso-divider">
            <span>或使用以下方式注册</span>
          </div>
          <div class="sso-buttons">
            <button @click="sso('github')" class="sso-btn">
              <span>文档</span>
            </button>
            <button @click="sso('gitlab')" class="sso-btn">
              <span>GitHub</span>
            </button>
            <button @click="sso('ldap')" class="sso-btn">
              <span>官网</span>
            </button>
          </div>
        </div>

        <div class="agreement-section">
          <ElCheckbox v-model="agree" size="small">
            我已阅读并同意
            <a href="#" class="agreement-link">工作区规则</a>
            与
            <a href="#" class="agreement-link">隐私政策</a>
          </ElCheckbox>
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
const agree = ref(true)

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
      if (!agree.value) {
        ElMessage.warning('请先阅读并同意相关协议')
        return
      }
      loading.value = true
      try {
        await registerAPI({
          username: formData.username,
          password: formData.password,
          name: formData.name,
          phone: formData.phone,
          gender: formData.gender
        })
        ElMessage.success('注册成功！正在跳转到登录页...')
        setTimeout(() => {
          router.push('/auth/login')
        }, 1500)
      } catch (error: any) {
        ElMessage.error(error.message || '注册失败，请重试')
      } finally {
        loading.value = false
      }
    }
  })
}

const sso = (provider: 'github' | 'gitlab' | 'ldap') => {
  const map: Record<string, string> = {
    github: 'GitHub 企业注册即将上线',
    gitlab: 'GitLab 企业注册即将上线',
    ldap: 'LDAP/AD 企业注册即将上线'
  }
  ElMessage.info(map[provider] || '敬请期待')
}
</script>

<style lang="scss" scoped>
.auth-container {
  position: relative;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  display: flex;
  align-items: stretch;
  justify-content: stretch;
}

/* 左侧 Hero 区域样式 - 严格50%宽度 */
.hero-section {
  position: relative;
  width: 50%;
  flex: 0 0 50%;
  min-width: 0;
  height: 100vh;
  overflow: hidden;
}

.hero-background {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, #722ed1 0%, #1890ff 50%, #13c2c2 100%);
  background-size: 400% 400%;
  animation: gradientShift 15s ease infinite;
}

@keyframes gradientShift {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

.hero-content {
  position: relative;
  z-index: 2;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 60px 40px;
  color: white;
}

.hero-title {
  text-align: center;
  margin-bottom: 60px;
}

.hero-title h1 {
  font-size: 32px;
  font-weight: 600;
  margin: 0 0 16px 0;
  line-height: 1.4;
  color: white;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.hero-title p {
  font-size: 16px;
  opacity: 0.9;
  margin: 0;
  line-height: 1.6;
  color: rgba(255, 255, 255, 0.85);
}

.hero-decoration {
  position: relative;
  width: 100%;
  height: 400px;
  margin: 60px 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.decoration-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-template-rows: 1fr 1fr;
  gap: 60px;
  max-width: 320px;
  z-index: 2;
  position: relative;
}

.grid-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  opacity: 0;
  transform: translateY(30px);
  animation: fadeInUp 0.8s ease forwards;
}

.grid-item:nth-child(1) { animation-delay: 0.1s; }
.grid-item:nth-child(2) { animation-delay: 0.2s; }
.grid-item:nth-child(3) { animation-delay: 0.3s; }
.grid-item:nth-child(4) { animation-delay: 0.4s; }

.icon-circle {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
  transition: all 0.3s ease;
  cursor: pointer;
}

.icon-circle:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-4px) scale(1.05);
  box-shadow: 0 8px 32px rgba(255, 255, 255, 0.1);
}

.item-text {
  color: rgba(255, 255, 255, 0.9);
  font-size: 14px;
  font-weight: 500;
  letter-spacing: 0.5px;
}

.geometric-shapes {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
  z-index: 1;
}

.shape {
  position: absolute;
  opacity: 0.1;
}

.circle-1 {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  border: 2px solid rgba(255, 255, 255, 0.3);
  top: 20px;
  left: 20px;
  animation: float 6s ease-in-out infinite;
}

.circle-2 {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.05);
  top: 60%;
  right: 10%;
  animation: float 8s ease-in-out infinite reverse;
}

.triangle-1 {
  width: 0;
  height: 0;
  border-left: 30px solid transparent;
  border-right: 30px solid transparent;
  border-bottom: 50px solid rgba(255, 255, 255, 0.08);
  top: 10%;
  right: 20%;
  animation: float 7s ease-in-out infinite;
}

.square-1 {
  width: 40px;
  height: 40px;
  background: rgba(255, 255, 255, 0.06);
  transform: rotate(45deg);
  bottom: 20%;
  left: 15%;
  animation: float 9s ease-in-out infinite reverse;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(-20px);
  }
}

@keyframes fadeInUp {
  0% {
    opacity: 0;
    transform: translateY(30px);
  }
  100% {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 右侧注册卡片样式 */
.auth-card {
  position: relative;
  z-index: 10;
  width: 50%;
  flex: 0 0 50%;
  max-width: 50vw;
  margin: 0;
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  box-shadow:
    0 20px 40px rgba(0, 0, 0, 0.1),
    0 0 0 1px rgba(255, 255, 255, 0.8);
  overflow: hidden;
  animation: cardSlideIn 0.8s ease-out;
  max-height: 100vh;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  justify-content: center;
  height: 100vh;
  padding: 0 80px;
}

@keyframes cardSlideIn {
  0% {
    opacity: 0;
    transform: translateX(50px);
  }
  100% {
    opacity: 1;
    transform: translateX(0);
  }
}

.card-header {
  padding: 60px 0 30px;
  text-align: center;
  position: relative;
  
  &::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 50%;
    transform: translateX(-50%);
    width: 60px;
    height: 2px;
    background: linear-gradient(90deg, transparent, #1890ff, transparent);
    border-radius: 1px;
  }
}

.brand-section {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 32px;
  gap: 16px;
}

.brand-icon {
  display: flex;
  align-items: center;
  justify-content: center;
}

.brand-text .brand-title {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
  position: relative;
  
  &::after {
    content: '';
    position: absolute;
    bottom: -6px;
    left: 0;
    width: 0;
    height: 2px;
    background: linear-gradient(90deg, #1890ff, #722ed1);
    transition: width 0.6s ease;
    animation: expandLine 2s ease 0.5s forwards;
  }
}

@keyframes expandLine {
  0% { width: 0; }
  100% { width: 60px; }
}

.welcome-section {
  margin-bottom: 24px;
}

.welcome-title {
  font-size: 28px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 8px 0;
  line-height: 1.3;
}

.welcome-subtitle {
  font-size: 16px;
  color: #6b7280;
  margin: 0;
}

/* 简洁装饰样式 */
.header-decoration {
  margin-top: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.decoration-line {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  max-width: 200px;
}

.line-segment {
  height: 2px;
  flex: 1;
  background: linear-gradient(90deg, transparent, #1890ff, transparent);
  border-radius: 1px;
  opacity: 0.6;
}

.line-segment:nth-child(2) {
  background: linear-gradient(90deg, transparent, #722ed1, transparent);
  opacity: 0.8;
}

.line-segment:nth-child(3) {
  background: linear-gradient(90deg, transparent, #13c2c2, transparent);
  opacity: 0.6;
}

.decoration-dots {
  display: flex;
  align-items: center;
  gap: 8px;
}

.dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: rgba(24, 144, 255, 0.3);
  transition: all 0.3s ease;
}

.dot.active {
  background: #1890ff;
  box-shadow: 0 0 8px rgba(24, 144, 255, 0.4);
}

.dot:nth-child(2) {
  background: rgba(114, 46, 209, 0.3);
}

.dot:nth-child(3) {
  background: rgba(19, 194, 194, 0.3);
}

.dot:nth-child(4) {
  background: rgba(82, 196, 26, 0.3);
}

.card-body {
  padding: 0 0 40px;
  max-width: 420px;
  margin: 0 auto;
  width: 100%;
}

.auth-form {
  .el-form-item {
    margin-bottom: 20px;
    opacity: 0;
    transform: translateY(20px);
    animation: fadeInUp 0.6s ease forwards;
  }
  
  .el-form-item:nth-child(1) { animation-delay: 0.1s; }
  .el-form-item:nth-child(2) { animation-delay: 0.2s; }
  .el-form-item:nth-child(3) { animation-delay: 0.3s; }
  .el-form-item:nth-child(4) { animation-delay: 0.4s; }
}

.form-row {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
  margin-bottom: 4px;

  @media (max-width: 640px) {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .el-form-item {
    width: 100%;
    margin-bottom: 20px;
  }
}

.input-group {
  position: relative;
  width: 100%;
}

.input-icon {
  position: absolute;
  left: 18px;
  top: 50%;
  transform: translateY(-50%);
  color: #64748b;
  z-index: 3;
  transition: all 0.3s ease;
  
  svg {
    width: 20px;
    height: 20px;
  }
}

.modern-input {
  padding-left: 52px;
  padding-right: 16px;
  height: 50px;
  border-radius: 14px;
  border: 2px solid #f1f5f9;
  background: #fafbfc;
  font-size: 15px;
  font-weight: 500;
  transition: all 0.3s ease;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.02);
  
  &:hover {
    border-color: #e2e8f0;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
    transform: translateY(-1px);
  }
  
  &:focus {
    border-color: #1890ff;
    box-shadow: 0 0 0 4px rgba(24, 144, 255, 0.12), 0 8px 24px rgba(24, 144, 255, 0.15);
    transform: translateY(-2px);
    background: #fff;
  }
  
  .el-input__inner {
    border: none;
    background: transparent;
    font-size: 15px;
    font-weight: 500;
    color: #1f2937;
    
    &::placeholder {
      color: #9ca3af;
      font-weight: 400;
    }
  }
}

.modern-input:focus-within .input-icon {
  color: #1890ff;
  transform: translateY(-50%) scale(1.1);
}

.modern-select {
  .el-select__wrapper {
    padding-left: 52px;
    padding-right: 16px;
    height: 50px;
    border-radius: 14px;
    border: 2px solid #f1f5f9;
    background: #fafbfc;
    font-size: 15px;
    font-weight: 500;
    transition: all 0.3s ease;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.02);
    
    &:hover {
      border-color: #e2e8f0;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
      transform: translateY(-1px);
    }
    
    &.is-focused {
      border-color: #1890ff;
      box-shadow: 0 0 0 4px rgba(24, 144, 255, 0.12), 0 8px 24px rgba(24, 144, 255, 0.15);
      transform: translateY(-2px);
      background: #fff;
    }
  }
  
  .el-select__placeholder {
    color: #9ca3af;
    font-weight: 400;
  }
}

.auth-button {
  width: 100%;
  height: 54px;
  border-radius: 16px;
  font-size: 16px;
  font-weight: 600;
  letter-spacing: 0.5px;
  background: linear-gradient(135deg, #1890ff 0%, #722ed1 100%);
  border: none;
  color: white;
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
  box-shadow: 0 8px 24px rgba(24, 144, 255, 0.3);
  margin: 12px auto 0;
  
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: -100%;
    width: 100%;
    height: 100%;
    background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
    transition: left 0.5s ease;
  }
  
  &:hover {
    background: linear-gradient(135deg, #0c7cd5 0%, #5b21b6 100%);
    transform: translateY(-3px);
    box-shadow: 
      0 12px 32px rgba(24, 144, 255, 0.4),
      0 8px 16px rgba(114, 46, 209, 0.3);
    
    &::before {
      left: 100%;
    }
  }
  
  &:active {
    transform: translateY(-1px);
  }
  
  &.is-loading {
    background: linear-gradient(135deg, #94a3b8, #64748b);
    cursor: not-allowed;
    
    &:hover {
      transform: none;
      box-shadow: 0 8px 24px rgba(24, 144, 255, 0.3);
    }
  }
}

.card-footer {
  padding: 30px 0 60px;
  text-align: center;
  border-top: 1px solid rgba(0, 0, 0, 0.06);
}

.register-link {
  margin-bottom: 20px;
  font-size: 14px;
  color: #6b7280;
  
  .link {
    color: #1890ff;
    text-decoration: none;
    font-weight: 500;
    margin-left: 4px;
    
    &:hover {
      color: #0c7cd5;
      text-decoration: underline;
    }
  }
}

.sso-section {
  margin-bottom: 24px;
}

.sso-divider {
  position: relative;
  margin-bottom: 16px;
  
  &::before {
    content: '';
    position: absolute;
    top: 50%;
    left: 0;
    right: 0;
    height: 1px;
    background: #e5e7eb;
    z-index: 1;
  }
  
  span {
    position: relative;
    z-index: 2;
    background: rgba(255, 255, 255, 0.98);
    padding: 0 16px;
    font-size: 12px;
    color: #9ca3af;
  }
}

.sso-buttons {
  display: flex;
  gap: 12px;
  justify-content: center;
}

.sso-btn {
  flex: 1;
  max-width: 80px;
  height: 36px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background: #fff;
  color: #6b7280;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &:hover {
    border-color: #1890ff;
    color: #1890ff;
    background: #f0f9ff;
    transform: translateY(-1px);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }
}

.agreement-section {
  margin-bottom: 20px;
  
  .el-checkbox {
    font-size: 12px;
    color: #6b7280;
  }
  
  .agreement-link {
    color: #1890ff;
    text-decoration: none;
    
    &:hover {
      text-decoration: underline;
    }
  }
}

.copyright {
  font-size: 12px;
  color: #9ca3af;
  margin-top: 16px;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .hero-section {
    display: none;
  }
  
  .auth-card {
    width: 100%;
    flex: 1 1 100%;
    max-width: none;
    margin: 0;
    border-radius: 0;
    justify-content: flex-start;
    padding-top: 40px;
  }
  
  .card-header {
    padding: 24px 20px 12px;
  }
  
  .card-body {
    padding: 0 20px 20px;
  }
  
  .card-footer {
    padding: 12px 20px 24px;
  }
  
  .brand-text .brand-title {
    font-size: 20px;
  }
  
  .welcome-title {
    font-size: 24px;
  }
}

@media (max-width: 768px) {
  .auth-card {
    background: rgba(17, 24, 39, 0.95);
    box-shadow:
      0 8px 32px rgba(0, 0, 0, 0.3),
      0 0 0 1px rgba(255, 255, 255, 0.1);
  }
  
  .brand-text .brand-title {
    color: #f9fafb;
  }
  
  .welcome-title {
    color: #f9fafb;
  }
  
  .welcome-subtitle {
    color: #d1d5db;
  }
}
</style>
