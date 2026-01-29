CREATE TABLE IF NOT EXISTS attendance (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL,
    subject_id INTEGER NOT NULL,
    schedule_id INTEGER,
    visit_day DATE NOT NULL,
    visited BOOLEAN NOT NULL DEFAULT false,
    remarks TEXT,
    marked_by INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_attendance_student FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
    CONSTRAINT fk_attendance_subject FOREIGN KEY (subject_id) REFERENCES subjects(id) ON DELETE CASCADE,
    CONSTRAINT fk_attendance_schedule FOREIGN KEY (schedule_id) REFERENCES schedule(id) ON DELETE SET NULL,
    CONSTRAINT fk_attendance_marker FOREIGN KEY (marked_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT unique_attendance UNIQUE (student_id, subject_id, visit_day)
);

CREATE INDEX IF NOT EXISTS idx_attendance_student ON attendance(student_id);
CREATE INDEX IF NOT EXISTS idx_attendance_subject ON attendance(subject_id);
CREATE INDEX IF NOT EXISTS idx_attendance_date ON attendance(visit_day);
CREATE INDEX IF NOT EXISTS idx_attendance_visited ON attendance(visited);
CREATE INDEX IF NOT EXISTS idx_attendance_schedule ON attendance(schedule_id);
CREATE INDEX IF NOT EXISTS idx_attendance_student_date ON attendance(student_id, visit_day);
CREATE INDEX IF NOT EXISTS idx_attendance_subject_date ON attendance(subject_id, visit_day);

INSERT INTO attendance (student_id, subject_id, visit_day, visited, remarks) VALUES
(1, 1, '2026-01-15', true, 'Present'),
(1, 1, '2026-01-16', true, 'Present'),
(2, 1, '2026-01-15', false, 'Absent'),
(2, 2, '2026-01-16', true, 'Present'),
(3, 3, '2026-01-17', true, 'Present')
ON CONFLICT (student_id, subject_id, visit_day) DO NOTHING;
