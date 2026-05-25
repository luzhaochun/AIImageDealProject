import { Layout, Menu } from 'antd'
import { Link, Outlet, useLocation } from 'react-router-dom'
import { AppHeaderBar } from '@/components/AppHeaderBar'

const { Header, Sider, Content } = Layout

const items = [
  { key: '/workspace/claim', label: <Link to="/workspace/claim">领图</Link> },
  { key: '/workspace/tasks', label: <Link to="/workspace/tasks">我的任务</Link> },
]

export function WorkspaceLayout() {
  const { pathname } = useLocation()
  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider theme="light" width={220}>
        <div style={{ padding: 16, fontWeight: 600 }}>ImageDeal · 作图</div>
        <Menu mode="inline" selectedKeys={[pathname]} items={items} />
      </Sider>
      <Layout>
        <Header style={{ background: '#fff', padding: '0 24px', lineHeight: '64px' }}>
          <AppHeaderBar title="用户工作台" />
        </Header>
        <Content style={{ margin: 24 }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}
