-- +goose Up
-- screenshotsテーブルの作成
CREATE TABLE IF NOT EXISTS screenshots (
    id TEXT PRIMARY KEY,
    path TEXT NOT NULL,
    thumb_path TEXT,
    captured_at DATETIME NOT NULL,
    mode TEXT NOT NULL DEFAULT 'manual' CHECK (mode IN ('manual', 'scheduled')),
    meta TEXT NOT NULL DEFAULT '{}',
    linked_task_id TEXT,
    linked_goal_id TEXT,
    FOREIGN KEY (linked_task_id) REFERENCES tasks(id) ON DELETE SET NULL,
    FOREIGN KEY (linked_goal_id) REFERENCES goals(id) ON DELETE SET NULL
);

-- インデックスの作成
CREATE INDEX IF NOT EXISTS idx_screenshots_captured_at ON screenshots(captured_at);
CREATE INDEX IF NOT EXISTS idx_screenshots_linked_task_id ON screenshots(linked_task_id);
CREATE INDEX IF NOT EXISTS idx_screenshots_mode ON screenshots(mode);

-- +goose Down
-- インデックスの削除
DROP INDEX IF EXISTS idx_screenshots_mode;
DROP INDEX IF EXISTS idx_screenshots_linked_task_id;
DROP INDEX IF EXISTS idx_screenshots_captured_at;

-- テーブルの削除
DROP TABLE IF EXISTS screenshots;

