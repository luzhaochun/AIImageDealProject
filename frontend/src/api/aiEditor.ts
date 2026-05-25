import { apiClient } from './client'
import type { AiStudioMode } from '@/constants/aiStudioModes'
import type { ApiResponse } from '@/types/api'

export interface AiStudioModeInfo {
  id: AiStudioMode
  label: string
  requires_mask: boolean
  mask_hint: string
  default_prompt: string
}

export interface AiConfig {
  mode: string
  has_key: boolean
  model: string
  studio_ready?: boolean
  studio_modes?: AiStudioModeInfo[]
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

/** Canvas 工作室：多模式 GPT Image（OpenAI Images Edits） */
export async function aiStudioEdit(
  imageBlob: Blob,
  maskBlob: Blob | null,
  prompt: string,
  mode: AiStudioMode,
) {
  const form = new FormData()
  form.append('image', imageBlob, 'image.png')
  if (maskBlob) {
    form.append('mask', maskBlob, 'mask.png')
  }
  form.append('prompt', prompt)
  form.append('mode', mode)
  const { data } = await apiClient.post<ApiResponse<AiInpaintResult>>('/admin/ai/studio-edit', form, {
    timeout: 180000,
  })
  return data.data!
}

/** @deprecated 使用 aiStudioEdit(..., 'erase') */
export async function aiStudioInpaint(imageBlob: Blob, maskBlob: Blob, prompt: string) {
  return aiStudioEdit(imageBlob, maskBlob, prompt, 'erase')
}
