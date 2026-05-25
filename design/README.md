# ImageDeal UI 静态设计稿

> **仅用于设计评审**，非生产代码。在浏览器中直接打开 HTML 即可预览。

## 打开方式

```bash
open /Users/timlu/Developer/go-project/ImageDealProject/design/index.html
```

或双击 `design/index.html`，从索引页进入各页面。

## 页面清单（共 14 屏）

| 页面 | 文件 | 角色 |
|------|------|------|
| 设计索引 | `index.html` | — |
| 登录 | `pages/login.html` | 公开 |
| 工作台首页 | `pages/home.html` | 用户 |
| CSV 导入 | `pages/admin-import.html` | 管理员 |
| 导入批次详情 | `pages/admin-import-detail.html` | 管理员 |
| 终稿导出 | `pages/admin-archive.html` | 管理员 |
| AI 消除 Demo | `pages/admin-ai-editor.html` | 管理员（面试演示） |
| 领图 | `pages/workspace-claim.html` | 用户 |
| 我的任务 | `pages/workspace-tasks.html` | 用户 |
| 在线编辑器 | `pages/workspace-editor.html` | 用户 |
| 一审队列 | `pages/review-first.html` | 审核员 |
| 二审队列 | `pages/review-second.html` | 审核员 |
| 审核对比 | `pages/review-compare.html` | 审核员 |
| 403 | `pages/403.html` | — |

## 设计规范

- 样式：`css/design-system.css`
- 主色：`#2563eb`
- 布局：左侧导航 + 顶栏 + 内容区（类 Ant Design Pro）
- 状态标签：与 `docs/06-STATE-MACHINE.md` 对齐

## 与实现的对应

正式实现见 `frontend/`（React + Ant Design），API 见 `docs/05-API.md`。
