import { RoutesAlias } from '../routesAlias'
import { AppRouteRecord } from '@/types/router'

/**
 * 菜单列表、异步路由
 *
 * 支持两种模式:
 * 1. 前端静态配置 - 直接使用本文件中定义的路由配置
 * 2. 后端动态配置 - 后端返回菜单数据，前端解析生成路由
 *
 * 菜单标题（title）:
 * 可以是 i18n 的 key，也可以是字符串，比如：'用户列表'
 */
export const asyncRoutes: AppRouteRecord[] = [
  {
    id: 1,
    name: 'Dashboard',
    path: '/dashboard',
    component: RoutesAlias.Home,
    meta: {
      title: 'DevOps 概览',
      icon: '&#xe721;',
      keepAlive: false,
      roles: ['超级管理员', '管理员']
    },
    children: [
      {
        id: 101,
        path: 'console',
        name: 'Console',
        component: RoutesAlias.Dashboard,
        meta: {
          title: '资产概览',
          keepAlive: false,
          roles: ['超级管理员', '管理员']
        }
      }
    ]
  },
  {
    id: 2,
    name: 'CMDB',
    path: '/cmdb',
    component: RoutesAlias.Home,
    meta: {
      title: '资源管理',
      icon: '&#xe67e;',
      keepAlive: false,
      roles: ['超级管理员', '管理员', '运维']
    },
    children: [
      {
        id: 201,
        path: 'host',
        name: 'HostList',
        component: () => import('@/views/cmdb/host/index.vue'),
        meta: {
          title: '主机管理',
          keepAlive: true,
          roles: ['超级管理员', '管理员', '运维']
        }
      },
      {
        id: 202,
        path: 'host/dashboard',
        name: 'HostDashboard',
        component: () => import('@/views/cmdb/host/Dashboard.vue'),
        meta: {
          title: '主机概览',
          keepAlive: false,
          hideInMenu: true,
          roles: ['超级管理员', '管理员', '运维']
        }
      },
      {
        id: 203,
        path: 'host/detail/:id',
        name: 'HostDetail',
        component: () => import('@/views/cmdb/host/HostDetail.vue'),
        meta: {
          title: '主机详情',
          keepAlive: false,
          hideInMenu: true,
          roles: ['超级管理员', '管理员', '运维']
        }
      },
      {
        id: 204,
        path: 'provider',
        name: 'ProviderList',
        component: () => import('@/views/cmdb/provider/index.vue'),
        meta: {
          title: '云账号管理',
          keepAlive: true,
          roles: ['超级管理员', '管理员']
        }
      }
    ]
  }
]
