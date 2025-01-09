{
  pkgs,
  lib,
  config,
  inputs,
  ...
}: let
  unstable-pkgs = import inputs.unstable-nixpkgs {
    inherit (pkgs.stdenv) system;
  };
in {
  languages = {
    nix.enable = true;
    go = {
      enable = true;
      package = unstable-pkgs.go;
    };
  };

  git-hooks = {
    hooks = {
      golangci-lint.enable = true;
    };
  };
  # https://devenv.sh/packages/
  packages = with pkgs; [
    git
    zsh
    unstable-pkgs.iferr
    go
    pprof
    gopls
    impl
    golangci-lint-langserver
    golangci-lint
    revive
    templ
    gomodifytags
    gotests
    gotools
  ];

  scripts = {
    generate.exec = ''
      go generate -v ./...
    '';

    tests.exec = ''
      go test -v -short ./...
    '';
    unit-tests.exec = ''
      go test -v ./...
    '';
    lint.exec = ''
      golangci-lint run
    '';
    dx.exec = ''
      $EDITOR $(git rev-parse --show-toplevel)/devenv.nix
    '';
  };

  enterShell = ''
    git status
  '';
  enterTest = ''
    echo "Running tests"
    git --version | grep --color=auto "${pkgs.git.version}"
  '';

  cachix.enable = true;
}
