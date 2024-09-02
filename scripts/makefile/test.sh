#!/bin/bash
# file: makefile.test.sh
# title: Test Script
# description: This script runs the test for the project.
#
# usage: make test

go test -race -v -timeout 30s ./...

go test -coverprofile=coverage.out ./...