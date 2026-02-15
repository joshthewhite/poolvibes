CREATE INDEX IF NOT EXISTS idx_chemistry_logs_user_tested_at ON chemistry_logs (user_id, tested_at DESC);
