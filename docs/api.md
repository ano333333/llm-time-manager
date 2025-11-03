# API 仕様

ローカルサーバ（Go）が提供する REST + WebSocket API。

ベース URL: `http://localhost:<port>`

## LLM

### POST /llm/chat

LLM とのチャット（Server-Sent Events または WebSocket でストリーム応答）。

#### request

```json
{
  "messages": [{ "role": "user", "content": "来週水曜にレポートを提出したい" }]
}
```

#### response: 200

```
data: {"type":"text","content":"わかりました"}
data: {"type":"text","content":"。"}
data: {"type":"entity","entity":{"type":"task","title":"レポート提出","due":"2025-11-05"}}
data: [DONE]
```

#### response: error

- `500 Internal Server Error` - LLM エンジンエラー

## タスク

### GET /tasks

タスク一覧取得。

#### query parameter

- `status` (optional): フィルタ（`todo|doing|paused|done`）
- `due` (optional): 期日フィルタ（`today|week|overdue`）
- `goalId` (optional): 目標 ID でフィルタ

#### response: 200

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

#### request

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

#### response: 200

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

#### response: 200

```json
{
  "task": {
    "id": "task-123",
    ...
  }
}
```

#### response: error

- `404 Not Found` - タスクが存在しない

### PATCH /tasks/:id

タスク更新。

#### request

```json
{
  "status": "doing",
  "priority": 4
}
```

#### response: 200

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

#### response

```json
{
  "message": "deleted"
}
```

## 目標

### GET /goals

目標一覧取得。

#### query parameter

- `status` (optional): フィルタ（`active|paused|done`）、カンマ区切り

```
GET /goals?status=active,paused
```

#### response: 200

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

#### request

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

#### response: 200

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

#### request

```json
{
  "status": "done"
}
```

#### response: 200

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

#### query parameter

- `limit` (optional): 取得件数（デフォルト: 50）
- `offset` (optional): オフセット
- `taskId` (optional): タスク ID でフィルタ

#### response: 200

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

#### response: 200

画像ファイル（PNG/JPEG）

#### response: error

- `404 Not Found` - ファイルが存在しない

### GET /capture/schedule

現行キャプチャスケジュールとして、アクティブなスケジュールを 1 つだけ取得する。

#### response: 200

- アクティブなスケジュール存在時

```json
{
  "schedule": {
    "id": "schedule-1",
    "active": true,
    "intervalMin": 5,
    "retention_max_items": 1000,
    "retention_max_days": 30
  }
}
```

- アクティブなスケジュール非存在時

```json
{
  "schedule": null
}
```

#### response: error

- 500: 内部エラー時

### PUT /capture/schedule

キャプチャスケジュール作成/更新。

#### request

```json
{
  "active": true,
  "intervalMin": 5,
  "retention_max_items": 1000,
  "retention_max_days": 30
}
```

#### response: 200

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

#### response: 200

```json
{
  "message": "started"
}
```

#### response: error

- `403 Forbidden` - 権限が未許可
  ```json
  { "code": "PERMISSION_DENIED", "message": "Screen capture not allowed" }
  ```

### POST /capture/schedule/stop

定期キャプチャを停止。

#### response: 200

```json
{
  "message": "stopped"
}
```

## 設定

### GET /settings

設定取得。

#### response: 200

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

#### request

```json
{
  "captureFormat": "jpg",
  "debugLog": true
}
```

#### response: 200

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
