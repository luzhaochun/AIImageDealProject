export type AiStudioMode = 'erase' | 'upscale' | 'repair' | 'background' | 'outpaint' | 'watermark'

export interface AiStudioModeConfig {
  id: AiStudioMode
  label: string
  requires_mask: boolean
  mask_hint: string
  default_prompt: string
}

export const AI_STUDIO_MODES: AiStudioModeConfig[] = [
  {
    id: 'erase',
    label: '消除',
    requires_mask: true,
    mask_hint: '必须用红色画笔涂抹要删除的区域',
    default_prompt:
      '移除蒙版区域内的物体，用与周围环境一致的自然背景无缝填充，保持真实照片质感。',
  },
  {
    id: 'upscale',
    label: '高清增强',
    requires_mask: false,
    mask_hint: '整图增强；可选涂蒙版仅增强局部',
    default_prompt: '提升整张图片的清晰度与细节，减少模糊与噪点，保持色彩自然、无过度锐化。',
  },
  {
    id: 'repair',
    label: '修复',
    requires_mask: false,
    mask_hint: '可整图修复；涂蒙版则只修复局部（划痕、破损、噪点）',
    default_prompt: '修复图片中的瑕疵、划痕、破损与噪点，使画面干净自然，保持原有构图与风格。',
  },
  {
    id: 'background',
    label: '换背景',
    requires_mask: false,
    mask_hint: '可整图换背景；涂蒙版可保留主体仅替换背景区域',
    default_prompt: '替换为简洁自然的背景，主体边缘清晰，光照与透视一致，整体像真实拍摄。',
  },
  {
    id: 'outpaint',
    label: '扩图',
    requires_mask: false,
    mask_hint: '整图扩展画布；涂蒙版可指定扩展方向或区域',
    default_prompt: '在保持原图主体不变的前提下自然扩展画面边缘，延伸场景与环境，透视与风格一致。',
  },
  {
    id: 'watermark',
    label: '去水印',
    requires_mask: false,
    mask_hint: '涂蒙版标出水印位置效果更好；也可整图自动查找去除',
    default_prompt: '去除图片中的水印、文字叠加与 logo，用周围纹理自然填充，不留明显修补痕迹。',
  },
]

export function getStudioModeConfig(id: AiStudioMode): AiStudioModeConfig {
  return AI_STUDIO_MODES.find((m) => m.id === id) ?? AI_STUDIO_MODES[0]
}
