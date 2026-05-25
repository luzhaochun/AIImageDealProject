# 部署与运维

---

## 1. 环境划分

| 环境 | 用途 | 基础设施 |
|------|------|----------|
| local | 本地开发 | Docker Compose |
| staging | 预发验证 | 单机 / 小集群 |
| production | 生产 | K8s 或云 VM + 托管 DB |

---

## 2. 本地开发架构

```yaml
# docker-compose.yml 服务
services:
  mysql:      # 3306
  redis:      # 6379
  minio:      # 9000 API, 9001 Console
```

**本地进程：**
```bash
# 终端 1
bee run                 # Beego API :8080

# 终端 2
go run ./worker         # Asynq consumer

# 终端 3
npm run dev             # :5173, proxy → API
```

**Vite 代理（vite.config.ts）：**
```typescript
server: {
  proxy: {
    '/api': 'http://localhost:8080'
  }
}
```

---

## 3. 生产部署拓扑

```
                    ┌─────────────┐
                    │   CDN       │ ← 缩略图/静态资源
                    └──────┬──────┘
                           │
┌──────────┐        ┌──────▼──────┐        ┌─────────────┐
│  用户     │ ──────►│ Nginx/Ingress│ ──────►│  API x N    │
└──────────┘        └──────┬──────┘        └──────┬──────┘
                           │                      │
                    ┌──────▼──────┐        ┌──────▼──────┐
                    │  SPA (OSS)   │        │ Worker x M  │
                    └─────────────┘        └──────┬──────┘
                                                  │
                    ┌─────────────┬───────────────┼───────────────┐
                    │             │               │               │
              ┌─────▼─────┐ ┌─────▼─────┐  ┌──────▼──────┐ ┌─────▼─────┐
              │  MySQL 8  │ │   Redis   │  │  OSS/MinIO  │ │ Prometheus│
              │ (RDS等)   │ │           │  │             │ │ + Grafana │
              └───────────┘ └───────────┘  └─────────────┘ └───────────┘
```

---

## 4. 容器化

### 4.1 API Dockerfile（多阶段）

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /imagedeal .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /imagedeal /imagedeal
COPY conf/ /conf/
EXPOSE 8080
CMD ["/imagedeal"]
```

Worker 镜像类似，入口 `./worker` 或独立模块。

### 4.2 前端

```bash
npm run build   # → dist/
# 上传 dist 至 OSS + CDN，或 Nginx 静态托管
```

---

## 5. K8s 资源（参考）

| 资源 | 副本 | 说明 |
|------|------|------|
| Deployment/api | 2+ | HPA CPU 70% |
| Deployment/worker | 2+ | 按队列深度 HPA |
| ConfigMap | 1 | 非敏感配置 |
| Secret | 1 | JWT、DB 密码 |
| Ingress | 1 | TLS 终止 |

**Worker 扩容信号：** Redis 队列 `asynq:{queue}:pending` 长度 > 100

---

## 6. 环境变量（生产）

| 变量 | 来源 |
|------|------|
| DB_* | Secret |
| JWT_SECRET | Secret（≥32 字符随机） |
| REDIS_ADDR | ConfigMap |
| MINIO_* / OSS_* | Secret |
| APP_ENV=production | ConfigMap |
| CORS_ORIGINS | 生产域名 |

---

## 7. 监控与日志

### 7.1 指标（Prometheus）

| 指标 | 类型 | 说明 |
|------|------|------|
| http_requests_total | counter | 按 route、status |
| http_request_duration_seconds | histogram | P99 延迟 |
| asynq_tasks_processed_total | counter | 任务处理数 |
| asynq_tasks_failed_total | counter | 任务失败数 |
| images_by_status | gauge | 各状态图片数量 |

### 7.2 日志

- 格式：JSON 结构化
- 字段：timestamp, level, request_id, tenant_id, user_id, msg
- 收集：Loki / ELK

### 7.3 告警

| 条件 | 级别 |
|------|------|
| API 5xx > 1% | P1 |
| Worker 失败率 > 5% | P1 |
| 导入批次 failed | P2 |
| 一审队列积压 > 500 | P2 |

---

## 8. CI/CD 流水线

```
Push → Lint/Test → Build Image → Push Registry
                         ↓
              Staging Deploy → Smoke Test
                         ↓
              Manual Approve → Prod Deploy
                         ↓
                   DB Migrate (job)
```

**GitHub Actions 阶段：**
1. `go test ./...`
2. `npm run build && npm run lint`
3. Docker build & push
4. kubectl apply / helm upgrade

---

## 9. 安全清单

- [ ] HTTPS 全站
- [ ] JWT 短过期 + 刷新 token（可选）
- [ ] 上传文件 MIME 校验（不仅看扩展名）
- [ ] SQL 注入：GORM 参数化
- [ ] 租户隔离集成测试
- [ ] MinIO bucket 私有，仅 presigned 访问
- [ ] 限流：登录 5/min，claim 10/min

---

## 10. 灾备 RTO/RPO

| 组件 | RPO | RTO |
|------|-----|-----|
| MySQL | 1h（binlog） | 2h |
| 对象存储 | 0（多副本） | 1h |
| Redis | 可重建队列 | 30min |

---

## 11. 成本估算（小规模 MVP）

| 资源 | 规格 | 月成本（参考） |
|------|------|----------------|
| 云 VM | 2C4G x2 | ¥200 |
| RDS PG | 1C2G | ¥150 |
| Redis | 1G | ¥50 |
| OSS | 100GB | ¥20 |
| CDN | 100GB 流量 | ¥30 |

*实际以云厂商为准。*
