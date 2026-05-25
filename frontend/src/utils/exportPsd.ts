import { writePsd } from 'ag-psd'
import { maskToGrayscaleCanvas } from './maskCanvas'

export interface PsdExportInput {
  width: number
  height: number
  original: HTMLCanvasElement
  result: HTMLCanvasElement
  mask: HTMLCanvasElement
  filename?: string
}

/** 导出标准 PSD：原图 / AI 结果 / 灰度蒙版 */
export function buildPsdBuffer(input: PsdExportInput): ArrayBuffer {
  const maskGray = maskToGrayscaleCanvas(input.mask)
  return writePsd({
    width: input.width,
    height: input.height,
    children: [
      { name: '原始图片', canvas: input.original },
      { name: 'AI 消除结果', canvas: input.result },
      { name: '蒙版', canvas: maskGray },
    ],
  })
}

export function downloadPsd(input: PsdExportInput) {
  const buffer = buildPsdBuffer(input)
  const blob = new Blob([buffer], { type: 'application/octet-stream' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = input.filename ?? `imagedeal-${Date.now()}.psd`
  a.click()
  URL.revokeObjectURL(url)
}
