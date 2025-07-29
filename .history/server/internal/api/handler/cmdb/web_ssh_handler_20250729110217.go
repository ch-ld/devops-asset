package cmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"api-server/internal/config"
	cmdbModel "api-server/internal/model/cmdb"
	cmdbRepo "api-server/internal/repository/cmdb"
	"api-server/pkg/crypto/encryption"
	"api-server/pkg/ssh"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

// HandleSSH 处理WebSSH连接
func (h *WebSshHandler) HandleWebSSH(c *gin.Context) {
	log.Printf("🔍 [WebSSH] Starting WebSSH connection for host_id: %s", c.Query("host_id"))

	hostID, err := strconv.ParseUint(c.Query("host_id"), 10, 64)
	if err != nil {
		log.Printf("❌ [WebSSH] Invalid host_id: %v", err)
		c.JSON(400, gin.H{"error": "Invalid host_id"})
		return
	}
	log.Printf("🔍 [WebSSH] Parsed host_id: %d", hostID)

	// 获取主机信息
	host, err := h.hostRepo.FindByID(uint(hostID))
	if err != nil {
		log.Printf("❌ [WebSSH] Failed to get host: %v", err)
		c.JSON(500, gin.H{"error": "Failed to get host"})
		return
	}
	log.Printf("🔍 [WebSSH] Found host: %s (ID: %d)", host.Name, host.ID)
	log.Printf("🔍 [WebSSH] Host details - Username: %s, Port: %d", host.Username, host.Port)
	log.Printf("🔍 [WebSSH] Host PublicIP: %s", string(host.PublicIP))
	log.Printf("🔍 [WebSSH] Host PrivateIP: %s", string(host.PrivateIP))

	// 升级为WebSocket连接
	wsConn, err := ssh.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("❌ [WebSSH] Failed to upgrade WebSocket: %v", err)
		return
	}
	defer wsConn.Close()
	log.Printf("✅ [WebSSH] WebSocket upgraded successfully")

	// 智能选择IP地址
	targetIP, err := h.selectTargetIP(host)
	if err != nil {
		log.Printf("❌ [WebSSH] Failed to select target IP: %v", err)
		errorMsg := fmt.Sprintf("连接失败: %v", err)
		wsConn.WriteMessage(websocket.TextMessage, []byte(errorMsg))
		wsConn.Close()
		return
	}
	log.Printf("🔍 [WebSSH] Selected target IP: %s", targetIP)

	// 解密密码
	decryptedPassword := host.Password
	if host.Password != "" {
		var err error
		decryptedPassword, err = h.decryptPassword(host.Password)
		if err != nil {
			log.Printf("❌ [WebSSH] Failed to decrypt password: %v", err)
			errorMsg := "连接失败: 密码解密失败"
			wsConn.WriteMessage(websocket.TextMessage, []byte(errorMsg))
			wsConn.Close()
			return
		}
		log.Printf("🔍 [WebSSH] Password decrypted successfully")
	}

	// 创建SSH配置
	sshConfig := &ssh.SSHClientConfig{
		Timeout:   time.Second * 30,
		IP:        targetIP,
		Port:      int(host.Port),
		UserName:  host.Username,
		Password:  decryptedPassword, // 使用解密后的密码
		PublicKey: host.PrivateKey,
		AuthModel: h.getAuthModel(host),
	}
	log.Printf("🔍 [WebSSH] SSH config - IP: %s, Port: %d, Username: %s, AuthModel: %s",
		sshConfig.IP, sshConfig.Port, sshConfig.UserName, sshConfig.AuthModel)
	log.Printf("🔍 [WebSSH] Auth details - Password set: %t, PrivateKey set: %t",
		len(decryptedPassword) > 0, len(host.PrivateKey) > 0)

	// 创建SSH客户端
	log.Printf("🔍 [WebSSH] Attempting SSH connection...")
	sshClient, err := ssh.NewSSHClient(sshConfig)
	if err != nil {
		log.Printf("❌ [WebSSH] SSH connection failed: %v", err)

		// 更新SSH状态为错误
		h.updateSSHStatus(uint(hostID), "error", nil)

		var errorMsg string
		if strings.Contains(err.Error(), "authentication failed") {
			errorMsg = "连接失败: 用户名或密码错误"
		} else if strings.Contains(err.Error(), "connection refused") {
			errorMsg = "连接失败: 目标主机拒绝连接，请检查主机是否在线及SSH服务是否启动"
		} else if strings.Contains(err.Error(), "timeout") {
			errorMsg = "连接失败: 连接超时，请检查主机网络连通性"
		} else if strings.Contains(err.Error(), "no route to host") {
			errorMsg = "连接失败: 网络不可达，请检查主机IP地址"
		} else {
			errorMsg = fmt.Sprintf("连接失败: %v", err)
		}
		wsConn.WriteMessage(websocket.TextMessage, []byte(errorMsg))
		wsConn.Close()
		return
	}
	defer sshClient.Close()
	log.Printf("✅ [WebSSH] SSH connection established successfully")

	// 更新SSH状态为在线
	now := time.Now()
	h.updateSSHStatus(uint(hostID), "online", &now)

	// 创建WebSSH会话
	log.Printf("🔍 [WebSSH] Creating SSH session...")
	turn, err := ssh.NewTurn(wsConn, sshClient)
	if err != nil {
		log.Printf("❌ [WebSSH] Failed to create SSH session: %v", err)
		errorMsg := fmt.Sprintf("连接失败: 无法创建终端会话 - %v", err)
		wsConn.WriteMessage(websocket.TextMessage, []byte(errorMsg))
		wsConn.Close()
		return
	}
	defer turn.Close()
	log.Printf("✅ [WebSSH] SSH session created successfully")

	// 启动会话处理
	log.Printf("🔍 [WebSSH] Starting session handlers...")
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(2)

	// 读取WebSocket数据并写入SSH
	go func() {
		defer wg.Done()
		log.Printf("🔍 [WebSSH] Starting LoopRead goroutine...")
		err := turn.LoopRead(ctx)
		if err != nil {
			log.Printf("❌ [WebSSH] LoopRead error: %v", err)
		} else {
			log.Printf("✅ [WebSSH] LoopRead completed")
		}
	}()

	// 等待SSH会话结束
	go func() {
		defer wg.Done()
		log.Printf("🔍 [WebSSH] Starting SessionWait goroutine...")
		err := turn.SessionWait()
		if err != nil {
			log.Printf("❌ [WebSSH] SessionWait error: %v", err)
		} else {
			log.Printf("✅ [WebSSH] SessionWait completed")
		}
		cancel()
	}()

	log.Printf("🔍 [WebSSH] Waiting for session to complete...")
	wg.Wait()
	log.Printf("✅ [WebSSH] Session completed")

	// 更新SSH状态为离线
	h.updateSSHStatus(uint(hostID), "offline", nil)
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
		log.Printf("❌ [WebSSH] Failed to find host %d for status update: %v", hostID, err)
		return
	}

	// 更新SSH状态
	host.SSHStatus = status
	if connectedAt != nil {
		host.LastConnectedAt = connectedAt
	}

	if err := h.hostRepo.Update(host); err != nil {
		log.Printf("❌ [WebSSH] Failed to update SSH status for host %d: %v", hostID, err)
	} else {
		log.Printf("✅ [WebSSH] Updated SSH status for host %d to %s", hostID, status)
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
