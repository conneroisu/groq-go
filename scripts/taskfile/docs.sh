#!/bin/bash
# Name: makefile/docs.sh
# Description: A script to generate the go docs for the project.
# 
# Usage: make docs

gum spin --spinner dot --title "Making Docs Folder" --show-output -- \
    mkdir docs

gum spin --spinner dot --title "Generating Docs" --show-output -- \
    golds -s -gen -wdpkgs-listing=promoted -dir=./docs -footer=verbose+qrcode

gum spin --spinner dot --title "Opening Docs Folder" --show-output -- \
    xdg-open ./docs/index.html
