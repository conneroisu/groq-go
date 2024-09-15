#!/bin/bash
# file: makefile.fmt.sh
# title: Formatting Go Files
# description: This script formats the Go files using gofmt and golines.
#
# Usage: make fmt

gofmt -w .

# golines -w --max-len=79 .
