ALTER TABLE task_notifications DROP CONSTRAINT IF EXISTS task_notifications_user_type_date_uq;

DELETE FROM task_notifications WHERE task_id IS NULL;

ALTER TABLE task_notifications ALTER COLUMN task_id SET NOT NULL;
ALTER TABLE task_notifications ALTER COLUMN task_id DROP DEFAULT;

ALTER TABLE task_notifications ADD CONSTRAINT task_notifications_task_id_type_due_date_key UNIQUE (task_id, type, due_date);

CREATE INDEX idx_task_notifications_task_id ON task_notifications(task_id);
