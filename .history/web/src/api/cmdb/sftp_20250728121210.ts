import request from '@/utils/http'

export interface FileInfo {
  name: string
  size: number
  mode: string
  modTime: string
  isDir: boolean
}

export interface SftpResponse<T = any> {
  code: number
  message: string
  data: T
}

export const sftpApi = {
  // 列出文件
  listFiles(hostId: string, path: string = '/'): Promise<SftpResponse<FileInfo[]>> {
    return request.get({
      url: `/api/v1/cmdb/sftp/list`,
      params: {
        host_id: hostId,
        path: path
      }
    })
  },

  // 下载文件
  downloadFile(hostId: string, filePath: string): Promise<Blob> {
    return request.get({
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

    return request.post({
      url: `/api/v1/cmdb/sftp/upload`,
      data: formData,
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  // 删除文件或文件夹
  deleteFile(hostId: string, filePath: string): Promise<SftpResponse> {
    return request.del({
      url: `/api/v1/cmdb/sftp/delete`,
      params: {
        host_id: hostId,
        path: filePath
      }
    })
  },

  // 创建文件夹
  createFolder(hostId: string, folderPath: string): Promise<SftpResponse> {
    const formData = new FormData()
    formData.append('host_id', hostId)
    formData.append('path', folderPath)

    return request({
      url: `/cmdb/sftp/mkdir`,
      method: 'post',
      data: formData,
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  // 重命名文件或文件夹
  renameFile(hostId: string, oldPath: string, newPath: string): Promise<SftpResponse> {
    const formData = new FormData()
    formData.append('host_id', hostId)
    formData.append('old_path', oldPath)
    formData.append('new_path', newPath)

    return request({
      url: `/cmdb/sftp/rename`,
      method: 'post',
      data: formData,
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  // 获取文件内容（用于编辑文本文件）
  getFileContent(hostId: string, filePath: string): Promise<SftpResponse<string>> {
    return request({
      url: `/cmdb/sftp/content`,
      method: 'get',
      params: {
        host_id: hostId,
        file_path: filePath
      }
    })
  },

  // 保存文件内容
  saveFileContent(hostId: string, filePath: string, content: string): Promise<SftpResponse> {
    return request({
      url: `/cmdb/sftp/content`,
      method: 'put',
      data: {
        host_id: hostId,
        file_path: filePath,
        content: content
      }
    })
  },

  // 获取文件权限
  getFilePermissions(hostId: string, filePath: string): Promise<SftpResponse<any>> {
    return request({
      url: `/cmdb/sftp/permissions`,
      method: 'get',
      params: {
        host_id: hostId,
        file_path: filePath
      }
    })
  },

  // 设置文件权限
  setFilePermissions(hostId: string, filePath: string, mode: string): Promise<SftpResponse> {
    return request({
      url: `/cmdb/sftp/permissions`,
      method: 'put',
      data: {
        host_id: hostId,
        file_path: filePath,
        mode: mode
      }
    })
  }
}
