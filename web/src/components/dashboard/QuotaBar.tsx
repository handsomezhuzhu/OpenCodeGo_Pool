interface QuotaBarProps {
  label: string
  percent: number
  resetSec: number
  limit?: number | null
}

function formatReset(sec: number): string {
  if (sec <= 0) return ''
  const days = Math.floor(sec / 86400)
  const hours = Math.floor((sec % 86400) / 3600)
  const mins = Math.floor((sec % 3600) / 60)
  if (days > 0) return `${days}d ${hours}h`
  if (hours > 0) return `${hours}h ${mins}m`
  return `${mins}m`
}

function barColor(percent: number): string {
  if (percent >= 90) return 'bg-quota-red'
  if (percent >= 50) return 'bg-quota-amber'
  return 'bg-quota-green'
}

function textColor(percent: number): string {
  if (percent >= 90) return 'text-quota-red'
  if (percent >= 50) return 'text-quota-amber'
  return 'text-quota-green'
}

export function QuotaBar({ label, percent, resetSec, limit }: QuotaBarProps) {
  return (
    <div className="space-y-1">
      <div className="flex items-center justify-between text-xs">
        <span className="text-muted-foreground">
          {label}
          {limit != null && <span className="ml-1 text-[10px] opacity-60">≤{limit}%</span>}
        </span>
        <span className={`font-medium ${textColor(percent)}`}>{percent}%</span>
      </div>
      <div className="relative h-1.5 overflow-hidden rounded-full bg-muted">
        <div
          className={`h-full rounded-full transition-all duration-500 ${barColor(percent)}`}
          style={{ width: `${Math.min(percent, 100)}%` }}
        />
        {limit != null && (
          <div
            className="absolute top-0 h-full w-px bg-foreground/40"
            style={{ left: `${Math.min(limit, 100)}%` }}
          />
        )}
      </div>
      {resetSec > 0 && (
        <p className="text-[10px] text-muted-foreground">Reset in {formatReset(resetSec)}</p>
      )}
    </div>
  )
}
