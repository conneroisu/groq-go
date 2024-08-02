#!/bin/bash
# file: taskfile/coverage.sh
# title: Running GoCovSh
# description: This script runs gocovsh to generate a coverage report.
#
# Usage: make coverage

task install

task test

# if gocovsh is executable
if [ -x "$(command -v gocovsh)" ]; then
    # if gocovsh is not empty
    if [ -s coverage.out ]; then
        # run gocovsh
        gocovsh
    else
        # if coverage.out is empty/not found
        echo "No coverage.out file found."
    fi
else
    # if gocovsh is not executable
    echo "gocovsh is not executable."
fi
