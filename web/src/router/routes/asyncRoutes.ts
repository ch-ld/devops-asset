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
  // 添加CMDB主机管理模块
  {
    id: 2,
    name: 'CMDB',
    path: '/cmdb',
    component: RoutesAlias.Home,
    meta: {
      title: '资产管理',
      icon: '&#xe73e;', // 服务器图标
      keepAlive: false,
      roles: ['超级管理员', '管理员', '运维人员']
    },
    children: [
      {
        id: 201,
        path: 'overview',
        name: 'ResourceOverview',
        component: 'cmdb/resource/index',
        meta: {
          title: '资源概览',
          keepAlive: true,
          roles: ['超级管理员', '管理员', '运维人员']
        }
      },
      {
        id: 202,
        path: 'host',
        name: 'HostManagement',
        component: 'cmdb/host/index',
        meta: {
          title: '主机管理',
          keepAlive: true,
          roles: ['超级管理员', '管理员', '运维人员']
        }
      },
      {
        id: 203,
        path: 'host/dashboard',
        name: 'HostDashboard',
        component: 'cmdb/host/dashboard/index',
        meta: {
          title: '主机概览',
          keepAlive: true,
          roles: ['超级管理员', '管理员', '运维人员']
        }
      },
      {
        id: 204,
        path: 'host/detail/:id',
        name: 'HostDetail',
        component: 'cmdb/host/detail',
        meta: {
          title: '主机详情',
          hidden: true,
          keepAlive: false,
          roles: ['超级管理员', '管理员', '运维人员']
        }
      },
      {
        id: 205,
        path: 'host/edit/:id',
        name: 'HostEdit',
        component: 'cmdb/host/edit',
        meta: {
          title: '编辑主机',
          hidden: true,
          keepAlive: false,
          roles: ['超级管理员', '管理员', '运维人员']
        }
      },
      {
        id: 206,
        path: 'host/create',
        name: 'HostCreate',
        component: 'cmdb/host/edit',
        meta: {
          title: '添加主机',
          hidden: true,
          keepAlive: false,
          roles: ['超级管理员', '管理员', '运维人员']
        }
      },
      {
        id: 207,
        path: 'providers',
        name: 'ProviderManagement',
        component: 'cmdb/provider/index',
        meta: {
          title: '云账号管理',
          keepAlive: true,
          roles: ['超级管理员', '管理员']
        }
      },
      {
        id: 208,
        path: 'alerts',
        name: 'AlertManagement',
        component: 'cmdb/alert/index',
        meta: {
          title: '告警管理',
          keepAlive: true,
          roles: ['超级管理员', '管理员', '运维人员']
        }
      }
    ]
  }
]
