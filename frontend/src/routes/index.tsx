import { Navigate, Route, Routes } from 'react-router-dom'
import { AdminLayout } from '@/layouts/AdminLayout'
import { WorkspaceLayout } from '@/layouts/WorkspaceLayout'
import { ReviewLayout } from '@/layouts/ReviewLayout'
import { ProtectedRoute } from './ProtectedRoute'
import { RoleRoute } from './RoleRoute'
import { LoginPage } from '@/features/auth/LoginPage'
import { HomePage } from '@/features/home/HomePage'
import { ForbiddenPage } from '@/features/common/ForbiddenPage'
import { ImportPage } from '@/features/admin/ImportPage'
import { ImportDetailPage } from '@/features/admin/ImportDetailPage'
import { AiEditorPage } from '@/features/admin/AiEditorPage'
import { KonvaStudioPage } from '@/features/admin/konva/KonvaStudioPage'
import { ArchiveExportPage } from '@/features/admin/ArchiveExportPage'
import { ClaimPage } from '@/features/workspace/ClaimPage'
import { TaskListPage } from '@/features/workspace/TaskListPage'
import { EditorPage } from '@/features/workspace/EditorPage'
import { FirstReviewQueue } from '@/features/review/FirstReviewQueue'
import { SecondReviewQueue } from '@/features/review/SecondReviewQueue'
import { ReviewComparePage } from '@/features/review/ReviewComparePage'

export function AppRoutes() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/403" element={<ForbiddenPage />} />
      <Route
        path="/"
        element={
          <ProtectedRoute>
            <HomePage />
          </ProtectedRoute>
        }
      />
      <Route
        path="/admin"
        element={
          <ProtectedRoute>
            <RoleRoute role="admin">
              <AdminLayout />
            </RoleRoute>
          </ProtectedRoute>
        }
      >
        <Route index element={<Navigate to="imports" replace />} />
        <Route path="imports" element={<ImportPage />} />
        <Route path="imports/:id" element={<ImportDetailPage />} />
        <Route path="archives" element={<ArchiveExportPage />} />
        <Route path="ai-editor" element={<AiEditorPage />} />
        <Route path="canvas-studio" element={<KonvaStudioPage />} />
      </Route>
      <Route
        path="/workspace"
        element={
          <ProtectedRoute>
            <RoleRoute role="user">
              <WorkspaceLayout />
            </RoleRoute>
          </ProtectedRoute>
        }
      >
        <Route index element={<Navigate to="tasks" replace />} />
        <Route path="claim" element={<ClaimPage />} />
        <Route path="tasks" element={<TaskListPage />} />
        <Route path="editor/:id" element={<EditorPage />} />
      </Route>
      <Route
        path="/review"
        element={
          <ProtectedRoute>
            <RoleRoute role="reviewer">
              <ReviewLayout />
            </RoleRoute>
          </ProtectedRoute>
        }
      >
        <Route index element={<Navigate to="first" replace />} />
        <Route path="first" element={<FirstReviewQueue />} />
        <Route path="second" element={<SecondReviewQueue />} />
        <Route path=":id" element={<ReviewComparePage />} />
      </Route>
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}
