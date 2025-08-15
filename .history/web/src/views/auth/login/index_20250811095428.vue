<template>
  <div class="auth-container">
    <!-- 左侧 DevOps 主题插画区域 -->
    <div class="hero-section">
      <!-- 背景渐变 -->
      <div class="hero-background"></div>
      
      <!-- DevOps 插画内容 -->
      <div class="hero-content">
        <div class="hero-title">
          <h1>DevOps Asset 资产管理平台</h1>
          <p>现代化 DevOps 工具链，助力企业数字化转型与高效运维</p>
        </div>
        
        <!-- DevOps 流程插画 -->
        <div class="devops-illustration">
          <svg viewBox="0 0 600 400" class="illustration-svg">
            <defs>
              <linearGradient id="primaryGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                <stop offset="0%" stop-color="#1890ff"/>
                <stop offset="100%" stop-color="#722ed1"/>
              </linearGradient>
              <linearGradient id="secondaryGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                <stop offset="0%" stop-color="#13c2c2"/>
                <stop offset="100%" stop-color="#52c41a"/>
              </linearGradient>
              <filter id="glow" x="-20%" y="-20%" width="140%" height="140%">
                <feGaussianBlur stdDeviation="3" result="blur"/>
                <feColorMatrix in="blur" type="matrix" values="0 0 0 0 0.2  0 0 0 0 0.6  0 0 0 0 1  0 0 0 0.8 0"/>
                <feMerge>
                  <feMergeNode/>
                  <feMergeNode in="SourceGraphic"/>
                </feMerge>
              </filter>
            </defs>
            
            <!-- 开发阶段 -->
            <g class="stage" transform="translate(70, 80)">
              <circle cx="0" cy="0" r="28" fill="url(#primaryGrad)" filter="url(#glow)"/>
              <g fill="white">
                <path d="M-10 -6 L-2 -6 L-2 -12 L10 0 L-2 12 L-2 6 L-10 6 Z"/>
                <circle cx="6" cy="-8" r="1.5"/>
                <circle cx="6" cy="8" r="1.5"/>
              </g>
              <text x="0" y="48" text-anchor="middle" fill="white" font-size="12" font-weight="500">开发</text>
            </g>
            
            <!-- 构建阶段 -->
            <g class="stage" transform="translate(200, 80)">
              <circle cx="0" cy="0" r="28" fill="url(#primaryGrad)" filter="url(#glow)"/>
              <g fill="white">
                <rect x="-8" y="-8" width="16" height="16" rx="2"/>
                <rect x="-6" y="-6" width="12" height="3" fill="url(#primaryGrad)"/>
                <rect x="-6" y="-2" width="8" height="2" fill="url(#primaryGrad)"/>
                <rect x="-6" y="1" width="10" height="2" fill="url(#primaryGrad)"/>
                <rect x="-6" y="4" width="6" height="2" fill="url(#primaryGrad)"/>
              </g>
              <text x="0" y="48" text-anchor="middle" fill="white" font-size="12" font-weight="500">构建</text>
            </g>
            
            <!-- 测试阶段 -->
            <g class="stage" transform="translate(330, 80)">
              <circle cx="0" cy="0" r="28" fill="url(#secondaryGrad)" filter="url(#glow)"/>
              <g stroke="white" stroke-width="2" fill="white">
                <circle cx="0" cy="0" r="9" fill="none"/>
                <path d="M-6 -1 L-2 3 L8 -7" fill="none" stroke-linecap="round"/>
                <circle cx="4" cy="-4" r="1"/>
                <circle cx="-4" cy="4" r="1"/>
              </g>
              <text x="0" y="48" text-anchor="middle" fill="white" font-size="12" font-weight="500">测试</text>
            </g>
            
            <!-- 部署阶段 -->
            <g class="stage" transform="translate(460, 80)">
              <circle cx="0" cy="0" r="28" fill="url(#secondaryGrad)" filter="url(#glow)"/>
              <g fill="white">
                <rect x="-7" y="-7" width="14" height="10" rx="2"/>
                <rect x="-5" y="-5" width="10" height="1.5" fill="url(#secondaryGrad)"/>
                <rect x="-5" y="-2.5" width="10" height="1.5" fill="url(#secondaryGrad)"/>
                <rect x="-5" y="0" width="10" height="1.5" fill="url(#secondaryGrad)"/>
                <circle cx="0" cy="6" r="2"/>
                <path d="M-2 8 L2 8" stroke="white" stroke-width="1"/>
              </g>
              <text x="0" y="48" text-anchor="middle" fill="white" font-size="12" font-weight="500">部署</text>
            </g>
            
            <!-- 连接线 -->
            <path d="M 98 80 Q 149 80 172 80" stroke="rgba(255,255,255,0.6)" stroke-width="2" fill="none" stroke-dasharray="5,5">
              <animate attributeName="stroke-dashoffset" values="0;-10" dur="1s" repeatCount="indefinite"/>
            </path>
            <path d="M 228 80 Q 279 80 302 80" stroke="rgba(255,255,255,0.6)" stroke-width="2" fill="none" stroke-dasharray="5,5">
              <animate attributeName="stroke-dashoffset" values="0;-10" dur="1s" repeatCount="indefinite"/>
            </path>
            <path d="M 358 80 Q 409 80 432 80" stroke="rgba(255,255,255,0.6)" stroke-width="2" fill="none" stroke-dasharray="5,5">
              <animate attributeName="stroke-dashoffset" values="0;-10" dur="1s" repeatCount="indefinite"/>
            </path>
            
            <!-- 服务器集群 -->
            <g class="servers" transform="translate(100, 200)">
              <rect x="0" y="0" width="40" height="60" rx="4" fill="rgba(255,255,255,0.2)" stroke="rgba(255,255,255,0.4)" stroke-width="1"/>
              <rect x="50" y="10" width="40" height="60" rx="4" fill="rgba(255,255,255,0.15)" stroke="rgba(255,255,255,0.3)" stroke-width="1"/>
              <rect x="100" y="5" width="40" height="60" rx="4" fill="rgba(255,255,255,0.1)" stroke="rgba(255,255,255,0.2)" stroke-width="1"/>
              
              <!-- 服务器指示灯 -->
              <circle cx="10" cy="15" r="2" fill="#52c41a"/>
              <circle cx="30" cy="15" r="2" fill="#52c41a"/>
              <circle cx="60" cy="25" r="2" fill="#faad14"/>
              <circle cx="80" cy="25" r="2" fill="#52c41a"/>
              <circle cx="110" cy="20" r="2" fill="#52c41a"/>
              <circle cx="130" cy="20" r="2" fill="#52c41a"/>
            </g>
            
            <!-- 云朵装饰 -->
            <g class="clouds" opacity="0.3">
              <path d="M 450 30 Q 460 20 470 30 Q 480 20 490 30 Q 500 25 510 35 Q 520 30 530 40 Q 525 50 510 50 L 460 50 Q 445 50 450 30" fill="white"/>
              <path d="M 80 300 Q 90 290 100 300 Q 110 290 120 300 Q 130 295 140 305 Q 150 300 160 310 Q 155 320 140 320 L 90 320 Q 75 320 80 300" fill="white"/>
            </g>
            
            <!-- 数据流动效果 -->
            <g class="data-flow">
              <circle r="3" fill="#1890ff" opacity="0.8">
                <animateMotion dur="3s" repeatCount="indefinite" path="M70,80 Q200,60 330,80 Q460,100 570,80"/>
              </circle>
              <circle r="2" fill="#722ed1" opacity="0.6">
                <animateMotion dur="4s" repeatCount="indefinite" path="M70,80 Q200,100 330,80 Q460,60 570,80"/>
              </circle>
              <circle r="2" fill="#13c2c2" opacity="0.7">
                <animateMotion dur="5s" repeatCount="indefinite" path="M70,80 Q265,90 330,80 Q395,70 570,80"/>
              </circle>
            </g>
            
            <!-- 额外装饰元素 -->
            <g class="decorations" opacity="0.3">
              <!-- 代码符号 -->
              <g transform="translate(120, 300)" fill="white" font-family="monospace" font-size="14">
                <text>&lt;/&gt;</text>
              </g>
              <g transform="translate(380, 320)" fill="white" font-family="monospace" font-size="12">
                <text>{ }</text>
              </g>
              <!-- 齿轮装饰 -->
              <g transform="translate(500, 300)" stroke="white" stroke-width="1" fill="none">
                <circle cx="0" cy="0" r="8"/>
                <path d="M-6 -6 L6 6 M6 -6 L-6 6"/>
              </g>
            </g>
          </svg>
        </div>
      </div>
    </div>
    <!-- 旧的动态背景（隐藏） -->
    <div class="background-wrapper hidden-bg">
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
    <!-- 右侧登录卡片 -->
    <div class="auth-card">
      <!-- 卡片头部 -->
      <div class="card-header">
        <div class="pipeline-strip" aria-hidden="true"></div>
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
          <h2 class="welcome-title">欢迎回来</h2>
          <p class="welcome-subtitle">请登录您的账户以继续</p>
        </div>

        <!-- DevOps 无文字视觉（升级版）：渐变流水线 + 节点光环 + 轨迹跑点 -->
        <div class="devops-visual narrow" aria-hidden="true">
          <svg viewBox="0 0 560 110" class="devops-svg">
            <defs>
              <linearGradient id="g1" x1="0%" y1="0%" x2="100%" y2="0%">
                <stop offset="0%" stop-color="#22d1ee"/>
                <stop offset="50%" stop-color="#5d87ff"/>
                <stop offset="100%" stop-color="#00dfd8"/>
              </linearGradient>
              <filter id="softGlow" x="-50%" y="-50%" width="200%" height="200%">
                <feGaussianBlur stdDeviation="3" result="blur"/>
                <feColorMatrix in="blur" type="matrix" values="0 0 0 0 0  0 0 0 0 0.4  0 0 0 0 1  0 0 0 0.7 0"/>
                <feMerge>
                  <feMergeNode/>
                  <feMergeNode in="SourceGraphic"/>
                </feMerge>
              </filter>
            </defs>

            <!-- 流水线主路径 -->
            <path id="pipelinePath" d="M16 66 C 96 18, 150 18, 210 60 S 330 100, 372 60 470 18, 544 60"
                  fill="none" stroke="url(#g1)" stroke-width="2.5" stroke-linecap="round"
                  stroke-dasharray="12 10" class="pipeline-path" />

            <!-- 轨迹跑点 -->
            <circle r="4" fill="url(#g1)" filter="url(#softGlow)">
              <animateMotion dur="8s" repeatCount="indefinite" rotate="auto">
                <mpath href="#pipelinePath" />
              </animateMotion>
            </circle>

            <!-- 节点（光环 + 图标描边）-->
            <g class="node" transform="translate(16,66)">
              <circle r="16" class="node-ring"/>
              <path d="M-6 2h12a6 6 0 0 0 0-12 7 7 0 0 0-13 -2 6 6 0 0 0 1 12z" class="icon-stroke"/>
            </g>
            <g class="node" transform="translate(120,30)">
              <circle r="16" class="node-ring"/>
              <path d="M0-9 l3 2 3-2 3 5 -3 2 0 4 3 2 -3 5 -3-2 -3 2 -3-5 3-2 0-4 -3-2z M0-4 a4 4 0 1 0 0.01 0" class="icon-stroke"/>
            </g>
            <g class="node" transform="translate(220,90)">
              <circle r="16" class="node-ring"/>
              <path d="M-7 -4 h14 v8 h-14 z m0 -6 h14 v4 h-14 z" class="icon-stroke"/>
            </g>
            <g class="node" transform="translate(320,30)">
              <circle r="16" class="node-ring"/>
              <path d="M0-9 L8-5 L8 5 L0 9 L-8 5 L-8 -5 Z" class="icon-stroke"/>
            </g>
            <g class="node" transform="translate(420,90)">
              <circle r="16" class="node-ring"/>
              <circle r="7" class="icon-stroke" fill="none"/>
              <path d="M-7 0 h14 M0 -7 a7 7 0 0 1 0 14 M0 -7 a7 7 0 0 0 0 14" class="icon-stroke"/>
            </g>
            <g class="node" transform="translate(544,60)">
              <circle r="16" class="node-ring"/>
              <path d="M-6 -2 v-3 a6 6 0 0 1 12 0 v3 M-6 -2 h12 v10 h-12 z" class="icon-stroke"/>
            </g>
          </svg>
        </div>
      </div>

      <!-- 登录表单 -->
      <div class="card-body">
        <ElForm
          ref="formRef"
          :model="formData"
          :rules="rules"
          @keyup.enter="handleSubmit"
          class="auth-form"
          size="large"
        >
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
                type="password"
                show-password
                v-model.trim="formData.password"
                class="modern-input"
              />
            </div>
          </ElFormItem>

          <ElFormItem prop="captcha">
            <div class="captcha-row">
              <div class="input-group captcha-input">
                <div class="input-icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                    <rect x="3" y="5" width="18" height="14" rx="2" stroke="currentColor" stroke-width="2"/>
                    <path d="M9 12L11 14L15 10" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                </div>
                <ElInput
                  placeholder="请输入验证码"
                  v-model.trim="formData.captcha"
                  class="modern-input"
                />
              </div>
              <div class="captcha-wrapper">
                <div class="captcha-image-container">
                  <img
                    :src="captchaImageUrl"
                    @click="refreshCaptcha"
                    class="captcha-image"
                    alt="验证码"
                    @error="handleCaptchaImageError"
                    @load="() => {}"
                  />
                </div>
                <div class="captcha-refresh" @click="refreshCaptcha" title="刷新验证码">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none">
                    <path d="M1 4V10H7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    <path d="M23 20V14H17" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                    <path d="M20.49 9A9 9 0 0 0 5.64 5.64L1 10M23 14L18.36 18.36A9 9 0 0 1 3.51 15" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                </div>
              </div>
            </div>
          </ElFormItem>

          <div class="form-options">
            <ElCheckbox v-model="formData.rememberPassword" class="remember-checkbox">
              记住密码
            </ElCheckbox>
            <router-link to="/auth/forget-password" class="forgot-link">忘记密码？</router-link>
          </div>

          <ElFormItem>
            <ElButton
              type="primary"
              size="large"
              @click="handleSubmit"
              class="auth-button"
              :loading="loading"
            >
              {{ loading ? '登录中...' : '立即登录' }}
            </ElButton>
          </ElFormItem>
        </ElForm>
      </div>

      <!-- 卡片底部 -->
      <div class="card-footer">
        <div class="register-link">
          <span>还没有账号？</span>
          <RouterLink to="/auth/register" class="link">立即注册</RouterLink>
        </div>

        <div class="footer-links">
          <a href="https://your-docs.com/devops-asset" target="_blank">文档</a>
          <span>·</span>
          <a href="https://github.com/ch-ld/devops-asset" target="_blank">GitHub</a>
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
  import AppConfig from '@/config'
  import { ElForm, ElMessage, ElNotification } from 'element-plus'
  import { useUserStore } from '@/store/modules/user'
  // ApiStatus已不再使用，直接比较状态码
  // import { ApiStatus } from '@/utils/http/status'
  import { SystemThemeEnum } from '@/enums/appEnum'
  import { useSettingStore } from '@/store/modules/setting'
  import type { FormInstance, FormRules } from 'element-plus'
  import { onMounted, ref, reactive, computed } from 'vue'
  import { getCaptcha, userLogin, getUserInfo } from '@/api/system/api'
  import defaultAvatar from '@/assets/img/user/avatar.png'
  import { storeToRefs } from 'pinia'

  defineOptions({ name: 'Login' })
  const settingStore = useSettingStore()
  const { isDark, systemThemeType } = storeToRefs(settingStore)

  const userStore = useUserStore()
  const router = useRouter()

  // 定义一些未使用的变量为_开头，避免lint警告
  const _systemName = AppConfig.systemInfo.name
  const formRef = ref<FormInstance>()
  const formData = reactive({
    username: import.meta.env.DEV ? 'admin' : '',
    password: import.meta.env.DEV ? '123456' : '',
    rememberPassword: true,
    captcha: ''
  })

  const rules = computed<FormRules>(() => ({
    username: [{ required: true, message: '请输入账号', trigger: 'blur' }],
    password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
    captcha: [{ required: true, message: '请输入验证码', trigger: 'blur' }]
  }))

  const loading = ref(false)
  const captchaImageUrl = ref('') // 验证码图片的URL
  const captchaImageID = ref('') // 验证码图片的ID

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

  const handleSubmit = async () => {
    if (!formRef.value) return

    await formRef.value.validate(async (valid) => {
      if (valid) {
        loading.value = true
        // 延时辅助函数
        const delay = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms))
        try {
          const res = await userLogin({
            username: formData.username,
            password: formData.password,
            captcha: formData.captcha,
            captcha_id: captchaImageID.value
          })
          // 调试查看登录响应
          console.log('Login response:', res)
          // 修改状态码判断为200，而非ApiStatus.success
          if (res.code === 200 && res.data) {
            // 设置 token
            userStore.setToken(res.data.access_token)
            // 获取用户信息
            const userRes = await getUserInfo()
            console.log('User info response:', userRes)
            // 修改状态码判断
            if (userRes.code === 200 && userRes.data) {
              console.log('获取用户信息成功:', userRes.data)

              // 转换为UserInfo类型
              const userInfo: Api.User.UserInfo = {
                userId: userRes.data.id,
                userName: userRes.data.username,
                roles: [userRes.data.role?.code || ''],
                buttons: [],
                avatar: defaultAvatar,
                email: userRes.data.email,
                phone: userRes.data.phone
              }

              userStore.setUserInfo(userInfo)
            } else {
              ElMessage.error(userRes.message)
              console.error('获取用户信息失败:', userRes.message)
            }
            // 设置登录状态
            userStore.setLoginStatus(true)
            // 延时辅助函数
            await delay(1000)
            // 登录成功提示
            showLoginSuccessNotice()
            // 跳转首页
            console.log('登录成功，跳转首页')
            router.push('/')
          } else {
            ElMessage.error(res.message)
            refreshCaptcha()
          }
        } finally {
          await delay(1000)
          loading.value = false
        }
      }
    })
  }

  // 保留简洁品牌与登录体验，去除占位的企业登录入口

  // 登录成功提示
  const showLoginSuccessNotice = () => {
    setTimeout(() => {
      ElNotification({
        title: '登录成功',
        type: 'success',
        showClose: true,
        duration: 2500,
        zIndex: 10000,
        message: `欢迎回来!`
      })
    }, 300)
  }

  // 切换主题
  import { useTheme } from '@/composables/useTheme'

  // 定义但不使用，添加_前缀
  const _toggleTheme = () => {
    let { LIGHT, DARK } = SystemThemeEnum
    useTheme().switchThemeStyles(systemThemeType.value === LIGHT ? DARK : LIGHT)
  }

  const refreshCaptcha = async () => {
    try {
      const captchaRes = await getCaptcha(80, 240)
      // 调试检查响应格式
      console.log('Captcha response:', captchaRes)

      // 全面检查各种可能的响应格式
      if (captchaRes) {
        if (captchaRes.code === 200 && captchaRes.data) {
          // 正常格式
          if (captchaRes.data.image && captchaRes.data.id) {
            captchaImageUrl.value = captchaRes.data.image
            captchaImageID.value = captchaRes.data.id
            console.log('验证码加载成功(标准格式):', {
              id: captchaImageID.value,
              imageLength: captchaImageUrl.value?.length
            })
            return
          }

          // 处理响应中直接包含image和id的情况
          if (typeof captchaRes.data === 'object') {
            const { id, image } = captchaRes.data as { id?: string; image?: string }
            if (id && image) {
              captchaImageUrl.value = image
              captchaImageID.value = id
              console.log('验证码加载成功(直接包含):', { id, imageLength: image.length })
              return
            }
          }
        }

        // 尝试在外层查找id和image
        // 使用类型断言避免类型错误
        const anyResponse = captchaRes as any
        if (anyResponse.image && anyResponse.id) {
          captchaImageUrl.value = anyResponse.image
          captchaImageID.value = anyResponse.id
          console.log('验证码加载成功(外层数据):', { id: anyResponse.id })
          return
        }
      }

      // 如果所有格式都不匹配
      console.error('验证码响应格式无法识别:', captchaRes)
      ElMessage.error('验证码格式错误')
    } catch (error) {
      console.error('Error refreshing captcha:', error)
      ElMessage.error('验证码获取失败')
    }
  }

  // 验证码图片加载错误处理
  const handleCaptchaImageError = () => {
    console.error('验证码图片加载失败，尝试重新获取')
    // 静默重新获取验证码，不显示错误提示
    setTimeout(() => {
      refreshCaptcha()
    }, 500)
  }

  onMounted(() => {
    refreshCaptcha() // 页面加载时获取验证码
  })
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
  background: linear-gradient(135deg, #1890ff 0%, #722ed1 50%, #13c2c2 100%);
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

.devops-illustration {
  width: 100%;
  max-width: 600px;
  height: 400px;
}

.illustration-svg {
  width: 100%;
  height: 100%;
}

.stage {
  cursor: pointer;
  transition: transform 0.3s ease;
}

.stage:hover {
  transform: scale(1.1);
}

.servers rect {
  transition: all 0.3s ease;
}

.servers:hover rect {
  fill: rgba(255, 255, 255, 0.3);
}

.hidden-bg { display: none; }

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
    #0a0f1c 0%,
    #0f1e3a 25%,
    #0b2a6f 50%,
    #0d3b8c 75%,
    #0f52ba 100%
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
  display: flex;
  flex-direction: column;
  justify-content: center;
  height: 100vh;
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
  padding: 48px 48px 24px;
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

.pipeline-strip {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 3px;
  background: linear-gradient(90deg, #22d1ee, #5d87ff, #00dfd8, #22d1ee);
  background-size: 300% 100%;
  animation: flow 3s linear infinite;
}

@keyframes flow {
  0% { background-position: 0% 50%; }
  100% { background-position: 100% 50%; }
}

.brand-section {
  margin-bottom: 30px;
}

.brand-icon {
  margin-bottom: 16px;
  display: flex;
  justify-content: center;
}

.brand-text {
  .brand-title {
    font-size: 28px;
    font-weight: 700;
    color: #1a1a1a;
    margin: 0 0 8px 0;
    background: linear-gradient(135deg, #1890ff, #722ed1);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    position: relative;
    
    &::after {
      content: '';
      position: absolute;
      bottom: -4px;
      left: 50%;
      transform: translateX(-50%);
      width: 0;
      height: 2px;
      background: linear-gradient(135deg, #1890ff, #722ed1);
      border-radius: 1px;
      animation: expandLine 2s ease-out 0.5s forwards;
    }
  }
}

@keyframes expandLine {
  to {
    width: 80px;
  }
}
/* DevOps 无文字视觉 */
.devops-visual { margin-top: 6px; }
.devops-visual.narrow { margin-top: 2px; }
.devops-svg { width: 100%; height: 110px; filter: url(#softGlow); }
.pipeline-path { animation: dashMove 9s linear infinite; opacity: 0.75; }
@keyframes dashMove { to { stroke-dashoffset: -320; } }
.node .node-ring { fill: rgba(99,102,241,0.06); stroke: url(#g1); stroke-width: 1.5; }
.node .icon-stroke { stroke: url(#g1); stroke-width: 1.8; fill: none; }
.node { transition: transform 0.25s ease; }
.node:hover { transform: scale(1.06); }

.welcome-section {
  .welcome-title {
    font-size: 24px;
    font-weight: 600;
    color: #1f2937;
    margin: 0 0 8px 0;
  }

  .welcome-subtitle {
    font-size: 14px;
    color: #6b7280;
    margin: 0;
  }
}

/* 已移除徽章与副标，保持更简洁专业的首屏 */

.logo-icon {
  font-size: 24px;
  color: #6366f1;
}

.logo-text {
  display: flex;
  flex-direction: column;
}

.logo-title {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.logo-version {
  font-size: 12px;
  color: #6b7280;
}

.card-body {
  padding: 0 80px 40px;
  max-width: 400px;
  margin: 0 auto;
}

.auth-form {
  .el-form-item {
    margin-bottom: 24px;
    opacity: 0;
    transform: translateY(20px);
    animation: fadeInUp 0.6s ease-out forwards;

    &:nth-child(1) { animation-delay: 0.1s; }
    &:nth-child(2) { animation-delay: 0.2s; }
    &:nth-child(3) { animation-delay: 0.3s; }
    &:nth-child(4) { animation-delay: 0.4s; }
    &:nth-child(5) { animation-delay: 0.5s; }

    &:last-child {
      margin-bottom: 0;
    }
  }
}

@keyframes fadeInUp {
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.input-group {
  width: 100%; // 保证输入框组占满行宽，左右对齐
  position: relative;
  display: flex;
  align-items: center;
  gap: 12px;
}

.input-group .input-icon {
  position: absolute;
  left: 18px;
  top: 50%;
  transform: translateY(-50%);
  z-index: 2;
  color: #64748b;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  
  svg {
    width: 20px;
    height: 20px;
  }
}

.input-group .modern-input {
  :deep(.el-input__wrapper) {
    padding-left: 52px;
    padding-right: 16px;
    height: 52px;
    border-radius: 16px;
    border: 2px solid #f1f5f9;
    background: #ffffff;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

    &:hover {
      border-color: #e2e8f0;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
      transform: translateY(-1px);
    }

    &.is-focus {
      border-color: #1890ff;
      box-shadow: 
        0 0 0 4px rgba(24, 144, 255, 0.12),
        0 8px 24px rgba(24, 144, 255, 0.15);
      transform: translateY(-2px);
    }
  }

  :deep(.el-input__inner) {
    color: #1f2937;
    font-size: 16px;
    font-weight: 500;
    line-height: 1.5;

    &::placeholder {
      color: #94a3b8;
      font-weight: 400;
    }
  }
}

.input-group:focus-within .input-icon {
  color: #1890ff;
  transform: translateY(-50%) scale(1.1);
}

.captcha-row {
  display: grid; // 采用网格便于对齐
  grid-template-columns: minmax(160px, 1fr) auto; // 取消刷新按钮列
  gap: 12px;
  align-items: center; // 垂直居中
}

.captcha-row .captcha-input {
  width: 100%;
  min-width: 160px; // 最小宽度
}
.captcha-row .captcha-input .modern-input {
  width: 100%;
}

.captcha-wrapper {
  display: flex;
  align-items: center;
  gap: 0; // 去除图片与按钮间距
}

.captcha-wrapper .captcha-image-container {
  position: relative;
  width: 180px; // 加宽验证码图片，实际宽度可根据图片自适应
  height: 48px;
  border-radius: 12px;
  overflow: hidden;
  border: 2px solid #e5e7eb;
  background: #ffffff;
  transition: border-color 0.3s ease;

  &:hover {
    border-color: #6366f1;
  }
}

.captcha-wrapper .captcha-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  cursor: pointer;
  display: block;
  background: #ffffff;
}

.captcha-wrapper .captcha-refresh {
  display: none !important;
}

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;

  .remember-checkbox {
    :deep(.el-checkbox__label) {
      font-size: 14px;
      color: #6b7280;
    }

    :deep(.el-checkbox__input.is-checked .el-checkbox__inner) {
      background-color: #6366f1;
      border-color: #6366f1;
    }
  }

  .forgot-link {
    font-size: 14px;
    color: #6366f1;
    text-decoration: none;
    transition: color 0.3s ease;

    &:hover {
      color: #4f46e5;
      text-decoration: underline;
    }
  }
}

.auth-button {
  width: 100%;
  height: 56px;
  border-radius: 16px;
  font-size: 16px;
  font-weight: 600;
  background: linear-gradient(135deg, #1890ff 0%, #722ed1 100%);
  border: none;
  color: white;
  letter-spacing: 0.5px;
  box-shadow: 
    0 4px 16px rgba(24, 144, 255, 0.3),
    0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: block;
  margin: 0 auto;
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: -100%;
    width: 100%;
    height: 100%;
    background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
    transition: left 0.5s;
  }

  &:hover {
    background: linear-gradient(135deg, #0c7cd5 0%, #5b21b6 100%);
    transform: translateY(-3px);
    box-shadow: 
      0 8px 32px rgba(24, 144, 255, 0.4),
      0 4px 16px rgba(0, 0, 0, 0.15);

    &::before {
      left: 100%;
    }
  }

  &:active {
    transform: translateY(-1px);
    transition-duration: 0.1s;
  }

  &.is-loading {
    background: linear-gradient(135deg, #94a3b8, #64748b);
    cursor: not-allowed;
    
    &:hover {
      transform: none;
      box-shadow: 
        0 4px 16px rgba(24, 144, 255, 0.2),
        0 2px 8px rgba(0, 0, 0, 0.1);
    }
  }
}

.card-footer {
  padding: 24px 48px 48px;
  text-align: center;
  border-top: 1px solid rgba(0, 0, 0, 0.06);
}

.register-link {
  margin-bottom: 24px;
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
  margin-bottom: 16px;
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
    margin: 0;
    border-radius: 16px;
  }

  .card-header {
    padding: 32px 24px 16px;
  }

  .card-body {
    padding: 16px 24px 32px;
  }

  .card-footer {
    padding: 16px 24px 32px;
  }

  .brand-text .brand-title {
    font-size: 24px;
  }

  .welcome-section .welcome-title {
    font-size: 20px;
  }
}

/* 响应式：1024 以下隐藏插画，表单占满 */
@media (max-width: 1024px) {
  .hero-section { display: none; }
  .auth-card {
    width: 100%;
    flex: 1 1 100%;
    max-width: none;
    margin: 0;
    border-radius: 0;
    justify-content: flex-start;
    padding-top: 60px;
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

  .input-group .modern-input :deep(.el-input__wrapper) {
    background: rgba(31, 41, 55, 0.8);
    border-color: #374151;

    &:hover {
      border-color: #4b5563;
    }

    &.is-focus {
      border-color: #6366f1;
    }
  }

  .input-group .modern-input :deep(.el-input__inner) {
    color: #f9fafb;

    &::placeholder {
      color: #9ca3af;
    }
  }

  .captcha-wrapper {
    border-color: #374151;

    &:hover {
      border-color: #4b5563;
    }
  }
}
</style>
