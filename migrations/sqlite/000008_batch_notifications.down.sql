CREATE TABLE task_notifications_old (
    id TEXT PRIMARY KEY,
    task_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    type TEXT NOT NULL,
    due_date TEXT NOT NULL,
    sent_at TEXT NOT NULL,
    UNIQUE(task_id, type, due_date)
);

INSERT OR IGNORE INTO task_notifications_old (id, task_id, user_id, type, due_date, sent_at)
    SELECT id, task_id, user_id, type, due_date, sent_at FROM task_notifications;

DROP TABLE task_notifications;
ALTER TABLE task_notifications_old RENAME TO task_notifications;

CREATE INDEX idx_task_notifications_task_id ON task_notifications(task_id);
CREATE INDEX idx_task_notifications_user_id ON task_notifications(user_id);
