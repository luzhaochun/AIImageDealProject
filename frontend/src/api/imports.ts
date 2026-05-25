import { apiClient } from './client'
import type { ApiResponse } from '@/types/api'

export interface ImportBatch {
  id: string
  file_name: string
  status: string
  total_rows: number
  success_count: number
  failed_count: number
  error_log?: string
  created_at: string
}

export interface ImportBatchDetail extends ImportBatch {
  pending_storage?: number
  pending_assign?: number
  downloaded?: number
}

export async function uploadCSV(file: File) {
  const form = new FormData()
  form.append('file', file)
  const { data } = await apiClient.post<ApiResponse<ImportBatch>>('/imports', form, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
  return data.data!
}

export async function listImports(page = 1, pageSize = 20) {
  const { data } = await apiClient.get<ApiResponse<{ items: ImportBatch[]; total: number }>>('/imports', {
    params: { page, page_size: pageSize },
  })
  return data.data!
}

export async function getImport(id: string) {
  const { data } = await apiClient.get<ApiResponse<ImportBatchDetail>>('/imports', { params: { id } })
  return data.data!
}
