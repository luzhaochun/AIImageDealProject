import { apiClient } from './client'
import type { ApiResponse, HealthData } from '@/types/api'

export async function fetchHealth() {
  const { data } = await apiClient.get<ApiResponse<HealthData>>('/health')
  return data.data
}
