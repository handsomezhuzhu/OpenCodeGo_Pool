import { useMemo } from 'react'
import {
  Bar,
  BarChart,
  CartesianGrid,
  Legend,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'

import type { UsageDailySummary } from '@/api/types'
import { Card, CardContent } from '@/components/ui/card'
import { useTheme } from '@/providers/ThemeProvider'

interface ChartDataPoint {
  date: string
  [key: string]: number | string
}

interface TooltipEntry {
  name: string
  value: number
  color: string
}

interface ChartTooltipProps {
  active?: boolean
  payload?: TooltipEntry[]
  label?: string
}

function formatXAxisDate(dateStr: string): string {
  const month = parseInt(dateStr.slice(5, 7), 10)
  const day = dateStr.slice(8, 10)
  return `${month}月 ${day}`
}

function formatDollar(usd: number): string {
  if (usd < 0.01) return `$${usd.toFixed(4)}`
  return `$${usd.toFixed(2)}`
}

function generateModelColors(keys: string[], isDark: boolean): Record<string, string> {
  const result: Record<string, string> = {}
  const n = Math.max(keys.length, 1)
  keys.forEach((key, i) => {
    const hue = Math.round((i * 360) / n)
    const s = isDark ? 52 : 65
    const l = isDark ? 58 : 70
    result[key] = `hsl(${hue}, ${s}%, ${l}%)`
  })
  return result
}

function buildChartData(summaries: UsageDailySummary[]): { data: ChartDataPoint[]; modelKeys: string[] } {
  if (summaries.length === 0) return { data: [], modelKeys: [] }

  const byDate = new Map<string, Map<string, number>>()
  const totals: Record<string, number> = {}

  for (const s of summaries) {
    const key = s.provider ? `${s.model} (${s.provider})` : s.model
    if (!byDate.has(s.date)) byDate.set(s.date, new Map())
    const dayMap = byDate.get(s.date)!
    dayMap.set(key, (dayMap.get(key) ?? 0) + s.cost)
    totals[key] = (totals[key] ?? 0) + s.cost
  }

  // highest cost at index 0 → bottom of stack; lowest cost at end → top
  const modelKeys = Object.keys(totals).sort((a, b) => totals[b] - totals[a])

  // fill calendar gaps using UTC noon to avoid DST edge cases
  const sortedDates = [...byDate.keys()].sort()
  const allDates: string[] = []
  let cursor = new Date(sortedDates[0] + 'T12:00:00Z')
  const last = new Date(sortedDates[sortedDates.length - 1] + 'T12:00:00Z')
  while (cursor <= last) {
    allDates.push(cursor.toISOString().slice(0, 10))
    cursor = new Date(cursor.getTime() + 86_400_000)
  }

  const data = allDates.map((date) => {
    const dayMap = byDate.get(date)
    const point: ChartDataPoint = { date }
    for (const key of modelKeys) {
      point[key] = (dayMap?.get(key) ?? 0) / 100_000_000
    }
    return point
  })

  return { data, modelKeys }
}

function ChartTooltip({ active, payload, label }: ChartTooltipProps) {
  if (!active || !payload?.length || !label) return null

  const sorted = [...payload].sort((a, b) => b.value - a.value)

  return (
    <div className="max-h-80 overflow-y-auto rounded-lg border bg-popover p-3 text-xs shadow-md">
      <p className="mb-2 font-medium text-foreground">{formatXAxisDate(label)}</p>
      <div className="space-y-1">
        {sorted.map((entry) => (
          <div key={entry.name} className="flex items-center gap-2">
            <span
              className="inline-block h-2.5 w-2.5 shrink-0 rounded-sm"
              style={{ background: entry.color }}
            />
            <span className="text-muted-foreground">{entry.name}</span>
            <span className="ml-auto pl-4 font-medium text-foreground">
              {formatDollar(entry.value)}
            </span>
          </div>
        ))}
      </div>
    </div>
  )
}

interface UsageChartProps {
  records: UsageDailySummary[]
}

export function UsageChart({ records }: UsageChartProps) {
  const { isDark } = useTheme()

  const { data, modelKeys } = useMemo(() => buildChartData(records), [records])
  const colors = useMemo(() => generateModelColors(modelKeys, isDark), [modelKeys, isDark])

  if (data.length === 0) return null

  const axisTickColor = isDark ? 'hsl(48, 7%, 55%)' : 'hsl(48, 4%, 50%)'
  const gridColor = isDark ? 'hsl(60, 2%, 24%)' : 'hsl(44, 20%, 85%)'
  const cursorColor = isDark ? 'rgba(255,255,255,0.04)' : 'rgba(0,0,0,0.04)'

  return (
    <Card>
      <CardContent className="pt-4">
        <ResponsiveContainer width="100%" height={280}>
          <BarChart data={data} margin={{ top: 4, right: 16, left: 0, bottom: 0 }} maxBarSize={36}>
            <CartesianGrid strokeDasharray="3 3" stroke={gridColor} vertical={false} />
            <XAxis
              dataKey="date"
              tickFormatter={formatXAxisDate}
              tick={{ fontSize: 11, fill: axisTickColor }}
              tickLine={false}
              axisLine={{ stroke: gridColor }}
              interval="preserveStartEnd"
              minTickGap={40}
            />
            <YAxis
              tickFormatter={(v: number) => `$${v.toFixed(2)}`}
              tick={{ fontSize: 11, fill: axisTickColor }}
              tickLine={false}
              axisLine={false}
              width={52}
            />
            <Tooltip
              content={<ChartTooltip />}
              cursor={{ fill: cursorColor }}
            />
            <Legend
              wrapperStyle={{ fontSize: 11, paddingTop: 12 }}
              formatter={(value: string) => (
                <span style={{ color: axisTickColor }}>{value}</span>
              )}
            />
            {modelKeys.map((key) => (
              <Bar
                key={key}
                dataKey={key}
                stackId="cost"
                fill={colors[key]}
                isAnimationActive={false}
              />
            ))}
          </BarChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  )
}
