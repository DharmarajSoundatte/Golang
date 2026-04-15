-- 006_create_classes.sql
-- Owner: Sanjana

CREATE TABLE IF NOT EXISTS classes (
    id                BIGSERIAL PRIMARY KEY,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ,

    name              VARCHAR(100) NOT NULL,
    section           VARCHAR(10)  NOT NULL,
    class_teacher_id  BIGINT REFERENCES teachers (id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_classes_deleted_at      ON classes (deleted_at);
CREATE INDEX IF NOT EXISTS idx_classes_class_teacher   ON classes (class_teacher_id);

-- Tracks which students belong to which class (one student → one class at a time)
CREATE TABLE IF NOT EXISTS class_students (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    class_id   BIGINT NOT NULL REFERENCES classes (id) ON DELETE CASCADE,
    student_id BIGINT NOT NULL,

    CONSTRAINT idx_class_student UNIQUE (student_id)
);

CREATE INDEX IF NOT EXISTS idx_class_students_class ON class_students (class_id);
