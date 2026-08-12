import type { ReactNode } from 'react'
import { NavLink } from 'react-router-dom'
import { LayoutDashboard, Users, Settings, RefreshCw, Menu, Sun, Moon, Monitor } from 'lucide-react'
import { useMutation, useQueryClient } from '@tanstack/react-query'

import { cn } from '@/lib/utils'
import { api } from '@/api/client'
import type { ThemeMode } from '@/lib/theme'
import { useTheme } from '@/providers/ThemeProvider'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetTrigger, SheetClose } from '@/components/ui/sheet'
import { SyncIndicator } from '@/components/sync/SyncIndicator'
import { ThemeToggle } from './ThemeToggle'

const navItems = [
  { to: '/', icon: LayoutDashboard, label: 'Dashboard' },
  { to: '/accounts', icon: Users, label: 'Accounts' },
  { to: '/settings', icon: Settings, label: 'Settings' },
]

type ThemeOption = { m: ThemeMode; Icon: typeof Sun; label: string }
const themeOptions: ThemeOption[] = [
  { m: 'light', Icon: Sun, label: 'Light' },
  { m: 'dark', Icon: Moon, label: 'Dark' },
  { m: 'system', Icon: Monitor, label: 'System' },
]

export function AppShell({ children }: { children: ReactNode }) {
  const queryClient = useQueryClient()
  const { mode, setMode } = useTheme()

  const refreshAll = useMutation({
    mutationFn: api.refreshAll,
    onSuccess: () => {
      setTimeout(() => queryClient.invalidateQueries({ queryKey: ['dashboard'] }), 3000)
    },
  })

  return (
    <div className="min-h-screen bg-background">
      <header className="sticky top-0 z-10 border-b border-border bg-card/95 backdrop-blur">
        <div className="mx-auto flex h-14 max-w-6xl items-center justify-between px-4">
          <div className="flex items-center gap-6">
            <h1 className="font-serif text-lg font-medium text-foreground">OpenCode Pool</h1>
            {/* Desktop nav — hidden on mobile */}
            <nav className="hidden items-center gap-1 md:flex">
              {navItems.map(({ to, icon: Icon, label }) => (
                <NavLink
                  key={to}
                  to={to}
                  end
                  className={({ isActive }) =>
                    cn(
                      'flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm transition-colors',
                      isActive
                        ? 'bg-accent font-medium text-foreground'
                        : 'text-muted-foreground hover:bg-accent/50 hover:text-foreground',
                    )
                  }
                >
                  <Icon size={15} />
                  {label}
                </NavLink>
              ))}
            </nav>
          </div>

          <div className="flex items-center gap-1">
            <SyncIndicator />
            <Tooltip>
              <TooltipTrigger asChild>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => refreshAll.mutate()}
                  disabled={refreshAll.isPending}
                >
                  <RefreshCw className={refreshAll.isPending ? 'animate-spin' : ''} />
                  <span className="hidden sm:inline">Refresh</span>
                </Button>
              </TooltipTrigger>
              <TooltipContent>Refresh all accounts</TooltipContent>
            </Tooltip>

            {/* Theme toggle — desktop only */}
            <span className="hidden md:inline-flex">
              <ThemeToggle />
            </span>

            {/* Mobile hamburger */}
            <Sheet>
              <SheetTrigger asChild>
                <Button variant="ghost" size="icon" className="md:hidden">
                  <Menu />
                </Button>
              </SheetTrigger>
              <SheetContent side="right" className="w-64">
                <SheetHeader>
                  <SheetTitle>OpenCode Pool</SheetTitle>
                </SheetHeader>

                {/* Theme switcher — 3 explicit buttons */}
                <div className="mt-5 flex gap-1 rounded-lg bg-muted p-1">
                  {themeOptions.map(({ m, Icon, label }) => (
                    <button
                      key={m}
                      onClick={() => setMode(m)}
                      className={cn(
                        'flex flex-1 items-center justify-center gap-1 rounded-md px-2 py-1.5 text-xs transition-colors',
                        mode === m
                          ? 'bg-background text-foreground shadow-sm'
                          : 'text-muted-foreground hover:text-foreground',
                      )}
                    >
                      <Icon size={13} />
                      {label}
                    </button>
                  ))}
                </div>

                {/* Nav items */}
                <nav className="mt-4 flex flex-col gap-1">
                  {navItems.map(({ to, icon: Icon, label }) => (
                    <SheetClose asChild key={to}>
                      <NavLink
                        to={to}
                        end
                        className={({ isActive }) =>
                          cn(
                            'flex items-center gap-2 rounded-lg px-3 py-2.5 text-sm transition-colors',
                            isActive
                              ? 'bg-accent font-medium text-foreground'
                              : 'text-muted-foreground hover:bg-accent/50 hover:text-foreground',
                          )
                        }
                      >
                        <Icon size={15} />
                        {label}
                      </NavLink>
                    </SheetClose>
                  ))}
                </nav>
              </SheetContent>
            </Sheet>
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-6xl px-4 py-6">{children}</main>
    </div>
  )
}
