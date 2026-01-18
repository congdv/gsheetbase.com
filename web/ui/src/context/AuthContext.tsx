import React, { createContext, useContext, useMemo, useState, useEffect } from 'react'
import api from '../lib/axios'

type User = {
  id: string
  email: string
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

  // Check for existing session on mount
  useEffect(() => {
    checkSession()
  }, [])

  const value = useMemo<AuthContextValue>(() => ({
    ...state,
    logout,
    checkSession
  }), [state])

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
