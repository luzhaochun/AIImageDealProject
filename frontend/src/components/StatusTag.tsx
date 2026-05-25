import { Tag } from 'antd'
import type { ImageStatus } from '@/types/image'
import { statusLabel } from '@/utils/status'

const colors: Partial<Record<ImageStatus, string>> = {
  pending_assign: 'default',
  assigned: 'processing',
  in_progress: 'processing',
  pending_1st_review: 'warning',
  rejected: 'error',
  '1st_review_passed': 'success',
  completed: 'success',
  discarded: 'default',
}

export function StatusTag({ status }: { status: ImageStatus }) {
  return <Tag color={colors[status]}>{statusLabel(status)}</Tag>
}
