DROP TABLE IF EXISTS task_notifications;

-- SQLite lacks DROP COLUMN before 3.35.0, recreate users table without notification columns
CREATE TABLE users_backup AS SELECT id, email, password_hash, is_admin, is_disabled, created_at, updated_at FROM users;
DROP TABLE users;
ALTER TABLE users_backup RENAME TO users;
CREATE UNIQUE INDEX idx_users_email ON users(email);
