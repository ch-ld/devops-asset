// 注册全局组件
import type { App } from 'vue'
import TableBar from './Table/TableBar.vue'
import ArtTable from './Table/ArtTable.vue'
import FormInput from './Form/FormInput.vue'
import FormSelect from './Form/FormSelect.vue'

// 导入布局组件
import ArtLayouts from './core/layouts/art-layouts/index.vue'
import ArtHeaderBar from './core/layouts/art-header-bar/index.vue'
import ArtSidebarMenu from './core/layouts/art-menus/art-sidebar-menu/index.vue'
import ArtPageContent from './core/layouts/art-page-content/index.vue'
import ArtSettingsPanel from './core/layouts/art-settings-panel/index.vue'
import ArtGlobalSearch from './core/layouts/art-global-search/index.vue'
import ArtScreenLock from './core/layouts/art-screen-lock/index.vue'
import ArtWatermark from './core/layouts/art-watermark/index.vue'
import ArtLogo from './core/layouts/art-logo/index.vue'
import ArtWorkTab from './core/layouts/art-work-tab/index.vue'
import ArtBreadcrumb from './core/layouts/art-breadcrumb/index.vue'

export function registerGlobComp(app: App) {
  // 注册表单和表格组件
  app.component('table-bar', TableBar)
  app.component('art-table', ArtTable)
  app.component('form-input', FormInput)
  app.component('form-select', FormSelect)

  // 注册布局组件
  app.component('ArtLayouts', ArtLayouts)
  app.component('ArtHeaderBar', ArtHeaderBar)
  app.component('ArtSidebarMenu', ArtSidebarMenu)
  app.component('ArtPageContent', ArtPageContent)
  app.component('ArtSettingsPanel', ArtSettingsPanel)
  app.component('ArtGlobalSearch', ArtGlobalSearch)
  app.component('ArtScreenLock', ArtScreenLock)
  app.component('ArtWatermark', ArtWatermark)
  app.component('ArtLogo', ArtLogo)
  app.component('ArtWorkTab', ArtWorkTab)
  app.component('ArtBreadcrumb', ArtBreadcrumb)
}
