import { toast } from '@/components/ui/sonner'

/** Get a user-facing error message */
export function getErrorMessage(err: unknown, fallback = 'Something went wrong, please try again.'): string {
  if (err instanceof Error && err.message) return err.message
  return fallback
}

export function notifyError(err: unknown, fallback?: string) {
  toast.error(getErrorMessage(err, fallback))
}

export function notifySuccess(message: string, description?: string) {
  toast.success(message, description ? { description } : undefined)
}
