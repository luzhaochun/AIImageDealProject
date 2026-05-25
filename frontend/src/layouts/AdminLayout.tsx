import { Layout, Menu } from 'antd'
import { Link, Outlet, useLocation } from 'react-router-dom'
import { AppHeaderBar } from '@/components/AppHeaderBar'

const { Header, Sider, Content } = Layout

const items = [
  { key: '/admin/imports', label: <Link to="/admin/imports">CSV 导入</Link> },
  { key: '/admin/archives', label: <Link to="/admin/archives">终稿导出</Link> },
  { key: '/admin/ai-editor', label: <Link to="/admin/ai-editor">AI 消除 Demo</Link> },
  { key: '/admin/canvas-studio', label: <Link to="/admin/canvas-studio">GPT 图像工作室</Link> },
]

export function AdminLayout() {
  const { pathname } = useLocation()
  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider theme="light" width={220}>
        <div style={{ padding: 16, fontWeight: 600 }}>ImageDeal · 管理</div>
        <Menu mode="inline" selectedKeys={[pathname]} items={items} />
      </Sider>
      <Layout>
        <Header style={{ background: '#fff', padding: '0 24px', lineHeight: '64px' }}>
          <AppHeaderBar title="管理员" />
        </Header>
        <Content style={{ margin: 24 }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}
