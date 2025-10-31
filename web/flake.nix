{
  description = "LLM時間管理ツール - WebView Development Environment";

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
            # Node.js環境
            nodejs_22
            
            # パッケージマネージャ
            nodePackages.pnpm
            
            # バージョン管理
            git
            gh
            
            # コードフォーマッター/リンター
            biome
            nodePackages.prettier
          ];

          shellHook = ''
            echo "🚀 LLM時間管理ツール - WebView Development Environment"
            echo ""
            echo "Available tools:"
            echo "  - Node.js: $(node --version)"
            echo "  - pnpm: $(pnpm --version)"
            echo "  - git: $(git --version | head -n 1)"
            echo "  - gh: $(gh --version | head -n 1)"
            echo "  - biome: $(biome --version)"
            echo "  - prettier: $(prettier --version)"
            echo ""
            echo "Biome適用範囲: .ts, .tsx, .json, .html, .css"
            echo "Prettier適用範囲: その他のファイル"
            echo ""
            
            # pnpm設定
            export PNPM_HOME="$HOME/.local/share/pnpm"
            export PATH="$PNPM_HOME:$PATH"
          '';
        };
      }
    );
}

