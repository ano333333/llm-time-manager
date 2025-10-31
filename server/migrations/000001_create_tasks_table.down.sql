-- トリガーの削除
DROP TRIGGER IF EXISTS update_tasks_updated_at;

-- インデックスの削除
DROP INDEX IF EXISTS idx_tasks_created_at;
DROP INDEX IF EXISTS idx_tasks_goal_id;
DROP INDEX IF EXISTS idx_tasks_due;
DROP INDEX IF EXISTS idx_tasks_status;

-- テーブルの削除
DROP TABLE IF EXISTS tasks;

