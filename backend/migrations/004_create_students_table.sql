CREATE TABLE IF NOT EXISTS students (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(50),
    gender VARCHAR(20) CHECK (gender IN ('Male', 'Female', 'Other')),
    birth_date DATE NOT NULL,
    group_id INTEGER NOT NULL,
    student_id_number VARCHAR(50) UNIQUE,
    enrollment_date DATE NOT NULL DEFAULT CURRENT_DATE,
    status VARCHAR(50) DEFAULT 'Active' CHECK (status IN ('Active', 'Inactive', 'Graduated', 'Suspended')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_student_group FOREIGN KEY (group_id) REFERENCES student_groups(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_students_group ON students(group_id);
CREATE INDEX IF NOT EXISTS idx_students_email ON students(email);
CREATE INDEX IF NOT EXISTS idx_students_student_id ON students(student_id_number);
CREATE INDEX IF NOT EXISTS idx_students_full_name ON students(full_name);
CREATE INDEX IF NOT EXISTS idx_students_status ON students(status);

INSERT INTO students (full_name, email, phone, gender, birth_date, group_id, student_id_number, status) VALUES
('Yernar Duzelbay', 'yernar.d@nu.edu.kz', '+77012345671', 'Male', '2002-12-05', 1, 'STU2025001', 'Active'),
('Aidana Ispayeva', 'aidana.i@nu.edu.kz', '+77012345672', 'Female', '2004-06-30', 1, 'STU2025002', 'Active'),
('Ualikhan Yertayev', 'ualikhan.y@nu.edu.kz', '+77012345673', 'Male', '2003-05-25', 2, 'STU2025003', 'Active')
ON CONFLICT (email) DO NOTHING;
