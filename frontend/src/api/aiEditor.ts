import { apiClient } from './client'
import type { ApiResponse } from '@/types/api'

export interface AiConfig {
  mode: string
  has_key: boolean
  model: string
  studio_ready?: boolean
}

export interface AiInpaintResult {
  result_url: string
  provider: string
  message?: string
}

export async function getAiConfig() {
  const { data } = await apiClient.get<ApiResponse<AiConfig>>('/admin/ai/config')
  return data.data!
}

export async function aiInpaint(imageBlob: Blob, maskBlob: Blob, prompt: string) {
  const form = new FormData()
  form.append('image', imageBlob, 'image.png')
  form.append('mask', maskBlob, 'mask.png')
  form.append('prompt', prompt)
  const { data } = await apiClient.post<ApiResponse<AiInpaintResult>>('/admin/ai/inpaint', form, {
    timeout: 120000,
  })
  return data.data!
}

/** Canvas 工作室：强制 GPT Image（OpenAI Images Edits） */
export async function aiStudioInpaint(imageBlob: Blob, maskBlob: Blob, prompt: string) {
  const form = new FormData()
  form.append('image', imageBlob, 'image.png')
  form.append('mask', maskBlob, 'mask.png')
  form.append('prompt', prompt)
  const { data } = await apiClient.post<ApiResponse<AiInpaintResult>>('/admin/ai/studio-inpaint', form, {
    timeout: 180000,
  })
  return data.data!
}
