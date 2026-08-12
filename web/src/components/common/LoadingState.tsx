import { Loader2 } from 'lucide-react'

import { cn } from '@/lib/utils'
import { Skeleton } from '@/components/ui/skeleton'

export function LoadingState({
  text = 'Loading…',
  className,
}: {
  text?: string
  className?: string
}) {
  return (
    <div
      className={cn(
        'flex flex-col items-center justify-center gap-3 py-16 text-muted-foreground',
        className,
      )}
    >
      <Loader2 className="h-6 w-6 animate-spin" />
      <span className="text-sm">{text}</span>
    </div>
  )
}

/** Skeleton grid shaped like AccountCard: status row + email line + three quota bars */
export function AccountCardSkeleton({ count = 6 }: { count?: number }) {
  return (
    <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
      {Array.from({ length: count }).map((_, i) => (
        <div key={i} className="space-y-4 rounded-xl border bg-card p-4">
          <div className="flex items-center justify-between">
            <Skeleton className="h-4 w-2/5" />
            <Skeleton className="h-5 w-16 rounded-full" />
          </div>
          <Skeleton className="h-3 w-3/5" />
          <div className="space-y-2">
            <Skeleton className="h-2.5 w-full rounded-full" />
            <Skeleton className="h-2.5 w-full rounded-full" />
            <Skeleton className="h-2.5 w-full rounded-full" />
          </div>
        </div>
      ))}
    </div>
  )
}

export function TableSkeleton({ rows = 6 }: { rows?: number }) {
  return (
    <div className="space-y-2">
      {Array.from({ length: rows }).map((_, i) => (
        <Skeleton key={i} className="h-12 w-full" />
      ))}
    </div>
  )
}
