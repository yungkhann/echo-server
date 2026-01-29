CREATE TABLE IF NOT EXISTS teachers (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(50),
    faculty_id INTEGER NOT NULL,
    position VARCHAR(100) DEFAULT 'Lecturer',
    degree VARCHAR(100),
    specialization TEXT,
    hire_date DATE NOT NULL DEFAULT CURRENT_DATE,
    status VARCHAR(50) DEFAULT 'Active' CHECK (status IN ('Active', 'Inactive', 'On Leave')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_teacher_faculty FOREIGN KEY (faculty_id) REFERENCES faculties(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_teachers_email ON teachers(email);
CREATE INDEX IF NOT EXISTS idx_teachers_faculty ON teachers(faculty_id);
CREATE INDEX IF NOT EXISTS idx_teachers_status ON teachers(status);

INSERT INTO teachers (full_name, email, phone, faculty_id, position, degree, specialization) VALUES
('Aigul Karimova', 'aigul.karimova@nu.edu.kz', '+77012345680', 1, 'Professor', 'PhD', 'Computer Science'),
('Bakytzhan Omarov', 'bakytzhan.omarov@nu.edu.kz', '+77012345681', 1, 'Associate Professor', 'PhD', 'Software Engineering'),
('Saule Nurmaganbetova', 'saule.n@nu.edu.kz', '+77012345682', 4, 'Professor', 'PhD', 'Mathematics'),
('Marat Aliyev', 'marat.aliyev@nu.edu.kz', '+77012345683', 2, 'Lecturer', 'MSc', 'Physics')
ON CONFLICT (email) DO NOTHING;
