import { canvasToBlob } from './maskCanvas'

export type RasterImageFormat = 'png' | 'jpeg' | 'webp'

export type StudioExportFormat = RasterImageFormat | 'psd'

const MIME: Record<RasterImageFormat, string> = {
  png: 'image/png',
  jpeg: 'image/jpeg',
  webp: 'image/webp',
}

const EXT: Record<RasterImageFormat, string> = {
  png: 'png',
  jpeg: 'jpg',
  webp: 'webp',
}

const DEFAULT_QUALITY: Record<'jpeg' | 'webp', number> = {
  jpeg: 0.92,
  webp: 0.9,
}

function triggerDownload(blob: Blob, filename: string) {
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  a.click()
  URL.revokeObjectURL(url)
}

/** JPEG 不支持透明通道，先铺白底再导出 */
function canvasForJpeg(source: HTMLCanvasElement): HTMLCanvasElement {
  const c = document.createElement('canvas')
  c.width = source.width
  c.height = source.height
  const ctx = c.getContext('2d')!
  ctx.fillStyle = '#ffffff'
  ctx.fillRect(0, 0, c.width, c.height)
  ctx.drawImage(source, 0, 0)
  return c
}

export async function downloadCanvasAsImage(
  canvas: HTMLCanvasElement,
  format: RasterImageFormat,
  options?: { filename?: string; quality?: number },
) {
  const src = format === 'jpeg' ? canvasForJpeg(canvas) : canvas
  const mime = MIME[format]
  const quality =
    format === 'png' ? undefined : (options?.quality ?? DEFAULT_QUALITY[format])
  const blob = await canvasToBlob(src, mime, quality)
  const ext = EXT[format]
  const filename = options?.filename ?? `imagedeal-${Date.now()}.${ext}`
  triggerDownload(blob, filename)
}

export const EXPORT_FORMAT_OPTIONS: { value: StudioExportFormat; label: string }[] = [
  { value: 'png', label: 'PNG' },
  { value: 'jpeg', label: 'JPEG' },
  { value: 'webp', label: 'WebP' },
  { value: 'psd', label: 'PSD（多图层）' },
]
