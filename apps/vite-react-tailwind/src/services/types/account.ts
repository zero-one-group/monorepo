import type { BaseApiResponse } from './base'

export interface User {
  id: number // Unique identifier for the user
  email: string // User's email address
  first_name: string // User's first name
  last_name?: string // User's last name
  avatar_url?: string // URL to user's avatar
  email_confirmed_at?: number // Timestamp when email was confirmed
  last_seen_at?: number // Timestamp of last user activity
  banned_until?: number // Timestamp when the user is banned until
  created_at: number // Timestamp when the user was created
  updated_at?: number // Timestamp when the user was last updated
}

export interface LoginResponse extends BaseApiResponse {
  error?: {
    hint?: string
    message: string
  }
  data?: {
    accessToken: string
    refreshToken: string
    role: 'admin' | 'user'
    user?: User
  }
}

export interface SignupResponse extends BaseApiResponse {
  data: {
    password: string
  }
}
