export interface User {
  id: string
  email: string
  name: string
  avatar_url?: string
  is_verified: boolean
  created_at: string
  updated_at: string
}

export interface Skill {
  id: string
  user_id: string
  name: string
  category: string
  color: string
  icon: string
  target_hours: number
  is_archived: boolean
  created_at: string
  updated_at: string
}

export interface Session {
  id: string
  skill_id: string
  user_id: string
  start_time: string
  end_time?: string
  duration_minutes: number
  notes?: string
  session_type: string
  quality_rating?: number
  created_at: string
  updated_at: string
}

export interface Goal {
  id: string
  user_id: string
  skill_id?: string
  type: 'daily' | 'weekly' | 'milestone'
  target_hours: number
  deadline?: string
  is_completed: boolean
  created_at: string
  completed_at?: string
}

export interface Notification {
  id: string
  user_id: string
  type: string
  message: string
  is_read: boolean
  created_at: string
}

export interface AuthResponse {
  user: User
  token: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
  name: string
}

export interface CreateSkillRequest {
  name: string
  category: string
  color: string
  icon: string
  target_hours?: number
}

export interface CreateSessionRequest {
  skill_id: string
  start_time: string
  end_time?: string
  duration_minutes?: number
  notes?: string
  session_type?: string
  quality_rating?: number
}

export interface CreateGoalRequest {
  skill_id?: string
  type: 'daily' | 'weekly' | 'milestone'
  target_hours: number
  deadline?: string
}
