import React, { createContext, useContext, useMemo, useState, useEffect } from 'react'
import api from '../lib/axios'

type User = {
  id: string
  email: string
  google_scopes?: string[]
  created_at: string
  updated_at: string
}

type AuthState = {
  user: User | null
  isLoading: boolean
}

type AuthContextValue = AuthState & {
  logout: () => Promise<void>
  checkSession: () => Promise<void>
  requestScopes: (scopes: string[], awaitCompletion?: boolean) => Promise<void>
  hasScope: (scope: string) => boolean
}

const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [state, setState] = useState<AuthState>({ user: null, isLoading: true })

  const checkSession = async () => {
    try {
      const response = await api.get('/auth/me')
      setState({
        user: response.data.user,
        isLoading: false,
      })
    } catch (error) {
      // No valid session, user needs to login
      setState({
        user: null,
        isLoading: false
      })
    }
  }

  const logout = async () => {
    try {
      await api.post('/auth/logout')
    } catch { }
    setState({ user: null, isLoading: false })
  }

  const hasScope = (scope: string): boolean => {
    return state.user?.google_scopes?.includes(scope) ?? false
  }

  const requestScopes = async (scopes: string[], awaitCompletion = true) => {
    try {
      const response = await api.post('/auth/google/request-scopes', { scopes })
      const authUrl = response.data.auth_url

      // Open in popup window instead of full redirect
      const width = 600
      const height = 700
      const left = window.screenX + (window.outerWidth - width) / 2
      const top = window.screenY + (window.outerHeight - height) / 2

      const popup = window.open(
        authUrl,
        'Google Authorization',
        `width=${width},height=${height},left=${left},top=${top},toolbar=0,menubar=0,location=0`
      )

      // If caller does not want to wait for completion, return immediately after opening popup
      if (!awaitCompletion) return

      // Return a promise that resolves only after popup closes and checkSession completes
      return await new Promise<void>((resolve, reject) => {
        const start = Date.now()
        const maxWait = 1000 * 60 * 3 // 3 minutes timeout

        const pollTimer = setInterval(async () => {
          try {
            if (!popup) {
              clearInterval(pollTimer)
              reject(new Error('Popup blocked or failed to open'))
              return
            }

            if (popup.closed) {
              clearInterval(pollTimer)
              try {
                await checkSession()
                resolve()
              } catch (err) {
                reject(err)
              }
              return
            }

            if (Date.now() - start > maxWait) {
              clearInterval(pollTimer)
              try { popup.close() } catch { }
              reject(new Error('Timed out waiting for OAuth consent'))
            }
          } catch (err) {
            clearInterval(pollTimer)
            reject(err as Error)
          }
        }, 500)
      })
    } catch (error) {
      console.error('Failed to request scopes:', error)
      throw error
    }
  }

  // Check for existing session on mount
  useEffect(() => {
    checkSession()
  }, [])

  const value = useMemo<AuthContextValue>(() => ({
    ...state,
    logout,
    checkSession,
    requestScopes,
    hasScope
  }), [state])

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
