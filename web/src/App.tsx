import { Routes, Route, Navigate } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { api } from './api/client'
import { AppShell } from './components/layout/AppShell'
import { LoginPage } from './components/auth/LoginPage'
import { DashboardPage } from './components/dashboard/DashboardPage'
import { AccountsPage } from './components/accounts/AccountsPage'
import { UsagePage } from './components/usage/UsagePage'
import { SettingsPage } from './components/settings/SettingsPage'

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { data, isLoading } = useQuery({
    queryKey: ['auth'],
    queryFn: api.auth.check,
    retry: false,
    refetchInterval: false,
  })

  if (isLoading) return null
  if (!data?.authenticated) return <Navigate to="/login" replace />
  return <>{children}</>
}

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route
        path="/*"
        element={
          <ProtectedRoute>
            <AppShell>
              <Routes>
                <Route path="/" element={<DashboardPage />} />
                <Route path="/accounts" element={<AccountsPage />} />
                <Route path="/accounts/:id/usage" element={<UsagePage />} />
                <Route path="/settings" element={<SettingsPage />} />
              </Routes>
            </AppShell>
          </ProtectedRoute>
        }
      />
    </Routes>
  )
}
