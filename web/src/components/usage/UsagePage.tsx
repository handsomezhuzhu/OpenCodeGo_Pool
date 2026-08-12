import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { ArrowLeft, Loader2 } from 'lucide-react'

import { api } from '@/api/client'
import type { UsageRecord } from '@/api/types'
import { EmptyState } from '@/components/common/EmptyState'
import { TableSkeleton } from '@/components/common/LoadingState'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { UsageChart } from './UsageChart'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

const PAGE_SIZE = 200

function formatCost(microUnits: number): string {
  const dollars = microUnits / 100_000_000
  if (dollars < 0.01) return `$${dollars.toFixed(4)}`
  return `$${dollars.toFixed(2)}`
}

function formatTokens(n: number): string {
  if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(1)}M`
  if (n >= 1_000) return `${(n / 1_000).toFixed(1)}K`
  return n.toString()
}

export function UsagePage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()

  const [offset, setOffset] = useState(0)
  const [allRecords, setAllRecords] = useState<UsageRecord[]>([])

  const { data: account } = useQuery({
    queryKey: ['account', id],
    queryFn: () => api.accounts.get(id!),
    enabled: !!id,
  })

  // daily summary for chart (all historical data, no record-count cap)
  const { data: dailySummaries = [] } = useQuery({
    queryKey: ['usage-daily', id],
    queryFn: () => api.accounts.usageDailyHistory(id!),
    enabled: !!id,
    select: (data) => data ?? [],
  })

  // paginated raw records for table
  const { data: pageRecords, isLoading, isFetching } = useQuery({
    queryKey: ['usage', id, offset],
    queryFn: () => api.accounts.usage(id!, PAGE_SIZE, offset),
    enabled: !!id,
    select: (data) => data ?? [],
  })

  useEffect(() => {
    if (!pageRecords) return
    if (offset === 0) {
      setAllRecords(pageRecords)
    } else {
      setAllRecords((prev) => [...prev, ...pageRecords])
    }
  }, [pageRecords, offset])

  // reset when account changes
  useEffect(() => {
    setOffset(0)
    setAllRecords([])
  }, [id])

  const hasMore = (pageRecords?.length ?? 0) === PAGE_SIZE

  const totalCost = dailySummaries.reduce((sum, r) => sum + r.cost, 0)

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <Tooltip>
          <TooltipTrigger asChild>
            <Button variant="ghost" size="icon" onClick={() => navigate(-1)}>
              <ArrowLeft />
            </Button>
          </TooltipTrigger>
          <TooltipContent>Back</TooltipContent>
        </Tooltip>
        <div>
          <h2 className="font-serif text-2xl font-medium text-foreground">Usage Details</h2>
          <p className="mt-0.5 text-sm text-muted-foreground">
            {account?.email || 'Loading...'} &middot; Total: {formatCost(totalCost)}
          </p>
        </div>
      </div>

      <UsageChart records={dailySummaries} />

      {isLoading && offset === 0 ? (
        <TableSkeleton />
      ) : allRecords.length === 0 ? (
        <EmptyState title="No usage records found." />
      ) : (
        <>
          <Card>
            <CardContent className="overflow-x-auto p-0">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Time</TableHead>
                    <TableHead>Model</TableHead>
                    <TableHead>Provider</TableHead>
                    <TableHead className="text-right">Input</TableHead>
                    <TableHead className="text-right">Output</TableHead>
                    <TableHead className="text-right">Reasoning</TableHead>
                    <TableHead className="text-right">Cache</TableHead>
                    <TableHead className="text-right">Cost</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {allRecords.map((r) => (
                    <TableRow key={r.id}>
                      <TableCell className="whitespace-nowrap text-xs text-muted-foreground">
                        {new Date(r.timeCreated).toLocaleString()}
                      </TableCell>
                      <TableCell className="font-mono text-xs text-foreground">{r.model}</TableCell>
                      <TableCell className="text-xs text-muted-foreground">{r.provider}</TableCell>
                      <TableCell className="text-right text-xs text-muted-foreground">
                        {formatTokens(r.inputTokens)}
                      </TableCell>
                      <TableCell className="text-right text-xs text-muted-foreground">
                        {formatTokens(r.outputTokens)}
                      </TableCell>
                      <TableCell className="text-right text-xs text-muted-foreground">
                        {r.reasoningTokens != null ? formatTokens(r.reasoningTokens) : '-'}
                      </TableCell>
                      <TableCell className="text-right text-xs text-muted-foreground">
                        {formatTokens(r.cacheReadTokens)}
                      </TableCell>
                      <TableCell className="text-right text-xs font-medium text-foreground">
                        {formatCost(r.cost)}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>

          {hasMore && (
            <div className="flex justify-center">
              <Button
                variant="outline"
                onClick={() => setOffset((prev) => prev + PAGE_SIZE)}
                disabled={isFetching}
              >
                {isFetching ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null}
                Load More
              </Button>
            </div>
          )}
        </>
      )}
    </div>
  )
}
