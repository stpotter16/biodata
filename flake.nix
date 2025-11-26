{
  description = "Dev environment for biodata";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";

    # 1.25.2 release
    go-nixpkgs.url = "github:NixOS/nixpkgs/01b6809f7f9d1183a2b3e081f0a1e6f8f415cb09";

    # 3.50.4 release
    sqlite-nixpkgs.url="github:NixOS/nixpkgs/6f374686605df381de8541c072038472a5ea2e2d";
  };

  outputs = {
    self,
    flake-utils,
    go-nixpkgs,
    sqlite-nixpkgs,
  } @ inputs:
    flake-utils.lib.eachDefaultSystem (system: let
      gopkg = go-nixpkgs.legacyPackages.${system};
      go = gopkg.go_1_25;
      sqlite = sqlite-nixpkgs.legacyPackages.${system}.sqlite;
    in {
      devShells.default = gopkg.mkShell {
        packages = [
          gopkg.gotools
          gopkg.gopls
          gopkg.go-outline
          gopkg.gopkgs
          gopkg.gocode-gomod
          gopkg.godef
          gopkg.golint
          go
          sqlite
        ];

        shellHook = ''
          PROJECT_NAME="$(basename "$PWD")"
          export GOPATH="$HOME/.local/share/go-workspaces/$PROJECT_NAME"
          export GOROOT="${go}/share/go"

          go version
          echo "sqlite" "$(sqlite3 --version | cut -d ' ' -f 1-2)"
        '';
      };

      formatter = gopkg.alejandra;
    });
}
