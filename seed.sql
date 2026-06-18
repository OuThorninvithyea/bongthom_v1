-- Seed data — runs AFTER migrations
-- Admin user with password: admin123

INSERT INTO tbl_users (user_name, password, role_id, role_name, is_admin, status_id)
SELECT 'ADMIN', '$2a$10$jgwUpnQQs6xVwoxIIswIXubE7uey4QdXmPtcHuKsS99cYGlIuQJqu', 0, 'SuperAdmin', true, 1
WHERE NOT EXISTS (SELECT 1 FROM tbl_users WHERE user_name = 'ADMIN');
