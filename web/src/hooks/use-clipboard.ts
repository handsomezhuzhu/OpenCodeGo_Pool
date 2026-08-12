import { useCallback, useState } from 'react'

/** Copy to clipboard with a fallback path and a transient `copied` state */
export function useClipboard(timeout = 1500) {
  const [copied, setCopied] = useState(false)

  const copy = useCallback(
    async (text: string): Promise<boolean> => {
      const done = () => {
        setCopied(true)
        window.setTimeout(() => setCopied(false), timeout)
      }
      try {
        await navigator.clipboard.writeText(text)
        done()
        return true
      } catch {
        try {
          const ta = document.createElement('textarea')
          ta.value = text
          ta.style.position = 'fixed'
          ta.style.opacity = '0'
          document.body.appendChild(ta)
          ta.select()
          document.execCommand('copy')
          document.body.removeChild(ta)
          done()
          return true
        } catch {
          return false
        }
      }
    },
    [timeout],
  )

  return { copied, copy }
}
