import { Alert } from 'antd'
import { PageHeader } from '@/components/PageHeader'

export function Placeholder({
  title,
  description = '功能开发中，当前为项目骨架占位页。',
}: {
  title: string
  description?: string
}) {
  return (
    <>
      <PageHeader title={title} />
      <Alert type="info" message={description} showIcon />
    </>
  )
}
