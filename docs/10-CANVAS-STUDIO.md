# GPT 图像工作室 · Canvas Studio

管理端：`/admin/canvas-studio`（本地 http://localhost:5173/admin/canvas-studio）

**Demo 与说明**：[README.md#gpt-image-studio](../README.md#gpt-image-studio)（EN）· [README.zh-CN.md#gpt-图像工作室](../README.zh-CN.md#gpt-图像工作室)（中文）

## 技术栈

| 能力 | 实现 |
|------|------|
| 画布 / 蒙版 / 变换 | Konva.js + react-konva |
| 多模式 AI | **OpenAI GPT Image**（`POST /api/v1/admin/ai/studio-edit`） |
| 兼容旧接口 | `POST /api/v1/admin/ai/studio-inpaint`（等同消除模式） |
| Mock Demo | `POST /api/v1/admin/ai/inpaint`（`/admin/ai-editor`） |
| PSD 导出 | ag-psd（三图层） |

## AI 功能模式

| 模式 ID | 名称 | 蒙版 | 说明 |
|---------|------|:----:|------|
| `erase` | 消除 | 必填 | 红色涂抹要删除的区域 |
| `upscale` | 高清增强 | 可选 | 可整图增强，或涂蒙版局部增强 |
| `repair` | 修复 | 可选 | 整图或局部修复划痕、破损、噪点 |
| `background` | 换背景 | 可选 | 整图换背景，或涂蒙版保留主体 |
| `outpaint` | 扩图 | 可选 | 整图扩展画布，或涂蒙版指定方向 |
| `watermark` | 去水印 | 可选 | 涂蒙版标出水印位置效果更佳 |

切换模式时自动切换默认「需求描述」；可手动编辑补充。

## 页面布局

```
┌──────────────────────────────────────────────────────────────────┐
│ GPT 图像工作室 · Konva + PSD          GPT Image · gpt-image-2     │
├────────────┬─────────────────────────────────────┬─────────────────┤
│ AI 功能    │ 画布（原图 / AI 结果 / 蒙版）        │ 图层 + 导出      │
│ 消除/增强… │ 当前模式提示条                       │ GPT {模式}      │
│ 工具       │                                     │ PNG/JPEG/PSD    │
│ 笔刷/剪裁  │                                     │                 │
│ 需求描述   │                                     │                 │
│ 上传图片   │                                     │                 │
└────────────┴─────────────────────────────────────┴─────────────────┘
```

## 配置（必做）

编辑 `backend/conf/app.conf`：

```ini
ai_edit_mode = openai
openai_api_key = sk-your-key-here
openai_image_model = gpt-image-2
```

保存后重启后端：`cd backend && go run .`

顶栏绿色 `GPT Image · gpt-image-2` 表示就绪。

## 推荐流程

1. 上传图片或加载演示图  
2. 左侧选择 **AI 功能**（如修复 / 消除）  
3. 消除、去水印等建议用画笔涂蒙版；高清、扩图可整图处理  
4. 编辑「需求描述」（留空则用该模式默认文案）  
5. 点击 **GPT {模式名}** → 等待约 20–90 秒  
6. 切换图层对比 → 导出 PNG / JPEG / WebP / PSD  

## PSD 图层

1. 原图  
2. AI 结果（当前模式）  
3. 蒙版（黑白灰度）

## 蒙版说明

- 画笔红色区域 = 编辑区（消除模式必填）  
- 后端转为 OpenAI 格式：透明区域 = 编辑区  
- **边缘羽化** 在送 API 前作用于蒙版  

## 性能

| 因素 | 说明 |
|------|------|
| OpenAI 云端 | `gpt-image-2` 单次通常 20–90 秒 |
| 默认 quality | `medium` |
| 长边上限 | ≤ 1024px 再送 API |

```ini
openai_image_quality = low
openai_image_max_long = 768
```

## 源码

- 页面：`frontend/src/features/admin/konva/KonvaStudioPage.tsx`
- 模式配置：`frontend/src/constants/aiStudioModes.ts`
- API：`frontend/src/api/aiEditor.ts` → `aiStudioEdit()`
