-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS goals (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    start_date DATE,
    end_date DATE,
    kpi_name TEXT,
    kpi_target FLOAT,
    kpi_unit TEXT,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'paused', 'done')),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_goals_status ON goals(status);

CREATE TRIGGER IF NOT EXISTS update_goals_updated_at
    BEFORE UPDATE ON goals
    FOR EACH ROW
    WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE goals SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_goals_updated_at;

DROP INDEX IF EXISTS idx_goals_status;

DROP TABLE IF EXISTS goals;

-- +goose StatementEnd
