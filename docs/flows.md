# 主要フロー

システムの主要な処理フローをシーケンスで示す。

## チャットからタスク作成

ユーザーがチャットでタスクの相談をし、LLMが提案したタスクを作成する流れ。

```mermaid
sequenceDiagram
    participant User
    participant WebView
    participant LocalServer
    participant LLM
    participant DB

    User->>WebView: 「来週水曜にレポート草案を提出」
    WebView->>LocalServer: POST /llm/chat
    LocalServer->>LLM: チャットメッセージ
    LLM-->>LocalServer: ストリーム応答（構造化提案含む）
    LocalServer-->>WebView: SSE: {task:{title,due,estimate}}
    WebView->>User: タスク作成モーダル表示（編集可）
    User->>WebView: 確認・編集後「保存」
    WebView->>LocalServer: POST /tasks
    LocalServer->>DB: INSERT task
    DB-->>LocalServer: OK
    LocalServer-->>WebView: タスクデータ
    WebView->>User: /tasks/:id へ遷移
```

### 処理ステップ

1. User→Chat: 「来週水曜にレポート草案を…」
2. LLM→Front: 構造化提案 `{task:{title,due,estimate}}`
3. Front: タスク作成モーダルを起動（編集可）
4. API: `POST /tasks` → 保存 → 画面遷移 `/tasks/:id`

## 定期キャプチャの開始/停止

定期キャプチャのスケジュール設定と実行の流れ。

```mermaid
sequenceDiagram
    participant User
    participant OS
    participant NativeBridge
    participant WebView
    participant LocalServer
    participant Scheduler
    participant DB

    User->>WebView: 「5分おきに実行」設定
    WebView->>LocalServer: PUT /capture/schedule {intervalMin:5, active:true}
    LocalServer->>DB: UPDATE schedule
    LocalServer->>Scheduler: タイマー起動
    LocalServer-->>WebView: OK

    loop 5分ごと
        Scheduler->>NativeBridge: captureScreenshot()
        NativeBridge->>OS: キャプチャAPI呼び出し
        OS-->>NativeBridge: 画像データ
        NativeBridge-->>Scheduler: {path, thumbPath}
        Scheduler->>DB: INSERT screenshot
        Scheduler->>WebView: WebSocket通知
    end

    User->>WebView: 「停止」
    WebView->>LocalServer: POST /capture/schedule/stop
    LocalServer->>Scheduler: タイマー停止
    LocalServer-->>WebView: OK
```

### 処理ステップ

1. User→Capture: 「5分おきに実行」→ `PUT /capture/schedule {intervalMin:5, active:true}`
2. Server: スケジューラ起動（Go でタイマー）
3. Server→NativeBridge: 間隔ごとに `captureScreenshot()`
4. User→Capture: 「停止」→ `POST /capture/schedule/stop`

## 目標→タスク化

既存の目標からタスクを作成する流れ。

```mermaid
sequenceDiagram
    participant User
    participant WebView
    participant LocalServer
    participant DB

    User->>WebView: /goals で目標選択
    User->>WebView: 「タスクを追加」ボタン
    WebView->>User: タスク作成モーダル表示（goalId付き）
    User->>WebView: タスク詳細入力後「保存」
    WebView->>LocalServer: POST /tasks {goalId: "goal-456", ...}
    LocalServer->>DB: INSERT task with goalId
    DB-->>LocalServer: OK
    LocalServer-->>WebView: タスクデータ
    WebView->>User: /tasks?goalId=goal-456 へ遷移（自動フィルタ）
```

### 処理ステップ

1. `/goals` で目標選択→「タスクを追加」
2. `POST /tasks` goalId 指定
3. タスク一覧で目標フィルタを自動適用

## 権限取得フロー

画面キャプチャ権限を取得する流れ。

```mermaid
sequenceDiagram
    participant User
    participant OS
    participant NativeBridge
    participant WebView
    participant LocalServer

    User->>WebView: /capture アクセス
    WebView->>LocalServer: GET /capture/schedule
    LocalServer-->>WebView: {active: false, ...}
    WebView->>NativeBridge: requestPermission('capture')
    NativeBridge->>OS: 権限確認
    
    alt 権限未許可
        OS-->>NativeBridge: 'not_determined'
        NativeBridge-->>WebView: 'not_determined'
        WebView->>User: 権限説明モーダル表示
        User->>WebView: 「許可ボタン」
        WebView->>NativeBridge: requestPermission('capture')
        NativeBridge->>OS: システムダイアログ表示
        OS->>User: 権限ダイアログ表示
        
        alt ユーザーが許可
            User->>OS: 許可
            OS-->>NativeBridge: 'granted'
            NativeBridge-->>WebView: 'granted'
            WebView->>User: キャプチャ設定UI有効化
        else ユーザーが拒否
            User->>OS: 拒否
            OS-->>NativeBridge: 'denied'
            NativeBridge-->>WebView: 'denied'
            WebView->>User: エラー表示＋設定遷移案内
        end
    else 既に許可済み
        OS-->>NativeBridge: 'granted'
        NativeBridge-->>WebView: 'granted'
        WebView->>User: キャプチャ設定UI有効化
    end
```

## タスク状態変更フロー

タスクのステータスを変更する流れ。

```mermaid
sequenceDiagram
    participant User
    participant WebView
    participant LocalServer
    participant DB

    User->>WebView: タスク詳細画面で「開始」
    WebView->>LocalServer: PATCH /tasks/:id {status: "doing"}
    LocalServer->>DB: UPDATE task SET status='doing'
    DB-->>LocalServer: OK
    LocalServer-->>WebView: 更新されたタスクデータ
    WebView->>User: UI更新（ステータス表示変更）

    Note over User,DB: 作業中...

    User->>WebView: 「完了」
    WebView->>LocalServer: PATCH /tasks/:id {status: "done"}
    LocalServer->>DB: UPDATE task SET status='done'
    DB-->>LocalServer: OK
    LocalServer-->>WebView: 更新されたタスクデータ
    WebView->>User: UI更新（完了マーク、進捗更新）
```

## チャットストリーミング

LLMとのリアルタイムチャット通信。

```mermaid
sequenceDiagram
    participant User
    participant WebView
    participant LocalServer
    participant LLM

    User->>WebView: メッセージ入力＆送信
    WebView->>LocalServer: POST /llm/chat (SSE接続)
    LocalServer->>LLM: リクエスト送信

    loop ストリーム応答
        LLM-->>LocalServer: トークンチャンク
        LocalServer-->>WebView: SSE: data: {"type":"text","content":"..."}
        WebView->>User: UI更新（逐次表示）
    end

    LLM-->>LocalServer: エンティティ情報
    LocalServer-->>WebView: SSE: data: {"type":"entity","entity":{...}}
    WebView->>User: アクションボタン表示（タスク化など）

    LLM-->>LocalServer: 完了
    LocalServer-->>WebView: SSE: data: [DONE]
    WebView->>LocalServer: 接続終了
```

## エラー処理フロー

キャプチャ失敗時のリトライ処理。

```mermaid
sequenceDiagram
    participant User
    participant NativeBridge
    participant WebView
    participant LocalServer
    participant Scheduler
    participant DB

    Scheduler->>NativeBridge: captureScreenshot()
    NativeBridge-->>Scheduler: Error: PERMISSION_DENIED
    Scheduler->>DB: INSERT error_log
    Scheduler->>WebView: WebSocket通知: {type: "error", code: "PERMISSION_DENIED"}
    WebView->>User: トースト表示「キャプチャ失敗」
    WebView->>User: リトライボタン表示

    User->>WebView: リトライボタン押下
    WebView->>NativeBridge: requestPermission('capture')
    
    alt 権限が回復
        NativeBridge-->>WebView: 'granted'
        WebView->>LocalServer: POST /capture/schedule/start
        LocalServer->>Scheduler: 再起動
    else まだ拒否状態
        NativeBridge-->>WebView: 'denied'
        WebView->>User: 設定画面への遷移案内
    end
```

