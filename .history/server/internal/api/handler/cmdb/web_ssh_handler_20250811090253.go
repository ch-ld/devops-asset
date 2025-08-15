package cmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"api-server/internal/config"
	cmdbModel "api-server/internal/model/cmdb"
	cmdbRepo "api-server/internal/repository/cmdb"
	"api-server/pkg/crypto/encryption"
	"api-server/pkg/ssh"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// WebSshHandler WebSSH处理器
type WebSshHandler struct {
	hostRepo *cmdbRepo.HostRepository
}

// NewWebSshHandler 创建WebSSH处理器
func NewWebSshHandler(hostRepo *cmdbRepo.HostRepository) *WebSshHandler {
	return &WebSshHandler{
		hostRepo: hostRepo,
	}
}

// HandleWebSSH 处理WebSSH连接
// @Summary 建立WebSSH连接
// @Description 通过WebSocket建立SSH连接到指定主机
// @Tags WebSSH
// @Accept json
// @Produce json
// @Param host_id query int true "主机ID"
// @Success 101 {string} string "WebSocket连接已建立"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /api/v1/cmdb/ssh/webssh [get]
func (h *WebSshHandler) HandleWebSSH(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.Query("host_id"), 10, 64)
	if err != nil {
		zap.L().Error("Invalid host_id parameter", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid host_id"})
		return
	}

	// 获取主机信息
	host, err := h.hostRepo.FindByID(uint(hostID))
	if err != nil {
		zap.L().Error("Failed to get host information",
			zap.Uint64("host_id", hostID),
			zap.Error(err))
		c.JSON(500, gin.H{"error": "Failed to get host"})
		return
	}

	// 升级为WebSocket连接
	wsConn, err := ssh.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zap.L().Error("Failed to upgrade WebSocket connection", zap.Error(err))
		return
	}
	defer wsConn.Close()

	// 智能选择IP地址
	targetIP, err := h.selectTargetIP(host)
	if err != nil {
		zap.L().Error("Failed to select target IP",
			zap.String("host", host.Name),
			zap.Error(err))
		errorMsg := fmt.Sprintf("连接失败: %v", err)
		wsConn.WriteMessage(websocket.TextMessage, []byte(errorMsg))
		wsConn.Close()
		return
	}

	// 解密密码
	decryptedPassword := host.Password
	if host.Password != "" {
		var err error
		decryptedPassword, err = h.decryptPassword(host.Password)
		if err != nil {
			zap.L().Error("Failed to decrypt password",
				zap.String("host", host.Name),
				zap.Error(err))
			errorMsg := "连接失败: 密码解密失败"
			wsConn.WriteMessage(websocket.TextMessage, []byte(errorMsg))
			wsConn.Close()
			return
		}
	}

	// 解密私钥
	decryptedPrivateKey := host.PrivateKey
	if host.PrivateKey != "" {
		var err error
		decryptedPrivateKey, err = h.decryptPassword(host.PrivateKey)
		if err != nil {
			zap.L().Error("Failed to decrypt private key",
				zap.String("host", host.Name),
				zap.Error(err))
			errorMsg := "连接失败: 私钥解密失败"
			wsConn.WriteMessage(websocket.TextMessage, []byte(errorMsg))
			wsConn.Close()
			return
		}
	}

	// 创建SSH配置
	sshConfig := &ssh.SSHClientConfig{
		Timeout:   time.Second * 30,
		IP:        targetIP,
		Port:      int(host.Port),
		UserName:  host.Username,
		Password:  decryptedPassword,
		PublicKey: decryptedPrivateKey,
		AuthModel: h.getAuthModel(host),
	}

	// 创建SSH客户端
	sshClient, err := ssh.NewSSHClient(sshConfig)
	if err != nil {
		zap.L().Error("SSH connection failed",
			zap.String("host", host.Name),
			zap.String("ip", targetIP),
			zap.Int("port", int(host.Port)),
			zap.String("username", host.Username),
			zap.Error(err))

		// 发送错误信息到前端
		errorMsg := "SSH连接失败"
		if err.Error() != "" {
			errorMsg = fmt.Sprintf("SSH连接失败: %s", err.Error())
		}
		wsConn.WriteMessage(websocket.TextMessage, []byte(errorMsg))
		wsConn.Close()

		// 更新主机SSH状态为失败
		h.updateSSHStatus(uint(hostID), "failed", nil)
		return
	}
	defer func() {
		if sshClient != nil {
			sshClient.Close()
		}
	}()

	// 更新主机SSH状态为已连接
	now := time.Now()
	h.updateSSHStatus(uint(hostID), "connected", &now)

	// 创建SSH会话处理器
	sshHandler, err := ssh.NewTurn(wsConn, sshClient)
	if err != nil {
		zap.L().Error("Failed to create SSH handler",
			zap.String("host", host.Name),
			zap.Error(err))
		wsConn.WriteMessage(websocket.TextMessage, []byte("处理程序创建失败"))
		wsConn.Close()
		h.updateSSHStatus(uint(hostID), "failed", nil)
		return
	}
	defer sshHandler.Close()

	// 创建context用于控制goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动数据读取循环
	go func() {
		defer func() {
			if r := recover(); r != nil {
				zap.L().Error("LoopRead panicked", zap.Any("panic", r))
			}
			cancel() // 确保context被取消
		}()
		if err := sshHandler.LoopRead(ctx); err != nil {
			zap.L().Warn("LoopRead ended", zap.Error(err))
		}
	}()

	// 启动会话等待
	go func() {
		defer func() {
			if r := recover(); r != nil {
				zap.L().Error("SessionWait panicked", zap.Any("panic", r))
			}
			cancel() // 确保context被取消
		}()
		if err := sshHandler.SessionWait(); err != nil {
			zap.L().Warn("SessionWait ended", zap.Error(err))
		}
	}()

	// 等待context被取消（任一goroutine结束都会取消context）
	<-ctx.Done()

	// 更新主机SSH状态为已断开
	h.updateSSHStatus(uint(hostID), "disconnected", nil)
}

// selectTargetIP 智能选择目标IP
func (h *WebSshHandler) selectTargetIP(host *cmdbModel.Host) (string, error) {
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
func (h *WebSshHandler) getAuthModel(host *cmdbModel.Host) string {
	// 优先使用数据库中设置的认证类型
	if host.AuthType != "" {
		switch host.AuthType {
		case "password":
			return "PASSWORD"
		case "privatekey":
			return "PUBLICKEY"
		case "both":
			// 如果同时支持两种方式，优先使用私钥
			if host.PrivateKey != "" {
				return "PUBLICKEY"
			}
			return "PASSWORD"
		}
	}

	// 兼容性处理：如果没有设置AuthType，按原来的逻辑判断
	if host.PrivateKey != "" {
		return "PUBLICKEY"
	}
	return "PASSWORD"
}

// updateSSHStatus 更新主机SSH状态
func (h *WebSshHandler) updateSSHStatus(hostID uint, status string, connectedAt *time.Time) {
	// 先获取主机信息
	host, err := h.hostRepo.FindByID(hostID)
	if err != nil {
		zap.L().Error("Failed to find host for status update",
			zap.Uint("host_id", hostID),
			zap.Error(err))
		return
	}

	// 更新SSH状态
	host.SSHStatus = status
	if connectedAt != nil {
		host.LastConnectedAt = connectedAt
	}

	if err := h.hostRepo.Update(host); err != nil {
		zap.L().Error("Failed to update SSH status for host",
			zap.Uint("host_id", hostID),
			zap.String("status", status),
			zap.Error(err))
	} else {
		zap.L().Info("Updated SSH status for host",
			zap.Uint("host_id", hostID),
			zap.String("status", status))
	}
}

// decryptPassword 解密主机密码
func (h *WebSshHandler) decryptPassword(encryptedPassword string) (string, error) {
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
