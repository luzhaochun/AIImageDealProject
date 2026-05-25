# MySQL 导出

本目录用于保存 **在线 MySQL**（本机 Docker / 远程实例）的逻辑备份。

## 快速导出

```bash
# 默认读取 backend/conf/app.conf 中的 sqlconn
./docs/DB/export-mysql.sh
```

导出文件：

- `dumps/imagedeal_YYYYMMDD_HHMMSS.sql` — 带时间戳
- `dumps/imagedeal_latest.sql` — 指向最近一次导出

## 环境变量（可选）

| 变量 | 说明 |
|------|------|
| `MYSQL_HOST` | 主机，默认从 app.conf 解析 |
| `MYSQL_PORT` | 端口 |
| `MYSQL_USER` / `MYSQL_PASSWORD` | 账号密码 |
| `MYSQL_DATABASE` | 库名，默认 `imagedeal` |
| `MYSQL_EXPORT_CONF` | 配置文件路径 |
| `MYSQLDUMP` | mysqldump 可执行文件路径 |

远程示例：

```bash
MYSQL_HOST=db.example.com MYSQL_PORT=3306 MYSQL_USER=root MYSQL_PASSWORD=secret \
  ./docs/DB/export-mysql.sh
```

## 依赖

- `mysqldump`（macOS: `brew install mysql-client`）
- MySQL 服务可连通（本机常见为 Docker `mysql8` → `127.0.0.1:3306`）

## 恢复

```bash
mysql -h127.0.0.1 -uroot -p123456 < docs/DB/dumps/imagedeal_latest.sql
```

## 说明

- `dumps/*.sql` 已加入 `.gitignore`，避免误提交大数据与敏感数据。
- 生产导出请在安全网络下执行，勿将含真实数据的 dump 提交到仓库。
