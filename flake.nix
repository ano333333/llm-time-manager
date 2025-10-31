{
  description = "LLM時間管理ツール - Unified Development Environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachSystem [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ] (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            # バージョン管理
            git
            gh
            
            # Go開発環境（server用）
            go
            gopls           # Go language server
            gotools         # goimports, godoc等のツール
            go-tools        # staticcheck等の追加ツール
            
            # Goフォーマッタ
            gofumpt         # gofmtの改善版
            golines         # 長い行を分割
            
            # Goリンター
            golangci-lint
            
            # データベース
            sqlite
            sqlc            # SQLからGoコード生成
            goose           # DBマイグレーション
            
            # Go開発ツール
            air             # ホットリロード
            delve           # Goデバッガー
            
            # Goテスト・ベンチマーク
            go-junit-report # JUnit形式のテストレポート
            
            # Node.js環境（web用）
            nodejs_22
            
            # パッケージマネージャ
            nodePackages.pnpm
            
            # コードフォーマッター/リンター
            biome
            nodePackages.prettier
            
            # Nix開発環境
            nil             # Nix language server
            nixpkgs-fmt     # Nixフォーマッター
            
            # その他便利ツール
            jq              # JSONパーサー
            curl            # API テスト
          ];

          shellHook = ''
            echo "🚀 LLM時間管理ツール - Unified Development Environment"
            echo ""
            echo "Go環境（server用）:"
            echo "  - Go: $(go version)"
            echo "  - gopls: Go language server"
            echo "  - gofumpt: コードフォーマッタ"
            echo "  - golangci-lint: $(golangci-lint version --format short 2>/dev/null || echo 'installed')"
            echo "  - air: ホットリロード"
            echo "  - delve: デバッガー"
            echo ""
            echo "Node.js環境（web用）:"
            echo "  - Node.js: $(node --version)"
            echo "  - pnpm: $(pnpm --version)"
            echo "  - biome: $(biome --version)"
            echo ""
            echo "データベース:"
            echo "  - sqlite: $(sqlite3 --version | cut -d' ' -f1)"
            echo "  - sqlc: SQLからGoコード生成"
            echo "  - goose: DBマイグレーション"
            echo ""
            echo "Nix環境:"
            echo "  - nil: Nix language server"
            echo "  - nixpkgs-fmt: Nixフォーマッター"
            echo ""
            echo "バージョン管理:"
            echo "  - git: $(git --version | head -n 1)"
            echo "  - gh: $(gh --version | head -n 1)"
            echo ""
            echo "その他ツール:"
            echo "  - prettier: $(prettier --version)"
            echo "  - jq: $(jq --version)"
            echo ""
            echo "サポートプラットフォーム: Linux (x86_64/aarch64), macOS (x86_64/aarch64)"
            echo ""
            
            # GOPATHの設定
            export GOPATH="$HOME/go"
            export PATH="$GOPATH/bin:$PATH"
            
            # pnpm設定
            export PNPM_HOME="$HOME/.local/share/pnpm"
            export PATH="$PNPM_HOME:$PATH"
          '';
        };
      }
    );
}

