import { useMutation, useQuery } from '@tanstack/react-query'
import { Alert, Button, Card, Input, List, Modal, Typography, message } from 'antd'
import { useNavigate, useParams, Link } from 'react-router-dom'
import { discardImage, getImage, listVersions, saveVersion, submitReview } from '@/api/images'
import { PageHeader } from '@/components/PageHeader'
import { StatusTag } from '@/components/StatusTag'
import { mediaUrl } from '@/utils/media'
import type { ImageStatus } from '@/types/image'

export function EditorPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { data: img, refetch } = useQuery({ queryKey: ['image', id], queryFn: () => getImage(id!), enabled: !!id })
  const { data: versions } = useQuery({ queryKey: ['versions', id], queryFn: () => listVersions(id!), enabled: !!id })

  const saveMut = useMutation({
    mutationFn: (layer: object) => saveVersion(id!, layer),
    onSuccess: () => {
      message.success('已保存版本')
      refetch()
    },
  })
  const submitMut = useMutation({
    mutationFn: () => submitReview(id!),
    onSuccess: () => {
      message.success('已提交一审')
      navigate('/workspace/tasks')
    },
  })

  const defaultLayer = '{"width":1920,"height":1080,"layers":[{"type":"text","content":"标题"}]}'

  return (
    <>
      <div style={{ marginBottom: 8, display: 'flex', gap: 12, alignItems: 'center' }}>
        <Link to="/workspace/tasks">← 返回任务</Link>
        {img?.status && <StatusTag status={img.status as ImageStatus} />}
      </div>
      <PageHeader title={`${img?.global_no ?? id} · ${img?.category ?? ''}`} subtitle="在线编辑器" />
      {img?.status === 'rejected' && img.last_review_comment && (
        <Alert type="error" message={`驳回意见：${img.last_review_comment}`} style={{ marginBottom: 16 }} />
      )}
      <div style={{ display: 'flex', gap: 8, marginBottom: 16, flexWrap: 'wrap' }}>
        <Button size="small">裁剪</Button>
        <Button size="small">文字</Button>
        <Button size="small">调色</Button>
        <span style={{ flex: 1 }} />
        <Button
          onClick={() => {
            const el = document.getElementById('layer-json') as HTMLTextAreaElement
            try {
              saveMut.mutate(JSON.parse(el.value))
            } catch {
              message.error('JSON 格式错误')
            }
          }}
          loading={saveMut.isPending}
        >
          保存版本
        </Button>
        <Button type="primary" loading={submitMut.isPending} onClick={() => submitMut.mutate()}>
          提交审核
        </Button>
      </div>
      <div style={{ display: 'grid', gridTemplateColumns: '220px 1fr', gap: 16 }}>
        <Card size="small" title="图层 / 版本">
          <Typography.Text type="secondary" style={{ fontSize: 12 }}>
            图层 JSON 编辑（MVP）
          </Typography.Text>
          <List
            size="small"
            style={{ marginTop: 8 }}
            dataSource={versions ?? []}
            renderItem={(v) => (
              <List.Item>
                v{v.version_no}
                {v.is_current ? '（当前）' : ''}
              </List.Item>
            )}
          />
          <Button
            size="small"
            danger
            block
            style={{ marginTop: 12 }}
            onClick={() => {
              Modal.confirm({
                title: '换图（废弃）',
                content: '确定废弃当前图片并重新领图？',
                onOk: async () => {
                  await discardImage(id!, '编辑器内换图')
                  message.success('已废弃')
                  navigate('/workspace/claim')
                },
              })
            }}
          >
            换图（废弃）
          </Button>
        </Card>
        <div>
          <div
            style={{
              background: '#1e293b',
              minHeight: 360,
              borderRadius: 8,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              overflow: 'hidden',
            }}
          >
            {img?.storage_path ? (
              <img src={mediaUrl(img.storage_path)} alt="原图" style={{ maxWidth: '100%', maxHeight: 480 }} />
            ) : (
              <Typography.Text style={{ color: '#94a3b8' }}>画布预览</Typography.Text>
            )}
          </div>
          <Input.TextArea
            id="layer-json"
            rows={6}
            style={{ marginTop: 12 }}
            placeholder="图层 JSON"
            defaultValue={defaultLayer}
          />
        </div>
      </div>
      <Alert
        type="info"
        message="提交审核前请确认成稿符合规范。提交后进入一审队列。"
        style={{ marginTop: 16 }}
      />
    </>
  )
}
