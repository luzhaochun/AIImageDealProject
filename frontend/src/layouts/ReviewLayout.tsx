import { Layout, Menu } from 'antd'
import { Link, Outlet, useLocation } from 'react-router-dom'
import { AppHeaderBar } from '@/components/AppHeaderBar'

const { Header, Sider, Content } = Layout

const items = [
  { key: '/review/first', label: <Link to="/review/first">一审队列</Link> },
  { key: '/review/second', label: <Link to="/review/second">二审队列</Link> },
]

export function ReviewLayout() {
  const { pathname } = useLocation()
  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider theme="light" width={220}>
        <div style={{ padding: 16, fontWeight: 600 }}>ImageDeal · 审核</div>
        <Menu mode="inline" selectedKeys={[pathname]} items={items} />
      </Sider>
      <Layout>
        <Header style={{ background: '#fff', padding: '0 24px', lineHeight: '64px' }}>
          <AppHeaderBar title="审核台" />
        </Header>
        <Content style={{ margin: 24 }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}
