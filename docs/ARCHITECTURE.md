# ImageDeal — 架构设计（索引）

> 完整架构文档已拆分至以下文件，请按需阅读。

## 文档目录

| 序号 | 文档 | 内容 |
|------|------|------|
| 01 | [OVERVIEW](./01-OVERVIEW.md) | 系统总览、架构选型、角色权限、NFR |
| 02 | [BACKEND](./02-BACKEND.md) | 后端分层、模块、Worker、存储 |
| 03 | [FRONTEND](./03-FRONTEND.md) | 前端结构、路由、页面、状态管理 |
| 04 | [DATABASE](./04-DATABASE.md) | ER 图、表结构、索引 |
| 05 | [API](./05-API.md) | REST API 完整定义 |
| 06 | [STATE-MACHINE](./06-STATE-MACHINE.md) | 状态机、13 步流程映射 |
| 07 | [DEPLOYMENT](./07-DEPLOYMENT.md) | 部署、监控、CI/CD |
| 08 | [BEEGO-RESEARCH](./08-BEEGO-RESEARCH.md) | Beego 技术调研、目录结构、选型结论 |

## 推荐阅读顺序

1. **01-OVERVIEW** — 建立全局认知
2. **06-STATE-MACHINE** — 理解核心业务
3. **04-DATABASE** + **05-API** — 数据与接口契约
4. **02-BACKEND** + **03-FRONTEND** — 实现细节
5. **07-DEPLOYMENT** — 上线运维

## 架构一句话

**前后端分离的模块化单体**：**Beego v2** API + 独立 Worker（Asynq）处理图片流水线，React SPA 按角色提供 Admin / Workspace / Review 三套界面，**MySQL 8** 管状态，Redis 管队列，MinIO 存图片。Beego 调研见 [08-BEEGO-RESEARCH.md](./08-BEEGO-RESEARCH.md)。
