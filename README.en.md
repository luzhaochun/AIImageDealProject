# ImageDeal

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev/)
[![React](https://img.shields.io/badge/React-19-61DAFB?logo=react)](https://react.dev/)
[![OpenAI](https://img.shields.io/badge/GPT_Image-gpt--image--2-412991?logo=openai)](https://platform.openai.com/)

**Team image workflow (CSV → edit → review → archive)** plus an **AI inpainting studio** powered by OpenAI GPT Image and Konva.

[中文 README](README.md)

---

## Highlights

- **AI Inpainting Studio** — Brush mask → GPT Image (`gpt-image-2`) → PNG / PSD (3 layers)
- **13-step workflow** — Bulk import, claim, edit, 1st/2nd review, final archive
- **Multi-role RBAC** — Admin, editor, reviewer
- **Async workers** — CSV parse & image download via Redis queues

## AI Studio (Demo)

| Before | After |
|:---:|:---:|
| ![before](docs/images/canvas-studio/demo-sofa-before.png) | ![after](docs/images/canvas-studio/demo-sofa-after.png) |

**Route:** `/admin/canvas-studio` · **API:** `POST /api/v1/admin/ai/studio-inpaint`

![Studio UI](docs/images/canvas-studio/ui-screenshot.png)

Configure `openai_api_key` in `backend/conf/app.conf` (not committed). Details: [docs/10-CANVAS-STUDIO.md](docs/10-CANVAS-STUDIO.md).

## Quick Start

```bash
docker compose up -d redis minio
cd backend && go run .          # API :8080
cd backend && go run ./worker   # required for CSV/download
cd frontend && npm install && npm run dev   # :5173
```

| Account | Role | Password |
|---------|------|----------|
| admin@example.com | admin | password123 |
| user@example.com | user | password123 |
| reviewer@example.com | reviewer | password123 |

## Docs

- [Requirements (13 steps)](docs/REQUIREMENTS.md)
- [API](docs/05-API.md) · [State machine](docs/06-STATE-MACHINE.md)

## Contact

**Tim Lu** · 15651616552 · tim.lu@lianwei.com.cn
