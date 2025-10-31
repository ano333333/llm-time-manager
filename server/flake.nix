{
  description = "LLM時間管理ツール - Go Server Development Environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            # Go開発環境
            go
            gopls           # Go language server
            gotools         # goimports, godoc等のツール
            go-tools        # staticcheck等の追加ツール
            
            # フォーマッタ
            gofumpt         # gofmtの改善版
            golines         # 長い行を分割
            
            # リンター
            golangci-lint
            
            # データベース
            sqlite
            sqlc            # SQLからGoコード生成
            goose           # DBマイグレーション
            
            # 開発ツール
            air             # ホットリロード
            delve           # Goデバッガー
            
            # テスト・ベンチマーク
            go-junit-report # JUnit形式のテストレポート
            
            # その他便利ツール
            jq              # JSONパーサー
            curl            # API テスト
          ];

          shellHook = ''
            echo "🚀 LLM時間管理ツール - Go Server Development Environment"
            echo ""
            echo "Go環境:"
            echo "  - Go: $(go version)"
            echo "  - gopls: Go language server installed"
            echo ""
            echo "開発ツール:"
            echo "  - gofumpt: コードフォーマッタ"
            echo "  - golangci-lint: $(golangci-lint version --format short 2>/dev/null || echo 'installed')"
            echo "  - air: ホットリロード"
            echo "  - delve: デバッガー"
            echo ""
            echo "データベース:"
            echo "  - sqlite: $(sqlite3 --version | cut -d' ' -f1)"
            echo "  - sqlc: SQLからGoコード生成"
            echo "  - goose: DBマイグレーション"
            echo ""
            echo "使用可能なコマンド:"
            echo "  - 'go mod init' でモジュールを初期化"
            echo "  - 'air' でホットリロード開発サーバーを起動"
            echo "  - 'golangci-lint run' でリンターを実行"
            echo "  - 'go test ./...' でテストを実行"
            echo ""
            
            # GOPATHの設定（オプション）
            export GOPATH="$HOME/go"
            export PATH="$GOPATH/bin:$PATH"
          '';
        };
      }
    );
}

