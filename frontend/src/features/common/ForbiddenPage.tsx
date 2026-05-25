import { Button, Result } from 'antd'
import { useNavigate } from 'react-router-dom'

export function ForbiddenPage() {
  const navigate = useNavigate()
  return (
    <Result
      status="403"
      title="403"
      subTitle="无权访问该页面"
      extra={<Button onClick={() => navigate('/')}>返回首页</Button>}
    />
  )
}
