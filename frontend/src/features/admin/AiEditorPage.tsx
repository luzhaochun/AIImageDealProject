import { useCallback, useRef, useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Alert, Button, Input, Radio, Slider, Space, Typography, Upload, message } from 'antd'
import { DownloadOutlined, PictureOutlined, RollbackOutlined, ThunderboltOutlined } from '@ant-design/icons'
import { aiInpaint, getAiConfig } from '@/api/aiEditor'
import axios from 'axios'

type Tool = 'brush' | 'eraser'

function maskHasPixels(canvas: HTMLCanvasElement): boolean {
  const ctx = canvas.getContext('2d')
  if (!ctx) return false
  const data = ctx.getImageData(0, 0, canvas.width, canvas.height).data
  for (let i = 3; i < data.length; i += 4) {
    if (data[i] > 20) return true
  }
  return false
}

export function AiEditorPage() {
  const imageCanvasRef = useRef<HTMLCanvasElement>(null)
  const maskCanvasRef = useRef<HTMLCanvasElement>(null)
  const wrapRef = useRef<HTMLDivElement>(null)
  const [tool, setTool] = useState<Tool>('brush')
  const [brushSize, setBrushSize] = useState(18)
  const [prompt, setPrompt] = useState('')
  const [imageLoaded, setImageLoaded] = useState(false)
  const [displaySize, setDisplaySize] = useState({ w: 0, h: 0 })
  const [resultUrl, setResultUrl] = useState<string | null>(null)
  const [originalSnapshot, setOriginalSnapshot] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const drawing = useRef(false)

  const { data: cfg } = useQuery({ queryKey: ['ai-config'], queryFn: getAiConfig })

  const setupImage = useCallback((img: HTMLImageElement) => {
    const maxW = 880
    const scale = Math.min(1, maxW / img.width)
    const w = Math.round(img.width * scale)
    const h = Math.round(img.height * scale)
    const ic = imageCanvasRef.current!
    const mc = maskCanvasRef.current!
    ic.width = mc.width = w
    ic.height = mc.height = h
    ic.getContext('2d')!.drawImage(img, 0, 0, w, h)
    mc.getContext('2d')!.clearRect(0, 0, w, h)
    setDisplaySize({ w, h })
    setImageLoaded(true)
    setResultUrl(null)
  }, [])

  const paint = (e: React.MouseEvent<HTMLCanvasElement>) => {
    const canvas = maskCanvasRef.current
    if (!canvas) return
    const rect = canvas.getBoundingClientRect()
    const scaleX = canvas.width / rect.width
    const scaleY = canvas.height / rect.height
    const x = (e.clientX - rect.left) * scaleX
    const y = (e.clientY - rect.top) * scaleY
    const ctx = canvas.getContext('2d')
    if (!ctx) return
    ctx.lineCap = 'round'
    ctx.lineJoin = 'round'
    ctx.lineWidth = brushSize
    if (tool === 'brush') {
      ctx.globalCompositeOperation = 'source-over'
      ctx.strokeStyle = 'rgba(255, 40, 40, 1)'
      ctx.fillStyle = 'rgba(255, 40, 40, 1)'
    } else {
      ctx.globalCompositeOperation = 'destination-out'
      ctx.strokeStyle = 'rgba(0,0,0,1)'
    }
    ctx.lineTo(x, y)
    ctx.stroke()
    ctx.beginPath()
    ctx.moveTo(x, y)
  }

  const startDraw = (e: React.MouseEvent<HTMLCanvasElement>) => {
    drawing.current = true
    maskCanvasRef.current?.getContext('2d')?.beginPath()
    paint(e)
  }
  const endDraw = () => {
    drawing.current = false
  }
  const moveDraw = (e: React.MouseEvent<HTMLCanvasElement>) => {
    if (!drawing.current) return
    paint(e)
  }

  const loadImage = (file: File) => {
    const url = URL.createObjectURL(file)
    const img = new Image()
    img.onload = () => {
      setupImage(img)
      URL.revokeObjectURL(url)
    }
    img.src = url
    return false
  }

  const loadDemoImage = () => {
    const img = new Image()
    img.onload = () => setupImage(img)
    img.onerror = () => message.error('示例图加载失败，请手动上传')
    img.src = '/demo-food.jpg'
  }

  const clearMask = () => {
    const m = maskCanvasRef.current
    if (!m) return
    m.getContext('2d')?.clearRect(0, 0, m.width, m.height)
  }

  const runInpaint = async () => {
    const ic = imageCanvasRef.current
    const mc = maskCanvasRef.current
    if (!ic || !mc || !imageLoaded) {
      message.warning('请先上传图片')
      return
    }
    if (!maskHasPixels(mc)) {
      message.warning('请用红色画笔涂抹要消除的区域')
      return
    }
    setLoading(true)
    try {
      setOriginalSnapshot(ic.toDataURL('image/png'))
      const imageBlob = await new Promise<Blob>((resolve, reject) => {
        ic.toBlob((b) => (b ? resolve(b) : reject(new Error('export failed'))), 'image/png')
      })
      const maskBlob = await new Promise<Blob>((resolve, reject) => {
        mc.toBlob((b) => (b ? resolve(b) : reject(new Error('mask failed'))), 'image/png')
      })
      const res = await aiInpaint(imageBlob, maskBlob, prompt)
      setResultUrl(res.result_url)
      message.success(`消除完成 · ${res.provider}`)
      if (res.message) message.info(res.message, 4)
    } catch (e) {
      let msg = 'AI 消除失败'
      if (axios.isAxiosError(e)) {
        if (!e.response) {
          msg = '无法连接后端，请在 backend 目录执行: go run .'
        } else {
          const body = e.response.data as { message?: string; code?: string }
          msg = body?.message || msg
          if (e.response.status === 403) msg = '请使用 admin 账号登录'
        }
      }
      message.error(msg)
    } finally {
      setLoading(false)
    }
  }

  const applyResultToCanvas = () => {
    if (!resultUrl) return
    const img = new Image()
    img.onload = () => {
      setupImage(img)
      clearMask()
      message.success('已应用为当前画布')
    }
    img.src = resultUrl
  }

  return (
    <div style={{ margin: -24, minHeight: 'calc(100vh - 64px)', background: '#0f172a', color: '#e2e8f0' }}>
      <div style={{ padding: '12px 20px', borderBottom: '1px solid #334155', display: 'flex', alignItems: 'center', gap: 12, flexWrap: 'wrap' }}>
        <Typography.Title level={4} style={{ color: '#fff', margin: 0 }}>
          AI 消除 · 面试 Demo
        </Typography.Title>
        <Typography.Text style={{ color: '#64748b' }}>
          模式：{cfg?.mode ?? 'mock'}
          {cfg?.mode === 'mock' || !cfg?.has_key ? ' · 本地模拟 · 零成本' : ` · ${cfg?.model}`}
        </Typography.Text>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '220px 1fr 240px', minHeight: 560 }}>
        <div style={{ padding: 16, borderRight: '1px solid #334155' }}>
          <Typography.Text strong style={{ color: '#94a3b8' }}>
            工具
          </Typography.Text>
          <Radio.Group
            value={tool}
            onChange={(e) => setTool(e.target.value)}
            style={{ display: 'flex', flexDirection: 'column', gap: 8, marginTop: 12 }}
            optionType="button"
            buttonStyle="solid"
            options={[
              { label: '画笔', value: 'brush' },
              { label: '橡皮', value: 'eraser' },
            ]}
          />
          <div style={{ marginTop: 20 }}>
            <Typography.Text style={{ color: '#94a3b8' }}>笔刷大小 {brushSize}px</Typography.Text>
            <Slider min={8} max={80} value={brushSize} onChange={setBrushSize} />
          </div>
          <Typography.Text style={{ color: '#64748b', fontSize: 12, display: 'block', marginTop: 16 }}>
            消除提示词（mock 下仅展示，不生效）
          </Typography.Text>
          <Input.TextArea
            rows={3}
            placeholder="描述要填充的背景…"
            value={prompt}
            onChange={(e) => setPrompt(e.target.value)}
            style={{ marginTop: 8, background: '#1e293b', color: '#e2e8f0' }}
          />
          <Button block style={{ marginTop: 12 }} onClick={clearMask}>
            清除蒙版
          </Button>
          <Upload showUploadList={false} beforeUpload={loadImage} accept="image/*">
            <Button block style={{ marginTop: 8 }}>
              上传图片
            </Button>
          </Upload>
          <Button block icon={<PictureOutlined />} style={{ marginTop: 8 }} onClick={loadDemoImage}>
            加载演示图（披萨）
          </Button>
        </div>

        <div style={{ padding: 16, display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 16 }}>
          <Alert
            type="warning"
            showIcon
            message="只涂要删的刀叉/小物体，不要涂整块披萨。涂完后点 AI 消除。mock 为模糊模拟，真 AI 需接云 API。"
            style={{ width: '100%', maxWidth: 920 }}
          />
          <div
            ref={wrapRef}
            style={{
              position: 'relative',
              lineHeight: 0,
              width: displaySize.w || undefined,
              maxWidth: '100%',
            }}
          >
            <canvas
              ref={imageCanvasRef}
              style={{
                width: displaySize.w || '100%',
                height: displaySize.h || 'auto',
                borderRadius: 8,
                display: 'block',
                background: '#1e293b',
              }}
            />
            <canvas
              ref={maskCanvasRef}
              style={{
                position: 'absolute',
                left: 0,
                top: 0,
                width: displaySize.w || '100%',
                height: displaySize.h || 'auto',
                borderRadius: 8,
                cursor: 'crosshair',
              }}
              onMouseDown={startDraw}
              onMouseUp={endDraw}
              onMouseLeave={endDraw}
              onMouseMove={moveDraw}
            />
          </div>

          {resultUrl && originalSnapshot && (
            <div style={{ width: '100%', maxWidth: 920 }}>
              <Typography.Text style={{ color: '#94a3b8' }}>原图 vs 消除结果</Typography.Text>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 12, marginTop: 8 }}>
                <div>
                  <Typography.Text type="secondary" style={{ fontSize: 12, color: '#94a3b8' }}>
                    原图
                  </Typography.Text>
                  <img src={originalSnapshot} alt="before" style={{ width: '100%', borderRadius: 8, marginTop: 4 }} />
                </div>
                <div>
                  <Typography.Text type="secondary" style={{ fontSize: 12, color: '#94a3b8' }}>
                    AI 消除结果
                  </Typography.Text>
                  <img src={resultUrl} alt="after" style={{ width: '100%', borderRadius: 8, marginTop: 4 }} />
                </div>
              </div>
            </div>
          )}
        </div>

        <div style={{ padding: 16, borderLeft: '1px solid #334155' }}>
          <Typography.Text strong style={{ color: '#94a3b8' }}>
            图层
          </Typography.Text>
          <div style={{ marginTop: 12, fontSize: 13 }}>
            <div style={{ padding: '8px 0', borderBottom: '1px solid #334155' }}>☑ 蒙版</div>
            <div style={{ padding: '8px 0', color: resultUrl ? '#38bdf8' : '#64748b' }}>AI 消除结果</div>
            <div style={{ padding: '8px 0' }}>原图</div>
          </div>
          <Space direction="vertical" style={{ width: '100%', marginTop: 24 }}>
            <Button type="primary" block size="large" icon={<ThunderboltOutlined />} loading={loading} onClick={runInpaint}>
              AI 消除
            </Button>
            {resultUrl && (
              <>
                <Button block icon={<DownloadOutlined />} href={resultUrl} download="ai-result.png">
                  下载结果
                </Button>
                <Button block onClick={applyResultToCanvas}>
                  应用结果到画布
                </Button>
              </>
            )}
            <Button block icon={<RollbackOutlined />} onClick={() => setResultUrl(null)}>
              隐藏对比
            </Button>
          </Space>
        </div>
      </div>
    </div>
  )
}
