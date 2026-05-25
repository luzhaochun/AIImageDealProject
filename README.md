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
| admin@example.com | admin | `/admin` CSV 导入、AI 消除 Demo |
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

## AI 消除 Demo（Admin 面试演示）

- 入口：`/admin/ai-editor`（仅 admin）
- 设计稿：`design/pages/admin-ai-editor.html`
- 说明：[docs/09-AI-EDITOR-DEMO.md](docs/09-AI-EDITOR-DEMO.md)
- 默认 **mock 模式** 无需 API Key；配置 `openai_api_key` 后可切换 GPT Image

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
