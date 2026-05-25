import { useQuery } from '@tanstack/react-query'
import { Alert, Button, Table, message } from 'antd'
import { Link } from 'react-router-dom'
import { claimSecondReview, reviewQueue, secondReviewPool } from '@/api/reviews'
import { ImageThumb } from '@/components/ImageThumb'
import { PageHeader } from '@/components/PageHeader'

export function SecondReviewQueue() {
  const { data: pool, refetch: refetchPool } = useQuery({
    queryKey: ['second-pool'],
    queryFn: () => secondReviewPool(),
    refetchInterval: 10000,
  })
  const { data: mine, refetch: refetchMine, isLoading } = useQuery({
    queryKey: ['review', 2],
    queryFn: () => reviewQueue(2),
    refetchInterval: 10000,
  })

  const claimOne = async () => {
    try {
      const img = await claimSecondReview()
      message.success(`已领取 ${img.global_no}`)
      refetchPool()
      refetchMine()
    } catch {
      message.error('领取失败')
    }
  }

  return (
    <>
      <PageHeader title="二审抽样" subtitle="从「一审通过库」领取进入二审" />
      <Alert
        type="info"
        message="当前策略：人工领取 · 二审通过后进入终稿库"
        style={{ marginBottom: 16 }}
      />
      <PageHeader title="一审通过库（可领取）" />
      <Table
        style={{ marginBottom: 24 }}
        rowKey="id"
        dataSource={pool?.items ?? []}
        pagination={false}
        columns={[
          { title: '缩略图', render: (_: unknown, r: { thumb_path?: string; storage_path?: string }) => <ImageThumb path={r.thumb_path || r.storage_path} /> },
          { title: '编号', dataIndex: 'global_no' },
          { title: '类目', dataIndex: 'category' },
          {
            title: '一审通过时间',
            dataIndex: 'updated_at',
            render: (t: string) => (t ? new Date(t).toLocaleDateString() : '—'),
          },
          {
            title: '操作',
            render: () => (
              <Button type="primary" size="small" onClick={() => claimOne()}>
                领取二审
              </Button>
            ),
          },
        ]}
      />
      <PageHeader title="我已领取（待二审）" />
      <Table
        loading={isLoading}
        rowKey="id"
        dataSource={mine?.items ?? []}
        locale={{ emptyText: '暂无 · 请从上方库中领取' }}
        columns={[
          { title: '编号', dataIndex: 'global_no' },
          { title: '类目', dataIndex: 'category' },
          {
            title: '操作',
            render: (_: unknown, r: { id: string }) => <Link to={`/review/${r.id}?round=2`}>审核</Link>,
          },
        ]}
      />
    </>
  )
}
