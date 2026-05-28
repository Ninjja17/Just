import axios from 'axios'
import type {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  Skill,
  CreateSkillRequest,
  Session,
  CreateSessionRequest,
  Goal,
  CreateGoalRequest,
  Notification,
} from '@/types'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

const api = axios.create({
  baseURL: `${API_URL}/api`,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor to add auth token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Response interceptor to handle errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// Auth API
export const authAPI = {
  register: (data: RegisterRequest) => 
    api.post<{ message: string; user: AuthResponse['user'] }>('/auth/register', data),
  
  login: (data: LoginRequest) => 
    api.post<AuthResponse>('/auth/login', data),
  
  verifyOTP: (email: string, code: string) => 
    api.post<AuthResponse>('/auth/verify-otp', { email, code }),
  
  logout: () => 
    api.post('/auth/logout'),
}

// Skills API
export const skillsAPI = {
  getAll: () => 
    api.get<Skill[]>('/skills'),
  
  getOne: (id: string) => 
    api.get<Skill>(`/skills/${id}`),
  
  create: (data: CreateSkillRequest) => 
    api.post<Skill>('/skills', data),
  
  update: (id: string, data: CreateSkillRequest) => 
    api.put<Skill>(`/skills/${id}`, data),
  
  delete: (id: string) => 
    api.delete(`/skills/${id}`),
  
  getStats: (id: string) => 
    api.get(`/skills/${id}/stats`),
}

// Sessions API
export const sessionsAPI = {
  getAll: (params?: { skill_id?: string; limit?: number; offset?: number }) => 
    api.get<Session[]>('/sessions', { params }),
  
  getOne: (id: string) => 
    api.get<Session>(`/sessions/${id}`),
  
  create: (data: CreateSessionRequest) => 
    api.post<Session>('/sessions', data),
  
  update: (id: string, data: Partial<CreateSessionRequest>) => 
    api.put<Session>(`/sessions/${id}`, data),
  
  delete: (id: string) => 
    api.delete(`/sessions/${id}`),
  
  startTimer: (skill_id: string) => 
    api.post<Session>('/sessions/start', { skill_id }),
  
  stopTimer: (session_id: string) => 
    api.post<Session>('/sessions/stop', { session_id }),
}

// Goals API
export const goalsAPI = {
  getAll: () => 
    api.get<Goal[]>('/goals'),
  
  create: (data: CreateGoalRequest) => 
    api.post<Goal>('/goals', data),
  
  update: (id: string, data: Partial<CreateGoalRequest>) => 
    api.put<Goal>(`/goals/${id}`, data),
  
  delete: (id: string) => 
    api.delete(`/goals/${id}`),
}

// Analytics API
export const analyticsAPI = {
  getOverview: () => 
    api.get('/analytics/overview'),
  
  getTrends: () => 
    api.get('/analytics/trends'),
  
  getPredictions: () => 
    api.get('/analytics/predictions'),
}

// Social API
export const socialAPI = {
  getProfile: (id: string) => 
    api.get(`/users/${id}/profile`),
  
  follow: (id: string) => 
    api.post(`/social/follow/${id}`),
  
  unfollow: (id: string) => 
    api.delete(`/social/follow/${id}`),
  
  getFollowers: () => 
    api.get('/social/followers'),
  
  getFollowing: () => 
    api.get('/social/following'),
  
  getLeaderboard: () => 
    api.get('/social/leaderboard'),
  
  getFeed: () => 
    api.get('/social/feed'),
}

// Notifications API
export const notificationsAPI = {
  getAll: (params?: { limit?: number; offset?: number }) => 
    api.get<Notification[]>('/notifications', { params }),
  
  markAsRead: (id: string) => 
    api.put(`/notifications/${id}/read`),
  
  delete: (id: string) => 
    api.delete(`/notifications/${id}`),
}

export default api
