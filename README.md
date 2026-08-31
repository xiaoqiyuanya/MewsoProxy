# MewsoProxy

基于 Go 与 Vue 3 的全栈面板系统，包含用户端与管理端。

## 技术栈

* **后端**：Go 1.22+、Gin、GORM、MySQL 8.4、Redis、JWT（双 Token）

* **前端**：Vue 3.5、TypeScript、Vite、Pinia

* **部署**：Docker Compose

## 目录结构

```
MewsoProxy/
├── server/              # Go 后端
│   ├── cmd/             # 入口引导（应用/配置/数据库/Redis）
│   ├── config/          # 配置定义
│   ├── database/        # DB 初始化与迁移
│   ├── dto/             # 数据传输对象
│   ├── handler/         # HTTP 处理器
│   ├── middleware/      # 认证/限流/幂等/日志等中间件
│   ├── model/           # GORM 模型
│   ├── pkg/             # 响应/错误码/Token/Redis 等公共包
│   ├── router/          # 路由注册（/api/v1）
│   ├── service/         # 业务逻辑层
│   ├── config.yaml      # 运行时配置
│   └── main.go
├── web/                 # 前端
│   └── src/
│       ├── api/         # 接口封装（用户端 + 管理端 admin）
│       ├── layouts/     # BasicLayout / AdminLayout
│       ├── pages/       # 用户端 + 管理端页面
│       ├── router/      # 前端路由与权限守卫
│       ├── store/       # Pinia 状态
│       └── utils/       # request / token
├── docker-compose.yml   # MySQL + Redis + Server 编排
└── .gitignore
```

## 功能进度

### 阶段一：基础骨架与用户端链路

* 用户注册 / 登录（Access JWT + Refresh httpOnly Cookie）

* 套餐列表与购买下单

* 订阅链接生成与轻量客户端下发

* 用户首页（基础布局）

### 阶段二：管理后台

* 管理员登录（仅管理员中间件拦截）

* 概览：用户总数 / 订单总数 / 今日实付 / 在线用户 / 数据库与 Redis 状态

* 用户管理：列表 / 详情 / 编辑 / 封禁 / 重置密钥

* 套餐管理：列表 / 新增 / 编辑 / 删除

* 订单管理：搜索 / 分页 / 详情 / 标记支付 / 取消

* 节点管理：节点分组 + 按协议（trojan / vmess / shadowsocks / hysteria）管理

* 优惠券管理：新增 / 编辑 / 删除 / 展示切换

* 公告管理：新增 / 编辑 / 删除 / 展示切换

* 支付渠道管理：新增 / 编辑 / 删除 / 启用切换

* 系统配置：站点名 / 注册开关 / 订阅地址 / secure\_path 等

## 启动方式

### Docker Compose（推荐）

从 ghcr.io 拉取镜像的完整命令：

```bash
# 1. 配置镜像（可选，默认使用 mewsoproxy/server:latest）
echo "MEWSO_SERVER_IMAGE=ghcr.io/<你的用户名>/mewsoproxy-server:latest" >> .env

# 2. 私有镜像需登录 ghcr（公开镜像可跳过）
docker login ghcr.io -u <你的用户名>   # 输入 GitHub 用户名 + PAT

# 3. 拉取镜像并启动
docker compose pull
docker compose up -d
```

> 镜像由 GitHub Actions 在推送代码后自动构建并发布到 ghcr.io，服务器环境无需本地构建（建议 2C2G 及更低配置直接采用此方式，避免构建期内存不足）。

* 启动后包含 3 个服务：MySQL、Redis、后端（`server`）。前端静态资源内嵌于后端，无需独立的 web/Nginx 容器，后端容器映射到宿主机 `8082` 端口，同时托管前端页面与 `/api/v1` API。

* **密钥自动生成**：首次启动自动生成 JWT / SSH / 节点 token 密钥并持久化，重启保持稳定；也可通过 `MEWSO_JWT_ACCESS_SECRET` 等环境变量固定。

* **默认管理员**：数据库无管理员时自动创建 `admin@test.com / admin123`（可用 `MEWSO_ADMIN_EMAIL` / `MEWSO_ADMIN_PASSWORD` 覆盖），请登录后立即修改密码。

* **订阅地址**：登录后台在「系统配置」中修改为你的公开域名，重启容器后生效。

### 本地开发

```bash
# 后端
cd server
go mod vendor
go run -mod=vendor .

# 前端
cd web
npm install
npm run dev
```

## 接口约定

* 业务路由统一挂在 `/api/v1` 前缀下。

* 统一响应结构：`{ "code": 0, "message": "ok", "data": ... }`，`code` 非 0 表示业务错误。

* UTC 落库，API 层按 RFC3339 输出时间。

## 默认管理员

| 项目 | 值                |
| -- | ---------------- |
| 账号 | `admin@test.com` |
| 密码 | `admin123`       |

> 生产环境请务必修改默认管理员密码，并通过系统配置调整站点信息与订阅地址。

## 许可与声明

* 本项目基于 **GNU AGPL-3.0** 协议发布（见 [LICENSE](LICENSE)）。任何对本项目的使用、修改、分发，或将其（含修改版）作为网络服务对外提供，均须完整遵循该协议并公开相应源代码。

* 本项目**仅供学习、研究与技术交流使用**，请勿将其用于任何违反所在国家/地区法律法规的用途；由此产生的一切法律责任与后果由使用者自行承担。

* 对于任何人因使用、滥用或分发本项目而产生的任何直接或间接损失与法律责任，项目作者不承担任何责任。
