CREATE TABLE IF NOT EXISTS faculties (
    id SERIAL PRIMARY KEY,
    faculty_name VARCHAR(255) NOT NULL,
    dean_name VARCHAR(255),
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_faculties_name ON faculties(faculty_name);

INSERT INTO faculties (faculty_name, dean_name, description) VALUES
('Computer Science and IT', 'Dr. Nurlan Abdullayev', 'Faculty of Computer Science and Information Technology'),
('Engineering', 'Dr. Asel Sarsenova', 'Faculty of Engineering and Applied Sciences'),
('Business and Economics', 'Dr. Erlan Zhumabayev', 'Faculty of Business and Economics'),
('Mathematics', 'Dr. Gulnara Temirova', 'Faculty of Mathematics and Statistics')
ON CONFLICT DO NOTHING;
