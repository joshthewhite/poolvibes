CREATE TABLE IF NOT EXISTS chemistry_logs (
    id UUID PRIMARY KEY,
    ph DOUBLE PRECISION NOT NULL,
    free_chlorine DOUBLE PRECISION NOT NULL,
    combined_chlorine DOUBLE PRECISION NOT NULL,
    total_alkalinity DOUBLE PRECISION NOT NULL,
    cya DOUBLE PRECISION NOT NULL,
    calcium_hardness DOUBLE PRECISION NOT NULL,
    temperature DOUBLE PRECISION NOT NULL DEFAULT 0,
    notes TEXT NOT NULL DEFAULT '',
    tested_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_chemistry_logs_tested_at ON chemistry_logs(tested_at);

CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    recurrence_frequency TEXT NOT NULL DEFAULT 'weekly',
    recurrence_interval INTEGER NOT NULL DEFAULT 1,
    due_date TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_due_date ON tasks(due_date);

CREATE TABLE IF NOT EXISTS equipment (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    category TEXT NOT NULL DEFAULT 'other',
    manufacturer TEXT NOT NULL DEFAULT '',
    model TEXT NOT NULL DEFAULT '',
    serial_number TEXT NOT NULL DEFAULT '',
    install_date TIMESTAMPTZ,
    warranty_expiry TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS service_records (
    id UUID PRIMARY KEY,
    equipment_id UUID NOT NULL REFERENCES equipment(id) ON DELETE CASCADE,
    service_date TIMESTAMPTZ NOT NULL,
    description TEXT NOT NULL,
    cost DOUBLE PRECISION NOT NULL DEFAULT 0,
    technician TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_service_records_equipment_id ON service_records(equipment_id);

CREATE TABLE IF NOT EXISTS chemicals (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL DEFAULT 'other',
    stock_amount DOUBLE PRECISION NOT NULL DEFAULT 0,
    stock_unit TEXT NOT NULL DEFAULT 'lbs',
    alert_threshold DOUBLE PRECISION NOT NULL DEFAULT 0,
    last_purchased TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_chemicals_stock_amount ON chemicals(stock_amount);
