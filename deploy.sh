#!/usr/bin/env bash
# MewsoProxy 服务器部署脚本（默认自适应的镜像分发模式）
# 用法:
#   bash deploy.sh                  # 自动：有镜像用镜像 → 无则加载离线包 → 都无则源码构建兜底(自动加 swap)
#   bash deploy.sh --build          # 强制源码构建（自动加 swap 兜底）
#   bash deploy.sh --pull           # 强制从镜像仓库拉取
# 说明：镜像默认由本地 docker-build.sh 构建导出，或推送到镜像仓库。
set -euo pipefail

cd "$(dirname "$0")"

SERVER_IMG="mewsoproxy/server:latest"
WEB_IMG="mewsoproxy/web:latest"
TAR_GZ="mewsoproxy-images.tar.gz"

log()  { printf '\033[1;32m[deploy]\033[0m %s\n' "$*"; }
warn() { printf '\033[1;33m[deploy]\033[0m %s\n' "$*"; }
err()  { printf '\033[1;31m[deploy]\033[0m %s\n' "$*" >&2; }

command -v docker >/dev/null 2>&1 || { err "未检测到 docker，请先安装 Docker Engine 与 Compose v2"; exit 1; }
docker compose version >/dev/null 2>&1 || { err "未检测到 docker compose v2"; exit 1; }

has_image() { docker image inspect "$1" >/dev/null 2>&1; }

# ---- 仅在 --build 源码构建时才需要内存兜底 ----
ensure_swap() {
  local total_mb swap_mb
  total_mb=$(awk '/^MemTotal:/ {print int($2/1024)}' /proc/meminfo)
  swap_mb=$(awk '/^SwapTotal:/ {print int($2/1024)}' /proc/meminfo)
  if [ "$total_mb" -lt 3072 ] && [ "$swap_mb" -lt 1024 ]; then
    warn "内存偏低(${total_mb}MB)且 swap 不足，创建 2G swapfile..."
    fallocate -l 2G /swapfile || dd if=/dev/zero of=/swapfile bs=1M count=2048
    chmod 600 /swapfile
    mkswap /swapfile >/dev/null
    swapon /swapfile 2>/dev/null || true
    grep -q '/swapfile' /etc/fstab 2>/dev/null || echo '/swapfile none swap sw 0 0' >> /etc/fstab
  fi
}

case "${1:-}" in
  --build)
    ensure_swap
    log "源码构建模式：构建 server 镜像 ..."
    docker compose build server
    log "源码构建模式：构建 web 镜像 ..."
    docker compose build web
    ;;
  --pull)
    log "从镜像仓库拉取镜像 ..."
    docker compose pull
    ;;
  *)
    if has_image "$SERVER_IMG" && has_image "$WEB_IMG"; then
      log "本机已有镜像，跳过加载"
    elif [ -f "$TAR_GZ" ]; then
      log "加载离线镜像包 ${TAR_GZ} ..."
      docker load < "$TAR_GZ"
    else
      warn "未找到本地镜像，执行源码构建兜底（低配机将自动加 swap）..."
      ensure_swap
      log "构建 server 镜像 ..."
      docker compose build server
      log "构建 web 镜像 ..."
      docker compose build web
    fi
    ;;
esac

log "启动全部服务 ..."
docker compose up -d
docker compose ps

cat <<'EOF'

========== 部署完成（镜像模式，零构建） ==========
前端:   http://<服务器IP>:8082
后端:   http://<服务器IP>:8081/healthz
默认管理员: admin@test.com / admin123  (登录后请立即修改密码)

后续步骤:
  1) 打开前端，用默认管理员登录后台
  2) 系统配置 -> 修改「订阅地址」-> 保存
  3) docker compose restart server   (使订阅地址生效)
  4) 将域名反向代理到 http://<服务器IP>:8082
==============================================
EOF
