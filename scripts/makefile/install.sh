#!/bin/bash
# file: makefile.install.sh
# title: Installing Development Requirements
# description: This script installs the required development tools for the project.

# Check if the command, gum, exists, if not install it
command -v gum >/dev/null 2>&1 || go install github.com/charmbracelet/gum@latest

# Check if the command, revive, exists, if not install it
command -v revive >/dev/null 2>&1 || go install github.com/mgechev/revive@latest

# Check if the command, golangci-lint, exists, if not install it
command -v golangci-lint >/dev/null 2>&1 || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Check if the command, staticcheck, exists, if not install it
command -v staticcheck >/dev/null 2>&1 || go install honnef.co/go/tools/cmd/staticcheck@latest

# Check if the command, gocovsh, exists, if not install it
command -v gocovsh >/dev/null 2>&1 || go install github.com/boumenot/gocovsh@latest

# Check if the command, protoc-gen-go, exists, if not install it
command -v protoc-gen-go >/dev/null 2>&1 || go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0

# Check if the command, buf, exists, if not install it
command -v buf >/dev/null 2>&1 || go install github.com/bufbuild/buf/cmd/buf@v1.21.0
