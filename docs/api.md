# API 仕様

ローカルサーバ（Go）が提供する REST + WebSocket API。

ベース URL: `http://localhost:<port>`

## LLM

### POST /llm/chat

LLM とのチャット（Server-Sent Events または WebSocket でストリーム応答）。

**リクエスト**

```json
{
  "messages": [{ "role": "user", "content": "来週水曜にレポートを提出したい" }]
}
```

**レスポンス（SSE）**

```
data: {"type":"text","content":"わかりました"}
data: {"type":"text","content":"。"}
data: {"type":"entity","entity":{"type":"task","title":"レポート提出","due":"2025-11-05"}}
data: [DONE]
```

**エラー**

- `500 Internal Server Error` - LLM エンジンエラー

## タスク

### GET /tasks

タスク一覧取得。

**クエリパラメータ**

- `status` (optional): フィルタ（`todo|doing|paused|done`）
- `due` (optional): 期日フィルタ（`today|week|overdue`）
- `goalId` (optional): 目標 ID でフィルタ

**レスポンス**

```json
{
  "tasks": [
    {
      "id": "task-123",
      "goalId": "goal-456",
      "title": "レポート提出",
      "description": "...",
      "due": "2025-11-05",
      "estimateMin": 120,
      "priority": 3,
      "status": "todo",
      "tags": ["重要"],
      "attachments": [],
      "createdAt": "2025-10-29T10:00:00Z",
      "updatedAt": "2025-10-29T10:00:00Z"
    }
  ]
}
```

### POST /tasks

タスク作成。

**リクエスト**

```json
{
  "goalId": "goal-456",
  "title": "レポート提出",
  "description": "...",
  "due": "2025-11-05",
  "estimateMin": 120,
  "priority": 3,
  "tags": ["重要"]
}
```

**レスポンス**

```json
{
  "task": {
    "id": "task-123",
    ...
  }
}
```

### GET /tasks/:id

タスク詳細取得。

**レスポンス**

```json
{
  "task": {
    "id": "task-123",
    ...
  }
}
```

**エラー**

- `404 Not Found` - タスクが存在しない

### PATCH /tasks/:id

タスク更新。

**リクエスト**

```json
{
  "status": "doing",
  "priority": 4
}
```

**レスポンス**

```json
{
  "task": {
    "id": "task-123",
    ...
  }
}
```

### DELETE /tasks/:id

タスク削除。

**レスポンス**

```json
{
  "message": "deleted"
}
```

## 目標

### GET /goals

目標一覧取得。

**クエリパラメータ**

- `status` (optional): フィルタ（`active|paused|done`）

**レスポンス**

```json
{
  "goals": [
    {
      "id": "goal-456",
      "title": "週10時間の集中作業",
      "description": "...",
      "startDate": "2025-10-01",
      "endDate": "2025-12-31",
      "kpi_name": "集中作業時間",
      "kpi_target": 10,
      "kpi_unit": "時間",
      "status": "active",
      "createdAt": "2025-10-01T00:00:00Z",
      "updatedAt": "2025-10-01T00:00:00Z"
    }
  ]
}
```

### POST /goals

目標作成。

**リクエスト**

```json
{
  "title": "週10時間の集中作業",
  "description": "...",
  "startDate": "2025-10-01",
  "endDate": "2025-12-31",
  "kpi_name": "集中作業時間",
  "kpi_target": 10,
  "kpi_unit": "時間"
}
```

**レスポンス**

```json
{
  "goal": {
    "id": "goal-456",
    ...
  }
}
```

### PATCH /goals/:id

目標更新。

**リクエスト**

```json
{
  "status": "done"
}
```

**レスポンス**

```json
{
  "goal": {
    "id": "goal-456",
    ...
  }
}
```

## キャプチャ

### GET /screenshots

スクリーンショット一覧取得。

**クエリパラメータ**

- `limit` (optional): 取得件数（デフォルト: 50）
- `offset` (optional): オフセット
- `taskId` (optional): タスク ID でフィルタ

**レスポンス**

```json
{
  "screenshots": [
    {
      "id": "screenshot-789",
      "path": "/path/to/screenshot.png",
      "thumbPath": "/path/to/thumbnail.jpg",
      "capturedAt": "2025-10-29T10:30:00Z",
      "mode": "scheduled",
      "meta": {},
      "linkedTaskId": "task-123"
    }
  ]
}
```

### GET /screenshots/:id

スクリーンショット画像ファイル配信。

**レスポンス**

画像ファイル（PNG/JPEG）

**エラー**

- `404 Not Found` - ファイルが存在しない

### GET /capture/schedule

現行キャプチャスケジュール取得。

**レスポンス**

```json
{
  "schedule": {
    "id": "schedule-1",
    "active": true,
    "intervalMin": 5,
    "retention_maxItems": 1000,
    "retention_maxDays": 30,
    "updatedAt": "2025-10-29T10:00:00Z"
  }
}
```

```json
{
  "schedule": null
}
```

### PUT /capture/schedule

キャプチャスケジュール作成/更新。

**リクエスト**

```json
{
  "active": true,
  "intervalMin": 5,
  "retention_maxItems": 1000,
  "retention_maxDays": 30
}
```

**レスポンス**

```json
{
  "schedule": {
    "id": "schedule-1",
    ...
  }
}
```

### POST /capture/schedule/start

定期キャプチャを有効化。

**レスポンス**

```json
{
  "message": "started"
}
```

**エラー**

- `403 Forbidden` - 権限が未許可
  ```json
  { "code": "PERMISSION_DENIED", "message": "Screen capture not allowed" }
  ```

### POST /capture/schedule/stop

定期キャプチャを停止。

**レスポンス**

```json
{
  "message": "stopped"
}
```

## 設定

### GET /settings

設定取得。

**レスポンス**

```json
{
  "settings": {
    "storagePath": "/home/user/llm-time-manager",
    "capturePath": "/home/user/screenshots",
    "captureFormat": "png",
    "thumbnailResolution": 320,
    "debugLog": false
  }
}
```

### PATCH /settings

設定更新。

**リクエスト**

```json
{
  "captureFormat": "jpg",
  "debugLog": true
}
```

**レスポンス**

```json
{
  "settings": {
    ...
  }
}
```

## 共通エラーレスポンス

### エラーフォーマット

```json
{
  "code": "ERROR_CODE",
  "message": "Human readable error message"
}
```

### エラーコード一覧

- `PERMISSION_DENIED` - 権限が拒否された
- `BRIDGE_UNAVAILABLE` - ネイティブブリッジが初期化されていない
- `NOT_FOUND` - リソースが見つからない
- `INVALID_REQUEST` - リクエストが不正
- `INTERNAL_ERROR` - サーバ内部エラー

## レート制限

プロトタイプのため、レート制限は実装しない。
