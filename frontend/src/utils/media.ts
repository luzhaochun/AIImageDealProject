/** 原图/缩略图静态路径（API 映射 /static → storage） */
export function mediaUrl(path?: string | null): string {
  if (!path) return ''
  if (path.startsWith('http://') || path.startsWith('https://')) return path
  return `/static/${path.replace(/^\//, '')}`
}
