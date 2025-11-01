# LLM 時間管理ツール - Go サーバー

LLM 時間管理ツールのバックエンドサーバーです。Go 言語で実装されており、REST API、WebSocket、SQLite
データベースを使用しています。

## 開発環境のセットアップ

### 前提条件

- [Nix](https://nixos.org/download.html) がインストールされていること
- Flakes が有効化されていること
- （推奨）[direnv](https://direnv.net/) がインストールされていること

### 開発環境に入る

#### direnv を使用する場合（推奨）

direnv を使用すると、ディレクトリに入るだけで自動的に Nix 環境がロードされます：

```bash
cd server
direnv allow  # 初回のみ
# 以降、serverディレクトリに入ると自動的にNix環境がロード
```

VSCode を使用する場合
、[direnv 拡張機能](https://marketplace.visualstudio.com/items?itemName=mkhl.direnv)をインストールす
ることで、エディタ内でも Nix 環境が自動的にロードされます。

#### 手動で devshell に入る場合

```bash
cd server
nix develop
```

devshell に入ると、以下のツールが利用可能になります：

- **Go**: Go 言語コンパイラ（v1.25+）
- **gopls**: Go 言語サーバー
- **gofumpt**: コードフォーマッタ
- **golangci-lint**: リンター
- **air**: ホットリロード開発サーバー
- **delve**: Go デバッガー
- **sqlite**: SQLite データベース
- **sqlc**: SQL から Go コード生成
- **goose**: DB マイグレーションツール

## プロジェクト構造

```text
server/
├── cmd/
│   └── api/              # メインエントリーポイント
├── internal/
│   ├── http/            # HTTPハンドラとルーター
│   ├── ws/              # WebSocket/SSEハンドラ
│   ├── store/           # データベースアクセス層
│   ├── capture/         # スクリーンショットキャプチャ管理
│   ├── llm/             # LLMクライアント
│   ├── config/          # 設定管理
│   └── logging/         # ロギング
├── migrations/          # データベースマイグレーション
├── go.mod
└── go.sum
```

## 開発コマンド

### ビルド

```bash
make build
```

または

```bash
go build -o bin/api ./cmd/api
```

### 実行

```bash
make run
```

または

```bash
go run ./cmd/api
```

### 開発モード（ホットリロード）

```bash
make dev
```

または

```bash
air
```

### テスト

```bash
make test
```

テストカバレッジ付き：

```bash
make test-coverage
```

### リンター

```bash
make lint
```

### フォーマッタ

```bash
make fmt
```

### クリーン

```bash
make clean
```

## 設定

`config.example.yaml` をコピーして `config.local.yaml` を作成し、環境に合わせて編集してください：

```bash
cp config.example.yaml config.local.yaml
```

主な設定項目：

- **server**: ホスト、ポート、タイムアウト設定
- **database**: SQLite のパス、接続プール設定
- **llm**: LLM エンドポイント、モデル設定
- **capture**: スクリーンショット保存先、キャプチャ間隔
- **logging**: ログレベル、出力先

## データベースマイグレーション

このプロジェクトでは [pressly/goose](https://github.com/pressly/goose) を使用してデータベースマイグ
レーションを管理しています。

### 新しいマイグレーションの作成

```bash
# goose CLI を使用（推奨）
goose -dir migrations create <migration_name> sql
```

これにより、タイムスタンプ付きのマイグレーションファイルが自動生成されます。

## API 仕様

詳細な API 仕様は [docs/api.md](../docs/api.md) を参照してください。

### 主なエンドポイント

- `GET /api/tasks` - タスク一覧取得
- `POST /api/tasks` - タスク作成
- `GET /api/tasks/:id` - タスク詳細取得
- `PATCH /api/tasks/:id` - タスク更新
- `DELETE /api/tasks/:id` - タスク削除
- `GET /api/goals` - 目標一覧取得
- `POST /api/goals` - 目標作成
- `GET /api/captures` - キャプチャ一覧取得
- `POST /api/captures` - 手動キャプチャ実行
- `POST /api/llm/chat` - LLM チャット（SSE ストリーム）
- `WS /api/ws` - WebSocket 接続

## テスト

### ユニットテスト

各パッケージ内でテストファイルを作成します：

```go
// example_test.go
package example

import "testing"

func TestExample(t *testing.T) {
    // テストコード
}
```

### 統合テスト

```bash
go test -tags=integration ./...
```

### カバレッジ

```bash
make test-coverage
# coverage.html がブラウザで開けます
```

## デバッグ

### Delve を使用したデバッグ

```bash
dlv debug ./cmd/api
```

### ログレベルの変更

`config.local.yaml` で `logging.level` を `debug` に設定：

```yaml
logging:
  level: "debug"
```

## パフォーマンス

### プロファイリング

```bash
go test -cpuprofile=cpu.prof -memprofile=mem.prof -bench=.
go tool pprof cpu.prof
```

### ベンチマーク

```bash
go test -bench=. -benchmem ./...
```

## コーディング規約

### フォーマット

- `gofumpt` を使用してフォーマット（`make fmt`）
- タブインデント（Go 標準）
- 1 行あたり最大 120 文字を推奨

### 命名規則

- パッケージ名: 小文字、単語区切りなし（例: `httphandler`）
- 関数名: キャメルケース（例: `GetTaskByID`）
- 変数名: キャメルケース（例: `taskID`）
- 定数名: キャメルケースまたは大文字スネークケース（例: `MaxRetries` または `MAX_RETRIES`）

### エラーハンドリング

- エラーは必ず処理する
- エラーメッセージは小文字で始める
- エラーラップには `fmt.Errorf("context: %w", err)` を使用

```go
if err != nil {
    return fmt.Errorf("failed to get task: %w", err)
}
```

### ロギング

- 構造化ログを使用（`zap` または `zerolog`）
- PII を含めない
- ログレベルを適切に設定

```go
logger.Info("task created",
    zap.String("task_id", task.ID),
    zap.String("status", task.Status),
)
```

## デプロイ

### バイナリのビルド

```bash
CGO_ENABLED=1 go build -o bin/api ./cmd/api
```

### 配布物

- `bin/api`: サーバーバイナリ
- `config.example.yaml`: 設定ファイルサンプル
- `migrations/`: データベースマイグレーション

## トラブルシューティング

### SQLite エラー

```text
database is locked
```

→ 接続プールの設定を確認してください（`database.max_open_conns`）

### ポート競合

```text
bind: address already in use
```

→ `config.local.yaml` の `server.port` を変更してください

### LLM 接続エラー

```text
dial tcp: connection refused
```

→ LLM サーバーが起動しているか確認し、`llm.endpoint` の設定を確認してください

## 参考資料

- [Go 公式ドキュメント](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)
- [プロジェクト全体ドキュメント](../docs/)

## ライセンス

MIT License - 詳細は [LICENSE](../LICENSE) を参照してください。
