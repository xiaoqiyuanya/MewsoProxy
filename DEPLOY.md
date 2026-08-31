# MewsoProxy 从零部署验证指南

本指南描述在一台**全新服务器/虚拟机**上，通过 Docker Compose 从零启动一套完整可用的 MewsoProxy 面板，并逐项验证其功能。目标是：**只上传源码 → `docker compose up -d --build` → 访问前端 → 登录后台 → 改订阅地址 → 用户可订阅、下单、支付**，全程无需手工建库、无需手工配置任何密钥。

## 0. 架构与端口

| 服务 | 容器名 | 镜像 | 宿主端口 | 作用 | 健康检查 |
| --- | --- | --- | --- | --- | --- |
| MySQL | `mewsoproxy-mysql` | `mysql:8.4` | `3306` | 数据库（自动建库） | `mysqladmin ping` |
| Redis | `mewsoproxy-redis` | `redis:7-alpine` | `6379` | 缓存/会话 | `redis-cli ping` |
| 后端 | `mewsoproxy-server` | 本地构建 | `8081` | Go API（`/api/v1`） | `GET /healthz` |
| 前端 | `mewsoproxy-web` | 本地构建（nginx） | `8082` | 静态站 + `/api` 反代 | `GET /` |

> 前端容器内的 nginx 已把 `/api/` 转发到后端容器（`server:8080`），因此**只要把域名反向代理到前端端口 `8082` 即可**，无需额外配置 API 转发。

## 1. 前置条件

- 一台 `Linux x86_64` 服务器（或本机），已安装：
  - **Docker Engine** 及 **Docker Compose v2**（`docker compose version` 能输出 N 开头版本）。
  - 能访问外部镜像仓库（`mysql:8.4`、`redis:7-alpine`、`node:22-alpine`、`golang:1.25-alpine`、`nginx:alpine`、`alpine:3.20`）。
- 一个准备解析到该服务器的**域名**（可选，但强烈建议，用于订阅地址与客户端连接）。
- 内存 ≥ 2GB（四个容器约 1GB+），磁盘 ≥ 10GB。

> **从零 = 无需手工建库**。MySQL 通过 `MYSQL_DATABASE: mewsoproxy` 自动建库，表由后端自动迁移创建，密钥由服务首次启动自动生成，默认管理员自动播种。`v2board-master/` 原始 SQL 已不再需要。

## 2. 放置源码

将本仓库上传到服务器，例如：

```bash
mkdir -p /opt/mewsoproxy && cd /opt/mewsoproxy
# 方式 A：git clone
git clone <你的仓库地址> .
# 方式 B：上传源码压缩包后解压
# unzip mewsoproxy.zip -d .
```

确认关键目录存在：

```bash
ls -1
# docker-compose.yml  server/  web/  DEPLOY.md  README.md
```

## 3. 生产配置（强烈建议先做）

`docker-compose.yml` 中默认值可直接启动，但生产建议修改以下项（直接编辑 `docker-compose.yml` 或通过覆盖环境变量）：

| 配置 | 默认值 | 建议 |
| --- | --- | --- |
| MySQL 密码 | `MYSQL_ROOT_PASSWORD: root`，且后端 `MEWSO_DATABASE_PASSWORD: root` | 两处同时改成强密码 |
| 默认管理员密码 | `MEWSO_ADMIN_PASSWORD: admin123` | 改成强密码 |
| 默认管理员邮箱 | `MEWSO_ADMIN_EMAIL: admin@test.com` | 改成你的邮箱 |
| 订阅地址 | 后台系统中配置 | 见第 6 步 |

> 密钥（JWT / SSH / 节点 token）**无需配置**：留空或为占位值时，服务首次启动会生成随机密钥并写入命名卷 `server_data:/app/run/secrets.env`，重启保持稳定。

## 4. 构建与分发镜像（本地/CI 一次完成，服务器零构建）

> **关键：镜像在本地/CI 构建，绝不要在低配服务器上构建。** 2C2G 服务器只需**加载镜像**即可，这也是 new-api/sub-store 等能在小机器上顺滑部署的原因——服务器只 `pull`/`load`，从不编译。

### 4.1 本地构建并导出离线包

在任一台有 Docker 的机器上（你的开发机/CI）：

```bash
cd mewsoproxy
bash ./docker-build.sh
```

- 产物：`mewsoproxy-images.tar.gz`（内含 `mewsoproxy/server:latest` 与 `mewsoproxy/web:latest` 两个镜像）。
- 若有镜像仓库：把 compose 里 `image` 改成仓库地址后，用 `docker compose build server web && docker compose push`，服务器端 `docker compose pull` 即可。

### 4.2 服务器加载并启动（零构建）

把 `mewsoproxy-images.tar.gz` 上传到服务器（与 `docker-compose.yml`、`deploy.sh` 放一起）：

```bash
cd /opt/mewsoproxy
docker load < mewsoproxy-images.tar.gz
bash ./deploy.sh
```

`deploy.sh` 的决策顺序（**默认永不触发源码构建**）：

1. 本机已存在镜像 → 直接启动；
2. 否则加载同目录 `mewsoproxy-images.tar.gz`；
3. 否则 `docker compose pull`（镜像仓库分发）；
4. 最后 `docker compose up -d`。

> 也可以只上传整个仓库，服务器执行 `bash ./deploy.sh`，脚本会自动优先使用本地镜像。

### 回退：源码构建（仅 4G+ 机器）

```bash
bash ./deploy.sh --build
```

该模式才会在服务器上编译：串行构建 + 自动创建 swap 兜底。**低配机器不要用。**

> 说明：`docker-compose.yml` 的 `server`/`web` 同时声明了 `image:` 与 `build:`，因此 `docker compose up -d`（不带 `--build`）会直接用已加载的镜像，只有显式 `docker compose build` 或 `--build` 才会在服务器上编译。

注意：

- 启动顺序由 Compose 依赖控制：`mysql`、`redis` 先健康，随后 `server`（依赖二者），最后 `web`（依赖 `server` 健康）。
- web 生产构建已默认**跳过 `vue-tsc` 类型检查**（仅生产不需要、只占内存），只执行 `vite build`，进一步降低构建内存。

## 5. 验证服务健康

```bash
# 1) 容器全部在运行（Up / healthy）
docker compose ps

# 2) 后端健康检查
curl -s http://127.0.0.1:8081/healthz

# 3) 关键日志无报错
docker compose logs -f server
# 预期看到：secrets initialized / 自动迁移完成 / SeedAdmin 成功 / server starting
```

预期 `docker compose ps` 输出：

| NAME | STATUS |
| --- | --- |
| `mewsoproxy-mysql` | Up (healthy) |
| `mewsoproxy-redis` | Up (healthy) |
| `mewsoproxy-server` | Up (healthy) |
| `mewsoproxy-web` | Up (healthy) |

## 6. 访问前端并登录后台

```bash
# 本机调试
open http://127.0.0.1:8082
```

- 用默认管理员登录：`admin@test.com` / `admin123`（或你在第 3 步改后的值）。
- 登录成功后进入后台首页，应能看到概览（用户/订单/今日实付/在线/DB/Redis 状态）。

> **进入后台后请立即修改管理员密码**。

## 7. 配置订阅地址（后台改，运行时生效）

1. 后台 → **系统配置**。
2. 将「订阅地址」改为你的公开域名，例如 `https://panel.example.com`（或 `http://<服务器IP>:8082`）。
3. 保存后**重启后端容器**使配置生效：

```bash
docker compose restart server
```

> 服务启动时会从数据库读取 `subscribe_url` / `register_enabled` 并覆盖运行时配置，无需再改 `config.yaml`。

## 8. 域名反向代理（可选但推荐）

将你的域名解析到服务器 IP，然后在 Web 服务器（nginx / Caddy / 宝塔等）把 `https://panel.example.com` 反代到：

```text
http://127.0.0.1:8082
```

因为前端 nginx 已转发 `/api`，所以只需这一个反代即可。若直接用 IP 访问，跳到第 9 步即可。

## 9. 端到端验证清单

| # | 验证项 | 操作 | 预期 |
| --- | --- | --- | --- |
| 1 | 注册开关 | 后台「系统配置」看注册开关 | 开启状态下用户可注册 |
| 2 | 用户注册 | 前端注册页注册新账号 | 成功创建并登录 |
| 3 | 用户下单 | 套餐页购买 | 生成订单 |
| 4 | 订阅链接 | 复制用户订阅地址并访问 | 返回客户端配置内容 |
| 5 | 结算支付 | 下单后走支付弹窗 | 选渠道后生成支付；模拟渠道可一键标记成功 |
| 6 | 订单回填 | 后台查看该订单 | 状态变为已支付，流量/到期已更新 |
| 7 | 订阅到期重置 | 后台对该用户重置流量/到期 | 用户端数据刷新 |

## 10. 常见问题（FAQ)

**Q1：日志提示 `database init failed` / 连不上 MySQL？**
确认 MySQL 容器已 healthy，且 `MEWSO_DATABASE_USER/PASSWORD/HOST` 与 compose 中 `MYSQL_ROOT_PASSWORD` 一致。若改过密码，**两处必须同步**。

**Q2：首次 `docker compose up` 后 `web` 一直启动中？**
`web` 依赖 `server` 的健康状态（`GET /healthz`）。请先看 `server` 是否 healthy；若 `server` 挂了，`web` 会一直等待。

**Q3：改了密码对不上 / 想重置密钥？**
删除命名卷后重新启动即可让系统重新生成密钥与默认管理员：
```bash
docker compose down -v && docker compose up -d --build
```
> `-v` 会清空数据库与密钥卷，**谨慎使用**（用于彻底重置整站）。

**Q4：用 IP 访问时订阅地址填什么？**
填 `http://<服务器IP>:8082`（前端端口）。若你走了域名反代，则填 `https://<你的域名>`。

**Q5：如何观察自动生成了哪些密钥？**
```bash
docker compose exec server cat /app/run/secrets.env
```
生成的文件权限为 `0600`（仅属主可读写），不要提交到版本库。
