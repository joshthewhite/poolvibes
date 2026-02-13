ALTER TABLE users ADD COLUMN phone TEXT NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN notify_email INTEGER NOT NULL DEFAULT 1;
ALTER TABLE users ADD COLUMN notify_sms INTEGER NOT NULL DEFAULT 0;

CREATE TABLE IF NOT EXISTS task_notifications (
    id TEXT PRIMARY KEY,
    task_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    type TEXT NOT NULL,
    due_date TEXT NOT NULL,
    sent_at TEXT NOT NULL,
    UNIQUE(task_id, type, due_date)
);

CREATE INDEX idx_task_notifications_task_id ON task_notifications(task_id);
CREATE INDEX idx_task_notifications_user_id ON task_notifications(user_id);
