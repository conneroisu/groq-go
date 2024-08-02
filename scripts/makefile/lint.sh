#!/bin/bash
# file: makefile.js.sh
# title: Running Webpack
# description: This script runs Webpack to build the JavaScript files.
#
# Usage: make js

staticcheck ./...

golangci-lint run

go vet ./...

revive -config .revive.toml ./...
