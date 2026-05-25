import { apiClient } from './client'
import type { ApiResponse } from '@/types/api'
import type { User } from '@/types/user'

export async function login(email: string, password: string) {
  const { data } = await apiClient.post<ApiResponse<{ token: string; expires_in: number; user: User }>>(
    '/auth/login',
    { email, password },
  )
  return data.data!
}

export async function fetchMe() {
  const { data } = await apiClient.get<ApiResponse<User>>('/auth/me')
  return data.data!
}
