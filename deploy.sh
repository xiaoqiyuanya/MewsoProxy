#!/usr/bin/env bash
# MewsoProxy 从零部署脚本（专为低配服务器优化，如 2C2G）
# 用法:   bash ./deploy.sh   （或 ./deploy.sh，需 chmod +x）
# 说明:
#   1) 若内存不足且无 swap，自动创建 2G swap 兜底，避免构建 OOM
#   2) 串行构建 server 与 web，避免两个镜像并行导致内存峰值爆掉
#   3) 一键启动全部服务
set -euo pipefail

cd "$(dirname "$0")"

log()  { printf '\033[1;32m[deploy]\033[0m %s\n' "$*"; }
warn() { printf '\033[1;33m[deploy]\033[0m %s\n' "$*"; }
err()  { printf '\033[1;31m[deploy]\033[0m %s\n' "$*" >&2; }

command -v docker >/dev/null 2>&1 || { err "未检测到 docker，请先安装 Docker Engine 与 Compose v2"; exit 1; }
docker compose version >/dev/null 2>&1 || { err "未检测到 docker compose v2"; exit 1; }

# ---- 1. 内存兜底：低配且无 swap 时自动加 2G swap ----
ensure_swap() {
  local total_mb swap_mb
  total_mb=$(awk '/^MemTotal:/ {print int($2/1024)}' /proc/meminfo)
  swap_mb=$(awk '/^SwapTotal:/ {print int($2/1024)}' /proc/meminfo)
  log "系统内存: ${total_mb}MB, swap: ${swap_mb}MB"

  if [ "$total_mb" -lt 3072 ] && [ "$swap_mb" -lt 1024 ]; then
    if swapon --show 2>/dev/null | grep -q swapfile; then
      warn "检测到已有 swapfile，跳过创建"
      return 0
    fi
    warn "内存偏低(${total_mb}MB)且 swap 不足，正在创建 2G swapfile..."
    fallocate -l 2G /swapfile || dd if=/dev/zero of=/swapfile bs=1M count=2048
    chmod 600 /swapfile
    mkswap /swapfile >/dev/null
    swapon /swapfile 2>/dev/null || true
    if ! grep -q '/swapfile' /etc/fstab 2>/dev/null; then
      echo '/swapfile none swap sw 0 0' >> /etc/fstab
    fi
    log "swap 已启用"
  fi
}
ensure_swap || warn "swap 创建失败（不影响后续构建，若 OOM 请手动检查）"

# ---- 2. 串行构建（避免两个镜像并行吃满内存） ----
log "构建后端镜像 (server)..."
docker compose build server

log "构建前端镜像 (web)..."
docker compose build web

# ---- 3. 启动 ----
log "启动全部服务..."
docker compose up -d
docker compose ps

cat <<'EOF'

========== 部署完成 ==========
前端:   http://<服务器IP>:8082
后端:   http://<服务器IP>:8081/healthz
默认管理员: admin@test.com / admin123  (登录后请立即修改密码)

后续步骤:
  1) 浏览器打开前端，用默认管理员登录后台
  2) 系统配置 -> 修改「订阅地址」-> 保存
  3) docker compose restart server   (使订阅地址生效)
  4) 将域名反向代理到 http://<服务器IP>:8082
==============================
EOF
