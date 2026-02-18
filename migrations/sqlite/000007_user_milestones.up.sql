CREATE TABLE IF NOT EXISTS user_milestones (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    milestone TEXT NOT NULL,
    earned_at TEXT NOT NULL,
    UNIQUE(user_id, milestone)
);
CREATE INDEX idx_user_milestones_user_id ON user_milestones(user_id);
