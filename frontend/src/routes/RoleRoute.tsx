import { Navigate } from 'react-router-dom'
import { useAuthStore } from '@/stores/authStore'
import type { UserRole } from '@/types/user'

export function RoleRoute({
  role,
  children,
}: {
  role: UserRole
  children: React.ReactNode
}) {
  const user = useAuthStore((s) => s.user)
  if (!user) return <Navigate to="/login" replace />
  if (user.role !== role) return <Navigate to="/403" replace />
  return <>{children}</>
}
