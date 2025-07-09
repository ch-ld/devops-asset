import { AppRouteRecord } from '@/types/router'
import { router } from '@/router'
import { useMenuStore } from '@/store/modules/menu'

// 打开外部链接
export const openExternalLink = (link: string) => {
  window.open(link, '_blank')
}

/**
 * 菜单跳转
 * @param item 菜单项
 * @param jumpToFirst 是否跳转到第一个子菜单
 * @returns
 */
export const handleMenuJump = async (item: AppRouteRecord, jumpToFirst: boolean = false) => {
  // 处理外部链接
  const { link, isIframe } = item.meta
  if (link && !isIframe) {
    return openExternalLink(link)
  }

  // 跳转目标
  const targetName = item.name;
  const targetPath = item.path;

  // 判断路由是否注册
  const isRegistered = typeof targetName === 'string' ? router.hasRoute(targetName) : !!router.getRoutes().find(r => r.path === targetPath);
  if (!isRegistered) {
    console.warn(`路由 ${targetPath} 未注册，等待注册完成...`);
    await waitForRouteRegistration(typeof targetName === 'string' ? targetName : targetPath);
  }

  // 跳转
  if (typeof targetName === 'string' && router.hasRoute(targetName)) {
    return router.push({ name: targetName });
  }
  return router.push(targetPath);
}

/**
 * 等待路由注册完成
 * @param nameOrPath 路由 name 或 path
 * @param maxWaitTime 最大等待时间（毫秒）
 */
const waitForRouteRegistration = (nameOrPath: string, maxWaitTime: number = 3000): Promise<void> => {
  return new Promise((resolve, reject) => {
    const startTime = Date.now();
    const checkRoute = () => {
      if (
        (typeof nameOrPath === 'string' && router.hasRoute(nameOrPath)) ||
        !!router.getRoutes().find(r => r.path === nameOrPath)
      ) {
        resolve();
        return;
      }
      if (Date.now() - startTime > maxWaitTime) {
        reject(new Error(`路由 ${nameOrPath} 注册超时`));
        return;
      }
      setTimeout(checkRoute, 100);
    };
    checkRoute();
  });
}
