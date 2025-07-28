<template>
  <div class="host-management">
    <div class="page-container">
      <!-- 左侧边栏 -->
      <div class="sidebar">
        <div class="sidebar-header">
          <div class="header-content">
            <h3 class="title">
              <i class="icon">📁</i>
              主机分组
            </h3>
            <div class="header-actions">
              <el-button size="small" @click="handleAddGroup" title="新建分组">
                ➕
              </el-button>
              <el-button size="small" type="primary" @click="refreshGroups" :loading="groupLoading">
                刷新
              </el-button>
            </div>
          </div>
        </div>
        
        <!-- 搜索框 -->
        <div class="search-section">
          <el-input
            v-model="groupSearchQuery"
            placeholder="搜索分组..."
            size="small"
            clearable
            class="search-input"
          >
            <template #prefix>
              <i class="search-icon">🔍</i>
            </template>
          </el-input>
        </div>

        <!-- 分组列表 -->
        <div class="group-list">
          <div 
            v-for="group in filteredGroups" 
            :key="group.id"
            :class="['group-item', { active: selectedGroupId === group.id }]"
            @click="handleGroupSelect(group)"
          >
            <div class="group-content" :style="{ paddingLeft: `${(group.level || 0) * 20}px` }">
              <div class="group-info">
                <i class="group-icon">{{ group.id === null ? '🏠' : '📂' }}</i>
                <span class="group-name">{{ group.name }}</span>
              </div>
              <div class="group-count">{{ group.host_count || 0 }}</div>
            </div>
            <div class="group-actions" v-if="group.id !== null">
              <el-button size="small" text @click.stop="handleEditGroup(group)" title="编辑">
                ✏️
              </el-button>
              <el-button size="small" text @click.stop="handleDeleteGroup(group)" title="删除">
                🗑️
              </el-button>
            </div>
          </div>
        </div>
      </div>

      <!-- 主内容区 -->
      <div class="main-content">
        <!-- 页面头部 -->
        <div class="content-header">
          <div class="header-left">
            <h1 class="page-title">主机管理</h1>
            <div class="breadcrumb">
              <span>主机管理</span>
              <span class="separator">/</span>
              <span class="current">{{ currentGroupName }}</span>
            </div>
          </div>
          <div class="header-actions">
            <el-button @click="refreshHosts" :loading="loading">
              刷新数据
            </el-button>
            <el-button @click="showBatchImportDialog" :icon="Upload">
              批量导入
            </el-button>
            <el-button @click="handleBatchExport" :icon="Download">
              批量导出
            </el-button>
            <el-button type="primary" @click="handleAdd">
              添加主机
            </el-button>
          </div>
        </div>

        <!-- 搜索工具栏 -->
        <div class="toolbar">
          <div class="toolbar-left">
            <el-input
              v-model="searchQuery"
              placeholder="搜索主机名称、IP地址..."
              size="default"
              clearable
              class="search-input"
            >
              <template #prefix>
                <i>🔍</i>
              </template>
            </el-input>
          </div>
          <div class="toolbar-right">
            <el-select v-model="statusFilter" placeholder="状态筛选" clearable size="default">
              <el-option
                v-for="status in filterOptions.statuses"
                :key="status"
                :label="getStatusLabel(status)"
                :value="status"
              />
            </el-select>
            <el-select v-model="regionFilter" placeholder="地域筛选" clearable size="default">
              <el-option
                v-for="region in filterOptions.regions"
                :key="region"
                :label="region"
                :value="region"
              />
            </el-select>
            <el-select v-model="providerFilter" placeholder="云厂商筛选" clearable size="default">
              <el-option
                v-for="provider in filterOptions.providers"
                :key="provider"
                :label="getProviderLabel(provider)"
                :value="provider"
              />
            </el-select>
          </div>
        </div>

        <!-- 主机列表 -->
        <div class="table-container">
          <el-table
            :data="hosts"
            v-loading="loading"
            class="host-table"
            stripe
            @selection-change="handleSelectionChange"
          >
            <!-- 选择框 -->
            <el-table-column type="selection" width="50" />
            
            <!-- 主机名称 -->
            <el-table-column prop="name" label="主机名称" min-width="150" show-overflow-tooltip>
              <template #default="{ row }">
                <div class="host-name">
                  <i class="host-icon">🖥️</i>
                  <span>{{ row.name }}</span>
                </div>
              </template>
            </el-table-column>

            <!-- IP地址 -->
            <el-table-column label="IP地址" min-width="180">
              <template #default="{ row }">
                <div class="ip-info">
                  <div class="ip-item" v-if="getDisplayIP(row.public_ip)">
                    <div class="ip-row">
                      <i class="el-icon-globe" style="color: #409eff; margin-right: 4px;"></i>
                      <span class="ip-label">公网:</span>
                      <span class="ip-value">{{ getDisplayIP(row.public_ip) }}</span>
                      <el-button
                        type="text"
                        size="mini"
                        @click="copyToClipboard(getDisplayIP(row.public_ip))"
                        title="复制IP"
                        class="copy-btn"
                      >
                        <i class="el-icon-copy-document"></i>
                      </el-button>
                    </div>
                  </div>
                  <div class="ip-item" v-if="getDisplayIP(row.private_ip)">
                    <div class="ip-row">
                      <i class="el-icon-house" style="color: #67c23a; margin-right: 4px;"></i>
                      <span class="ip-label">私网:</span>
                      <span class="ip-value">{{ getDisplayIP(row.private_ip) }}</span>
                      <el-button
                        type="text"
                        size="mini"
                        @click="copyToClipboard(getDisplayIP(row.private_ip))"
                        title="复制IP"
                        class="copy-btn"
                      >
                        <i class="el-icon-copy-document"></i>
                      </el-button>
                    </div>
                  </div>
                  <div class="ip-item" v-if="!getDisplayIP(row.public_ip) && !getDisplayIP(row.private_ip)">
                    <span class="ip-empty">暂无IP</span>
                  </div>
                </div>
              </template>
            </el-table-column>

            <!-- 状态 -->
            <el-table-column label="状态" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)" size="small">
                  {{ getStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>

            <!-- 配置 -->
            <el-table-column label="配置" width="180" align="center">
              <template #default="{ row }">
                <div class="config-specs">
                  <div class="spec-item">
                    <i class="el-icon-cpu" style="color: #409eff; margin-right: 4px;"></i>
                    <span class="spec-label">CPU:</span>
                    <span class="spec-value">{{ formatConfiguration(row.configuration).cpu }}</span>
                  </div>
                  <div class="spec-item">
                    <i class="el-icon-memory" style="color: #67c23a; margin-right: 4px;"></i>
                    <span class="spec-label">内存:</span>
                    <span class="spec-value">{{ formatConfiguration(row.configuration).memory }}</span>
                  </div>
                  <div class="spec-item">
                    <i class="el-icon-hard-disk" style="color: #e6a23c; margin-right: 4px;"></i>
                    <span class="spec-label">磁盘:</span>
                    <span class="spec-value">{{ formatConfiguration(row.configuration).disk }}</span>
                  </div>
                </div>
              </template>
            </el-table-column>

            <!-- 系统 -->
            <el-table-column prop="os" label="系统" width="120" show-overflow-tooltip>
              <template #default="{ row }">
                <span class="os-info">{{ row.os || '-' }}</span>
              </template>
            </el-table-column>

            <!-- 地域 -->
            <el-table-column prop="region" label="地域" width="120">
              <template #default="{ row }">
                <span class="region-info">{{ row.region || '-' }}</span>
              </template>
            </el-table-column>

            <!-- 云厂商 -->
            <el-table-column label="云厂商" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="getProviderType(row.provider_type)" size="small">
                  {{ getProviderText(row.provider_type) }}
                </el-tag>
              </template>
            </el-table-column>

            <!-- 操作 -->
            <el-table-column label="操作" width="240" align="center" fixed="right">
              <template #default="{ row }">
                <div class="action-buttons">
                  <el-button type="primary" size="small" @click="handleView(row)" title="查看详情">
                    <i class="el-icon-view"></i>
                    查看
                  </el-button>
                  <el-button type="success" size="small" @click="handleTerminal(row)" title="SSH终端">
                    <i class="el-icon-monitor"></i>
                    终端
                  </el-button>
                  <el-button type="warning" size="small" @click="handleEdit(row)" title="编辑主机">
                    <i class="el-icon-edit"></i>
                    编辑
                  </el-button>
                  <el-button type="danger" size="small" @click="handleDelete(row)" title="删除主机">
                    <i class="el-icon-delete"></i>
                    删除
                  </el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <!-- 分页 -->
        <div class="pagination-container">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[5, 10, 15, 20]"
            :total="total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </div>

    <!-- 现代化添加/编辑主机对话框 -->
    <el-dialog
      v-model="addHostDialogVisible"
      :title="null"
      width="80%"
      :close-on-click-modal="false"
      class="modern-host-form-dialog"
      :show-close="false"
    >
      <div class="modern-form-container">
        <!-- 自定义头部 -->
        <div class="modern-form-header">
          <div class="header-left">
            <div class="form-avatar">
              <i class="el-icon-plus"></i>
            </div>
            <div class="form-info">
              <h2 class="form-title">{{ isEditMode ? '编辑主机' : '添加主机' }}</h2>
              <div class="form-subtitle">{{ isEditMode ? '修改主机配置信息' : '创建新的主机实例' }}</div>
            </div>
          </div>
          <div class="header-actions">
            <el-button
              type="info"
              :icon="Close"
              circle
              @click="closeHostDialog"
              class="close-btn"
            />
          </div>
        </div>

        <!-- 表单内容 -->
        <div class="form-content">
          <el-form
            :model="hostFormData"
            :rules="hostRules"
            ref="hostFormRef"
            class="modern-form"
            label-width="120px"
          >
            <!-- 基本信息区域 -->
            <div class="form-section">
              <div class="section-header">
                <div class="section-icon basic-icon">
                  <i class="el-icon-monitor"></i>
                </div>
                <div class="section-title">
                  <h3>基本信息</h3>
                  <p>主机的基本标识和分类信息</p>
                </div>
              </div>
              <div class="section-content">
                <el-row :gutter="24">
                  <el-col :span="8">
                    <el-form-item label="主机名称" prop="name">
                      <el-input
                        v-model="hostFormData.name"
                        placeholder="请输入主机名称"
                        size="large"
                      />
                    </el-form-item>
                  </el-col>
                  <el-col :span="8">
                    <el-form-item label="实例ID">
                      <el-input
                        v-model="hostFormData.instance_id"
                        placeholder="留空自动生成"
                        size="large"
                      >
                        <template #append>
                          <el-button @click="generateInstanceId" type="primary" size="small">
                            生成
                          </el-button>
                        </template>
                      </el-input>
                    </el-form-item>
                  </el-col>
                  <el-col :span="8">
                    <el-form-item label="主机类型" prop="host_type">
                      <el-select
                        v-model="hostFormData.host_type"
                        placeholder="请选择主机类型"
                        size="large"
                        style="width: 100%"
                        @change="handleHostTypeChange"
                      >
                        <el-option label="云服务器" value="cloud">
                          <div class="provider-option">
                            <span class="provider-icon">☁️</span>
                            <span>云服务器</span>
                          </div>
                        </el-option>
                        <el-option label="自建服务器" value="self-hosted">
                          <div class="provider-option">
                            <span class="provider-icon">🖥️</span>
                            <span>自建服务器</span>
                          </div>
                        </el-option>
                      </el-select>
                    </el-form-item>
                  </el-col>
                  <el-col :span="8">
                    <el-form-item label="主机组" prop="host_group">
                      <el-select
                        v-model="hostFormData.host_group"
                        placeholder="请选择主机组"
                        size="large"
                        style="width: 100%"
                      >
                        <el-option
                          v-for="group in hostGroups"
                          :key="group.id"
                          :label="group.name"
                          :value="group.id"
                        >
                          <div class="provider-option">
                            <span class="provider-icon">📂</span>
                            <span>{{ group.name }}</span>
                          </div>
                        </el-option>
                      </el-select>
                    </el-form-item>
                  </el-col>
                </el-row>

                <el-row :gutter="24" v-if="hostFormData.host_type === 'cloud'">
                  <el-col :span="8">
                    <el-form-item label="云厂商" prop="provider_type">
                      <el-select
                        v-model="hostFormData.provider_type"
                        placeholder="请选择云厂商"
                        size="large"
                        style="width: 100%"
                      >
                        <el-option label="AWS" value="aws">
                          <div class="provider-option">
                            <span class="provider-icon">🟠</span>
                            <span>Amazon Web Services</span>
                          </div>
                        </el-option>
                        <el-option label="阿里云" value="aliyun">
                          <div class="provider-option">
                            <span class="provider-icon">🟠</span>
                            <span>阿里云</span>
                          </div>
                        </el-option>
                        <el-option label="腾讯云" value="tencent">
                          <div class="provider-option">
                            <span class="provider-icon">🔵</span>
                            <span>腾讯云</span>
                          </div>
                        </el-option>
                        <el-option label="华为云" value="huawei">
                          <div class="provider-option">
                            <span class="provider-icon">🔴</span>
                            <span>华为云</span>
                          </div>
                        </el-option>
                      </el-select>
                    </el-form-item>
                  </el-col>
                  <el-col :span="8">
                    <el-form-item label="实例类型" prop="instance_type">
                      <el-input
                        v-model="hostFormData.configuration.instance_type"
                        placeholder="如：t3.medium"
                        size="large"
                      />
                    </el-form-item>
                  </el-col>
                  <el-col :span="8">
                    <el-form-item label="操作系统" prop="os">
                      <el-input
                        v-model="hostFormData.os"
                        placeholder="如：Ubuntu 20.04 LTS"
                        size="large"
                      />
                    </el-form-item>
                  </el-col>
                </el-row>

                <el-row :gutter="24">
                  <el-col :span="8">
                    <el-form-item label="地域" prop="region">
                      <el-input
                        v-model="hostFormData.region"
                        placeholder="如：us-east-1"
                        size="large"
                      />
                    </el-form-item>
                  </el-col>
                </el-row>
              </div>
            </div>

            <!-- 网络配置区域 -->
            <div class="form-section">
              <div class="section-header">
                <div class="section-icon network-icon">
                  <i class="el-icon-connection"></i>
                </div>
                <div class="section-title">
                  <h3>网络配置</h3>
                  <p>主机的网络连接信息</p>
                </div>
              </div>
              <div class="section-content">
                <el-row :gutter="24">
                  <el-col :span="12">
                    <el-form-item label="公网IP" prop="public_ip">
                      <el-input
                        v-model="hostFormData.public_ip[0]"
                        placeholder="请输入公网IP地址"
                        size="large"
                      />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="私网IP" prop="private_ip">
                      <el-input
                        v-model="hostFormData.private_ip[0]"
                        placeholder="请输入私网IP地址"
                        size="large"
                      />
                    </el-form-item>
                  </el-col>
                </el-row>
              </div>
            </div>

            <!-- SSH连接配置区域 -->
            <div class="form-section">
              <div class="section-header">
                <div class="section-icon ssh-icon">
                  <i class="el-icon-key"></i>
                </div>
                <div class="section-title">
                  <h3>SSH连接配置</h3>
                  <p>配置SSH连接信息以便远程管理</p>
                </div>
              </div>
              <div class="section-content">
                <el-row :gutter="24">
                  <el-col :span="8">
                    <el-form-item label="SSH端口" prop="ssh_port">
                      <el-input-number
                        v-model="hostFormData.ssh_config.port"
                        :min="1"
                        :max="65535"
                        size="large"
                        style="width: 100%"
                        placeholder="22"
                      />
                    </el-form-item>
                  </el-col>
                  <el-col :span="8">
                    <el-form-item label="SSH用户名" prop="ssh_username">
                      <el-input
                        v-model="hostFormData.ssh_config.username"
                        placeholder="如：root, ubuntu"
                        size="large"
                      />
                    </el-form-item>
                  </el-col>
                  <el-col :span="8">
                    <el-form-item label="认证方式" prop="auth_type">
                      <el-select
                        v-model="hostFormData.ssh_config.auth_type"
                        placeholder="选择认证方式"
                        size="large"
                        style="width: 100%"
                      >
                        <el-option label="密码认证" value="password" />
                        <el-option label="密钥认证" value="key" />
                      </el-select>
                    </el-form-item>
                  </el-col>
                </el-row>

                <!-- 密码认证 -->
                <el-row :gutter="24" v-if="hostFormData.ssh_config.auth_type === 'password'">
                  <el-col :span="12">
                    <el-form-item label="SSH密码" prop="ssh_config.password">
                      <el-input
                        v-model="hostFormData.ssh_config.password"
                        type="password"
                        placeholder="请输入SSH密码"
                        size="large"
                        show-password
                      />
                    </el-form-item>
                  </el-col>
                </el-row>

                <!-- 密钥认证 -->
                <div v-if="hostFormData.ssh_config.auth_type === 'key'">
                  <el-row :gutter="24">
                    <el-col :span="24">
                      <el-form-item label="私钥内容" prop="ssh_private_key">
                        <el-input
                          v-model="hostFormData.ssh_config.private_key"
                          type="textarea"
                          :rows="6"
                          placeholder="请粘贴SSH私钥内容，或点击下方按钮上传密钥文件"
                          size="large"
                        />
                      </el-form-item>
                    </el-col>
                  </el-row>
                  <el-row :gutter="24">
                    <el-col :span="12">
                      <el-upload
                        class="key-upload"
                        :before-upload="handleKeyUpload"
                        :show-file-list="false"
                        accept=".pem,.key,.pub"
                      >
                        <el-button size="large" type="primary" plain>
                          <i class="el-icon-upload"></i>
                          上传密钥文件
                        </el-button>
                      </el-upload>
                    </el-col>
                    <el-col :span="12">
                      <el-form-item label="密钥密码" prop="key_passphrase">
                        <el-input
                          v-model="hostFormData.ssh_config.passphrase"
                          type="password"
                          placeholder="如果密钥有密码请输入"
                          size="large"
                          show-password
                        />
                      </el-form-item>
                    </el-col>
                  </el-row>
                </div>
              </div>
            </div>

            <!-- 硬件配置区域 -->
            <div class="form-section">
              <div class="section-header">
                <div class="section-icon hardware-icon">
                  <i class="el-icon-cpu"></i>
                </div>
                <div class="section-title">
                  <h3>硬件配置</h3>
                  <p>主机的硬件资源信息</p>
                </div>
              </div>
              <div class="section-content">
                <el-row :gutter="24">
                  <el-col :span="8">
                    <el-form-item label="CPU核数" prop="cpu_cores">
                      <el-input-number
                        v-model="hostFormData.configuration.cpu_cores"
                        :min="1"
                        :max="128"
                        size="large"
                        style="width: 100%"
                        controls-position="right"
                      />
                    </el-form-item>
                  </el-col>
                  <el-col :span="8">
                    <el-form-item label="内存(GB)" prop="memory_size">
                      <el-input-number
                        v-model="hostFormData.configuration.memory_size"
                        :min="1"
                        :max="1024"
                        size="large"
                        style="width: 100%"
                        controls-position="right"
                      />
                    </el-form-item>
                  </el-col>
                  <el-col :span="8">
                    <el-form-item label="磁盘(GB)" prop="disk_size">
                      <el-input-number
                        v-model="hostFormData.configuration.disk_size"
                        :min="1"
                        :max="10240"
                        size="large"
                        style="width: 100%"
                        controls-position="right"
                      />
                    </el-form-item>
                  </el-col>
                </el-row>
              </div>
            </div>

            <!-- 备注信息区域 -->
            <div class="form-section">
              <div class="section-header">
                <div class="section-icon note-icon">
                  <i class="el-icon-document"></i>
                </div>
                <div class="section-title">
                  <h3>备注信息</h3>
                  <p>主机的描述和标签信息</p>
                </div>
              </div>
              <div class="section-content">
                <el-row :gutter="24">
                  <el-col :span="12">
                    <el-form-item label="主机描述" prop="description">
                      <el-input
                        v-model="hostFormData.description"
                        type="textarea"
                        :rows="3"
                        placeholder="请输入主机描述信息（可选）"
                        size="large"
                      />
                    </el-form-item>
                  </el-col>
                  <el-col :span="12">
                    <el-form-item label="标签" prop="tags">
                      <el-input
                        v-model="hostFormData.tags"
                        placeholder="请输入标签，多个标签用逗号分隔（可选）"
                        size="large"
                      />
                    </el-form-item>
                  </el-col>
                </el-row>
              </div>
            </div>
          </el-form>
        </div>

        <!-- 底部操作按钮 -->
        <div class="form-footer">
          <div class="footer-actions">
            <el-button
              @click="closeHostDialog"
              size="large"
            >
              取消
            </el-button>
            <el-button
              type="primary"
              @click="handleSubmitHost"
              size="large"
              :loading="submitLoading"
            >
              {{ isEditMode ? '保存修改' : '创建主机' }}
            </el-button>
          </div>
        </div>
      </div>
    </el-dialog>

    <!-- 新建主机组对话框 -->
    <el-dialog
      v-model="addGroupDialogVisible"
      title="新建主机组"
      width="500px"
      :before-close="() => addGroupDialogVisible = false"
    >
      <el-form :model="groupFormData" label-width="100px" class="group-form">
        <el-form-item label="分组名称" required>
          <el-input v-model="groupFormData.name" placeholder="请输入分组名称" />
        </el-form-item>
        <el-form-item label="描述信息">
          <el-input
            v-model="groupFormData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入描述信息"
          />
        </el-form-item>
        <el-form-item label="父分组">
          <el-select v-model="groupFormData.parent_id" placeholder="选择父分组" clearable>
            <el-option
              v-for="group in flattenGroups(groups)"
              :key="group.id"
              :label="group.name"
              :value="group.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="groupFormData.sort" :min="0" :max="999" />
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="addGroupDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveGroup">确定</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 编辑主机组对话框 -->
    <el-dialog
      v-model="editGroupDialogVisible"
      title="编辑主机组"
      width="500px"
      :before-close="() => editGroupDialogVisible = false"
    >
      <el-form :model="groupFormData" label-width="100px" class="group-form">
        <el-form-item label="分组名称" required>
          <el-input v-model="groupFormData.name" placeholder="请输入分组名称" />
        </el-form-item>
        <el-form-item label="描述信息">
          <el-input
            v-model="groupFormData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入描述信息"
          />
        </el-form-item>
        <el-form-item label="父分组">
          <el-select v-model="groupFormData.parent_id" placeholder="选择父分组" clearable>
            <el-option
              v-for="group in flattenGroups(groups).filter(g => g.id !== currentEditGroup?.id)"
              :key="group.id"
              :label="group.name"
              :value="group.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="groupFormData.sort" :min="0" :max="999" />
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="editGroupDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveGroup">确定</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 现代化主机详情模态窗口 -->
    <el-dialog
      v-model="showDetailDialog"
      :title="null"
      width="90%"
      :close-on-click-modal="false"
      class="modern-host-detail-dialog"
      :show-close="false"
    >
      <div v-if="selectedHost" class="modern-detail-container">
        <!-- 自定义头部 -->
        <div class="modern-header">
          <div class="header-left">
            <div class="host-avatar">
              <i class="el-icon-monitor"></i>
            </div>
            <div class="host-info">
              <h2 class="host-name">{{ selectedHost.name }}</h2>
              <div class="host-meta">
                <el-tag
                  :type="getStatusType(selectedHost.status)"
                  class="status-tag"
                  effect="dark"
                >
                  <i class="status-dot"></i>
                  {{ getStatusText(selectedHost.status) }}
                </el-tag>
                <span class="provider-badge">
                  <i class="provider-icon"></i>
                  {{ getProviderText(selectedHost.provider_type) }}
                </span>
              </div>
            </div>
          </div>
          <div class="header-actions">
            <el-button
              type="primary"
              :icon="Edit"
              circle
              @click="handleEdit(selectedHost)"
              class="action-btn"
            />
            <el-button
              type="success"
              :icon="Monitor"
              circle
              @click="handleTerminal(selectedHost)"
              class="action-btn"
            />
            <el-button
              type="info"
              :icon="Close"
              circle
              @click="showDetailDialog = false"
              class="close-btn"
            />
          </div>
        </div>

        <!-- 详情内容 -->
        <div class="modern-content">
          <!-- 快速信息栏 -->
          <div class="quick-info-bar">
            <div class="quick-info-item">
              <div class="info-icon cpu-icon">
                <i class="el-icon-cpu"></i>
              </div>
              <div class="info-content">
                <div class="info-label">CPU</div>
                <div class="info-value">{{ formatConfiguration(selectedHost.configuration).cpu }}</div>
              </div>
            </div>
            <div class="quick-info-item">
              <div class="info-icon memory-icon">
                <i class="el-icon-memory"></i>
              </div>
              <div class="info-content">
                <div class="info-label">内存</div>
                <div class="info-value">{{ formatConfiguration(selectedHost.configuration).memory }}</div>
              </div>
            </div>
            <div class="quick-info-item">
              <div class="info-icon disk-icon">
                <i class="el-icon-hard-disk"></i>
              </div>
              <div class="info-content">
                <div class="info-label">磁盘</div>
                <div class="info-value">{{ formatConfiguration(selectedHost.configuration).disk }}</div>
              </div>
            </div>
            <div class="quick-info-item">
              <div class="info-icon location-icon">
                <i class="el-icon-location"></i>
              </div>
              <div class="info-content">
                <div class="info-label">地域</div>
                <div class="info-value">{{ selectedHost.region }}</div>
              </div>
            </div>
          </div>

          <!-- 详细信息卡片组 -->
          <div class="detail-cards-grid">
            <!-- 基本信息 -->
            <div class="modern-card">
              <div class="card-header">
                <div class="card-icon basic-icon">
                  <i class="el-icon-info"></i>
                </div>
                <h3 class="card-title">基本信息</h3>
              </div>
              <div class="card-content">
                <div class="info-row">
                  <span class="label">实例ID</span>
                  <span class="value">{{ selectedHost.instance_id || '-' }}</span>
                </div>
                <div class="info-row">
                  <span class="label">操作系统</span>
                  <span class="value">{{ selectedHost.os }}</span>
                </div>
                <div class="info-row">
                  <span class="label">可用区</span>
                  <span class="value">{{ selectedHost.availability_zone || '-' }}</span>
                </div>
                <div class="info-row">
                  <span class="label">实例类型</span>
                  <span class="value">{{ getInstanceType(selectedHost.configuration) }}</span>
                </div>
              </div>
            </div>

            <!-- 网络信息 -->
            <div class="modern-card">
              <div class="card-header">
                <div class="card-icon network-icon">
                  <i class="el-icon-connection"></i>
                </div>
                <h3 class="card-title">网络信息</h3>
              </div>
              <div class="card-content">
                <div class="info-row">
                  <span class="label">公网IP</span>
                  <div class="ip-container">
                    <template v-if="getDisplayIP(selectedHost.public_ip)">
                      <span class="ip-value">{{ getDisplayIP(selectedHost.public_ip) }}</span>
                      <el-button
                        size="small"
                        type="primary"
                        :icon="DocumentCopy"
                        circle
                        class="copy-btn"
                        @click="copyToClipboard(getDisplayIP(selectedHost.public_ip))"
                      />
                    </template>
                    <span v-else class="value">-</span>
                  </div>
                </div>
                <div class="info-row">
                  <span class="label">私网IP</span>
                  <div class="ip-container">
                    <template v-if="getDisplayIP(selectedHost.private_ip)">
                      <span class="ip-value">{{ getDisplayIP(selectedHost.private_ip) }}</span>
                      <el-button
                        size="small"
                        type="info"
                        :icon="DocumentCopy"
                        circle
                        class="copy-btn"
                        @click="copyToClipboard(getDisplayIP(selectedHost.private_ip))"
                      />
                    </template>
                    <span v-else class="value">-</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- 时间信息 -->
            <div class="modern-card">
              <div class="card-header">
                <div class="card-icon time-icon">
                  <i class="el-icon-time"></i>
                </div>
                <h3 class="card-title">时间信息</h3>
              </div>
              <div class="card-content">
                <div class="info-row">
                  <span class="label">创建时间</span>
                  <span class="value">{{ formatDateTime(selectedHost.created_at) }}</span>
                </div>
                <div class="info-row">
                  <span class="label">更新时间</span>
                  <span class="value">{{ formatDateTime(selectedHost.updated_at) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>

    <!-- 批量导入对话框 -->
    <BatchImportModal
      v-model:visible="batchImportVisible"
      @success="handleBatchImportSuccess"
    />

    <!-- 导出选择对话框 -->
    <ExportDialog
      v-model="exportDialogVisible"
      :current-filters="currentFilters"
      @export="handleExportConfirm"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { DocumentCopy, Edit, Monitor, Close, Upload, Download } from '@element-plus/icons-vue'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'
import { useHostStore } from '@/store/modules/host'
import { getHostGroupTree } from '@/api/system/host'
import * as hostApi from '@/api/system/host'
import BatchImportModal from './components/BatchImportModal.vue'
import ExportDialog from './components/ExportDialog.vue'
import dayjs from 'dayjs'

// 类型定义
interface HostGroup {
  id?: number | null
  name: string
  description?: string
  parent_id?: number | null
  host_count?: number
  children?: HostGroup[]
  level?: number
}

// Router
const router = useRouter()

// Store
const hostStore = useHostStore()
const { hostList, hostGroupTree, isLoading, pagination } = storeToRefs(hostStore)

// 响应式数据
const loading = computed(() => isLoading.value)
const groupLoading = ref(false)
const selectedGroupId = ref(null)
const searchQuery = ref('')
const groupSearchQuery = ref('')
const statusFilter = ref('')
const regionFilter = ref('')
const providerFilter = ref('')
const addHostDialogVisible = ref(false)
const batchImportVisible = ref(false)

// 筛选选项
const filterOptions = ref({
  statuses: [] as string[],
  regions: [] as string[],
  providers: [] as string[]
})

// 主机详情模态窗口状态
const showDetailDialog = ref(false)
const selectedHost = ref<any>(null)

// 主机组对话框状态
const addGroupDialogVisible = ref(false)
const editGroupDialogVisible = ref(false)
const currentEditGroup = ref<HostGroup | null>(null)

// 使用store中的数据
const allHosts = computed(() => hostList.value)
const groups = computed(() => hostGroupTree.value)

// 前端分页逻辑
const currentPage = ref(1)
const pageSize = ref(10) // 调整为10条/页，更适合当前显示环境

// 表单数据
const hostFormData = reactive({
  name: '',
  instance_id: '', // 实例ID，如果为空则自动生成
  host_type: 'cloud', // cloud 或 self-hosted
  host_group: '', // 主机组ID
  provider_type: 'aws',
  public_ip: [''],
  private_ip: [''],
  os: '',
  region: '',
  description: '',
  tags: '',
  configuration: {
    cpu_cores: 1,
    memory_size: 1,
    disk_size: 20,
    instance_type: 't2.micro'
  },
  ssh_config: {
    port: 22,
    username: 'root',
    auth_type: 'password', // password 或 key
    password: '',
    private_key: '',
    passphrase: ''
  }
})

// 表单相关数据
const isEditMode = ref(false)
const submitLoading = ref(false)
const hostFormRef = ref()
const currentEditHostId = ref(null)

// 主机组数据
const hostGroups = ref<HostGroup[]>([])

// 获取主机组列表
const loadHostGroups = async () => {
  try {
    const response = await getHostGroupTree()
    if (response.code === 200) {
      // 将树形结构扁平化为列表，方便下拉框使用
      const flattenGroups = (groups: HostGroup[], result: HostGroup[] = []): HostGroup[] => {
        groups.forEach(group => {
          result.push({
            id: group.id,
            name: group.name,
            description: group.description,
            parent_id: group.parent_id
          })
          if (group.children && group.children.length > 0) {
            flattenGroups(group.children, result)
          }
        })
        return result
      }
      hostGroups.value = flattenGroups(response.data || [])
    }
  } catch (error) {
    console.error('获取主机组列表失败:', error)
    ElMessage.error('获取主机组列表失败')
  }
}

// 表单验证规则
const hostRules = {
  name: [
    { required: true, message: '请输入主机名称', trigger: 'blur' }
  ],
  host_type: [
    { required: true, message: '请选择主机类型', trigger: 'change' }
  ],
  host_group: [
    { required: true, message: '请选择主机组', trigger: 'change' }
  ],
  provider_type: [
    {
      required: true,
      message: '请选择云厂商',
      trigger: 'change',
      validator: (rule, value, callback) => {
        if (hostFormData.host_type === 'cloud' && !value) {
          callback(new Error('云服务器必须选择云厂商'))
        } else {
          callback()
        }
      }
    }
  ],
  ssh_username: [
    { required: true, message: '请输入SSH用户名', trigger: 'blur' }
  ],
  'ssh_config.password': [
    {
      validator: (rule, value, callback) => {
        if (hostFormData.ssh_config?.auth_type === 'password' && !value) {
          callback(new Error('密码认证方式必须输入密码'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  ssh_private_key: [
    {
      validator: (rule, value, callback) => {
        if (hostFormData.ssh_config?.auth_type === 'key' && !value) {
          callback(new Error('密钥认证方式必须提供私钥'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 主机组表单数据
const groupFormData = reactive({
  name: '',
  description: '',
  parent_id: null as number | null,
  sort: 0
})

// 计算属性
const currentGroupName = computed(() => {
  const group = groups.value.find(g => g.id === selectedGroupId.value)
  return group ? group.name : '全部主机'
})

// 将树形结构扁平化为列表
const flattenGroups = (groups: HostGroup[], level = 0): HostGroup[] => {
  const result: HostGroup[] = []
  for (const group of groups) {
    result.push({ ...group, level })
    if (group.children && group.children.length > 0) {
      result.push(...flattenGroups(group.children, level + 1))
    }
  }
  return result
}

// 将平铺数组转换成树形结构
const buildGroupTree = (flatGroups: HostGroup[]): HostGroup[] => {
  const groupMap = new Map()
  const rootGroups: HostGroup[] = []

  // 创建所有分组的映射
  flatGroups.forEach(group => {
    groupMap.set(group.id, { ...group, children: [] })
  })

  // 构建树形结构
  flatGroups.forEach(group => {
    const groupNode = groupMap.get(group.id)
    if (group.parent_id === null || group.parent_id === undefined) {
      // 根节点
      rootGroups.push(groupNode)
    } else {
      // 子节点，添加到父节点的children中
      const parentNode = groupMap.get(group.parent_id)
      if (parentNode) {
        parentNode.children.push(groupNode)
      }
    }
  })

  return rootGroups
}

// 递归展开分组树，保持层级结构
const expandGroupTree = (groups: HostGroup[], level = 0): HostGroup[] => {
  const result: HostGroup[] = []

  groups.forEach(group => {
    // 添加当前分组，设置层级
    const groupWithLevel = { ...group, level }
    result.push(groupWithLevel)

    // 递归添加子分组
    if (group.children && group.children.length > 0) {
      result.push(...expandGroupTree(group.children, level + 1))
    }
  })

  return result
}

const filteredGroups = computed(() => {
  // 添加"全部主机"选项
  const allHostsGroup: HostGroup = {
    id: null,
    name: '全部主机',
    host_count: allHosts.value.length,
    level: 0
  }

  // 先将平铺数组转换成树形结构，再展开显示
  const treeGroups = buildGroupTree(groups.value)
  const expandedGroups = expandGroupTree(treeGroups)

  const allGroups = [allHostsGroup, ...expandedGroups]

  if (!groupSearchQuery.value) return allGroups
  return allGroups.filter(group =>
    group.name.toLowerCase().includes(groupSearchQuery.value.toLowerCase())
  )
})

// 先进行搜索和筛选，再分页
const filteredHosts = computed(() => {
  let result = allHosts.value

  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(host => {
      // 主机名称匹配
      if (host.name.toLowerCase().includes(query)) {
        return true
      }

      // 公网IP匹配 - 处理数组和字符串两种情况
      if (host.public_ip) {
        if (Array.isArray(host.public_ip)) {
          if (host.public_ip.some(ip => ip.includes(query))) {
            return true
          }
        } else if (typeof host.public_ip === 'string') {
          if (host.public_ip.includes(query)) {
            return true
          }
        }
      }

      // 私网IP匹配 - 处理数组和字符串两种情况
      if (host.private_ip) {
        if (Array.isArray(host.private_ip)) {
          if (host.private_ip.some(ip => ip.includes(query))) {
            return true
          }
        } else if (typeof host.private_ip === 'string') {
          if (host.private_ip.includes(query)) {
            return true
          }
        }
      }

      return false
    })
  }

  // 状态过滤
  if (statusFilter.value) {
    result = result.filter(host => host.status === statusFilter.value)
  }

  // 地域过滤
  if (regionFilter.value) {
    result = result.filter(host => host.region === regionFilter.value)
  }

  // 云厂商过滤
  if (providerFilter.value) {
    result = result.filter(host => host.provider_type === providerFilter.value)
  }

  return result
})

// 计算总数（基于筛选后的数据）
const total = computed(() => filteredHosts.value.length)

// 计算当前页显示的数据（基于筛选后的数据进行分页）
const hosts = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredHosts.value.slice(start, end)
})

// 监听搜索条件变化，重置分页
watch([searchQuery, statusFilter, regionFilter, providerFilter], () => {
  currentPage.value = 1
})

const regions = computed(() => {
  const regionSet = new Set()
  allHosts.value.forEach(host => {
    if (host.region) regionSet.add(host.region)
  })
  return Array.from(regionSet)
})

// 工具函数
const getDisplayIP = (ip) => {
  if (!ip) return ''
  if (Array.isArray(ip)) {
    return ip.length > 0 ? ip[0] : ''
  }
  if (typeof ip === 'string') {
    try {
      // 尝试解析JSON字符串格式的IP
      const parsed = JSON.parse(ip)
      if (Array.isArray(parsed)) {
        return parsed.length > 0 ? parsed[0] : ''
      }
      return parsed
    } catch {
      // 如果不是JSON格式，直接返回
      return ip
    }
  }
  return ''
}

const formatConfiguration = (config) => {
  if (!config) return { cpu: '-', memory: '-', disk: '-' }

  try {
    let configObj

    if (typeof config === 'string') {
      // 尝试解析JSON字符串，可能需要多次解析（处理双重编码）
      configObj = JSON.parse(config)

      // 如果解析后仍然是字符串，说明是双重编码，再解析一次
      if (typeof configObj === 'string') {
        configObj = JSON.parse(configObj)
      }
    } else {
      configObj = config
    }

    const result = {
      cpu: configObj.cpu_cores ? `${configObj.cpu_cores}核` : (configObj.instance_type || '-'),
      memory: configObj.memory_size ? `${configObj.memory_size}GB` : '-',
      disk: configObj.disk_size ? `${configObj.disk_size}GB` : (configObj.storage_size ? `${configObj.storage_size}GB` : '-')
    }

    return result
  } catch (error) {
    console.error('配置信息解析错误:', error, config)
    return { cpu: '-', memory: '-', disk: '-' }
  }
}

// 获取原始配置数据（用于编辑）
const getRawConfiguration = (config) => {
  if (!config) return { cpu_cores: 1, memory_size: 1, disk_size: 20, instance_type: 't2.micro' }

  try {
    const configObj = typeof config === 'string' ? JSON.parse(config) : config
    return {
      cpu_cores: configObj.cpu_cores || 1,
      memory_size: configObj.memory_size || 1,
      disk_size: configObj.disk_size || configObj.storage_size || 20,
      instance_type: configObj.instance_type || 't2.micro'
    }
  } catch (error) {
    console.error('配置信息解析错误:', error, config)
    return { cpu_cores: 1, memory_size: 1, disk_size: 20, instance_type: 't2.micro' }
  }
}

// 获取实例类型
const getInstanceType = (config) => {
  if (!config) return '-'

  try {
    const configObj = typeof config === 'string' ? JSON.parse(config) : config
    return configObj.instance_type || '-'
  } catch (error) {
    return '-'
  }
}

const getStatusType = (status) => {
  const statusMap = {
    'running': 'success',
    'stopped': 'danger',
    'restarting': 'warning',
    'pending': 'info'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status) => {
  const statusMap = {
    'running': '运行中',
    'stopped': '已停止',
    'restarting': '重启中',
    'pending': '待启动'
  }
  return statusMap[status] || status
}

const getProviderType = (provider) => {
  const providerMap = {
    'aws': 'warning',
    'aliyun': 'success',
    'tencent': 'primary',
    'huawei': 'info'
  }
  return providerMap[provider] || 'info'
}

const getProviderText = (provider) => {
  const providerMap = {
    'aws': 'AWS',
    'aliyun': '阿里云',
    'tencent': '腾讯云',
    'huawei': '华为云'
  }
  return providerMap[provider] || provider
}

// 格式化日期时间
const formatDateTime = (dateTime) => {
  if (!dateTime) return '-'

  try {
    let date

    // 处理各种时间格式
    if (typeof dateTime === 'number') {
      // Unix时间戳（秒）
      if (dateTime < 10000000000) {
        date = new Date(dateTime * 1000)
      } else {
        // Unix时间戳（毫秒）
        date = new Date(dateTime)
      }
    } else if (typeof dateTime === 'string') {
      // 字符串格式
      date = new Date(dateTime)
    } else {
      date = new Date(dateTime)
    }

    // 检查日期是否有效
    if (isNaN(date.getTime())) {
      return '-'
    }

    // 检查是否是1970年（通常表示无效时间戳）
    if (date.getFullYear() === 1970) {
      return '-'
    }

    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  } catch (error) {
    console.warn('时间格式化错误:', error, dateTime)
    return '-'
  }
}

// 事件处理
const refreshGroups = async () => {
  groupLoading.value = true
  try {
    await hostStore.fetchHostGroupTree()
    ElMessage.success('分组列表刷新成功')
  } catch (error) {
    console.error('获取主机组失败:', error)
    ElMessage.error('刷新失败')
  } finally {
    groupLoading.value = false
  }
}

const refreshHosts = async () => {
  try {
    const params = {
      keyword: searchQuery.value,
      status: statusFilter.value,
      region: regionFilter.value,
      group_id: selectedGroupId.value
    }

    await hostStore.fetchHosts(params)
    currentPage.value = 1  // 重置到第一页
    ElMessage.success('主机列表刷新成功')
  } catch (error) {
    console.error('获取主机列表失败:', error)
    ElMessage.error('刷新失败')
  }
}

// 获取筛选选项
const fetchFilterOptions = async () => {
  try {
    const response = await hostApi.getHostFilterOptions()
    filterOptions.value = response.data
  } catch (error) {
    console.error('获取筛选选项失败:', error)
  }
}

// 批量导入相关方法
const showBatchImportDialog = () => {
  batchImportVisible.value = true
}

const handleBatchImportSuccess = () => {
  batchImportVisible.value = false
  refreshHosts()
  ElMessage.success('批量导入成功')
}

// 批量导出主机
const handleBatchExport = async () => {
  try {
    const loading = ElLoading.service({
      lock: true,
      text: '正在准备导出数据...',
      background: 'rgba(0, 0, 0, 0.7)'
    })

    // 构建导出参数，包含当前的筛选条件
    const exportParams = {
      format: 'excel' as const,
      name: searchQuery.value || undefined,
      status: statusFilter.value || undefined,
      region: regionFilter.value || undefined,
      provider: providerFilter.value || undefined
    }

    const response = await hostApi.batchExportHosts(exportParams)

    // 创建下载链接
    const blob = new Blob([response as BlobPart], {
      type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
    })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `主机列表_${dayjs().format('YYYY-MM-DD_HH-mm-ss')}.xlsx`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)

    ElMessage.success('导出成功')
    loading.close()
  } catch (error) {
    console.error('导出失败:', error)
    ElMessage.error('导出失败')
  }
}

// 状态标签转换
const getStatusLabel = (status: string) => {
  const statusMap: Record<string, string> = {
    'running': '运行中',
    'stopped': '已停止',
    'restarting': '重启中',
    'pending': '启动中',
    'stopping': '停止中',
    'error': '错误',
    'unknown': '未知'
  }
  return statusMap[status] || status
}

// 云厂商标签转换
const getProviderLabel = (provider: string) => {
  const providerMap: Record<string, string> = {
    'aws': 'AWS',
    'aliyun': '阿里云',
    'tencent': '腾讯云',
    'huawei': '华为云',
    'manual': '自建'
  }
  return providerMap[provider] || provider
}

const handleGroupSelect = (group) => {
  selectedGroupId.value = group.id
  // 根据分组筛选主机
  refreshHosts()
}

const handleAdd = () => {
  // 重置表单
  Object.assign(hostFormData, {
    name: '',
    host_type: 'cloud',
    host_group: '',
    provider_type: 'aws',
    public_ip: [''],
    private_ip: [''],
    os: '',
    region: '',
    description: '',
    tags: '',
    configuration: {
      cpu_cores: 1,
      memory_size: 1,
      disk_size: 20,
      instance_type: 't2.micro'
    },
    ssh_config: {
      port: 22,
      username: 'root',
      auth_type: 'password',
      password: '',
      private_key: '',
      passphrase: ''
    }
  })
  // 重置状态
  isEditMode.value = false
  submitLoading.value = false
  addHostDialogVisible.value = true
}

// 处理主机类型变化
const handleHostTypeChange = (value) => {
  if (value === 'self-hosted') {
    // 自建服务器，清空云厂商相关字段
    hostFormData.provider_type = ''
    hostFormData.configuration.instance_type = ''
  } else {
    // 云服务器，设置默认值
    hostFormData.provider_type = 'aws'
    hostFormData.configuration.instance_type = 't2.micro'
  }
}

// 处理密钥文件上传
const handleKeyUpload = (file) => {
  const reader = new FileReader()
  reader.onload = (e) => {
    hostFormData.ssh_config.private_key = e.target.result
    ElMessage.success('密钥文件上传成功')
  }
  reader.onerror = () => {
    ElMessage.error('密钥文件读取失败')
  }
  reader.readAsText(file)
  return false // 阻止自动上传
}

const closeHostDialog = () => {
  addHostDialogVisible.value = false
  isEditMode.value = false
  submitLoading.value = false
}

const handleSubmitHost = async () => {
  submitLoading.value = true
  try {
    await saveHost()
  } finally {
    submitLoading.value = false
  }
}

// 生成实例ID
const generateInstanceId = () => {
  const timestamp = Date.now()
  const random = Math.random().toString(36).substring(2, 8)
  const instanceId = `${hostFormData.host_type}-${timestamp}-${random}`
  hostFormData.instance_id = instanceId
  ElMessage.success('实例ID已生成')
}

const saveHost = async () => {
  try {
    // 验证表单
    if (!hostFormData.name) {
      ElMessage.error('请输入主机名称')
      return
    }

    if (!hostFormData.host_type) {
      ElMessage.error('请选择主机类型')
      return
    }

    if (!hostFormData.host_group) {
      ElMessage.error('请选择主机组')
      return
    }

    if (hostFormData.host_type === 'cloud' && !hostFormData.provider_type) {
      ElMessage.error('云服务器必须选择云厂商')
      return
    }

    if (!hostFormData.ssh_config.username) {
      ElMessage.error('请输入SSH用户名')
      return
    }

    if (hostFormData.ssh_config.auth_type === 'password' && !hostFormData.ssh_config.password) {
      ElMessage.error('密码认证方式必须输入密码')
      return
    }

    if (hostFormData.ssh_config.auth_type === 'key' && !hostFormData.ssh_config.private_key) {
      ElMessage.error('密钥认证方式必须提供私钥')
      return
    }

    // 生成实例ID（如果没有提供的话）
    const generateInstanceIdInternal = () => {
      const timestamp = Date.now()
      const random = Math.random().toString(36).substring(2, 8)
      return `${hostFormData.host_type}-${timestamp}-${random}`
    }

    // 准备保存数据，确保符合后端接口要求
    const saveData = {
      // 必填字段
      instance_id: hostFormData.instance_id || generateInstanceIdInternal(), // 生成实例ID
      name: hostFormData.name,

      // 主机类型相关
      resource_type: hostFormData.host_type === 'cloud' ? 'cloud' : 'manual',
      provider_type: hostFormData.provider_type || 'manual',

      // 网络配置
      public_ip: Array.isArray(hostFormData.public_ip) ? hostFormData.public_ip : [hostFormData.public_ip].filter(Boolean),
      private_ip: Array.isArray(hostFormData.private_ip) ? hostFormData.private_ip : [hostFormData.private_ip].filter(Boolean),

      // 基本信息
      os: hostFormData.os || '',
      region: hostFormData.region || '',

      // SSH配置
      username: hostFormData.ssh_config.username,
      password: hostFormData.ssh_config.password || '',

      // 配置信息
      configuration: JSON.stringify(hostFormData.configuration),

      // 分组和其他信息
      group_id: hostFormData.host_group,
      remark: hostFormData.description || '',
      tags: JSON.stringify(hostFormData.tags || []),

      // 状态
      status: 'running' // 默认状态
    }

    // 调试：打印发送的数据
    console.log('🔍 准备发送的主机数据:', saveData)
    console.log('🔍 主机组ID:', hostFormData.host_group)
    console.log('🔍 主机组选项:', hostGroups.value)

    if (isEditMode.value && currentEditHostId.value) {
      // 编辑模式
      await hostStore.updateHost(currentEditHostId.value, saveData)
      ElMessage.success('主机更新成功')
    } else {
      // 添加模式
      await hostStore.createHost(saveData)
      ElMessage.success('主机添加成功')
    }

    // 关闭对话框并刷新数据
    closeHostDialog()
    refreshHosts()
  } catch (error) {
    console.error('保存失败:', error)
    ElMessage.error(isEditMode.value ? '更新失败' : '添加失败')
  }
}

// 复制到剪贴板
const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success(`已复制: ${text}`)
  } catch (error) {
    // 降级方案
    const textArea = document.createElement('textarea')
    textArea.value = text
    document.body.appendChild(textArea)
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
    ElMessage.success(`已复制: ${text}`)
  }
}

const handleView = (row) => {
  // 显示主机详情模态窗口
  selectedHost.value = { ...row }
  showDetailDialog.value = true
}

const handleTerminal = (row) => {
  // TODO: 实现SSH终端功能
  ElMessage.info(`连接SSH终端: ${row.name}`)
  console.log('连接SSH终端:', row)
}



const handleEdit = (row) => {
  // 填充表单数据
  const rawConfig = getRawConfiguration(row.configuration)
  Object.assign(hostFormData, {
    name: row.name,
    host_type: row.host_type || 'cloud',
    host_group: row.host_group || '',
    provider_type: row.provider_type,
    public_ip: Array.isArray(row.public_ip) ? row.public_ip : [row.public_ip || ''],
    private_ip: Array.isArray(row.private_ip) ? row.private_ip : [row.private_ip || ''],
    os: row.os,
    region: row.region,
    description: row.description || '',
    tags: row.tags || '',
    configuration: rawConfig,
    ssh_config: {
      port: row.ssh_config?.port || 22,
      username: row.ssh_config?.username || 'root',
      auth_type: row.ssh_config?.auth_type || 'password',
      password: row.ssh_config?.password || '',
      private_key: row.ssh_config?.private_key || '',
      passphrase: row.ssh_config?.passphrase || ''
    }
  })

  // 设置编辑模式
  isEditMode.value = true
  submitLoading.value = false

  // 存储当前编辑的主机ID
  currentEditHostId.value = row.id

  // 打开对话框
  addHostDialogVisible.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(`确定要删除主机 "${row.name}" 吗？`, '确认删除', {
      type: 'warning'
    })

    // 调用删除API
    await hostStore.deleteHost(row.id)
    ElMessage.success('删除成功')
    refreshHosts()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败: ' + (error.message || error))
    }
  }
}

// 主机组管理功能
const handleAddGroup = () => {
  // 重置表单
  Object.assign(groupFormData, {
    name: '',
    description: '',
    parent_id: null,
    sort: 0
  })
  addGroupDialogVisible.value = true
}

const handleEditGroup = (group: HostGroup) => {
  currentEditGroup.value = group
  Object.assign(groupFormData, {
    name: group.name,
    description: group.description || '',
    parent_id: group.parent_id,
    sort: 0
  })
  editGroupDialogVisible.value = true
}

const handleDeleteGroup = async (group: HostGroup) => {
  if (!group.id) return

  try {
    await ElMessageBox.confirm(`确定要删除分组 "${group.name}" 吗？`, '确认删除', {
      type: 'warning'
    })

    await hostStore.deleteHostGroup(group.id)
    ElMessage.success('删除成功')

    // 自动刷新主机组列表
    await refreshGroups()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('删除主机组失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

const saveGroup = async () => {
  try {
    if (!groupFormData.name) {
      ElMessage.error('请输入分组名称')
      return
    }

    if (currentEditGroup.value && currentEditGroup.value.id) {
      // 编辑模式
      await hostStore.updateHostGroup(currentEditGroup.value.id, groupFormData)
      ElMessage.success('更新成功')
      editGroupDialogVisible.value = false
    } else {
      // 新建模式
      await hostStore.addHostGroup(groupFormData)
      ElMessage.success('创建成功')
      addGroupDialogVisible.value = false
    }

    // 自动刷新主机组列表
    await refreshGroups()
  } catch (error) {
    console.error('保存主机组失败:', error)
    ElMessage.error('保存失败')
  }
}

const handleSelectionChange = (selection) => {
  console.log('选中的主机:', selection)
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
}

// 初始化
onMounted(async () => {
  selectedGroupId.value = null

  // 初始化数据
  await loadHostGroups()
  await refreshGroups()
  await refreshHosts()
  await fetchFilterOptions()
})
</script>

<style scoped>
/* 主容器 */
.host-management {
  height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.page-container {
  display: flex;
  height: 100%;
  width: 100%;
  background: white;
  box-shadow: 0 0 30px rgba(0, 0, 0, 0.1);
}

/* 左侧边栏 */
.sidebar {
  width: 240px;
  background: linear-gradient(180deg, #ffffff 0%, #f8f9fa 100%);
  border-right: 1px solid #e9ecef;
  display: flex;
  flex-direction: column;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.05);
}

.sidebar-header {
  padding: 16px;
  border-bottom: 1px solid #e9ecef;
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  color: white;
  position: relative;
  overflow: hidden;
}

.sidebar-header::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="grid" width="10" height="10" patternUnits="userSpaceOnUse"><path d="M 10 0 L 0 0 0 10" fill="none" stroke="rgba(255,255,255,0.1)" stroke-width="0.5"/></pattern></defs><rect width="100" height="100" fill="url(%23grid)"/></svg>');
  opacity: 0.3;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: relative;
  z-index: 1;
}

.header-actions {
  display: flex;
  gap: 6px;
}

.title {
  margin: 0;
  font-size: 16px;
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 8px;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.icon {
  font-size: 18px;
  filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.1));
}

.search-section {
  padding: 12px 16px;
  border-bottom: 1px solid #f0f0f0;
  background: #fafbfc;
}

.search-input {
  width: 100%;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: 20px !important;
  border: 1px solid #e2e8f0 !important;
  background: white !important;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1) !important;
  transition: all 0.3s ease !important;
}

.search-input :deep(.el-input__wrapper):hover {
  border-color: #4f46e5 !important;
  box-shadow: 0 2px 8px rgba(79, 70, 229, 0.15) !important;
}

.search-icon {
  color: #6b7280;
  font-size: 14px;
}

/* 分组列表 */
.group-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px 12px;
}

.group-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  margin: 4px 0;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  background: white;
  border: 1px solid #e2e8f0;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.group-item::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 4px;
  background: transparent;
  border-radius: 12px 0 0 12px;
  transition: all 0.3s ease;
}

.group-item:hover {
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  transform: translateX(2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-color: #cbd5e1;
}

.group-item:hover::before {
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
}

.group-item.active {
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  color: white;
  box-shadow: 0 8px 25px rgba(79, 70, 229, 0.3);
  border-color: transparent;
  transform: translateX(4px);
}

.group-item.active::before {
  background: rgba(255, 255, 255, 0.3);
}

.group-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.group-info {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  min-width: 0; /* 允许flex子元素收缩 */
}

.group-icon {
  font-size: 14px;
  filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.1));
  flex-shrink: 0; /* 图标不收缩 */
}

.group-name {
  font-weight: 600;
  font-size: 12px;
  letter-spacing: 0.025em;
  white-space: nowrap; /* 不换行 */
  overflow: hidden; /* 隐藏溢出 */
  text-overflow: ellipsis; /* 显示省略号 */
  flex: 1;
}

.group-count {
  background: rgba(0, 0, 0, 0.08);
  color: inherit;
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 11px;
  font-weight: 700;
  min-width: 24px;
  text-align: center;
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.group-item.active .group-count {
  background: rgba(255, 255, 255, 0.25);
  border-color: rgba(255, 255, 255, 0.2);
  color: white;
}

.group-actions {
  display: flex;
  gap: 2px;
  opacity: 0;
  transition: all 0.3s ease;
  transform: translateX(8px);
}

.group-item:hover .group-actions {
  opacity: 1;
  transform: translateX(0);
}

.group-actions .el-button {
  padding: 4px 6px !important;
  border-radius: 6px !important;
  font-size: 12px !important;
  background: rgba(255, 255, 255, 0.9) !important;
  border: 1px solid rgba(0, 0, 0, 0.1) !important;
  color: #6b7280 !important;
}

.group-actions .el-button:hover {
  background: white !important;
  color: #374151 !important;
  transform: scale(1.05) !important;
}

/* 主内容区 */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #f8fafc;
  min-height: 100vh;
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 20px 24px;
  background: linear-gradient(135deg, #ffffff 0%, #f8fafc 100%);
  border-bottom: 1px solid #e2e8f0;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  position: relative;
}

.content-header::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg, transparent 0%, #4f46e5 50%, transparent 100%);
  opacity: 0.3;
}

.header-left {
  flex: 1;
}

.page-title {
  margin: 0 0 6px 0;
  font-size: 24px;
  font-weight: 800;
  color: #1e293b;
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: -0.025em;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.breadcrumb {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #64748b;
  font-size: 13px;
  font-weight: 500;
}

.separator {
  color: #cbd5e1;
  font-weight: 300;
}

.current {
  color: #475569;
  font-weight: 600;
  background: linear-gradient(135deg, #e2e8f0 0%, #f1f5f9 100%);
  padding: 2px 8px;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.header-actions .el-button {
  border-radius: 8px !important;
  font-weight: 600 !important;
  padding: 8px 16px !important;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1) !important;
  transition: all 0.3s ease !important;
}

.header-actions .el-button:hover {
  transform: translateY(-1px) !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15) !important;
}

/* 工具栏 */
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background: white;
  border-bottom: 1px solid #e2e8f0;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.toolbar-left {
  flex: 1;
  max-width: 400px;
}

.toolbar-left .search-input :deep(.el-input__wrapper) {
  border-radius: 24px !important;
  border: 1px solid #e2e8f0 !important;
  background: #f8fafc !important;
  box-shadow: inset 0 1px 3px rgba(0, 0, 0, 0.05) !important;
  transition: all 0.3s ease !important;
}

.toolbar-left .search-input :deep(.el-input__wrapper):hover {
  background: white !important;
  border-color: #4f46e5 !important;
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1) !important;
}

.toolbar-left .search-input :deep(.el-input__wrapper.is-focus) {
  background: white !important;
  border-color: #4f46e5 !important;
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1) !important;
}

.toolbar-right {
  display: flex;
  gap: 8px;
}

.toolbar-right .el-select {
  width: 140px;
}

.toolbar-right .el-select :deep(.el-select__wrapper) {
  border-radius: 8px !important;
  border: 1px solid #e2e8f0 !important;
  background: white !important;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05) !important;
  transition: all 0.3s ease !important;
}

.toolbar-right .el-select :deep(.el-select__wrapper):hover {
  border-color: #4f46e5 !important;
  box-shadow: 0 2px 8px rgba(79, 70, 229, 0.15) !important;
}

/* 表格容器 */
.table-container {
  flex: 1;
  margin: 16px 24px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 25px rgba(0, 0, 0, 0.08);
  overflow: auto; /* 改为auto，允许滚动 */
  border: 1px solid #e2e8f0;
  position: relative;
  min-height: 0; /* 确保flex子元素可以收缩 */
}

.table-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, #4f46e5 0%, #7c3aed 50%, #06b6d4 100%);
  opacity: 0.8;
}

.host-table {
  width: 100%;
}

/* 主机名称 */
.host-name {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
}

.host-icon {
  font-size: 16px;
}

/* IP信息 */
.ip-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.ip-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
}

.ip-label {
  background: #f8f9fa;
  color: #495057;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
  min-width: 32px;
  text-align: center;
}

.ip-value {
  font-family: 'Monaco', 'Menlo', monospace;
  color: #2c3e50;
  font-weight: 500;
}

.ip-empty {
  color: #adb5bd;
  font-style: italic;
}

/* 配置信息 */
.config-specs {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 8px;
  background: #fafbfc;
  border-radius: 6px;
  border: 1px solid #e1e4e8;
}

.spec-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  padding: 3px 0;
}

.spec-label {
  color: #586069;
  font-weight: 500;
  min-width: 40px;
  font-size: 12px;
}

.spec-value {
  color: #24292e;
  font-weight: 600;
  background: #ffffff;
  padding: 2px 8px;
  border-radius: 4px;
  border: 1px solid #e1e4e8;
  font-size: 12px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  border: 1px solid #e9ecef;
}

/* IP地址样式优化 */
.ip-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.ip-item {
  font-size: 12px;
}

.ip-row {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-bottom: 2px;
}

.ip-label {
  color: #666;
  font-weight: 500;
  min-width: 35px;
}

.ip-value {
  color: #333;
  font-weight: 600;
  background: #f0f9ff;
  padding: 2px 6px;
  border-radius: 4px;
  border: 1px solid #e1f5fe;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

.ip-empty {
  color: #999;
  font-style: italic;
}

.copy-btn {
  opacity: 0;
  transition: opacity 0.2s;
  padding: 2px 4px !important;
  margin-left: 4px;
  color: #409eff;
}

.ip-item:hover .copy-btn {
  opacity: 1;
}

.copy-btn:hover {
  color: #66b1ff;
}

/* 系统和地域信息 */
.os-info, .region-info {
  color: #495057;
  font-size: 13px;
}

/* 操作按钮 */
.action-buttons {
  display: flex;
  gap: 4px;
}

.action-buttons .el-button {
  padding: 4px 8px;
  font-size: 12px;
}

/* 分页 */
.pagination-container {
  display: flex;
  justify-content: center;
  padding: 20px 32px;
  background: white;
  border-top: 1px solid #f0f0f0;
}

/* 对话框 */
.host-form, .group-form {
  padding: 20px 0;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* 表格样式增强 */
.host-table :deep(.el-table__header) {
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
}

.host-table :deep(.el-table__header th) {
  background: transparent !important;
  color: #2c3e50 !important;
  font-weight: 600 !important;
  border-bottom: 2px solid #dee2e6 !important;
  font-size: 13px !important;
}

.host-table :deep(.el-table__row) {
  transition: all 0.3s ease;
}

.host-table :deep(.el-table__row):hover {
  background: linear-gradient(135deg, #f8f9fa 0%, #e3f2fd 100%) !important;
}

/* 标签美化 */
.el-tag {
  border: none !important;
  font-weight: 500 !important;
  border-radius: 6px !important;
}

/* 按钮美化 */
.el-button {
  border-radius: 6px !important;
  font-weight: 500 !important;
  transition: all 0.3s ease !important;
}

.el-button:hover {
  transform: translateY(-1px) !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15) !important;
}

/* 输入框美化 */
.el-input :deep(.el-input__wrapper) {
  border-radius: 8px !important;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05) !important;
  transition: all 0.3s ease !important;
}

.el-input :deep(.el-input__wrapper):hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1) !important;
}

.el-input :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2) !important;
}

/* 选择器美化 */
.el-select :deep(.el-select__wrapper) {
  border-radius: 8px !important;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05) !important;
}

/* 滚动条美化 */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(135deg, #5a6fd8 0%, #6a4190 100%);
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .sidebar {
    width: 220px;
  }
}

@media (max-width: 768px) {
  .sidebar {
    display: none;
  }

  .content-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }

  .toolbar {
    flex-direction: column;
    gap: 16px;
  }

  .toolbar-right {
    flex-wrap: wrap;
  }

  .table-container {
    margin: 16px;
  }

  .action-buttons {
    flex-direction: column;
    gap: 2px;
  }
}

/* 现代化主机详情模态窗口样式 */
.modern-host-detail-dialog {
  .el-dialog {
    border-radius: 20px;
    overflow: hidden;
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  }

  .el-dialog__body {
    padding: 0;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }
}

.modern-detail-container {
  min-height: 600px;
}

/* 现代化头部 */
.modern-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 30px 40px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  position: relative;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="grain" width="100" height="100" patternUnits="userSpaceOnUse"><circle cx="25" cy="25" r="1" fill="white" opacity="0.1"/><circle cx="75" cy="75" r="1" fill="white" opacity="0.1"/><circle cx="50" cy="10" r="0.5" fill="white" opacity="0.1"/><circle cx="10" cy="90" r="0.5" fill="white" opacity="0.1"/></pattern></defs><rect width="100" height="100" fill="url(%23grain)"/></svg>');
    pointer-events: none;
  }
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
  z-index: 1;
}

.host-avatar {
  width: 80px;
  height: 80px;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  color: white;
  border: 2px solid rgba(255, 255, 255, 0.3);
}

.host-info {
  .host-name {
    font-size: 28px;
    font-weight: 700;
    margin: 0 0 10px 0;
    color: white;
    text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
  }

  .host-meta {
    display: flex;
    align-items: center;
    gap: 15px;
  }
}

.status-tag {
  border-radius: 20px;
  padding: 8px 16px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 8px;

  .status-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: currentColor;
    animation: pulse 2s infinite;
  }
}

.provider-badge {
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  padding: 8px 16px;
  border-radius: 20px;
  font-weight: 600;
  border: 1px solid rgba(255, 255, 255, 0.3);
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-actions {
  display: flex;
  gap: 12px;
  z-index: 1;

  .action-btn, .close-btn {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    backdrop-filter: blur(10px);
    border: 2px solid rgba(255, 255, 255, 0.3);
    transition: all 0.3s ease;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 8px 25px rgba(0, 0, 0, 0.3);
    }
  }
}

/* 现代化内容区域 */
.modern-content {
  background: #f8fafc;
  padding: 40px;
  min-height: 500px;
}

/* 快速信息栏 */
.quick-info-bar {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 40px;
}

.quick-info-item {
  background: white;
  border-radius: 16px;
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  border: 1px solid #e2e8f0;

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
  }
}

.info-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  color: white;

  &.cpu-icon {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }

  &.memory-icon {
    background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  }

  &.disk-icon {
    background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  }

  &.location-icon {
    background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
  }
}

.info-content {
  .info-label {
    font-size: 14px;
    color: #64748b;
    margin-bottom: 4px;
  }

  .info-value {
    font-size: 18px;
    font-weight: 700;
    color: #1e293b;
  }
}

/* 详细信息卡片组 */
.detail-cards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 24px;
}

.modern-card {
  background: white;
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e2e8f0;
  transition: all 0.3s ease;

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  }
}

.card-header {
  padding: 24px 24px 16px 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  border-bottom: 1px solid #f1f5f9;
}

.card-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  color: white;

  &.basic-icon {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }

  &.network-icon {
    background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  }

  &.time-icon {
    background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  }
}

.card-title {
  font-size: 18px;
  font-weight: 700;
  color: #1e293b;
  margin: 0;
}

.card-content {
  padding: 16px 24px 24px 24px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f8fafc;

  &:last-child {
    border-bottom: none;
  }

  .label {
    font-weight: 500;
    color: #64748b;
    font-size: 14px;
  }

  .value {
    font-weight: 600;
    color: #1e293b;
    font-size: 14px;
  }
}

.ip-container {
  display: flex;
  align-items: center;
  gap: 8px;

  .ip-value {
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
    background: #f1f5f9;
    padding: 4px 8px;
    border-radius: 6px;
    font-size: 13px;
    font-weight: 600;
    color: #1e293b;
  }

  .copy-btn {
    width: 24px;
    height: 24px;
    padding: 0;
  }
}

/* 动画效果 */
@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

/* 现代化表单对话框样式 */
.modern-host-form-dialog {
  .el-dialog {
    border-radius: 20px;
    overflow: hidden;
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  }

  .el-dialog__body {
    padding: 0;
    background: #f8fafc;
  }
}

.modern-form-container {
  min-height: 700px;
  max-height: 80vh;
  overflow-y: auto;
}

/* 现代化表单头部 */
.modern-form-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 30px 40px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  position: relative;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><defs><pattern id="grain" width="100" height="100" patternUnits="userSpaceOnUse"><circle cx="25" cy="25" r="1" fill="white" opacity="0.1"/><circle cx="75" cy="75" r="1" fill="white" opacity="0.1"/><circle cx="50" cy="10" r="0.5" fill="white" opacity="0.1"/><circle cx="10" cy="90" r="0.5" fill="white" opacity="0.1"/></pattern></defs><rect width="100" height="100" fill="url(%23grain)"/></svg>');
    pointer-events: none;
  }
}

.form-avatar {
  width: 80px;
  height: 80px;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  color: white;
  border: 2px solid rgba(255, 255, 255, 0.3);
}

.form-info {
  .form-title {
    font-size: 28px;
    font-weight: 700;
    margin: 0 0 10px 0;
    color: white;
    text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
  }

  .form-subtitle {
    color: rgba(255, 255, 255, 0.8);
    font-size: 16px;
  }
}

/* 步骤指示器 */
.step-indicator {
  display: flex;
  justify-content: center;
  padding: 30px 40px;
  background: white;
  border-bottom: 1px solid #e2e8f0;
}

.step-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin: 0 30px;
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;

  &:not(:last-child)::after {
    content: '';
    position: absolute;
    top: 20px;
    left: 100%;
    width: 60px;
    height: 2px;
    background: #e2e8f0;
    transition: all 0.3s ease;
  }

  &.completed::after {
    background: #67c23a;
  }

  &.active::after {
    background: linear-gradient(90deg, #667eea 0%, #e2e8f0 100%);
  }
}

.step-number {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  margin-bottom: 8px;
  transition: all 0.3s ease;

  .step-item.completed & {
    background: #67c23a;
    color: white;
  }

  .step-item.active & {
    background: #667eea;
    color: white;
  }

  .step-item:not(.active):not(.completed) & {
    background: #e2e8f0;
    color: #64748b;
  }
}

.step-label {
  font-size: 14px;
  font-weight: 500;
  color: #64748b;
  transition: all 0.3s ease;

  .step-item.active & {
    color: #667eea;
  }

  .step-item.completed & {
    color: #67c23a;
  }
}

/* 表单内容 */
.form-content {
  padding: 30px 40px;
  min-height: 500px;
}

/* 表单区域样式 */
.form-section {
  background: white;
  border-radius: 16px;
  margin-bottom: 24px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  border: 1px solid #e2e8f0;
  overflow: hidden;
}

.section-header {
  display: flex;
  align-items: center;
  padding: 24px 30px;
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  border-bottom: 1px solid #e2e8f0;
}

.section-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  color: white;
  margin-right: 16px;
  flex-shrink: 0;

  &.basic-icon {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }

  &.network-icon {
    background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  }

  &.ssh-icon {
    background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  }

  &.hardware-icon {
    background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
  }

  &.note-icon {
    background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
  }
}

.section-title {
  h3 {
    font-size: 20px;
    font-weight: 700;
    color: #1e293b;
    margin: 0 0 4px 0;
  }

  p {
    color: #64748b;
    font-size: 14px;
    margin: 0;
  }
}

.section-content {
  padding: 30px;

  .el-form-item {
    margin-bottom: 32px;

    &:last-child {
      margin-bottom: 0;
    }
  }

  .el-form-item__label {
    font-weight: 600;
    color: #374151;
    font-size: 14px;
    line-height: 1.5;
    margin-bottom: 8px;
  }

  .el-form-item__error {
    position: absolute;
    top: 100%;
    left: 0;
    font-size: 12px;
    color: #f56565;
    line-height: 1;
    padding-top: 4px;
    z-index: 1;
  }
}

.step-content {
  animation: fadeInUp 0.5s ease;
}

.step-header {
  text-align: center;
  margin-bottom: 40px;

  h3 {
    font-size: 24px;
    font-weight: 700;
    color: #1e293b;
    margin: 0 0 8px 0;
  }

  p {
    color: #64748b;
    font-size: 16px;
    margin: 0;
  }
}

.form-grid {
  display: flex;
  justify-content: center;
}

.form-card {
  background: white;
  border-radius: 20px;
  padding: 40px;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e2e8f0;
  max-width: 500px;
  width: 100%;
  display: flex;
  align-items: flex-start;
  gap: 24px;
}

.card-icon {
  width: 60px;
  height: 60px;
  border-radius: 15px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: white;
  flex-shrink: 0;

  &.basic-icon {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }

  &.network-icon {
    background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  }

  &.system-icon {
    background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  }
}

.card-content {
  flex: 1;

  .el-form-item {
    margin-bottom: 24px;

    &:last-child {
      margin-bottom: 0;
    }
  }

  .el-form-item__label {
    font-weight: 600;
    color: #1e293b;
  }
}

/* 硬件配置网格 */
.hardware-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 24px;
  max-width: 800px;
  margin: 0 auto;
}

.hardware-card {
  background: white;
  border-radius: 20px;
  padding: 30px 20px;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e2e8f0;
  text-align: center;
  transition: all 0.3s ease;

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  }
}

/* 硬件配置卡片图标样式 */
.hardware-card .hardware-icon {
  width: 60px;
  height: 60px;
  border-radius: 15px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: white;
  margin: 0 auto 20px auto;

  &.cpu-icon {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }

  &.memory-icon {
    background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  }

  &.disk-icon {
    background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  }
}

.hardware-content {
  h4 {
    font-size: 18px;
    font-weight: 700;
    color: #1e293b;
    margin: 0 0 16px 0;
  }

  .el-form-item {
    margin-bottom: 0;
  }

  .el-input-number {
    width: 100%;
  }
}

/* 云厂商选项样式 */
.provider-option {
  display: flex;
  align-items: center;
  gap: 8px;

  .provider-icon {
    font-size: 16px;
  }
}

/* 密钥上传样式 */
.key-upload {
  .el-upload {
    width: 100%;
  }

  .el-button {
    width: 100%;
    border-style: dashed;
    border-color: #d1d5db;
    background: #f9fafb;

    &:hover {
      border-color: #667eea;
      background: #f0f4ff;
      color: #667eea;
    }
  }
}

/* 表单底部 */
.form-footer {
  padding: 30px 40px;
  background: white;
  border-top: 1px solid #e2e8f0;
  position: sticky;
  bottom: 0;
  z-index: 10;
}

.footer-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* 动画效果 */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .modern-host-detail-dialog {
    .el-dialog {
      width: 95% !important;
      margin: 2vh auto !important;
    }
  }

  .modern-header {
    padding: 20px;
    flex-direction: column;
    gap: 20px;
    text-align: center;
  }

  .modern-content {
    padding: 20px;
  }

  .quick-info-bar {
    grid-template-columns: 1fr;
  }

  .detail-cards-grid {
    grid-template-columns: 1fr;
  }

  .host-info .host-name {
    font-size: 24px;
  }

  .modern-host-form-dialog {
    .el-dialog {
      width: 95% !important;
      margin: 2vh auto !important;
    }
  }

  .modern-form-header {
    padding: 20px;
    flex-direction: column;
    gap: 20px;
    text-align: center;
  }

  .form-content {
    padding: 20px;
  }

  .step-indicator {
    padding: 20px;
    flex-wrap: wrap;
    gap: 20px;
  }

  .step-item {
    margin: 0;

    &:not(:last-child)::after {
      display: none;
    }
  }

  .hardware-grid {
    grid-template-columns: 1fr;
  }

  .form-footer {
    padding: 20px;
    flex-direction: column;
    gap: 16px;
  }

  .footer-left, .footer-right {
    width: 100%;
    justify-content: center;
  }
}
</style>
