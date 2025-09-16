#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Ensure dist directory exists
mkdir -p dist

# Build the Go project and output the binary into dist/adb-keep-screen-on
if ! go build -o dist/adb-keep-screen-on; then
  echo "Go build failed."
  exit 1
fi

echo "Build succeeded. Binary is in dist directory."