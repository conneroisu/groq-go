#!/bin/bash
# Name: makefile/docs.sh
# Description: A script to generate the go docs for the project.
# 
# Usage: make docs

gum spin --spinner dot --title "Generating Docs" --show-output -- \
    gomarkdoc -o README.md -e .
