-- Add student_id reference to users table to link users to students
ALTER TABLE users ADD COLUMN IF NOT EXISTS student_id INTEGER;
ALTER TABLE users ADD CONSTRAINT fk_users_student_id FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_users_student_id ON users(student_id);

-- Note: Admin should manually link existing student users to their student records
-- Example: UPDATE users SET student_id = <student_id> WHERE id = <user_id> AND role = 'student';
