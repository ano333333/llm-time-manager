-- tasksテーブルの作成
CREATE TABLE IF NOT EXISTS tasks (
    id TEXT PRIMARY KEY,
    goal_id TEXT,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    due DATE,
    estimate_min INTEGER NOT NULL DEFAULT 0,
    priority INTEGER NOT NULL DEFAULT 3 CHECK (priority >= 1 AND priority <= 5),
    status TEXT NOT NULL DEFAULT 'todo' CHECK (status IN ('todo', 'doing', 'paused', 'done', 'archived')),
    tags TEXT NOT NULL DEFAULT '[]',
    attachments TEXT NOT NULL DEFAULT '[]',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (goal_id) REFERENCES goals(id) ON DELETE SET NULL
);

-- インデックスの作成
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_due ON tasks(due);
CREATE INDEX IF NOT EXISTS idx_tasks_goal_id ON tasks(goal_id);
CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);

-- updated_atの自動更新トリガー
CREATE TRIGGER IF NOT EXISTS update_tasks_updated_at
    AFTER UPDATE ON tasks
    FOR EACH ROW
BEGIN
    UPDATE tasks SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

