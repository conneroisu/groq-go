#!/bin/bash
# Name: makefile.docs.sh
# Url: 
# 
# Description: A script to generate the go docs for the project.
# 
# Usage: make docs

mkdir docs
golds -s -gen -wdpkgs-listing=promoted -dir=./docs -footer=verbose+qrcode
xdg-open ./docs/index.html
