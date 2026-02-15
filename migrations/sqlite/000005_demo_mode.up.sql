ALTER TABLE users ADD COLUMN is_demo INTEGER NOT NULL DEFAULT 0;
ALTER TABLE users ADD COLUMN demo_expires_at TEXT;

CREATE INDEX idx_users_demo_expires ON users(is_demo, demo_expires_at);
