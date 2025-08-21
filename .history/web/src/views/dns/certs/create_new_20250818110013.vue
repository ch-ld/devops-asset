<template>
  <div class="cert-create">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-left">
          <el-button 
            text 
            @click="goBack" 
            class="back-btn"
            size="large"
          >
            <el-icon><ArrowLeft /></el-icon>
            返回
          </el-button>
          <div class="header-title">
            <h1>申请SSL证书</h1>
            <p>为您的域名申请免费的Let's Encrypt SSL证书，保护网站安全</p>
          </div>
        </div>
        <div class="header-actions">
          <el-button @click="goBack" size="large">取消</el-button>
          <el-button 
            type="primary" 
            :loading="applying" 
            @click="handleApply"
            :disabled="!canSubmit"
            size="large"
          >
            <el-icon><Lock /></el-icon>
            申请证书
          </el-button>
        </div>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="main-content">
      <!-- 步骤指示器 -->
      <div class="steps-container">
        <el-steps :active="currentStep" align-center>
          <el-step title="选择申请方式" description="选择证书申请类型"></el-step>
          <el-step title="配置证书信息" description="填写证书相关信息"></el-step>
          <el-step title="完成申请" description="确认信息并提交申请"></el-step>
        </el-steps>
      </div>

      <!-- 表单内容 -->
      <div class="form-container">
        <el-form 
          ref="formRef" 
          :model="formData" 
          :rules="formRules" 
          label-width="140px"
          size="large"
          class="cert-form"
        >
          <!-- 第一步：申请方式选择 -->
          <div class="step-section" v-show="currentStep >= 0">
            <div class="section-header">
              <h3>
                <el-icon><Setting /></el-icon>
                选择申请方式
              </h3>
              <p>请选择您希望的证书申请方式</p>
            </div>
            
            <el-form-item prop="applyType">
              <div class="apply-type-cards">
                <div 
                  class="apply-type-card"
                  :class="{ active: formData.applyType === 'auto' }"
                  @click="selectApplyType('auto')"
                >
                  <div class="card-icon">
                    <el-icon><Lightning /></el-icon>
                  </div>
                  <div class="card-content">
                    <div class="card-title">
                      自动申请
                      <el-tag type="success" size="small">推荐</el-tag>
                    </div>
                    <div class="card-desc">
                      使用Let's Encrypt免费申请SSL证书，支持自动续期
                    </div>
                    <div class="card-features">
                      <span>• 完全免费</span>
                      <span>• 自动续期</span>
                      <span>• 90天有效期</span>
                    </div>
                  </div>
                  <div class="card-radio">
                    <el-radio v-model="formData.applyType" label="auto"></el-radio>
                  </div>
                </div>

                <div 
                  class="apply-type-card"
                  :class="{ active: formData.applyType === 'upload' }"
                  @click="selectApplyType('upload')"
                >
                  <div class="card-icon">
                    <el-icon><FolderAdd /></el-icon>
                  </div>
                  <div class="card-content">
                    <div class="card-title">导入证书</div>
                    <div class="card-desc">
                      导入已有的SSL证书文件
                    </div>
                    <div class="card-features">
                      <span>• 支持现有证书</span>
                      <span>• 快速导入</span>
                      <span>• 灵活管理</span>
                    </div>
                  </div>
                  <div class="card-radio">
                    <el-radio v-model="formData.applyType" label="upload"></el-radio>
                  </div>
                </div>

                <div 
                  class="apply-type-card"
                  :class="{ active: formData.applyType === 'csr' }"
                  @click="selectApplyType('csr')"
                >
                  <div class="card-icon">
                    <el-icon><Upload /></el-icon>
                  </div>
                  <div class="card-content">
                    <div class="card-title">自定义CSR</div>
                    <div class="card-desc">
                      上传CSR文件申请证书
                    </div>
                    <div class="card-features">
                      <span>• 自定义配置</span>
                      <span>• 高级选项</span>
                      <span>• 专业用户</span>
                    </div>
                  </div>
                  <div class="card-radio">
                    <el-radio v-model="formData.applyType" label="csr"></el-radio>
                  </div>
                </div>
              </div>
            </el-form-item>
          </div>

          <!-- 第二步：自动申请配置 -->
          <div class="step-section" v-if="formData.applyType === 'auto'">
            <div class="section-header">
              <h3>
                <el-icon><Connection /></el-icon>
                配置证书信息
              </h3>
              <p>请选择域名和证书类型</p>
            </div>

            <!-- 域名选择 -->
            <el-form-item label="选择域名" prop="domainId" required>
              <el-select 
                v-model="formData.domainId" 
                placeholder="请选择要申请证书的域名"
                filterable
                @change="handleDomainChange"
                style="width: 100%"
                size="large"
              >
                <el-option 
                  v-for="domain in domainOptions" 
                  :key="domain.id"
                  :label="domain.name" 
                  :value="domain.id"
                >
                  <div class="domain-option">
                    <div class="domain-info">
                      <span class="domain-name">{{ domain.name }}</span>
                      <el-tag 
                        v-if="domain.provider" 
                        type="primary" 
                        size="small"
                        effect="light"
                      >
                        {{ domain.provider.name }}
                      </el-tag>
                    </div>
                    <div class="domain-status">
                      <el-icon color="#67c23a"><Connection /></el-icon>
                      <span>已配置DNS</span>
                    </div>
                  </div>
                </el-option>
              </el-select>
              <div class="form-tip">
                <el-icon><InfoFilled /></el-icon>
                只显示已配置DNS Provider的域名，确保能够自动验证域名所有权
              </div>
            </el-form-item>

            <!-- 证书类型 -->
            <el-form-item label="证书类型" prop="certType">
              <div class="cert-type-cards">
                <div 
                  class="cert-type-card"
                  :class="{ active: formData.certType === 'single' }"
                  @click="selectCertType('single')"
                >
                  <div class="type-icon">
                    <el-icon><Document /></el-icon>
                  </div>
                  <div class="type-content">
                    <div class="type-title">单域名证书</div>
                    <div class="type-desc">仅保护选择的域名</div>
                    <div class="type-example">例如：example.com</div>
                  </div>
                  <div class="type-radio">
                    <el-radio v-model="formData.certType" label="single"></el-radio>
                  </div>
                </div>

                <div 
                  class="cert-type-card"
                  :class="{ active: formData.certType === 'wildcard' }"
                  @click="selectCertType('wildcard')"
                >
                  <div class="type-icon">
                    <el-icon><Connection /></el-icon>
                  </div>
                  <div class="type-content">
                    <div class="type-title">通配符证书</div>
                    <div class="type-desc">保护域名及其所有子域名</div>
                    <div class="type-example">例如：*.example.com</div>
                  </div>
                  <div class="type-radio">
                    <el-radio v-model="formData.certType" label="wildcard"></el-radio>
                  </div>
                </div>

                <div 
                  class="cert-type-card"
                  :class="{ active: formData.certType === 'multi' }"
                  @click="selectCertType('multi')"
                >
                  <div class="type-icon">
                    <el-icon><Plus /></el-icon>
                  </div>
                  <div class="type-content">
                    <div class="type-title">多域名证书</div>
                    <div class="type-desc">保护多个不同的域名</div>
                    <div class="type-example">例如：多个域名</div>
                  </div>
                  <div class="type-radio">
                    <el-radio v-model="formData.certType" label="multi"></el-radio>
                  </div>
                </div>
              </div>
            </el-form-item>

            <!-- 邮箱地址 -->
            <el-form-item label="邮箱地址" prop="email" required>
              <el-input 
                v-model="formData.email" 
                placeholder="请输入邮箱地址，用于接收证书相关通知"
                size="large"
              >
                <template #prefix>
                  <el-icon><Message /></el-icon>
                </template>
              </el-input>
              <div class="form-tip">
                <el-icon><InfoFilled /></el-icon>
                邮箱用于接收证书到期提醒和重要通知
              </div>
            </el-form-item>
          </div>
        </el-form>
      </div>
    </div>
  </div>
</template>
