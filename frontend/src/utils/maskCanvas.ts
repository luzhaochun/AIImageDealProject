/** 蒙版画布：笔刷绘制、羽化、导出灰度与 API 用 PNG */

export function createMaskCanvas(w: number, h: number): HTMLCanvasElement {
  const c = document.createElement('canvas')
  c.width = w
  c.height = h
  return c
}

export function maskHasPixels(canvas: HTMLCanvasElement): boolean {
  const data = canvas.getContext('2d')!.getImageData(0, 0, canvas.width, canvas.height).data
  for (let i = 3; i < data.length; i += 4) {
    if (data[i] > 20) return true
  }
  return false
}

export function clearMask(canvas: HTMLCanvasElement) {
  const ctx = canvas.getContext('2d')!
  ctx.clearRect(0, 0, canvas.width, canvas.height)
}

/** 红色半透明预览层（Konva 叠加用） */
export function buildMaskPreview(mask: HTMLCanvasElement): HTMLCanvasElement {
  const w = mask.width
  const h = mask.height
  const out = document.createElement('canvas')
  out.width = w
  out.height = h
  const src = mask.getContext('2d')!.getImageData(0, 0, w, h)
  const dst = out.getContext('2d')!.createImageData(w, h)
  for (let i = 0; i < src.data.length; i += 4) {
    const a = src.data[i + 3]
    if (a > 20) {
      dst.data[i] = 255
      dst.data[i + 1] = 60
      dst.data[i + 2] = 60
      dst.data[i + 3] = Math.min(200, a)
    }
  }
  out.getContext('2d')!.putImageData(dst, 0, 0)
  return out
}

/** 黑白灰度蒙版（PSD 第 3 图层） */
export function maskToGrayscaleCanvas(mask: HTMLCanvasElement): HTMLCanvasElement {
  const w = mask.width
  const h = mask.height
  const out = document.createElement('canvas')
  out.width = w
  out.height = h
  const src = mask.getContext('2d')!.getImageData(0, 0, w, h)
  const dst = out.getContext('2d')!.createImageData(w, h)
  for (let i = 0; i < src.data.length; i += 4) {
    const a = src.data[i + 3]
    const r = src.data[i]
    const v = a > 20 || r > 180 ? 255 : 0
    dst.data[i] = v
    dst.data[i + 1] = v
    dst.data[i + 2] = v
    dst.data[i + 3] = 255
  }
  out.getContext('2d')!.putImageData(dst, 0, 0)
  return out
}

/** 羽化后副本（送 AI API） */
export function featherMaskCopy(mask: HTMLCanvasElement, radius: number): HTMLCanvasElement {
  if (radius <= 0) {
    const c = createMaskCanvas(mask.width, mask.height)
    c.getContext('2d')!.drawImage(mask, 0, 0)
    return c
  }
  const c = createMaskCanvas(mask.width, mask.height)
  const ctx = c.getContext('2d')!
  ctx.filter = `blur(${radius}px)`
  ctx.drawImage(mask, 0, 0)
  ctx.filter = 'none'
  return c
}

export function canvasToBlob(
  canvas: HTMLCanvasElement,
  type = 'image/png',
  quality?: number,
): Promise<Blob> {
  return new Promise((resolve, reject) => {
    canvas.toBlob(
      (b) => (b ? resolve(b) : reject(new Error('canvas export failed'))),
      type,
      quality,
    )
  })
}

export function imageToCanvas(img: HTMLImageElement, w: number, h: number): HTMLCanvasElement {
  const c = document.createElement('canvas')
  c.width = w
  c.height = h
  c.getContext('2d')!.drawImage(img, 0, 0, w, h)
  return c
}

export async function loadImageFromFile(file: File): Promise<HTMLImageElement> {
  const url = URL.createObjectURL(file)
  try {
    return await loadImageFromUrl(url)
  } finally {
    URL.revokeObjectURL(url)
  }
}

/** 裁剪画布区域（坐标为画布像素） */
export function cropCanvas(
  src: HTMLCanvasElement,
  x: number,
  y: number,
  w: number,
  h: number,
): HTMLCanvasElement {
  const ix = Math.max(0, Math.round(x))
  const iy = Math.max(0, Math.round(y))
  const iw = Math.min(Math.round(w), src.width - ix)
  const ih = Math.min(Math.round(h), src.height - iy)
  if (iw < 1 || ih < 1) {
    throw new Error('裁剪区域过小')
  }
  const out = createMaskCanvas(iw, ih)
  out.getContext('2d')!.drawImage(src, ix, iy, iw, ih, 0, 0, iw, ih)
  return out
}

/** 送 API 前缩小长边，加快 GPT Image 推理 */
export const API_MAX_LONG_EDGE = 1024

export function downscaleCanvasForApi(
  canvas: HTMLCanvasElement,
  maxLong: number = API_MAX_LONG_EDGE,
): HTMLCanvasElement {
  const w = canvas.width
  const h = canvas.height
  const long = Math.max(w, h)
  if (long <= maxLong) return canvas
  const scale = maxLong / long
  const nw = Math.round(w * scale)
  const nh = Math.round(h * scale)
  const out = createMaskCanvas(nw, nh)
  out.getContext('2d')!.drawImage(canvas, 0, 0, w, h, 0, 0, nw, nh)
  return out
}

export function canvasToImage(canvas: HTMLCanvasElement): Promise<HTMLImageElement> {
  return loadImageFromUrl(canvas.toDataURL('image/png'))
}

export function loadImageFromUrl(url: string): Promise<HTMLImageElement> {
  return new Promise((resolve, reject) => {
    const img = new Image()
    img.crossOrigin = 'anonymous'
    img.onload = () => resolve(img)
    img.onerror = () => reject(new Error('图片加载失败'))
    img.src = url
  })
}
