import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Button, Input, Modal, Space, Table, message } from 'antd'
import { Link } from 'react-router-dom'
import { discardImage, listImages, startEditing, submitReview } from '@/api/images'
import { ImageThumb } from '@/components/ImageThumb'
import { PageHeader } from '@/components/PageHeader'
import { StatusTag } from '@/components/StatusTag'
import type { ImageItem, ImageStatus } from '@/types/image'

const filters: { key: string; label: string }[] = [
  { key: '', label: '全部' },
  { key: 'in_progress', label: '进行中' },
  { key: 'pending_1st_review', label: '待审核' },
  { key: 'rejected', label: '已驳回' },
]

export function TaskListPage() {
  const [status, setStatus] = useState('')
  const { data, refetch, isLoading } = useQuery({
    queryKey: ['tasks', status],
    queryFn: () => listImages({ library: 'tasks', status: status || undefined }),
  })

  const onDiscard = (id: string) => {
    let reason = '质量不符合要求'
    Modal.confirm({
      title: '换图',
      content: (
        <Input.TextArea
          defaultValue={reason}
          rows={3}
          onChange={(e) => {
            reason = e.target.value
          }}
        />
      ),
      onOk: async () => {
        await discardImage(id, reason)
        message.success('已换图')
        refetch()
      },
    })
  }

  const columns = [
    {
      title: '缩略图',
      render: (_: unknown, r: ImageItem) => <ImageThumb path={r.thumb_path || r.storage_path} />,
    },
    { title: '全局编号', dataIndex: 'global_no' },
    { title: '类目', dataIndex: 'category' },
    { title: '状态', dataIndex: 'status', render: (s: ImageStatus) => <StatusTag status={s} /> },
    {
      title: '更新时间',
      dataIndex: 'updated_at',
      render: (t: string) => (t ? new Date(t).toLocaleString() : '—'),
    },
    {
      title: '操作',
      render: (_: unknown, r: ImageItem) => (
        <Space wrap>
          {(r.status === 'assigned' || r.status === 'in_progress' || r.status === 'rejected') && (
            <Link to={`/workspace/editor/${r.id}`}>
              {r.status === 'rejected' ? '修改' : r.status === 'in_progress' ? '继续编辑' : '开始编辑'}
            </Link>
          )}
          {r.status === 'assigned' && (
            <Button
              size="small"
              onClick={async () => {
                await startEditing(r.id)
                refetch()
              }}
            >
              开始编辑
            </Button>
          )}
          {r.status === 'in_progress' && (
            <Button
              size="small"
              type="primary"
              onClick={async () => {
                await submitReview(r.id)
                message.success('已提交审核')
                refetch()
              }}
            >
              提交审核
            </Button>
          )}
          {r.status === 'assigned' && (
            <Button size="small" danger onClick={() => onDiscard(r.id)}>
              换图
            </Button>
          )}
          {r.status === 'pending_1st_review' && <span style={{ color: '#888' }}>等待审核</span>}
        </Space>
      ),
    },
  ]

  return (
    <>
      <PageHeader title="任务列表" subtitle="用户任务区：已领取、作图中、已驳回、待审核" />
      <Space style={{ marginBottom: 16 }} wrap>
        {filters.map((f) => (
          <Button key={f.key} type={status === f.key ? 'primary' : 'default'} size="small" onClick={() => setStatus(f.key)}>
            {f.label}
          </Button>
        ))}
        <Link to="/workspace/claim" style={{ marginLeft: 'auto' }}>
          <Button type="primary">+ 领图</Button>
        </Link>
      </Space>
      <Table loading={isLoading} rowKey="id" columns={columns} dataSource={data?.items ?? []} />
    </>
  )
}
