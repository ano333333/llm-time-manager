-- +goose Up
-- chat_messagesテーブルの作成
CREATE TABLE IF NOT EXISTS chat_messages (
    id TEXT PRIMARY KEY,
    role TEXT NOT NULL CHECK (role IN ('user', 'assistant', 'system')),
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- インデックスの作成
CREATE INDEX IF NOT EXISTS idx_chat_messages_created_at ON chat_messages(created_at);
CREATE INDEX IF NOT EXISTS idx_chat_messages_role ON chat_messages(role);

-- +goose Down
-- インデックスの削除
DROP INDEX IF EXISTS idx_chat_messages_role;
DROP INDEX IF EXISTS idx_chat_messages_created_at;

-- テーブルの削除
DROP TABLE IF EXISTS chat_messages;

