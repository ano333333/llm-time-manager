-- +goose Up
-- +goose StatementBegin
ALTER TABLE capture_schedules DROP COLUMN retention_max_items;
ALTER TABLE capture_schedules DROP COLUMN retention_max_days;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE capture_schedules ADD COLUMN retention_max_items INTEGER NOT NULL DEFAULT 100;
ALTER TABLE capture_schedules ADD COLUMN retention_max_days INTEGER NOT NULL DEFAULT 30;
-- +goose StatementEnd
