# GitHub 仓库展示优化

本文说明如何设置 **Social Preview（分享封面）** 与 **演示 GIF**，提升仓库在 GitHub / 链接预览中的观感。

---

## 1. Social Preview 封面图

### 推荐尺寸

- **1280 × 640** 像素（GitHub 官方推荐比例 2:1）
- 格式 PNG 或 JPG

### 仓库内素材

项目已生成默认封面：

```
docs/images/social-preview.png
```

重新生成（基于现有 Demo 图拼接）：

```bash
pip3 install Pillow
python3 scripts/generate-readme-assets.py
```

### 在 GitHub 上设置

1. 打开仓库 → **Settings** → **General**
2. 滚动到 **Social preview**
3. 点击 **Upload an image**，选择 `docs/images/social-preview.png`
4. 保存

设置后，在 Slack、微信、Twitter 等分享仓库链接时会显示该封面。

### 自定义封面建议

- 左侧 Before / 右侧 After（已有消除对比）
- 中间或角标加上 **ImageDeal** 与 **GPT Image**
- 背景色与 Canvas 工作室一致（`#0f172a`）更统一

可用 Figma / Canva / Keynote 导出 1280×640 后覆盖 `docs/images/social-preview.png`。

---

## 2. Canvas 操作演示 GIF

### 仓库内素材

```
docs/images/canvas-studio/demo.gif
```

README AI 章节已引用。当前 GIF 为 **Before/After 渐变切换**（自动生成）。

### 重新生成（静态图动画）

```bash
pip3 install Pillow
python3 scripts/generate-readme-assets.py
```

### 录制真实操作 GIF（推荐）

1. 启动前后端（见 README 快速开始）
2. 使用 macOS **⌘ + Shift + 5** 或 [ScreenToGif](https://www.screentogif.com/) / [Kap](https://getkap.co/) 录制：
   - 登录 admin → 打开 `/admin/canvas-studio`
   - 选择 AI 功能（如修复）→ 按需涂蒙版 → 点击 GPT 处理 → 切换图层 → 导出
3. 导出为 GIF，建议：
   - 宽度 **800–960px**
   - 时长 **15–30 秒**
   - 文件 **&lt; 10MB**（GitHub README 加载更快）
4. 覆盖保存为 `docs/images/canvas-studio/demo.gif`
5. 提交并 push

---

## 3. README 引用位置

| 资源 | README 位置 |
|------|-------------|
| `demo.gif` | [GPT 图像工作室](https://github.com/luzhaochun/AIImageDealProject#gpt-图像工作室)（中文 README） |
| `social-preview.png` | 本文档 + 需手动上传到 GitHub Settings |

---

## 4. 检查清单

- [ ] `docs/images/canvas-studio/demo.gif` 可在 README 中正常播放
- [ ] GitHub Settings 已上传 Social Preview
- [ ] 分享仓库链接预览显示封面图
