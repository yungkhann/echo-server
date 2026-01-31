CREATE TABLE IF NOT EXISTS subjects (
    id SERIAL PRIMARY KEY,
    subject_name VARCHAR(255) NOT NULL,
    subject_code VARCHAR(50) UNIQUE NOT NULL,
    credits INTEGER NOT NULL DEFAULT 3,
    description TEXT,
    faculty_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_subject_faculty FOREIGN KEY (faculty_id) REFERENCES faculties(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_subjects_name ON subjects(subject_name);
CREATE INDEX IF NOT EXISTS idx_subjects_code ON subjects(subject_code);
CREATE INDEX IF NOT EXISTS idx_subjects_faculty ON subjects(faculty_id);

INSERT INTO subjects (subject_name, subject_code, credits, description, faculty_id) VALUES
('Data Structures and Algorithms', 'CSCI231', 5, 'Advanced algorithms and data structures', 1),
('Database Systems', 'CSCI235', 5, 'Relational and NoSQL databases', 1),
('Calculus II', 'MATH231', 5, 'Differential and integral calculus', 4),
('Calculus I', 'MATH162', 5, 'Introduction to calculus', 4),
('Physics II', 'PHYS162', 5, 'Classical mechanics and thermodynamics', 2)
ON CONFLICT (subject_code) DO NOTHING;
