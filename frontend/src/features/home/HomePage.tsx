import { useQuery } from '@tanstack/react-query'
import { Card, Col, Layout, Row, Statistic, Typography } from 'antd'
import { Link } from 'react-router-dom'
import { fetchHealth } from '@/api/health'
import { useAuthStore } from '@/stores/authStore'
import { PageHeader } from '@/components/PageHeader'
import { AppHeaderBar } from '@/components/AppHeaderBar'

const { Header, Content } = Layout

export function HomePage() {
  const user = useAuthStore((s) => s.user)
  const { data, isLoading, isError } = useQuery({
    queryKey: ['health'],
    queryFn: fetchHealth,
  })

  const homeLink =
    user?.role === 'admin' ? '/admin' : user?.role === 'reviewer' ? '/review' : '/workspace'

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ background: '#fff', padding: '0 24px', lineHeight: '64px' }}>
        <AppHeaderBar title="ImageDeal 首页" />
      </Header>
      <Content style={{ padding: 24, maxWidth: 1200, margin: '0 auto', width: '100%' }}>
        <PageHeader
          title={`欢迎，${user?.display_name || user?.email || '用户'}`}
          subtitle="ImageDeal SaaS 图片协作平台"
        />
        <Row gutter={16}>
          <Col span={8}>
            <Card loading={isLoading}>
              <Statistic title="后端状态" value={data?.status ?? (isError ? '异常' : '-')} />
            </Card>
          </Col>
          <Col span={8}>
            <Card loading={isLoading}>
              <Statistic title="MySQL" value={data?.mysql ? '已连接' : '未连接'} />
            </Card>
          </Col>
          <Col span={8}>
            <Card loading={isLoading}>
              <Statistic title="Redis" value={data?.redis ? '已连接' : '未连接'} />
            </Card>
          </Col>
        </Row>
        <Typography.Paragraph type="secondary" style={{ marginTop: 24 }}>
          角色：{user?.role} ·{' '}
          <Link to={homeLink}>进入工作台</Link>
        </Typography.Paragraph>
      </Content>
    </Layout>
  )
}
