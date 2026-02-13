DROP INDEX IF EXISTS idx_chemicals_user_id;
DROP INDEX IF EXISTS idx_service_records_user_id;
DROP INDEX IF EXISTS idx_equipment_user_id;
DROP INDEX IF EXISTS idx_tasks_user_id;
DROP INDEX IF EXISTS idx_chemistry_logs_user_id;

DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;

-- Recreate tables without user_id (SQLite lacks DROP COLUMN before 3.35.0)
CREATE TABLE chemistry_logs_backup AS SELECT id, ph, free_chlorine, combined_chlorine, total_alkalinity, cya, calcium_hardness, temperature, notes, tested_at, created_at, updated_at FROM chemistry_logs;
DROP TABLE chemistry_logs;
ALTER TABLE chemistry_logs_backup RENAME TO chemistry_logs;
CREATE INDEX idx_chemistry_logs_tested_at ON chemistry_logs(tested_at);

CREATE TABLE tasks_backup AS SELECT id, name, description, recurrence_frequency, recurrence_interval, due_date, status, completed_at, created_at, updated_at FROM tasks;
DROP TABLE tasks;
ALTER TABLE tasks_backup RENAME TO tasks;
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_due_date ON tasks(due_date);

CREATE TABLE equipment_backup AS SELECT id, name, category, manufacturer, model, serial_number, install_date, warranty_expiry, created_at, updated_at FROM equipment;
DROP TABLE equipment;
ALTER TABLE equipment_backup RENAME TO equipment;

CREATE TABLE service_records_backup AS SELECT id, equipment_id, service_date, description, cost, technician, created_at, updated_at FROM service_records;
DROP TABLE service_records;
ALTER TABLE service_records_backup RENAME TO service_records;
CREATE INDEX idx_service_records_equipment_id ON service_records(equipment_id);

CREATE TABLE chemicals_backup AS SELECT id, name, type, stock_amount, stock_unit, alert_threshold, last_purchased, created_at, updated_at FROM chemicals;
DROP TABLE chemicals;
ALTER TABLE chemicals_backup RENAME TO chemicals;
CREATE INDEX idx_chemicals_stock_amount ON chemicals(stock_amount);
