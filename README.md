# PCMS - 个人认知管理系统

基于 AI + RAG 的个人知识库，帮助思考、整理、检索、生成和迭代。

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3 + TypeScript + Element Plus + Pinia |
| 后端 | Go + Gin + GORM |
| 数据库 | PostgreSQL 16 |
| 认证 | JWT |

## 快速开始

### 1. 启动数据库

```bash
docker-compose up -d
```

### 2. 启动后端

```bash
cd server
go mod tidy
go run main.go
```

服务启动在 `http://localhost:8080`

### 3. 启动前端

```bash
cd web
npm install
npm run dev
```

开发服务启动在 `http://localhost:3000`

## 项目结构

```
pcms/
├── design.md              # 开发需求细则
├── docker-compose.yml      # PostgreSQL 容器
├── init.sql                # 数据库初始化脚本
├── server/                 # Go 后端
│   ├── main.go
│   ├── config/             # 配置
│   ├── models/             # 数据模型
│   ├── middleware/          # 中间件 (JWT, CORS)
│   ├── handlers/           # 请求处理
│   ├── services/           # 业务逻辑
│   ├── routes/             # 路由
│   └── utils/              # 工具函数
└── web/                    # Vue 3 前端
    └── src/
        ├── api/            # API 请求
        ├── router/         # 路由
        ├── stores/         # Pinia 状态
        ├── views/          # 页面
        └── components/     # 组件
```

## MVP 功能

- [x] 用户注册/登录 (JWT)
- [x] 无限级分类管理 (树形结构 + 拖拽)
- [x] 文档 CRUD (Markdown 编辑)
- [x] 文件上传 (图片)
- [x] 全文搜索 (标题 + 内容)
- [x] 标签系统
- [x] 文档版本管理
- [ ] AI 润色 (下一步)
- [ ] RAG 检索 (下一步)

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
