# Beego 技术调研报告 — ImageDeal 后端选型

> **调研结论**：ImageDeal 后端采用 **Beego v2.3.x + Beego ORM + MySQL 8**，异步重任务采用 **独立 Worker 进程 + Redis 队列（Asynq）** 与 Beego 主进程配合，而非仅依赖 Beego 内置 `task` 模块。  
> **版本基准**：`github.com/beego/beego/v2` v2.3.10（2026-03，持续维护）  
> **状态**：技术调研 / 架构设计阶段，**暂不编码实现**

---

## 1. 调研背景与目标

### 1.1 业务特征（来自业务流程图）

| 特征 | 技术影响 |
|------|----------|
| 多角色 RBAC（Admin/User/Reviewer） | 需要 Filter 中间件、路由分组 |
| 复杂状态机（11+ 状态） | Service 层集中管理，ORM 事务 |
| CSV 批量导入 + 逐张下载 | **长耗时异步**，非 HTTP 同步 |
| 在线编辑 + 版本/图层 | JSON 字段、版本表 |
| 审核对比、图库视图 | REST API + 分页查询 |
| SaaS 多租户 | 全表 `tenant_id` 隔离 |

### 1.2 调研问题

1. Beego v2 是否仍活跃、适合 2026 年新项目？  
2. Beego MVC 如何映射本项目的 domain 划分？  
3. Beego ORM 对 **MySQL 8** 的支持是否满足？  
4. 内置 `task` 能否替代 Worker/消息队列？  
5. 与 Gin + GORM 方案相比的利弊？

---

## 2. Beego 框架概览

### 2.1 定位

Beego 是 Go 语言的 **全栈企业级 Web 框架**（类比 Java Spring MVC / Python Django），强调：

- **MVC 分层**：Router → Controller → Model（ORM）
- **模块化**：ORM、Cache、Logs、Config、Session、Task、Httplib 等可独立使用
- **快速交付**：`bee` CLI 生成项目、注解路由、热编译

官方仓库：[github.com/beego/beego](https://github.com/beego/beego)（32k+ stars，2026 年仍有版本发布）

### 2.2 核心模块（v2）

| 模块 | 包路径 | ImageDeal 用途 |
|------|--------|----------------|
| **web** | `server/web` | HTTP、Controller、Filter、路由 |
| **orm** | `client/orm` | MySQL 模型、CRUD、事务、QuerySet |
| **cache** | `client/cache` | Redis 缓存、可选 Session |
| **logs** | `core/logs` | 结构化日志 |
| **config** | `core/config` | `app.conf` / 环境配置 |
| **task** | `task` | 定时任务（cron），**非消息队列** |
| **httplib** | `client/httplib` | Worker 内下载外链图片 |
| **toolbox** | `server/web/toolbox` | 可选：健康检查、性能监控 |

### 2.3 与 Gin 的定位差异

| 维度 | Gin | Beego v2 |
|------|-----|----------|
| 类型 | 微框架（路由 + Context） | 全栈框架（MVC + 内置组件） |
| 架构 | 自由组合（常配 GORM） | 约定 MVC + 官方 ORM |
| 性能 | 更高（radix 路由） | 中等（正则路由为主） |
| 学习曲线 | 低，需自选 ORM/日志等 | 中，文档与模块一体 |
| 适用 | 高 QPS API、微服务 | **企业后台、复杂业务、快速交付** |
| 生态 | 第三方中间件极多 | 官方模块为主，第三方较少 |

**对本项目的意义**：ImageDeal 是 **多模块业务后台 + 工作流**，不是极限 QPS 网关；Beego 的「开箱即用」与 MVC 约定更契合团队按业务域拆分 Controller/Model。

---

## 3. 推荐架构：Beego MVC + 独立 Worker（混合）

### 3.1 为什么不只用 Beego 内置能力？

| 需求 | Beego 能力 | 是否足够 |
|------|------------|----------|
| REST API | ✅ Controller + 路由 | 足够 |
| MySQL CRUD | ✅ ORM | 足够 |
| JWT / RBAC | ✅ Filter | 足够 |
| 定时清理日志 | ✅ `task` 模块（cron） | 足够 |
| **CSV 解析后批量下载 N 张图** | ❌ `task` 仅 cron，无 Redis 队列、无重试/backoff | **不足** |
| **导入失败重试、并发度控制** | 需消息队列 | 建议 **Asynq + Redis** |
| **领图分布式锁** | ORM 事务 `FOR UPDATE SKIP LOCKED` | 足够（MySQL 8） |

**结论**：主进程用 Beego；**重异步**单独 `cmd/worker`（仍可用 Go 标准库 + Asynq，与 Beego 解耦），通过 Redis 与 API 通信。Beego 的 `task` 仅用于轻量定时任务（如清理临时文件、统计报表）。

### 3.2 逻辑架构（Beego 版）

```
┌─────────────────────────────────────────────────────────────┐
│                    Beego Application (cmd/main)              │
│  ┌─────────────┐  ┌──────────────┐  ┌─────────────────────┐ │
│  │   Router    │→ │  Controller  │→ │      Service        │ │
│  │ NS + Filter │  │ (按业务分包)  │  │ (状态机/业务编排)    │ │
│  └─────────────┘  └──────────────┘  └──────────┬──────────┘ │
│                                                 │            │
│                                    ┌────────────▼──────────┐ │
│                                    │   Model (Beego ORM)   │ │
│                                    └────────────┬──────────┘ │
└─────────────────────────────────────────────────┼──────────┘
                                                  │
         ┌────────────────────────────────────────┼────────────┐
         │                                        │            │
  ┌──────▼──────┐                          ┌──────▼──────┐ ┌───▼───┐
  │  MySQL 8    │                          │    Redis    │ │ MinIO │
  └─────────────┘                          └──────┬──────┘ └───────┘
                                                  │
                                           ┌──────▼──────┐
                                           │ cmd/worker  │
                                           │ Asynq 消费   │
                                           └─────────────┘
```

### 3.3 与原文档（Gin + GORM）的映射

| 原分层（Gin 方案） | Beego 方案 |
|-------------------|------------|
| `handler/` | `controllers/`（嵌入 `web.Controller`） |
| `service/` | `services/`（纯 Go 包，无框架依赖） |
| `repository/` + GORM | `models/` + Beego ORM，或 Service 内 `orm.NewOrm()` |
| `middleware/` | `filters/`（`web.FilterFunc`） |
| `router/router.go` | `routers/router.go` + `NSRouter` 命名空间 |
| `domain/` | `models/` 实体 + `services/` 领域逻辑 |
| `cmd/worker` | **保留**，不强行用 Beego task 做队列 |

---

## 4. 推荐目录结构（Beego 约定 + 业务扩展）

```
imagedeal/
├── main.go                      # 入口：注册路由、ORM、Filter
├── go.mod
├── conf/
│   ├── app.conf                 # Beego 主配置
│   └── app.local.conf           # 本地覆盖
├── routers/
│   └── router.go                # 路由、命名空间、Filter 链
├── controllers/
│   ├── base.go                  # 基类：取 JWT、tenant_id、统一 JSON 响应
│   ├── auth/
│   │   └── auth.go
│   ├── admin/
│   │   ├── import.go            # CSV 导入
│   │   └── archive.go
│   ├── workspace/
│   │   ├── image.go             # 领图、换图、列表
│   │   └── editor.go            # 版本、提审
│   └── review/
│       └── review.go
├── models/
│   ├── init.go                  # RegisterModel、RegisterDataBase
│   ├── image.go
│   ├── user.go
│   ├── import_batch.go
│   ├── image_version.go
│   ├── image_layer.go
│   └── review_record.go
├── services/                    # 业务核心（推荐保留，避免 Controller 臃肿）
│   ├── image_service.go
│   ├── import_service.go
│   ├── review_service.go
│   └── archive_service.go
├── filters/
│   ├── auth.go                  # JWT 校验
│   ├── rbac.go
│   └── cors.go
├── tasks/                       # Beego cron（轻量）
│   └── cleanup.go
├── worker/                      # 独立进程（重任务）
│   ├── main.go
│   └── handlers/
│       ├── parse_csv.go
│       ├── download_image.go
│       └── export_archive.go
├── pkg/
│   ├── storage/                 # MinIO
│   ├── queue/                   # Asynq client
│   └── response/                # 统一 API 响应
├── sql/
│   └── migrations/              # 手工 SQL 迁移（推荐）
├── static/                      # 可选
└── views/                       # 本项目为前后端分离，可留空
```

**说明**：前后端分离场景下，`views/` 几乎不用；API 项目可用 `bee api` 生成骨架。

---

## 5. Beego ORM + MySQL 8

### 5.1 注册与连接

```go
import (
    "github.com/beego/beego/v2/client/orm"
    _ "github.com/go-sql-driver/mysql"
)

func init() {
    orm.RegisterModel(
        new(User), new(Image), new(ImportBatch),
        new(ImageVersion), new(ImageLayer), new(ReviewRecord),
    )
    orm.RegisterDataBase("default", "mysql",
        "root:123456@tcp(127.0.0.1:3306)/imagedeal?charset=utf8mb4&parseTime=True&loc=Local",
    )
    // 生产环境：orm.RunSyncdb 关闭，改用 sql/migrations
}
```

### 5.2 类型映射（ImageDeal）

| 业务字段 | Go 类型 | ORM 标签建议 |
|----------|---------|--------------|
| 主键 UUID | `string` 或自定义类型 | `orm:"pk;size(36)"` |
| tenant_id | `string` | `orm:"size(36);index"` |
| status | `string` | `orm:"size(30);index"` |
| 图层 JSON | `string` 或 `json.RawMessage` | `orm:"type(json)"` |
| 时间 | `time.Time` | `orm:"auto_now_add;type(datetime)"` |
| 布尔 | `bool` | 映射 MySQL `TINYINT(1)` |

> **注意**：Beego ORM 对「非自增联合主键」支持较弱，**推荐单字段 UUID 主键**（应用层 `uuid.New()` 生成）。

### 5.3 领图事务（与 MySQL 8 文档一致）

```go
o := orm.NewOrm()
err := o.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
    var img Image
  // Raw 或 QueryRow + FOR UPDATE SKIP LOCKED
  // 更新 status、assigned_to、global_no
    return nil
})
```

GORM 的 `clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}` 在 Beego 中通常用 **Raw SQL** 或 **QuerySeter + Raw** 实现。

### 5.4 ORM 使用建议

| 实践 | 说明 |
|------|------|
| 生产禁用 `RunSyncdb` | 表结构用 `sql/migrations` 版本管理 |
| 开发可 `orm.Debug = true` | 打印 SQL |
| 复杂报表用 Raw SQL | 避免 ORM 过度嵌套 |
| Service 层封装 Orm | Controller 不直接 `orm.NewOrm()` |

---

## 6. 路由、Controller、Filter 设计

### 6.1 命名空间（与 API 文档对齐）

```go
// routers/router.go 示意
ns := beego.NewNamespace("/api/v1",
    beego.NSNamespace("/auth",
        beego.NSRouter("/login", &controllers.AuthController{}, "post:Login"),
    ),
    beego.NSNamespace("/imports",
        beego.NSBefore(filters.RequireAdmin),
        beego.NSRouter("/", &controllers.ImportController{}, "post:Create"),
        beego.NSRouter("/:id", &controllers.ImportController{}, "get:Get"),
    ),
    beego.NSNamespace("/images",
        beego.NSBefore(filters.RequireAuth),
        beego.NSRouter("/claim", &controllers.ImageController{}, "post:Claim"),
        // ...
    ),
)
beego.AddNamespace(ns)
```

### 6.2 Filter 链（中间件）

```
全局：CORS → RequestID → RecoverPanic
/api/v1/*：JWTAuth → InjectTenant
/admin/*：RequireRole("admin")
/review/*：RequireRole("reviewer")
```

JWT 可使用 `github.com/golang-jwt/jwt/v5`，在 Filter 中解析后写入 `context.Context`，基类 Controller 读取。

### 6.3 Controller 基类

```go
type BaseController struct {
    web.Controller
}

func (c *BaseController) Success(data interface{}) {
    c.Data["json"] = map[string]interface{}{"data": data}
    c.ServeJSON()
}

func (c *BaseController) TenantID() string { /* 从 context */ }
func (c *BaseController) UserID() string   { /* 从 context */ }
```

---

## 7. 异步与任务分工

| 场景 | 方案 | 说明 |
|------|------|------|
| Admin 上传 CSV | API 存文件 → 写 batch → **Asynq 入队** | 立即返回 batch_id |
| 解析 CSV、批量下载 | Worker 消费 | 并发度可配置 |
| 生成缩略图 | Worker | 下载完成后链式任务 |
| 导出 ZIP | Worker | 可能运行数十分钟 |
| 每日清理临时文件 | **Beego task** cron | `0 0 3 * * *` |
| 队列积压监控 | Prometheus + Redis | 与框架无关 |

**Beego `task` 示例（仅适合定时）：**

```go
tk := task.NewTask("cleanup", "0 0 3 * * *", func(ctx context.Context) error {
    return services.CleanupTempFiles(ctx)
})
task.AddTask("cleanup", tk)
task.StartTask()
```

---

## 8. 业务模块 → Beego 实现对照

| 业务模块 | Controller | Service | Model | 异步 |
|----------|------------|---------|-------|------|
| Auth | AuthController | AuthService | User | — |
| Import | ImportController | ImportService | ImportBatch, Image | parse_csv |
| Image | ImageController | ImageService | Image | download |
| Editor | EditorController | EditorService | ImageVersion, ImageLayer | — |
| Review | ReviewController | ReviewService | ReviewRecord | — |
| Archive | ArchiveController | ArchiveService | Image | export |

**状态机**：仅在 `ImageService` / `ReviewService` 中变更 `Image.Status`，与框架无关。

---

## 9. 优劣势与风险（针对 ImageDeal）

### 9.1 优势

1. **MVC 约定清晰**：适合多角色、多模块后台。  
2. **ORM 内置**：减少 GORM + Gin 的胶水代码。  
3. **配置/日志/缓存统一**：降低企业级项目集成成本。  
4. **团队若熟悉 Beego**：上手快，`bee` 工具链提升效率。  
5. **社区与文档**：中文资料多，v2 持续维护。

### 9.2 劣势与风险

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| 性能低于 Gin | 常规后台 QPS 足够 | 热点加 Redis 缓存；静态资源 CDN |
| ORM 能力弱于 GORM | 复杂查询、hook | 复杂 SQL 用 Raw；关键路径写集成测试 |
| `task` 非消息队列 | 批量下载不能仅靠 cron | **独立 Worker + Asynq** |
| 路由正则 | 大量路由时略慢 | 路由表不宜过深；避免正则滥用 |
| 非自增主键支持一般 | UUID 主键 | 应用层生成 UUID，避免联合主键 |
| 框架绑定 | 迁移成本 | Service/Model 与 web 解耦，便于将来拆服务 |

### 9.3 不推荐的做法

- ❌ 用 Beego `task` 模拟 Redis 队列做批量下载  
- ❌ 在 Controller 里直接改 `image.Status`  
- ❌ 生产环境 `RunSyncdb` 自动改表  
- ❌ 忽略 `tenant_id` 的全局 Filter  

---

## 10. 技术栈定稿（Beego 版）

| 层级 | 选型 |
|------|------|
| Web 框架 | **Beego v2.3.x** |
| ORM | **Beego ORM** + `go-sql-driver/mysql` |
| 数据库 | **MySQL 8.0**（utf8mb4） |
| 队列 | **Redis + Asynq**（Worker 进程） |
| 对象存储 | MinIO / 云 OSS |
| 缓存 | Redis（Beego cache 适配器） |
| 认证 | JWT + Filter |
| 迁移 | golang-migrate / 手工 SQL |
| CLI | bee v2（可选） |
| 前端 | React + TS + Vite（不变） |

---

## 11. 开发工具链

```bash
# 安装框架
go get github.com/beego/beego/v2@latest

# 安装 bee CLI（项目生成、热编译）
go install github.com/beego/bee/v2@latest

# 创建 API 项目骨架（前后端分离推荐）
bee api imagedeal -tables=""

# 本地运行
bee run
```

**配置示例 `conf/app.conf`：**

```ini
appname = imagedeal
httpport = 8080
runmode = dev
copyrequestbody = true

# MySQL
mysqluser = imagedeal
mysqlpass = imagedeal
mysqlurls = 127.0.0.1:3306
mysqldb = imagedeal

# Redis
redis_host = 127.0.0.1:6379
```

---

## 12. 实施阶段建议（仍不写业务代码）

| 阶段 | 内容 |
|------|------|
| P0 | `bee api` 初始化、MySQL 迁移脚本、Model 注册、健康检查 |
| P1 | Auth Filter + 登录 API |
| P2 | Import API + Worker 解析/下载 |
| P3 | 领图/换图/状态机 Service |
| P4 | 编辑器版本 API |
| P5 | 审核 API |
| P6 | 归档导出 |

---

## 13. 参考资料

- [Beego 官方文档 v2](https://beegodoc.com/en-US/v2.2.x/)
- [Beego ORM Quick Start](https://beegodoc.com/en-US/v2.2.x/orm/)
- [Beego Task 模块](https://beegodoc.com/en-US/developing/task/)
- [GitHub beego/beego](https://github.com/beego/beego)
- 项目内：[01-OVERVIEW.md](./01-OVERVIEW.md)、[02-BACKEND.md](./02-BACKEND.md)、[04-DATABASE.md](./04-DATABASE.md)

---

## 14. 调研结论（决策记录）

| 决策项 | 结论 |
|--------|------|
| Go Web 框架 | **Beego v2**（用户指定） |
| 数据库 | **MySQL 8**（用户指定） |
| ORM | **Beego ORM**（与框架一体） |
| 异步重任务 | **独立 Worker + Asynq**，不用 Beego task 替代 |
| 架构风格 | **MVC + Service 层** + 前后端分离 |
| 是否微服务 | **否**（模块化单体） |

**下一步（待你确认后）**：按本文目录结构初始化 Beego 项目，并编写 MySQL migration，仍可按模块分 PR 交付。
