-- 005_create_teachers.sql
-- Owner: Sanjana

CREATE TABLE IF NOT EXISTS teachers (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    name       VARCHAR(255) NOT NULL,
    email      VARCHAR(255) NOT NULL,
    phone      VARCHAR(20),
    subject    VARCHAR(255),
    is_active  BOOLEAN NOT NULL DEFAULT TRUE,

    CONSTRAINT teachers_email_unique UNIQUE (email)
);

CREATE INDEX IF NOT EXISTS idx_teachers_deleted_at ON teachers (deleted_at);
