import axios from 'axios'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

export const api = axios.create({
  baseURL: `${API_URL}/api/admin`,
  headers: { 'Content-Type': 'application/json' },
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('admin_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

api.interceptors.response.use(
  (r) => r,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem('admin_token')
      localStorage.removeItem('admin_user')
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }
    return Promise.reject(err)
  }
)

export interface AdminUser {
  id: string
  email: string
  name: string
  role: 'admin' | 'superadmin'
}

export interface UserRow {
  id: string
  email: string
  name: string
  is_verified: boolean
  created_at: string
  skill_count: number
  session_count: number
  total_minutes: number
}

export interface Stats {
  total_users: number
  verified_users: number
  total_skills: number
  total_sessions: number
  total_minutes: number
  new_users_7d: number
  active_users_7d: number
  sessions_per_day: { date: string; value: number }[]
  new_users_per_day: { date: string; value: number }[]
  top_skills: { name: string; count: number }[]
}

export const adminAPI = {
  login: (email: string, password: string) =>
    api.post<{ admin: AdminUser; token: string }>('/login', { email, password }),
  me: () => api.get<AdminUser>('/me'),
  stats: () => api.get<Stats>('/stats'),
  listUsers: (params: { search?: string; limit?: number; offset?: number }) =>
    api.get<{ users: UserRow[]; total: number; limit: number; offset: number }>('/users', { params }),
  getUser: (id: string) => api.get(`/users/${id}`),
  deleteUser: (id: string) => api.delete(`/users/${id}`),
  exportUsersCSV: () => `${API_URL}/api/admin/export/users.csv`,
  exportSessionsCSV: () => `${API_URL}/api/admin/export/sessions.csv`,
}

export function downloadCSV(path: string, filename: string) {
  const token = localStorage.getItem('admin_token')
  fetch(path, { headers: { Authorization: `Bearer ${token}` } })
    .then((r) => {
      if (!r.ok) throw new Error(`HTTP ${r.status}`)
      return r.blob()
    })
    .then((blob) => {
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = filename
      a.click()
      URL.revokeObjectURL(url)
    })
}
