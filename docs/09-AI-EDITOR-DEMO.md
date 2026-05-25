# Admin AI 消除 Demo

面试演示功能，**仅 admin** 可访问，与领图/提审流程独立。

## 入口

- 前端：`/admin/ai-editor`
- 设计稿：`design/pages/admin-ai-editor.html`

## 使用步骤

1. `admin@example.com` / `password123` 登录
2. 侧栏 **AI 消除 Demo**
3. **加载演示图（披萨）** 或上传自有图片
4. 红色 **画笔** 涂抹要消除的区域（刀叉等）
5. 点击 **AI 消除** → 原图/结果对比 → 可下载或「应用结果到画布」

> `ai_edit_mode = mock` 时零成本，提示词仅展示不生效。

## 后端配置（`backend/conf/app.conf`）

```ini
# mock = 本地模拟（默认，无需 Key）
ai_edit_mode = mock

# openai = 调用 OpenAI Images Edits API
# ai_edit_mode = openai
# openai_api_key = sk-...
# openai_image_model = gpt-image-1.5
```

## API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/admin/ai/config` | 当前模式 |
| POST | `/api/v1/admin/ai/inpaint` | multipart: `image`, `mask`, `prompt` |

## 技术说明

- **mock**：用蒙版边缘像素均值填充，演示流程，无外部费用
- **openai**：`POST https://api.openai.com/v1/images/edits`，模型如 `gpt-image-1.5` / `gpt-image-1`
- 结果保存在 `storage/ai-demo/{uuid}.png`，通过 `/static/...` 访问
