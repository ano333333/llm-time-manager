-- +goose Up
-- capture_schedulesテーブルの作成
CREATE TABLE IF NOT EXISTS capture_schedules (
    id TEXT PRIMARY KEY,
    active INTEGER NOT NULL DEFAULT 0,
    interval_min INTEGER NOT NULL DEFAULT 10,
    retention_max_items INTEGER NOT NULL DEFAULT 100,
    retention_max_days INTEGER NOT NULL DEFAULT 30,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- updated_atの自動更新トリガー
-- +goose StatementBegin
CREATE TRIGGER IF NOT EXISTS update_capture_schedules_updated_at
    AFTER UPDATE ON capture_schedules
    FOR EACH ROW
    WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE capture_schedules SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;
-- +goose StatementEnd

-- +goose Down
-- トリガーの削除
DROP TRIGGER IF EXISTS update_capture_schedules_updated_at;

-- テーブルの削除
DROP TABLE IF EXISTS capture_schedules;

