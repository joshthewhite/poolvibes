DROP INDEX IF EXISTS idx_users_demo_expires;

ALTER TABLE users DROP COLUMN demo_expires_at;
ALTER TABLE users DROP COLUMN is_demo;
