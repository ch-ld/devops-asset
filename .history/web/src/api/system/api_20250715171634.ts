import api from '@/api/client'
import { ApiResponse } from '@/api/client'

interface CaptchaResponse {
  id: string
  image: string
}

export const getCaptcha = (
  height: number,
  width: number
): Promise<ApiResponse<CaptchaResponse>> => {
  return api.get({ url: '/api/v1/auth/captcha', params: { height, width } })
}

// 登录响应
interface LoginResponse {
  access_token: string
  expires_at: number
  token_type: string
}

export const userLogin = (data: {
  username: string
  password: string
  captcha: string
  captcha_id: string
}): Promise<ApiResponse<LoginResponse>> => {
  return api.post({ url: '/api/v1/auth/login', data })
}

// 用户信息响应
interface UserInfoResponse {
  id: number
  username: string
  nickname: string
  avatar: string
  email: string
  phone: string
  role_id: number
  department_id: number
  status: number
  remark: string
  role: {
    id: number
    name: string
    code: string
  }
  department: {
    id: number
    name: string
  }
}

export const getUserInfo = (): Promise<ApiResponse<UserInfoResponse>> => {
  return api.get({ url: '/api/v1/users/info' })
}

export const updateUserInfo = (data: any): Promise<ApiResponse<ApiResponse>> => {
  return api.put({ url: '/api/v1/users/info', data })
}

export const changePassword = (data: {
  old_password: string
  new_password: string
}): Promise<ApiResponse<ApiResponse>> => {
  return api.post({ url: '/api/v1/users/password', data })
}

export const getUserMenu = (): Promise<ApiResponse<ApiResponse>> => {
  return api.get({ url: '/api/v1/users/menu' })
}

export const getAllMenu = (): Promise<ApiResponse<ApiResponse>> => {
  return api.get({ url: '/api/v1/menus' })
}

// 添加菜单
export const addMenu = (data: any): Promise<ApiResponse<ApiResponse>> => {
  return api.post({ url: '/api/v1/menus', data })
}

// 更新菜单
export const updateMenu = (data: any): Promise<ApiResponse<ApiResponse>> => {
  return api.put({ url: `/api/v1/menus`, data })
}

// 删除菜单
export const deleteMenu = (id: string | number): Promise<ApiResponse<ApiResponse>> => {
  return api.del({ url: `/api/v1/menus`, data: { id } })
}

// 新增权限
export const addAuth = (data: any): Promise<ApiResponse<ApiResponse>> => {
  return api.post({ url: '/api/v1/menu-auths', data })
}

// 更新权限
export const updateAuth = (data: any): Promise<ApiResponse<ApiResponse>> => {
  return api.put({ url: `/api/v1/menu-auths`, data })
}

// 删除权限
export const deleteAuth = (id: number): Promise<ApiResponse<ApiResponse>> => {
  return api.del({ url: `/api/v1/menu-auths`, data: { id } })
}

// 获取权限列表
export const getAuthList = (menuID: number): Promise<ApiResponse<ApiResponse>> => {
  return api.get({ url: `/api/v1/menu-auths?menu_id=${menuID}` })
}

export const getDepartmentList = (params?: any): Promise<ApiResponse<ApiResponse>> => {
  return api.get({ url: '/api/v1/departments', params })
}

export const addDepartment = (data: any): Promise<ApiResponse<ApiResponse>> => {
  return api.post({ url: '/api/v1/departments', data })
}

export const updateDepartment = (data: any): Promise<ApiResponse<ApiResponse>> => {
  return api.put({ url: '/api/v1/departments', data })
}

export const deleteDepartment = (id: number): Promise<ApiResponse<ApiResponse>> => {
  return api.del({ url: `/api/v1/departments`, data: { id } })
}

export const getRoleList = (params?: any): Promise<ApiResponse<ApiResponse>> => {
  return api.get({ url: '/api/v1/roles', params })
}

export const addRole = (data: any): Promise<ApiResponse<ApiResponse>> => {
  return api.post({ url: '/api/v1/roles', data })
}

export const updateRole = (data: any): Promise<ApiResponse<ApiResponse>> => {
  return api.put({ url: '/api/v1/roles', data })
}

export const deleteRole = (id: number): Promise<ApiResponse<ApiResponse>> => {
  return api.del({ url: `/api/v1/roles`, data: { id } })
}

export const getAllMenuByRole = (roleID: number): Promise<ApiResponse<ApiResponse>> => {
  return api.get({ url: `/api/v1/menus/role?role_id=${roleID}` })
}

export const saveRolePermission = (data: any): Promise<ApiResponse<ApiResponse>> => {
  return api.put({ url: '/api/v1/menus/role', data })
}

export const getUserList = (params: any): Promise<ApiResponse<ApiResponse>> => {
  return api.get({ url: '/api/v1/users', params })
}

export const addUser = (data: any): Promise<ApiResponse<ApiResponse>> => {
  return api.post({ url: '/api/v1/users', data })
}

export const updateUser = (data: any): Promise<ApiResponse<ApiResponse>> => {
  return api.put({ url: '/api/v1/users', data })
}

export const deleteUser = (id: number): Promise<ApiResponse<ApiResponse>> => {
  return api.del({ url: `/api/v1/users/${id}` })
}
