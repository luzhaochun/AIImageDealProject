import { useQueryClient } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '@/stores/authStore'

/** 清除本地登录态并跳转登录页 */
export function useLogout() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const logout = useAuthStore((s) => s.logout)

  return () => {
    logout()
    queryClient.clear()
    navigate('/login', { replace: true })
  }
}
