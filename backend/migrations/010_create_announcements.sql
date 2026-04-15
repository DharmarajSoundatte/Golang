-- 010_create_announcements.sql
-- Owner: Sanjana

CREATE TYPE target_audience AS ENUM ('all', 'class', 'role');

CREATE TABLE IF NOT EXISTS announcements (
    id              BIGSERIAL PRIMARY KEY,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,

    title           VARCHAR(255)    NOT NULL,
    content         TEXT            NOT NULL,
    author_id       BIGINT          NOT NULL,
    target_audience target_audience NOT NULL DEFAULT 'all',
    target_class_id BIGINT          REFERENCES classes (id) ON DELETE SET NULL,
    target_role     VARCHAR(50)
);

CREATE INDEX IF NOT EXISTS idx_announcements_deleted_at ON announcements (deleted_at);
