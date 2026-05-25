import type { ImageStatus } from '@/types/image'

const labels: Record<ImageStatus, string> = {
  pending_storage: '待入库',
  pending_assign: '待分配',
  assigned: '已领取',
  discarded: '已废弃',
  in_progress: '作图中',
  pending_1st_review: '待一审',
  rejected: '已驳回',
  '1st_review_passed': '一审通过',
  pending_2nd_review: '待二审',
  '2nd_review_passed': '二审通过',
  completed: '已归档',
}

export function statusLabel(status: ImageStatus): string {
  return labels[status] ?? status
}
