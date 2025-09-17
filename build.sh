#!/bin/bash

# Always build a release binary with Go optimization flags
# Usage: ./build.sh

# Exit immediately if a command exits with a non-zero status
set -e

# Ensure dist directory exists
mkdir -p dist

echo ""
echo "🚀 Building release binary..."
echo ""

if ! go build -ldflags "-s -w" -o dist/adb-keep-screen-on; then
  echo "\n❌ Release build failed."
  exit 1
fi

echo "✅ Release build succeeded."
echo ""
echo "📦 Binary is in dist directory."
echo ""