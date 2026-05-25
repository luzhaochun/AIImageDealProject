import { apiClient } from './client'
import type { ApiResponse } from '@/types/api'
import type { ImageItem } from '@/types/image'

export interface CategoryCount {
  name: string
  available_count: number
}

export async function listCategories() {
  const { data } = await apiClient.get<ApiResponse<CategoryCount[]>>('/categories')
  return data.data!
}

export async function claimImage(category: string) {
  const { data } = await apiClient.post<ApiResponse<ImageItem>>('/images/claim', { category })
  return data.data!
}

export async function listImages(params: { library?: string; status?: string; page?: number; page_size?: number }) {
  const { data } = await apiClient.get<ApiResponse<{ items: ImageItem[]; total: number }>>('/images', { params })
  return data.data!
}

export async function discardImage(id: string, reason: string) {
  await apiClient.post(`/images/${id}/discard`, { reason })
}

export async function startEditing(id: string) {
  const { data } = await apiClient.post<ApiResponse<ImageItem>>(`/images/${id}/start-editing`)
  return data.data!
}

export async function saveVersion(id: string, layerData: unknown) {
  const { data } = await apiClient.put<ApiResponse<unknown>>(`/images/${id}/versions`, { layer_data: layerData })
  return data.data!
}

export async function submitReview(id: string) {
  const { data } = await apiClient.post<ApiResponse<ImageItem>>(`/images/${id}/submit`)
  return data.data!
}

export interface ImageDetail extends ImageItem {
  last_review_comment?: string
  assignee_name?: string
  version_no?: number
}

export async function getImage(id: string) {
  const { data } = await apiClient.get<ApiResponse<ImageDetail>>(`/images/${id}`)
  return data.data!
}

export async function listVersions(id: string) {
  const { data } = await apiClient.get<ApiResponse<{ id: string; version_no: number; created_at: string; is_current: number }[]>>(
    `/images/${id}/versions`,
  )
  return data.data!
}
