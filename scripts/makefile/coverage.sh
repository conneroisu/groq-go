#!/bin/bash
# file: makefile.coverage.sh
# title: Test Script
# description: This script runs the test for the project.
#
# usage: make test

go test -race -timeout 30s ./...

go test -coverprofile=coverage.out ./...
