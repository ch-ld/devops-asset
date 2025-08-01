/**
 * 动态路由处理
 * 根据接口返回的菜单列表注册动态路由
 */
import type { Router, RouteRecordRaw } from 'vue-router'
import type { AppRouteRecord } from '@/types/router'
import { saveIframeRoutes } from './menuToRouter'
import { RoutesAlias } from '../routesAlias'
import { h } from 'vue'
import { useMenuStore } from '@/store/modules/menu'

/**
 * 动态导入 views 目录下所有 .vue 组件
 */
const modules: Record<string, () => Promise<any>> = import.meta.glob('../../views/**/*.vue')

/**
 * 注册异步路由
 * 将接口返回的菜单列表转换为 Vue Router 路由配置，并添加到传入的 router 实例中
 * @param router Vue Router 实例
 * @param menuList 接口返回的菜单列表
 */
export function registerDynamicRoutes(router: Router, menuList: AppRouteRecord[]): void {
  console.log('registerDynamicRoutes 开始注册路由，菜单列表:', menuList)

  /**
   * 递归注册所有有 component 的节点
   */
  function register(menu: AppRouteRecord[]) {
    menu.forEach((route) => {
      console.log('处理路由:', route.name, 'component:', route.component)

      // 只要有 component 就注册
      if (route.component) {
        try {
          const routeConfig = convertRouteComponent(route, [])
          console.log('转换后的路由配置:', routeConfig)
          router.addRoute(routeConfig as RouteRecordRaw)
          console.log('成功注册路由:', route.name)
        } catch (error) {
          console.error('注册路由失败:', route.name, error)
        }
      } else {
        console.log('跳过没有component的路由:', route.name)
      }

      // 递归注册 children
      if (Array.isArray(route.children) && route.children.length > 0) {
        console.log('注册子路由:', route.name, 'children count:', route.children.length)
        register(route.children)
      }
    })
  }

  register(menuList)
  console.log('registerDynamicRoutes 注册完成')
}

/**
 * 路径解析函数：处理父路径和子路径的拼接
 */
function resolvePath(parent: string, child: string): string {
  return [parent.replace(/\/$/, ''), child.replace(/^\//, '')].filter(Boolean).join('/')
}

/**
 * 检测菜单中的重复路由（包括子路由）
 */
function checkDuplicateRoutes(routes: AppRouteRecord[], parentPath = ''): void {
  // 用于检测动态路由中的重复项
  const routeNameMap = new Map<string, string>() // 路由名称 -> 路径
  const componentPathMap = new Map<string, string>() // 组件路径 -> 路由信息

  const checkRoutes = (routes: AppRouteRecord[], parentPath = '') => {
    routes.forEach((route) => {
      // 处理路径拼接
      const currentPath = route.path || ''
      const fullPath = resolvePath(parentPath, currentPath)

      // 名称重复检测
      if (route.name) {
        if (routeNameMap.has(String(route.name))) {
          console.warn(`[路由警告] 名称重复: "${String(route.name)}"`)
        } else {
          routeNameMap.set(String(route.name), fullPath)
        }
      }

      // 组件路径重复检测
      if (route.component) {
        const componentPath = getComponentPathString(route.component)

        if (componentPath && componentPath !== RoutesAlias.Home) {
          const componentKey = `${parentPath}:${componentPath}`

          if (componentPathMap.has(componentKey)) {
            console.warn(`[路由警告] 路径重复: "${componentPath}"`)
          } else {
            componentPathMap.set(componentKey, fullPath)
          }
        }
      }

      // 递归处理子路由
      if (route.children?.length) {
        checkRoutes(route.children, fullPath)
      }
    })
  }

  checkRoutes(routes, parentPath)
}

/**
 * 获取组件路径的字符串表示
 */
function getComponentPathString(component: any): string {
  if (typeof component === 'string') {
    return component
  }

  // 对于其他别名路由，获取组件名称
  for (const key in RoutesAlias) {
    if (RoutesAlias[key as keyof typeof RoutesAlias] === component) {
      return `RoutesAlias.${key}`
    }
  }

  return ''
}

/**
 * 根据组件路径动态加载组件
 * @param componentPath 组件路径（不包含 ../../views 前缀和 .vue 后缀）
 * @param routeName 当前路由名称（用于错误提示）
 * @returns 组件加载函数
 */
function loadComponent(componentPath: string, routeName: string): () => Promise<any> {
  console.log(`尝试加载组件: ${routeName}, 路径: ${componentPath}`)

  // 如果路径为空，直接返回一个空的组件
  if (componentPath === '') {
    console.warn(`组件路径为空: ${routeName}`)
    return () =>
      Promise.resolve({
        render() {
          return h('div', {})
        }
      })
  }

  // 处理以/开头的组件路径
  if (componentPath.startsWith('/')) {
    componentPath = componentPath.substring(1)
    console.log(`调整后的组件路径: ${componentPath}`)
  }

  // 构建可能的路径
  const fullPath = `../../views/${componentPath}.vue`
  const fullPathWithIndex = `../../views/${componentPath}/index.vue`

  console.log(`尝试查找组件路径: ${fullPath}`)
  console.log(`或者: ${fullPathWithIndex}`)

  // 先尝试直接路径，再尝试添加/index的路径
  const module = modules[fullPath] || modules[fullPathWithIndex]

  if (!module) {
    console.error(
      `[路由错误] 未找到组件：${routeName}，尝试过的路径: ${fullPath} 和 ${fullPathWithIndex}`
    )
    // 列出所有可用的模块路径，方便调试
    console.log('可用的模块路径:', Object.keys(modules))

    return () =>
      Promise.resolve({
        render() {
          return h('div', `组件未找到: ${routeName}, 路径: ${componentPath}`)
        }
      })
  }

  console.log(`成功找到组件: ${routeName}`)
  return module
}

/**
 * 转换后的路由配置类型
 */
interface ConvertedRoute extends Omit<RouteRecordRaw, 'children'> {
  id?: number
  children?: ConvertedRoute[]
  component?: RouteRecordRaw['component'] | (() => Promise<any>)
}

/**
 * 转换路由组件配置
 */
function convertRouteComponent(
  route: AppRouteRecord,
  iframeRoutes: AppRouteRecord[],
  depth = 0
): ConvertedRoute {
  const { component, children, ...routeConfig } = route

  // 基础路由配置
  const converted: ConvertedRoute = {
    ...routeConfig,
    component: undefined
  }

  // 是否为一级菜单
  const isFirstLevel = depth === 0 && route.children?.length === 0

  if (route.meta.isIframe) {
    handleIframeRoute(converted, route, iframeRoutes)
  } else if (isFirstLevel) {
    handleLayoutRoute(converted, route, component as string)
  } else {
    handleNormalRoute(converted, component as string, String(route.name))
  }

  // 递归时增加深度
  if (children?.length) {
    converted.children = children.map((child) =>
      convertRouteComponent(child, iframeRoutes, depth + 1)
    )
  }

  return converted
}

/**
 * 处理 iframe 类型路由
 */
function handleIframeRoute(
  converted: ConvertedRoute,
  route: AppRouteRecord,
  iframeRoutes: AppRouteRecord[]
): void {
  converted.path = `/outside/iframe/${String(route.name)}`
  converted.component = () => import('@/views/outside/Iframe.vue')
  iframeRoutes.push(route)
}

/**
 * 处理一级菜单路由
 */
function handleLayoutRoute(
  converted: ConvertedRoute,
  route: AppRouteRecord,
  component: string | undefined
): void {
  converted.component = () => import('@/views/index/index.vue')
  converted.path = `/${(route.path?.split('/')[1] || '').trim()}`
  converted.name = ''
  route.meta.isFirstLevel = true

  converted.children = [
    {
      id: route.id,
      path: route.path,
      name: route.name,
      component: loadComponent(component as string, String(route.name)),
      meta: route.meta
    }
  ]
}

/**
 * 处理普通路由
 */
function handleNormalRoute(
  converted: ConvertedRoute,
  component: string | undefined,
  routeName: string
): void {
  if (component) {
    // 检查是否是 RoutesAlias 的值
    const aliasValues = Object.values(RoutesAlias)
    if (aliasValues.includes(component as any)) {
      // 如果是 RoutesAlias 的值，直接使用 loadComponent
      converted.component = loadComponent(component as string, routeName)
    } else {
      // 尝试作为 RoutesAlias 的键来查找
      const aliasComponent = RoutesAlias[
        component as keyof typeof RoutesAlias
      ] as unknown as RouteRecordRaw['component']
      converted.component = aliasComponent || loadComponent(component as string, routeName)
    }
  }
}
