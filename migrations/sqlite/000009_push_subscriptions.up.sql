CREATE TABLE IF NOT EXISTS push_subscriptions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    endpoint TEXT NOT NULL,
    p256dh TEXT NOT NULL,
    auth TEXT NOT NULL,
    created_at TEXT NOT NULL,
    UNIQUE(user_id, endpoint)
);

ALTER TABLE users ADD COLUMN notify_push INTEGER NOT NULL DEFAULT 0;
