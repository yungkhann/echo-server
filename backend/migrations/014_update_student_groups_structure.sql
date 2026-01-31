ALTER TABLE student_groups ADD COLUMN IF NOT EXISTS course_year INTEGER DEFAULT 1;
ALTER TABLE student_groups ADD COLUMN IF NOT EXISTS academic_year VARCHAR(20) DEFAULT '2025-2026';
ALTER TABLE student_groups ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE student_groups ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

UPDATE student_groups 
SET created_at = CURRENT_TIMESTAMP 
WHERE created_at IS NULL;

UPDATE student_groups 
SET updated_at = CURRENT_TIMESTAMP 
WHERE updated_at IS NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'fk_faculty' AND table_name = 'student_groups'
    ) THEN
        ALTER TABLE student_groups 
        ADD CONSTRAINT fk_faculty 
        FOREIGN KEY (faculty_id) REFERENCES faculties(id) ON DELETE CASCADE;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_student_groups_faculty ON student_groups(faculty_id);
CREATE INDEX IF NOT EXISTS idx_student_groups_name ON student_groups(group_name);
CREATE INDEX IF NOT EXISTS idx_student_groups_year ON student_groups(course_year);
