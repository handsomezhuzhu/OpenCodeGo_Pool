import { Check, Copy } from 'lucide-react'

import { cn } from '@/lib/utils'
import { useClipboard } from '@/hooks/use-clipboard'
import { toast } from '@/components/ui/sonner'
import { Button, type ButtonProps } from '@/components/ui/button'

interface CopyButtonProps extends Omit<ButtonProps, 'onClick'> {
  value: string
  label?: string
  toastMessage?: string
}

export function CopyButton({
  value,
  label,
  toastMessage = 'Copied to clipboard',
  variant = 'outline',
  size = 'sm',
  className,
  ...props
}: CopyButtonProps) {
  const { copied, copy } = useClipboard()

  return (
    <Button
      type="button"
      variant={variant}
      size={size}
      className={cn(className)}
      onClick={async () => {
        const ok = await copy(value)
        if (ok) toast.success(toastMessage)
        else toast.error('Copy failed, please copy manually')
      }}
      {...props}
    >
      {copied ? <Check className="text-success" /> : <Copy />}
      {label ?? (copied ? 'Copied' : 'Copy')}
    </Button>
  )
}
