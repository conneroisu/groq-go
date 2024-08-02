#!/bin/bash
# file: taskfile.test.sh
# title: Test Script
# description: This script runs the test for the project.
#
# usage: make test

gum spin --spinner dot --title "Running Lint Checks" --show-output -- \
    gum spin --spinner dot --title "Running Staticcheck" --show-output -- \
        staticcheck ./...
    gum spin --spinner dot --title "Running Go Lint CI" --show-output -- \
        golangci-lint run
    gum spin --spinner dot --title "Running Revive" --show-output -- \
        revive -config .revive.toml ./...
    gum spin --spinner dot --title "Running go vet" --show-output -- \
        go vet ./...
#
#
# go vet ./...
#
# revive -config .revive.toml ./...
