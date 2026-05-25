#!/usr/bin/env bash
set -euo pipefail
export PATH="/Users/timlu/.nvm/versions/node/v22.22.1/bin:$PATH"
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
DESIGN="$ROOT/design"
OUT="$ROOT/docs/images/design"
mkdir -p "$OUT"
cd "$DESIGN"
python3 -m http.server 9876 &
SERVER_PID=$!
trap 'kill $SERVER_PID 2>/dev/null || true' EXIT
sleep 2
BASE="http://127.0.0.1:9876"
capture() {
  local path="$1" name="$2"
  echo "Capturing $name..."
  npx -y playwright screenshot "$BASE/$path" "$OUT/$name.png" --full-page --viewport-size=1440,900 --wait-for-timeout=500
}
capture index.html index
capture pages/login.html login
capture pages/home.html home
capture pages/403.html 403
capture pages/admin-import.html admin-import
capture pages/admin-import-detail.html admin-import-detail
capture pages/admin-archive.html admin-archive
capture pages/admin-ai-editor.html admin-ai-editor
capture pages/workspace-claim.html workspace-claim
capture pages/workspace-tasks.html workspace-tasks
capture pages/workspace-editor.html workspace-editor
capture pages/review-first.html review-first
capture pages/review-compare.html review-compare
capture pages/workflow-overview.html workflow-overview
echo "Exporting README workflow image..."
npx -y playwright screenshot "$BASE/pages/workflow-overview.html" "$ROOT/docs/images/business-workflow.png" --full-page --viewport-size=1400,900 --wait-for-timeout=500
echo "Done: $(ls "$OUT" | wc -l) design files + business-workflow.png"
