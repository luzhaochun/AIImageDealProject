export type ImageStatus =
  | 'pending_storage'
  | 'pending_assign'
  | 'assigned'
  | 'discarded'
  | 'in_progress'
  | 'pending_1st_review'
  | 'rejected'
  | '1st_review_passed'
  | 'pending_2nd_review'
  | '2nd_review_passed'
  | 'completed'

export interface ImageItem {
  id: string
  tenant_id?: string
  global_no?: string
  image_url?: string
  storage_path?: string
  thumb_path?: string
  category?: string
  status: ImageStatus
  assigned_to?: string
  discard_reason?: string
  created_at?: string
  updated_at?: string
}
