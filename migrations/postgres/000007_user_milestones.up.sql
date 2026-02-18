CREATE TABLE IF NOT EXISTS user_milestones (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    milestone TEXT NOT NULL,
    earned_at TIMESTAMPTZ NOT NULL,
    UNIQUE(user_id, milestone)
);
CREATE INDEX idx_user_milestones_user_id ON user_milestones(user_id);
