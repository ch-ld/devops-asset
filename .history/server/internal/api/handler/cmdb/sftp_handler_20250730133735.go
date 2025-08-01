package cmdb

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"api-server/internal/config"
	cmdbModel "api-server/internal/model/cmdb"
	cmdbRepo "api-server/internal/repository/cmdb"
	"api-server/internal/response/response"
	"api-server/pkg/crypto/encryption"
	"api-server/pkg/ssh"

	"github.com/gin-gonic/gin"
)

// SftpHandler SFTP处理器
type SftpHandler struct {
	hostRepo *cmdbRepo.HostRepository
}

// NewSftpHandler 创建SFTP处理器
func NewSftpHandler(hostRepo *cmdbRepo.HostRepository) *SftpHandler {
	return &SftpHandler{
		hostRepo: hostRepo,
	}
}

// List 列出文件
func (h *SftpHandler) List(c *gin.Context) {
	hostIDStr := c.Query("host_id")
	path := c.DefaultQuery("path", "/")

	hostID, err := strconv.Atoi(hostIDStr)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "Invalid host_id")
		return
	}

	// 查询主机信息
	host, err := h.hostRepo.FindByID(uint(hostID))
	if err != nil {
		response.ReturnError(c, response.NOT_FOUND, "Host not found")
		return
	}

	// 创建SFTP客户端
	sftpClient, err := h.createSFTPClient(host)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, fmt.Sprintf("Failed to create SFTP client: %v", err))
		return
	}
	defer sftpClient.Close()

	// 列出文件
	files, err := sftpClient.ListDir(path)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, fmt.Sprintf("Failed to list files: %v", err))
		return
	}

	response.ReturnData(c, files)
}

// Upload 上传文件
func (h *SftpHandler) Upload(c *gin.Context) {
	hostIDStr := c.PostForm("host_id")
	remotePath := c.PostForm("path")

	hostID, err := strconv.Atoi(hostIDStr)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "Invalid host_id")
		return
	}

	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "No file uploaded")
		return
	}
	defer file.Close()

	// 查询主机信息
	host, err := h.hostRepo.FindByID(uint(hostID))
	if err != nil {
		response.ReturnError(c, response.NOT_FOUND, "Host not found")
		return
	}

	// 创建临时文件
	tempFile, err := os.CreateTemp("", "upload_*")
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "Failed to create temp file")
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// 复制文件内容到临时文件
	_, err = io.Copy(tempFile, file)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "Failed to save file")
		return
	}

	// 创建SFTP客户端
	sftpClient, err := h.createSFTPClient(host)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, fmt.Sprintf("Failed to create SFTP client: %v", err))
		return
	}
	defer sftpClient.Close()

	// 构建目标路径
	targetPath := filepath.Join(remotePath, header.Filename)

	// 上传文件
	err, _ = sftpClient.UploadFile(tempFile.Name(), targetPath)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, fmt.Sprintf("Failed to upload file: %v", err))
		return
	}

	response.ReturnSuccess(c)
}

// Download 下载文件
func (h *SftpHandler) Download(c *gin.Context) {
	hostIDStr := c.Query("host_id")
	filePath := c.Query("file_path")

	hostID, err := strconv.Atoi(hostIDStr)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "Invalid host_id")
		return
	}

	// 查询主机信息
	host, err := h.hostRepo.FindByID(uint(hostID))
	if err != nil {
		response.ReturnError(c, response.NOT_FOUND, "Host not found")
		return
	}

	// 创建SFTP客户端
	sftpClient, err := h.createSFTPClient(host)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, fmt.Sprintf("Failed to create SFTP client: %v", err))
		return
	}
	defer sftpClient.Close()

	// 创建临时文件
	tempFile, err := os.CreateTemp("", "download_*")
	if err != nil {
		response.ReturnError(c, response.INTERNAL, "Failed to create temp file")
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// 下载文件
	err = sftpClient.DownloadFile(filePath, tempFile.Name())
	if err != nil {
		response.ReturnError(c, response.INTERNAL, fmt.Sprintf("Failed to download file: %v", err))
		return
	}

	// 返回文件
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
	c.Header("Content-Type", "application/octet-stream")
	c.File(tempFile.Name())
}

// Delete 删除文件或目录
func (h *SftpHandler) Delete(c *gin.Context) {
	hostIDStr := c.Query("host_id")
	filePath := c.Query("file_path")

	hostID, err := strconv.Atoi(hostIDStr)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "Invalid host_id")
		return
	}

	// 查询主机信息
	host, err := h.hostRepo.FindByID(uint(hostID))
	if err != nil {
		response.ReturnError(c, response.NOT_FOUND, "Host not found")
		return
	}

	// 创建SFTP客户端
	sftpClient, err := h.createSFTPClient(host)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, fmt.Sprintf("Failed to create SFTP client: %v", err))
		return
	}
	defer sftpClient.Close()

	// 获取文件信息以判断是文件还是目录
	fileInfo, err := sftpClient.Stat(filePath)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, fmt.Sprintf("Failed to get file info: %v", err))
		return
	}

	// 根据文件类型选择删除方法
	if fileInfo.IsDir() {
		// 删除目录
		err = sftpClient.RemoveDirectory(filePath)
	} else {
		// 删除文件
		err = sftpClient.RemoveFile(filePath)
	}

	if err != nil {
		response.ReturnError(c, response.INTERNAL, fmt.Sprintf("Failed to delete: %v", err))
		return
	}

	response.ReturnSuccess(c)
}

// MakeDir 创建目录
func (h *SftpHandler) MakeDir(c *gin.Context) {
	hostIDStr := c.PostForm("host_id")
	dirPath := c.PostForm("path")

	hostID, err := strconv.Atoi(hostIDStr)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "Invalid host_id")
		return
	}

	// 查询主机信息
	host, err := h.hostRepo.FindByID(uint(hostID))
	if err != nil {
		response.ReturnError(c, response.NOT_FOUND, "Host not found")
		return
	}

	// 创建SFTP客户端
	sftpClient, err := h.createSFTPClient(host)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, fmt.Sprintf("Failed to create SFTP client: %v", err))
		return
	}
	defer sftpClient.Close()

	// 创建目录
	err = sftpClient.MakeDir(dirPath)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, fmt.Sprintf("Failed to create directory: %v", err))
		return
	}

	response.ReturnSuccess(c)
}

// Rename 重命名文件
func (h *SftpHandler) Rename(c *gin.Context) {
	hostIDStr := c.PostForm("host_id")
	oldPath := c.PostForm("old_path")
	newPath := c.PostForm("new_path")

	hostID, err := strconv.Atoi(hostIDStr)
	if err != nil {
		response.ReturnError(c, response.INVALID_ARGUMENT, "Invalid host_id")
		return
	}

	// 查询主机信息
	host, err := h.hostRepo.FindByID(uint(hostID))
	if err != nil {
		response.ReturnError(c, response.NOT_FOUND, "Host not found")
		return
	}

	// 创建SFTP客户端
	sftpClient, err := h.createSFTPClient(host)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, fmt.Sprintf("Failed to create SFTP client: %v", err))
		return
	}
	defer sftpClient.Close()

	// 重命名文件
	err = sftpClient.RenameFile(oldPath, newPath)
	if err != nil {
		response.ReturnError(c, response.INTERNAL, fmt.Sprintf("Failed to rename file: %v", err))
		return
	}

	response.ReturnSuccess(c)
}

// createSFTPClient 创建SFTP客户端
func (h *SftpHandler) createSFTPClient(host *cmdbModel.Host) (*ssh.SFTPClient, error) {
	// 智能选择IP地址
	targetIP, err := h.selectTargetIP(host)
	if err != nil {
		return nil, err
	}

	// 解密密码
	decryptedPassword := host.Password
	if host.Password != "" {
		var err error
		decryptedPassword, err = h.decryptPassword(host.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt password: %v", err)
		}
	}

	// 解密私钥
	decryptedPrivateKey := host.PrivateKey
	if host.PrivateKey != "" {
		var err error
		decryptedPrivateKey, err = h.decryptPassword(host.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt private key: %v", err)
		}
	}

	// 创建SSH配置
	sshConfig := &ssh.SSHClientConfig{
		Timeout:   time.Second * 30,
		IP:        targetIP,
		Port:      int(host.Port),
		UserName:  host.Username,
		Password:  decryptedPassword, // 使用解密后的密码
		PublicKey: decryptedPrivateKey, // 使用解密后的私钥
		AuthModel: h.getAuthModel(host),
	}

	return ssh.NewSFTPClient(sshConfig)
}

// selectTargetIP 智能选择目标IP
func (h *SftpHandler) selectTargetIP(host *cmdbModel.Host) (string, error) {
	// 优先使用公网IP
	if len(host.PublicIP) > 2 { // 检查是否为非空JSON数组
		var publicIPs []string
		if err := json.Unmarshal(host.PublicIP, &publicIPs); err == nil && len(publicIPs) > 0 {
			// 过滤空字符串，找到第一个有效IP
			for _, ip := range publicIPs {
				if ip != "" && strings.TrimSpace(ip) != "" {
					return ip, nil
				}
			}
		}
	}

	// 其次使用私网IP
	if len(host.PrivateIP) > 2 { // 检查是否为非空JSON数组
		var privateIPs []string
		if err := json.Unmarshal(host.PrivateIP, &privateIPs); err == nil && len(privateIPs) > 0 {
			// 过滤空字符串，找到第一个有效IP
			for _, ip := range privateIPs {
				if ip != "" && strings.TrimSpace(ip) != "" {
					return ip, nil
				}
			}
		}
	}

	return "", fmt.Errorf("no available IP address for host %s", host.Name)
}

// getAuthModel 获取认证模式
func (h *SftpHandler) getAuthModel(host *cmdbModel.Host) string {
	if host.PrivateKey != "" {
		return "PUBLICKEY"
	}
	return "PASSWORD"
}

// RegisterRoutes 注册SFTP相关路由
func (h *SftpHandler) RegisterRoutes(r *gin.RouterGroup) {
	sftpGroup := r.Group("/sftp")
	{
		sftpGroup.GET("/list", h.List)
		sftpGroup.POST("/upload", h.Upload)
		sftpGroup.GET("/download", h.Download)
		sftpGroup.DELETE("/delete", h.Delete)
		sftpGroup.POST("/mkdir", h.MakeDir)
		sftpGroup.POST("/rename", h.Rename)
	}
}

// decryptPassword 解密主机密码
func (h *SftpHandler) decryptPassword(encryptedPassword string) (string, error) {
	if encryptedPassword == "" {
		return "", nil
	}

	// 获取AES密钥
	var keys [][]byte
	if len(config.GlobalConfig.App.AesKeys) > 0 {
		// 使用多密钥
		for _, keyStr := range config.GlobalConfig.App.AesKeys {
			keys = append(keys, []byte(keyStr))
		}
		return encryption.DecryptAESWithKeys(encryptedPassword, keys)
	} else if config.GlobalConfig.App.AesKey != "" {
		// 使用单密钥
		return encryption.DecryptAES(encryptedPassword, []byte(config.GlobalConfig.App.AesKey))
	}

	return "", fmt.Errorf("no encryption key configured")
}
