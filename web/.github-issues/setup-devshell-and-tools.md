---
title: '開発環境のセットアップ（Nix devshell + フォーマッター + React Router）'
labels: ['setup', 'infrastructure', 'frontend']
assignees: []
---

## 概要

webviewリポジトリの開発環境を構築し、必要なツールとライブラリを導入する。

## タスク

### 1. Nix devshellのセットアップ

- [x] `flake.nix`の作成
  - Node.js (v22)
  - pnpm
  - git
  - gh (GitHub CLI)
  - biome
  - prettier
- [x] `.gitignore`の作成
- [x] devshellの動作確認
  ```bash
  nix develop
  ```

### 2. コードフォーマッター/リンターのセットアップ

#### Biome（TypeScript/TSX/JSON/HTML/CSS）

- [x] `biome.json`の作成
  - 適用対象: `.ts`, `.tsx`, `.json`, `.html`, `.css`
  - フォーマット設定: 2スペースインデント、行幅100、セミコロン有効
  - リンタールール: recommended有効、suspicious/styleカスタマイズ
- [x] package.jsonにスクリプト追加
  ```json
  {
    "scripts": {
      "lint": "biome check .",
      "lint:fix": "biome check --write .",
      "format": "biome format --write ."
    }
  }
  ```

#### Prettier（その他のファイル）

- [x] `.prettierrc`の作成
  - Markdown, YAML, その他のファイル用
  - Biomeと同じフォーマット設定（一貫性維持）
- [x] `.prettierignore`の作成
  - Biome管轄ファイルを除外（重複を避ける）
- [x] package.jsonにスクリプト追加
  ```json
  {
    "scripts": {
      "format:other": "prettier --write '**/*.{md,yaml,yml}'"
    }
  }
  ```

### 3. React Routerのセットアップ

#### 依存関係のインストール

- [x] 依存関係のインストール完了
  ```bash
  cd web
  pnpm init
  pnpm add react react-dom react-router-dom
  pnpm add -D @types/react @types/react-dom typescript vite @vitejs/plugin-react
  pnpm add -D @biomejs/biome prettier
  ```

#### プロジェクト構造の作成

- [x] プロジェクト構造作成完了
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
  ├── index.html
  ├── vite.config.ts
  ├── tsconfig.json
  └── package.json
  ```

#### ルーティング設定

- [x] `src/main.tsx`でBrowserRouterを設定
- [x] `src/App.tsx`でRoute定義
  - `/` - ホーム画面
  - `/chat` - LLMチャット
  - `/goals` - 目標一覧
  - `/tasks` - タスク一覧
  - `/capture` - キャプチャ設定
  - `/settings/local` - ローカル設定
- [x] 各ルート用の基本コンポーネント作成

#### Vite設定

- [x] `vite.config.ts`作成
- [x] `tsconfig.json`作成
- [x] ビルドスクリプト確認
  ```json
  {
    "scripts": {
      "dev": "vite",
      "build": "tsc && vite build",
      "preview": "vite preview"
    }
  }
  ```

### 4. 動作確認

- [x] `nix develop`でdevshell起動
- [x] `cd web && pnpm install`で依存関係インストール
- [x] `pnpm dev`で開発サーバー起動
- [x] ブラウザで`http://localhost:5173`アクセス
- [x] ルーティング動作確認（各ページへの遷移）
- [x] `pnpm lint`でBiomeの動作確認
- [x] `pnpm format:other`でPrettierの動作確認

## 受け入れ基準

- [x] Nix devshellが正常に起動し、すべてのツールが利用可能
- [x] BiomeとPrettierが正しく設定され、重複なく動作
- [x] React Routerが設定され、すべてのルートが表示される
- [x] 開発サーバーが起動し、Hot Module Replacement (HMR)が動作
- [x] コードフォーマット/リンタースクリプトが正常に実行される

## 参考資料

- [Biome公式ドキュメント](https://biomejs.dev/)
- [Prettier公式ドキュメント](https://prettier.io/)
- [React Router公式ドキュメント](https://reactrouter.com/)
- [Vite公式ドキュメント](https://vitejs.dev/)
- [Nix Flakes](https://nixos.wiki/wiki/Flakes)
