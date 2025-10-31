{
  description = "LLM時間管理ツール - Root Development Environment";

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
            # バージョン管理
            git
            gh
            
            # コードフォーマッター
            nodePackages.prettier
          ];

          shellHook = ''
            echo "🚀 LLM時間管理ツール - Root Development Environment"
            echo ""
            echo "Available tools:"
            echo "  - git: $(git --version | head -n 1)"
            echo "  - gh: $(gh --version | head -n 1)"
            echo "  - prettier: $(prettier --version)"
            echo ""
            echo "Prettier適用範囲: markdown, yaml, ymlなど"
            echo ""
          '';
        };
      }
    );
}

