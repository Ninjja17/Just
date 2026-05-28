import { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { User, AuthResponse } from '@/types'
import { authAPI } from '@/services/api'
import toast from 'react-hot-toast'

interface AuthContextType {
  user: User | null
  token: string | null
  loading: boolean
  login: (email: string, password: string) => Promise<void>
  register: (name: string, email: string, password: string) => Promise<void>
  verifyOTP: (email: string, code: string) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    // Check for stored token on mount
    const storedToken = localStorage.getItem('token')
    const storedUser = localStorage.getItem('user')

    if (storedToken && storedUser) {
      setToken(storedToken)
      setUser(JSON.parse(storedUser))
    }
    setLoading(false)
  }, [])

  const login = async (email: string, password: string) => {
    try {
      const response = await authAPI.login({ email, password })
      const { user, token } = response.data

      setUser(user)
      setToken(token)
      localStorage.setItem('token', token)
      localStorage.setItem('user', JSON.stringify(user))
      toast.success('Logged in successfully!')
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Login failed')
      throw error
    }
  }

  const register = async (name: string, email: string, password: string) => {
    try {
      await authAPI.register({ name, email, password })
      toast.success('Registration successful! Please check your email for OTP.')
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Registration failed')
      throw error
    }
  }

  const verifyOTP = async (email: string, code: string) => {
    try {
      const response = await authAPI.verifyOTP(email, code)
      const { user, token } = response.data

      setUser(user)
      setToken(token)
      localStorage.setItem('token', token)
      localStorage.setItem('user', JSON.stringify(user))
      toast.success('Email verified successfully!')
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Verification failed')
      throw error
    }
  }

  const logout = () => {
    authAPI.logout().catch(() => {})
    setUser(null)
    setToken(null)
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    toast.success('Logged out successfully')
  }

  return (
    <AuthContext.Provider value={{ user, token, loading, login, register, verifyOTP, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
