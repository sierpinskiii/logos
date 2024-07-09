let
  sources = import ./nix/sources.nix;
  pkgs = import sources.nixpkgs {};
in pkgs.mkShell {
  nativeBuildInputs = [
    pkgs.niv
    pkgs.zlib
    pkgs.pkg-config 
    pkgs.sqlite
    pkgs.go
    pkgs.gopls
    pkgs.delve
    pkgs.go-tools
    pkgs.imagemagick
    pkgs.ghostscript
  ];
}
