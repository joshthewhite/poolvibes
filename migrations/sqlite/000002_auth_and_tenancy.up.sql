CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    is_admin INTEGER NOT NULL DEFAULT 0,
    is_disabled INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE UNIQUE INDEX idx_users_email ON users(email);

CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TEXT NOT NULL,
    created_at TEXT NOT NULL
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

ALTER TABLE chemistry_logs ADD COLUMN user_id TEXT NOT NULL DEFAULT '';
ALTER TABLE tasks ADD COLUMN user_id TEXT NOT NULL DEFAULT '';
ALTER TABLE equipment ADD COLUMN user_id TEXT NOT NULL DEFAULT '';
ALTER TABLE service_records ADD COLUMN user_id TEXT NOT NULL DEFAULT '';
ALTER TABLE chemicals ADD COLUMN user_id TEXT NOT NULL DEFAULT '';

CREATE INDEX idx_chemistry_logs_user_id ON chemistry_logs(user_id);
CREATE INDEX idx_tasks_user_id ON tasks(user_id);
CREATE INDEX idx_equipment_user_id ON equipment(user_id);
CREATE INDEX idx_service_records_user_id ON service_records(user_id);
CREATE INDEX idx_chemicals_user_id ON chemicals(user_id);
