import { useUserStore } from '@/stores'
import { $gettext } from '@/utils/gettext'

export function resolveResError(code: number | string | undefined, msg = ''): string {
  switch (code) {
    case 400:
    case 422:
      msg = msg || $gettext('Request parameter error')
      break
    case 401:
      msg = msg || $gettext('Login has expired')
      useUserStore().logout()
      break
    case 418:
      msg = msg || $gettext('Login has expired')
      useUserStore().refresh()
      break
    case 403:
      msg = msg || $gettext('Permission denied')
      break
    case 404:
      msg = msg || $gettext('Resource or API does not exist')
      break
    case 500:
      msg = msg || $gettext('Server exception')
      break
    default:
      msg = msg || $gettext('【{code}】: Unknown exception!', { code: String(code) })
      break
  }
  return msg
}
