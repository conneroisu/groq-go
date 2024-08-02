#!/bin/bash
# file: taskfile.clean.sh
# title: Cleaning Script
# description: This script cleans the project

task install

# if there is a tmp folder, delete it
if [ -d "tmp" ]; then
    rm -rf tmp
fi

# if there is a bin folder, delete it
if [ -d "bin" ]; then
    rm -rf bin
fi

# if there is a node_modules folder, delete it
if [ -d "node_modules" ]; then
    rm -rf node_modules
fi

# if there is a node_modules in a subfolder, delete it
if [ -d "data/javascript/node_modules" ]; then
    rm -rf data/javascript/node_modules
fi

# if there is a coverage.out file, delete it
if [ -f "coverage.out" ]; then
    rm -rf coverage.out
fi
