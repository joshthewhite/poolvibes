ALTER TABLE users ADD COLUMN is_demo BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE users ADD COLUMN demo_expires_at TIMESTAMPTZ;

CREATE INDEX idx_users_demo_expires ON users(is_demo, demo_expires_at);
