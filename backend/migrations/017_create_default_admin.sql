-- Create default admin user if not exists
-- Email: admin@gmail.com
-- Password: admin
INSERT INTO users (email, password, role, full_name, created_at, is_active)
VALUES (
    'admin@gmail.com',
    '$2a$10$Jp6UyXVLnLUxsbDut/psGu/YokeOHojSl8RTJUl0uHd9O8XkXD8DK',
    'admin',
    'System Administrator',
    CURRENT_TIMESTAMP,
    true
)
ON CONFLICT (email) DO UPDATE SET
    role = 'admin',
    full_name = 'System Administrator',
    is_active = true;
