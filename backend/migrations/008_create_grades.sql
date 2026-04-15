-- 008_create_grades.sql
-- Owner: Sanjana

CREATE TABLE IF NOT EXISTS grades (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    student_id BIGINT          NOT NULL,
    class_id   BIGINT          NOT NULL REFERENCES classes (id) ON DELETE CASCADE,
    subject    VARCHAR(100)    NOT NULL,
    exam_name  VARCHAR(100)    NOT NULL,
    marks      NUMERIC(6,2)    NOT NULL,
    max_marks  NUMERIC(6,2)    NOT NULL DEFAULT 100,
    teacher_id BIGINT          NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_grades_student_id ON grades (student_id);
CREATE INDEX IF NOT EXISTS idx_grades_class_id   ON grades (class_id);
