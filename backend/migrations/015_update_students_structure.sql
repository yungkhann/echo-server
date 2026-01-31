ALTER TABLE students ADD COLUMN IF NOT EXISTS email VARCHAR(255);
ALTER TABLE students ADD COLUMN IF NOT EXISTS phone VARCHAR(50);
ALTER TABLE students ADD COLUMN IF NOT EXISTS student_id_number VARCHAR(50);
ALTER TABLE students ADD COLUMN IF NOT EXISTS enrollment_date DATE DEFAULT CURRENT_DATE;
ALTER TABLE students ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'Active';
ALTER TABLE students ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE students ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

UPDATE students 
SET created_at = CURRENT_TIMESTAMP 
WHERE created_at IS NULL;

UPDATE students 
SET updated_at = CURRENT_TIMESTAMP 
WHERE updated_at IS NULL;

UPDATE students 
SET enrollment_date = CURRENT_DATE 
WHERE enrollment_date IS NULL;

UPDATE students 
SET student_id_number = 'STU' || TO_CHAR(CURRENT_DATE, 'YYYY') || LPAD(id::TEXT, 4, '0')
WHERE student_id_number IS NULL;

ALTER TABLE students ADD CONSTRAINT IF NOT EXISTS students_email_key UNIQUE (email);
ALTER TABLE students ADD CONSTRAINT IF NOT EXISTS students_student_id_number_key UNIQUE (student_id_number);
ALTER TABLE students ADD CONSTRAINT IF NOT EXISTS students_status_check CHECK (status IN ('Active', 'Inactive', 'Graduated', 'Suspended'));

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'fk_student_group' AND table_name = 'students'
    ) THEN
        ALTER TABLE students 
        ADD CONSTRAINT fk_student_group 
        FOREIGN KEY (group_id) REFERENCES student_groups(id) ON DELETE RESTRICT;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_students_group ON students(group_id);
CREATE INDEX IF NOT EXISTS idx_students_email ON students(email);
CREATE INDEX IF NOT EXISTS idx_students_student_id ON students(student_id_number);
CREATE INDEX IF NOT EXISTS idx_students_full_name ON students(full_name);
CREATE INDEX IF NOT EXISTS idx_students_status ON students(status);
