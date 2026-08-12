import { cn } from '@/lib/utils'

export type Tone = 'default' | 'success' | 'warning' | 'info' | 'danger' | 'muted'

const toneClass: Record<Tone, string> = {
  default: 'border-border bg-secondary text-secondary-foreground',
  success: 'border-success/30 bg-success/10 text-success',
  warning: 'border-warning/30 bg-warning/10 text-warning',
  info: 'border-info/30 bg-info/10 text-info',
  danger: 'border-destructive/30 bg-destructive/10 text-destructive',
  muted: 'border-border bg-muted text-muted-foreground',
}

const dotClass: Record<Tone, string> = {
  default: 'bg-secondary-foreground',
  success: 'bg-success',
  warning: 'bg-warning',
  info: 'bg-info',
  danger: 'bg-destructive',
  muted: 'bg-muted-foreground',
}

interface StatusBadgeProps {
  tone?: Tone
  label: string
  withDot?: boolean
  className?: string
}

export function StatusBadge({ tone = 'muted', label, withDot = true, className }: StatusBadgeProps) {
  return (
    <span
      className={cn(
        'inline-flex items-center gap-1.5 whitespace-nowrap rounded-md border px-2 py-0.5 text-xs font-medium',
        toneClass[tone],
        className,
      )}
    >
      {withDot && <span className={cn('h-1.5 w-1.5 rounded-full', dotClass[tone])} />}
      {label}
    </span>
  )
}
