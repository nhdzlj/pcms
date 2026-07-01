-- PCMS 数据库初始化脚本
-- 注意：所有表结构由 GORM AutoMigrate 管理，此脚本仅初始化扩展

-- 全文搜索扩展（后续 RAG 也会用到）
CREATE EXTENSION IF NOT EXISTS pg_trgm;
