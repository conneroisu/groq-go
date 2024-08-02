#!/bin/bash
# file: makefile.tidy.sh
# title: Running Go Mod Tidy
# description: This script runs go mod tidy to clean up the go.mod and go.sum files.
#
# Usage: make tidy

gum spin --spinner dot --title "Running Go Mod Tidy" --show-output -- \
    go mod tidy
