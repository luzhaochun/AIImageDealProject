import { apiClient } from './client'
import type { ApiResponse } from '@/types/api'
import type { ImageItem } from '@/types/image'

export interface ReviewRecord {
  id: string
  round: number
  result: string
  comment?: string
  reviewer_id: string
  created_at: string
}

export interface ReviewStats {
  queue_count: number
  today_reviewed: number
  pass_rate: number
}

export async function reviewQueue(round: 1 | 2, page = 1, pageSize = 20) {
  const { data } = await apiClient.get<ApiResponse<{ items: ImageItem[]; total: number }>>('/reviews/queue', {
    params: { round, page, page_size: pageSize },
  })
  return data.data!
}

export async function reviewStats(round: 1 | 2) {
  const { data } = await apiClient.get<ApiResponse<ReviewStats>>('/reviews/stats', { params: { round } })
  return data.data!
}

export async function secondReviewPool(page = 1, pageSize = 20) {
  const { data } = await apiClient.get<ApiResponse<{ items: ImageItem[]; total: number }>>('/reviews/second-pool', {
    params: { page, page_size: pageSize },
  })
  return data.data!
}

export async function submitReview(imageId: string, round: number, result: 'pass' | 'reject', comment: string) {
  const { data } = await apiClient.post<ApiResponse<ImageItem>>(`/images/${imageId}/reviews`, {
    round,
    result,
    comment,
  })
  return data.data!
}

export async function claimSecondReview() {
  const { data } = await apiClient.post<ApiResponse<ImageItem>>('/reviews/claim-second')
  return data.data!
}

export async function listImageReviews(imageId: string) {
  const { data } = await apiClient.get<ApiResponse<ReviewRecord[]>>(`/images/${imageId}/reviews`)
  return data.data!
}
