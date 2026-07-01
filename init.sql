-- PCMS 数据库初始化脚本
-- PostgreSQL 16

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id          BIGSERIAL PRIMARY KEY,
    username    VARCHAR(64)  UNIQUE NOT NULL,
    password    VARCHAR(256) NOT NULL,
    email       VARCHAR(128),
    avatar      VARCHAR(512),
    created_at  TIMESTAMPTZ  DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  DEFAULT NOW()
);

-- 分类表（无限层级，通过 parent_id 建立树形结构）
CREATE TABLE IF NOT EXISTS categories (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(128) NOT NULL,
    parent_id   BIGINT       REFERENCES categories(id) ON DELETE CASCADE,
    sort_order  INTEGER      DEFAULT 0,
    user_id     BIGINT       NOT NULL REFERENCES users(id),
    icon        VARCHAR(64)  DEFAULT 'folder',
    created_at  TIMESTAMPTZ  DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_categories_parent ON categories(parent_id);
CREATE INDEX IF NOT EXISTS idx_categories_user ON categories(user_id);

-- 文档表
CREATE TABLE IF NOT EXISTS documents (
    id          BIGSERIAL PRIMARY KEY,
    title       VARCHAR(256) NOT NULL,
    content     TEXT         DEFAULT '',
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

CREATE INDEX IF NOT EXISTS idx_documents_category ON documents(category_id);
CREATE INDEX IF NOT EXISTS idx_documents_user ON documents(user_id);
CREATE INDEX IF NOT EXISTS idx_documents_updated ON documents(updated_at DESC);

-- 全文搜索索引（需要先启用 pg_trgm 扩展）
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX IF NOT EXISTS idx_documents_title_trgm ON documents USING gin(title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_documents_content_trgm ON documents USING gin(content gin_trgm_ops);

-- 标签表
CREATE TABLE IF NOT EXISTS tags (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(64) NOT NULL,
    user_id     BIGINT      NOT NULL REFERENCES users(id),
    UNIQUE(name, user_id)
);

CREATE INDEX IF NOT EXISTS idx_tags_user ON tags(user_id);

-- 文档-标签关联表
CREATE TABLE IF NOT EXISTS document_tags (
    document_id BIGINT NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    tag_id      BIGINT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (document_id, tag_id)
);

-- 文档版本表
CREATE TABLE IF NOT EXISTS document_versions (
    id          BIGSERIAL PRIMARY KEY,
    document_id BIGINT       NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    version     INTEGER      NOT NULL,
    title       VARCHAR(256),
    content     TEXT,
    created_at  TIMESTAMPTZ  DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_doc_versions_doc ON document_versions(document_id);

-- 附件表
CREATE TABLE IF NOT EXISTS attachments (
    id          BIGSERIAL PRIMARY KEY,
    document_id BIGINT       REFERENCES documents(id) ON DELETE SET NULL,
    file_name   VARCHAR(256),
    file_path   VARCHAR(512),
    file_size   BIGINT,
    mime_type   VARCHAR(128),
    user_id     BIGINT       NOT NULL REFERENCES users(id),
    created_at  TIMESTAMPTZ  DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_attachments_doc ON attachments(document_id);
