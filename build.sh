#!/bin/bash

# Build release binaries for multiple platforms
# Usage: ./build.sh

set -e

# Ensure dist directory exists
mkdir -p dist

# Derive version from git (tag or short SHA), overridable via $VERSION
VERSION="${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo dev)}"
echo "📦 Version: $VERSION"

echo ""
echo "🚀 Building release binaries for all major platforms..."
echo ""

PLATFORMS=(
  "linux/amd64"
  "macos/amd64"
  "windows/amd64"
  "linux/arm64"
  "macos/arm64"
)

for PLATFORM in "${PLATFORMS[@]}"; do
  IFS="/" read -r OS ARCH <<< "$PLATFORM"
  EXT=""
  GOOS="$OS"
  if [ "$OS" == "macos" ]; then
    GOOS="darwin"
  fi
  if [ "$OS" == "windows" ]; then
    EXT=".exe"
  fi
  OUT="dist/adb-keep-screen-on-${OS}-${ARCH}${EXT}"
  echo "🔨 Building for $OS/$ARCH -> $OUT"
  if ! env CGO_ENABLED=0 GOOS=$GOOS GOARCH=$ARCH go build -trimpath -ldflags "-s -w -X main.version=$VERSION" -o "$OUT"; then
    echo "❌ Build failed for $OS/$ARCH"
    exit 1
  fi
  echo "✅ Build succeeded for $OS/$ARCH"
  echo ""
done

echo "📦 All binaries are in the dist directory."
echo ""