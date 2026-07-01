# PCMS - 个人认知管理系统

基于 AI + RAG 的个人知识库，帮助思考、整理、检索、生成和迭代。

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3 + TypeScript + Element Plus + Pinia |
| 后端 | Go + Gin + GORM |
| 数据库 | PostgreSQL 16 |
| 认证 | JWT |

---

## Docker 部署（推荐）

### 前置条件

- Docker 20.10+ 和 Docker Compose v2

### 一键部署

```bash
# Linux / macOS
chmod +x docker.sh
./docker.sh

# Windows PowerShell
.\docker.ps1
```

部署成功后访问 `http://localhost:83`，首次使用需注册账号。

### Docker 脚本命令

```bash
./docker.sh [命令]
```

| 命令 | 说明 |
|------|------|
| `up`（默认） | 构建镜像并启动所有服务 |
| `build` | 仅构建镜像（不启动） |
| `start` | 启动已有服务 |
| `stop` | 停止所有服务 |
| `restart` | 重启所有服务 |
| `logs [容器名]` | 查看日志，如 `./docker.sh logs server` |
| `status` | 查看服务运行状态 |
| `exec [容器名]` | 进入容器 shell，默认 server |

Windows 下将 `./docker.sh` 替换为 `.\docker.ps1`。

### 架构

```
浏览器 → :83 (nginx)  → 静态文件
                      → /api/* → Go Server :8080 → PostgreSQL :5432
```

3 个容器：`pcms-web`、`pcms-server`、`pcms-db`

### 端口配置

如需修改端口，直接设置环境变量或编辑 `.env` 后再启动：

```bash
# 例如：数据库用 7532，前端用 8888
export DB_PORT=7532
export WEB_PORT=8888
./docker.sh
```

### 数据持久化

- `pgdata` 卷 — 数据库文件
- `uploads` 卷 — 上传文件

重建容器不会丢失数据。如需彻底清除：

```bash
docker compose down -v
```

### 代码更新后重新部署

```bash
# 方式一：先拉代码再部署
./scripts/pull.sh deploy

# 方式二：本地直接重新构建
./docker.sh
```

---

## 从 GitHub 部署

```bash
# 设置仓库地址
export GITHUB_REPO="https://github.com/你的用户名/pcms.git"

# 克隆
./scripts/pull.sh clone

# 拉取最新代码 + 自动部署
./scripts/pull.sh deploy
```

---

## 本地开发

### 1. 启动数据库

```bash
docker compose up db -d
```

### 2. 启动后端

```bash
cd server
go mod tidy
go run main.go
```

服务在 `http://localhost:8080`

### 3. 启动前端

```bash
cd web
npm install
npm run dev
```

开发服务在 `http://localhost:3000`，API 自动代理到后端。

---

## 项目结构

```
pcms/
├── design.md                # 需求细则
├── docker-compose.yml       # Docker 编排（db + server + web）
├── init.sql                 # 数据库扩展初始化
├── docker.sh / docker.ps1   # Docker 管理脚本
├── deploy.sh / deploy.ps1   # 完整部署脚本
├── scripts/
│   ├── pull.sh              # GitHub 拉取脚本
│   └── pull.ps1             # GitHub 拉取脚本 (Windows)
├── server/                  # Go 后端
│   ├── main.go
│   ├── Dockerfile
│   ├── config/              # 配置
│   ├── models/              # 数据模型
│   ├── middleware/           # JWT、CORS
│   ├── handlers/            # 请求处理
│   ├── services/            # 业务逻辑
│   ├── routes/              # 路由
│   └── utils/               # 工具
└── web/                     # Vue 3 前端
    ├── Dockerfile
    ├── nginx.conf            # Nginx 配置
    └── src/
        ├── api/              # API 请求层
        ├── router/           # 路由 + 守卫
        ├── stores/           # Pinia 状态
        ├── views/            # 页面
        └── components/       # 组件
```

## MVP 功能

- [x] 用户注册/登录 (JWT)
- [x] 无限级分类管理（树形 + 拖拽）
- [x] 文档 CRUD（Markdown 编辑 + 预览）
- [x] 文件上传（图片）
- [x] 全文搜索
- [x] 标签系统
- [x] 文档版本管理
- [ ] AI 润色（下一步）
- [ ] RAG 检索（下一步）

## API 概览

| Method | Path | 说明 |
|--------|------|------|
| POST | /api/v1/auth/register | 注册 |
| POST | /api/v1/auth/login | 登录 |
| GET | /api/v1/auth/me | 当前用户 |
| GET | /api/v1/categories | 分类树 |
| POST | /api/v1/categories | 创建分类 |
| PUT | /api/v1/categories/:id | 更新分类 |
| DELETE | /api/v1/categories/:id | 删除分类 |
| PUT | /api/v1/categories/:id/move | 移动分类 |
| GET | /api/v1/documents | 文档列表 |
| POST | /api/v1/documents | 创建文档 |
| GET | /api/v1/documents/:id | 文档详情 |
| PUT | /api/v1/documents/:id | 更新文档 |
| DELETE | /api/v1/documents/:id | 删除文档 |
| GET | /api/v1/documents/search | 搜索 |
| POST | /api/v1/files/upload | 上传文件 |
| GET/POST | /api/v1/tags | 标签管理 |
