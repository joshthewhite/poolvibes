ALTER TABLE users ADD COLUMN phone TEXT NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN notify_email BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE users ADD COLUMN notify_sms BOOLEAN NOT NULL DEFAULT FALSE;

CREATE TABLE IF NOT EXISTS task_notifications (
    id UUID PRIMARY KEY,
    task_id UUID NOT NULL,
    user_id UUID NOT NULL,
    type TEXT NOT NULL,
    due_date DATE NOT NULL,
    sent_at TIMESTAMPTZ NOT NULL,
    UNIQUE(task_id, type, due_date)
);

CREATE INDEX idx_task_notifications_task_id ON task_notifications(task_id);
CREATE INDEX idx_task_notifications_user_id ON task_notifications(user_id);
