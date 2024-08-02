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

# Check if the command, protoc, exists, if not install it
gum spin --spinner dot --title "Installing Protoc" --show-output -- \
    command -v protoc >/dev/null 2>&1 || go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 &&  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# Check if the command, sqlite3, exists, if not install it
gum spin --spinner dot --title "Installing SQLite3" --show-output -- \
    command -v sqlite3 >/dev/null 2>&1 || brew install sqlite

# Check if the command, sqldiff, exists, if not install it
gum spin --spinner dot --title "Installing SQLDiff" --show-output -- \
    command -v sqldiff >/dev/null 2>&1 || brew install sqldiff

# Check if the command, sqlc, exists, if not install it
gum spin --spinner dot --title "Installing SQLC" --show-output -- \
    command -v sqlc >/dev/null 2>&1 || brew install sqlc
