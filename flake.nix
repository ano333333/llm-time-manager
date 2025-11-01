{
  description = "LLM時間管理ツール - Unified Development Environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachSystem [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ] (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = import ./nix/packages.nix { inherit pkgs; };
          shellHook = import ./nix/shell-hook.nix;
        };

        apps.fmt = {
          type = "app";
          program = toString (
            pkgs.writeShellScript "fmt" ''
              # 現在のディレクトリを保存
              ROOT_DIR="$PWD"
              
              echo "🎨 Formatting web..."
              cd "$ROOT_DIR/web"
              ${pkgs.biome}/bin/biome check --write . || true
              ${pkgs.nodePackages.prettier}/bin/prettier --write '**/*.{md,yaml,yml}' || true
              
              echo ""
              echo "🎨 Formatting server..."
              cd "$ROOT_DIR/server"
              ${pkgs.gofumpt}/bin/gofumpt -l -w . || true
              ${pkgs.gotools}/bin/goimports -w . || true
              
              echo ""
              echo "✅ Formatting completed!"
            ''
          );
        };
      }
    );
}
