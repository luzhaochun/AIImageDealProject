# 后端架构设计

> **框架选型**：Beego v2 + Beego ORM + MySQL 8（详见 [08-BEEGO-RESEARCH.md](./08-BEEGO-RESEARCH.md)）  
> **异步方案**：独立 Worker + Redis/Asynq（批量下载、导出）；Beego `task` 仅用于轻量定时任务

---

## 1. 总体原则

- **Beego MVC**：Router → Controller → Service → Model（ORM）
- **领域逻辑在 Service**：Controller 保持薄层；状态机禁止在 Controller 直接改 `status`
- **租户隔离**：Filter 注入 `tenant_id`，所有 ORM 查询带租户条件
- **同步/异步分离**：Beego 处理 HTTP；重 IO 由 `cmd/worker` + Asynq 消费

---

## 2. 进程模型

```
┌─────────────────────────┐     ┌─────────────────────────┐
│   main.go (Beego)       │     │   worker/main.go        │
│   HTTP :8080            │     │   Asynq Consumer        │
└───────────┬─────────────┘     └───────────┬─────────────┘
            │                               │
            └───────────────┬───────────────┘
                            ▼
            共享 models/ services/ pkg/
            共享 MySQL 8 / Redis / MinIO
```

| 进程 | 职责 | 扩容 |
|------|------|------|
| **Beego API** | REST、鉴权、业务编排、入队 | 水平扩容 |
| **Worker** | CSV 解析、下载、缩略图、ZIP | 按队列深度扩容 |

---

## 3. 目录结构（Beego 约定）

```
imagedeal/
├── main.go
├── conf/app.conf
├── routers/router.go
├── controllers/          # 按角色/业务分包
├── models/               # Beego ORM 模型 + init 注册
├── services/             # 状态机、领图、审核（纯 Go）
├── filters/              # JWT、RBAC、CORS
├── tasks/                # Beego cron（清理等）
├── worker/               # Asynq 消费者（独立 main）
├── pkg/                  # storage、queue、response
└── sql/migrations/
```

完整说明见 [08-BEEGO-RESEARCH.md](./08-BEEGO-RESEARCH.md) 第 4 节。

---

## 4. 分层职责

### 4.1 Controller 层

- 嵌入 `web.Controller` 或项目 `BaseController`
- 参数绑定：`ParseForm` / `UnmarshalJSON`
- 调用 Service，通过 `ServeJSON` 返回统一结构
- **禁止**：`orm.NewOrm()` 复杂事务、直接改 status

### 4.2 Service 层

- 无 Beego 依赖，便于单测
- 状态机、`ClaimImage`、`SubmitReview` 等命名方法
- 调用 `models` 或 Raw SQL；需要时调用 `pkg/queue` 入队

### 4.3 Model 层（Beego ORM）

- `models/init.go` 中 `RegisterModel` / `RegisterDataBase`
- 实体 struct + orm tag
- 简单 CRUD；复杂查询可由 Service 使用 Raw

### 4.4 Filter 层

```
CORS → RequestID → JWTAuth → InjectTenant → RBAC → Controller
```

---

## 5. 核心业务模块

| 模块 | Controller | Service | 关键 API |
|------|------------|---------|----------|
| **auth** | AuthController | AuthService | `POST /api/v1/auth/login` |
| **import** | ImportController | ImportService | `POST /api/v1/imports` |
| **image** | ImageController | ImageService | `POST /api/v1/images/claim` |
| **editor** | EditorController | EditorService | `PUT /api/v1/images/:id/versions` |
| **review** | ReviewController | ReviewService | `POST /api/v1/images/:id/reviews` |
| **archive** | ArchiveController | ArchiveService | `POST /api/v1/archives/export` |

API 明细见 [05-API.md](./05-API.md)。

---

## 6. 异步任务

| 任务 | 执行方 | 说明 |
|------|--------|------|
| `import:parse_csv` | Worker (Asynq) | 解析 CSV、写入 images |
| `image:download` | Worker | 下载、校验、上传 MinIO |
| `image:thumb` | Worker | 缩略图 |
| `archive:export` | Worker | 打包 ZIP |
| `cleanup:temp` | Beego task (cron) | 每日清理临时文件 |

---

## 7. 对象存储布局

```
{bucket}/
├── imports/{tenant_id}/{batch_id}/source.csv
├── originals/{tenant_id}/{image_id}.jpg
├── thumbs/{tenant_id}/{image_id}_thumb.jpg
├── versions/{tenant_id}/{image_id}/v{n}.png
└── exports/{tenant_id}/{export_id}.zip
```

---

## 8. 配置项

| 变量 / conf | 说明 |
|-------------|------|
| `httpport` | API 端口 |
| MySQL DSN | `RegisterDataBase` |
| `redis_host` | 缓存 + Asynq |
| `jwt_secret` | 自定义配置项 |
| `MINIO_*` | 对象存储 |

---

## 9. 测试策略

| 层级 | 方式 |
|------|------|
| Service | 单元测试 + mock ORM |
| Controller | `httptest` + Beego 测试模式 |
| Worker | Asynq handler 单测 |
| 集成 | testcontainers MySQL + Redis |

---

## 10. 演进路径

1. **MVP**：Beego API + Worker，完成主流程  
2. **增强**：编辑器、二审、导出  
3. **扩展**：Media/Import 拆独立服务；事件驱动
