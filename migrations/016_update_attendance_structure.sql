ALTER TABLE attendance ADD COLUMN IF NOT EXISTS schedule_id INTEGER;
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS remarks TEXT;
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS marked_by INTEGER;
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE attendance ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

UPDATE attendance 
SET created_at = CURRENT_TIMESTAMP 
WHERE created_at IS NULL;

UPDATE attendance 
SET updated_at = CURRENT_TIMESTAMP 
WHERE updated_at IS NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'fk_attendance_student' AND table_name = 'attendance'
    ) THEN
        ALTER TABLE attendance 
        ADD CONSTRAINT fk_attendance_student 
        FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE;
    END IF;
    
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'fk_attendance_subject' AND table_name = 'attendance'
    ) THEN
        ALTER TABLE attendance 
        ADD CONSTRAINT fk_attendance_subject 
        FOREIGN KEY (subject_id) REFERENCES subjects(id) ON DELETE CASCADE;
    END IF;
    
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'fk_attendance_schedule' AND table_name = 'attendance'
    ) THEN
        ALTER TABLE attendance 
        ADD CONSTRAINT fk_attendance_schedule 
        FOREIGN KEY (schedule_id) REFERENCES schedule(id) ON DELETE SET NULL;
    END IF;
    
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'fk_attendance_marker' AND table_name = 'attendance'
    ) THEN
        ALTER TABLE attendance 
        ADD CONSTRAINT fk_attendance_marker 
        FOREIGN KEY (marked_by) REFERENCES users(id) ON DELETE SET NULL;
    END IF;
    
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'unique_attendance' AND table_name = 'attendance'
    ) THEN
        ALTER TABLE attendance 
        ADD CONSTRAINT unique_attendance 
        UNIQUE (student_id, subject_id, visit_day);
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_attendance_student ON attendance(student_id);
CREATE INDEX IF NOT EXISTS idx_attendance_subject ON attendance(subject_id);
CREATE INDEX IF NOT EXISTS idx_attendance_date ON attendance(visit_day);
CREATE INDEX IF NOT EXISTS idx_attendance_visited ON attendance(visited);
CREATE INDEX IF NOT EXISTS idx_attendance_schedule ON attendance(schedule_id);
CREATE INDEX IF NOT EXISTS idx_attendance_student_date ON attendance(student_id, visit_day);
CREATE INDEX IF NOT EXISTS idx_attendance_subject_date ON attendance(subject_id, visit_day);
