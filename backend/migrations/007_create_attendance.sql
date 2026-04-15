-- 007_create_attendance.sql
-- Owner: Sanjana

CREATE TYPE attendance_status AS ENUM ('present', 'absent', 'late');

CREATE TABLE IF NOT EXISTS attendances (
    id           BIGSERIAL PRIMARY KEY,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    class_id     BIGINT NOT NULL REFERENCES classes (id) ON DELETE CASCADE,
    student_id   BIGINT NOT NULL,
    date         TIMESTAMPTZ NOT NULL,
    status       attendance_status NOT NULL DEFAULT 'present',
    marked_by_id BIGINT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_attendances_class_id   ON attendances (class_id);
CREATE INDEX IF NOT EXISTS idx_attendances_student_id ON attendances (student_id);
CREATE INDEX IF NOT EXISTS idx_attendances_date       ON attendances (date);
