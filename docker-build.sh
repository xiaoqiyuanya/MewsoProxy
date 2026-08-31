#!/usr/bin/env bash
# 在本地/CI 构建 server + web 镜像，并导出为离线包。
# 用途：把构建从低配服务器(VPS)上移走——本地构建好后，服务器只需 docker load。
# 用法:  bash ./docker-build.sh
#   产物：mewsoproxy-images.tar.gz（上传到服务器后执行 docker load < mewsoproxy-images.tar.gz）
# 若你有镜像仓库，可改为 docker tag / docker push 分发（见文件末尾注释）。
set -euo pipefail

cd "$(dirname "$0")"

SERVER_IMG="mewsoproxy/server:latest"
WEB_IMG="mewsoproxy/web:latest"

command -v docker >/dev/null 2>&1 || { echo "未检测到 docker，请先安装 Docker；本脚本在本地/CI 运行，不要在低配服务器上执行。"; exit 1; }

echo "[build] 构建后端镜像 (${SERVER_IMG}) ..."
docker compose build server

echo "[build] 构建前端镜像 (${WEB_IMG}) ..."
docker compose build web

echo "[save] 导出离线镜像包 mewsoproxy-images.tar.gz ..."
docker save "$SERVER_IMG" "$WEB_IMG" | gzip > mewsoproxy-images.tar.gz
echo "完成：$(pwd)/mewsoproxy-images.tar.gz"
echo "下一步：上传到服务器后执行  docker load < mewsoproxy-images.tar.gz，再 bash deploy.sh"
echo
echo "（若要走镜像仓库分发，可把 docker compose.yml 的 image 改成 registry 地址后：）"
echo "docker compose push && 服务器端：docker compose pull && bash deploy.sh"
