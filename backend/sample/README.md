# 导入测试 CSV

管理员登录后，在 **导入管理** 页上传以下任一文件。需确保 **Worker** 已启动，图片才会自动下载并进入「待领图」状态。

## 格式要求

| 列名 | 必填 | 说明 |
|------|------|------|
| `image_url` 或 `url` | 是 | 可公网访问的图片直链（HTTP/HTTPS） |
| `category` | 是 | 分类名，用于领图时按类目筛选（任意字符串，同一租户内保持一致即可） |

示例：

```csv
image_url,category
https://picsum.photos/seed/demo/800/600,banner
```

## 文件说明

| 文件 | 行数 | 用途 |
|------|------|------|
| `import_minimal.csv` | 1 | 最快冒烟：导入 → 下载 → 领 1 张图 |
| `import_sample.csv` | 3 | 默认示例：banner×2 + poster×1 |
| `import_batch_medium.csv` | 15 | 多类目压测：banner×5、poster×4、icon×3、social×3 |
| `import_url_alias.csv` | 2 | 表头使用 `url` 列名（兼容别名） |
| `import_100.csv` | 100 | 大批量压测：banner×30、poster×30、icon×20、social×20 |

## 推荐测试路径

1. **admin** 上传 `import_minimal.csv` → 等批次 `success_count = 1`
2. **user** → 领图 → 选类目 `banner` → 作图保存 → 提审
3. **reviewer** → 一审通过 → 二审领图 → 二审通过 → 归档

多用户领图：上传 `import_batch_medium.csv`，用多个 **user** 账号分别领不同 `category`。

## 图片来源

当前样例使用 [Lorem Picsum](https://picsum.photos/) 公网图床，无需 API Key，适合本地联调。生产环境请替换为你们 CDN / OSS 上的真实商品图 URL。

> **为何不用百度图片链接？** 百度搜索结果里的图片 URL 多数带防盗链、会过期或被 403，Worker 批量下载失败率很高。`import_100.csv` 使用可直链下载的公网图床，专为本项目导入/下载流程设计。若必须用自有图库，请把图片放到 OSS/CDN 后写入 `image_url`。

## 自定义 CSV

复制任一样本，修改 `category` 与 `image_url` 即可。注意：

- 空行、缺列、空 URL/分类会计入 `failed_count`
- 下载失败（超时、404）不会增加 `success_count`，可在批次详情查看 `error_log`
