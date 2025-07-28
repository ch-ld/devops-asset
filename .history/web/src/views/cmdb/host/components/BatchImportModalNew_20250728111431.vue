<template>
  <el-dialog
    v-model="dialogVisible"
    title=""
    width="900px"
    :close-on-click-modal="false"
    class="batch-import-dialog"
    :show-close="false"
  >
    <!-- 自定义头部 -->
    <template #header>
      <div class="dialog-header">
        <div class="header-left">
          <div class="header-icon">
            <el-icon size="24"><Upload /></el-icon>
          </div>
          <div class="header-content">
            <h2 class="dialog-title">批量导入主机</h2>
            <p class="dialog-subtitle">快速导入多台主机到系统中</p>
          </div>
        </div>
        <el-button 
          type="text" 
          size="large" 
          @click="handleCancel"
          class="close-btn"
        >
          <el-icon size="20"><Close /></el-icon>
        </el-button>
      </div>
    </template>

    <div class="import-container">
      <!-- 进度指示器 -->
      <div class="progress-indicator">
        <div class="steps-wrapper">
          <div 
            v-for="(step, index) in steps" 
            :key="index"
            :class="['step-item', { 
              active: currentStep === index, 
              completed: currentStep > index 
            }]"
          >
            <div class="step-circle">
              <el-icon v-if="currentStep > index" class="check-icon">
                <Check />
              </el-icon>
              <span v-else>{{ index + 1 }}</span>
            </div>
            <div class="step-content">
              <div class="step-title">{{ step.title }}</div>
              <div class="step-desc">{{ step.desc }}</div>
            </div>
            <div v-if="index < steps.length - 1" class="step-line"></div>
          </div>
        </div>
      </div>

      <!-- 步骤内容 -->
      <div class="step-content-wrapper">
        <!-- 步骤 0: 选择文件 -->
        <div v-if="currentStep === 0" class="step-panel">
          <div class="panel-header">
            <h3>选择导入文件</h3>
            <p>支持 Excel (.xlsx) 和 CSV (.csv) 格式文件</p>
          </div>

          <div class="upload-section">
            <el-upload
              class="upload-dragger"
              drag
              :auto-upload="false"
              :on-change="handleFileChange"
              :show-file-list="false"
              accept=".xlsx,.csv"
              :class="{ 'has-file': selectedFile }"
            >
              <div class="upload-content">
                <div class="upload-icon">
                  <el-icon size="48" color="#409EFF">
                    <UploadFilled />
                  </el-icon>
                </div>
                <div class="upload-text">
                  <div class="primary-text">点击或拖拽文件到此处</div>
                  <div class="secondary-text">支持 .xlsx / .csv 格式，文件大小不超过 10MB</div>
                </div>
              </div>
            </el-upload>

            <!-- 文件信息 -->
            <div v-if="selectedFile" class="file-info">
              <div class="file-card">
                <div class="file-icon">
                  <el-icon size="24" color="#67C23A">
                    <Document />
                  </el-icon>
                </div>
                <div class="file-details">
                  <div class="file-name">{{ selectedFile.name }}</div>
                  <div class="file-size">{{ formatFileSize(selectedFile.size) }}</div>
                </div>
                <el-button 
                  type="text" 
                  @click="selectedFile = null"
                  class="remove-file"
                >
                  <el-icon><Close /></el-icon>
                </el-button>
              </div>
            </div>
          </div>

          <!-- 模板下载区域 -->
          <div class="template-section">
            <div class="template-card">
              <div class="template-header">
                <div class="template-icon">
                  <el-icon size="20" color="#E6A23C">
                    <Download />
                  </el-icon>
                </div>
                <div class="template-content">
                  <h4>需要导入模板？</h4>
                  <p>下载标准模板，按照格式填写数据后导入</p>
                </div>
              </div>
              <el-button type="primary" plain @click="downloadTemplate">
                下载模板
              </el-button>
            </div>
          </div>
        </div>

        <!-- 步骤 1: 数据预览 -->
        <div v-if="currentStep === 1" class="step-panel">
          <div class="panel-header">
            <h3>数据预览</h3>
            <p>请检查解析的数据是否正确，确认无误后继续导入</p>
          </div>

          <div v-if="previewData.length > 0" class="preview-section">
            <!-- 统计信息 -->
            <div class="stats-cards">
              <div class="stat-card total">
                <div class="stat-icon">
                  <el-icon size="24"><Document /></el-icon>
                </div>
                <div class="stat-content">
                  <div class="stat-number">{{ previewData.length }}</div>
                  <div class="stat-label">总记录数</div>
                </div>
              </div>
              <div class="stat-card valid">
                <div class="stat-icon">
                  <el-icon size="24" color="#67C23A"><CircleCheck /></el-icon>
                </div>
                <div class="stat-content">
                  <div class="stat-number">{{ validCount }}</div>
                  <div class="stat-label">有效记录</div>
                </div>
              </div>
              <div class="stat-card invalid">
                <div class="stat-icon">
                  <el-icon size="24" color="#F56C6C"><CircleClose /></el-icon>
                </div>
                <div class="stat-content">
                  <div class="stat-number">{{ invalidCount }}</div>
                  <div class="stat-label">无效记录</div>
                </div>
              </div>
            </div>

            <!-- 数据表格 -->
            <div class="preview-table-wrapper">
              <el-table
                :data="previewData.slice(0, 10)"
                border
                size="default"
                max-height="400"
                class="preview-table"
                :row-class-name="getRowClassName"
              >
                <el-table-column type="index" label="#" width="50" />
                <el-table-column prop="name" label="主机名称" min-width="120">
                  <template #default="{ row }">
                    <div class="cell-content">
                      <span :class="{ 'error-text': !row.name }">
                        {{ row.name || '缺失' }}
                      </span>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column prop="public_ip" label="公网IP" min-width="120">
                  <template #default="{ row }">
                    <div class="cell-content">
                      <span :class="{ 'error-text': !row.public_ip }">
                        {{ row.public_ip || '缺失' }}
                      </span>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column prop="private_ip" label="私网IP" min-width="120" />
                <el-table-column prop="username" label="用户名" min-width="100">
                  <template #default="{ row }">
                    <div class="cell-content">
                      <span :class="{ 'error-text': !row.username }">
                        {{ row.username || '缺失' }}
                      </span>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column prop="password" label="密码" min-width="100">
                  <template #default="{ row }">
                    <div class="cell-content">
                      <span v-if="row.password" class="password-mask">••••••</span>
                      <span v-else class="error-text">缺失</span>
                    </div>
                  </template>
                </el-table-column>
                <el-table-column prop="os" label="操作系统" min-width="100" />
                <el-table-column prop="provider_type" label="提供商" min-width="100" />
                <el-table-column label="状态" width="80" fixed="right">
                  <template #default="{ row }">
                    <el-tag 
                      :type="row._valid ? 'success' : 'danger'" 
                      size="small"
                      effect="light"
                    >
                      {{ row._valid ? '有效' : '无效' }}
                    </el-tag>
                  </template>
                </el-table-column>
              </el-table>

              <div v-if="previewData.length > 10" class="table-footer">
                <el-text type="info">
                  <el-icon><InfoFilled /></el-icon>
                  仅显示前 10 条记录，共 {{ previewData.length }} 条数据
                </el-text>
              </div>
            </div>

            <!-- 错误提示 -->
            <div v-if="invalidCount > 0" class="error-summary">
              <el-alert
                title="数据验证警告"
                type="warning"
                show-icon
                :closable="false"
              >
                <template #default>
                  <p>发现 {{ invalidCount }} 条无效记录，这些记录将被跳过。</p>
                  <p>请确保必填字段（主机名称、公网IP、用户名、密码）不为空。</p>
                </template>
              </el-alert>
            </div>
          </div>

          <div v-else class="empty-data">
            <el-empty 
              description="没有解析到有效数据"
              :image-size="120"
            >
              <el-button type="primary" @click="currentStep = 0">
                重新选择文件
              </el-button>
            </el-empty>
          </div>
        </div>

        <!-- 步骤 2: 导入结果 -->
        <div v-if="currentStep === 2" class="step-panel">
          <div class="result-section">
            <div class="result-icon">
              <el-icon size="64" color="#67C23A">
                <CircleCheck />
              </el-icon>
            </div>
            <h3 class="result-title">导入完成</h3>
            <p class="result-desc">
              成功导入 {{ importResult?.success || 0 }} 条主机记录
              <span v-if="importResult?.failed > 0">
                ，跳过 {{ importResult.failed }} 条无效记录
              </span>
            </p>
            
            <!-- 导入统计 -->
            <div class="import-stats">
              <div class="stat-item">
                <span class="stat-label">总计：</span>
                <span class="stat-value">{{ importResult?.total || 0 }}</span>
              </div>
              <div class="stat-item success">
                <span class="stat-label">成功：</span>
                <span class="stat-value">{{ importResult?.success || 0 }}</span>
              </div>
              <div v-if="importResult?.failed > 0" class="stat-item failed">
                <span class="stat-label">失败：</span>
                <span class="stat-value">{{ importResult.failed }}</span>
              </div>
            </div>

            <!-- 错误信息 -->
            <div v-if="importResult?.failed_msg?.length" class="error-messages">
              <el-collapse>
                <el-collapse-item title="查看错误详情" name="errors">
                  <ul class="error-list">
                    <li v-for="(error, index) in importResult.failed_msg" :key="index">
                      {{ error }}
                    </li>
                  </ul>
                </el-collapse-item>
              </el-collapse>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部操作栏 -->
    <template #footer>
      <div class="dialog-footer">
        <div class="footer-left">
          <el-button 
            v-if="currentStep > 0 && currentStep < 2" 
            @click="handlePrevious"
            :icon="ArrowLeft"
          >
            上一步
          </el-button>
        </div>
        <div class="footer-right">
          <el-button @click="handleCancel">
            {{ currentStep === 2 ? '关闭' : '取消' }}
          </el-button>
          <el-button
            v-if="currentStep === 0"
            type="primary"
            :disabled="!selectedFile"
            @click="handleNext"
            :icon="ArrowRight"
          >
            下一步
          </el-button>
          <el-button 
            v-if="currentStep === 1" 
            type="primary" 
            :loading="loading" 
            @click="handleImport"
            :icon="Upload"
          >
            {{ loading ? '导入中...' : '确认导入' }}
          </el-button>
          <el-button 
            v-if="currentStep === 2" 
            type="primary" 
            @click="handleFinish"
            :icon="Check"
          >
            完成
          </el-button>
        </div>
      </div>
    </template>
  </el-dialog>
</template>
