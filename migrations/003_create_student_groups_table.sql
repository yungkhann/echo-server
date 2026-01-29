CREATE TABLE IF NOT EXISTS student_groups (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(100) NOT NULL UNIQUE,
    faculty_id INTEGER NOT NULL,
    course_year INTEGER NOT NULL DEFAULT 1,
    academic_year VARCHAR(20) NOT NULL DEFAULT '2025-2026',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_faculty FOREIGN KEY (faculty_id) REFERENCES faculties(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_student_groups_faculty ON student_groups(faculty_id);
CREATE INDEX IF NOT EXISTS idx_student_groups_name ON student_groups(group_name);
CREATE INDEX IF NOT EXISTS idx_student_groups_year ON student_groups(course_year);

INSERT INTO student_groups (group_name, faculty_id, course_year, academic_year) VALUES
('CSCI-231', 1, 2, '2025-2026'),
('CSCI-235', 1, 2, '2025-2026'),
('MATH-231', 4, 2, '2025-2026'),
('MATH-162', 4, 1, '2025-2026'),
('PHYS-162', 2, 1, '2025-2026')
ON CONFLICT (group_name) DO NOTHING;
