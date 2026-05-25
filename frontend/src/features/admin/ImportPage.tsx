import { useQuery } from '@tanstack/react-query'
import { Progress, Table, Tag, Upload, message } from 'antd'
import { useState } from 'react'
import { Link } from 'react-router-dom'
import { listImports, uploadCSV, type ImportBatch } from '@/api/imports'
import { PageHeader } from '@/components/PageHeader'

export function ImportPage() {
  const [uploading, setUploading] = useState(false)
  const { data, refetch, isLoading } = useQuery({
    queryKey: ['imports'],
    queryFn: () => listImports(),
    refetchInterval: 5000,
  })

  const onUpload = async (file: File) => {
    setUploading(true)
    try {
      await uploadCSV(file)
      message.success('已创建导入批次，Worker 正在处理')
      refetch()
    } catch {
      message.error('上传失败')
    } finally {
      setUploading(false)
    }
    return false
  }

  const statusTag = (s: string) => {
    if (s === 'completed') return <Tag color="success">已完成</Tag>
    if (s === 'failed') return <Tag color="error">失败</Tag>
    return <Tag color="processing">处理中</Tag>
  }

  const columns = [
    { title: '批次', dataIndex: 'id', render: (id: string) => <Link to={`/admin/imports/${id}`}>{id.slice(0, 8)}…</Link> },
    { title: '文件', dataIndex: 'file_name' },
    { title: '状态', dataIndex: 'status', render: statusTag },
    {
      title: '进度',
      render: (_: unknown, r: ImportBatch) =>
        r.total_rows > 0 ? (
          <Progress percent={Math.round((r.success_count / r.total_rows) * 100)} size="small" />
        ) : (
          '—'
        ),
    },
    { title: '成功/失败', render: (_: unknown, r: ImportBatch) => `${r.success_count} / ${r.failed_count}` },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      render: (t: string) => (t ? new Date(t).toLocaleString() : '—'),
    },
    {
      title: '操作',
      render: (_: unknown, r: ImportBatch) => (
        <>
          <Link to={`/admin/imports/${r.id}`}>详情</Link>
          {r.error_log && (
            <>
              {' · '}
              <a
                href="#"
                onClick={(e) => {
                  e.preventDefault()
                  const blob = new Blob([r.error_log!], { type: 'text/plain' })
                  const a = document.createElement('a')
                  a.href = URL.createObjectURL(blob)
                  a.download = `errors-${r.id}.txt`
                  a.click()
                }}
              >
                错误日志
              </a>
            </>
          )}
        </>
      ),
    },
  ]

  return (
    <>
      <PageHeader title="CSV 导入" subtitle="上传后由 Redis 队列异步解析并下载入库" />
      <Upload.Dragger beforeUpload={onUpload} showUploadList={false} disabled={uploading}>
        <p>拖拽或点击上传 CSV（列：image_url, category）</p>
      </Upload.Dragger>
      <Table style={{ marginTop: 24 }} loading={isLoading} rowKey="id" columns={columns} dataSource={data?.items ?? []} />
    </>
  )
}
