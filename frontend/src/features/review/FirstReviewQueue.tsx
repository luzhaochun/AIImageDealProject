import { useQuery } from '@tanstack/react-query'
import { Button, Card, Col, Row, Statistic, Table } from 'antd'
import { Link } from 'react-router-dom'
import { reviewQueue, reviewStats } from '@/api/reviews'
import { ImageThumb } from '@/components/ImageThumb'
import { PageHeader } from '@/components/PageHeader'
import { StatusTag } from '@/components/StatusTag'
import type { ImageItem, ImageStatus } from '@/types/image'

export function FirstReviewQueue() {
  const { data, isLoading } = useQuery({
    queryKey: ['review', 1],
    queryFn: () => reviewQueue(1),
    refetchInterval: 10000,
  })
  const { data: stats } = useQuery({
    queryKey: ['review-stats', 1],
    queryFn: () => reviewStats(1),
    refetchInterval: 10000,
  })

  return (
    <>
      <PageHeader title="待一审" subtitle="用户已提交成稿，请对照原图审核通过或驳回" />
      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={8}>
          <Card>
            <Statistic title="队列中" value={stats?.queue_count ?? 0} />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic title="今日已审" value={stats?.today_reviewed ?? 0} />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic title="通过率" value={stats?.pass_rate ?? 0} suffix="%" precision={0} />
          </Card>
        </Col>
      </Row>
      <Table
        loading={isLoading}
        rowKey="id"
        dataSource={data?.items ?? []}
        columns={[
          { title: '缩略图', render: (_: unknown, r: ImageItem) => <ImageThumb path={r.thumb_path || r.storage_path} /> },
          { title: '编号', dataIndex: 'global_no' },
          { title: '类目', dataIndex: 'category' },
          { title: '状态', dataIndex: 'status', render: (s: ImageStatus) => <StatusTag status={s} /> },
          {
            title: '提交时间',
            dataIndex: 'updated_at',
            render: (t: string) => (t ? new Date(t).toLocaleString() : '—'),
          },
          {
            title: '操作',
            render: (_: unknown, r: { id: string }) => (
              <Link to={`/review/${r.id}`}>
                <Button type="primary" size="small">
                  开始审核
                </Button>
              </Link>
            ),
          },
        ]}
      />
    </>
  )
}
