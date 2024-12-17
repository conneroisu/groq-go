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

  languages.go.enable = true;
  languages.nix.enable = true;
  # https://devenv.sh/packages/
  packages = with pkgs; [
    git
    zsh
    revive
    unstable-pkgs.iferr
    go
    gopls
    impl
    golangci-lint-langserver
    golangci-lint
    templ
    gomodifytags
    gotests
    gotools
  ];

  # https://devenv.sh/languages/
  # languages.rust.enable = true;

  # https://devenv.sh/processes/
  # processes.cargo-watch.exec = "cargo-watch";

  # https://devenv.sh/services/
  # services.postgres.enable = true;

  # https://devenv.sh/scripts/
  scripts.generate.exec = ''
    go generate -v ./...
  '';

  enterShell = ''
    git status
  '';

  # https://devenv.sh/tasks/
  # tasks = {
  #   "myproj:setup".exec = "mytool build";
  #   "devenv:enterShell".after = [ "myproj:setup" ];
  # };

  # https://devenv.sh/tests/
  enterTest = ''
    echo "Running tests"
    git --version | grep --color=auto "${pkgs.git.version}"
  '';

  # https://devenv.sh/pre-commit-hooks/
  # pre-commit.hooks.shellcheck.enable = true;

  # See full reference at https://devenv.sh/reference/options/
  cachix.enable = true;
}
