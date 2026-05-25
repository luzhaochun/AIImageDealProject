import { mediaUrl } from '@/utils/media'

export function ImageThumb({ path, alt }: { path?: string; alt?: string }) {
  const src = mediaUrl(path)
  if (!src) {
    return (
      <div
        style={{
          width: 56,
          height: 56,
          background: '#e2e8f0',
          borderRadius: 4,
        }}
      />
    )
  }
  return (
    <img
      src={src}
      alt={alt ?? ''}
      style={{ width: 56, height: 56, objectFit: 'cover', borderRadius: 4 }}
      onError={(e) => {
        ;(e.target as HTMLImageElement).style.display = 'none'
      }}
    />
  )
}
