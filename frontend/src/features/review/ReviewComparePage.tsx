import { useMutation, useQuery } from '@tanstack/react-query'
import { Button, Card, Input, List, Typography, message } from 'antd'
import { useNavigate, useParams, useSearchParams, Link } from 'react-router-dom'
import { getImage } from '@/api/images'
import { listImageReviews, submitReview } from '@/api/reviews'
import { PageHeader } from '@/components/PageHeader'
import { StatusTag } from '@/components/StatusTag'
import { mediaUrl } from '@/utils/media'
import type { ImageStatus } from '@/types/image'

export function ReviewComparePage() {
  const { id } = useParams<{ id: string }>()
  const [search] = useSearchParams()
  const round = Number(search.get('round') || '1') as 1 | 2
  const navigate = useNavigate()
  const { data: img } = useQuery({ queryKey: ['image', id], queryFn: () => getImage(id!), enabled: !!id })
  const { data: history } = useQuery({
    queryKey: ['reviews', id],
    queryFn: () => listImageReviews(id!),
    enabled: !!id,
  })

  const reviewMut = useMutation({
    mutationFn: ({ result, comment }: { result: 'pass' | 'reject'; comment: string }) =>
      submitReview(id!, round, result, comment),
    onSuccess: () => {
      message.success('审核完成')
      navigate(round === 2 ? '/review/second' : '/review/first')
    },
  })

  return (
    <>
      <div style={{ marginBottom: 8 }}>
        <Link to={round === 2 ? '/review/second' : '/review/first'}>← 返回队列</Link>
      </div>
      <PageHeader
        title={`${round === 1 ? '一审' : '二审'} · ${img?.global_no ?? id}`}
        subtitle={`提交人：${img?.assignee_name ?? '—'} · 版本 v${img?.version_no ?? 1} · ${img?.category ?? ''}`}
      />
      {img?.status && <StatusTag status={img.status as ImageStatus} />}
      <Typography.Paragraph type="secondary">原图 vs 成稿对比</Typography.Paragraph>
      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 16, marginBottom: 24 }}>
        <Card title="原图（只读）" size="small">
          <div style={{ background: '#334155', minHeight: 240, borderRadius: 8, overflow: 'hidden' }}>
            {img?.storage_path && (
              <img src={mediaUrl(img.storage_path)} alt="原图" style={{ width: '100%', display: 'block' }} />
            )}
          </div>
        </Card>
        <Card title={`成稿（当前版本 v${img?.version_no ?? 1}）`} size="small">
          <div style={{ background: '#334155', minHeight: 240, borderRadius: 8, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
            {img?.image_url ? (
              <img src={img.image_url} alt="外链成稿" style={{ maxWidth: '100%', maxHeight: 240 }} />
            ) : (
              <Typography.Text style={{ color: '#94a3b8' }}>成稿预览（与左同源，MVP）</Typography.Text>
            )}
          </div>
        </Card>
      </div>
      <Card>
        <Typography.Text>审核意见</Typography.Text>
        <Input.TextArea rows={4} placeholder="通过可留空；驳回请填写修改建议" id="review-comment" style={{ marginTop: 8 }} />
        <div style={{ marginTop: 16, display: 'flex', gap: 8 }}>
          <Button
            type="primary"
            onClick={() => reviewMut.mutate({ result: 'pass', comment: '' })}
            loading={reviewMut.isPending}
          >
            通过
          </Button>
          <Button
            danger
            loading={reviewMut.isPending}
            onClick={() => {
              const c = (document.getElementById('review-comment') as HTMLTextAreaElement).value
              reviewMut.mutate({ result: 'reject', comment: c })
            }}
          >
            驳回
          </Button>
        </div>
        <Typography.Text type="secondary" style={{ fontSize: 12, display: 'block', marginTop: 12 }}>
          通过 → 一审通过库 / 终稿库 · 驳回 → 用户任务区（状态已驳回）
        </Typography.Text>
      </Card>
      <Card title="历史审核记录" style={{ marginTop: 16 }}>
        {(history?.length ?? 0) === 0 ? (
          <Typography.Text type="secondary">首次审核，无历史记录</Typography.Text>
        ) : (
          <List
            dataSource={history}
            renderItem={(r) => (
              <List.Item>
                第 {r.round} 轮 · {r.result === 'pass' ? '通过' : '驳回'} · {r.comment || '无意见'} ·{' '}
                {new Date(r.created_at).toLocaleString()}
              </List.Item>
            )}
          />
        )}
      </Card>
    </>
  )
}
