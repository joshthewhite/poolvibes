-- Change notification dedup from per-task to per-user-per-day.
-- Users now receive at most one notification per channel per day.
CREATE TABLE task_notifications_new (
    id TEXT PRIMARY KEY,
    task_id TEXT NOT NULL DEFAULT '',
    user_id TEXT NOT NULL,
    type TEXT NOT NULL,
    due_date TEXT NOT NULL,
    sent_at TEXT NOT NULL,
    UNIQUE(user_id, type, due_date)
);

INSERT OR IGNORE INTO task_notifications_new (id, task_id, user_id, type, due_date, sent_at)
    SELECT id, task_id, user_id, type, due_date, sent_at FROM task_notifications;

DROP TABLE task_notifications;
ALTER TABLE task_notifications_new RENAME TO task_notifications;

CREATE INDEX idx_task_notifications_user_id ON task_notifications(user_id);
