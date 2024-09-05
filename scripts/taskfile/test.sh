#!/bin/bash
# file: taskfile.test.sh
# title: Test Script
# description: This script runs the test for the project.
#
# usage: make test

go clean -testcache
gum spin --spinner dot --title "Running Tests" --show-output -- \
    go test -race -timeout 60s ./...

gum spin --spinner dot --title "Generating Coverage" --show-output -- \
    go test -coverprofile=coverage.out ./...

gocovsh coverage.out
