package cmdb

import (
	"api-server/internal/config"
	"api-server/internal/repository/cmdb"
	"api-server/pkg/crypto/encryption"
	"api-server/pkg/ssh"
	"fmt"
	"os"
	"time"

	gossh "golang.org/x/crypto/ssh"
)

// SftpService SFTP服务
type SftpService struct {
	hostRepo *cmdb.HostRepository
}

// NewSftpService 创建SFTP服务实例
func NewSftpService(hostRepo *cmdb.HostRepository) *SftpService {
	return &SftpService{hostRepo: hostRepo}
}

// getSSHClient 创建SSH客户端连接
func (s *SftpService) getSSHClient(hostID uint) (*gossh.Client, error) {
	// 获取主机信息
	host, err := s.hostRepo.FindByID(hostID)
	if err != nil {
		return nil, fmt.Errorf("主机不存在: %w", err)
	}

	// 智能选择IP地址 (默认自动选择)
	targetIP, err := s.selectTargetIP(host, "auto")
	if err != nil {
		return nil, fmt.Errorf("IP地址选择失败: %w", err)
	}

	// 检查认证信息
	if host.Username == "" {
		return nil, fmt.Errorf("主机SSH用户名未配置")
	}

	var password string
	var privateKey string
	port := host.Port
	if port == 0 {
		port = 22 // 默认端口
	}

	// 解密密码（如果存在）
	if host.Password != "" {
		key := []byte(config.GlobalConfig.App.AesKey)
		var err error
		password, err = encryption.DecryptAES(host.Password, key)
		if err != nil {
			return nil, fmt.Errorf("解密密码失败: %w", err)
		}
	}

	// 解密私钥（如果存在）
	if host.PrivateKey != "" {
		key := []byte(config.GlobalConfig.App.AesKey)
		var err error
		privateKey, err = encryption.DecryptAES(host.PrivateKey, key)
		if err != nil {
			return nil, fmt.Errorf("解密私钥失败: %w", err)
		}
	}

	// 创建SSH客户端配置
	sshConfig := &gossh.ClientConfig{
		User:            host.Username,
		Auth:            []gossh.AuthMethod{},
		HostKeyCallback: gossh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	// 添加认证方法
	if privateKey != "" {
		// 优先使用私钥认证
		signer, err := gossh.ParsePrivateKey([]byte(privateKey))
		if err == nil {
			sshConfig.Auth = append(sshConfig.Auth, gossh.PublicKeys(signer))
		}
	}

	if password != "" {
		// 添加密码认证
		sshConfig.Auth = append(sshConfig.Auth, gossh.Password(password))
	}

	// 检查是否有可用的认证方法
	if len(sshConfig.Auth) == 0 {
		return nil, fmt.Errorf("主机SSH凭证未配置")
	}

	// 连接SSH服务器，使用动态端口
	client, err := gossh.Dial("tcp", fmt.Sprintf("%s:%d", targetIP, port), sshConfig)
	if err != nil {
		return nil, fmt.Errorf("SSH连接失败: %w", err)
	}

	return client, nil
}

// getSftpClient 创建SFTP客户端
func (s *SftpService) getSftpClient(hostID uint) (*ssh.SftpClient, *gossh.Client, error) {
	// 创建SSH客户端
	sshClient, err := s.getSSHClient(hostID)
	if err != nil {
		return nil, nil, err
	}

	// 创建SFTP客户端
	sftpClient, err := ssh.NewSftpClient(sshClient)
	if err != nil {
		sshClient.Close()
		return nil, nil, fmt.Errorf("创建SFTP客户端失败: %w", err)
	}

	return sftpClient, sshClient, nil
}

// ListFiles 列出指定路径下的文件和目录
func (s *SftpService) ListFiles(hostID uint, path string) ([]os.FileInfo, error) {
	sftpClient, sshClient, err := s.getSftpClient(hostID)
	if err != nil {
		return nil, err
	}
	defer sshClient.Close()
	defer sftpClient.Close()

	// 列出文件
	files, err := sftpClient.List(path)
	if err != nil {
		return nil, fmt.Errorf("列出文件失败: %w", err)
	}

	return files, nil
}

// UploadFile 上传文件到主机
func (s *SftpService) UploadFile(hostID uint, localPath, remotePath string) error {
	sftpClient, sshClient, err := s.getSftpClient(hostID)
	if err != nil {
		return err
	}
	defer sshClient.Close()
	defer sftpClient.Close()

	// 上传文件
	if err := sftpClient.Upload(localPath, remotePath); err != nil {
		return fmt.Errorf("上传文件失败: %w", err)
	}

	return nil
}

// DownloadFile 从主机下载文件
func (s *SftpService) DownloadFile(hostID uint, remotePath, localPath string) error {
	sftpClient, sshClient, err := s.getSftpClient(hostID)
	if err != nil {
		return err
	}
	defer sshClient.Close()
	defer sftpClient.Close()

	// 下载文件
	if err := sftpClient.Download(remotePath, localPath); err != nil {
		return fmt.Errorf("下载文件失败: %w", err)
	}

	return nil
}

// DeleteFile 删除主机上的文件
func (s *SftpService) DeleteFile(hostID uint, remotePath string) error {
	sftpClient, sshClient, err := s.getSftpClient(hostID)
	if err != nil {
		return err
	}
	defer sshClient.Close()
	defer sftpClient.Close()

	// 删除文件
	if err := sftpClient.Delete(remotePath); err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	return nil
}

// MakeDir 在主机上创建目录
func (s *SftpService) MakeDir(hostID uint, remotePath string) error {
	sftpClient, sshClient, err := s.getSftpClient(hostID)
	if err != nil {
		return err
	}
	defer sshClient.Close()
	defer sftpClient.Close()

	// 创建目录
	if err := sftpClient.Mkdir(remotePath); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	return nil
}

// RenameFile 重命名主机上的文件或目录
func (s *SftpService) RenameFile(hostID uint, oldPath, newPath string) error {
	sftpClient, sshClient, err := s.getSftpClient(hostID)
	if err != nil {
		return err
	}
	defer sshClient.Close()
	defer sftpClient.Close()

	// 重命名文件
	if err := sftpClient.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("重命名失败: %w", err)
	}

	return nil
}
