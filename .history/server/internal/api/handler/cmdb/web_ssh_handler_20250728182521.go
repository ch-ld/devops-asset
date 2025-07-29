package cmdb

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	cmdbModel "api-server/internal/model/cmdb"
	cmdbRepo "api-server/internal/repository/cmdb"
	"api-server/pkg/ssh"

	"github.com/gin-gonic/gin"
)

// WebSshHandler WebSSH处理器
type WebSshHandler struct {
	hostRepo cmdbRepo.HostRepository
}

// NewWebSshHandler 创建WebSSH处理器
func NewWebSshHandler(hostRepo cmdbRepo.HostRepository) *WebSshHandler {
	return &WebSshHandler{
		hostRepo: hostRepo,
	}
}

// HandleSSH 处理WebSSH连接
func (h *WebSshHandler) HandleSSH(c *gin.Context) {
	// 升级为WebSocket连接
	wsConn, err := ssh.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer wsConn.Close()

	// 获取主机ID
	hostIDStr := c.Query("host_id")
	hostID, err := strconv.Atoi(hostIDStr)
	if err != nil {
		wsConn.WriteMessage(1, []byte("Error: Invalid host_id parameter"))
		return
	}

	// 查询主机信息
	host, err := h.hostRepo.FindByID(uint(hostID))
	if err != nil {
		wsConn.WriteMessage(1, []byte("Error: Host not found"))
		return
	}

	// 智能选择IP地址
	targetIP, err := h.selectTargetIP(host)
	if err != nil {
		wsConn.WriteMessage(1, []byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	// 创建SSH配置
	sshConfig := &ssh.SSHClientConfig{
		Timeout:   time.Second * 30,
		IP:        targetIP,
		Port:      int(host.SSHPort),
		UserName:  host.SSHUser,
		Password:  host.SSHPassword,
		PublicKey: host.PrivateKey,
		AuthModel: h.getAuthModel(host),
	}

	// 创建SSH客户端
	sshClient, err := ssh.NewSSHClient(sshConfig)
	if err != nil {
		wsConn.WriteMessage(1, []byte(fmt.Sprintf("SSH connection failed: %v", err)))
		return
	}
	defer sshClient.Close()

	// 创建WebSSH会话
	turn, err := ssh.NewTurn(wsConn, sshClient)
	if err != nil {
		wsConn.WriteMessage(1, []byte(fmt.Sprintf("Failed to create terminal: %v", err)))
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
	if host.PublicIP != "" {
		return host.PublicIP, nil
	}

	// 其次使用私网IP
	if host.PrivateIP != "" {
		return host.PrivateIP, nil
	}

	// 最后使用管理IP
	if host.ManagementIP != "" {
		return host.ManagementIP, nil
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
