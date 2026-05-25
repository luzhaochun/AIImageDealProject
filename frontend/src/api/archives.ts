import { apiClient } from './client'
import type { ApiResponse } from '@/types/api'
import type { ImageItem } from '@/types/image'

export interface ExportJob {
  id: string
  status: string
  filter_label: string
  file_name: string
  download_path?: string
  image_count: number
  created_at: string
}

export async function listArchives(params?: { category?: string; page?: number; page_size?: number }) {
  const { data } = await apiClient.get<ApiResponse<{ items: ImageItem[]; total: number }>>('/archives', { params })
  return data.data!
}

export async function completeArchive(id: string) {
  const { data } = await apiClient.post<ApiResponse<ImageItem>>(`/archives/${id}/complete`)
  return data.data!
}

export async function createExport(body: {
  category?: string
  status?: string
  date_from?: string
  date_to?: string
}) {
  const { data } = await apiClient.post<ApiResponse<ExportJob>>('/archives/export', body)
  return data.data!
}

export async function listExports() {
  const { data } = await apiClient.get<ApiResponse<ExportJob[]>>('/archives/exports')
  return data.data!
}
