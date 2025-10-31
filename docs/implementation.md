# 実装ガイド

## 実装メモ

### フロントエンド（React）

- **状態管理**: Redux / Pinia / Zustand 等を使用
- **ルーティング**: React Router（SPA）
- **SSE/WS 対応**: EventSource API（SSE）または WebSocket でストリーム応答を受信
- **スタイリング**: Tailwind CSS / CSS Modules 推奨
- **型安全**: TypeScript 必須
- **バリデーション**: Zod 等でスキーマバリデーション

```ts
// Zod スキーマ例
import { z } from 'zod';

export const TaskSchema = z.object({
  id: z.string().uuid(),
  goalId: z.string().uuid().optional(),
  title: z.string().min(1).max(200),
  description: z.string().max(2000),
  due: z.string().datetime(),
  estimateMin: z.number().int().min(0),
  priority: z.number().int().min(1).max(5),
  status: z.enum(['todo', 'doing', 'paused', 'done', 'archived']),
  tags: z.array(z.string()),
  attachments: z.array(z.string()),
  createdAt: z.string().datetime(),
  updatedAt: z.string().datetime(),
});

export type Task = z.infer<typeof TaskSchema>;
```

### バックエンド（Go）

- **ルータ**: chi または gin
- **データベース**: SQLite + sqlc / gorm
- **マイグレーション**: golang-migrate
- **WebSocket**: gorilla/websocket
- **構造化ログ**: zap / zerolog

```go
// ハンドラ例
package http

import (
    "net/http"
    "github.com/go-chi/chi/v5"
)

func (s *Server) TaskRoutes() chi.Router {
    r := chi.NewRouter()
    
    r.Get("/", s.handleGetTasks)
    r.Post("/", s.handleCreateTask)
    r.Get("/{id}", s.handleGetTask)
    r.Patch("/{id}", s.handleUpdateTask)
    r.Delete("/{id}", s.handleDeleteTask)
    
    return r
}

func (s *Server) handleGetTasks(w http.ResponseWriter, r *http.Request) {
    // クエリパラメータ取得
    status := r.URL.Query().Get("status")
    due := r.URL.Query().Get("due")
    
    // DB から取得
    tasks, err := s.store.GetTasks(r.Context(), status, due)
    if err != nil {
        s.respondError(w, err)
        return
    }
    
    s.respondJSON(w, http.StatusOK, map[string]interface{}{
        "tasks": tasks,
    })
}
```

### 画像保存

- OS 推奨ディレクトリを使用（ユーザー指定可）
- **macOS**: `~/Library/Application Support/LLMTimeManager/Screenshots/`
- **Windows**: `%APPDATA%\LLMTimeManager\Screenshots\`
- **Linux**: `~/.local/share/llm-time-manager/screenshots/`
- **iOS**: アプリ Documents ディレクトリ

### LLM 統合

- ストリーム応答対応（SSE）
- JSON モードでエンティティ抽出
- システムプロンプト例:

```
あなたは時間管理アシスタントです。
ユーザーの発言から、タスクや目標を抽出してください。

出力形式:
- 通常の応答: テキスト
- エンティティ抽出: JSON { "type": "task"|"goal", ... }

例:
ユーザー: 「来週水曜にレポートを提出したい」
→ {"type":"task","title":"レポート提出","due":"2025-11-05","estimateMin":120}
```

## バリデーション/エラー UX

### 入力バリデーション

- **期日 < 今日**: 警告を表示（保存は可能）
- **タイトル空欄**: エラー（保存不可）
- **見積時間が負数**: エラー
- **優先度が範囲外**: エラー

### エラー表示

```tsx
// Toast 通知例
import { toast } from 'react-hot-toast';

try {
  await createTask(data);
  toast.success('タスクを作成しました');
} catch (error) {
  if (error.code === 'VALIDATION_ERROR') {
    toast.error(error.message);
  } else {
    toast.error('エラーが発生しました');
  }
}
```

### キャプチャ失敗時

- リトライボタンを表示
- エラーログをクリップボードにコピー可能
- 権限エラーの場合は設定画面への誘導

```tsx
function CaptureErrorView({ error }: { error: CaptureError }) {
  const handleRetry = async () => {
    try {
      await bridge.captureScreenshot();
    } catch (e) {
      // エラー処理
    }
  };

  const copyLog = () => {
    navigator.clipboard.writeText(JSON.stringify(error, null, 2));
    toast.success('ログをコピーしました');
  };

  return (
    <div className="error-view">
      <p>キャプチャに失敗しました</p>
      <p>{error.message}</p>
      <button onClick={handleRetry}>リトライ</button>
      <button onClick={copyLog}>ログをコピー</button>
      {error.code === 'PERMISSION_DENIED' && (
        <Link to="/settings/local">設定を確認</Link>
      )}
    </div>
  );
}
```

### ブリッジ未初期化

```tsx
function BridgeErrorView() {
  return (
    <div className="error-view">
      <p>ネイティブブリッジが初期化されていません</p>
      <button onClick={() => window.location.reload()}>
        ページをリロード
      </button>
    </div>
  );
}
```

## 記録/監査

### 主要イベントのロギング

- タスク作成/更新/完了
- キャプチャ実行
- 権限状態変化
- エラー発生

```go
// ログ例
type AuditLog struct {
    ID        string    `json:"id"`
    EventType string    `json:"event_type"` // task_created, capture_executed, etc.
    EntityID  string    `json:"entity_id"`
    UserID    string    `json:"user_id"` // プロトタイプでは "local"
    Metadata  string    `json:"metadata"` // JSON
    CreatedAt time.Time `json:"created_at"`
}

func (s *Store) LogEvent(ctx context.Context, event AuditLog) error {
    _, err := s.db.ExecContext(ctx, `
        INSERT INTO audit_logs (id, event_type, entity_id, user_id, metadata, created_at)
        VALUES (?, ?, ?, ?, ?, ?)
    `, event.ID, event.EventType, event.EntityID, event.UserID, event.Metadata, event.CreatedAt)
    return err
}
```

### 保持期間

- ローカル 30 日（プロトタイプ）
- 古いログは自動削除（cron または起動時）

```go
func (s *Store) CleanupOldLogs(ctx context.Context, days int) error {
    _, err := s.db.ExecContext(ctx, `
        DELETE FROM audit_logs
        WHERE created_at < datetime('now', '-' || ? || ' days')
    `, days)
    return err
}
```

## メトリクス（プロトタイプ）

### 計測項目

- 1 日あたり作成タスク数
- 完了率（完了タスク数 / 全タスク数）
- キャプチャ回数
- チャットからのタスク化比率

```sql
-- 1日あたりのタスク作成数
SELECT DATE(created_at) as date, COUNT(*) as count
FROM tasks
GROUP BY DATE(created_at)
ORDER BY date DESC
LIMIT 30;

-- 完了率
SELECT
  COUNT(*) FILTER (WHERE status = 'done') * 100.0 / COUNT(*) as completion_rate
FROM tasks
WHERE created_at >= datetime('now', '-30 days');

-- チャットからのタスク化比率
SELECT
  COUNT(*) FILTER (WHERE source = 'chat') * 100.0 / COUNT(*) as chat_ratio
FROM tasks
WHERE created_at >= datetime('now', '-30 days');
```

### ダッシュボード表示

```tsx
function MetricsDashboard() {
  const [metrics, setMetrics] = useState<Metrics | null>(null);

  useEffect(() => {
    fetch('/metrics')
      .then(res => res.json())
      .then(setMetrics);
  }, []);

  if (!metrics) return <div>Loading...</div>;

  return (
    <div className="metrics-dashboard">
      <MetricCard title="今日のタスク" value={metrics.tasksToday} />
      <MetricCard title="完了率" value={`${metrics.completionRate}%`} />
      <MetricCard title="キャプチャ回数" value={metrics.captureCount} />
    </div>
  );
}
```

## テスト観点（抜粋）

### ユニットテスト

- **フロント**: React Testing Library + Vitest
- **バック**: Go の標準 `testing` パッケージ

```go
// テスト例
func TestCreateTask(t *testing.T) {
    store := NewTestStore(t)
    defer store.Close()

    task := &Task{
        ID:          "task-123",
        Title:       "テストタスク",
        Description: "説明",
        Due:         time.Now().Add(24 * time.Hour),
        Status:      "todo",
    }

    err := store.CreateTask(context.Background(), task)
    assert.NoError(t, err)

    // 取得して確認
    got, err := store.GetTask(context.Background(), task.ID)
    assert.NoError(t, err)
    assert.Equal(t, task.Title, got.Title)
}
```

### 統合テスト

- 権限未許可時のキャプチャ動作
- WebSocket 切断時のチャット再接続
- 定期キャプチャのインターバル精度とジッター動作
- 大量キャプチャ通知の一覧パフォーマンス

```ts
// E2E テスト例（Playwright）
import { test, expect } from '@playwright/test';

test('権限未許可時にエラーが表示される', async ({ page }) => {
  // ブリッジをモック
  await page.addInitScript(() => {
    window.bridge = {
      captureScreenshot: async () => {
        throw new Error('PERMISSION_DENIED');
      },
      requestPermission: async () => 'denied',
    };
  });

  await page.goto('/capture');
  await page.click('button:has-text("キャプチャ開始")');

  // エラーメッセージが表示されることを確認
  await expect(page.locator('text=権限が拒否されています')).toBeVisible();
});
```

### パフォーマンステスト

- 1000件のタスク一覧表示
- 100件のスクリーンショットサムネイル表示
- チャットのストリーム応答の遅延測定

## デバッグ

### デバッグログ

設定画面でトグル可能。

```ts
// ログユーティリティ
const logger = {
  debug: (message: string, data?: any) => {
    if (settings.debugLog) {
      console.log(`[DEBUG] ${message}`, data);
    }
  },
  error: (message: string, error?: any) => {
    console.error(`[ERROR] ${message}`, error);
  },
};

// 使用例
logger.debug('Fetching tasks', { status: 'todo' });
```

### 開発ツール

- React DevTools
- Go pprof（プロファイリング）
- SQLite DB Browser（DB確認）

## デプロイ/配布

### ビルド手順

```bash
# フロントエンド
cd web
npm run build  # dist/ に静的ファイル生成

# バックエンド
cd server
go build -o bin/llm-time-manager ./cmd/api

# クライアント（例: macOS）
cd clients/macos-swift
xcodebuild -scheme LLMTimeManager -configuration Release
```

### パッケージング

- **macOS**: DMG または PKG
- **Windows**: インストーラー（Inno Setup / WiX）
- **Linux**: Deb / RPM / AppImage
- **iOS**: TestFlight → App Store

### 配布物

- サーバーバイナリ（`server/bin/`）
- 静的ファイル（`web/dist/` をバイナリに同梱）
- クライアントアプリ（各OS向けバンドル）

