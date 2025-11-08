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

### GET /goal

目標一覧取得。

#### query parameter

- `status` (optional): フィルタ（`active|paused|done`）、カンマ区切り。指定されない、または空白文字列の場合、空配列を返す。

```
GET /goal?status=active,paused
```

#### response: 200

```json
{
  "goals": [
    {
      "id": "goal-456",
      "title": "週10時間の集中作業",
      "description": "...",
      "start_date": "2025-10-01",
      "end_date": "2025-12-31",
      "kpi_name": "集中作業時間",
      "kpi_target": 10,
      "kpi_unit": "時間",
      "status": "active",
      "created_at": "2025-10-01T00:00:00Z",
      "updated_at": "2025-10-01T00:00:00Z"
    }
  ]
}
```

```ts
{
  goals: Array<{
    id: string,
    title: string,
    description: string,
    start_date: string,
    end_date: string,
    kpi_name: string | null,
    kpi_target: number | null,
    kpi_unit: string | null,
    status: "active" | "paused" | "done",
    created_at: string,
    updated_at: string,
  }>,
}
```

- start_date, end_date は`YYYY-MM-DD`形式である
- created_at, updated_at は ISO8601 形式である
- kpi_name, kpi_target, kpi_unit はすべて null かすべて非 null かのいずれかである
- goals の要素は id 昇順

#### response: error

- 400: query parameter が不正
- 500: 内部エラー

### POST /goal

目標作成。

#### request

```json
{
  "title": "週10時間の集中作業",
  "description": "...",
  "start_date": "2025-10-01",
  "end_date": "2025-12-31",
  "kpi_name": "集中作業時間",
  "kpi_target": 10,
  "kpi_unit": "時間",
  "status": "active"
}
```

```ts
{
  title: string,
  description: string,
  start_date: string,
  end_date: string,
  kpi_name: string | null,
  kpi_target: number | null,
  kpi_unit: string | null,
  status: "active"|"paused"|"done",
}
```

- title は空白文字(`\s`)のみで構成されてはならない
- start_date と end_date は`"YYYY-MM-DD"`形式の string
- start_date は end_date 以下
- kpi_name, kpi_target, kpi_unit のいずれかが非 null ならば、それ以外の値もすべて非 null である
- kpi_name と kpi_unit は、string ならば空白文字のみで構成されてはならない

#### response: 200

```json
{
  "goal": {
    "id": "goal-456",
    ...
  }
}
```

```ts
{
  id: string,
  title: string,
  description: string,
  start_date: string,
  end_date: string,
  kpi_name: string | null,
  kpi_target: number | null,
  kpi_unit: string | null,
  status: "active" | "paused" | "done",
  created_at: string,
  updated_at: string,
}
```

- start_date, end_date は`YYYY-MM-DD`形式である
- created_at, updated_at は ISO8601 形式である
- kpi_name, kpi_target, kpi_unit はすべて null かすべて非 null かのいずれかである

### PATCH /goal/:id

目標更新。

#### request

```json
{
  "status": "done"
}
```

```ts
{
  title?: string,
  description?: string,
  start_data?: string,
  end_date?: string,
  kpi_name?: string,
  kpi_target?: number,
  kpi_unit?: string,
  status?: "active"|"paused"|"done",
}
```

- title は空白文字(`\s`)のみで構成されてはならない
- start_date, end_date は`"YYYY-MM-DD"`形式の string
- start_date が end_date より大きくなる更新は適用されない
- kpi_name, kpi_target, kpi_unit のうち 1 つ以上が null かつ 1 つ以上が非 null になる更新は適用されない
- kpi_name, kpi_unit は string ならば空白文字のみで構成されてはならない

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

### POST /capture/screenshot

スクリーンショットをサーバーに送信し、LLM で分析する。

#### request

- `Content-Type: multipart/form-data`
- `image`: スクリーンショット画像ファイル（PNG/JPEG）

#### response: 200

```json
{
  "message": "LLM による分析結果のテキスト",
  "analysis": {
    "summary": "作業内容の要約",
    "suggestions": ["提案1", "提案2"]
  }
}
```

#### response: error

- `400 Bad Request` - リクエストが不正な場合（画像が含まれていない等）
  ```json
  { "code": "INVALID_REQUEST", "message": "Image file is required" }
  ```
- `403 Forbidden` - 権限が未許可の場合
  ```json
  { "code": "PERMISSION_DENIED", "message": "Screen capture not allowed" }
  ```
- `500 Internal Server Error` - LLM エラーまたは内部エラー
  ```json
  { "code": "INTERNAL_ERROR", "message": "Failed to analyze screenshot" }
  ```

### GET /capture/schedule

現行キャプチャスケジュールとして、アクティブなスケジュールを 1 つだけ取得する。

#### response: 200

- アクティブなスケジュール存在時

```json
{
  "schedule": {
    "id": "schedule-1",
    "active": true,
    "intervalMin": 5
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
  "interval_min": 5
}
```

```ts
{
  active: boolean,
  interval_min: number,
}
```

- `interval_min`は正整数

#### response: 200

```json
{
  "schedule": {
    "id": "schedule-1",
    "active": true,
    "interval_min": 5
  }
}
```

#### response: error

- `400 Bad Request` - リクエストパラメータが不正な場合

```json
{
  "message": "invalid parameter",
  "target": "id"
}
```

- `400 Bad Request` - アクティブなスケジュールがない場合

```json
{
  "message": "no active capture schedule"
}
```

- `500 Internal Server Error` - 内部エラー時

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
