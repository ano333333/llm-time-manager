{ pkgs }:

with pkgs;
[
  # バージョン管理
  git
  gh

  # Go開発環境（server用）
  go
  gopls # Go language server
  gotools # goimports, godoc等のツール
  go-tools # staticcheck等の追加ツール

  # Goフォーマッタ
  gofumpt # gofmtの改善版
  golines # 長い行を分割

  # Goリンター
  golangci-lint

  # データベース
  sqlite
  sqlc # SQLからGoコード生成
  goose # DBマイグレーション

  # Go開発ツール
  air # ホットリロード
  delve # Goデバッガー

  # Goテスト・ベンチマーク
  go-junit-report # JUnit形式のテストレポート

  # Node.js環境（web用）
  nodejs_22

  # パッケージマネージャ
  nodePackages.pnpm

  # コードフォーマッター/リンター
  biome
  nodePackages.prettier

  # Playwright（E2E/コンポーネントテスト用）
  playwright-driver.browsers # Playwright browsers

  # Nix開発環境
  nil # Nix language server
  nixpkgs-fmt # Nixフォーマッター

  # その他便利ツール
  jq # JSONパーサー
  curl # API テスト
]
