INSERT INTO subjects (subject_name, subject_code, credits, faculty_id)
SELECT DISTINCT 
    subject_name,
    UPPER(SUBSTRING(subject_name FROM 1 FOR 2)) || LPAD(FLOOR(RANDOM() * 1000)::TEXT, 3, '0'),
    3,
    1
FROM schedule
WHERE subject_name IS NOT NULL
ON CONFLICT (subject_code) DO NOTHING;

ALTER TABLE schedule ADD COLUMN IF NOT EXISTS subject_id INTEGER;
ALTER TABLE schedule ADD COLUMN IF NOT EXISTS teacher_id INTEGER;
ALTER TABLE schedule ADD COLUMN IF NOT EXISTS day_of_week VARCHAR(20);
ALTER TABLE schedule ADD COLUMN IF NOT EXISTS room_number VARCHAR(50);
ALTER TABLE schedule ADD COLUMN IF NOT EXISTS semester VARCHAR(20) DEFAULT 'Spring 2026';
ALTER TABLE schedule ADD COLUMN IF NOT EXISTS academic_year VARCHAR(20) DEFAULT '2025-2026';
ALTER TABLE schedule ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE schedule ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

UPDATE schedule s
SET subject_id = sub.id
FROM subjects sub
WHERE s.subject_name = sub.subject_name
AND s.subject_id IS NULL;

INSERT INTO teachers (full_name, email, phone, faculty_id, position, degree)
VALUES ('Default Instructor', 'default@nu.edu.kz', '+77000000000', 1, 'Lecturer', 'MSc')
ON CONFLICT (email) DO NOTHING;

UPDATE schedule
SET teacher_id = (SELECT id FROM teachers WHERE email = 'default@nu.edu.kz' LIMIT 1)
WHERE teacher_id IS NULL;

UPDATE schedule
SET day_of_week = 'Monday'
WHERE day_of_week IS NULL;

ALTER TABLE schedule ALTER COLUMN subject_id SET NOT NULL;
ALTER TABLE schedule ALTER COLUMN teacher_id SET NOT NULL;
ALTER TABLE schedule ALTER COLUMN day_of_week SET NOT NULL;

ALTER TABLE schedule DROP CONSTRAINT IF EXISTS fk_schedule_subject;
ALTER TABLE schedule DROP CONSTRAINT IF EXISTS fk_schedule_teacher;
ALTER TABLE schedule DROP CONSTRAINT IF EXISTS fk_schedule_group;

ALTER TABLE schedule ADD CONSTRAINT fk_schedule_subject 
    FOREIGN KEY (subject_id) REFERENCES subjects(id) ON DELETE CASCADE;
ALTER TABLE schedule ADD CONSTRAINT fk_schedule_teacher 
    FOREIGN KEY (teacher_id) REFERENCES teachers(id) ON DELETE RESTRICT;
ALTER TABLE schedule ADD CONSTRAINT fk_schedule_group 
    FOREIGN KEY (group_id) REFERENCES student_groups(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_schedule_subject ON schedule(subject_id);
CREATE INDEX IF NOT EXISTS idx_schedule_teacher ON schedule(teacher_id);
CREATE INDEX IF NOT EXISTS idx_schedule_day ON schedule(day_of_week);
CREATE INDEX IF NOT EXISTS idx_schedule_semester ON schedule(semester);
CREATE INDEX IF NOT EXISTS idx_schedule_group_day ON schedule(group_id, day_of_week);
