-- +goose Up
-- chat_messagesテーブルの作成
CREATE TABLE IF NOT EXISTS chat_messages (
    id TEXT PRIMARY KEY,
    role TEXT NOT NULL CHECK (role IN ('user', 'assistant', 'system')),
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- インデックスの作成
CREATE INDEX IF NOT EXISTS idx_chat_messages_created_at ON chat_messages(created_at);
CREATE INDEX IF NOT EXISTS idx_chat_messages_role ON chat_messages(role);

-- updated_atの自動更新トリガー
-- +goose StatementBegin
CREATE TRIGGER IF NOT EXISTS update_chat_messages_updated_at
    AFTER UPDATE ON chat_messages
    FOR EACH ROW
BEGIN
    UPDATE chat_messages SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
-- +goose StatementEnd

-- +goose Down
-- インデックスの削除
DROP INDEX IF EXISTS idx_chat_messages_role;
DROP INDEX IF EXISTS idx_chat_messages_created_at;

-- テーブルの削除
DROP TABLE IF EXISTS chat_messages;

