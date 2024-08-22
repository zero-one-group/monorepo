import type ApiClient from '../client'
import type { LoginResponse, SignupResponse } from '../types/account'

export default class AuthService {
  constructor(private apiClient: ApiClient) {}

  login(username: string, password: string) {
    return this.apiClient._request<LoginResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    })
  }

  signup(username: string, password: string) {
    return this.apiClient._request<SignupResponse>('/auth/signup', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    })
  }
}
