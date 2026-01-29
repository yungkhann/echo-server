CREATE TABLE IF NOT EXISTS schedule (
    id SERIAL PRIMARY KEY,
    subject_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    teacher_id INTEGER NOT NULL,
    day_of_week VARCHAR(20) NOT NULL CHECK (day_of_week IN ('Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday')),
    time_slot VARCHAR(50) NOT NULL,
    room_number VARCHAR(50),
    semester VARCHAR(20) NOT NULL DEFAULT 'Fall 2025',
    academic_year VARCHAR(20) NOT NULL DEFAULT '2025-2026',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_schedule_subject FOREIGN KEY (subject_id) REFERENCES subjects(id) ON DELETE CASCADE,
    CONSTRAINT fk_schedule_group FOREIGN KEY (group_id) REFERENCES student_groups(id) ON DELETE CASCADE,
    CONSTRAINT fk_schedule_teacher FOREIGN KEY (teacher_id) REFERENCES teachers(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_schedule_subject ON schedule(subject_id);
CREATE INDEX IF NOT EXISTS idx_schedule_group ON schedule(group_id);
CREATE INDEX IF NOT EXISTS idx_schedule_teacher ON schedule(teacher_id);
CREATE INDEX IF NOT EXISTS idx_schedule_day ON schedule(day_of_week);
CREATE INDEX IF NOT EXISTS idx_schedule_semester ON schedule(semester);
CREATE INDEX IF NOT EXISTS idx_schedule_group_day ON schedule(group_id, day_of_week);

ALTER TABLE schedule ADD COLUMN IF NOT EXISTS subject_name VARCHAR(255);

INSERT INTO schedule (subject_id, group_id, teacher_id, day_of_week, time_slot, room_number, semester) VALUES
(1, 1, 1, 'Monday', '09:00-10:30', 'Room 301', 'Spring 2026'),
(2, 1, 2, 'Tuesday', '10:45-12:15', 'Room 302', 'Spring 2026'),
(5, 3, 4, 'Wednesday', '15:00-16:30', 'Gym Hall', 'Spring 2026'),
(3, 3, 3, 'Thursday', '09:00-10:30', 'Lecture Hall 101', 'Spring 2026')
ON CONFLICT DO NOTHING;
