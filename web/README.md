# LLM時間管理ツール - WebView

React + React Routerを使用したフロントエンドアプリケーション

## 開発環境のセットアップ

### 1. Nix devshellに入る

```bash
cd web
nix develop
```

### 2. GitHub CLI認証（初回のみ）

```bash
gh auth login
```

### 3. Issueを作成（認証後）

```bash
gh issue create \
  --title "開発環境のセットアップ（Nix devshell + フォーマッター + React Router）" \
  --body-file .github-issues/setup-devshell-and-tools.md \
  --label "setup,infrastructure,frontend"
```

## 開発ツール

devshellには以下のツールが含まれています：

- **Node.js** (v22.20.0) - JavaScript/TypeScript実行環境
- **pnpm** (10.19.0) - パッケージマネージャ
- **git** (2.51.0) - バージョン管理
- **gh** (2.82.1) - GitHub CLI
- **biome** (2.2.7) - TypeScript/TSX/JSON/HTML/CSS用フォーマッター/リンター
- **prettier** (3.6.2) - その他ファイル用フォーマッター

## コードフォーマット/リント

### Biome（TypeScript/TSX/JSON/HTML/CSS）

設定ファイル: `biome.json`

```bash
# チェックのみ
pnpm lint

# 自動修正
pnpm lint:fix

# フォーマット
pnpm format
```

### Prettier（Markdown/YAML等）

設定ファイル: `.prettierrc`, `.prettierignore`

```bash
# フォーマット
pnpm format:other
```

## プロジェクト構成

```
web/
├── src/
│   ├── main.tsx              # エントリーポイント
│   ├── App.tsx               # ルートコンポーネント
│   ├── routes/               # ページコンポーネント
│   │   ├── index.tsx         # / (ホーム)
│   │   ├── chat.tsx          # /chat (LLMチャット)
│   │   ├── goals.tsx         # /goals (目標一覧)
│   │   ├── tasks.tsx         # /tasks (タスク一覧)
│   │   ├── capture.tsx       # /capture (キャプチャ設定)
│   │   └── settings.tsx      # /settings/local (設定)
│   ├── components/           # 共通コンポーネント
│   ├── lib/                  # ユーティリティ
│   └── styles/               # グローバルスタイル
├── public/
├── flake.nix                 # Nix開発環境定義
├── biome.json                # Biome設定
├── .prettierrc               # Prettier設定
└── package.json              # Node.jsプロジェクト定義
```

## 次のステップ

1. 依存関係のインストール

   ```bash
   pnpm init
   pnpm add react react-dom react-router-dom
   pnpm add -D @types/react @types/react-dom typescript vite @vitejs/plugin-react
   pnpm add -D @biomejs/biome prettier
   ```

2. Vite設定とTypeScript設定の作成

3. React Routerのセットアップ

詳細は [.github-issues/setup-devshell-and-tools.md](.github-issues/setup-devshell-and-tools.md)
を参照してください。
