# REST API 设计

> Base URL: `/api/v1`  
> 认证: `Authorization: Bearer <JWT>`（登录接口除外）  
> Content-Type: `application/json`（文件上传除外）

---

## 1. 通用约定

### 1.1 成功响应

```json
{
  "data": { ... },
  "request_id": "req-abc123"
}
```

### 1.2 分页响应

```json
{
  "data": {
    "items": [ ... ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

### 1.3 错误响应

```json
{
  "code": "INVALID_TRANSITION",
  "message": "当前状态不允许提交审核",
  "request_id": "req-abc123"
}
```

### 1.4 查询参数

| 参数 | 说明 |
|------|------|
| page | 页码，从 1 开始 |
| page_size | 每页条数，默认 20，最大 100 |
| category | 分类过滤 |
| status | 状态过滤 |

---

## 2. Auth

### POST /auth/login

**角色：** 公开

**Request:**
```json
{
  "email": "user@example.com",
  "password": "secret"
}
```

**Response:**
```json
{
  "data": {
    "token": "eyJhbG...",
    "expires_in": 86400,
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "display_name": "张三",
      "role": "user",
      "tenant_id": "uuid"
    }
  }
}
```

### GET /auth/me

**角色：** 已登录

**Response:** 当前用户信息

### POST /auth/logout

**角色：** 已登录（可选：服务端 token 黑名单）

---

## 3. Import（Admin）

### POST /imports

**角色：** admin

**Request:** `multipart/form-data`
- `file`: CSV 文件

**Response:**
```json
{
  "data": {
    "batch_id": "uuid",
    "status": "processing",
    "file_name": "images.csv"
  }
}
```

### GET /imports

**角色：** admin

**Query:** page, page_size

**Response:** 批次列表

### GET /imports/:id

**角色：** admin

**Response:**
```json
{
  "data": {
    "id": "uuid",
    "status": "processing",
    "total_rows": 1000,
    "success_count": 850,
    "failed_count": 5,
    "created_at": "2026-05-23T10:00:00Z"
  }
}
```

---

## 4. Categories

### GET /categories

**角色：** user, admin, reviewer

**Response:**
```json
{
  "data": [
    { "name": "banner", "available_count": 42 },
    { "name": "poster", "available_count": 18 }
  ]
}
```

---

## 5. Images

### GET /images

**角色：** 按 library 参数限制

**Query:**
| 参数 | 说明 |
|------|------|
| library | tasks / pending / first_review / second_review / archive |
| category | 可选 |
| assigned_to | admin 可查他人 |

**library 映射：**

| library | status 条件 |
|---------|-------------|
| pending | pending_assign |
| tasks | assigned_to=me, status in (assigned,in_progress,rejected) |
| first_review | pending_1st_review |
| passed_first | 1st_review_passed |
| second_review | pending_2nd_review |
| archive | 2nd_review_passed, completed |

### GET /images/:id

**Response:** 图片详情 + 当前版本 + 最近审核记录

### POST /images/claim

**角色：** user

**Request:**
```json
{
  "category": "banner"
}
```

**Response:**
```json
{
  "data": {
    "id": "uuid",
    "global_no": "ACME-20260523-000042",
    "category": "banner",
    "status": "assigned",
    "thumb_url": "https://..."
  }
}
```

**Errors:**
- `409 NO_IMAGE_AVAILABLE` — 该类目无可用图

### POST /images/:id/discard

**角色：** user（须为 assigned_to）

**Request:**
```json
{
  "reason": "图片质量差，无法使用"
}
```

**Response:** 更新后的 image，`status=discarded`

### POST /images/:id/start-editing

**角色：** user

**说明：** assigned → in_progress

---

## 6. Editor / Versions

### GET /images/:id/versions

**Response:** 版本列表

### GET /images/:id/versions/:versionId

**Response:** 版本详情 + layers

### PUT /images/:id/versions

**角色：** user

**Request:**
```json
{
  "layer_data": { "width": 1920, "height": 1080, "layers": [...] },
  "render_file": "base64 或 presigned upload 后的 path"
}
```

**Response:** 新版本信息

### POST /images/:id/submit

**角色：** user

**说明：** in_progress → pending_1st_review

**Request:**
```json
{
  "version_id": "uuid"
}
```

---

## 7. Review

### GET /reviews/queue

**角色：** reviewer

**Query:** round=1|2, page, page_size

### POST /images/:id/reviews

**角色：** reviewer

**Request:**
```json
{
  "round": 1,
  "result": "pass",
  "comment": "符合规范"
}
```

`result`: `pass` | `reject`

**状态变更：**
- round=1, pass → 1st_review_passed
- round=1, reject → rejected
- round=2, pass → 2nd_review_passed
- round=2, reject → rejected（或退回一审通过库，按产品定）

### POST /reviews/claim-second

**角色：** reviewer

**说明：** 从 1st_review_passed 库领取进入 pending_2nd_review

**Response:** 领取到的 image

---

## 8. Archive（Admin）

### POST /archives/export

**角色：** admin

**Request:**
```json
{
  "category": "banner",
  "date_from": "2026-05-01",
  "date_to": "2026-05-31",
  "status": "completed"
}
```

**Response:**
```json
{
  "data": {
    "export_id": "uuid",
    "status": "processing"
  }
}
```

### GET /archives/exports/:id

**Response:** 导出任务状态 + download_url（完成后）

### POST /images/:id/complete

**角色：** admin 或 system

**说明：** 2nd_review_passed → completed（若跳过二审则 1st_review_passed → completed）

---

## 9. Media

### GET /media/presigned

**Query:** path, operation=get|put

**Response:**
```json
{
  "data": {
    "url": "https://minio...",
    "expires_at": "2026-05-23T11:00:00Z"
  }
}
```

---

## 10. 错误码一览

| code | HTTP | 说明 |
|------|------|------|
| VALIDATION_ERROR | 400 | 参数校验失败 |
| UNAUTHORIZED | 401 | 未登录或 token 无效 |
| FORBIDDEN | 403 | 角色无权限 |
| NOT_FOUND | 404 | 资源不存在 |
| NO_IMAGE_AVAILABLE | 409 | 无图可领 |
| INVALID_TRANSITION | 409 | 非法状态转换 |
| DISCARD_LIMIT_EXCEEDED | 429 | 换图次数超限 |
| INTERNAL_ERROR | 500 | 服务器错误 |

---

## 11. Webhook / 事件（Phase 2 可选）

| 事件 | 触发 | Payload |
|------|------|---------|
| image.assigned | 领图 | image_id, user_id |
| image.rejected | 一审驳回 | image_id, comment |
| import.completed | 批次完成 | batch_id, counts |

可用于对接企业微信/钉钉通知。

---

## 12. OpenAPI

建议使用 **swaggo/swag** 从 Go handler 注释生成 OpenAPI 3.0，供前端 codegen 与 Postman 导入。
