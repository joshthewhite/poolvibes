CREATE TABLE IF NOT EXISTS chemistry_logs (
    id TEXT PRIMARY KEY,
    ph REAL NOT NULL,
    free_chlorine REAL NOT NULL,
    combined_chlorine REAL NOT NULL,
    total_alkalinity REAL NOT NULL,
    cya REAL NOT NULL,
    calcium_hardness REAL NOT NULL,
    temperature REAL NOT NULL DEFAULT 0,
    notes TEXT NOT NULL DEFAULT '',
    tested_at TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE INDEX idx_chemistry_logs_tested_at ON chemistry_logs(tested_at);

CREATE TABLE IF NOT EXISTS tasks (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    recurrence_frequency TEXT NOT NULL DEFAULT 'weekly',
    recurrence_interval INTEGER NOT NULL DEFAULT 1,
    due_date TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    completed_at TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_due_date ON tasks(due_date);

CREATE TABLE IF NOT EXISTS equipment (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    category TEXT NOT NULL DEFAULT 'other',
    manufacturer TEXT NOT NULL DEFAULT '',
    model TEXT NOT NULL DEFAULT '',
    serial_number TEXT NOT NULL DEFAULT '',
    install_date TEXT,
    warranty_expiry TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS service_records (
    id TEXT PRIMARY KEY,
    equipment_id TEXT NOT NULL REFERENCES equipment(id) ON DELETE CASCADE,
    service_date TEXT NOT NULL,
    description TEXT NOT NULL,
    cost REAL NOT NULL DEFAULT 0,
    technician TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE INDEX idx_service_records_equipment_id ON service_records(equipment_id);

CREATE TABLE IF NOT EXISTS chemicals (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL DEFAULT 'other',
    stock_amount REAL NOT NULL DEFAULT 0,
    stock_unit TEXT NOT NULL DEFAULT 'lbs',
    alert_threshold REAL NOT NULL DEFAULT 0,
    last_purchased TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE INDEX idx_chemicals_stock_amount ON chemicals(stock_amount);
