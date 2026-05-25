import { useQuery } from '@tanstack/react-query'
import { Alert, Button, Card, Col, Row, Table, Typography, message } from 'antd'
import { useNavigate } from 'react-router-dom'
import { apiClient } from '@/api/client'
import { claimImage, listCategories, listImages } from '@/api/images'
import type { ApiResponse } from '@/types/api'
import { PageHeader } from '@/components/PageHeader'
import { StatusTag } from '@/components/StatusTag'
import type { ImageStatus } from '@/types/image'

interface UserProfile {
  discard_used_today: number
  max_discard_per_day: number
}

export function ClaimPage() {
  const navigate = useNavigate()
  const { data, refetch, isLoading } = useQuery({ queryKey: ['categories'], queryFn: listCategories })
  const { data: recent } = useQuery({
    queryKey: ['recent-claims'],
    queryFn: () => listImages({ library: 'tasks', page: 1, page_size: 5 }),
  })
  const { data: profile } = useQuery({
    queryKey: ['me'],
    queryFn: async () => {
      const { data } = await apiClient.get<ApiResponse<UserProfile>>('/auth/me')
      return data.data!
    },
  })

  const onClaim = async (category: string) => {
    try {
      const img = await claimImage(category)
      message.success(`已领取 ${img.global_no}`)
      refetch()
      navigate('/workspace/tasks')
    } catch {
      message.error('领取失败，可能无库存')
    }
  }

  return (
    <>
      <PageHeader title="领取图片" subtitle="从原图库按类目领取一张，系统自动分配全局编号 global_no" />
      <Row gutter={16}>
        <Col span={14}>
          <Card title="选择类目">
            <Table
              loading={isLoading}
              rowKey="name"
              dataSource={data ?? []}
              pagination={false}
              columns={[
                { title: '类目', dataIndex: 'name' },
                { title: '可领数量', dataIndex: 'available_count' },
                {
                  title: '',
                  render: (_: unknown, r: { name: string; available_count: number }) =>
                    r.available_count === 0 ? (
                      <Typography.Text type="secondary">暂无库存</Typography.Text>
                    ) : (
                      <Button type="primary" size="small" onClick={() => onClaim(r.name)}>
                        领取一张
                      </Button>
                    ),
                },
              ]}
            />
            <Typography.Text type="secondary" style={{ fontSize: 12, display: 'block', marginTop: 12 }}>
              每日换图上限 {profile?.max_discard_per_day ?? 10} 次（今日已用 {profile?.discard_used_today ?? 0}）· 进行中任务上限 5 个
            </Typography.Text>
          </Card>
        </Col>
        <Col span={10}>
          <Card title="领图说明">
            <ul style={{ paddingLeft: 20, fontSize: 13, color: '#64748b', lineHeight: 1.8 }}>
              <li>领取后图片进入「我的任务」，状态为已领取</li>
              <li>若原图不可用，可「换图」废弃当前图并重新领取</li>
              <li>废弃图片进入废料堆，不可再次领取</li>
            </ul>
            <Alert type="warning" message="换图将记录原因并计入每日限额" style={{ marginTop: 16 }} />
          </Card>
        </Col>
      </Row>
      <Card title="最近领取" style={{ marginTop: 16 }}>
        <Table
          rowKey="id"
          dataSource={recent?.items ?? []}
          pagination={false}
          columns={[
            { title: '编号', dataIndex: 'global_no' },
            { title: '类目', dataIndex: 'category' },
            { title: '状态', dataIndex: 'status', render: (s: ImageStatus) => <StatusTag status={s} /> },
          ]}
        />
      </Card>
    </>
  )
}
