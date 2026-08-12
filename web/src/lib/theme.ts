/*
 * Theme (dark mode) logic layer.
 * Three states: light / dark / system (follows OS preference, default).
 * Color variables live in index.css's :root and .dark blocks; this file only
 * toggles the .dark class on <html>, persists the choice, and resolves system preference.
 *
 * Note: the storage key and resolution logic here must stay in sync with the
 * inline bootstrap script in index.html, otherwise a reload will flash the
 * wrong theme before React mounts (FOUC).
 */

export type ThemeMode = 'light' | 'dark' | 'system'

export const THEME_STORAGE_KEY = 'ocp-theme'

const DARK_QUERY = '(prefers-color-scheme: dark)'

/** Read the persisted theme mode; missing or invalid values fall back to system */
export function getStoredTheme(): ThemeMode {
  try {
    const v = localStorage.getItem(THEME_STORAGE_KEY)
    if (v === 'light' || v === 'dark' || v === 'system') return v
  } catch {
    // localStorage unavailable (private mode, etc.) — fall back silently
  }
  return 'system'
}

/** Persist the theme mode */
export function storeTheme(mode: ThemeMode): void {
  try {
    localStorage.setItem(THEME_STORAGE_KEY, mode)
  } catch {
    // ignore write failures
  }
}

/** Whether the OS currently prefers dark */
export function systemPrefersDark(): boolean {
  return typeof window !== 'undefined' && window.matchMedia(DARK_QUERY).matches
}

/** Whether the given mode should resolve to dark */
export function resolveDark(mode: ThemeMode): boolean {
  return mode === 'dark' || (mode === 'system' && systemPrefersDark())
}

/** Apply the mode to <html>: toggle .dark class and sync color-scheme (native controls/scrollbars) */
export function applyTheme(mode: ThemeMode): void {
  const dark = resolveDark(mode)
  const root = document.documentElement
  root.classList.toggle('dark', dark)
  root.style.colorScheme = dark ? 'dark' : 'light'
}

/** Three-state cycle: light -> dark -> system -> light */
export function nextMode(mode: ThemeMode): ThemeMode {
  return mode === 'light' ? 'dark' : mode === 'dark' ? 'system' : 'light'
}
