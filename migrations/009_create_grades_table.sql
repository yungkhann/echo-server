CREATE TABLE IF NOT EXISTS grades (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL,
    subject_id INTEGER NOT NULL,
    grade_type VARCHAR(50) NOT NULL CHECK (grade_type IN ('Midterm', 'Final', 'Assignment', 'Quiz', 'Project', 'Lab')),
    grade_value DECIMAL(5,2) NOT NULL CHECK (grade_value >= 0 AND grade_value <= 100),
    max_grade DECIMAL(5,2) NOT NULL DEFAULT 100,
    weight DECIMAL(5,2) DEFAULT 100,
    exam_date DATE,
    remarks TEXT,
    graded_by INTEGER,
    semester VARCHAR(20) NOT NULL DEFAULT 'Spring 2026',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_grades_student FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE,
    CONSTRAINT fk_grades_subject FOREIGN KEY (subject_id) REFERENCES subjects(id) ON DELETE CASCADE,
    CONSTRAINT fk_grades_grader FOREIGN KEY (graded_by) REFERENCES teachers(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_grades_student ON grades(student_id);
CREATE INDEX IF NOT EXISTS idx_grades_subject ON grades(subject_id);
CREATE INDEX IF NOT EXISTS idx_grades_type ON grades(grade_type);
CREATE INDEX IF NOT EXISTS idx_grades_semester ON grades(semester);
CREATE INDEX IF NOT EXISTS idx_grades_student_subject ON grades(student_id, subject_id);

INSERT INTO grades (student_id, subject_id, grade_type, grade_value, exam_date, graded_by, semester) VALUES
(1, 1, 'Midterm', 85.5, '2026-01-10', 1, 'Spring 2026'),
(1, 1, 'Assignment', 90.0, '2026-01-12', 1, 'Spring 2026'),
(2, 1, 'Midterm', 78.0, '2026-01-10', 1, 'Spring 2026'),
(3, 3, 'Quiz', 92.5, '2026-01-15', 2, 'Spring 2026')
ON CONFLICT DO NOTHING;
