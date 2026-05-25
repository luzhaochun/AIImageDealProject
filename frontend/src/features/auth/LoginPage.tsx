import { Button, Card, Form, Input, message } from 'antd'
import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { login } from '@/api/auth'
import { useAuthStore } from '@/stores/authStore'

export function LoginPage() {
  const navigate = useNavigate()
  const token = useAuthStore((s) => s.token)
  const user = useAuthStore((s) => s.user)
  const setAuth = useAuthStore((s) => s.setAuth)

  useEffect(() => {
    if (!token || !user) return
    const dest =
      user.role === 'admin' ? '/admin' : user.role === 'reviewer' ? '/review' : '/workspace'
    navigate(dest, { replace: true })
  }, [token, user, navigate])

  const onFinish = async (values: { email: string; password: string }) => {
    try {
      const res = await login(values.email, values.password)
      setAuth(res.token, res.user)
      const dest =
        res.user.role === 'admin' ? '/admin' : res.user.role === 'reviewer' ? '/review' : '/workspace'
      navigate(dest)
    } catch {
      message.error('登录失败，请检查账号密码')
    }
  }

  return (
    <div style={{ minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center', background: '#f5f5f5' }}>
      <Card style={{ width: 400 }}>
        <h2>ImageDeal 登录</h2>
        <p style={{ color: '#888', marginBottom: 24 }}>测试账号 password123：admin@ / user@ / reviewer@ example.com</p>
        <Form layout="vertical" onFinish={onFinish} initialValues={{ email: 'user@example.com', password: 'password123' }}>
          <Form.Item name="email" label="邮箱" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="password" label="密码" rules={[{ required: true }]}>
            <Input.Password />
          </Form.Item>
          <Button type="primary" htmlType="submit" block>
            登录
          </Button>
        </Form>
      </Card>
    </div>
  )
}
