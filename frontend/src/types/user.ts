export type UserRole = 'admin' | 'user' | 'reviewer'

export interface User {
  id: string
  email: string
  display_name?: string
  role: UserRole
  tenant_id: string
}
