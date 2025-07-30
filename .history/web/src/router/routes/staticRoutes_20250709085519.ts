import { AppRouteRecordRaw } from '../utils/utils'
import { RoutesAlias } from '../routesAlias'

/**
 * 静态路由配置
 * 不需要权限就能访问的路由
 */
export const staticRoutes: AppRouteRecordRaw[] = [
  {
    path: '/auth/login',
    name: 'Login',
    component: () => import('@/views/auth/login/index.vue'),
    meta: { title: '登录 - DevOps Asset', setTheme: true }
  },
  {
    path: '/auth/forget-password',
    name: 'ForgetPassword',
    component: () => import('@/views/auth/forget-password/index.vue'),
    meta: { title: '找回密码 - DevOps Asset', setTheme: true }
  },
  {
    path: '/exception',
    component: () => import('@views/index/index.vue'),
    name: 'Exception',
    meta: { title: 'menus.exception.title' },
    children: [
      {
        path: RoutesAlias.Exception403,
        name: 'Exception403',
        component: () => import('@views/exception/403/index.vue'),
        meta: { title: '403' }
      },
      {
        path: '/:catchAll(.*)',
        name: 'Exception404',
        component: () => import('@views/exception/404/index.vue'),
        meta: { title: '404' }
      },
      {
        path: RoutesAlias.Exception500,
        name: 'Exception500',
        component: () => import('@views/exception/500/index.vue'),
        meta: { title: '500' }
      }
    ]
  },
  {
    path: '/outside',
    component: () => import('@views/index/index.vue'),
    name: 'Outside',
    meta: { title: 'menus.outside.title' },
    children: [
      {
        path: '/outside/iframe/:path',
        name: 'Iframe',
        component: () => import('@/views/outside/Iframe.vue'),
        meta: { title: 'iframe' }
      }
    ]
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/views/dashboard/index/index.vue'),
    meta: { title: '仪表盘' }
  },
  {
    path: '/workbench',
    name: 'Workbench',
    component: () => import('@/views/dashboard/console/index.vue'),
    meta: { title: '工作台' }
  },
  {
    path: '/system-management',
    name: 'SystemManagement',
    component: () => import('@/views/system/user/index.vue'),
    meta: { title: '系统管理' }
  },
  {
    path: '/system',
    component: () => import('@views/index/index.vue'),
    name: 'System',
    meta: { title: '系统管理' },
    children: [
      {
        path: 'user-center',
        name: 'UserCenter',
        component: () => import('@/views/system/user-center/index.vue'),
        meta: { title: '个人中心', icon: 'user' }
      }
    ]
  }
]
