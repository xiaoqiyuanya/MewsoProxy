# 运行阶段的 nginx 基础镜像，可按环境覆盖
ARG NGINX_BASE=nginx:alpine

# 构建阶段
FROM node:22-alpine AS builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build:docker

# 运行阶段（nginx 托管静态资源，监听 8080）
FROM ${NGINX_BASE}
COPY --from=builder --chown=nginx:nginx /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO- http://127.0.0.1:8080/ || exit 1
