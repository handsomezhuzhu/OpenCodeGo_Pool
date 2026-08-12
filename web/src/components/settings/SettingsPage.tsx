import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Save, Plus, X, Loader2 } from 'lucide-react'

import { api } from '@/api/client'
import { notifyError, notifySuccess } from '@/lib/errors'
import { PageHeader } from '@/components/common/PageHeader'
import { LoadingState } from '@/components/common/LoadingState'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Form, FormControl, FormField, FormItem, FormLabel } from '@/components/ui/form'

const formSchema = z.object({
  endpoint: z.string(),
  bearer_token: z.string(),
  provider_name: z.string(),
  base_url: z.string(),
})

type FormValues = z.infer<typeof formSchema>

export function SettingsPage() {
  const queryClient = useQueryClient()
  const { data: settings, isLoading } = useQuery({
    queryKey: ['cpa-settings'],
    queryFn: api.settings.getCPA,
    refetchInterval: false,
  })

  const [models, setModels] = useState<string[]>([])
  const [newModel, setNewModel] = useState('')

  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: { endpoint: '', bearer_token: '', provider_name: '', base_url: '' },
  })

  useEffect(() => {
    if (settings) {
      form.reset({
        endpoint: settings.endpoint,
        bearer_token: settings.bearer_token,
        provider_name: settings.provider_name,
        base_url: settings.base_url,
      })
      setModels(settings.models || [])
    }
  }, [settings, form])

  const saveMutation = useMutation({
    mutationFn: (values: FormValues) => api.settings.updateCPA({ ...values, models }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cpa-settings'] })
      notifySuccess('Settings saved')
    },
    onError: (err) => notifyError(err),
  })

  const addModel = () => {
    const m = newModel.trim()
    if (m && !models.includes(m)) {
      setModels([...models, m])
      setNewModel('')
    }
  }

  const removeModel = (name: string) => {
    setModels(models.filter((m) => m !== name))
  }

  if (isLoading) {
    return <LoadingState />
  }

  return (
    <div className="space-y-6">
      <PageHeader title="Settings" description="Configure CPA (ClipProxyAPI) integration" />

      <Card>
        <CardContent className="p-6">
          <Form {...form}>
            <form
              onSubmit={form.handleSubmit((values) => saveMutation.mutate(values))}
              className="space-y-4"
            >
              <FormField
                control={form.control}
                name="endpoint"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>CPA Endpoint</FormLabel>
                    <FormControl>
                      <Input className="font-mono" {...field} />
                    </FormControl>
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="bearer_token"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Bearer Token</FormLabel>
                    <FormControl>
                      <Input type="password" className="font-mono" {...field} />
                    </FormControl>
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="provider_name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Provider Name</FormLabel>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="base_url"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Base URL</FormLabel>
                    <FormControl>
                      <Input className="font-mono" {...field} />
                    </FormControl>
                  </FormItem>
                )}
              />

              <div className="space-y-2">
                <Label>Models</Label>
                <div className="flex flex-wrap gap-1.5">
                  {models.map((m) => (
                    <Badge key={m} variant="secondary" className="gap-1">
                      {m}
                      <button
                        type="button"
                        onClick={() => removeModel(m)}
                        className="text-muted-foreground hover:text-destructive"
                      >
                        <X className="h-3 w-3" />
                      </button>
                    </Badge>
                  ))}
                </div>
                <div className="flex gap-2">
                  <Input
                    value={newModel}
                    onChange={(e) => setNewModel(e.target.value)}
                    onKeyDown={(e) => {
                      if (e.key === 'Enter') {
                        e.preventDefault()
                        addModel()
                      }
                    }}
                    placeholder="Add model name..."
                    className="flex-1 font-mono"
                  />
                  <Button type="button" variant="secondary" size="icon" onClick={addModel}>
                    <Plus />
                  </Button>
                </div>
              </div>

              <Button type="submit" disabled={saveMutation.isPending} className="mt-2">
                {saveMutation.isPending ? <Loader2 className="animate-spin" /> : <Save />}
                {saveMutation.isPending ? 'Saving...' : 'Save Settings'}
              </Button>
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  )
}
