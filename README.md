# ImageDeal — SaaS 图片协作平台

SaaS 图片协作平台：从 CSV 批量导入外链图片，经入库分类、用户领图作图、一审/二审质检，到终稿归档导出的全流程系统。核心驱动为 **图片状态机** + **逻辑图库（虚拟队列）** + **多角色 RBAC**。

---

## 项目需求

### 业务背景

平台解决 **大批量外链图片素材** 在团队内的协作生产问题：

- 管理员通过 CSV 批量导入图片 URL 与分类
- 系统自动下载、校验、生成缩略图并分类入库
- 用户按类目领图、在线编辑、提交审核
- 审核员一审（可选二审）进行质量控制
- 通过后进入终稿库，支持导出打包

### 业务流程总图（13 步）

![SaaS 图片协作平台 — 业务全流程](docs/images/business-workflow.jpg)

### 13 步需求明细

| 步骤 | 操作人 | 操作 | 系统行为 | 数据/状态变更 | 图库归属 |
|:---:|--------|------|----------|---------------|----------|
| **1 导入** | 管理员 | 上传 CSV | 解析 CSV，读取图片下载链接，批量下载原图 | `image_urls` 写入 `images` 表 | — |
| **2 入库** | 系统 | 自动下载 | 逐张下载、校验格式/尺寸、生成缩略图 | 插入 N 条 `images` 记录 | 原始图库（待分配） |
| **3 分类** | 系统 | 按 CSV 分类 | 将 CSV 类目写入图片元数据 | 填充 `images.category` | 原始图库（按类别分好） |
| **4 领图** | 用户 | 选类目领取 1 张 | 从该类目待分配库取 1 张，分配全局唯一编号 | `status=已分配`，`assigned_to=用户`，生成 `global_no` | 原始图库 → 已领出 |
| **5 换图** | 用户 | 评估后不可用 | 当前图入废料堆，用户重新领图 | 原图 `status=废弃`；重复步骤 4 | 已领出 → 废图堆 |
| **6 作图** | 用户 | 在线编辑 | 图层/裁剪/调色/文字等，保存进度 | 创建 `image_versions` + `image_layers` | 工作中（用户任务区） |
| **7 提审** | 用户 | 编辑完成提交一审 | 图片进入一审队列 | `status=一审待审` | 用户任务区 → 一审待审队列 |
| **8 一审** | 审核员 | 审核 | 对比原图，通过/驳回/打回修改 | 写入 `review_records` | 一审队列 |
| **9 一审通过** | 系统 | 审核通过 | 进入「一审通过图库」 | `status=一审通过` | 一审通过图库 |
| **10 一审驳回** | 审核员 | 审核不通过 | 退回用户并附修改意见 | `status=驳回`；用户重新编辑 | 退回用户任务区 |
| **11 二审（可选）** | 审核员 | 从一审通过库领取二审 | 抽样或人工领取 | `status=二审待审` | 一审通过图库 → 二审队列 |
| **12 二审通过** | 系统 | 二审通过 | 进入「二审通过图库（终稿）」 | `status=二审通过` | 二审通过图库（终稿） |
| **13 终稿归档** | 系统 | 终稿标记 | 已完成图可导出/打包 | `status=已完成` | 终稿归档库 |

### 角色与权限

| 功能 | Admin | User | Reviewer | System |
|------|:-----:|:----:|:--------:|:------:|
| CSV 导入 / 查看批次 | ✅ | — | — | — |
| 按类目领图 / 换图 | — | ✅ | — | — |
| 在线编辑 / 保存版本 / 提审 | — | ✅ | — | — |
| 一审 / 二审审核 | — | — | ✅ | — |
| 终稿导出 | ✅ | — | — | — |
| 批量下载 / 缩略图 | — | — | — | ✅ |
| GPT 图像消除（Canvas 工作室） | ✅ | — | — | — |

测试账号见下文「测试账号」；JWT 携带 `role` 做路由守卫。

### 逻辑图库（虚拟队列）

图库由 **status + 角色 + 过滤条件** 定义，非物理文件夹：

| 图库名称 | 过滤条件 | 可见角色 |
|----------|----------|----------|
| 原图库（待分配） | `pending_assign` | User |
| 原图库（已分类） | `pending_assign` + `category` | User |
| 已领取 / 用户任务区 | `assigned_to = me`，status 含 assigned / in_progress / rejected | User |
| 废料堆 | `discarded` | Admin |
| 一审队列 | `pending_1st_review` | Reviewer |
| 一审通过库 | `1st_review_passed` | Reviewer、Admin |
| 二审队列 | `pending_2nd_review` | Reviewer |
| 终稿库 | `2nd_review_passed` / `completed` | Admin |

### 功能模块拆分

| 模块 | 路由前缀 | 主要页面 | 对应需求步骤 |
|------|----------|----------|--------------|
| **管理端 Admin** | `/admin` | CSV 导入、批次详情、终稿导出、AI Demo、Canvas 工作室 | 1、13；扩展 AI 消除 |
| **工作台 Workspace** | `/workspace` | 领图、我的任务、在线编辑器 | 4–7 |
| **审核端 Review** | `/review` | 一审队列、二审队列、审核对比 | 8–12 |
| **认证** | `/login` | 登录、403 | 全角色入口 |

**前端页面 ↔ 设计稿：**

| 页面 | 实现路由 | 设计稿 |
|------|----------|--------|
| CSV 导入 | `/admin/imports` | `design/pages/admin-import.html` |
| 导入批次详情 | `/admin/imports/:id` | `design/pages/admin-import-detail.html` |
| 终稿导出 | `/admin/archives` | `design/pages/admin-archive.html` |
| 领图 | `/workspace/claim` | `design/pages/workspace-claim.html` |
| 我的任务 | `/workspace/tasks` | `design/pages/workspace-tasks.html` |
| 在线编辑器 | `/workspace/editor/:id` | `design/pages/workspace-editor.html` |
| 一审 / 二审 | `/review/first`、`/review/second` | `design/pages/review-first.html` 等 |
| 审核对比 | `/review/:id` | `design/pages/review-compare.html` |
| GPT Canvas 工作室 | `/admin/canvas-studio` | React + Konva（见下文专节） |

### 状态机（实现对照）

| 状态值 | 中文 | 所属图库 |
|--------|------|----------|
| `pending_storage` | 待入库 | 系统处理中 |
| `pending_assign` | 待分配 | 原图库 |
| `assigned` / `in_progress` | 已领取 / 作图中 | 用户任务区 |
| `discarded` | 已废弃 | 废料堆 |
| `pending_1st_review` | 待一审 | 一审队列 |
| `rejected` | 一审驳回 | 用户任务区 |
| `1st_review_passed` | 一审通过 | 一审通过库 |
| `pending_2nd_review` | 待二审 | 二审队列 |
| `2nd_review_passed` | 二审通过 | 终稿库 |
| `completed` | 已归档 | 终稿归档库 |

状态流转与合法转换详见 [docs/06-STATE-MACHINE.md](docs/06-STATE-MACHINE.md)。

### 业务规则（默认）

| 规则 | 默认值 |
|------|--------|
| 用户同时进行中任务上限 | 5 |
| 每日换图上限 | 10（Redis 计数） |
| 领图并发 | `FOR UPDATE SKIP LOCKED`，一张图仅一人可领 |
| 二审模式 | 人工从一审通过库领取（MVP） |
| `global_no` 生成时机 | 领图成功时 |

### 扩展需求：GPT 图像消除

管理端 **Canvas 工作室**（`/admin/canvas-studio`）：Konva 蒙版编辑 + OpenAI GPT Image 消除 + PSD 三图层导出。用于素材预处理或演示，不纳入 13 步主流程状态机。详见下文「GPT 图像消除 · Canvas 工作室」。

### 需求 ↔ 实现索引

| 步骤 | API / 组件要点 |
|------|----------------|
| 1 导入 | `POST /api/v1/imports` + Redis `parse_csv` |
| 2–3 入库/分类 | Worker `download` → `storage/originals/`；`images.category` |
| 4 领图 | `POST /images/claim` |
| 5 换图 | `POST /images/:id/discard` |
| 6 作图 | `PUT /images/:id/versions` |
| 7 提审 | `POST /images/:id/submit` |
| 8–12 审核 | `POST /images/:id/reviews`；`POST /reviews/claim-second` |
| 13 归档 | `GET /archives`；`POST /archives/:id/complete` |

完整 API：[docs/05-API.md](docs/05-API.md) · 架构总览：[docs/01-OVERVIEW.md](docs/01-OVERVIEW.md) · 文档索引：[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)

---

## 项目结构

```
ImageDealProject/
├── backend/          # Go API (Beego v2) + Worker
├── frontend/         # React + Vite + Ant Design
├── design/           # 静态 HTML 设计稿
├── docs/             # 架构文档
└── docker-compose.yml
```

## 快速开始

### 1. 基础设施

```bash
# Redis + MinIO（MySQL 使用本机）
docker compose up -d redis minio
```

**MySQL 8**（本机）：

| 项 | 值 |
|----|-----|
| 地址 | `127.0.0.1:3306` |
| 用户/密码 | `root` / `123456` |
| 库名 | `imagedeal` |

```sql
CREATE DATABASE IF NOT EXISTS imagedeal CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

执行迁移：

```bash
mysql -h127.0.0.1 -uroot -p123456 imagedeal < backend/sql/migrations/001_init.sql
```

### 2. 后端 API

```bash
cd backend
go mod tidy
go run .
```

- API: `http://localhost:8080`
- 健康检查: `GET /api/v1/health`
- 静态原图: `/static/originals/...`

### 3. Worker（必须，否则 CSV 不会解析/下载）

```bash
cd backend
go run ./worker
```

Redis 队列：`parse_csv`、`download`（`LPUSH` / `BRPOP`）。

### 4. 前端

```bash
cd frontend
npm install
npm run dev
```

- 前端: `http://localhost:5173`（`/api` 代理到 8080）

## 测试账号

密码均为 **`password123`**：

| 邮箱 | 角色 | 入口 |
|------|------|------|
| admin@example.com | admin | `/admin` CSV 导入、Canvas 工作室、AI 消除 Demo |
| user@example.com | user | `/workspace` 领图/作图 |
| reviewer@example.com | reviewer | `/review` 一审/二审 |

## 测试数据（CSV）

目录：`backend/sample/`

| 文件 | 说明 |
|------|------|
| `import_minimal.csv` | 1 条，快速冒烟 |
| `import_sample.csv` | 3 条，默认示例 |
| `import_batch_medium.csv` | 15 条，多类目（banner/poster/icon/social） |
| `import_url_alias.csv` | 表头为 `url` 的兼容测试 |

详见 [backend/sample/README.md](backend/sample/README.md)。列：`image_url`（或 `url`）+ `category`。上传后需 Worker 下载完成才可领图。

<a id="canvas-studio"></a>

## GPT 图像消除 · Canvas 工作室

管理端单页应用，路由：**`/admin/canvas-studio`**（登录 `admin@example.com` 后，侧栏「GPT 图像消除」）。

> 说明：无独立 `canvas-studio.design` 路由；静态设计稿见 `design/pages/admin-ai-editor.html`（旧 mock Demo），Canvas 工作室为 React + Konva 实现。详细技术说明见 [docs/10-CANVAS-STUDIO.md](docs/10-CANVAS-STUDIO.md)。

### 访问

| 项 | 值 |
|----|-----|
| 本地地址 | http://localhost:5173/admin/canvas-studio |
| 后端接口 | `POST /api/v1/admin/ai/studio-inpaint` |
| 实现文件 | `frontend/src/features/admin/konva/KonvaStudioPage.tsx` |

### 界面截图

完整页面（左工具栏 · 中央画布 · 右图层与导出）：

![Canvas 工作室界面](docs/images/canvas-studio/ui-screenshot.png)

### 页面布局（三栏）

```
┌─────────────────────────────────────────────────────────────┐
│ 顶栏：GPT 图像消除 · Konva + PSD  │  GPT Image · gpt-image-2 │
├──────────┬──────────────────────────────┬───────────────────┤
│ 左：工具  │ 中：Konva 画布                │ 右：图层 + 导出    │
│ 画笔/橡皮 │ 原图 + 蒙版 + AI 结果叠层      │ 原图 / 结果 / 蒙版 │
│ 剪裁/移动 │ 剪裁框 / 变换手柄              │ GPT 消除 / 导出    │
│ 笔刷/羽化 │                              │ PNG·JPEG·WebP·PSD │
│ 提示词    │                              │                   │
│ 上传/演示 │                              │                   │
└──────────┴──────────────────────────────┴───────────────────┘
```

### 功能状态一览

Canvas 工作室为**单路由多状态**（非多子页面），下表对应界面上的各工作阶段：

| 状态 | 触发条件 | 界面表现 |
|------|----------|----------|
| **空画布** | 未上传图片 | 中央占位「请上传图片开始编辑」 |
| **未配置 Key** | `openai_api_key` 为空 | 顶栏黄色「未配置 API Key」+ 黄色 Alert |
| **已就绪** | Key 已配置且后端已重启 | 顶栏绿色 `GPT Image · gpt-image-2` + 蓝色提示 |
| **编辑中** | 已上传原图 | 画布显示原图；可用画笔涂红色蒙版 |
| **剪裁模式** | 工具选「剪裁」 | 蓝色虚线框 + 拖拽缩放；点「应用剪裁」同步裁剪原图/蒙版/结果 |
| **移动模式** | 工具选「移动」 | 原图可拖拽、旋转、缩放；「重置变换」恢复 |
| **AI 处理中** | 点击「GPT 消除」 | 按钮 loading，底部显示耗时提示（约 20–90 秒） |
| **有结果** | 消除成功 | 「AI 消除结果」图层可开关；可导出 PNG/JPEG/WebP/PSD |
| **导出 PSD** | 格式选 PSD | 三图层：原图 / AI 结果 / 蒙版（ag-psd） |

### 工具与图层

| 区域 | 功能 |
|------|------|
| **画笔** | 红色涂抹 = 要消除的区域 |
| **橡皮** | 擦除蒙版笔迹 |
| **剪裁** | 框选区域 →「应用剪裁」 |
| **移动** | 调整原图位置/旋转/缩放 |
| **边缘羽化** | 送 API 前柔化蒙版边缘（0–24px） |
| **消除提示词** | 传给 GPT Image 的文本指令 |
| **图层开关** | 原图 / AI 消除结果 / 蒙版预览 独立显示 |
| **导出** | 结果图：PNG、JPEG、WebP；工程：PSD 三图层 |

### 推荐操作流程

1. 用 **admin** 登录 → 打开 http://localhost:5173/admin/canvas-studio  
2. 配置 `backend/conf/app.conf`（见 [docs/10-CANVAS-STUDIO.md](docs/10-CANVAS-STUDIO.md)）并重启后端  
3. **上传图片** 或 **加载演示图**（`frontend/public/demo-food.jpg`）  
4. （可选）**剪裁** 缩小编辑区域，加快推理  
5. **画笔** 涂抹要消除的物体 → 调整 **边缘羽化** / **提示词**  
6. 点击 **GPT 消除**，等待完成  
7. 切换 **图层** 对比原图与结果 → **导出** PNG 或 PSD  

### Demo 效果（GPT Image 实机消除）

#### 示例 1：沙发看球 — 消除人物

| 原图（3 人） | 消除 1 人（蒙版涂左侧） | 消除全部人物 |
|:---:|:---:|:---:|
| ![原图](docs/images/canvas-studio/demo-sofa-before.png) | ![消除一人](docs/images/canvas-studio/demo-sofa-remove-one.png) | ![空沙发](docs/images/canvas-studio/demo-sofa-after.png) |

同一素材另一角度消除 1 人：

![消除一人（侧视）](docs/images/canvas-studio/demo-sofa-remove-one-alt.png)

#### 示例 2：牧场 — 消除背景牛群

| 原图（多头牛） | AI 消除后（仅保留前景牛） |
|:---:|:---:|
| ![牧场原图](docs/images/canvas-studio/demo-cows-before.png) | ![消除后](docs/images/canvas-studio/demo-cows-after.png) |

> 操作要点：蒙版只涂要去掉的区域；区域越小、图越短边，OpenAI 推理越快。默认 `quality=medium`、长边 ≤ 1024px。

### 配置速查

```ini
# backend/conf/app.conf（勿提交到 git）
ai_edit_mode = openai
openai_api_key = sk-your-key
openai_image_model = gpt-image-2
openai_image_quality = medium
openai_image_max_long = 1024
```

---

## AI 消除 Demo（旧版 Mock，Admin 面试演示）

- 入口：`/admin/ai-editor`（仅 admin）
- 设计稿：`design/pages/admin-ai-editor.html`
- 说明：[docs/09-AI-EDITOR-DEMO.md](docs/09-AI-EDITOR-DEMO.md)
- 默认 **mock 模式** 无需 API Key；配置 `openai_api_key` 后 mock 与 Canvas 工作室共用同一 Key

## UI 设计稿

```bash
open design/index.html
```

## 技术栈

- **后端**: Go, Beego v2, Beego ORM, MySQL 8, Redis（队列/计数）, JWT
- **前端**: React, TypeScript, Vite, Ant Design, TanStack Query, Zustand
- **存储**: 本地 `backend/storage/`（MVP；MinIO 容器已预留）
