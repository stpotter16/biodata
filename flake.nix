{
  description = "Dev environment for biodata";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";

    # 1.25.2 release
    go-nixpkgs.url = "github:NixOS/nixpkgs/01b6809f7f9d1183a2b3e081f0a1e6f8f415cb09";
  };

  outputs = {
    self,
    flake-utils,
    go-nixpkgs,
  } @ inputs:
    flake-utils.lib.eachDefaultSystem (system: let
      gopkg = go-nixpkgs.legacyPackages.${system};
      go = gopkg.go_1_23;
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
        ];

        shellHook = ''
          export GOROOT="${go}/share/go"

          go version
        '';
      };

      formatter = gopkg.alejandra;
    });
}
