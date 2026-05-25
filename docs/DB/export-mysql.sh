#!/usr/bin/env bash
# 从 backend/conf/app.conf 读取 sqlconn，导出 MySQL 到 docs/DB/dumps/
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
DUMP_DIR="${MYSQL_DUMP_DIR:-$ROOT/docs/DB/dumps}"
CONF="${MYSQL_EXPORT_CONF:-$ROOT/backend/conf/app.conf}"
TIMESTAMP="$(date +%Y%m%d_%H%M%S)"

find_mysqldump() {
  if [[ -n "${MYSQLDUMP:-}" ]] && command -v "$MYSQLDUMP" &>/dev/null; then
    echo "$MYSQLDUMP"
    return
  fi
  local candidates=(
    /opt/homebrew/opt/mysql-client/bin/mysqldump
    /opt/homebrew/opt/mysql-client@9.6/bin/mysqldump
    /usr/local/mysql/bin/mysqldump
    mysqldump
  )
  for p in "${candidates[@]}"; do
    if command -v "$p" &>/dev/null; then
      echo "$p"
      return
    fi
  done
  echo "错误: 未找到 mysqldump。可安装: brew install mysql-client" >&2
  exit 1
}

parse_conf() {
  python3 - "$CONF" <<'PY'
import os, re, sys
from pathlib import Path

def env(key, default=""):
    return os.environ.get(key, default)

conf = Path(sys.argv[1])
text = conf.read_text(encoding="utf-8") if conf.is_file() else ""
m = re.search(r"^sqlconn\s*=\s*(.+)$", text, re.M)
dsn = (m.group(1).strip() if m else "") or ""

user = env("MYSQL_USER")
password = env("MYSQL_PASSWORD")
host = env("MYSQL_HOST", "127.0.0.1")
port = env("MYSQL_PORT", "3306")
db = env("MYSQL_DATABASE", "imagedeal")

if dsn:
    mm = re.match(r"([^:]+):([^@]+)@tcp\(([^:]+):(\d+)\)/([^?\s]+)", dsn)
    if mm:
        if not user:
            user = mm.group(1)
        if not password:
            password = mm.group(2)
        if env("MYSQL_HOST", "") == "":
            host = mm.group(3)
        if env("MYSQL_PORT", "") == "":
            port = mm.group(4)
        if env("MYSQL_DATABASE", "") == "":
            db = mm.group(5)

if not user:
    user = "root"
if not password:
    password = "123456"

print(f"export MYSQL_USER={user!r}")
print(f"export MYSQL_PASSWORD={password!r}")
print(f"export MYSQL_HOST={host!r}")
print(f"export MYSQL_PORT={port!r}")
print(f"export MYSQL_DATABASE={db!r}")
PY
}

mkdir -p "$DUMP_DIR"
eval "$(parse_conf)"

MYSQLDUMP_BIN="$(find_mysqldump)"
OUT_FILE="$DUMP_DIR/${MYSQL_DATABASE}_${TIMESTAMP}.sql"
LATEST_LINK="$DUMP_DIR/${MYSQL_DATABASE}_latest.sql"

echo "==> 导出 ${MYSQL_USER}@${MYSQL_HOST}:${MYSQL_PORT}/${MYSQL_DATABASE}"
echo "    工具: $MYSQLDUMP_BIN"
echo "    输出: $OUT_FILE"

export MYSQL_PWD="$MYSQL_PASSWORD"
"$MYSQLDUMP_BIN" \
  -h"$MYSQL_HOST" \
  -P"$MYSQL_PORT" \
  -u"$MYSQL_USER" \
  --single-transaction \
  --routines \
  --triggers \
  --events \
  --set-gtid-purged=OFF \
  --default-character-set=utf8mb4 \
  --databases "$MYSQL_DATABASE" \
  > "$OUT_FILE"
unset MYSQL_PWD

ln -sf "$(basename "$OUT_FILE")" "$LATEST_LINK"

BYTES="$(wc -c < "$OUT_FILE" | tr -d ' ')"
echo "==> 完成: $(du -h "$OUT_FILE" | cut -f1) ($BYTES bytes)"
echo "    最新链接: $LATEST_LINK"
