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

	cmdbModel "api-server/internal/model/cmdb"
	cmdbRepo "api-server/internal/repository/cmdb"
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

	// 创建SSH配置
	sshConfig := &ssh.SSHClientConfig{
		Timeout:   time.Second * 30,
		IP:        targetIP,
		Port:      int(host.Port),
		UserName:  host.Username,
		Password:  host.Password,
		PublicKey: host.PrivateKey,
		AuthModel: h.getAuthModel(host),
	}
	log.Printf("🔍 [WebSSH] SSH config - IP: %s, Port: %d, Username: %s, AuthModel: %s",
		sshConfig.IP, sshConfig.Port, sshConfig.UserName, sshConfig.AuthModel)

	// 创建SSH客户端
	log.Printf("🔍 [WebSSH] Attempting SSH connection...")
	sshClient, err := ssh.NewSSHClient(sshConfig)
	if err != nil {
		log.Printf("❌ [WebSSH] SSH connection failed: %v", err)
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

	// 创建WebSSH会话
	turn, err := ssh.NewTurn(wsConn, sshClient)
	if err != nil {
		errorMsg := fmt.Sprintf("连接失败: 无法创建终端会话 - %v", err)
		wsConn.WriteMessage(websocket.TextMessage, []byte(errorMsg))
		wsConn.Close()
		return
	}
	defer turn.Close()

	// 启动会话处理
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(2)

	// 读取WebSocket数据并写入SSH
	go func() {
		defer wg.Done()
		err := turn.LoopRead(ctx)
		if err != nil {
			log.Printf("LoopRead error: %v", err)
		}
	}()

	// 等待SSH会话结束
	go func() {
		defer wg.Done()
		err := turn.SessionWait()
		if err != nil {
			log.Printf("SessionWait error: %v", err)
		}
		cancel()
	}()

	wg.Wait()
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
	if host.PrivateKey != "" {
		return "PUBLICKEY"
	}
	return "PASSWORD"
}
