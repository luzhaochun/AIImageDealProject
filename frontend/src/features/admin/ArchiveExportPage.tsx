import { useQuery } from '@tanstack/react-query'
import { Button, Col, Form, Input, Row, Select, Table, message } from 'antd'
import { completeArchive, createExport, listArchives, listExports } from '@/api/archives'
import { PageHeader } from '@/components/PageHeader'
import { StatusTag } from '@/components/StatusTag'
import { mediaUrl } from '@/utils/media'
import type { ImageStatus } from '@/types/image'

export function ArchiveExportPage() {
  const [form] = Form.useForm()
  const { data, refetch, isLoading } = useQuery({
    queryKey: ['archives'],
    queryFn: () => listArchives(),
  })
  const { data: exports, refetch: refetchExports } = useQuery({
    queryKey: ['exports'],
    queryFn: listExports,
  })

  const onExport = async () => {
    const v = form.getFieldsValue()
    try {
      const job = await createExport({
        category: v.category || undefined,
        status: v.status || 'completed',
        date_from: v.date_from || undefined,
        date_to: v.date_to || undefined,
      })
      message.success(`导出完成，共 ${job.image_count} 张`)
      refetchExports()
    } catch {
      message.error('导出失败')
    }
  }

  return (
    <>
      <PageHeader title="终稿导出" subtitle="筛选终稿库图片，打包 ZIP + manifest.csv" />
      <Form form={form} layout="vertical" style={{ marginBottom: 24, maxWidth: 720 }}>
        <Row gutter={16}>
          <Col span={12}>
            <Form.Item name="category" label="分类">
              <Select allowClear placeholder="全部" options={[{ value: 'banner' }, { value: 'poster' }, { value: 'icon' }, { value: 'social' }]} />
            </Form.Item>
          </Col>
          <Col span={12}>
            <Form.Item name="status" label="状态" initialValue="completed">
              <Select
                options={[
                  { value: 'completed', label: '已归档 completed' },
                  { value: '2nd_review_passed', label: '二审通过' },
                ]}
              />
            </Form.Item>
          </Col>
          <Col span={12}>
            <Form.Item name="date_from" label="开始日期">
              <Input type="date" />
            </Form.Item>
          </Col>
          <Col span={12}>
            <Form.Item name="date_to" label="结束日期">
              <Input type="date" />
            </Form.Item>
          </Col>
        </Row>
        <Button type="primary" onClick={onExport}>
          创建导出任务
        </Button>
      </Form>

      <PageHeader title="导出任务历史" />
      <Table
        style={{ marginBottom: 32 }}
        rowKey="id"
        dataSource={exports ?? []}
        pagination={false}
        columns={[
          { title: '任务 ID', dataIndex: 'id', render: (id: string) => id.slice(0, 8) },
          { title: '条件', dataIndex: 'filter_label' },
          { title: '状态', dataIndex: 'status', render: (s: string) => (s === 'completed' ? '已完成' : s) },
          { title: '张数', dataIndex: 'image_count' },
          {
            title: '操作',
            render: (_: unknown, r: { download_path?: string; file_name: string }) =>
              r.download_path ? (
                <a href={mediaUrl(r.download_path)} download={r.file_name}>
                  下载 ZIP
                </a>
              ) : (
                '—'
              ),
          },
        ]}
      />

      <PageHeader title="终稿库" subtitle="二审通过后可标记为已归档" />
      <Table
        loading={isLoading}
        rowKey="id"
        dataSource={data?.items ?? []}
        columns={[
          { title: '编号', dataIndex: 'global_no' },
          { title: '类目', dataIndex: 'category' },
          { title: '状态', dataIndex: 'status', render: (s: ImageStatus) => <StatusTag status={s} /> },
          {
            title: '操作',
            render: (_: unknown, r: { id: string; status: string }) =>
              r.status !== 'completed' ? (
                <Button
                  size="small"
                  onClick={async () => {
                    await completeArchive(r.id)
                    message.success('已归档')
                    refetch()
                  }}
                >
                  归档
                </Button>
              ) : (
                '已归档'
              ),
          },
        ]}
      />
    </>
  )
}
