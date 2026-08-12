import { Monitor, Moon, Sun } from 'lucide-react'

import { useTheme } from '@/providers/ThemeProvider'
import type { ThemeMode } from '@/lib/theme'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'

const LABEL: Record<ThemeMode, string> = {
  light: 'Light',
  dark: 'Dark',
  system: 'System',
}

const ICON: Record<ThemeMode, typeof Sun> = {
  light: Sun,
  dark: Moon,
  system: Monitor,
}

/**
 * Theme toggle button: a single button cycling Light -> Dark -> System.
 * The icon follows the current mode; hovering shows the mode name.
 */
export function ThemeToggle() {
  const { mode, cycle } = useTheme()
  const Icon = ICON[mode]

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <Button
          variant="ghost"
          size="icon"
          onClick={cycle}
          aria-label={`Theme: ${LABEL[mode]} (click to switch)`}
        >
          <Icon />
        </Button>
      </TooltipTrigger>
      <TooltipContent>Theme: {LABEL[mode]}</TooltipContent>
    </Tooltip>
  )
}
