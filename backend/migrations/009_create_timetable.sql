-- 009_create_timetable.sql
-- Owner: Sanjana

CREATE TABLE IF NOT EXISTS timetables (
    id            BIGSERIAL PRIMARY KEY,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    class_id      BIGINT       NOT NULL REFERENCES classes (id) ON DELETE CASCADE,
    day_of_week   VARCHAR(20)  NOT NULL,
    period_number INT          NOT NULL,
    subject       VARCHAR(100) NOT NULL,
    teacher_id    BIGINT       NOT NULL REFERENCES teachers (id) ON DELETE RESTRICT,
    start_time    VARCHAR(10)  NOT NULL,
    end_time      VARCHAR(10)  NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_timetables_class_id ON timetables (class_id);
