# ImageDeal — SaaS 图片协作平台

基于业务流程图（13 步）实现的图片导入、领图作图、一审/二审、归档全流程。

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

## 业务流程 ↔ 实现

| 步骤 | 实现要点 |
|------|----------|
| 导入 | `POST /api/v1/imports` + Redis `parse_csv` |
| 入库/下载 | Worker `download` → 本地 `storage/originals/` |
| 分类 | CSV `category` 写入 `images.category` |
| 领图 | `POST /images/claim`（MySQL `FOR UPDATE SKIP LOCKED`） |
| 换图 | `POST /images/:id/discard` + Redis 每日限额 |
| 作图 | `PUT /images/:id/versions` |
| 提审 | `POST /images/:id/submit` |
| 一审/二审 | `POST /images/:id/reviews`，二审领图 `POST /reviews/claim-second` |
| 归档 | `GET /archives`，`POST /archives/:id/complete` |

状态机详见 [docs/06-STATE-MACHINE.md](docs/06-STATE-MACHINE.md)，API 详见 [docs/05-API.md](docs/05-API.md)。

## UI 设计稿

```bash
open design/index.html
```

## 技术栈

- **后端**: Go, Beego v2, Beego ORM, MySQL 8, Redis（队列/计数）, JWT
- **前端**: React, TypeScript, Vite, Ant Design, TanStack Query, Zustand
- **存储**: 本地 `backend/storage/`（MVP；MinIO 容器已预留）
