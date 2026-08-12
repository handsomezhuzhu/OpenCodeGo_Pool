import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react'
import type { ReactNode } from 'react'

import {
  applyTheme,
  getStoredTheme,
  nextMode,
  resolveDark,
  storeTheme,
  type ThemeMode,
} from '@/lib/theme'

interface ThemeContextValue {
  /** User-selected mode: light / dark / system */
  mode: ThemeMode
  /** Whether the resolved theme is currently dark (system preference already resolved) */
  isDark: boolean
  /** Set the mode directly */
  setMode: (mode: ThemeMode) => void
  /** Three-state cycle: light -> dark -> system -> light */
  cycle: () => void
}

const ThemeContext = createContext<ThemeContextValue | null>(null)

export function ThemeProvider({ children }: { children: ReactNode }) {
  const [mode, setModeState] = useState<ThemeMode>(getStoredTheme)
  const [isDark, setIsDark] = useState<boolean>(() => resolveDark(getStoredTheme()))

  const setMode = useCallback((next: ThemeMode) => {
    setModeState(next)
    storeTheme(next)
    applyTheme(next)
    setIsDark(resolveDark(next))
  }, [])

  const cycle = useCallback(() => {
    setMode(nextMode(getStoredTheme()))
  }, [setMode])

  // Re-apply to <html> whenever mode changes (runtime switch after the bootstrap script)
  useEffect(() => {
    applyTheme(mode)
    setIsDark(resolveDark(mode))
  }, [mode])

  // Only in system mode, follow live OS theme changes
  useEffect(() => {
    if (mode !== 'system') return
    const mql = window.matchMedia('(prefers-color-scheme: dark)')
    const onChange = () => {
      applyTheme('system')
      setIsDark(resolveDark('system'))
    }
    mql.addEventListener('change', onChange)
    return () => mql.removeEventListener('change', onChange)
  }, [mode])

  const value = useMemo<ThemeContextValue>(
    () => ({ mode, isDark, setMode, cycle }),
    [mode, isDark, setMode, cycle],
  )

  return <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>
}

export function useTheme(): ThemeContextValue {
  const ctx = useContext(ThemeContext)
  if (!ctx) throw new Error('useTheme must be used within a ThemeProvider')
  return ctx
}
