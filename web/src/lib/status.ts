import type { Account } from '@/api/types'
import type { Tone } from '@/components/common/StatusBadge'

interface StatusMeta {
  tone: Tone
  label: string
}

const ACCOUNT_STATUS: Record<Account['status'], StatusMeta> = {
  active: { tone: 'success', label: 'Active' },
  rate_limited: { tone: 'warning', label: 'Rate Limited' },
  error: { tone: 'danger', label: 'Error' },
  disabled: { tone: 'muted', label: 'Disabled' },
}

export function accountStatusMeta(status: Account['status']): StatusMeta {
  return ACCOUNT_STATUS[status] ?? { tone: 'muted', label: status }
}

const SYNC_STATUS: Record<string, StatusMeta> = {
  success: { tone: 'success', label: 'Synced' },
  error: { tone: 'danger', label: 'Error' },
}

export function syncStatusMeta(status: string | undefined): StatusMeta {
  if (!status) return { tone: 'muted', label: 'N/A' }
  return SYNC_STATUS[status] ?? { tone: 'muted', label: status }
}
