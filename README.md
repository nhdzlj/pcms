# PCMS - 个人认知管理系统

基于 AI + RAG 的个人知识库，帮助思考、整理、检索、生成和迭代。

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3 + TypeScript + Element Plus + Pinia |
| 后端 | Go + Gin + GORM（Store 抽象层，支持 MemStore 测试） |
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

### 4. 运行测试

```bash
cd server
go test ./...
```

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
│   ├── store/               # 数据访问抽象层（GORM / MemStore）
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
- [x] 摘要编辑
- [x] 标签系统（选择 + 新建 + 筛选）
- [x] 文件上传（图片 + 附件管理）
- [x] 全文搜索（支持分类/标签筛选）
- [x] 文档版本管理（历史查看）
- [ ] AI 润色（下一步）
- [ ] RAG 检索（下一步）

## 单元测试

后端全部接口均有单元测试覆盖，无需外部数据库依赖（使用内存 Store）。

```bash
# 运行所有测试
cd server
go test ./...

# 按模块运行
go test ./services   # 业务逻辑测试
go test ./utils      # 工具测试（JWT / 分页 / 响应）
go test .            # HTTP 接口测试（Handler 层）
```

| 测试范围 | 测试文件 | 覆盖接口 |
|---------|---------|---------|
| Auth | `services/auth_test.go` + `handler_test.go` | 注册、登录、当前用户 |
| Category | `services/category_test.go` + `handler_test.go` | CRUD、移动、树形结构 |
| Document | `services/document_test.go` + `handler_test.go` | CRUD、搜索、版本、标签筛选 |
| Tag | `services/tag_test.go` + `handler_test.go` | 创建、去重、删除、隔离 |
| Attachment | `services/attachment_test.go` + `handler_test.go` | 附件 CRUD、绑定文档 |
| File | `handler_test.go` | 文件上传、静态服务、目录遍历防护 |
| Utils | `utils/utils_test.go` | JWT 签发/解析、分页、响应格式 |

## API 概览

### 认证
| Method | Path | 说明 |
|--------|------|------|
| POST | /api/v1/auth/register | 注册 |
| POST | /api/v1/auth/login | 登录 |
| GET | /api/v1/auth/me | 当前用户 |

### 分类
| Method | Path | 说明 |
|--------|------|------|
| GET | /api/v1/categories | 分类树 |
| POST | /api/v1/categories | 创建分类 |
| PUT | /api/v1/categories/:id | 更新分类 |
| DELETE | /api/v1/categories/:id | 删除分类 |
| PUT | /api/v1/categories/:id/move | 移动分类 |

### 文档
| Method | Path | 说明 |
|--------|------|------|
| GET | /api/v1/documents | 文档列表（支持 `category_id`、`status`、`tag_id` 筛选） |
| POST | /api/v1/documents | 创建文档 |
| GET | /api/v1/documents/search | 搜索（支持 `keyword`、`tag_id`、`category_id` 筛选） |
| GET | /api/v1/documents/:id | 文档详情 |
| PUT | /api/v1/documents/:id | 更新文档 |
| DELETE | /api/v1/documents/:id | 删除文档 |
| GET | /api/v1/documents/:id/versions | 版本列表 |
| GET | /api/v1/documents/:id/versions/:vid | 版本详情 |

### 标签
| Method | Path | 说明 |
|--------|------|------|
| GET | /api/v1/tags | 标签列表 |
| POST | /api/v1/tags | 创建标签 |
| DELETE | /api/v1/tags/:id | 删除标签 |

### 文件 & 附件
| Method | Path | 说明 |
|--------|------|------|
| POST | /api/v1/files/upload | 上传文件 |
| GET | /api/v1/attachments | 附件列表 |
| POST | /api/v1/attachments | 创建附件记录 |
| DELETE | /api/v1/attachments/:id | 删除附件 |
| PUT | /api/v1/attachments/:id/bind | 绑定附件到文档 |
