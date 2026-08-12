import { useNavigate } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { Plus } from 'lucide-react'

import { api } from '@/api/client'
import { PageHeader } from '@/components/common/PageHeader'
import { EmptyState } from '@/components/common/EmptyState'
import { ErrorState } from '@/components/common/ErrorState'
import { AccountCardSkeleton } from '@/components/common/LoadingState'
import { Button } from '@/components/ui/button'
import { AccountCard } from './AccountCard'

export function DashboardPage() {
  const navigate = useNavigate()
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['dashboard'],
    queryFn: api.dashboard,
  })

  const accounts = data?.accounts ?? []

  return (
    <div className="space-y-6">
      <PageHeader
        title="Dashboard"
        description={isLoading ? undefined : `${accounts.length} account(s) in pool`}
        actions={
          <Button onClick={() => navigate('/accounts')}>
            <Plus />
            Add Account
          </Button>
        }
      />

      {isLoading ? (
        <AccountCardSkeleton />
      ) : error ? (
        <ErrorState error={error} onRetry={() => refetch()} />
      ) : accounts.length === 0 ? (
        <EmptyState
          title="No accounts yet"
          description="Add your first OpenCode Go account to get started."
        />
      ) : (
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
          {accounts.map((acc, i) => (
            <div key={acc.id} className="animate-fade-in" style={{ animationDelay: `${i * 50}ms` }}>
              <AccountCard account={acc} />
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
