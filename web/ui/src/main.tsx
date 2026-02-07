import React from 'react'
import ReactDOM from 'react-dom/client'
import { QueryClientProvider } from '@tanstack/react-query'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import { queryClient } from './lib/queryClient'
import { configureAxios } from './lib/axios'
import { AuthProvider } from './context/AuthContext'
import { ConfigProvider } from './context/ConfigContext'
import { fetchConfig } from './services/configService'
import App from './App'
import './index.css'

const router = createBrowserRouter([
  { path: '/*', element: <App /> },
])

function ConfigErrorPage({ error }: { error: string }) {
  return (
    <div style={{
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      justifyContent: 'center',
      minHeight: '100vh',
      backgroundColor: '#f8f9fa',
      fontFamily: 'system-ui, -apple-system, sans-serif',
      padding: '20px',
      textAlign: 'center',
    }}>
      <h1 style={{ color: '#d32f2f', marginBottom: '16px' }}>
        Configuration Error
      </h1>
      <p style={{ color: '#666', maxWidth: '500px', lineHeight: '1.6' }}>
        {error}
      </p>
      <p style={{ color: '#999', marginTop: '24px', fontSize: '14px' }}>
        Please ensure the API server is running and accessible.
      </p>
    </div>
  )
}

async function initializeApp() {
  try {
    const config = await fetchConfig()

    // Configure axios with runtime config before mounting app
    configureAxios(config.apiBaseUrl)

    const root = ReactDOM.createRoot(document.getElementById('root')!)
    root.render(
      <React.StrictMode>
        <ConfigProvider config={config}>
          <QueryClientProvider client={queryClient}>
            <AuthProvider>
              <RouterProvider router={router} />
            </AuthProvider>
          </QueryClientProvider>
        </ConfigProvider>
      </React.StrictMode>
    )
  } catch (error) {
    const errorMessage =
      error instanceof Error ? error.message : 'Unknown error occurred'
    const root = ReactDOM.createRoot(document.getElementById('root')!)
    root.render(<ConfigErrorPage error={errorMessage} />)
  }
}

initializeApp()
