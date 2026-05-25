# ImageDeal

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![React](https://img.shields.io/badge/React-19-61DAFB?logo=react&logoColor=black)](https://react.dev/)
[![OpenAI](https://img.shields.io/badge/GPT_Image-gpt--image--2-412991?logo=openai&logoColor=white)](https://platform.openai.com/)
[![Vite](https://img.shields.io/badge/Vite-8-646CFF?logo=vite&logoColor=white)](https://vitejs.dev/)

**SaaS image collaboration platform** + **GPT inpainting studio** — bulk CSV import, multi-role claim/edit/review/archive workflow, with Konva masks, OpenAI image removal, and PSD export on the admin side.

*[中文 README](README.md)*

---

## Table of Contents

- [Highlights](#highlights)
- [AI Image Inpainting](#ai-image-inpainting)
- [Team Collaboration Workflow](#team-collaboration-workflow)
- [Quick Start](#quick-start)
- [Features & Routes](#features--routes)
- [UI Design Mockups](#ui-design-mockups)
- [Documentation](#documentation)
- [GitHub Presentation](#github-presentation)
- [Tech Stack](#tech-stack)
- [Contact](#contact)

---

## Highlights

| | |
|---|---|
| **GPT inpainting** | Brush mask → OpenAI `gpt-image-2` inpainting → PNG / PSD (3 layers) |
| **13-step workflow** | CSV import → storage & classification → claim & edit → 1st/2nd review → archive |
| **Multi-role RBAC** | Admin / editor / reviewer, JWT + route guards |
| **Virtual libraries** | Queues driven by `status` (pending, task area, review queues, final archive) |
| **Async worker** | Redis queues for CSV parsing, bulk download, thumbnails |

**Local URLs:** Frontend http://localhost:5173 · API http://localhost:8080

---

<a id="ai-image-inpainting"></a>

## AI Image Inpainting

> **Paint mask → GPT remove → compare layers → export**  
> Admin [Canvas Studio](http://localhost:5173/admin/canvas-studio) · `POST /api/v1/admin/ai/studio-inpaint`

### Real Results (Before → After)

| Remove people (sofa scene) | Remove background cattle (pasture) |
|:---:|:---:|
| ![Before](docs/images/canvas-studio/demo-sofa-before.png) → ![After](docs/images/canvas-studio/demo-sofa-after.png) | ![Before](docs/images/canvas-studio/demo-cows-before.png) → ![After](docs/images/canvas-studio/demo-cows-after.png) |

*Left: original · Right: after GPT Image (live API output)*

### Demo (GIF)

![Canvas Studio demo](docs/images/canvas-studio/demo.gif)

*Before/After animation; to record a real screen capture see [docs/GITHUB-PREVIEW.md](docs/GITHUB-PREVIEW.md)*

### Studio UI

![Canvas Studio](docs/images/canvas-studio/ui-screenshot.png)

### Capabilities

| Capability | Description |
|------------|-------------|
| Smart removal | OpenAI GPT Image (`gpt-image-2`, Images Edits API) |
| Mask editing | Konva brush / eraser / edge feather / crop |
| Prompts | Custom inpainting instructions (Chinese default optimized) |
| Export | PNG · JPEG · WebP; **PSD with 3 layers** (original / AI result / mask) |
| Performance | Default `quality=medium`, max long edge ≤ 1024px, typically 20–90s per run |

### Use Cases

- E-commerce / poster assets — remove clutter or people  
- Background fill and local inpainting  
- Single-image touch-up before CSV pipeline ingestion  

### Try in 30 Seconds

1. Log in as `admin@example.com` / `password123` → open `/admin/canvas-studio`  
2. Set `openai_api_key` in `backend/conf/app.conf` (do not commit) and restart the backend  
3. Upload an image → paint red over areas to remove → **GPT Remove** → export  

```ini
ai_edit_mode = openai
openai_api_key = sk-your-key
openai_image_model = gpt-image-2
```

Technical details: [docs/10-CANVAS-STUDIO.md](docs/10-CANVAS-STUDIO.md) · Legacy mock demo: `/admin/ai-editor` · [docs/09-AI-EDITOR-DEMO.md](docs/09-AI-EDITOR-DEMO.md)

<details open>
<summary><b>More demo comparisons</b></summary>

**Example 1 — Remove one person / remove all people**

| Original (3 people) | Remove one | Remove all |
|:---:|:---:|:---:|
| ![Original](docs/images/canvas-studio/demo-sofa-before.png) | ![Remove one](docs/images/canvas-studio/demo-sofa-remove-one.png) | ![Empty sofa](docs/images/canvas-studio/demo-sofa-after.png) |

![Remove one — side view](docs/images/canvas-studio/demo-sofa-remove-one-alt.png)

</details>

---

## Team Collaboration Workflow

13-step closed loop: **import → storage → classify → claim → edit → submit → review → archive** (ImageDeal product specification).

![ImageDeal full business workflow](docs/images/business-workflow.png)

| Steps | Role | Summary |
|:---:|------|---------|
| 1–3 | Admin / System | CSV import, download to storage, classify by category |
| 4–7 | User | Claim image, discard & re-claim, online edit, submit for 1st review |
| 8–12 | Reviewer | 1st / optional 2nd review, pass or reject |
| 13 | System / Admin | Final archive, bundle export |

**Full requirements, virtual libraries, state machine** → [docs/REQUIREMENTS.md](docs/REQUIREMENTS.md)

---

## Quick Start

### Prerequisites

```bash
docker compose up -d redis minio   # MySQL on localhost 127.0.0.1:3306
mysql -h127.0.0.1 -uroot -p123456 -e "CREATE DATABASE IF NOT EXISTS imagedeal ..."
mysql -h127.0.0.1 -uroot -p123456 imagedeal < backend/sql/migrations/001_init.sql
```

| MySQL | `root` / `123456` · database `imagedeal` |

### Run

```bash
# Terminal 1 — API
cd backend && go mod tidy && go run .

# Terminal 2 — Worker (required for CSV parse & image download)
cd backend && go run ./worker

# Terminal 3 — Frontend
cd frontend && npm install && npm run dev
```

| Service | URL |
|---------|-----|
| Frontend | http://localhost:5173 |
| API | http://localhost:8080 · `GET /api/v1/health` |
| Canvas Studio | http://localhost:5173/admin/canvas-studio |

### Test Accounts

Password for all: **`password123`**

| Email | Role | Entry |
|-------|------|-------|
| admin@example.com | admin | `/admin` · Canvas Studio |
| user@example.com | user | `/workspace` |
| reviewer@example.com | reviewer | `/review` |

Sample CSV files: `backend/sample/` (`import_sample.csv`, etc.) · See [backend/sample/README.md](backend/sample/README.md)

---

## Features & Routes

| Module | Route | Description |
|--------|-------|-------------|
| Admin | `/admin/imports` | CSV import & batches |
| | `/admin/archives` | Final export |
| | `/admin/canvas-studio` | **GPT image inpainting** |
| | `/admin/ai-editor` | AI demo (mock) |
| Workspace | `/workspace/claim` | Claim image |
| | `/workspace/tasks` | My tasks |
| | `/workspace/editor/:id` | Online editor |
| Review | `/review/first` · `/review/second` | 1st / 2nd review queues |
| | `/review/:id` | Original vs. final comparison |

### Project Structure

```
AIImageDealProject/
├── backend/     # Go API (Beego) + Worker
├── frontend/    # React + Vite + Ant Design
├── design/      # Static HTML mockups
├── docs/        # Architecture & requirements
└── docker-compose.yml
```

---

## UI Design Mockups

Static HTML prototypes (`design/`). Local preview: `open design/index.html` · See [design/README.md](design/README.md)

| Group | Screens | Example routes |
|-------|:-------:|----------------|
| Public | 4 | `/login` · `/` · `/403` |
| Admin | 5 | `/admin/imports` · `/admin/canvas-studio` |
| Workspace | 3 | `/workspace/claim` · `/workspace/tasks` |
| Review | 3 | `/review/first` · `/review/:id` |

<details open>
<summary><b>UI design screenshots (14 static screens + Canvas live UI)</b></summary>

#### Public

| Index | Login | Home | 403 |
|:---:|:---:|:---:|:---:|
| ![index](docs/images/design/index.png) | ![login](docs/images/design/login.png) | ![home](docs/images/design/home.png) | ![403](docs/images/design/403.png) |

#### Admin

| CSV import | Batch detail | Archive export | AI demo | Canvas Studio |
|:---:|:---:|:---:|:---:|:---:|
| ![admin-import](docs/images/design/admin-import.png) | ![detail](docs/images/design/admin-import-detail.png) | ![archive](docs/images/design/admin-archive.png) | ![ai-editor](docs/images/design/admin-ai-editor.png) | ![canvas](docs/images/design/canvas-studio.png) |

#### Workspace (editors)

| Claim | My tasks | Editor |
|:---:|:---:|:---:|
| ![claim](docs/images/design/workspace-claim.png) | ![tasks](docs/images/design/workspace-tasks.png) | ![editor](docs/images/design/workspace-editor.png) |

#### Review

| 1st review queue | 2nd review queue | Compare |
|:---:|:---:|:---:|
| ![first](docs/images/design/review-first.png) | ![second](docs/images/design/review-second.png) | ![compare](docs/images/design/review-compare.png) |

</details>

---

## Documentation

| Doc | Content |
|-----|---------|
| [docs/REQUIREMENTS.md](docs/REQUIREMENTS.md) | **Full 13-step requirements**, libraries, state machine |
| [docs/10-CANVAS-STUDIO.md](docs/10-CANVAS-STUDIO.md) | GPT Canvas config & performance |
| [docs/05-API.md](docs/05-API.md) | REST API |
| [docs/06-STATE-MACHINE.md](docs/06-STATE-MACHINE.md) | State transitions |
| [docs/01-OVERVIEW.md](docs/01-OVERVIEW.md) | Architecture overview |
| [docs/07-DEPLOYMENT.md](docs/07-DEPLOYMENT.md) | Deployment |

Regenerate design screenshots: `bash scripts/capture-design-screenshots.sh`

---

## GitHub Presentation

When sharing this repo on GitHub, we recommend:

| Item | File | Action |
|------|------|--------|
| **Social preview image** | `docs/images/social-preview.png` | Upload under **Settings → General → Social preview** (1280×640) |
| **Demo GIF** | `docs/images/canvas-studio/demo.gif` | Embedded in README; replace with a screen recording if desired |

Regenerate cover image & GIF from existing demo assets:

```bash
pip3 install Pillow
python3 scripts/generate-readme-assets.py
```

Full guide (including screen recording): [docs/GITHUB-PREVIEW.md](docs/GITHUB-PREVIEW.md)

---

## Tech Stack

| Layer | Technologies |
|-------|----------------|
| Backend | Go · Beego v2 · MySQL 8 · Redis · JWT |
| Frontend | React 19 · TypeScript · Vite · Ant Design · TanStack Query |
| AI | OpenAI GPT Image · Konva · ag-psd |
| Storage | Local `backend/storage/` (MinIO container reserved) |

---

## Contact

| | |
|---|---|
| **Name** | Tim Lu |
| **Phone** | 15651616552 |
| **Email** | [tim.lu@lianwei.com.cn](mailto:tim.lu@lianwei.com.cn) |
