export interface ApiResponse<T> {
  data?: T
  code?: string
  message?: string
  request_id?: string
}

export interface HealthData {
  status: string
  service: string
  mysql: boolean
  redis: string
  features: string
}
