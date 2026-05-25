import { LogoutOutlined, UserOutlined } from '@ant-design/icons'
import { Button, Space, Typography } from 'antd'
import { useAuthStore } from '@/stores/authStore'
import { useLogout } from '@/hooks/useLogout'

export function AppHeaderBar({ title }: { title?: string }) {
  const user = useAuthStore((s) => s.user)
  const doLogout = useLogout()
  const displayName = user?.display_name || user?.email || '用户'

  return (
    <div
      style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        width: '100%',
        height: '100%',
      }}
    >
      <Typography.Text strong>{title}</Typography.Text>
      <Space size="middle">
        <Space size={4}>
          <UserOutlined />
          <Typography.Text>{displayName}</Typography.Text>
          {user?.role && (
            <Typography.Text type="secondary">({user.role})</Typography.Text>
          )}
        </Space>
        <Button type="link" icon={<LogoutOutlined />} onClick={doLogout}>
          退出登录
        </Button>
      </Space>
    </div>
  )
}
