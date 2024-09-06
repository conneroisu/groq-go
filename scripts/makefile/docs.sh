#!/bin/bash
# Name: makefile.docs.sh
# Url: 
# 
# Description: A script to generate the go docs for the project.
# 
# Usage: make docs

gomarkdoc -o README.md -e .
