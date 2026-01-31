DO $$ 
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'faculties' AND column_name = 'name'
    ) THEN
        ALTER TABLE faculties RENAME COLUMN name TO faculty_name;
    END IF;
END $$;

ALTER TABLE faculties ADD COLUMN IF NOT EXISTS dean_name VARCHAR(255);
ALTER TABLE faculties ADD COLUMN IF NOT EXISTS description TEXT;
ALTER TABLE faculties ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE faculties ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

UPDATE faculties 
SET created_at = CURRENT_TIMESTAMP 
WHERE created_at IS NULL;

UPDATE faculties 
SET updated_at = CURRENT_TIMESTAMP 
WHERE updated_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_faculties_name ON faculties(faculty_name);
