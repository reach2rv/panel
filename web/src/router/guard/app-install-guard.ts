import type { Router } from 'vue-router'

import app from '@/api/panel/app'
import { $gettext } from '@/utils/gettext'

// 防止重复显示错误消息
let lastErrorMsg = ''
let lastErrorTime = 0
const ERROR_COOLDOWN = 2000

function showErrorMessage(message: string) {
  const now = Date.now()
  if (lastErrorMsg !== message || now - lastErrorTime > ERROR_COOLDOWN) {
    window.$message.error(message)
    lastErrorMsg = message
    lastErrorTime = now
  }
}

export function createAppInstallGuard(router: Router) {
  router.beforeEach(async (to) => {
    const slug = to.path.split('/').pop()
    if (to.path.startsWith('/apps/') && slug) {
      await useRequest(app.isInstalled(slug)).onSuccess(({ data }) => {
        if (!data) {
          showErrorMessage($gettext('App is not installed'))
          return router.push({ name: 'app-index' })
        }
      })
    }

    // 网站
    if (to.path.startsWith('/website')) {
      await useRequest(app.isInstalled('nginx,openresty,apache,openlitespeed,caddy')).onSuccess(
        ({ data }) => {
          if (!data) {
            showErrorMessage($gettext('Web server is not installed'))
            return router.push({ name: 'app-index' })
          }
        },
      )
    }

    // 容器
    if (to.path.startsWith('/container')) {
      await useRequest(app.isInstalled('docker,podman')).onSuccess(({ data }) => {
        if (!data) {
          showErrorMessage($gettext('Container engine is not installed'))
          return router.push({ name: 'app-index' })
        }
      })
    }
  })
}
