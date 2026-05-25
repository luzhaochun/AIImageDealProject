import { useQuery } from '@tanstack/react-query'
import { Alert, Button, Card, Col, Progress, Row, Statistic, Typography } from 'antd'
import { Link, useParams } from 'react-router-dom'
import { getImport } from '@/api/imports'
import { PageHeader } from '@/components/PageHeader'

export function ImportDetailPage() {
  const { id } = useParams<{ id: string }>()
  const { data, isLoading } = useQuery({
    queryKey: ['import', id],
    queryFn: () => getImport(id!),
    enabled: !!id,
    refetchInterval: 3000,
  })

  const progress =
    data && data.total_rows > 0
      ? Math.round(((data.downloaded ?? data.success_count) / data.total_rows) * 100)
      : 0

  const downloadErrorLog = () => {
    if (!data?.error_log) return
    const blob = new Blob([data.error_log], { type: 'text/plain' })
    const a = document.createElement('a')
    a.href = URL.createObjectURL(blob)
    a.download = `import-${id}-errors.txt`
    a.click()
  }

  return (
    <>
      <div style={{ marginBottom: 8 }}>
        <Link to="/admin/imports">← 返回导入列表</Link>
      </div>
      <PageHeader title={data?.file_name ?? '批次详情'} subtitle="自动刷新 · 每 3 秒更新进度" />
      {data?.status === 'processing' && (
        <Alert type="info" showIcon message="系统正在下载图片，完成后状态将变为已完成" style={{ marginBottom: 16 }} />
      )}
      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col span={6}>
          <Card loading={isLoading}>
            <Statistic title="总行数" value={data?.total_rows ?? 0} />
          </Card>
        </Col>
        <Col span={6}>
          <Card loading={isLoading}>
            <Statistic title="已成功" value={data?.success_count ?? 0} valueStyle={{ color: '#16a34a' }} />
          </Card>
        </Col>
        <Col span={6}>
          <Card loading={isLoading}>
            <Statistic title="已失败" value={data?.failed_count ?? 0} valueStyle={{ color: '#dc2626' }} />
          </Card>
        </Col>
        <Col span={6}>
          <Card loading={isLoading}>
            <Statistic title="待下载" value={data?.pending_storage ?? 0} />
          </Card>
        </Col>
      </Row>
      <Card title="整体进度" style={{ marginBottom: 16 }}>
        <Progress percent={progress} />
        <Typography.Text type="secondary">
          已入库 {data?.downloaded ?? 0} / {data?.total_rows ?? 0} · 待领图 {data?.pending_assign ?? 0}
        </Typography.Text>
      </Card>
      {data?.error_log && (
        <Card
          title="错误日志"
          extra={
            <Button size="small" onClick={downloadErrorLog}>
              下载完整错误日志
            </Button>
          }
        >
          <pre style={{ whiteSpace: 'pre-wrap', fontSize: 12, maxHeight: 200, overflow: 'auto' }}>{data.error_log}</pre>
        </Card>
      )}
    </>
  )
}
