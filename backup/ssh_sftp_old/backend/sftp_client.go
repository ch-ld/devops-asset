package ssh

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// SftpClient wraps the sftp.Client to provide a simpler interface.
type SftpClient struct {
	sftpClient *sftp.Client
}

// NewSftpClient creates a new SftpClient.
func NewSftpClient(sshClient *ssh.Client) (*SftpClient, error) {
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create sftp client: %w", err)
	}
	return &SftpClient{sftpClient: sftpClient}, nil
}

// Client returns the underlying sftp.Client.
func (c *SftpClient) Client() *sftp.Client {
	return c.sftpClient
}

// Close closes the sftp client.
func (c *SftpClient) Close() {
	if c.sftpClient != nil {
		c.sftpClient.Close()
	}
	// Note: The SSH client is managed by the caller, so we don't close it here.
}

// List returns a list of file entries for the given path.
func (c *SftpClient) List(remotePath string) ([]os.FileInfo, error) {
	return c.sftpClient.ReadDir(remotePath)
}

// Upload uploads a local file to a remote path.
func (c *SftpClient) Upload(localPath, remotePath string) error {
	localFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %w", err)
	}
	defer localFile.Close()

	remoteFile, err := c.sftpClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("failed to create remote file: %w", err)
	}
	defer remoteFile.Close()

	_, err = io.Copy(remoteFile, localFile)
	if err != nil {
		return fmt.Errorf("failed to upload file content: %w", err)
	}

	return nil
}

// Download downloads a remote file to a local path.
func (c *SftpClient) Download(remotePath, localPath string) error {
	remoteFile, err := c.sftpClient.Open(remotePath)
	if err != nil {
		return fmt.Errorf("failed to open remote file: %w", err)
	}
	defer remoteFile.Close()

	// Ensure local directory exists
	if err := os.MkdirAll(filepath.Dir(localPath), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create local directory: %w", err)
	}

	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		return fmt.Errorf("failed to download file content: %w", err)
	}

	return nil
}

// Delete deletes a remote file.
func (c *SftpClient) Delete(remotePath string) error {
	return c.sftpClient.Remove(remotePath)
}

// Mkdir creates a new remote directory.
func (c *SftpClient) Mkdir(remotePath string) error {
	return c.sftpClient.Mkdir(remotePath)
}

// Rename renames or moves a remote file.
func (c *SftpClient) Rename(oldPath, newPath string) error {
	return c.sftpClient.Rename(oldPath, newPath)
}
