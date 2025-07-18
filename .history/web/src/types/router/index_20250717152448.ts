/**
 * 路由相关类型定义
 */

import type { RouteRecordRaw, RouteComponent } from 'vue-router'
import { MenuTypeEnum } from '@/enums/appEnum'

// 扩展元数据
export interface RouteMeta {
  title?: string
  icon?: string
  hidden?: boolean
  isHide?: boolean
  isLink?: string
  isIframe?: boolean
  isFirstLevel?: boolean
  activePath?: string
  keepAlive?: boolean
  roles?: string[]
  auths?: string[]
  showBadge?: boolean
  showTextBadge?: string | boolean
}

// 应用路由记录
export interface AppRouteRecord extends Omit<RouteRecordRaw, 'meta' | 'children'> {
  id?: number
  path: string
  component?: RouteComponent | string
  name?: string | symbol
  redirect?: string
  meta: RouteMeta
  children?: AppRouteRecord[]
}

// 菜单配置
export interface MenuConfig {
  type?: MenuTypeEnum
  width?: number
  miniWidth?: number
  theme?: 'light' | 'dark'
  openWidth?: number
}

// 菜单选项
export interface MenuOption {
  id?: number | string
  name: string
  path: string
  icon?: string
  children?: MenuOption[]
}

// 选项卡项目
export interface TabItem {
  name: string
  path: string
  title: string
  icon?: string
  close?: boolean
  isIframe?: boolean
}
