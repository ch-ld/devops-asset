package cmdb

import "time"

// BatchChangeStatusRequest 批量状态变更请求
type BatchChangeStatusRequest struct {
	IDs    []uint `json:"ids"`
	Status string `json:"status"`
}

// BatchLifecycleRequest 批量生命周期请求
type BatchLifecycleRequest struct {
	IDs       []uint     `json:"ids"`
	ExpiredAt *time.Time `json:"expired_at"`
	Status    string     `json:"status"`
	Recycle   bool       `json:"recycle"`
}

// BatchSetCustomFieldsRequest 批量自定义字段请求
type BatchSetCustomFieldsRequest struct {
	IDs         []uint                 `json:"ids"`
	ExtraFields map[string]interface{} `json:"extra_fields"`
}

// BatchSSHRequest 批量SSH请求
type BatchSSHRequest struct {
	IDs     []uint `json:"ids"`
	Cmd     string `json:"cmd"`
	Timeout int    `json:"timeout"`
}

// BatchSFTPRequest 批量SFTP请求
type BatchSFTPRequest struct {
	IDs        []uint `json:"ids"`
	LocalPath  string `json:"local_path"`
	RemotePath string `json:"remote_path"`
	Operation  string `json:"operation"` // upload, download, delete, mkdir, list
}
