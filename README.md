# MewsoProxy

基于 V2Board 数据模型重构的面板系统：Go 后端 + TDesign 前端。采用 Go (Gin + GORM + MySQL + Redis + JWT) 重写后端，前端使用 Vue 3.5 + tdesign-vue-next 构建。

## 技术栈

- **后端**：Go 1.22+、Gin、GORM、MySQL 8.4、Redis、JWT（双 Token）
- **前端**：Vue 3.5、TypeScript、Vite、Pinia、tdesign-vue-next
- **部署**：Docker Compose

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
│   ├── model/           # GORM 模型（复用 v2_* 表）
│   ├── pkg/             # 响应/错误码/Token/Redis 等公共包
│   ├── router/          # 路由注册（/api/v1）
│   ├── service/         # 业务逻辑层
│   ├── config.yaml      # 运行时配置
│   └── main.go
├── web/                 # TDesign 前端
│   └── src/
│       ├── api/         # 接口封装（用户端 + 管理端 admin）
│       ├── layouts/     # BasicLayout / AdminLayout
│       ├── pages/       # 用户端 + 管理端页面
│       ├── router/      # 前端路由与权限守卫
│       ├── store/       # Pinia 状态
│       └── utils/       # request / token
├── docker-compose.yml   # MySQL + Redis + Server + Web 编排
└── .gitignore
```

## 功能进度

### 阶段一：基础骨架与用户端链路
- 用户注册 / 登录（Access JWT + Refresh httpOnly Cookie）
- 套餐列表与购买下单
- 订阅链接生成与轻量客户端下发
- 用户首页（基础布局）

### 阶段二：管理后台
- 管理员登录（仅管理员中间件拦截）
- 概览：用户总数 / 订单总数 / 今日实付 / 在线用户 / 数据库与 Redis 状态
- 用户管理：列表 / 详情 / 编辑 / 封禁 / 重置密钥
- 套餐管理：列表 / 新增 / 编辑 / 删除
- 订单管理：搜索 / 分页 / 详情 / 标记支付 / 取消
- 节点管理：节点分组 + 按协议（trojan / vmess / shadowsocks / hysteria）管理
- 优惠券管理：新增 / 编辑 / 删除 / 展示切换
- 公告管理：新增 / 编辑 / 删除 / 展示切换
- 支付渠道管理：新增 / 编辑 / 删除 / 启用切换
- 系统配置：站点名 / 注册开关 / 订阅地址 / secure_path 等

## 启动方式

### Docker Compose（推荐）

```bash
docker compose up -d
```

- 启动后包含 3 个服务：MySQL、Redis、后端（`server`）。
- **前端已用 `go:embed` 内嵌进后端二进制**，因此无需独立的 web/Nginx 容器。后端容器映射到宿主机 `8082` 端口，同时托管「前端静态页 + `/api/v1` API」，域名反向代理指向 `http://<服务器IP>:8082` 即可。
- 服务器**不构建前端**（vite/node 已在本地完成），2C2G 也只需一次 Go 编译，不会 OOM。
- **密钥自动生成**：首次启动会自动生成 JWT / SSH / 节点 token 密钥并写入命名卷 `server_data:/app/run/secrets.env`，重启保持稳定；也可通过环境变量 `MEWSO_JWT_ACCESS_SECRET` 等固定。
- **默认管理员**：数据库无管理员时自动创建账号（默认 `admin@test.com` / `admin123`，可用 `MEWSO_ADMIN_EMAIL` / `MEWSO_ADMIN_PASSWORD` 覆盖），请务必在后台修改密码。
- **订阅地址**：登录后台在「系统配置」中修改「订阅地址」为你的公开域名后生效（重启容器后生效）。

> 低配服务器建议直接执行 `bash ./deploy.sh`（会自适应：有镜像用镜像、否则加载离线包、否则自动源码构建仅编译 Go）。详细步骤见 [DEPLOY.md](DEPLOY.md)。

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

- 业务路由统一挂在 `/api/v1` 前缀下。
- 统一响应结构：`{ "code": 0, "message": "ok", "data": ... }`，`code` 非 0 表示业务错误。
- UTC 落库，API 层按 RFC3339 输出时间。

## 默认管理员

| 项目 | 值 |
| --- | --- |
| 账号 | `admin@test.com` |
| 密码 | `admin123` |

> 生产环境请务必修改默认管理员密码，并通过系统配置调整站点信息与订阅地址。
