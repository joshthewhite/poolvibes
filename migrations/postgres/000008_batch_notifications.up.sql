-- Change notification dedup from per-task to per-user-per-day.
-- Users now receive at most one notification per channel per day.
ALTER TABLE task_notifications DROP CONSTRAINT IF EXISTS task_notifications_task_id_type_due_date_key;
DROP INDEX IF EXISTS task_notifications_task_id_type_due_date_key;

ALTER TABLE task_notifications ALTER COLUMN task_id DROP NOT NULL;
ALTER TABLE task_notifications ALTER COLUMN task_id SET DEFAULT NULL;

ALTER TABLE task_notifications ADD CONSTRAINT task_notifications_user_type_date_uq UNIQUE (user_id, type, due_date);

DROP INDEX IF EXISTS idx_task_notifications_task_id;
