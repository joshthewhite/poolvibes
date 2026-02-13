DROP INDEX IF EXISTS idx_chemicals_user_id;
DROP INDEX IF EXISTS idx_service_records_user_id;
DROP INDEX IF EXISTS idx_equipment_user_id;
DROP INDEX IF EXISTS idx_tasks_user_id;
DROP INDEX IF EXISTS idx_chemistry_logs_user_id;

DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;

ALTER TABLE chemistry_logs DROP COLUMN user_id;
ALTER TABLE tasks DROP COLUMN user_id;
ALTER TABLE equipment DROP COLUMN user_id;
ALTER TABLE service_records DROP COLUMN user_id;
ALTER TABLE chemicals DROP COLUMN user_id;
