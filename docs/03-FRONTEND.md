# 前端架构设计

---

## 1. 总体原则

- **SPA 单页应用**：React + TypeScript + Vite
- **按角色分 Feature**：admin / workspace / review 三大功能区
- **服务端状态与 UI 状态分离**：TanStack Query 管 API 数据，Zustand 管 token/用户信息
- **权限路由守卫**：未登录跳转登录；角色不匹配跳转 403
- **与后端契约对齐**：`types/` 目录镜像后端 domain 与 API DTO

---

## 2. 技术选型

| 类别 | 选型 | 理由 |
|------|------|------|
| 框架 | React 18 | 生态成熟，编辑器组件丰富 |
| 语言 | TypeScript | 与后端 DTO 类型安全对齐 |
| 构建 | Vite 5 | 快速冷启动与 HMR |
| 路由 | React Router 6 | 嵌套路由 + loader 可选 |
| UI | Ant Design 5 | 表格、表单、Upload、Steps 开箱即用 |
| HTTP | Axios | 拦截器统一注入 JWT |
| 服务端状态 | TanStack Query 5 | 缓存、重试、轮询导入进度 |
| 客户端状态 | Zustand | 轻量 auth store |
| 表单 | React Hook Form + Zod | 校验与类型推导 |
| 图片编辑 | Konva + react-konva | 二期；MVP 可先占位页 |
| 国际化 | 可选 i18next | 初期中文即可 |

---

## 3. 应用结构

```
frontend/
├── public/
├── src/
│   ├── main.tsx                 # 入口
│   ├── App.tsx                  # Provider 挂载
│   ├── api/                     # API 层
│   │   ├── client.ts            # Axios 实例、拦截器
│   │   ├── auth.ts
│   │   ├── imports.ts
│   │   ├── images.ts
│   │   ├── reviews.ts
│   │   └── archives.ts
│   ├── types/                   # 与后端对齐的类型
│   │   ├── image.ts
│   │   ├── user.ts
│   │   └── api.ts               # 统一响应包装
│   ├── stores/
│   │   └── authStore.ts         # token, user, login/logout
│   ├── routes/
│   │   ├── index.tsx            # 路由表
│   │   ├── ProtectedRoute.tsx   # 登录守卫
│   │   └── RoleRoute.tsx        # 角色守卫
│   ├── layouts/
│   │   ├── AdminLayout.tsx      # 侧边栏 + 顶栏
│   │   ├── WorkspaceLayout.tsx
│   │   └── ReviewLayout.tsx
│   ├── features/
│   │   ├── auth/
│   │   │   └── LoginPage.tsx
│   │   ├── admin/
│   │   │   ├── ImportPage.tsx       # CSV 上传
│   │   │   ├── ImportBatchList.tsx  # 批次列表与进度
│   │   │   └── ArchiveExportPage.tsx
│   │   ├── workspace/
│   │   │   ├── ClaimPage.tsx        # 选类目领图
│   │   │   ├── TaskListPage.tsx     # 我的任务区
│   │   │   ├── EditorPage.tsx       # 在线编辑器
│   │   │   └── ImageDetailPage.tsx
│   │   └── review/
│   │       ├── FirstReviewQueue.tsx
│   │       ├── SecondReviewQueue.tsx
│   │       └── ReviewComparePage.tsx # 原图 vs 成稿对比
│   ├── components/              # 跨 feature 复用
│   │   ├── ImageThumb.tsx
│   │   ├── StatusTag.tsx        # status 映射颜色/文案
│   │   ├── CategorySelect.tsx
│   │   └── PageHeader.tsx
│   ├── hooks/
│   │   ├── useAuth.ts
│   │   ├── useImages.ts         # TanStack Query 封装
│   │   └── useImportProgress.ts
│   └── utils/
│       ├── status.ts            # 状态文案映射
│       └── format.ts
├── index.html
├── vite.config.ts
├── tsconfig.json
├── .env.example                 # VITE_API_BASE_URL
└── package.json
```

---

## 4. 路由设计

| 路径 | 页面 | 角色 | 说明 |
|------|------|------|------|
| `/login` | LoginPage | 公开 | 登录 |
| `/admin/imports` | ImportPage | admin | CSV 导入 |
| `/admin/imports/:id` | ImportBatchDetail | admin | 批次详情/进度 |
| `/admin/archives` | ArchiveExportPage | admin | 终稿导出 |
| `/workspace/claim` | ClaimPage | user | 领图 |
| `/workspace/tasks` | TaskListPage | user | 任务列表 |
| `/workspace/editor/:id` | EditorPage | user | 编辑器 |
| `/review/first` | FirstReviewQueue | reviewer | 一审队列 |
| `/review/second` | SecondReviewQueue | reviewer | 二审队列 |
| `/review/:id` | ReviewComparePage | reviewer | 审核详情 |
| `/403` | Forbidden | — | 无权限 |
| `/` | — | — | 按角色 redirect |

**默认跳转逻辑：**
```
admin    → /admin/imports
user     → /workspace/tasks
reviewer → /review/first
```

---

## 5. 页面与交互设计

### 5.1 Admin — CSV 导入

**组件：** Ant Design `Upload.Dragger` + `Table` + `Progress`

**交互流程：**
1. 拖拽/选择 CSV 文件
2. 上传 → 显示 batch_id
3. TanStack Query 每 3s 轮询批次状态（processing → completed/failed）
4. 表格展示：总行数、成功数、失败数、失败原因下载

**状态展示：**
- processing：Spin + Progress
- completed：绿色 Tag
- failed：红色 Tag + 错误日志 Modal

### 5.2 User — 领图

**布局：**
- 左侧：类目树/Select（从 API 获取可用类目及库存数量）
- 右侧：领图按钮 + 最近领取记录

**交互：**
1. 选择 category
2. 点击「领取一张」→ POST claim
3. 成功：跳转 Editor 或 TaskList
4. 无库存：Toast 提示

### 5.3 User — 任务区

**表格列：** 缩略图 | global_no | category | status | 更新时间 | 操作

**操作按钮（随 status 变化）：**

| status | 操作 |
|--------|------|
| assigned | 开始编辑、换图 |
| in_progress | 继续编辑、提交审核 |
| rejected | 继续编辑（显示驳回意见） |
| pending_1st_review | 只读，等待审核 |

### 5.4 User — 编辑器（二期完整 / MVP 简化）

**MVP 方案：** 上传替换图 + 基础裁剪占位，先打通提审流程

**完整方案（Konva）：**
```
┌──────────────────────────────────────────────────┐
│ Toolbar: 裁剪 | 文字 | 调色 | 撤销 | 保存 | 提审  │
├────────────┬─────────────────────────────────────┤
│ Layer Panel│           Canvas (Konva Stage)       │
│  - 背景图   │                                     │
│  - 文字层   │                                     │
│  - 装饰层   │                                     │
└────────────┴─────────────────────────────────────┘
```

**自动保存：** debounce 2s → PUT version API

**提审：** Modal 确认 → POST submit → 跳转 TaskList

### 5.5 Reviewer — 审核对比

**布局：**
```
┌─────────────────┬─────────────────┐
│   原图 (只读)     │   成稿 (当前版本)  │
│   thumb/original │   thumb/version  │
└─────────────────┴─────────────────┘
┌───────────────────────────────────┐
│ 审核意见 TextArea                    │
│ [通过]  [驳回]                       │
└───────────────────────────────────┘
```

**一审队列：** 表格 + 「领取审核」或直接进入详情（按产品定）

**二审：** 从「一审通过库」抽样列表领取

---

## 6. 状态管理策略

### 6.1 Zustand — Auth Store

```typescript
interface AuthState {
  token: string | null;
  user: User | null;
  setAuth: (token: string, user: User) => void;
  logout: () => void;
}
```

持久化：`localStorage` 存 token（或 httpOnly cookie 方案需后端配合）

### 6.2 TanStack Query — 服务端数据

| Query Key | 用途 | staleTime |
|-----------|------|-----------|
| `['images', 'tasks', userId]` | 我的任务 | 30s |
| `['images', 'library', status]` | 图库列表 | 30s |
| `['imports', batchId]` | 导入进度 | 0（轮询） |
| `['reviews', 'queue', round]` | 审核队列 | 15s |
| `['categories']` | 类目列表 | 5min |

**Mutation 后 invalidate：**
- claim → invalidate tasks
- submit review → invalidate queue + tasks

### 6.3 编辑器本地状态

- 画布状态放 EditorPage 组件内或 `useReducer`
- 不放入全局 store（避免污染）

---

## 7. API 客户端设计

```typescript
// api/client.ts
const client = axios.create({ baseURL: import.meta.env.VITE_API_BASE_URL });

client.interceptors.request.use((config) => {
  const token = useAuthStore.getState().token;
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});

client.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401) {
      useAuthStore.getState().logout();
      window.location.href = '/login';
    }
    return Promise.reject(err);
  }
);
```

**统一响应类型：**
```typescript
interface ApiResponse<T> {
  data: T;
  request_id?: string;
}

interface ApiError {
  code: string;
  message: string;
}
```

---

## 8. 组件设计规范

| 规范 | 说明 |
|------|------|
| 命名 | 页面 `XxxPage`，容器 `XxxContainer`，展示 `XxxCard` |
| StatusTag | 统一映射 11 种 status → 颜色/中文 |
| 列表 | 统一分页：`page`, `page_size`，Ant Design Table |
| 空状态 | 每列表配 Empty + 引导操作 |
| 错误 | Query error → Result 组件 + 重试按钮 |

**StatusTag 映射示例：**

| status | 颜色 | 文案 |
|--------|------|------|
| pending_assign | default | 待分配 |
| assigned | processing | 已领取 |
| in_progress | processing | 作图中 |
| pending_1st_review | warning | 待一审 |
| rejected | error | 已驳回 |
| 1st_review_passed | success | 一审通过 |
| completed | success | 已归档 |
| discarded | default | 已废弃 |

---

## 9. 权限与安全

- 路由级：`RoleRoute requiredRole="admin"`
- 按钮级：根据 `user.role` 条件渲染
- 图片 URL：使用后端 presigned URL，不暴露 MinIO 内网地址
- XSS：图层 JSON 渲染时消毒文本层内容

---

## 10. 性能优化

| 场景 | 方案 |
|------|------|
| 列表缩略图 | 懒加载 + 小图 thumb |
| 编辑器 | 动态 import Konva（code splitting） |
| 大列表 | 虚拟滚动（rc-virtual-list） |
| 构建 | Vite manualChunks 分离 antd、konva |

---

## 11. 环境变量

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_APP_TITLE=ImageDeal
```

---

## 12. MVP 与完整版切分

| 阶段 | 前端交付 |
|------|----------|
| **MVP** | 登录、导入页、领图、任务列表、简化编辑（上传覆盖）、提审、一审页、状态 Tag |
| **V1.1** | Konva 完整编辑器、版本历史、撤销重做 |
| **V1.2** | 二审队列、导出页、WebSocket 待办通知 |

---

## 13. 与后端协作约定

1. 所有 datetime 使用 ISO 8601 UTC，前端按本地时区展示
2. UUID 作为主键，列表接口返回 string
3. 分页：`{ items, total, page, page_size }`
4. 枚举 status/role 前后端共用同一字符串值
5. OpenAPI/Swagger 由后端生成，前端可用 openapi-typescript 自动生成 types
