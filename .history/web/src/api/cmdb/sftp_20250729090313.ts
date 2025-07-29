import api from '@/api/client'

// 文件信息接口
export interface FileInfo {
  name: string
  size: number
  mode: string
  modTime: string
  isDir: boolean
}

// SFTP响应接口
export interface SftpResponse<T = any> {
  code: number
  message: string
  data: T
}

export const sftpApi = {
  // 列出文件
  listFiles(hostId: string, path: string = '/'): Promise<SftpResponse<FileInfo[]>> {
    return api.get({
      url: `/api/v1/cmdb/sftp/list`,
      params: {
        host_id: hostId,
        path: path
      }
    })
  },

  // 下载文件
  downloadFile(hostId: string, filePath: string): Promise<Blob> {
    return api.get({
      url: `/api/v1/cmdb/sftp/download`,
      params: {
        host_id: hostId,
        file_path: filePath
      },
      responseType: 'blob'
    })
  },

  // 上传文件
  uploadFile(hostId: string, remotePath: string, file: File): Promise<SftpResponse> {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('host_id', hostId)
    formData.append('path', remotePath)

    return api.post({
      url: `/api/v1/cmdb/sftp/upload`,
      data: formData,
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  // 删除文件
  deleteFile(hostId: string, filePath: string): Promise<SftpResponse> {
    return api.delete({
      url: `/api/v1/cmdb/sftp/delete`,
      params: {
        host_id: hostId,
        file_path: filePath
      }
    })
  },

  // 创建目录
  createDirectory(hostId: string, dirPath: string): Promise<SftpResponse> {
    const formData = new FormData()
    formData.append('host_id', hostId)
    formData.append('path', dirPath)

    return api.post({
      url: `/api/v1/cmdb/sftp/mkdir`,
      data: formData
    })
  },

  // 重命名文件
  renameFile(hostId: string, oldPath: string, newPath: string): Promise<SftpResponse> {
    const formData = new FormData()
    formData.append('host_id', hostId)
    formData.append('old_path', oldPath)
    formData.append('new_path', newPath)

    return api.post({
      url: `/api/v1/cmdb/sftp/rename`,
      data: formData
    })
  }
}
