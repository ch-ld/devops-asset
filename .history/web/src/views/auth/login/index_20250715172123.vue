<template>
  <div class="login-bg-animate">
    <svg
      class="bg-anim-svg"
      width="100vw"
      height="100vh"
      viewBox="0 0 1920 1080"
      preserveAspectRatio="none"
    >
      <circle cx="300" cy="300" r="40" fill="#5D87FF22">
        <animate attributeName="r" values="40;60;40" dur="4s" repeatCount="indefinite" />
      </circle>
      <circle cx="1600" cy="200" r="30" fill="#38C0FC22">
        <animate attributeName="r" values="30;50;30" dur="5s" repeatCount="indefinite" />
      </circle>
      <circle cx="500" cy="800" r="25" fill="#5D87FF11">
        <animate attributeName="r" values="25;40;25" dur="6s" repeatCount="indefinite" />
      </circle>
      <line x1="300" y1="300" x2="1600" y2="200" stroke="#5D87FF33" stroke-width="2">
        <animate attributeName="x2" values="1600;1400;1600" dur="6s" repeatCount="indefinite" />
      </line>
      <line x1="500" y1="800" x2="300" y2="300" stroke="#38C0FC33" stroke-width="2">
        <animate attributeName="x2" values="300;400;300" dur="5s" repeatCount="indefinite" />
      </line>
      <circle cx="1200" cy="700" r="18" fill="#38C0FC22">
        <animate attributeName="r" values="18;32;18" dur="7s" repeatCount="indefinite" />
      </circle>
      <line x1="1200" y1="700" x2="1600" y2="200" stroke="#5D87FF22" stroke-width="2">
        <animate attributeName="y2" values="200;300;200" dur="7s" repeatCount="indefinite" />
      </line>
    </svg>
    <div class="login-center-card">
      <div class="logo-row">
        <ArtLogo class="logo-svg" :size="48" />
        <span class="logo-title">DevOps Asset</span>
      </div>
      <div class="slogan">统一的 DevOps 资产管理平台</div>
      <ElForm
        ref="formRef"
        :model="formData"
        :rules="rules"
        @keyup.enter="handleSubmit"
        class="login-form"
      >
        <ElFormItem prop="username">
          <ElInput placeholder="请输入账号" size="large" v-model.trim="formData.username" />
        </ElFormItem>
        <ElFormItem prop="password">
          <ElInput
            placeholder="请输入密码"
            size="large"
            v-model.trim="formData.password"
            type="password"
            autocomplete="off"
            show-password
          />
        </ElFormItem>
        <ElFormItem prop="captcha">
          <ElRow :gutter="5">
            <ElCol :span="16">
              <ElInput placeholder="请输入验证码" size="large" v-model.trim="formData.captcha" />
            </ElCol>
            <ElCol :push="1" :span="8">
              <img :src="captchaImageUrl" @click="refreshCaptcha" class="captcha-image" />
            </ElCol>
          </ElRow>
        </ElFormItem>
        <div class="login-options">
          <ElCheckbox v-model="formData.rememberPassword">记住密码</ElCheckbox>
          <router-link to="/auth/forget-password">忘记密码</router-link>
        </div>
        <ElButton
          class="login-btn"
          size="large"
          type="primary"
          @click="handleSubmit"
          :loading="loading"
        >
          登录
        </ElButton>
      </ElForm>
      <div class="login-footer">
        <a href="https://your-docs.com/devops-asset" target="_blank">文档</a> ·
        <a href="https://github.com/your-repo/devops-asset" target="_blank">Github</a> ·
        <a href="https://your-homepage.com" target="_blank">官网</a>
        <br />
        Copyright © 2024 DevOps Asset
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
    username: '',
    password: '',
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
      // 这里需要修改，后端API状态码为200而不是20000
      if (captchaRes.code === 200 && captchaRes.data) {
        captchaImageUrl.value = captchaRes.data.image
        captchaImageID.value = captchaRes.data.id
        console.log('验证码加载成功:', {
          id: captchaImageID.value,
          imageLength: captchaImageUrl.value?.length
        })
      } else {
        console.error('Invalid captcha response format:', captchaRes)
        ElMessage.error('验证码格式错误')
      }
    } catch (error) {
      console.error('Error refreshing captcha:', error)
      ElMessage.error('验证码获取失败')
    }
  }

  onMounted(() => {
    refreshCaptcha() // 页面加载时获取验证码
  })
</script>

<style lang="scss" scoped>
  @use './index';
</style>
