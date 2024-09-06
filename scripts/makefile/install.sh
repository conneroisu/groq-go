#!/bin/bash
# file: makefile.install.sh
# title: Installing Development Requirements
# description: This script installs the required development tools for the project.

# Check if the command, brew, exists, if not install it
command -v brew >/dev/null 2>&1 || /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Check if the command, go, exists, if not install it
command -v go >/dev/null 2>&1 || brew install go

# Check if the command, gum, exists, if not install it
command -v gum >/dev/null 2>&1 || go install github.com/charmbracelet/gum@latest
