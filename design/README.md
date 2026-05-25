# ImageDeal UI 静态设计稿

> **仅用于设计评审**，非生产代码。在浏览器中直接打开 HTML 即可预览。

## 打开方式

```bash
open design/index.html
```

或双击 `design/index.html`，从索引页进入各页面。

---

## 页面清单

### 公共（3 屏）

| 页面 | 文件 | 说明 |
|------|------|------|
| 设计索引 | `index.html` | 按模块展示全部设计稿入口 |
| 登录 | `pages/login.html` | 账号密码登录，按角色进入不同工作台 |
| 工作台首页 | `pages/home.html` | 待办统计、快捷入口、系统状态 |
| 403 无权限 | `pages/403.html` | 角色无权访问时的提示页 |

### 管理员 Admin（4 屏 + 1 实现专页）

| 页面 | 文件 | 说明 |
|------|------|------|
| CSV 导入 | `pages/admin-import.html` | 上传 CSV、批次列表、导入进度 |
| 导入批次详情 | `pages/admin-import-detail.html` | 解析进度、成功/失败数、错误日志 |
| 终稿导出 | `pages/admin-archive.html` | 筛选终稿库、打包 ZIP 导出 |
| AI 消除 Demo（mock） | `pages/admin-ai-editor.html` | 画笔蒙版、面试演示（独立于主流程） |
| GPT Canvas 工作室 | — | 无静态稿；见 `frontend` `/admin/canvas-studio` |

### 作图用户 Workspace（3 屏）

| 页面 | 文件 | 说明 |
|------|------|------|
| 领图 | `pages/workspace-claim.html` | 按类目查看库存、领取一张、换图说明 |
| 我的任务 | `pages/workspace-tasks.html` | 任务列表、状态筛选、操作随状态变化 |
| 在线编辑器 | `pages/workspace-editor.html` | 图层、工具栏、保存、提交审核 |

### 审核员 Review（3 屏）

| 页面 | 文件 | 说明 |
|------|------|------|
| 一审队列 | `pages/review-first.html` | 待一审列表、领取/进入审核 |
| 二审队列 | `pages/review-second.html` | 从一审通过库抽样领取 |
| 审核对比 | `pages/review-compare.html` | 原图 vs 成稿、通过/驳回、审核意见 |

---

## 设计截图

各页预览图已生成至 `docs/images/design/`（README 中可直接查看）：

| 页面 | 截图 |
|------|------|
| 设计索引 | `docs/images/design/index.png` |
| 登录 | `docs/images/design/login.png` |
| 工作台首页 | `docs/images/design/home.png` |
| 403 | `docs/images/design/403.png` |
| CSV 导入 | `docs/images/design/admin-import.png` |
| 导入批次详情 | `docs/images/design/admin-import-detail.png` |
| 终稿导出 | `docs/images/design/admin-archive.png` |
| AI 消除 Demo | `docs/images/design/admin-ai-editor.png` |
| GPT Canvas 工作室 | `docs/images/design/canvas-studio.png` |
| 领图 | `docs/images/design/workspace-claim.png` |
| 我的任务 | `docs/images/design/workspace-tasks.png` |
| 在线编辑器 | `docs/images/design/workspace-editor.png` |
| 一审队列 | `docs/images/design/review-first.png` |
| 二审队列 | `docs/images/design/review-second.png` |
| 审核对比 | `docs/images/design/review-compare.png` |

重新生成截图：

```bash
bash scripts/capture-design-screenshots.sh
```

---

## 与前端实现路由对照

| 设计稿 | 实现路由 |
|--------|----------|
| `login.html` | `/login` |
| `home.html` | `/` |
| `403.html` | `/403` |
| `admin-import.html` | `/admin/imports` |
| `admin-import-detail.html` | `/admin/imports/:id` |
| `admin-archive.html` | `/admin/archives` |
| `admin-ai-editor.html` | `/admin/ai-editor` |
| — | `/admin/canvas-studio` |
| `workspace-claim.html` | `/workspace/claim` |
| `workspace-tasks.html` | `/workspace/tasks` |
| `workspace-editor.html` | `/workspace/editor/:id` |
| `review-first.html` | `/review/first` |
| `review-second.html` | `/review/second` |
| `review-compare.html` | `/review/:id` |

---

## 设计规范

- 样式：`css/design-system.css`
- 主色：`#2563eb`
- 布局：左侧导航 + 顶栏 + 内容区（类 Ant Design Pro）
- 状态标签：与 `docs/06-STATE-MACHINE.md` 对齐

## 与实现的对应

正式实现见 `frontend/`（React + Ant Design）。README 中含折叠截图预览；完整需求见 [docs/REQUIREMENTS.md](../docs/REQUIREMENTS.md)。
