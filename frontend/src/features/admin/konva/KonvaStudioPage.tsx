import { useCallback, useEffect, useRef, useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Alert, Button, Input, Radio, Select, Slider, Space, Switch, Typography, Upload, message } from 'antd'
import {
  ClearOutlined,
  DownloadOutlined,
  FileImageOutlined,
  ScissorOutlined,
  ThunderboltOutlined,
  UndoOutlined,
} from '@ant-design/icons'
import { Layer, Rect, Stage, Image as KonvaImage, Transformer } from 'react-konva'
import type Konva from 'konva'
import axios from 'axios'
import { aiStudioInpaint, getAiConfig } from '@/api/aiEditor'
import {
  EXPORT_FORMAT_OPTIONS,
  downloadCanvasAsImage,
  type StudioExportFormat,
} from '@/utils/exportImage'
import { downloadPsd } from '@/utils/exportPsd'
import {
  buildMaskPreview,
  canvasToBlob,
  canvasToImage,
  clearMask,
  createMaskCanvas,
  cropCanvas,
  downscaleCanvasForApi,
  featherMaskCopy,
  imageToCanvas,
  loadImageFromFile,
  loadImageFromUrl,
  maskHasPixels,
} from '@/utils/maskCanvas'

type Tool = 'brush' | 'eraser' | 'move' | 'crop'

const MAX_DISPLAY_W = 920

function fitSize(nw: number, nh: number) {
  const scale = Math.min(1, MAX_DISPLAY_W / nw)
  return { w: Math.round(nw * scale), h: Math.round(nh * scale), scale }
}

export function KonvaStudioPage() {
  const stageRef = useRef<Konva.Stage>(null)
  const transformerRef = useRef<Konva.Transformer>(null)
  const cropRectRef = useRef<Konva.Rect>(null)
  const cropTransformerRef = useRef<Konva.Transformer>(null)
  const originalNodeRef = useRef<Konva.Image>(null)

  const maskCanvasRef = useRef<HTMLCanvasElement | null>(null)
  const originalCanvasRef = useRef<HTMLCanvasElement | null>(null)
  const resultCanvasRef = useRef<HTMLCanvasElement | null>(null)

  const [tool, setTool] = useState<Tool>('brush')
  const [brushSize, setBrushSize] = useState(16)
  const [feather, setFeather] = useState(4)
  const [prompt, setPrompt] = useState(
    '移除蒙版区域内的物体，用与周围环境一致的自然背景无缝填充，保持真实照片质感。',
  )
  const [stageSize, setStageSize] = useState({ w: 0, h: 0 })
  const [konvaImage, setKonvaImage] = useState<HTMLImageElement | null>(null)
  const [resultKonvaImage, setResultKonvaImage] = useState<HTMLImageElement | null>(null)
  const [maskPreview, setMaskPreview] = useState<HTMLCanvasElement | null>(null)
  const [showOriginal, setShowOriginal] = useState(true)
  const [showResult, setShowResult] = useState(true)
  const [showMask, setShowMask] = useState(true)
  const [hasResult, setHasResult] = useState(false)
  const [exportFormat, setExportFormat] = useState<StudioExportFormat>('png')
  const [loading, setLoading] = useState(false)
  const [providerMsg, setProviderMsg] = useState('')
  const [sessionKey, setSessionKey] = useState(0)
  const drawing = useRef(false)

  const { data: cfg } = useQuery({ queryKey: ['ai-config'], queryFn: getAiConfig })

  const refreshMaskPreview = useCallback(() => {
    const m = maskCanvasRef.current
    if (!m) return
    setMaskPreview(buildMaskPreview(m))
  }, [])

  const initFromImage = useCallback(
    (img: HTMLImageElement) => {
      const { w, h } = fitSize(img.naturalWidth, img.naturalHeight)
      drawing.current = false
      maskCanvasRef.current = createMaskCanvas(w, h)
      originalCanvasRef.current = imageToCanvas(img, w, h)
      resultCanvasRef.current = imageToCanvas(img, w, h)
      setKonvaImage(img)
      setResultKonvaImage(null)
      setHasResult(false)
      setShowResult(false)
      setStageSize({ w, h })
      setProviderMsg('')
      setTool('brush')
      setSessionKey((k) => k + 1)
      refreshMaskPreview()
    },
    [refreshMaskPreview],
  )

  const handleUpload = async (file: File) => {
    try {
      const img = await loadImageFromFile(file)
      initFromImage(img)
    } catch {
      message.error('图片加载失败')
    }
    return false
  }

  const loadDemo = async () => {
    try {
      const img = await loadImageFromUrl('/demo-food.jpg')
      initFromImage(img)
    } catch {
      message.warning('请上传图片（可将演示图放到 frontend/public/demo-food.jpg）')
    }
  }

  const clearMaskAction = () => {
    const m = maskCanvasRef.current
    if (!m) return
    clearMask(m)
    refreshMaskPreview()
  }

  const getPointerOnMask = () => {
    const stage = stageRef.current
    const m = maskCanvasRef.current
    if (!stage || !m || m.width === 0 || m.height === 0) return null
    const pos = stage.getPointerPosition()
    if (!pos) return null
    const sx = stage.width()
    const sy = stage.height()
    if (sx === 0 || sy === 0) return null
    return { x: pos.x * (m.width / sx), y: pos.y * (m.height / sy) }
  }

  const strokeOnMask = (x: number, y: number) => {
    const ctx = maskCanvasRef.current?.getContext('2d')
    if (!ctx) return
    ctx.lineCap = 'round'
    ctx.lineJoin = 'round'
    ctx.lineWidth = brushSize
    if (tool === 'brush') {
      ctx.globalCompositeOperation = 'source-over'
      ctx.strokeStyle = 'rgba(255, 50, 50, 1)'
      ctx.fillStyle = 'rgba(255, 50, 50, 1)'
    } else {
      ctx.globalCompositeOperation = 'destination-out'
      ctx.strokeStyle = 'rgba(0,0,0,1)'
    }
    ctx.lineTo(x, y)
    ctx.stroke()
    ctx.beginPath()
    ctx.moveTo(x, y)
    refreshMaskPreview()
  }

  const onPointerDown = () => {
    if (!maskCanvasRef.current || tool === 'move' || tool === 'crop') return
    const p = getPointerOnMask()
    if (!p) return
    drawing.current = true
    const ctx = maskCanvasRef.current.getContext('2d')
    ctx?.beginPath()
    ctx?.moveTo(p.x, p.y)
    strokeOnMask(p.x, p.y)
  }

  const onPointerMove = () => {
    if (!drawing.current || tool === 'move' || tool === 'crop') return
    const p = getPointerOnMask()
    if (!p) return
    strokeOnMask(p.x, p.y)
  }

  const onPointerUp = () => {
    drawing.current = false
  }

  useEffect(() => {
    if (tool !== 'move' || !transformerRef.current || !originalNodeRef.current) {
      transformerRef.current?.nodes([])
    } else {
      transformerRef.current.nodes([originalNodeRef.current])
      transformerRef.current.getLayer()?.batchDraw()
    }
    if (tool !== 'crop' || !cropTransformerRef.current || !cropRectRef.current) {
      cropTransformerRef.current?.nodes([])
    } else {
      cropTransformerRef.current.nodes([cropRectRef.current])
      cropTransformerRef.current.getLayer()?.batchDraw()
    }
  }, [tool, konvaImage, stageSize.w, sessionKey])

  useEffect(() => {
    const node = originalNodeRef.current
    if (!node) return
    node.scaleX(1)
    node.scaleY(1)
    node.rotation(0)
    node.position({ x: 0, y: 0 })
    node.getLayer()?.batchDraw()
  }, [sessionKey])

  const applyCrop = async () => {
    const node = cropRectRef.current
    const orig = originalCanvasRef.current
    const mask = maskCanvasRef.current
    if (!node || !orig || !mask) {
      message.warning('请先上传图片')
      return
    }
    const x = Math.max(0, Math.round(node.x()))
    const y = Math.max(0, Math.round(node.y()))
    const w = Math.min(stageSize.w - x, Math.max(20, Math.round(node.width() * node.scaleX())))
    const h = Math.min(stageSize.h - y, Math.max(20, Math.round(node.height() * node.scaleY())))
    try {
      const croppedOrig = cropCanvas(orig, x, y, w, h)
      const croppedMask = cropCanvas(mask, x, y, w, h)
      const img = await canvasToImage(croppedOrig)
      maskCanvasRef.current = croppedMask
      originalCanvasRef.current = croppedOrig
      if (hasResult && resultCanvasRef.current) {
        const croppedResult = cropCanvas(resultCanvasRef.current, x, y, w, h)
        resultCanvasRef.current = croppedResult
        setResultKonvaImage(await canvasToImage(croppedResult))
      } else {
        resultCanvasRef.current = imageToCanvas(img, w, h)
        setResultKonvaImage(null)
        setHasResult(false)
      }
      setKonvaImage(img)
      setStageSize({ w, h })
      node.scaleX(1)
      node.scaleY(1)
      node.position({ x: 0, y: 0 })
      node.size({ width: w, height: h })
      refreshMaskPreview()
      setTool('brush')
      message.success('剪裁已应用')
    } catch (e) {
      message.error(e instanceof Error ? e.message : '剪裁失败')
    }
  }

  const studioReady = cfg?.studio_ready === true

  const runInpaint = async () => {
    const orig = originalCanvasRef.current
    const mask = maskCanvasRef.current
    if (!orig || !mask) {
      message.warning('请先上传图片')
      return
    }
    if (!studioReady) {
      message.error('未配置 GPT Image：请在 backend/conf/app.conf 设置 ai_edit_mode=openai 并填写 openai_api_key，然后重启后端')
      return
    }
    if (!maskHasPixels(mask)) {
      message.warning('请用画笔涂抹要消除的区域')
      return
    }
    setLoading(true)
    setProviderMsg('正在调用 GPT Image（quality=medium，通常 20–90 秒）…')
    const t0 = Date.now()
    try {
      const maskForApi = feather > 0 ? featherMaskCopy(mask, feather) : mask
      const apiOrig = downscaleCanvasForApi(orig)
      const apiMask = downscaleCanvasForApi(maskForApi)
      const imageBlob = await canvasToBlob(apiOrig)
      const maskBlob = await canvasToBlob(apiMask)
      const res = await aiStudioInpaint(imageBlob, maskBlob, prompt)
      const img = await loadImageFromUrl(res.result_url)
      resultCanvasRef.current = imageToCanvas(img, stageSize.w, stageSize.h)
      setResultKonvaImage(img)
      setHasResult(true)
      setShowResult(true)
      setProviderMsg(res.message ?? '')
      const sec = ((Date.now() - t0) / 1000).toFixed(0)
      message.success(`AI 消除完成 · ${res.provider}（${sec}s）`)
    } catch (e) {
      let msg = 'AI 消除失败'
      if (axios.isAxiosError(e)) {
        if (!e.response) msg = '无法连接后端，请在 backend 目录执行: go run .'
        else {
          const body = e.response.data as { message?: string }
          msg = body?.message || msg
          if (msg.includes('billing_hard_limit')) {
            msg = 'OpenAI 账户已达计费硬上限，请到 platform.openai.com → Billing 提高限额或充值后再试'
          }
        }
      }
      setProviderMsg(msg)
      message.error(msg, 8)
    } finally {
      setLoading(false)
    }
  }

  const handleExport = async () => {
    const orig = originalCanvasRef.current
    const result = resultCanvasRef.current
    const mask = maskCanvasRef.current
    if (!orig || !result || !mask) {
      message.warning('请先上传图片')
      return
    }
    const stamp = Date.now()
    try {
      if (exportFormat === 'psd') {
        if (!hasResult) {
          message.warning('建议先执行 AI 消除再导出 PSD（当前第 2 图层将使用原图副本）')
        }
        downloadPsd({
          width: stageSize.w,
          height: stageSize.h,
          original: orig,
          result,
          mask,
          filename: `imagedeal-studio-${stamp}.psd`,
        })
        message.success('PSD 已导出（3 图层：原图 / AI 结果 / 蒙版）')
        return
      }
      if (!hasResult) {
        message.warning('尚未执行 AI 消除，将导出当前画布（与原图相同）')
      }
      await downloadCanvasAsImage(result, exportFormat, {
        filename: `imagedeal-studio-${stamp}.${exportFormat === 'jpeg' ? 'jpg' : exportFormat}`,
      })
      const label = EXPORT_FORMAT_OPTIONS.find((o) => o.value === exportFormat)?.label ?? exportFormat
      message.success(`已导出 ${label}（AI 消除结果）`)
    } catch (e) {
      message.error(e instanceof Error ? e.message : '导出失败')
    }
  }

  const resetTransform = () => {
    const node = originalNodeRef.current
    if (!node) return
    node.scaleX(1)
    node.scaleY(1)
    node.rotation(0)
    node.x(0)
    node.y(0)
    node.getLayer()?.batchDraw()
  }

  return (
    <div style={{ margin: -24, minHeight: 'calc(100vh - 64px)', background: '#0f172a', color: '#e2e8f0' }}>
      <div
        style={{
          padding: '12px 20px',
          borderBottom: '1px solid #334155',
          display: 'flex',
          alignItems: 'center',
          gap: 12,
          flexWrap: 'wrap',
        }}
      >
        <Typography.Title level={4} style={{ color: '#fff', margin: 0 }}>
          GPT 图像消除 · Konva + PSD
        </Typography.Title>
        <Typography.Text style={{ color: studioReady ? '#22c55e' : '#f59e0b' }}>
          {studioReady ? `GPT Image · ${cfg?.model ?? 'gpt-image-2'}` : '未配置 API Key'}
        </Typography.Text>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: '240px 1fr 220px', minHeight: 560 }}>
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
              { label: '剪裁', value: 'crop' },
              { label: '移动', value: 'move' },
            ]}
          />
          <div style={{ marginTop: 16 }}>
            <Typography.Text style={{ color: '#94a3b8' }}>笔刷 {brushSize}px</Typography.Text>
            <Slider
              min={4}
              max={80}
              value={brushSize}
              onChange={setBrushSize}
              disabled={tool === 'move' || tool === 'crop'}
            />
          </div>
          <div style={{ marginTop: 8 }}>
            <Typography.Text style={{ color: '#94a3b8' }}>边缘羽化 {feather}px</Typography.Text>
            <Slider min={0} max={24} value={feather} onChange={setFeather} />
          </div>
          <Typography.Text style={{ color: '#64748b', fontSize: 12, display: 'block', marginTop: 12 }}>
            消除提示词（GPT Image）
          </Typography.Text>
          <Input.TextArea
            rows={3}
            value={prompt}
            onChange={(e) => setPrompt(e.target.value)}
            style={{ marginTop: 6, background: '#1e293b', color: '#e2e8f0' }}
          />
          <Space direction="vertical" style={{ width: '100%', marginTop: 12 }}>
            <Button block icon={<ClearOutlined />} onClick={clearMaskAction}>
              清除笔迹
            </Button>
            <Button
              block
              icon={<ScissorOutlined />}
              onClick={applyCrop}
              disabled={!konvaImage || tool !== 'crop'}
              type={tool === 'crop' ? 'primary' : 'default'}
            >
              应用剪裁
            </Button>
            <Button block icon={<UndoOutlined />} onClick={resetTransform} disabled={!konvaImage}>
              重置变换
            </Button>
            <Upload showUploadList={false} beforeUpload={handleUpload} accept="image/*">
              <Button block icon={<FileImageOutlined />}>
                上传图片
              </Button>
            </Upload>
            <Button block onClick={loadDemo}>
              加载演示图
            </Button>
          </Space>
        </div>

        <div style={{ padding: 16, display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 12 }}>
          {!studioReady && (
            <Alert
              type="warning"
              showIcon
              message="请配置 OpenAI Key"
              description="编辑 backend/conf/app.conf：ai_edit_mode=openai，openai_api_key=sk-...，openai_image_model=gpt-image-2，保存后重启 go run ."
              style={{ width: '100%', maxWidth: 960 }}
            />
          )}
          {studioReady && (
            <Alert
              type="info"
              showIcon
              message="剪裁后涂蒙版再消除更快。GPT Image 耗时主要在 OpenAI 云端（quality=medium，长边≤1024）。"
              style={{ width: '100%', maxWidth: 960 }}
            />
          )}
          {stageSize.w > 0 && konvaImage ? (
            <div
              style={{
                borderRadius: 8,
                overflow: 'hidden',
                boxShadow: '0 4px 24px rgba(0,0,0,.4)',
                cursor: tool === 'brush' || tool === 'eraser' ? 'crosshair' : 'default',
              }}
            >
              <Stage
                key={sessionKey}
                ref={stageRef}
                width={stageSize.w}
                height={stageSize.h}
              >
                <Layer listening={false}>
                  {showOriginal && (
                    <KonvaImage
                      ref={originalNodeRef}
                      image={konvaImage}
                      width={stageSize.w}
                      height={stageSize.h}
                      listening={tool === 'move'}
                      draggable={tool === 'move'}
                    />
                  )}
                  {showResult && resultKonvaImage && (
                    <KonvaImage
                      image={resultKonvaImage}
                      width={stageSize.w}
                      height={stageSize.h}
                    />
                  )}
                  {showMask && maskPreview && (
                    <KonvaImage
                      key={`mask-${sessionKey}-${maskPreview.width}x${maskPreview.height}`}
                      image={maskPreview}
                      width={stageSize.w}
                      height={stageSize.h}
                    />
                  )}
                </Layer>
                <Layer>
                  {(tool === 'brush' || tool === 'eraser') && (
                    <Rect
                      width={stageSize.w}
                      height={stageSize.h}
                      fill="rgba(0,0,0,0.001)"
                      onMouseDown={onPointerDown}
                      onMousemove={onPointerMove}
                      onMouseup={onPointerUp}
                      onMouseLeave={onPointerUp}
                    />
                  )}
                  {tool === 'crop' && stageSize.w > 0 && (
                    <>
                      <Rect
                        ref={cropRectRef}
                        x={0}
                        y={0}
                        width={stageSize.w}
                        height={stageSize.h}
                        stroke="#38bdf8"
                        strokeWidth={2}
                        dash={[10, 6]}
                        draggable
                      />
                      <Transformer
                        ref={cropTransformerRef}
                        rotateEnabled={false}
                        boundBoxFunc={(oldBox, newBox) => {
                          if (newBox.width < 24 || newBox.height < 24) return oldBox
                          return newBox
                        }}
                      />
                    </>
                  )}
                  {tool === 'move' && <Transformer ref={transformerRef} rotateEnabled bounded />}
                </Layer>
              </Stage>
            </div>
          ) : (
            <div
              style={{
                width: 480,
                height: 320,
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                background: '#1e293b',
                borderRadius: 8,
                color: '#64748b',
              }}
            >
              请上传图片开始编辑
            </div>
          )}
          {providerMsg && (
            <Typography.Text style={{ color: '#94a3b8', fontSize: 12 }}>{providerMsg}</Typography.Text>
          )}
        </div>

        <div style={{ padding: 16, borderLeft: '1px solid #334155' }}>
          <Typography.Text strong style={{ color: '#94a3b8' }}>
            图层
          </Typography.Text>
          <div style={{ marginTop: 12, fontSize: 13 }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', padding: '6px 0' }}>
              <span>原图</span>
              <Switch size="small" checked={showOriginal} onChange={setShowOriginal} />
            </div>
            <div style={{ display: 'flex', justifyContent: 'space-between', padding: '6px 0' }}>
              <span style={{ color: hasResult ? '#38bdf8' : '#64748b' }}>AI 消除结果</span>
              <Switch size="small" checked={showResult} onChange={setShowResult} disabled={!hasResult} />
            </div>
            <div style={{ display: 'flex', justifyContent: 'space-between', padding: '6px 0' }}>
              <span>蒙版预览</span>
              <Switch size="small" checked={showMask} onChange={setShowMask} />
            </div>
          </div>
          <Typography.Text type="secondary" style={{ fontSize: 11, color: '#64748b', display: 'block', marginTop: 8 }}>
            PNG/JPEG/WebP 导出 AI 处理结果；PSD 含原图、结果、蒙版三图层
          </Typography.Text>
          <Space direction="vertical" style={{ width: '100%', marginTop: 24 }}>
            <Button
              type="primary"
              block
              size="large"
              icon={<ThunderboltOutlined />}
              loading={loading}
              onClick={runInpaint}
              disabled={!konvaImage || !studioReady}
            >
              GPT 消除
            </Button>
            <Select
              value={exportFormat}
              onChange={setExportFormat}
              options={EXPORT_FORMAT_OPTIONS}
              style={{ width: '100%' }}
              disabled={!konvaImage}
            />
            <Button block icon={<DownloadOutlined />} onClick={() => void handleExport()} disabled={!konvaImage}>
              导出图片
            </Button>
          </Space>
        </div>
      </div>
    </div>
  )
}
