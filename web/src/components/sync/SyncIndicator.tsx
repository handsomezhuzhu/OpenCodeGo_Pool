import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Cloud, CloudOff, Loader2 } from 'lucide-react'

import { api } from '@/api/client'
import { syncStatusMeta } from '@/lib/status'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { StatusBadge } from '@/components/common/StatusBadge'
import { DropdownMenu, DropdownMenuContent, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'

export function SyncIndicator() {
  const queryClient = useQueryClient()

  const { data: syncLog } = useQuery({
    queryKey: ['sync-status'],
    queryFn: api.sync.status,
    refetchInterval: 30_000,
  })

  const syncMutation = useMutation({
    mutationFn: api.sync.trigger,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['sync-status'] })
    },
  })

  const isOk = syncLog?.status === 'success'
  const meta = syncStatusMeta(syncLog?.status)

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <button
          className={cn(
            'flex items-center gap-1.5 rounded-lg px-2.5 py-1 text-xs transition-colors',
            isOk
              ? 'bg-success/10 text-success'
              : syncLog
                ? 'bg-destructive/10 text-destructive'
                : 'bg-muted text-muted-foreground',
          )}
        >
          {isOk ? <Cloud size={13} /> : <CloudOff size={13} />}
          CPA {syncLog ? (isOk ? 'Synced' : 'Error') : 'N/A'}
        </button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-72 p-4">
        <p className="mb-2 text-sm font-medium text-foreground">CPA Sync Status</p>
        {syncLog ? (
          <>
            <p className="mb-1.5 flex items-center gap-1.5 text-xs text-muted-foreground">
              Status: <StatusBadge tone={meta.tone} label={meta.label} />
            </p>
            <p className="mb-1 text-xs text-muted-foreground">Keys: {syncLog.key_count}</p>
            <p className="mb-3 text-xs text-muted-foreground/70">
              {new Date(syncLog.synced_at).toLocaleString()}
            </p>
            {syncLog.message && (
              <p className="mb-3 max-h-20 overflow-auto break-all rounded bg-muted p-2 text-xs text-muted-foreground">
                {syncLog.message}
              </p>
            )}
          </>
        ) : (
          <p className="mb-3 text-xs text-muted-foreground">No sync performed yet</p>
        )}
        <Button
          size="sm"
          className="w-full"
          onClick={() => syncMutation.mutate()}
          disabled={syncMutation.isPending}
        >
          {syncMutation.isPending && <Loader2 className="animate-spin" />}
          {syncMutation.isPending ? 'Syncing...' : 'Sync Now'}
        </Button>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
