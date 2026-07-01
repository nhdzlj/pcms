# PCMS 开发需求细则 & MVP 实现方案

## 一、MVP 范围

### 1.1 包含模块

| 模块 | 说明 | 优先级 |
|------|------|--------|
| 用户认证 | JWT 登录/注册 | P0 |
| 分类管理 | 无限级树形目录、拖拽排序 | P0 |
| 文档管理 | 富文本/Markdown编辑、CRUD | P0 |
| 文件上传 | 图片、附件上传 | P0 |
| 标签系统 | 文档标签管理 | P1 |
| 全文搜索 | 标题+内容搜索 | P1 |
| 版本管理 | 文档历史版本 | P2 |
| AI 润色 | 接入 LLM 进行内容优化 | P2 |
| RAG 检索 | Embedding + 向量检索 | P2 |

### 1.2 暂不包含

- 知识图谱可视化
- AI 主动提醒
- 语音/OCR/网页抓取
- 移动端 PWA
- 团队协作/权限发布

---

## 二、技术选型

| 层级 | 技术 | 版本 |
|------|------|------|
| 前端框架 | Vue 3 + TypeScript | ^3.4 |
| 构建工具 | Vite | ^5 |
| UI 组件 | Element Plus | ^2.5 |
| 富文本编辑器 | Tiptap | ^2.x |
| Markdown 编辑器 | @kangc/v-md-editor | ^3.x |
| 状态管理 | Pinia | ^2.1 |
| 路由 | Vue Router | ^4.2 |
| HTTP 客户端 | Axios | ^1.6 |
| 后端框架 | Gin | ^1.9 |
| ORM | GORM | ^1.25 |
| 数据库 | PostgreSQL 16 | - |
| 认证 | JWT (golang-jwt) | ^5 |
| 文件存储 | 本地存储 (MVP) / MinIO (后续) | - |
| 向量扩展 | pgvector (后续) | - |

---

## 三、数据库设计

```sql
-- 用户表
CREATE TABLE users (
    id          BIGSERIAL PRIMARY KEY,
    username    VARCHAR(64)  UNIQUE NOT NULL,
    password    VARCHAR(256) NOT NULL,
    email       VARCHAR(128),
    avatar      VARCHAR(512),
    created_at  TIMESTAMPTZ  DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  DEFAULT NOW()
);

-- 分类表（无限层级）
CREATE TABLE categories (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(128) NOT NULL,
    parent_id   BIGINT       REFERENCES categories(id) ON DELETE CASCADE,
    sort_order  INTEGER      DEFAULT 0,
    user_id     BIGINT       NOT NULL REFERENCES users(id),
    icon        VARCHAR(64),
    created_at  TIMESTAMPTZ  DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  DEFAULT NOW()
);

CREATE INDEX idx_categories_parent ON categories(parent_id);
CREATE INDEX idx_categories_user ON categories(user_id);

-- 文档表
CREATE TABLE documents (
    id          BIGSERIAL PRIMARY KEY,
    title       VARCHAR(256) NOT NULL,
    content     TEXT,
    summary     TEXT,
    category_id BIGINT       REFERENCES categories(id) ON DELETE SET NULL,
    user_id     BIGINT       NOT NULL REFERENCES users(id),
    status      VARCHAR(20)  DEFAULT 'draft',
    view_count  INTEGER      DEFAULT 0,
    is_favorite BOOLEAN      DEFAULT FALSE,
    version     INTEGER      DEFAULT 1,
    created_at  TIMESTAMPTZ  DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  DEFAULT NOW()
);

CREATE INDEX idx_documents_category ON documents(category_id);
CREATE INDEX idx_documents_user ON documents(user_id);
CREATE INDEX idx_documents_title_content ON documents USING gin(to_tsvector('simple', title || ' ' || coalesce(content, '')));

-- 标签表
CREATE TABLE tags (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(64)  NOT NULL,
    user_id     BIGINT       NOT NULL REFERENCES users(id),
    UNIQUE(name, user_id)
);

-- 文档-标签关联表
CREATE TABLE document_tags (
    document_id BIGINT NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    tag_id      BIGINT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (document_id, tag_id)
);

-- 文档版本表
CREATE TABLE document_versions (
    id          BIGSERIAL PRIMARY KEY,
    document_id BIGINT  NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    version     INTEGER NOT NULL,
    title       VARCHAR(256),
    content     TEXT,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_doc_versions_doc ON document_versions(document_id);

-- 附件表
CREATE TABLE attachments (
    id          BIGSERIAL PRIMARY KEY,
    document_id BIGINT       REFERENCES documents(id) ON DELETE SET NULL,
    file_name   VARCHAR(256),
    file_path   VARCHAR(512),
    file_size   BIGINT,
    mime_type   VARCHAR(128),
    user_id     BIGINT       NOT NULL REFERENCES users(id),
    created_at  TIMESTAMPTZ  DEFAULT NOW()
);
```

---

## 四、API 设计

### 4.1 认证

| Method | Path | 说明 |
|--------|------|------|
| POST | /api/v1/auth/register | 注册 |
| POST | /api/v1/auth/login | 登录 |
| GET | /api/v1/auth/me | 获取当前用户 |

### 4.2 分类

| Method | Path | 说明 |
|--------|------|------|
| GET | /api/v1/categories | 获取分类树 |
| POST | /api/v1/categories | 创建分类 |
| PUT | /api/v1/categories/:id | 更新分类 |
| DELETE | /api/v1/categories/:id | 删除分类 |
| PUT | /api/v1/categories/:id/move | 移动分类 |

### 4.3 文档

| Method | Path | 说明 |
|--------|------|------|
| GET | /api/v1/documents | 文档列表（分页、分类筛选） |
| POST | /api/v1/documents | 创建文档 |
| GET | /api/v1/documents/:id | 文档详情 |
| PUT | /api/v1/documents/:id | 更新文档 |
| DELETE | /api/v1/documents/:id | 删除文档 |
| GET | /api/v1/documents/search | 全文搜索 |

### 4.4 文件

| Method | Path | 说明 |
|--------|------|------|
| POST | /api/v1/files/upload | 上传文件 |

### 4.5 标签

| Method | Path | 说明 |
|--------|------|------|
| GET | /api/v1/tags | 获取标签列表 |
| POST | /api/v1/tags | 创建标签 |

---

## 五、前端路由

```
/login                          # 登录页
/                               # 首页（最近文档）
/documents                      # 文档列表
/documents/:id                  # 文档详情/编辑器
/categories                     # 分类管理
/search                         # 搜索页
```

---

## 六、目录结构

```
pcms/
├── design.md                   # 本文档
├── docker-compose.yml          # PostgreSQL 容器
├── server/                     # Go 后端
│   ├── main.go
│   ├── go.mod
│   ├── config/
│   │   └── config.go
│   ├── models/
│   │   ├── user.go
│   │   ├── category.go
│   │   ├── document.go
│   │   ├── tag.go
│   │   └── attachment.go
│   ├── middleware/
│   │   ├── auth.go
│   │   └── cors.go
│   ├── handlers/
│   │   ├── auth.go
│   │   ├── category.go
│   │   ├── document.go
│   │   ├── tag.go
│   │   └── file.go
│   ├── services/
│   │   ├── auth.go
│   │   ├── category.go
│   │   ├── document.go
│   │   └── tag.go
│   ├── routes/
│   │   └── router.go
│   └── utils/
│       ├── jwt.go
│       ├── response.go
│       └── pagination.go
├── web/                        # Vue 3 前端
│   ├── index.html
│   ├── package.json
│   ├── vite.config.ts
│   ├── tsconfig.json
│   └── src/
│       ├── main.ts
│       ├── App.vue
│       ├── api/
│       │   ├── request.ts
│       │   ├── auth.ts
│       │   ├── category.ts
│       │   ├── document.ts
│       │   ├── tag.ts
│       │   └── file.ts
│       ├── router/
│       │   └── index.ts
│       ├── stores/
│       │   ├── auth.ts
│       │   ├── category.ts
│       │   └── document.ts
│       ├── views/
│       │   ├── LoginView.vue
│       │   ├── LayoutView.vue
│       │   ├── HomeView.vue
│       │   ├── DocumentList.vue
│       │   ├── DocumentEditor.vue
│       │   └── SearchView.vue
│       ├── components/
│       │   ├── AppHeader.vue
│       │   ├── AppSidebar.vue
│       │   ├── CategoryTree.vue
│       │   ├── DocumentCard.vue
│       │   └── MarkdownEditor.vue
│       └── styles/
│           └── global.scss
└── init.sql                    # 数据库初始化脚本
```
