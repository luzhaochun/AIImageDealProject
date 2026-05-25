import { chromium } from 'playwright'
import { createServer } from 'http'
import { readFileSync, existsSync, mkdirSync } from 'fs'
import { join, extname } from 'path'
import { fileURLToPath } from 'url'

const __dirname = fileURLToPath(new URL('.', import.meta.url))
const root = join(__dirname, '..')
const designDir = join(root, 'design')
const outDir = join(root, 'docs/images/design')
mkdirSync(outDir, { recursive: true })

const mime = {
  '.html': 'text/html; charset=utf-8',
  '.css': 'text/css; charset=utf-8',
  '.js': 'application/javascript',
  '.png': 'image/png',
  '.jpg': 'image/jpeg',
  '.svg': 'image/svg+xml',
}

const pages = [
  { file: 'index.html', name: 'index', title: '设计稿索引' },
  { file: 'pages/login.html', name: 'login', title: '登录' },
  { file: 'pages/home.html', name: 'home', title: '工作台首页' },
  { file: 'pages/403.html', name: '403', title: '403 无权限' },
  { file: 'pages/admin-import.html', name: 'admin-import', title: 'CSV 导入' },
  { file: 'pages/admin-import-detail.html', name: 'admin-import-detail', title: '导入批次详情' },
  { file: 'pages/admin-archive.html', name: 'admin-archive', title: '终稿导出' },
  { file: 'pages/admin-ai-editor.html', name: 'admin-ai-editor', title: 'AI 消除 Demo' },
  { file: 'pages/workspace-claim.html', name: 'workspace-claim', title: '领图' },
  { file: 'pages/workspace-tasks.html', name: 'workspace-tasks', title: '我的任务' },
  { file: 'pages/workspace-editor.html', name: 'workspace-editor', title: '在线编辑器' },
  { file: 'pages/review-first.html', name: 'review-first', title: '一审队列' },
  { file: 'pages/review-second.html', name: 'review-second', title: '二审队列' },
  { file: 'pages/review-compare.html', name: 'review-compare', title: '审核对比' },
]

function startServer() {
  return new Promise((resolve) => {
    const server = createServer((req, res) => {
      const path = req.url === '/' ? '/index.html' : req.url.split('?')[0]
      const filePath = join(designDir, path)
      if (!filePath.startsWith(designDir) || !existsSync(filePath)) {
        res.writeHead(404)
        res.end('Not found')
        return
      }
      const body = readFileSync(filePath)
      res.writeHead(200, { 'Content-Type': mime[extname(filePath)] || 'application/octet-stream' })
      res.end(body)
    })
    server.listen(0, '127.0.0.1', () => {
      const { port } = server.address()
      resolve({ server, port })
    })
  })
}

const { server, port } = await startServer()
const browser = await chromium.launch()
const context = await browser.newContext({ viewport: { width: 1440, height: 900 } })
const page = await context.newPage()

for (const item of pages) {
  const url = `http://127.0.0.1:${port}/${item.file}`
  await page.goto(url, { waitUntil: 'networkidle' })
  await page.waitForTimeout(300)
  const out = join(outDir, `${item.name}.png`)
  await page.screenshot({ path: out, fullPage: true })
  console.log('saved', out)
}

await browser.close()
server.close()
