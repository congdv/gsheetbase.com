import { Navigate, Route, Routes } from 'react-router-dom'
import { ConfigProvider, Layout, Spin } from 'antd'
import LoginPage from './pages/LoginPage'
import ProtectedRoute from './components/ProtectedRoute'
import DashboardLayout from './components/Layout/DashboardLayout'
import { useAuth } from './context/AuthContext'
import OAuthCallbackPage from './pages/OAuthCallbackPage'
import HomePage from './pages/home'
import BillingPage from './pages/billing'
import SheetDetailPage from './pages/sheet-detail'
import { ROUTES } from './constants/routes'

const { Content } = Layout

export default function App() {
  const { user, isLoading } = useAuth()

  if (isLoading) {
    return (
      <ConfigProvider>
        <Layout style={{ minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
          <Spin size="large" />
        </Layout>
      </ConfigProvider>
    )
  }

  return (
    <ConfigProvider>
      <Routes>
        <Route path={ROUTES.ROOT} element={!user ? <Navigate to={ROUTES.LOGIN} replace /> : <Navigate to={ROUTES.HOME} replace />} />
        <Route path={ROUTES.LOGIN} element={!user ? <LoginPage /> : <Navigate to={ROUTES.HOME} replace />} />
        <Route path={ROUTES.OAUTH_CALLBACK} element={<OAuthCallbackPage />} />

        <Route element={
          <ProtectedRoute>
            <DashboardLayout />
          </ProtectedRoute>
        }>
          <Route path={ROUTES.HOME} element={<HomePage />} />
          <Route path={ROUTES.BILLING} element={<BillingPage />} />
          <Route path={ROUTES.SHEET_DETAIL} element={<SheetDetailPage />} />
        </Route>
      </Routes>
    </ConfigProvider>
  )
}