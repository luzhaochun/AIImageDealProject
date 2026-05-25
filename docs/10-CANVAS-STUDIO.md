# GPT 图像消除 · Canvas 工作室

管理端：`/admin/canvas-studio`

## 技术栈

| 能力 | 实现 |
|------|------|
| 画布 / 蒙版 / 变换 | Konva.js + react-konva |
| AI 消除 | **OpenAI GPT Image**（`POST /api/v1/admin/ai/studio-inpaint`） |
| 旧 Demo mock | `POST /api/v1/admin/ai/inpaint`（`/admin/ai-editor`） |
| PSD 导出 | ag-psd（三图层） |

## 配置（必做）

编辑 `backend/conf/app.conf`：

```ini
ai_edit_mode = openai
openai_api_key = sk-your-key-here
openai_image_model = gpt-image-2
```

保存后重启后端：`cd backend && go run .`

前端顶栏显示绿色 `GPT Image · gpt-image-2` 即表示就绪。

## PSD 图层

1. 原始图片  
2. AI 消除结果  
3. 蒙版（黑白灰度）

## 蒙版说明

- 画笔涂红色区域 = 要消除的内容  
- 后端转为 OpenAI 格式：**透明区域 = 编辑区**（已修正，与 API 一致）  
- **边缘羽化** 在送 API 前作用于蒙版

## 剪裁

1. 工具栏选 **剪裁** → 拖动蓝色虚线框 → **应用剪裁**  
2. 原图、蒙版、AI 结果（若有）同步裁剪

## 速度说明（GPT 消除慢的原因）

| 因素 | 说明 |
|------|------|
| OpenAI 云端推理 | `gpt-image-2` 单次通常 **20–90 秒**，属正常 |
| 原 `quality=high` | 已改为默认 **`medium`**（更快） |
| 大图 | 后端/前端长边限制 **1024px** 再送 API |
| 网络 | 上传 PNG + 等待生成 |

可调 `app.conf`：

```ini
openai_image_quality = low      # 更快，略降画质
openai_image_max_long = 768     # 更小更快
```

## 费用

调用 OpenAI Images Edits **按量计费**，请自行控制 Key 与用量。
