# LLM時間管理ツール

v1.0.0

## 概要

各クライアントデバイス（モバイル/デスクトップ）から、LLMチャットを通じたタスク・目標設定とスクリーンショット取得（以降「キャプチャ」）を行い、ユーザーの時間管理を支援する個人用ローカルサーバープロトタイプです。

## 主な特徴

- **完全ローカル動作**: オフライン環境でも機能継続
- **LLMによるタスク・目標管理**: 自然言語での対話的な時間管理
- **定期スクリーンショット**: 作業状況の自動記録
- **マルチプラットフォーム対応**: Linux/Windows/macOS/iOS/iPadOS

## 前提

- クライアントアプリは **WebView を共通 UI レイヤ**とし、ローカルサーバー（同一端末上、`http://localhost:<port>`）に接続
- **個人用ローカルサーバープロトタイプ**。アカウント登録/ログイン/ヘルプは不要
- LLM 推論はローカルor LAN 上のエンジンを想定（外部連携は任意）
- 端末固有のキャプチャ許可/実行はネイティブブリッジ経由で行う（WebView ⇄ ネイティブ）

## 技術スタック

- **フロントエンド**: React + React Router (SSG)
- **バックエンド**: Go (net/http + chi/gin)
- **データベース**: SQLite
- **通信**: REST + WebSocket/SSE
- **ネイティブ層**:
  - Linux (Ubuntu): Rust
  - Windows: C++ (Win32/COM、Graphics Capture API)
  - macOS: Swift (ScreenCaptureKit)
  - iOS/iPadOS: Swift (ReplayKit / UIScreenCapture API)

## ドキュメント

- [アーキテクチャ](docs/architecture.md) - システム構成、非機能要件
- [UI設計](docs/ui-design.md) - 画面一覧、遷移、ワイヤーフレーム
- [データモデル](docs/data-model.md) - テーブル設計、ER図
- [API仕様](docs/api.md) - エンドポイント一覧
- [主要フロー](docs/flows.md) - シーケンス図
- [プラットフォーム固有実装](docs/platform-specific.md) - OS依存挙動、権限管理
- [状態機械](docs/state-machines.md) - タスク・キャプチャの状態遷移
- [ブリッジインターフェース](docs/bridge-interface.md) - WebView-ネイティブ間I/F
- [実装ガイド](docs/implementation.md) - 実装メモ、テスト観点

## リポジトリ構成

```
repo/
  server/                # Go バックエンド
    cmd/api/            # main パッケージ
    internal/
      http/             # ルータ/ハンドラ（chi/gin）
      ws/               # WebSocket/SSE
      store/            # SQLite アクセス（sqlc/gorm）
      capture/          # Bridge 呼び出し・スケジューラ
      llm/              # LLM クライアント
      config/
      logging/
    migrations/         # DB マイグレーション
    go.mod go.sum

  web/                  # React + React Router (SSG)
    src/
      routes/           # /, /chat, /goals, /tasks, /capture, /settings/local
      components/
      lib/
      styles/
    public/
    package.json
    vite.config.ts      # or Next.js の静的輸出

  clients/
    linux-rust/
      src/              # Wayland/X11 キャプチャ、bridge 実装
      Cargo.toml
    windows-cpp/
      src/              # Graphics Capture API、bridge DLL/EXE
      CMakeLists.txt
    ios-swift/
      App/              # WebView + Bridge（WKScriptMessageHandler）
      Shared/
      Package.swift
    macos-swift/
      App/              # ScreenCaptureKit + WebView
      Shared/
      Package.swift

  shared/
    schemas/            # Zod/JSON Schema（Task/Goal 等）
    ui/                 # 共通デザイン/アイコン

  tools/
    scripts/            # 開発用スクリプト

  docs/                 # ドキュメント
  README.md
```

## 開発・配布

- **server**: 単体バイナリ（Go）
- **web**: `web/dist` を server の静的配信へ同梱
- **clients**: 各 OS 向けバンドル（Deb/DMG/EXE/TestFlight 等）

## 将来拡張

- 外部カレンダー連携（読み取りのみ）
- 音声要約→タスク化
- 目標の自動進捗推定（キャプチャ/完了イベントから）
- エクスポート（Markdown/CSV）

## ライセンス

TBD

